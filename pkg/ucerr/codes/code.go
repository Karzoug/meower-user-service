package codes

// Code for service/usecase level error.
type Code uint32

// Codes used for service/usecase level errors.
const (
	// OK is returned on success.
	OK Code = 0

	// Canceled indicates the operation was canceled.
	Canceled Code = 1

	// Unknown error.
	Unknown Code = 2

	// InvalidArgument indicates client specified an invalid argument.
	// Note that this differs from FailedPrecondition. It indicates arguments
	// that are problematic regardless of the state of the system
	// (e.g., a malformed file name).
	InvalidArgument Code = 3

	// DeadlineExceeded means operation expired before completion.
	// For operations that change the state of the system, this error may be
	// returned even if the operation has completed successfully. For
	// example, a successful response from a server could have been delayed
	// long enough for the deadline to expire.
	DeadlineExceeded Code = 4

	// NotFound means some requested entity (e.g., file or directory) was
	// not found.
	NotFound Code = 5

	// AlreadyExists means an attempt to create an entity failed because one
	// already exists.
	AlreadyExists Code = 6

	// PermissionDenied indicates the caller does not have permission to
	// execute the specified operation.
	PermissionDenied Code = 7

	// ResourceExhausted indicates some resource has been exhausted, perhaps
	// a per-user quota, or perhaps the entire file system is out of space.
	ResourceExhausted Code = 8

	// FailedPrecondition indicates operation was rejected because the
	// system is not in a state required for the operation's execution.
	// For example, directory to be deleted may be non-empty, an rmdir
	// operation is applied to a non-directory, etc.
	FailedPrecondition Code = 9

	// Aborted indicates the operation was aborted, typically due to a
	// concurrency issue like sequencer check failures, transaction aborts,
	// etc.
	Aborted Code = 10

	// OutOfRange means operation was attempted past the valid range.
	// E.g., seeking or reading past end of file.
	//
	// Unlike InvalidArgument, this error indicates a problem that may
	// be fixed if the system state changes.
	OutOfRange Code = 11

	// Unimplemented indicates operation is not implemented or not
	// supported/enabled in this service.
	Unimplemented Code = 12

	// Internal errors. Means some invariants expected by underlying
	// system has been broken.
	Internal Code = 13

	// Unavailable indicates the service is currently unavailable.
	Unavailable Code = 14

	// DataLoss indicates unrecoverable data loss or corruption.
	DataLoss Code = 15

	// Unauthenticated indicates the request does not have valid
	// authentication credentials for the operation.
	Unauthenticated Code = 16
)
