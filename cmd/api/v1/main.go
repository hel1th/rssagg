package main

import (
	"context"
	"database/sql"
	"log"
	"net/http"
	"os"
	"time"

	"github.com/go-chi/chi/v5"
	"github.com/go-chi/cors"
	"github.com/joho/godotenv"
	_ "github.com/lib/pq"

	"github.com/hel1th/rssagg/api/v1/handlers"
	"github.com/hel1th/rssagg/api/v1/middleware"
	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
	"github.com/hel1th/rssagg/internal/service"
)

func main() {
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("No .env file found, using system environment variables")
	}

	port := os.Getenv("PORT")
	if port == "" {
		log.Fatal("PORT environment variable is not set")
	}

	dbURL := os.Getenv("DB_URL")
	if dbURL == "" {
		log.Fatal("DB_URL environment variable is not set")
	}

	conn, err := sql.Open("postgres", dbURL)
	if err != nil {
		log.Fatal("Failed to connect to database:", err)
	}
	defer conn.Close()

	if err := conn.Ping(); err != nil {
		log.Fatal("Failed to ping database:", err)
	}
	log.Println("Database connection established")

	db := database.New(conn)

	// Initialize repositories
	userRepo := repository.NewUserRepository(db)
	feedRepo := repository.NewFeedRepository(db)
	feedFollowRepo := repository.NewFeedFollowRepository(db)
	postRepo := repository.NewPostRepository(db)

	// Initialize services
	userService := service.NewUserService(userRepo)
	feedService := service.NewFeedService(feedRepo)
	feedFollowService := service.NewFeedFollowService(feedFollowRepo)
	postService := service.NewPostService(postRepo)
	rssService := service.NewRSSService(postRepo, feedRepo)

	// Initialize handlers
	userHandler := handlers.NewUserHandler(userService)
	feedHandler := handlers.NewFeedHandler(feedService)
	feedFollowHandler := handlers.NewFeedFollowHandler(feedFollowService)
	postHandler := handlers.NewPostHandler(postService)
	rssHandler := handlers.NewRSSHandler(rssService, feedService)

	// Initialize middleware
	authMiddleware := middleware.NewAuthMiddleware(userService)

	// Start background RSS scraper
	go startScraper(db, feedService, rssService, 10, time.Minute)

	router := setupRouter(
		userHandler,
		feedHandler,
		feedFollowHandler,
		postHandler,
		rssHandler,
		authMiddleware,
	)

	srv := &http.Server{
		Handler:      router,
		Addr:         ":" + port,
		ReadTimeout:  15 * time.Second,
		WriteTimeout: 15 * time.Second,
		IdleTimeout:  60 * time.Second,
	}

	log.Printf("Server starting on port %s", port)
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal("Server failed to start:", err)
	}
}

func setupRouter(
	userHandler *handlers.UserHandler,
	feedHandler *handlers.FeedHandler,
	feedFollowHandler *handlers.FeedFollowHandler,
	postHandler *handlers.PostHandler,
	rssHandler *handlers.RSSHandler,
	authMiddleware *middleware.AuthMiddleware,
) http.Handler {
	router := chi.NewRouter()

	router.Use(cors.Handler(cors.Options{
		AllowedOrigins:   []string{"https://*", "http://*"},
		AllowedMethods:   []string{"GET", "POST", "PUT", "DELETE", "OPTIONS"},
		AllowedHeaders:   []string{"*"},
		ExposedHeaders:   []string{"Link"},
		AllowCredentials: false,
		MaxAge:           300,
	}))

	v1Router := chi.NewRouter()

	v1Router.Get("/healthz", healthCheck)
	v1Router.Get("/err", errorHandler)

	v1Router.Post("/users", userHandler.CreateUser)
	v1Router.With(authMiddleware.Require).Get("/users", adaptAuthHandler(userHandler.GetUser))

	v1Router.With(authMiddleware.Require).Post("/feeds", adaptAuthHandler(feedHandler.CreateFeed))
	v1Router.Get("/feeds", feedHandler.GetAllFeeds)

	v1Router.With(authMiddleware.Require).Post("/feed_follows", adaptAuthHandler(feedFollowHandler.FollowFeed))
	v1Router.With(authMiddleware.Require).Get("/feed_follows", adaptAuthHandler(feedFollowHandler.GetUserFeedFollows))
	v1Router.With(authMiddleware.Require).Delete("/feed_follows", adaptAuthHandler(feedFollowHandler.UnfollowFeed))

	v1Router.With(authMiddleware.Require).Get("/posts", adaptAuthHandler(postHandler.GetPostsForUser))

	v1Router.With(authMiddleware.Require).Post("/rss/fetch", adaptAuthHandler(rssHandler.FetchFeed))

	router.Mount("/v1", v1Router)

	return router
}

func adaptAuthHandler(handler func(http.ResponseWriter, *http.Request, *domain.User)) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		user, ok := middleware.GetUserFromContext(r.Context())
		if !ok {
			http.Error(w, "Unauthorized", http.StatusUnauthorized)
			return
		}
		handler(w, r, user)
	}
}

func healthCheck(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusOK)
	w.Write([]byte(`{"status":"ok"}`))
}

func errorHandler(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusInternalServerError)
	w.Write([]byte(`{"error":"Internal Server Error"}`))
}

func startScraper(
	db *database.Queries,
	feedService service.FeedService,
	rssService service.RSSService,
	concurrency int,
	interval time.Duration,
) {
	log.Printf("Starting RSS scraper: interval=%v, concurrency=%d", interval, concurrency)
	ticker := time.NewTicker(interval)
	defer ticker.Stop()

	for range ticker.C {
		ctx := context.Background()
		feeds, err := feedService.GetNextFeedsToFetch(ctx, concurrency)
		if err != nil {
			log.Printf("Error fetching feeds to scrape: %v", err)
			continue
		}

		if len(feeds) == 0 {
			continue
		}

		log.Printf("Scraping %d feeds", len(feeds))

		feedValues := make([]domain.Feed, len(feeds))
		for i, feed := range feeds {
			feedValues[i] = *feed
		}

		if err := rssService.FetchAndStoreFeeds(ctx, feedValues); err != nil {
			log.Printf("Error during feed scraping: %v", err)
		}
	}
}
