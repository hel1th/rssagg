package domain

import (
	"net/url"
	"time"

	"github.com/google/uuid"
)

type Post struct {
	ID          uuid.UUID
	CreatedAt   time.Time
	UpdatedAt   time.Time
	Title       string
	Description *string
	PublishedAt time.Time
	URL         string
	FeedID      uuid.UUID
}

func NewPost(title, postURL string, publishedAt time.Time, feedID uuid.UUID, description *string) *Post {
	now := time.Now().UTC()
	return &Post{
		ID:          uuid.New(),
		CreatedAt:   now,
		UpdatedAt:   now,
		Title:       title,
		Description: description,
		PublishedAt: publishedAt,
		URL:         postURL,
		FeedID:      feedID,
	}
}

func (p *Post) Validate() error {
	if p.Title == "" {
		return ErrInvalidPostTitle
	}
	if p.URL == "" {
		return ErrInvalidPostURL
	}

	parsedURL, err := url.ParseRequestURI(p.URL)
	if err != nil {
		return ErrInvalidPostURL
	}

	if parsedURL.Scheme != "http" && parsedURL.Scheme != "https" {
		return ErrInvalidPostURL
	}

	if p.FeedID == uuid.Nil {
		return ErrInvalidFeedID
	}

	if p.PublishedAt.IsZero() {
		return ErrInvalidPublishedAt
	}

	return nil
}

func (p *Post) HasDescription() bool {
	return p.Description != nil && *p.Description != ""
}
