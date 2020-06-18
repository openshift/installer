package packet

import (
	"net/http"
	"strings"

	"github.com/packethost/packngo"
)

func friendlyError(err error) error {
	if e, ok := err.(*packngo.ErrorResponse); ok {
		errors := Errors(e.Errors)
		// if packngo gives us blank error strings, populate them with something useful
		// this is useful so the user gets some sort of indication of a failure rather than a blank message
		if 0 == len(errors) {
			errors = Errors{e.SingleError}
		}
		return &ErrorResponse{
			StatusCode: e.Response.StatusCode,
			Errors:     errors,
		}
	}
	return err
}

func isForbidden(err error) bool {
	if r, ok := err.(*ErrorResponse); ok {
		return r.StatusCode == http.StatusForbidden
	}
	return false
}

func isNotFound(err error) bool {
	if r, ok := err.(*ErrorResponse); ok {
		return r.StatusCode == http.StatusNotFound
	}
	return false
}

type Errors []string

func (e Errors) Error() string {
	return strings.Join(e, "; ")
}

type ErrorResponse struct {
	StatusCode int
	Errors
}
