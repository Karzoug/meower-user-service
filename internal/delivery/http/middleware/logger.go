package middleware

import (
	"context"
	"net/http"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
	"github.com/rs/zerolog"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/http/gen/user/v1"
)

// Logger is a middleware that logs incoming requests.
func Logger(logger zerolog.Logger) gen.StrictMiddlewareFunc {
	return func(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
		return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
			cw := newWrapResponseWriter(w)

			//nolint:zerologlint
			event := logger.Info().
				Str("remote_addr", r.RemoteAddr).
				Str("path", r.URL.Path).
				Str("handler_id", operationID).
				Str("proto", r.Proto).
				Str("method", r.Method).
				Str("user_agent", r.UserAgent()).
				Ctx(ctx) // for trace_id

			defer func() {
				event.
					Int("status_code", cw.Status()).
					Msg("incoming_request")
			}()

			return f(ctx, cw, r, request)
		}
	}
}
