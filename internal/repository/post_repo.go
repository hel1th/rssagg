package repository

import (
	"context"

	"github.com/hel1th/rssagg/internal/database"
)

type PostRepository interface {
	Create(ctx context.Context, params database.CreatePostParams) error
	GetForUser(ctx context.Context, params database.GetPostsForUserParams) ([]database.Post, error)
}

type postRepository struct {
	db *database.Queries
}

func NewPostRepository(db *database.Queries) PostRepository {
	return &postRepository{
		db: db,
	}
}

func (r *postRepository) Create(ctx context.Context, params database.CreatePostParams) error {
	return r.db.CreatePost(ctx, params)
}

func (r *postRepository) GetForUser(ctx context.Context, params database.GetPostsForUserParams) ([]database.Post, error) {
	return r.db.GetPostsForUser(ctx, params)
}
