package errfunc

import (
	"errors"
	"net/http"
	"syscall"

	"go.opentelemetry.io/otel/trace"

	gen "github.com/Karzoug/meower-user-service/internal/delivery/http/gen/user/v1"
	"github.com/Karzoug/meower-user-service/internal/delivery/http/response"
	"github.com/Karzoug/meower-user-service/pkg/trace/otlp"

	"github.com/rs/zerolog"
)

func JSONRequest(logger zerolog.Logger, tracer trace.Tracer) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		ctx := otlp.InjectTracing(r.Context(), tracer)
		logger.Info().
			Str("remote_addr", r.RemoteAddr).
			Str("path", r.URL.Path).
			Str("proto", r.Proto).
			Str("method", r.Method).
			Str("user_agent", r.UserAgent()).
			Ctx(ctx). // for trace_id
			Int("status_code", http.StatusBadRequest).
			Msg("incoming_request")

		respErr := response.JSON(w,
			http.StatusBadRequest,
			gen.ErrorResponse{Error: err.Error()},
		)

		// if we can't write error response, log it, unless it's network error
		if respErr != nil && !isNetworkError(respErr) {
			logger.Error().
				Err(respErr).
				Ctx(ctx). // for trace_id
				Msg("error handler: couldn't write error response")
		}
	}
}

func JSONResponse(logger zerolog.Logger) func(w http.ResponseWriter, r *http.Request, err error) {
	return func(w http.ResponseWriter, r *http.Request, err error) {
		logger.Error().Err(err).Msg("error handler: uncaught error")
		_ = response.JSON(w,
			http.StatusInternalServerError,
			gen.ErrorResponse{Error: http.StatusText(http.StatusInternalServerError)})
	}
}

func isNetworkError(err error) bool {
	// Ignore syscall.EPIPE and syscall.ECONNRESET errors which occurs
	// when a write operation happens on the http.ResponseWriter that
	// has simultaneously been disconnected by the client (TCP
	// connections is broken). For instance, when large amounts of
	// data is being written or streamed to the client.
	// https://blog.cloudflare.com/the-complete-guide-to-golang-net-http-timeouts/
	// https://gosamples.dev/broken-pipe/
	// https://gosamples.dev/connection-reset-by-peer/

	switch {
	case errors.Is(err, syscall.EPIPE):

		// Usually, you get the broken pipe error when you write to the connection after the
		// RST (TCP RST Flag) is sent.
		// The broken pipe is a TCP/IP error occurring when you write to a stream where the
		// other end (the peer) has closed the underlying connection. The first write to the
		// closed connection causes the peer to reply with an RST packet indicating that the
		// connection should be terminated immediately. The second write to the socket that
		// has already received the RST causes the broken pipe error.
		return true

	case errors.Is(err, syscall.ECONNRESET):

		// Usually, you get connection reset by peer error when you read from the
		// connection after the RST (TCP RST Flag) is sent.
		// The connection reset by peer is a TCP/IP error that occurs when the other end (peer)
		// has unexpectedly closed the connection. It happens when you send a packet from your
		// end, but the other end crashes and forcibly closes the connection with the RST
		// packet instead of the TCP FIN, which is used to close a connection under normal
		// circumstances.
		return true
	}

	return false
}
