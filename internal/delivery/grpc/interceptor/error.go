package interceptor

import (
	"context"
	"errors"
	"syscall"

	"github.com/rs/zerolog"
	"google.golang.org/grpc"
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"

	"github.com/Karzoug/meower-user-service/pkg/auth"
	"github.com/Karzoug/meower-user-service/pkg/ucerr"
	ucodes "github.com/Karzoug/meower-user-service/pkg/ucerr/codes"
)

func Error(logger zerolog.Logger) grpc.UnaryServerInterceptor {
	return func(ctx context.Context, req any, info *grpc.UnaryServerInfo, handler grpc.UnaryHandler) (resp any, err error) {
		defer func() {
			if nil == err {
				return
			}

			switch e := err.(type) { //nolint:errorlint
			// it's a service/usecase layer error (= trusted), just return it
			case ucerr.Error:
				logServiceError(ctx, e, info.FullMethod, logger)
				err = status.Error(e.GrpcCode(), e.Error())
			default:
				switch {
				// we could try to write response and got network error,
				// it's ok, so ignore it
				case isNetworkError(e):

				// it's an authN error on delivery layer (= trusted)
				case auth.IsAuthNError(e):
					err = status.Error(codes.Unauthenticated, e.Error())

				// it's unknown (= untrusted) error,
				// log it and return internal server error
				default:
					logger.Error().
						Err(e).
						Ctx(ctx). // for trace_id
						Str("method", info.FullMethod).
						Msg("error handler")
					err = status.Error(codes.Internal, codes.Internal.String())
				}
			}
		}()

		return handler(ctx, req)
	}
}

func logServiceError(ctx context.Context, err ucerr.Error, method string, logger zerolog.Logger) {
	var ev *zerolog.Event

	switch err.Code() {
	case ucodes.OK,
		ucodes.Canceled,
		ucodes.InvalidArgument,
		ucodes.DeadlineExceeded,
		ucodes.NotFound,
		ucodes.AlreadyExists,
		ucodes.PermissionDenied,
		ucodes.FailedPrecondition,
		ucodes.OutOfRange,
		ucodes.Unimplemented, ucodes.Unauthenticated:
		return

	case ucodes.ResourceExhausted, ucodes.Aborted:
		ev = logger.Warn()

	case ucodes.Internal, ucodes.Unavailable, ucodes.Unknown, ucodes.DataLoss:
		ev = logger.Error()

	}

	ev.Err(err.Unwrap()).
		Ctx(ctx). // for trace_id
		Str("method", method).
		Msg("error handler")
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
