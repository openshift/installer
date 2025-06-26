package aws

import (
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
)

// handleErrorCode takes the error and extracts the error code if it was successfully cast as an API Error.
func handleErrorCode(err error) string {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode()
	}
	return ""
}
