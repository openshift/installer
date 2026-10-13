package gcp

import (
	"errors"
	"fmt"
	"net/http"
	"testing"

	"google.golang.org/api/googleapi"
)

func Test_newQuota(t *testing.T) {
	cases := []struct {
		name string
	}{{
		name: "missing usage",
	}, {
		name: "usage in single zone",
	}, {
		name: "usage in multiple zones",
	}, {
		name: "",
	}}
	for _, test := range cases {
		t.Run(test.name, func(t *testing.T) {

		})
	}
}

func TestIsTransient(t *testing.T) {
	tests := []struct {
		name string
		err  error
		want bool
	}{
		{name: "nil"},
		{name: "unrelated error", err: errors.New("unrelated")},
		{name: "bad request", err: &googleapi.Error{Code: http.StatusBadRequest}},
		{name: "rate limited", err: &googleapi.Error{Code: http.StatusTooManyRequests}, want: true},
		{name: "internal server error", err: &googleapi.Error{Code: http.StatusInternalServerError}, want: true},
		{name: "service unavailable", err: &googleapi.Error{Code: http.StatusServiceUnavailable}, want: true},
		{name: "wrapped transient error", err: fmt.Errorf("loading quota: %w", &googleapi.Error{Code: http.StatusServiceUnavailable}), want: true},
	}

	for _, test := range tests {
		t.Run(test.name, func(t *testing.T) {
			if got := IsTransient(test.err); got != test.want {
				t.Errorf("IsTransient() = %t, want %t", got, test.want)
			}
		})
	}
}
