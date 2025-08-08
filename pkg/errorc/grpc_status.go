package errorc

import (
	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

// GRPCCanceledError 1
func GRPCCanceledError(msg string) error {
	return status.Error(codes.Canceled, msg)
}

// GRPCUnknownError 2
func GRPCUnknownError(msg string) error {
	return status.Error(codes.Unknown, msg)
}

// GRPCInvalidArgumentError 3
func GRPCInvalidArgumentError(msg string) error {
	return status.Error(codes.InvalidArgument, msg)
}

// GRPCDeadlineExceededError 4
func GRPCDeadlineExceededError(msg string) error {
	return status.Error(codes.DeadlineExceeded, msg)
}

// GRPCNotFoundError 5
func GRPCNotFoundError(msg string) error {
	return status.Error(codes.NotFound, msg)
}

// GRPCAlreadyExistsError 6
func GRPCAlreadyExistsError(msg string) error {
	return status.Error(codes.AlreadyExists, msg)
}

// GRPCPermissionDeniedError 7
func GRPCPermissionDeniedError(msg string) error {
	return status.Error(codes.PermissionDenied, msg)
}

// GRPCResourceExhaustedError 8
func GRPCResourceExhaustedError(msg string) error {
	return status.Error(codes.ResourceExhausted, msg)
}

// GRPCFailedPreconditionError 9
func GRPCFailedPreconditionError(msg string) error {
	return status.Error(codes.FailedPrecondition, msg)
}

// GRPCAbortedError 10
func GRPCAbortedError(msg string) error {
	return status.Error(codes.Aborted, msg)
}

// GRPCOutOfRangeError 11
func GRPCOutOfRangeError(msg string) error {
	return status.Error(codes.OutOfRange, msg)
}

// GRPCUnimplementedError 12
func GRPCUnimplementedError(msg string) error {
	return status.Error(codes.Unimplemented, msg)
}

// GRPCInternalError 13
func GRPCInternalError(msg string) error {
	return status.Error(codes.Internal, msg)
}

// GRPCUnavailableError 14
func GRPCUnavailableError(msg string) error {
	return status.Error(codes.Unavailable, msg)
}

// GRPCDataLossError 15
func GRPCDataLossError(msg string) error {
	return status.Error(codes.DataLoss, msg)
}

// GRPCUnauthenticatedError 16
func GRPCUnauthenticatedError(msg string) error {
	return status.Error(codes.Unauthenticated, msg)
}
