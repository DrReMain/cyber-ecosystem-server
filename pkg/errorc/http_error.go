package errorc

import "fmt"

type HTTPError struct {
	Status  int
	Detail  string
	Code    int
	Message string
}

func (e *HTTPError) Error() string {
	return e.Detail
}

func NewHTTPError(status int, message, detail string) (e *HTTPError) {
	return &HTTPError{
		Status:  status,
		Detail:  fmt.Sprintf("%d %s", status, detail),
		Code:    status,
		Message: message,
	}
}
