package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/domain"
)

type CreateFeedRequest struct {
	Name string `json:"name"`
	URL  string `json:"url"`
}

type FeedResponse struct {
	ID            uuid.UUID  `json:"id"`
	CreatedAt     time.Time  `json:"created_at"`
	UpdatedAt     time.Time  `json:"updated_at"`
	Name          string     `json:"name"`
	URL           string     `json:"url"`
	UserID        uuid.UUID  `json:"user_id"`
	LastFetchedAt *time.Time `json:"last_fetched_at,omitempty"`
}

func FeedToResponse(feed *domain.Feed) FeedResponse {
	return FeedResponse{
		ID:            feed.ID,
		CreatedAt:     feed.CreatedAt,
		UpdatedAt:     feed.UpdatedAt,
		Name:          feed.Name,
		URL:           feed.URL,
		UserID:        feed.UserID,
		LastFetchedAt: feed.LastFetchedAt,
	}
}

func FeedsToResponse(feeds []*domain.Feed) []FeedResponse {
	responses := make([]FeedResponse, len(feeds))
	for i, feed := range feeds {
		responses[i] = FeedToResponse(feed)
	}
	return responses
}
