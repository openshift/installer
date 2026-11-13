package powervs

import (
	"errors"
	"net/http"
	"testing"

	cosawserr "github.com/IBM/ibm-cos-sdk-go/aws/awserr"
	ibms3 "github.com/IBM/ibm-cos-sdk-go/service/s3"
	"github.com/stretchr/testify/assert"
)

func TestIsCOSS3NotFound(t *testing.T) {
	cases := []struct {
		name     string
		err      error
		expected bool
	}{
		{
			name:     "nil error",
			err:      nil,
			expected: false,
		},
		{
			name:     "generic error containing 404 in text must not match",
			err:      errors.New("failed with status 4040: connection timeout"),
			expected: false,
		},
		{
			name:     "generic error containing NoSuchBucket in text must not match",
			err:      errors.New("something NoSuchBucket unexpected happened"),
			expected: false,
		},
		{
			name:     "generic error containing NotFound in text must not match",
			err:      errors.New("resource NotFound on server"),
			expected: false,
		},
		{
			name:     "cosawserr RequestFailure with 404 status",
			err:      cosawserr.NewRequestFailure(cosawserr.New("NotFound", "Not Found", nil), http.StatusNotFound, "req-1"),
			expected: true,
		},
		{
			name:     "cosawserr RequestFailure with non-404 status",
			err:      cosawserr.NewRequestFailure(cosawserr.New("InternalError", "Internal Server Error", nil), http.StatusInternalServerError, "req-2"),
			expected: false,
		},
		{
			name:     "cosawserr Error with ErrCodeNoSuchBucket",
			err:      cosawserr.New(ibms3.ErrCodeNoSuchBucket, "The specified bucket does not exist", nil),
			expected: true,
		},
		{
			name:     "cosawserr Error with ErrCodeNoSuchKey",
			err:      cosawserr.New(ibms3.ErrCodeNoSuchKey, "The specified key does not exist", nil),
			expected: true,
		},
		{
			name:     "cosawserr Error with NotFound code (non-RequestFailure returns false)",
			err:      cosawserr.New("NotFound", "Not Found", nil),
			expected: false,
		},
		{
			name:     "cosawserr Error with 404 code (non-RequestFailure returns false)",
			err:      cosawserr.New("404", "Not Found", nil),
			expected: false,
		},
		{
			name:     "cosawserr Error with non-404 error code",
			err:      cosawserr.New("AccessDenied", "Access Denied", nil),
			expected: false,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			actual := isCOSS3NotFound(tc.err)
			assert.Equal(t, tc.expected, actual, "isCOSS3NotFound mismatch for %v", tc.err)
		})
	}
}

// TestIsCOSS3NotFound_Non404ErrorContaining404 verifies the regression case where a non-404 error
// has a message containing "404" (e.g. port 4040, resource ID with 404, or internal error 404-XYZ).
// It must NOT be treated as a not-found error so destroy does not falsely report success.
func TestIsCOSS3NotFound_Non404ErrorContaining404(t *testing.T) {
	testCases := []struct {
		name string
		err  error
	}{
		{
			name: "generic error with 404 substring in message",
			err:  errors.New("connection failed to host:4040: timeout"),
		},
		{
			name: "generic error with status code 404 in description but 500 actual failure",
			err:  errors.New("upstream gateway returned 502 while querying item 404"),
		},
		{
			name: "cosawserr RequestFailure with HTTP 500 whose message mentions 404",
			err: cosawserr.NewRequestFailure(
				cosawserr.New("InternalError", "internal server error: failed on item 404_bucket", nil),
				http.StatusInternalServerError,
				"req-500",
			),
		},
		{
			name: "cosawserr RequestFailure with HTTP 403 whose message mentions 404",
			err: cosawserr.NewRequestFailure(
				cosawserr.New("AccessDenied", "access denied to partition 404", nil),
				http.StatusForbidden,
				"req-403",
			),
		},
		{
			name: "cosawserr Error with non-404 code whose message mentions 404",
			err:  cosawserr.New("InternalError", "database error while processing object 404", nil),
		},
	}

	for _, tc := range testCases {
		t.Run(tc.name, func(t *testing.T) {
			assert.False(t, isCOSS3NotFound(tc.err), "expected isCOSS3NotFound to return false for non-404 error with '404' in message: %v", tc.err)
		})
	}
}
