package repo

import "errors"

var (
	ErrNoAffected          = errors.New("no affected")
	ErrRecordAlreadyExists = errors.New("record already exists")
	ErrRecordNotFound      = errors.New("record not found")
)
