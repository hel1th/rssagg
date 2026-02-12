package domain

import (
	"time"

	"github.com/google/uuid"
)

type User struct {
	ID        uuid.UUID
	CreatedAt time.Time
	UpdatedAt time.Time
	Name      string
	APIKey    string
}

func NewUser(name string) *User {
	now := time.Now().UTC()
	return &User{
		ID:        uuid.New(),
		CreatedAt: now,
		UpdatedAt: now,
		Name:      name,
	}
}

func (u *User) Validate() error {
	if u.Name == "" {
		return ErrInvalidUserName
	}
	if len(u.Name) > 255 {
		return ErrUserNameTooLong
	}
	return nil
}
