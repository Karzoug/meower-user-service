package middleware

import (
	"context"
	"fmt"
	"net/http"
	"runtime/debug"

	"github.com/oapi-codegen/runtime/strictmiddleware/nethttp"
)

// Recoverer is a middleware that prevents panics and logs them.
func Recover(f nethttp.StrictHTTPHandlerFunc, operationID string) nethttp.StrictHTTPHandlerFunc {
	return func(ctx context.Context, w http.ResponseWriter, r *http.Request, request any) (response any, err error) {
		defer func() {
			if rec := recover(); rec != nil {
				if rec == http.ErrAbortHandler { //nolint:errorlint
					// we don't recover http.ErrAbortHandler so the response
					// to the client is aborted, this should not be logged
					panic(rec)
				}

				stack := debug.Stack()
				err = fmt.Errorf("recovered panic: %v; stack: %s", rec, string(stack))
			}
		}()

		return f(ctx, w, r, request)
	}
}
