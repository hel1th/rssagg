package domain

import "errors"

var (
	ErrInternalServer = errors.New("internal server error")
	ErrUnauthorized   = errors.New("unauthorized")
	ErrForbidden      = errors.New("forbidden")
)

var (
	ErrInvalidUserName = errors.New("invalid user name")
	ErrUserNameTooLong = errors.New("user name is too long")
	ErrUserNotFound    = errors.New("user not found")
	ErrInvalidAPIKey   = errors.New("invalid API key")
	ErrInvalidUserID   = errors.New("invalid user ID")
)

var (
	ErrInvalidFeedName = errors.New("invalid feed name")
	ErrInvalidFeedURL  = errors.New("invalid feed URL")
	ErrFeedNotFound    = errors.New("feed not found")
	ErrInvalidFeedID   = errors.New("invalid feed ID")
	ErrDuplicateFeed   = errors.New("feed already exists")
)

var (
	ErrFeedFollowNotFound  = errors.New("feed follow not found")
	ErrDuplicateFeedFollow = errors.New("already following this feed")
	ErrCannotUnfollowFeed  = errors.New("cannot unfollow feed")
)

var (
	ErrInvalidPostTitle   = errors.New("invalid post title")
	ErrInvalidPostURL     = errors.New("invalid post URL")
	ErrInvalidPublishedAt = errors.New("invalid published date")
	ErrPostNotFound       = errors.New("post not found")
	ErrDuplicatePost      = errors.New("post already exists")
)
