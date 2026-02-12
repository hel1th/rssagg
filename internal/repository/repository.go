package repository

import (
	"github.com/hel1th/rssagg/internal/database"
)

type Repositories struct {
	User       UserRepository
	Feed       FeedRepository
	FeedFollow FeedFollowRepository
	Post       PostRepository
}

func NewRepositories(db *database.Queries) *Repositories {
	return &Repositories{
		User:       NewUserRepository(db),
		Feed:       NewFeedRepository(db),
		FeedFollow: NewFeedFollowRepository(db),
		Post:       NewPostRepository(db),
	}
}
