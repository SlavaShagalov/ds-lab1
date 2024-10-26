package errors

import "net/http"

var httpCodes = map[error]int{
	// Common repository
	ErrDb: http.StatusInternalServerError,

	// Users
	ErrPersonNotFound:      http.StatusNotFound,
	ErrPersonAlreadyExists: http.StatusConflict,

	// HTTP
	ErrReadBody: http.StatusBadRequest,
}

func GetHTTPCodeByError(err error) (int, bool) {
	httpCode, exist := httpCodes[err]
	if !exist {
		httpCode = http.StatusInternalServerError
	}
	return httpCode, exist
}
