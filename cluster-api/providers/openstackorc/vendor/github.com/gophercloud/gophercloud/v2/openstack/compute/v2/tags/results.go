package tags

import (
	"net/http"

	"github.com/gophercloud/gophercloud/v2"
)

type commonResult struct {
	gophercloud.Result
}

// Extract is a function that accepts a result and extracts a tags resource.
func (r commonResult) Extract() ([]string, error) {
	var s struct {
		Tags []string `json:"tags"`
	}
	err := r.ExtractInto(&s)
	return s.Tags, err
}

type ListResult struct {
	commonResult
}

// CheckResult is the result from the Check operation.
type CheckResult struct {
	gophercloud.Result
}

func (r CheckResult) Extract() (bool, error) {
	exists := r.Err == nil

	if gophercloud.ResponseCodeIs(r.Err, http.StatusNotFound) {
		r.Err = nil
	}

	return exists, r.Err
}

// ReplaceAllResult is the result from the ReplaceAll operation.
type ReplaceAllResult struct {
	commonResult
}

// AddResult is the result from the Add operation.
type AddResult struct {
	gophercloud.ErrResult
}

// DeleteResult is the result from the Delete operation.
type DeleteResult struct {
	gophercloud.ErrResult
}
