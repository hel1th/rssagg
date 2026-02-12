package repository

import (
	"context"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
)

type FeedRepository interface {
	Create(ctx context.Context, params database.CreateFeedParams) (database.Feed, error)
	GetAll(ctx context.Context) ([]database.Feed, error)
	GetByID(ctx context.Context, id uuid.UUID) (database.Feed, error)
	GetNextToFetch(ctx context.Context, limit int32) ([]database.Feed, error)
	MarkAsFetched(ctx context.Context, id uuid.UUID) (database.Feed, error)
}

type feedRepository struct {
	db *database.Queries
}

func NewFeedRepository(db *database.Queries) FeedRepository {
	return &feedRepository{
		db: db,
	}
}

func (r *feedRepository) Create(ctx context.Context, params database.CreateFeedParams) (database.Feed, error) {
	return r.db.CreateFeed(ctx, params)
}

func (r *feedRepository) GetAll(ctx context.Context) ([]database.Feed, error) {
	return r.db.GetFeeds(ctx)
}

func (r *feedRepository) GetByID(ctx context.Context, id uuid.UUID) (database.Feed, error) {
	return r.db.GetFeedByID(ctx, id)
}

func (r *feedRepository) GetNextToFetch(ctx context.Context, limit int32) ([]database.Feed, error) {
	return r.db.GetNextFeedsToFetch(ctx, limit)
}

func (r *feedRepository) MarkAsFetched(ctx context.Context, id uuid.UUID) (database.Feed, error) {
	return r.db.MarkFeedAsFetched(ctx, id)
}
