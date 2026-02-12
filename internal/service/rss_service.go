package service

import (
	"context"
	gosql "database/sql"
	"fmt"
	"log"
	"strings"
	"sync"
	"time"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
	"github.com/hel1th/rssagg/internal/rss"
)

type RSSService interface {
	FetchAndStoreFeeds(ctx context.Context, feeds []domain.Feed) error
	FetchSingleFeed(ctx context.Context, feed domain.Feed) (int, error)
}

type rssService struct {
	postRepo repository.PostRepository
	feedRepo repository.FeedRepository
	fetcher  rss.Fetcher
}

func NewRSSService(postRepo repository.PostRepository, feedRepo repository.FeedRepository) RSSService {
	return NewRSSServiceWithFetcher(postRepo, feedRepo, rss.NewFetcher())
}

func NewRSSServiceWithFetcher(postRepo repository.PostRepository, feedRepo repository.FeedRepository, fetcher rss.Fetcher) RSSService {
	return &rssService{
		postRepo: postRepo,
		feedRepo: feedRepo,
		fetcher:  fetcher,
	}
}

func (s *rssService) FetchAndStoreFeeds(ctx context.Context, feeds []domain.Feed) error {
	var wg sync.WaitGroup

	for _, feed := range feeds {
		wg.Add(1)
		go func(f domain.Feed) {
			defer wg.Done()

			newPosts, err := s.FetchSingleFeed(ctx, f)
			if err != nil {
				log.Printf("Error fetching feed %s: %v", f.Name, err)
				return
			}

			log.Printf("Feed %s collected. %d new posts", f.Name, newPosts)
		}(feed)
	}

	wg.Wait()
	return nil
}

func (s *rssService) FetchSingleFeed(ctx context.Context, feed domain.Feed) (int, error) {
	_, err := s.feedRepo.MarkAsFetched(ctx, feed.ID)
	if err != nil {
		return 0, fmt.Errorf("failed to mark feed as fetched: %w", err)
	}

	if s.fetcher == nil {
		return 0, fmt.Errorf("no RSS fetcher configured")
	}

	rssFeed, err := s.fetcher.Fetch(feed.URL)
	if err != nil {
		return 0, fmt.Errorf("failed to fetch RSS from URL: %w", err)
	}

	newPostCount := 0
	for _, item := range rssFeed.Items {
		postData, err := s.parseRSSItem(item, feed.ID)
		if err != nil {
			log.Printf("Error parsing RSS item: %v", err)
			continue
		}

		err = s.postRepo.Create(ctx, postData)
		if err != nil {
			if s.isDuplicateError(err) {
				continue
			}
			log.Printf("Error creating post: %v", err)
			continue
		}

		newPostCount++
	}

	return newPostCount, nil
}

func (s *rssService) parseRSSItem(item domain.RSSItemData, feedID uuid.UUID) (database.CreatePostParams, error) {
	var pubAt time.Time
	var parseErr error
	layouts := []string{time.RFC1123Z, time.RFC1123, time.RFC822, time.RFC3339}
	for _, l := range layouts {
		pubAt, parseErr = time.Parse(l, item.PubDate)
		if parseErr == nil {
			break
		}
	}
	if parseErr != nil {
		return database.CreatePostParams{}, fmt.Errorf("failed to parse publish date %s: %w", item.PubDate, parseErr)
	}
	desc := gosql.NullString{}
	if item.Description != "" {
		desc = gosql.NullString{String: item.Description, Valid: true}
	}

	now := time.Now().UTC()
	return database.CreatePostParams{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       item.Title,
		Description: desc,
		Url:         item.Link,
		FeedID:      feedID,
		PublishedAt: pubAt,
	}, nil
}

func (s *rssService) isDuplicateError(err error) bool {
	return strings.Contains(err.Error(), "23505")
}
