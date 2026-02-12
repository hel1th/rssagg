package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
)

type FeedService interface {
	CreateFeed(ctx context.Context, name, url string, userID uuid.UUID) (*domain.Feed, error)
	GetAllFeeds(ctx context.Context) ([]*domain.Feed, error)
	GetFeedByID(ctx context.Context, id uuid.UUID) (*domain.Feed, error)
	GetNextFeedsToFetch(ctx context.Context, limit int) ([]*domain.Feed, error)
	MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (*domain.Feed, error)
}

type feedService struct {
	repo repository.FeedRepository
}

func NewFeedService(repo repository.FeedRepository) FeedService {
	return &feedService{
		repo: repo,
	}
}

func (s *feedService) CreateFeed(ctx context.Context, name, url string, userID uuid.UUID) (*domain.Feed, error) {
	feed := domain.NewFeed(name, url, userID)

	if err := feed.Validate(); err != nil {
		return nil, err
	}

	dbFeed, err := s.repo.Create(ctx, database.CreateFeedParams{
		ID:        feed.ID,
		CreatedAt: feed.CreatedAt,
		UpdatedAt: feed.UpdatedAt,
		Name:      feed.Name,
		Url:       feed.URL,
		UserID:    feed.UserID,
	})
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			return nil, domain.ErrDuplicateFeed
		}
		return nil, err
	}

	return domain.MapFeedFromDB(dbFeed), nil
}

func (s *feedService) GetAllFeeds(ctx context.Context) ([]*domain.Feed, error) {
	dbFeeds, err := s.repo.GetAll(ctx)
	if err != nil {
		return nil, err
	}

	return domain.MapFeedsFromDB(dbFeeds), nil
}

func (s *feedService) GetFeedByID(ctx context.Context, id uuid.UUID) (*domain.Feed, error) {
	dbFeed, err := s.repo.GetByID(ctx, id)
	if err != nil {
		return nil, domain.ErrFeedNotFound
	}

	return domain.MapFeedFromDB(dbFeed), nil
}

func (s *feedService) GetNextFeedsToFetch(ctx context.Context, limit int) ([]*domain.Feed, error) {
	if limit <= 0 {
		limit = 10 
	}

	dbFeeds, err := s.repo.GetNextToFetch(ctx, int32(limit))
	if err != nil {
		return nil, err
	}

	return domain.MapFeedsFromDB(dbFeeds), nil
}

func (s *feedService) MarkFeedAsFetched(ctx context.Context, id uuid.UUID) (*domain.Feed, error) {
	dbFeed, err := s.repo.MarkAsFetched(ctx, id)
	if err != nil {
		return nil, domain.ErrFeedNotFound
	}

	return domain.MapFeedFromDB(dbFeed), nil
}
