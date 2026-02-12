package service

import (
	"context"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
)

type PostService interface {
	GetPostsForUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Post, error)
}

type postService struct {
	repo repository.PostRepository
}

func NewPostService(repo repository.PostRepository) PostService {
	return &postService{repo: repo}
}

func (s *postService) GetPostsForUser(ctx context.Context, userID uuid.UUID, limit, offset int) ([]*domain.Post, error) {
	if userID == uuid.Nil {
		return nil, domain.ErrInvalidUserID
	}

	if limit <= 0 {
		limit = 10
	}
	if limit > 100 {
		limit = 100
	}
	if offset < 0 {
		offset = 0
	}

	dbPosts, err := s.repo.GetForUser(ctx, database.GetPostsForUserParams{
		UserID: userID,
		Limit:  int32(limit),
		Offset: int32(offset),
	})
	if err != nil {
		return nil, err
	}

	return domain.MapPostsFromDB(dbPosts), nil
}
