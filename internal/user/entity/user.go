package entity

import "time"

type ShortUserInfo struct {
	ID         string `db:"id" validate:"required,min=1,max=50"`
	Name       string `db:"name" validate:"required,min=1,max=50"`
	ImageURL   string `db:"image_url" validate:"url"`
	StatusText string `db:"status_text" validate:"max=200"`
}

type User struct {
	ShortUserInfo
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
}

func (u User) Validate() error {
	return validatorError(validate.Struct(u))
}

func NewUser(id string) User {
	return User{
		ShortUserInfo: ShortUserInfo{
			ID:   id,
			Name: id,
		},
		CreatedAt: time.Now(),
		UpdatedAt: time.Now(),
	}
}
