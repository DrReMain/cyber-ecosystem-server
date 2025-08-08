package errorc

type HTTPError struct {
	Status  int
	Detail  string
	Code    string
	Message string
}

func (e *HTTPError) Error() string {
	return e.Detail
}

func NewHTTPError(status int, detail, code, message string) *HTTPError {
	return &HTTPError{
		Status:  status,
		Detail:  detail,
		Code:    code,
		Message: message,
	}
}
