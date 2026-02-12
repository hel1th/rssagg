package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
)

type FeedFollowRepository interface {
	Create(ctx context.Context, params database.CreateFeedFollowParams) (database.FeedFollow, error)
	GetByUser(ctx context.Context, userID uuid.UUID) ([]database.FeedFollow, error)
	Delete(ctx context.Context, params database.DeleteFeedFollowParams) error
}

type feedFollowRepository struct {
	db *database.Queries
}

func NewFeedFollowRepository(db *database.Queries) FeedFollowRepository {
	return &feedFollowRepository{
		db: db,
	}
}

func (r *feedFollowRepository) Create(ctx context.Context, params database.CreateFeedFollowParams) (database.FeedFollow, error) {
	return r.db.CreateFeedFollow(ctx, params)
}

func (r *feedFollowRepository) GetByUser(ctx context.Context, userID uuid.UUID) ([]database.FeedFollow, error) {
	return r.db.GetFeedFollows(ctx, userID)
}

func (r *feedFollowRepository) Delete(ctx context.Context, params database.DeleteFeedFollowParams) error {
	return r.db.DeleteFeedFollow(ctx, params)
}
