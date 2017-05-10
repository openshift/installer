package api

import (
	"fmt"
	"net/http"
)

// httpError represents an error that occurred in one of the API handlers.
type httpError struct {
	status  int
	message string
}

func newError(status int, format string, a ...interface{}) *httpError {
	return &httpError{
		status:  status,
		message: fmt.Sprintf(format, a...),
	}
}

func newNotFoundError(format string, a ...interface{}) *httpError {
	return newError(http.StatusNotFound, format, a...)
}

func newBadRequestError(format string, a ...interface{}) *httpError {
	return newError(http.StatusInternalServerError, format, a...)
}

func newInternalServerError(format string, a ...interface{}) *httpError {
	return newError(http.StatusBadRequest, format, a...)
}

func (he *httpError) Error() string {
	return he.message
}
