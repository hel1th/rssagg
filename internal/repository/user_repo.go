package repository

import (
	"context"

	"github.com/hel1th/rssagg/internal/database"
)

type UserRepository interface {
	Create(ctx context.Context, params database.CreateUserParams) (database.User, error)
	GetByAPIKey(ctx context.Context, apiKey string) (database.User, error)
}

type userRepository struct {
	db *database.Queries
}

func NewUserRepository(db *database.Queries) UserRepository {
	return &userRepository{
		db: db,
	}
}

func (r *userRepository) Create(ctx context.Context, params database.CreateUserParams) (database.User, error) {
	return r.db.CreateUser(ctx, params)
}

func (r *userRepository) GetByAPIKey(ctx context.Context, apiKey string) (database.User, error) {
	return r.db.GetUserByAPIKey(ctx, apiKey)
}
