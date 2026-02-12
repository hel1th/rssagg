package dto

import (
	"time"

	"github.com/google/uuid"
	"github.com/hel1th/rssagg/internal/domain"
)

type CreateUserRequest struct {
	Name string `json:"name"`
}


type UserResponse struct {
	ID        uuid.UUID `json:"id"`
	CreatedAt time.Time `json:"created_at"`
	UpdatedAt time.Time `json:"updated_at"`
	Name      string    `json:"name"`
	APIKey    string    `json:"api_key"`
}

type PostResponse struct {
	ID          uuid.UUID `json:"id"`
	CreatedAt   time.Time `json:"created_at"`
	UpdatedAt   time.Time `json:"updated_at"`
	Title       string    `json:"title"`
	Description *string   `json:"description,omitempty"` 
	PublishedAt time.Time `json:"published_at"`
	URL         string    `json:"url"`
	FeedID      uuid.UUID `json:"feed_id"`
}


func UserToResponse(user *domain.User) UserResponse {
	return UserResponse{
		ID:        user.ID,
		CreatedAt: user.CreatedAt,
		UpdatedAt: user.UpdatedAt,
		Name:      user.Name,
		APIKey:    user.APIKey,
	}
}

func PostToResponse(post domain.Post) PostResponse {
	return PostResponse{
		ID:          post.ID,
		CreatedAt:   post.CreatedAt,
		UpdatedAt:   post.UpdatedAt,
		Title:       post.Title,
		Description: post.Description, 
		PublishedAt: post.PublishedAt,
		URL:         post.URL,
		FeedID:      post.FeedID,
	}
}

func PostsToResponse(posts []domain.Post) []PostResponse {
	responses := make([]PostResponse, len(posts))
	for i, post := range posts {
		responses[i] = PostToResponse(post)
	}
	return responses
}
