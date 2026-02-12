package service

import (
	"context"

	"github.com/hel1th/rssagg/internal/database"
	"github.com/hel1th/rssagg/internal/domain"
	"github.com/hel1th/rssagg/internal/repository"
)

type UserService interface {
	CreateUser(ctx context.Context, name string) (*domain.User, error)
	GetUserByAPIKey(ctx context.Context, apiKey string) (*domain.User, error)
}

type userService struct {
	repo repository.UserRepository
}

func NewUserService(repo repository.UserRepository) UserService {
	return &userService{
		repo: repo,
	}
}

func (s *userService) CreateUser(ctx context.Context, name string) (*domain.User, error) {
	user := domain.NewUser(name)
	
	if err := user.Validate(); err != nil {
		return nil, err
	}
	
	dbUser, err := s.repo.Create(ctx, database.CreateUserParams{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
	})
	if err != nil {
		return nil, err
	}
	
	return domain.MapUserFromDB(dbUser), nil
}

func (s *userService) GetUserByAPIKey(ctx context.Context, apiKey string) (*domain.User, error) {
	if apiKey == "" {
		return nil, domain.ErrInvalidAPIKey
	}
	
	dbUser, err := s.repo.GetByAPIKey(ctx, apiKey)
	if err != nil {
		return nil, domain.ErrUserNotFound
	}
	
	return domain.MapUserFromDB(dbUser), nil
}