package middleware

import (
	"context"
	"encoding/json"
	"net/http"

	"github.com/hel1th/rssagg/internal/auth"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/service"
)

type contextKey string

const userContextKey contextKey = "user"

type AuthMiddleware struct {
	userService service.UserService
}

func NewAuthMiddleware(userService service.UserService) *AuthMiddleware {
	return &AuthMiddleware{userService: userService}
}

func (m *AuthMiddleware) Require(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		apiKey, err := auth.GetAPIKey(r.Header)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Unauthorized")
			return
		}

		user, err := m.userService.GetUserByAPIKey(r.Context(), apiKey)
		if err != nil {
			respondWithError(w, http.StatusUnauthorized, "Invalid API key")
			return
		}

		ctx := context.WithValue(r.Context(), userContextKey, user)
		next.ServeHTTP(w, r.WithContext(ctx))
	})
}

func GetUserFromContext(ctx context.Context) (*domain.User, bool) {
	user, ok := ctx.Value(userContextKey).(*domain.User)
	return user, ok
}

func respondWithError(w http.ResponseWriter, code int, msg string) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	json.NewEncoder(w).Encode(map[string]string{"error": msg})
}
