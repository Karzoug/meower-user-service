package ucerr

import (
	"net/http"

	"google.golang.org/grpc/codes"

	uc "github.com/Karzoug/meower-user-service/pkg/ucerr/codes"
)

const statusClientClosedRequest = 499

// Error is a service/usecase level error.
type Error struct {
	msg  string
	err  error
	code uc.Code
}

func NewError(err error, msg string, code uc.Code) Error {
	return Error{
		msg:  msg,
		err:  err,
		code: code,
	}
}

func NewInternalError(err error) Error {
	return Error{
		msg:  "Internal error",
		err:  err,
		code: uc.Internal,
	}
}

// Error returns error message which can be returned to the client.
func (e Error) Error() string {
	return e.msg
}

func (e Error) Code() uc.Code {
	return e.code
}

func (e Error) Unwrap() error {
	return e.err
}

func (e Error) HTTPCode() int {
	switch e.code {
	case uc.Aborted, uc.AlreadyExists:
		return http.StatusConflict
	case uc.Canceled:
		return statusClientClosedRequest
	case uc.DeadlineExceeded:
		return http.StatusGatewayTimeout
	case uc.InvalidArgument, uc.FailedPrecondition, uc.OutOfRange:
		return http.StatusBadRequest
	case uc.NotFound:
		return http.StatusNotFound
	case uc.OK:
		return http.StatusOK
	case uc.PermissionDenied:
		return http.StatusForbidden
	case uc.ResourceExhausted:
		return http.StatusTooManyRequests
	case uc.Unauthenticated:
		return http.StatusUnauthorized
	case uc.Unavailable:
		return http.StatusServiceUnavailable
	case uc.Unimplemented:
		return http.StatusNotImplemented
	default: // uc.Unknown, uc.Internal, uc.DataLoss
		return http.StatusInternalServerError
	}
}

func (e Error) GrpcCode() codes.Code {
	return codes.Code(e.code)
}
