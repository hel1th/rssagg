package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
)

type FeedFollowService interface {
	FollowFeed(ctx context.Context, userID, feedID uuid.UUID) (*domain.FeedFollow, error)
	GetUserFeedFollows(ctx context.Context, userID uuid.UUID) ([]*domain.FeedFollow, error)
	UnfollowFeed(ctx context.Context, feedFollowID, userID uuid.UUID) error
}

type feedFollowService struct {
	repo repository.FeedFollowRepository
}

func NewFeedFollowService(repo repository.FeedFollowRepository) FeedFollowService {
	return &feedFollowService{
		repo: repo,
	}
}

func (s *feedFollowService) FollowFeed(ctx context.Context, userID, feedID uuid.UUID) (*domain.FeedFollow, error) {
	feedFollow := domain.NewFeedFollow(userID, feedID)
	
	if err := feedFollow.Validate(); err != nil {
		return nil, err
	}
	
	dbFeedFollow, err := s.repo.Create(ctx, database.CreateFeedFollowParams{
		ID:        feedFollow.ID,
		CreatedAt: feedFollow.CreatedAt,
		UpdatedAt: feedFollow.UpdatedAt,
		UserID:    feedFollow.UserID,
		FeedID:    feedFollow.FeedID,
	})
	if err != nil {
		if err.Error() == "duplicate key value violates unique constraint" {
			return nil, domain.ErrDuplicateFeedFollow
		}
		return nil, err
	}
	
	return domain.MapFeedFollowFromDB(dbFeedFollow), nil
}

func (s *feedFollowService) GetUserFeedFollows(ctx context.Context, userID uuid.UUID) ([]*domain.FeedFollow, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}
	
	dbFeedFollows, err := s.repo.GetByUser(ctx, userID)
	if err != nil {
		return nil, err
	}
	
	return domain.MapFeedFollowsFromDB(dbFeedFollows), nil
}

func (s *feedFollowService) UnfollowFeed(ctx context.Context, feedFollowID, userID uuid.UUID) error {
	if feedFollowID == uuid.Nil {
		return domain.ErrFeedFollowNotFound
	}
	if userID == uuid.Nil {
		return domain.ErrInvalidUserID
	}
	
	err := s.repo.Delete(ctx, database.DeleteFeedFollowParams{
		ID:     feedFollowID,
		UserID: userID,
	})
	if err != nil {
		return domain.ErrCannotUnfollowFeed
	}
	
	return nil
}