package errorc

import "net/http"

// NewHTTPBadRequest 400
func NewHTTPBadRequest(msgs ...string) *HTTPError {
	status := http.StatusBadRequest
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPUnauthorized 401
func NewHTTPUnauthorized(msgs ...string) *HTTPError {
	status := http.StatusUnauthorized
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPForbidden 403
func NewHTTPForbidden(msgs ...string) *HTTPError {
	status := http.StatusForbidden
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPNotFound 404
func NewHTTPNotFound(msgs ...string) *HTTPError {
	status := http.StatusNotFound
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPMethodNotAllowed 405
func NewHTTPMethodNotAllowed(msgs ...string) *HTTPError {
	status := http.StatusMethodNotAllowed
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPRequestTimeout 408
func NewHTTPRequestTimeout(msgs ...string) *HTTPError {
	status := http.StatusRequestTimeout
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPConflict 409
func NewHTTPConflict(msgs ...string) *HTTPError {
	status := http.StatusConflict
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPGone 410
func NewHTTPGone(msgs ...string) *HTTPError {
	status := http.StatusGone
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPLengthRequired 411
func NewHTTPLengthRequired(msgs ...string) *HTTPError {
	status := http.StatusLengthRequired
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPPreconditionFailed 412
func NewHTTPPreconditionFailed(msgs ...string) *HTTPError {
	status := http.StatusPreconditionFailed
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPRequestEntityTooLarge 413
func NewHTTPRequestEntityTooLarge(msgs ...string) *HTTPError {
	status := http.StatusRequestEntityTooLarge
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPRequestURITooLong 414
func NewHTTPRequestURITooLong(msgs ...string) *HTTPError {
	status := http.StatusRequestURITooLong
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPUnsupportedMediaType 415
func NewHTTPUnsupportedMediaType(msgs ...string) *HTTPError {
	status := http.StatusUnsupportedMediaType
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPTooManyRequests 429
func NewHTTPTooManyRequests(msgs ...string) *HTTPError {
	status := http.StatusTooManyRequests
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}

// NewHTTPInternal 500
func NewHTTPInternal(msgs ...string) *HTTPError {
	status := http.StatusInternalServerError
	message := http.StatusText(status)
	detail := message
	if len(msgs) > 0 {
		message = msgs[0]
	}
	if len(msgs) > 1 {
		detail = msgs[1]
	}
	return NewHTTPError(status, message, detail)
}
