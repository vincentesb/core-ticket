package error_code

import "net/http"

type ErrorCode string

const EC_ = "EC031"

const (
	Success            ErrorCode = "00000"
	Unauthorized       ErrorCode = "00001"
	Forbidden          ErrorCode = "00002"
	ValidationError    ErrorCode = "00003"
	NotFound           ErrorCode = "00004"
	ServiceUnavailable ErrorCode = "00005"
	UnknownError       ErrorCode = "00032"
)

func (ec ErrorCode) String() string {
	return EC_ + string(ec)
}

func (ec ErrorCode) Message() (message string) {
	switch ec {
	case Success:
		message = "OK"
	case Unauthorized:
		message = "Unauthorized"
	case Forbidden:
		message = "Forbidden"
	case ValidationError:
		message = "Validation Error"
	case NotFound:
		message = "Not Found"
	case ServiceUnavailable:
		message = "Service Unavailable"
	case UnknownError:
		message = "Unknown Error"
	}

	return message
}

func (ec ErrorCode) HttpStatusCode() int {
	switch ec {
	case NotFound, ValidationError:
		return http.StatusBadRequest
	case Unauthorized:
		return http.StatusUnauthorized
	case Forbidden:
		return http.StatusForbidden
	case ServiceUnavailable:
		return http.StatusServiceUnavailable
	default:
		return http.StatusInternalServerError
	}
}
