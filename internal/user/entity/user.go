package entity

import "time"

type ShortUserInfo struct {
	ID         string `db:"id"`
	Name       string `db:"name"`
	ImageUrl   string `db:"image_url"`
	StatusText string `db:"status_text"`
}

type User struct {
	ShortUserInfo
	CreatedAt time.Time `db:"created_at"`
	UpdatedAt time.Time `db:"updated_at"`
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
