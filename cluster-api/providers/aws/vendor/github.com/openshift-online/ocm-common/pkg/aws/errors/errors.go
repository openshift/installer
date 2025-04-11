package aws

import (
	"errors"

	iamtypes "github.com/aws/aws-sdk-go-v2/service/iam/types"
	"github.com/aws/smithy-go"
)

const (
	SignatureDoesNotMatch        = "SignatureDoesNotMatch"
	InvalidClientTokenID         = "InvalidClientTokenId"
	AccessDenied                 = "AccessDenied"
	Forbidden                    = "Forbidden"
	DryRunOperation              = "DryRunOperation"
	UnauthorizedOperation        = "UnauthorizedOperation"
	AuthFailure                  = "AuthFailure"
	OptInRequired                = "OptInRequired"
	VpcLimitExceeded             = "VpcLimitExceeded"
	LimitExceeded                = "LimitExceeded"
	UnrecognizedClientException  = "UnrecognizedClientException"
	IncompleteSignature          = "IncompleteSignature"
	AccessDeniedException        = "AccessDeniedException"
	NoSuchResourceException      = "NoSuchResourceException"
	Throttling                   = "Throttling"
	SubnetNotFound               = "InvalidSubnetID.NotFound"
	VolumeTypeNotAvailableInZone = "VolumeTypeNotAvailableInZone"
	InvalidParameterValue        = "InvalidParameterValue"
	NoSuchHostedZone             = "NoSuchHostedZone"
	DependencyViolation          = "DependencyViolation"
	NoSuchEntity                 = "NoSuchEntity"
	InvalidRouteTableID          = "InvalidRouteTableID.NotFound"
	InvalidInternetGatewayID     = "InvalidInternetGatewayID.NotFound"
	InvalidVpcID                 = "InvalidVpcID.NotFound"
	InvalidAllocationID          = "InvalidAllocationID.NotFound"
	InvalidGroup                 = "InvalidGroup.NotFound"
	InvalidSubnetID              = "InvalidSubnetId.NotFound"
)

func IsErrorCode(err error, code string) bool {
	var apiErr smithy.APIError
	return errors.As(err, &apiErr) && apiErr.ErrorCode() == code
}

func IsSubnetNotFoundError(err error) bool {
	return IsErrorCode(err, SubnetNotFound)
}

func IsThrottle(err error) bool {
	return IsErrorCode(err, Throttling)
}

func IsEntityAlreadyExistsException(err error) bool {
	var entityAlreadyExists *iamtypes.EntityAlreadyExistsException
	return errors.As(err, &entityAlreadyExists)
}

func IsNoSuchEntityException(err error) bool {
	var noSuchEntity *iamtypes.NoSuchEntityException
	return errors.As(err, &noSuchEntity)
}

func IsAccessDeniedException(err error) bool {
	return IsErrorCode(err, AccessDenied)
}

func IsForbiddenException(err error) bool {
	return IsErrorCode(err, Forbidden)
}

func IsLimitExceededException(err error) bool {
	return IsErrorCode(err, LimitExceeded)
}

func IsInvalidTokenException(err error) bool {
	return IsErrorCode(err, InvalidClientTokenID)
}

func IsDeleteConfictException(err error) bool {
	var deleteConflict *iamtypes.DeleteConflictException
	return errors.As(err, &deleteConflict)
}
