package middleware

import (
	"context"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/http/gen/user/v1"
	"github.com/Karzoug/meower-user-service/pkg/auth"
)

const authHeader = "X-User"

// AuthN is a middleware that adds an username from the request "X-User" Header to the context.
// (!) The middleware doesn't check if the token is valid - it's up to the gateway.
func AuthN(next nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (resp any, err error) {
		sub := r.Header.Get(authHeader)
		if sub != "" {
			ctx = auth.WithUserID(ctx, sub)
		}

		// if spec claim authentification
		if ctx.Value(gen.OAuthScopes) != nil && sub == "" {
			return nil, auth.NewAuthNError()
		}

		return next(ctx, w, r, request)
	}
}
