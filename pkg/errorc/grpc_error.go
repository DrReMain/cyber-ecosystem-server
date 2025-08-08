package errorc

import (
	"errors"

	"google.golang.org/grpc/codes"
	"google.golang.org/grpc/status"
)

type GRPCError struct {
	Status  int
	Detail  string
	Code    codes.Code
	Message string
}

func (e *GRPCError) Error() string {
	return e.Detail
}

func NewGRPCError(err error) (e *GRPCError) {
	if !errors.As(err, &e) {
		switch s := status.Convert(err); s.Code() {
		case codes.Canceled:
			e = &GRPCError{499, err.Error(), s.Code(), s.Message()}
		case codes.Unknown:
			e = &GRPCError{500, err.Error(), s.Code(), s.Message()}
		case codes.InvalidArgument:
			e = &GRPCError{400, err.Error(), s.Code(), s.Message()}
		case codes.DeadlineExceeded:
			e = &GRPCError{504, err.Error(), s.Code(), s.Message()}
		case codes.NotFound:
			e = &GRPCError{404, err.Error(), s.Code(), s.Message()}
		case codes.AlreadyExists:
			e = &GRPCError{409, err.Error(), s.Code(), s.Message()}
		case codes.PermissionDenied:
			e = &GRPCError{403, err.Error(), s.Code(), s.Message()}
		case codes.ResourceExhausted:
			e = &GRPCError{503, err.Error(), s.Code(), s.Message()}
		case codes.FailedPrecondition:
			e = &GRPCError{412, err.Error(), s.Code(), s.Message()}
		case codes.Aborted:
			e = &GRPCError{409, err.Error(), s.Code(), s.Message()}
		case codes.OutOfRange:
			e = &GRPCError{400, err.Error(), s.Code(), s.Message()}
		case codes.Unimplemented:
			e = &GRPCError{501, err.Error(), s.Code(), s.Message()}
		case codes.Internal:
			e = &GRPCError{500, err.Error(), s.Code(), s.Message()}
		case codes.Unavailable:
			e = &GRPCError{503, err.Error(), s.Code(), s.Message()}
		case codes.DataLoss:
			e = &GRPCError{500, err.Error(), s.Code(), s.Message()}
		case codes.Unauthenticated:
			e = &GRPCError{401, err.Error(), s.Code(), s.Message()}
		default:
			e = &GRPCError{500, err.Error(), s.Code(), s.Message()}
		}
	}
	return
}
