package entity

import (
	"time"

	"github.com/rs/xid"
)

type UserShortProjection struct {
	ID         xid.ID
	Username   string
	Name       string
	ImageURL   string
	StatusText string
}

type User struct {
	UserShortProjection
	UpdatedAt time.Time
}

func (u User) CreatedAt() time.Time {
	return u.ID.Time()
}

func NewUser(username string) User {
	id := xid.New()
	return User{
		UserShortProjection: UserShortProjection{
			ID:       id,
			Username: username,
			Name:     username,
		},
		UpdatedAt: id.Time(),
	}
}
