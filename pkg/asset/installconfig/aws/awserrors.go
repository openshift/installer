package aws

import (
	"errors"
	"net/http"
	"strings"

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
	InvalidCarrierGatewayNotFound              = "InvalidCarrierGateway.NotFound"
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
	NatGatewayNotFound                         = "NatGatewayNotFound"
	GatewayNotAttached                         = "Gateway.NotAttached"
)

// ParseAPIErrorCode returns the AWS error code from an error, if one is present.
// If the error is not an AWS API error, an empty string is returned.
func ParseAPIErrorCode(err error) string {
	if err == nil {
		return ""
	}
	var apiErr smithy.APIError
	if errors.As(err, &apiErr) {
		return apiErr.ErrorCode()
	}
	return ""
}

// IsUnauthorized checks if the error is due to lacking permissions.
func IsUnauthorized(err error) bool {
	if err == nil {
		return false
	}
	errCode := ParseAPIErrorCode(err)
	// see reference:
	// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetServiceQuota.html
	// https://docs.aws.amazon.com/servicequotas/2019-06-24/apireference/API_GetAWSDefaultServiceQuota.html
	return errCode == AccessDeniedException || errCode == NoSuchResourceException
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

// IsNoSuchEntity returns true if the error is an AWS API error and the error code indicates
// that the AWS API request references a resource entity that does not exist.
func IsNoSuchEntity(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), NoSuchEntity)
}

// IsInvalidParameter returns true if the error is an AWS API error and the error code indicates
// that a parameter specified in the AWS API request is not valid, is unsupported, or cannot be used.
func IsInvalidParameter(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), InvalidParameter)
}

// IsInvalidNotFound returns true if the error is an AWS API error and the error code indicates
// that the specified EC2 resource entity can not be found.
//
// NotFound errors generally follow the form "InvalidResource.NotFound", for example, "InvalidVpcID.NotFound".
// Except some cases, for example, "NatGatewayNotFound".
// Reference: https://docs.aws.amazon.com/AWSEC2/latest/APIReference/errors-overview.html#CommonErrors.
func IsInvalidNotFound(err error) bool {
	errCode := ParseAPIErrorCode(err)

	switch errCode {
	case InvalidDhcpOptionsNotFound,
		InvalidAMINotFound,
		InvalidAllocationNotFound,
		InvalidInstanceNotFound,
		InvalidCarrierGatewayNotFound,
		InvalidPlacementGroupNotFound,
		InvalidRouteTableIDNotFound,
		InvalidGroupNotFound,
		InvalidSnapshotNotFound,
		InvalidNetworkInterfaceIDNotFound,
		InvalidSubnetIDNotFound,
		InvalidVolumeNotFound,
		InvalidVpcIDNotFound,
		InvalidVpcPeeringConnectionNotFound,
		InvalidVpcEndpointServiceNotFound,
		InvalidEgressOnlyInternetGatewayIDNotFound:
		return true
	}

	return false
}

// IsNATGatewayNotFound returns true if the error is an AWS API error and error code indicates
// that the specified NAT gateway can not be found.
// Note: The request to delete a NAT gateway, if not found, will return error code "NatGatewayNotFound", not "InvalidNatGatewayID.NotFound".
func IsNATGatewayNotFound(err error) bool {
	return ParseAPIErrorCode(err) == NatGatewayNotFound
}

// IsGateWayNotAttached returns true if the error is an AWS API error and error code indicates
// that the specified internet gateway is not attached to a VPC.
func IsGateWayNotAttached(err error) bool {
	return ParseAPIErrorCode(err) == GatewayNotAttached
}

// IsListenerNotFound returns true if the error is an AWS API error and error code indicates
// that the specified ELB listener does not exist.
func IsListenerNotFound(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), ListenerNotFound)
}

// IsHostedZoneNotFound returns true if the error is an AWS API error and error code indicates
// that the specified Route53 hosted zone does not exist.
func IsHostedZoneNotFound(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), NoSuchHostedZone)
}

// IsFileSystemNotFound returns true if the error is an AWS API error and error code indicates
// that the specified EFS filesystem does not exist.
func IsFileSystemNotFound(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), FileSystemNotFound)
}

// IsAccessPointNotFound returns true if the error is an AWS API error and error code indicates
// that the specified EFS access point does not exist.
func IsAccessPointNotFound(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), AccessPointNotFound)
}

// IsMountTargetNotFound returns true if the error is an AWS API error and error code indicates
// that the specified EFS mount target does not exist.
func IsMountTargetNotFound(err error) bool {
	return strings.Contains(ParseAPIErrorCode(err), MountTargetNotFound)
}
