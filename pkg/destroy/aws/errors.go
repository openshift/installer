package aws

import (
	"github.com/aws/smithy-go"
	"github.com/pkg/errors"
)

// HandleErrorCode takes the error and extracts the error code if it was successfully cast as an API Error.
func HandleErrorCode(err error) string {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode()
	}
	return ""
}
