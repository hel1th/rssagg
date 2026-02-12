package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/domain"
)

type CreateFeedFollowRequest struct {
	FeedID uuid.UUID `json:"feed_id"`
}

type FeedFollowResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	UserID    uuid.UUID `json:"user_id"`
	FeedID    uuid.UUID `json:"feed_id"`
}

func FeedFollowToResponse(ff *domain.FeedFollow) FeedFollowResponse {
	return FeedFollowResponse{
		ID:        ff.ID,
		CreatedAt: ff.CreatedAt,
		UpdatedAt: ff.UpdatedAt,
		UserID:    ff.UserID,
		FeedID:    ff.FeedID,
	}
}

func FeedFollowsToResponse(ffs []*domain.FeedFollow) []FeedFollowResponse {
	responses := make([]FeedFollowResponse, len(ffs))
	for i, ff := range ffs {
		responses[i] = FeedFollowToResponse(ff)
	}
	return responses
}
