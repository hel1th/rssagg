package domain

import (
	"time"

	"github.com/google/uuid"
)

type FeedFollow struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	UserID    uuid.UUID
	FeedID    uuid.UUID
}

func NewFeedFollow(userID, feedID uuid.UUID) *FeedFollow {
	now := time.Now().UTC()
	return &FeedFollow{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		UserID:    userID,
		FeedID:    feedID,
	}
}

func (ff *FeedFollow) Validate() error {
	if ff.UserID == uuid.Nil {
		return ErrInvalidUserID
	}
	if ff.FeedID == uuid.Nil {
		return ErrInvalidFeedID
	}
	return nil
}
