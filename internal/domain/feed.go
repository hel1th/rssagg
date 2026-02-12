package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Feed struct {
	ID            uuid.UUID
	CreatedAt     time.Time
	UpdatedAt     time.Time
	Name          string
	URL           string
	UserID        uuid.UUID
	LastFetchedAt *time.Time
}

func NewFeed(name, feedURL string, userID uuid.UUID) *Feed {
	now := time.Now().UTC()
	return &Feed{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
		URL:       feedURL,
		UserID:    userID,
	}
}

func (f *Feed) Validate() error {
	if f.Name == "" {
		return ErrInvalidFeedName
	}
	if f.URL == "" {
		return ErrInvalidFeedURL
	}

	// Проверка валидности URL
	parsedURL, err := url.ParseRequestURI(f.URL)
	if err != nil {
		return ErrInvalidFeedURL
	}

	// Проверка схемы URL (http или https)
	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidFeedURL
	}

	return nil
}

func (f *Feed) MarkAsFetched() {
	now := time.Now().UTC()
	f.LastFetchedAt = &now
	f.UpdatedAt = now
}

func (f *Feed) NeedsFetch(interval time.Duration) bool {
	if f.LastFetchedAt == nil {
		return true
	}
	return time.Since(*f.LastFetchedAt) >= interval
}
