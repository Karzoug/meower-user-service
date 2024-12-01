package entity

import (
	"time"

	"github.com/rs/xid"
)

type UserShortProjection struct {
	ID         xid.ID
	Username   string `validate:"required,min=1,max=50"`
	Name       string `validate:"required,min=1,max=50"`
	ImageURL   string `validate:"omitempty,url,max=255"`
	StatusText string `validate:"omitempty,max=200"`
}

type User struct {
	UserShortProjection
	UpdatedAt time.Time
}

func (u User) CreatedAt() time.Time {
	return u.ID.Time()
}

func (u User) Validate() error {
	return validatorError(validate.Struct(u))
}

// NewUser creates a new user by given username.
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
