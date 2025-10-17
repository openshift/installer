/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package genericarmclient

import (
	"encoding/json"

	"github.com/pkg/errors"

	"github.com/Azure/azure-service-operator/v2/pkg/genruntime/core"
)

// CloudError - An error response for a resource management request
// We have to support two different formats for the error as some services do things differently.
//
// The ARM spec says that error details should be nested inside an `error` element.
// See https://github.com/Azure/azure-resource-manager-rpc/blob/master/v1.0/common-api-details.md#error-response-content
//
// However, some services put the code & message at the top level instead.
// This is common enough that the Azure Python SDK has specific handling to promote a nested error to the top level.
// See https://github.com/Azure/azure-sdk-for-python/blob/9791fb5bc4cb6001768e6e1fb986b8d8f8326c43/sdk/core/azure-core/azure/core/exceptions.py#L153
type CloudError struct {
	error error

	code    *string
	message *string
	target  *string
	details []*ErrorResponse
}

// NewCloudError returns a new CloudError
func NewCloudError(err error) *CloudError {
	return &CloudError{
		error: err,
	}
}

func NewTestCloudError(code string, message string) *CloudError {
	return &CloudError{
		code:    &code,
		message: &message,
	}
}

// Error implements the error interface for type CloudError.
// The contents of the error text are not contractual and subject to change.
func (e CloudError) Error() string {
	return e.error.Error()
}

// Code returns the error code from the message, if present, or UnknownErrorCode if not.
func (e CloudError) Code() string {
	if e.code != nil && *e.code != "" {
		return *e.code
	}

	return core.UnknownErrorCode
}

// Message returns the message from the error, if present, or UnknownErrorMessage if not.
func (e CloudError) Message() string {
	if e.message != nil && *e.message != "" {
		return *e.message
	}

	return core.UnknownErrorMessage
}

// Target returns the target of the error, if present, or an empty string if not.
func (e CloudError) Target() string {
	if e.target != nil && *e.target != "" {
		return *e.target
	}

	return ""
}

// Details returns the details of the error, if present, or an empty slice if not
func (e CloudError) Details() []*ErrorResponse {
	return e.details
}

func (e *CloudError) UnmarshalJSON(data []byte) error {
	var content struct {
		Code       *string          `json:"code,omitempty"`
		Message    *string          `json:"message,omitempty"`
		Target     *string          `json:"target,omitempty"`
		Details    []*ErrorResponse `json:"details,omitempty"`
		InnerError *struct {
			Code    *string          `json:"code,omitempty"`
			Message *string          `json:"message,omitempty"`
			Target  *string          `json:"target,omitempty"`
			Details []*ErrorResponse `json:"details,omitempty"`
		} `json:"error,omitempty"`
	}

	err := json.Unmarshal(data, &content)
	if err != nil {
		return errors.Wrap(err, "unmarshalling JSON for CloudError")
	}

	if content.InnerError != nil {
		e.code = content.InnerError.Code
		e.message = content.InnerError.Message
		e.target = content.InnerError.Target
		e.details = content.InnerError.Details
	} else {
		e.code = content.Code
		e.message = content.Message
		e.target = content.Target
		e.details = content.Details
	}

	return nil
}

func (e CloudError) Unwrap() error {
	return e.error
}
