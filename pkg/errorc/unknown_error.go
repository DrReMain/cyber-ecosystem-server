package errorc

type UnknownError struct {
	Detail string
}

func (e *UnknownError) Error() string {
	return e.Detail
}

func NewUnknownError(err error) *UnknownError {
	return &UnknownError{
		Detail: err.Error(),
	}
}
