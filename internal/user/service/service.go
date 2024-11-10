package service

import (
	"github.com/rs/zerolog"
)

type UserService struct {
	logger zerolog.Logger
}

func NewUserService(logger zerolog.Logger) UserService {
	logger = logger.With().
		Str("component", "user service").
		Logger()

	return UserService{
		logger: logger,
	}
}
