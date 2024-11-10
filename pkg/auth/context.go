package auth

import (
	"context"
	"errors"
)

type authKey struct{}

var authUserIDKey authKey

var errAuthN = errors.New("authentication required")

func UserIDFromContext(ctx context.Context) string {
	username, ok := ctx.Value(authUserIDKey).(string)
	if !ok {
		return ""
	}
	return username
}

func WithUserID(ctx context.Context, username string) context.Context {
	return context.WithValue(ctx, authUserIDKey, username)
}

func NewAuthNError() error {
	return errAuthN
}

func IsAuthNError(err error) bool {
	return errors.Is(err, errAuthN)
}
