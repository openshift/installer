package waiter

import (
	"strconv"
	"time"

	"github.com/aws/aws-sdk-go/service/ec2"
	"github.com/hashicorp/terraform-plugin-sdk/helper/resource"
)

const (
	// Maximum amount of time to wait for a LocalGatewayRouteTableVpcAssociation to return Associated
	LocalGatewayRouteTableVpcAssociationAssociatedTimeout = 5 * time.Minute

	// Maximum amount of time to wait for a LocalGatewayRouteTableVpcAssociation to return Disassociated
	LocalGatewayRouteTableVpcAssociationDisassociatedTimeout = 5 * time.Minute

	// General timeout for EC2 resource creations to propagate
	PropagationTimeout = 2 * time.Minute

	VpcPropagationTimeout = 2 * time.Minute
)

// LocalGatewayRouteTableVpcAssociationAssociated waits for a LocalGatewayRouteTableVpcAssociation to return Associated
func LocalGatewayRouteTableVpcAssociationAssociated(conn *ec2.EC2, localGatewayRouteTableVpcAssociationID string) (*ec2.LocalGatewayRouteTableVpcAssociation, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.RouteTableAssociationStateCodeAssociating},
		Target:  []string{ec2.RouteTableAssociationStateCodeAssociated},
		Refresh: LocalGatewayRouteTableVpcAssociationState(conn, localGatewayRouteTableVpcAssociationID),
		Timeout: LocalGatewayRouteTableVpcAssociationAssociatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.LocalGatewayRouteTableVpcAssociation); ok {
		return output, err
	}

	return nil, err
}

// LocalGatewayRouteTableVpcAssociationDisassociated waits for a LocalGatewayRouteTableVpcAssociation to return Disassociated
func LocalGatewayRouteTableVpcAssociationDisassociated(conn *ec2.EC2, localGatewayRouteTableVpcAssociationID string) (*ec2.LocalGatewayRouteTableVpcAssociation, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.RouteTableAssociationStateCodeDisassociating},
		Target:  []string{ec2.RouteTableAssociationStateCodeDisassociated},
		Refresh: LocalGatewayRouteTableVpcAssociationState(conn, localGatewayRouteTableVpcAssociationID),
		Timeout: LocalGatewayRouteTableVpcAssociationAssociatedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.LocalGatewayRouteTableVpcAssociation); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnEndpointDeletedTimout = 5 * time.Minute
)

func ClientVpnEndpointDeleted(conn *ec2.EC2, id string) (*ec2.ClientVpnEndpoint, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnEndpointStatusCodeDeleting},
		Target:  []string{},
		Refresh: ClientVpnEndpointStatus(conn, id),
		Timeout: ClientVpnEndpointDeletedTimout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.ClientVpnEndpoint); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnAuthorizationRuleActiveTimeout = 1 * time.Minute

	ClientVpnAuthorizationRuleRevokedTimeout = 1 * time.Minute
)

func ClientVpnAuthorizationRuleAuthorized(conn *ec2.EC2, authorizationRuleID string) (*ec2.AuthorizationRule, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnAuthorizationRuleStatusCodeAuthorizing},
		Target:  []string{ec2.ClientVpnAuthorizationRuleStatusCodeActive},
		Refresh: ClientVpnAuthorizationRuleStatus(conn, authorizationRuleID),
		Timeout: ClientVpnAuthorizationRuleActiveTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.AuthorizationRule); ok {
		return output, err
	}

	return nil, err
}

func ClientVpnAuthorizationRuleRevoked(conn *ec2.EC2, authorizationRuleID string) (*ec2.AuthorizationRule, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnAuthorizationRuleStatusCodeRevoking},
		Target:  []string{},
		Refresh: ClientVpnAuthorizationRuleStatus(conn, authorizationRuleID),
		Timeout: ClientVpnAuthorizationRuleRevokedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.AuthorizationRule); ok {
		return output, err
	}

	return nil, err
}

const (
	ClientVpnRouteDeletedTimeout = 1 * time.Minute
)

func ClientVpnRouteDeleted(conn *ec2.EC2, routeID string) (*ec2.ClientVpnRoute, error) {
	stateConf := &resource.StateChangeConf{
		Pending: []string{ec2.ClientVpnRouteStatusCodeActive, ec2.ClientVpnRouteStatusCodeDeleting},
		Target:  []string{},
		Refresh: ClientVpnRouteStatus(conn, routeID),
		Timeout: ClientVpnRouteDeletedTimeout,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.ClientVpnRoute); ok {
		return output, err
	}

	return nil, err
}

const (
	VpcAttributePropagationTimeout = 5 * time.Minute
)

func VpcAttributeUpdated(conn *ec2.EC2, vpcID string, attribute string, expectedValue bool) (*ec2.Vpc, error) {
	stateConf := &resource.StateChangeConf{
		Target:     []string{strconv.FormatBool(expectedValue)},
		Refresh:    VpcAttribute(conn, vpcID, attribute),
		Timeout:    VpcAttributePropagationTimeout,
		Delay:      10 * time.Second,
		MinTimeout: 3 * time.Second,
	}

	outputRaw, err := stateConf.WaitForState()

	if output, ok := outputRaw.(*ec2.Vpc); ok {
		return output, err
	}

	return nil, err
}
