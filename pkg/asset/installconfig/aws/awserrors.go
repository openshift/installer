package aws

import (
	"errors"
	"net/http"

	awshttp "github.com/aws/aws-sdk-go-v2/aws/transport/http"
	"github.com/aws/smithy-go"
)

// Error constants for AWS error codes.
const (
	// Common error codes.
	AccessDeniedException   = "AccessDeniedException"
	NoSuchResourceException = "NoSuchResourceException"
	InvalidParameter        = "InvalidParameter"
	NoSuchEntity            = "NoSuchEntity"

	// Route53 error codes.
	NoSuchHostedZone = "NoSuchHostedZone"

	// EFS error codes.
	FileSystemNotFound  = "FileSystemNotFound"
	AccessPointNotFound = "AccessPointNotFound"
	MountTargetNotFound = "MountTargetNotFound"

	// ELB and ELBv2 error codes.
	ListenerNotFound = "ListenerNotFound"

	// EC2 error codes.
	InvalidDhcpOptionsNotFound                 = "InvalidDhcpOptions.NotFound"
	InvalidAMINotFound                         = "InvalidAMI.NotFound"
	InvalidAllocationNotFound                  = "InvalidAllocation.NotFound"
	InvalidInstanceNotFound                    = "InvalidInstance.NotFound"
	GatewayNotAttached                         = "Gateway.NotAttached"
	InvalidCarrierGatewayNotFound              = "InvalidCarrierGateway.NotFound"
	NatGatewayNotFound                         = "NatGateway.NotFound"
	InvalidPlacementGroupNotFound              = "InvalidPlacementGroup.NotFound"
	InvalidRouteTableIDNotFound                = "InvalidRouteTableID.NotFound"
	InvalidGroupNotFound                       = "InvalidGroup.NotFound"
	InvalidSnapshotNotFound                    = "InvalidSnapshot.NotFound"
	InvalidNetworkInterfaceIDNotFound          = "InvalidNetworkInterfaceID.NotFound"
	InvalidSubnetIDNotFound                    = "InvalidSubnetID.NotFound"
	InvalidVolumeNotFound                      = "InvalidVolume.NotFound"
	InvalidVpcIDNotFound                       = "InvalidVpcID.NotFound"
	InvalidVpcPeeringConnectionNotFound        = "InvalidVpcPeeringConnection.NotFound"
	InvalidVpcEndpointServiceNotFound          = "InvalidVpcEndpointService.NotFound"
	InvalidEgressOnlyInternetGatewayIDNotFound = "InvalidEgressOnlyInternetGatewayId.NotFound"
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

// GetAWSErrorCode takes the error and extracts the AWS error code if it is an AWS API Error.
func GetAWSErrorCode(err error) string {
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode()
	}
	return ""
}
