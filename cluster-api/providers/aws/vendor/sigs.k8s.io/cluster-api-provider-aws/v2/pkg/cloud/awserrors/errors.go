/*
Copyright 2018 The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

	http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// Package awserrors provides a way to generate AWS errors.
package awserrors

import (
	"net/http"

	"github.com/aws/aws-sdk-go/aws/awserr"
	"github.com/aws/aws-sdk-go/service/ssm"
)

// Error singletons for AWS errors.
const (
	AssociationIDNotFound             = "InvalidAssociationID.NotFound"
	AuthFailure                       = "AuthFailure"
	BucketAlreadyOwnedByYou           = "BucketAlreadyOwnedByYou"
	EIPNotFound                       = "InvalidElasticIpID.NotFound"
	GatewayNotFound                   = "InvalidGatewayID.NotFound"
	GroupNotFound                     = "InvalidGroup.NotFound"
	InternetGatewayNotFound           = "InvalidInternetGatewayID.NotFound"
	InvalidCarrierGatewayNotFound     = "InvalidCarrierGatewayID.NotFound"
	EgressOnlyInternetGatewayNotFound = "InvalidEgressOnlyInternetGatewayID.NotFound"
	InUseIPAddress                    = "InvalidIPAddress.InUse"
	InvalidAccessKeyID                = "InvalidAccessKeyId"
	InvalidClientTokenID              = "InvalidClientTokenId"
	InvalidInstanceID                 = "InvalidInstanceID.NotFound"
	InvalidSubnet                     = "InvalidSubnet"
	LaunchTemplateNameNotFound        = "InvalidLaunchTemplateName.NotFoundException"
	LoadBalancerNotFound              = "LoadBalancerNotFound"
	NATGatewayNotFound                = "InvalidNatGatewayID.NotFound"
	//nolint:gosec
	NoCredentialProviders                   = "NoCredentialProviders"
	NoSuchKey                               = "NoSuchKey"
	PermissionNotFound                      = "InvalidPermission.NotFound"
	ResourceExists                          = "ResourceExistsException"
	ResourceNotFound                        = "InvalidResourceID.NotFound"
	RouteTableNotFound                      = "InvalidRouteTableID.NotFound"
	SubnetNotFound                          = "InvalidSubnetID.NotFound"
	UnrecognizedClientException             = "UnrecognizedClientException"
	UnauthorizedOperation                   = "UnauthorizedOperation"
	VPCNotFound                             = "InvalidVpcID.NotFound"
	VPCMissingParameter                     = "MissingParameter"
	ErrCodeRepositoryAlreadyExistsException = "RepositoryAlreadyExistsException"
	ASGNotFound                             = "AutoScalingGroup.NotFound"
)

var _ error = &EC2Error{}

// Code returns the AWS error code as a string.
func Code(err error) (string, bool) {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Code(), true
	}
	return "", false
}

// Message returns the AWS error message as a string.
func Message(err error) string {
	if awserr, ok := err.(awserr.Error); ok {
		return awserr.Message()
	}
	return ""
}

// EC2Error is an error exposed to users of this library.
type EC2Error struct {
	msg string

	Code int
}

// Error implements the Error interface.
func (e *EC2Error) Error() string {
	return e.msg
}

// NewNotFound returns an error which indicates that the resource of the kind and the name was not found.
func NewNotFound(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusNotFound,
	}
}

// NewConflict returns an error which indicates that the request cannot be processed due to a conflict.
func NewConflict(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusConflict,
	}
}

// IsBucketAlreadyOwnedByYou checks if the bucket is already owned.
func IsBucketAlreadyOwnedByYou(err error) bool {
	if code, ok := Code(err); ok {
		return code == BucketAlreadyOwnedByYou
	}
	return false
}

// IsResourceExists checks the state of the resource.
func IsResourceExists(err error) bool {
	if code, ok := Code(err); ok {
		return code == ResourceExists
	}
	return false
}

// IsRepositoryExists checks if there is already a repository with the same name.
func IsRepositoryExists(err error) bool {
	if code, ok := Code(err); ok {
		return code == ErrCodeRepositoryAlreadyExistsException
	}
	return false
}

// NewFailedDependency returns an error which indicates that a dependency failure status.
func NewFailedDependency(msg string) error {
	return &EC2Error{
		msg:  msg,
		Code: http.StatusFailedDependency,
	}
}

// IsFailedDependency checks if the error is pf http.StatusFailedDependency.
func IsFailedDependency(err error) bool {
	return ReasonForError(err) == http.StatusFailedDependency
}

// IsNotFound returns true if the error was created by NewNotFound.
func IsNotFound(err error) bool {
	if ReasonForError(err) == http.StatusNotFound {
		return true
	}
	return IsInvalidNotFoundError(err)
}

// IsConflict returns true if the error was created by NewConflict.
func IsConflict(err error) bool {
	return ReasonForError(err) == http.StatusConflict
}

// IsSDKError returns true if the error is of type awserr.Error.
func IsSDKError(err error) (ok bool) {
	_, ok = err.(awserr.Error)
	return
}

// IsInvalidNotFoundError tests for common aws not found errors.
func IsInvalidNotFoundError(err error) bool {
	if code, ok := Code(err); ok {
		switch code {
		case VPCNotFound:
			return true
		case InvalidInstanceID:
			return true
		case ssm.ErrCodeParameterNotFound:
			return true
		case LaunchTemplateNameNotFound:
			return true
		case ASGNotFound:
			return true
		}
	}

	return false
}

// IsPermissionsError tests for common aws permission errors.
func IsPermissionsError(err error) bool {
	if code, ok := Code(err); ok {
		return code == AuthFailure || code == UnauthorizedOperation
	}

	return false
}

// ReasonForError returns the HTTP status for a particular error.
func ReasonForError(err error) int {
	if t, ok := err.(*EC2Error); ok {
		return t.Code
	}

	return -1
}

// IsIgnorableSecurityGroupError checks for errors in SG that can be ignored and then return nil.
func IsIgnorableSecurityGroupError(err error) error {
	if code, ok := Code(err); ok {
		switch code {
		case GroupNotFound, PermissionNotFound:
			return nil
		default:
			return err
		}
	}
	return nil
}

// IsPermissionNotFoundError returns whether the error is InvalidPermission.NotFound.
func IsPermissionNotFoundError(err error) bool {
	if code, ok := Code(err); ok {
		switch code {
		case PermissionNotFound:
			return true
		default:
			return false
		}
	}
	return false
}
