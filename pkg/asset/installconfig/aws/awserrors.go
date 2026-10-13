package aws

import (
	"errors"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
)

// Error constants for AWS error codes.
const (
	AccessDeniedException   = "AccessDeniedException"
	NoSuchResourceException = "NoSuchResourceException"
)

// IsUnauthorized checks if the error is due to lacking permissions.
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		// see reference:
		// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetServiceQuota.html
		// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetAWSDefaultServiceQuota.html
		return apiErr.ErrorCode() == AccessDeniedException || apiErr.ErrorCode() == NoSuchResourceException
	}
	return false
}

// IsHTTPForbidden returns true if and only if the error is an HTTP
// 403 error from the AWS API.
func IsHTTPForbidden(err error) bool {
	if err == nil {
		return false
	}

	var respErr *awshttp.ResponseError
	if errors.As(err, &respErr) {
		return respErr.HTTPStatusCode() == http.StatusForbidden
	}
	return false
}
