package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VPNConnectionId{}

// VPNConnectionId is a struct representing the Resource ID for a V P N Connection
type VPNConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	GatewayName    string
	ConnectionName string
}

// NewVPNConnectionID returns a new VPNConnectionId struct
func NewVPNConnectionID(subscriptionId string, resourceGroup string, gatewayName string, connectionName string) VPNConnectionId {
	return VPNConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		GatewayName:    gatewayName,
		ConnectionName: connectionName,
	}
}

// ParseVPNConnectionID parses 'input' into a VPNConnectionId
func ParseVPNConnectionID(input string) (*VPNConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VPNConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VPNConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, fmt.Errorf("the segment 'gatewayName' was not found in the resource id %q", input)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVPNConnectionIDInsensitively parses 'input' case-insensitively into a VPNConnectionId
// note: this method should only be used for API response data and not user input
func ParseVPNConnectionIDInsensitively(input string) (*VPNConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VPNConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VPNConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.GatewayName, ok = parsed.Parsed["gatewayName"]; !ok {
		return nil, fmt.Errorf("the segment 'gatewayName' was not found in the resource id %q", input)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVPNConnectionID checks that 'input' can be parsed as a V P N Connection ID
func ValidateVPNConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVPNConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted V P N Connection ID
func (id VPNConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/vpnGateways/%s/vpnConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.GatewayName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this V P N Connection ID
func (id VPNConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("vpnGateways", "vpnGateways", "vpnGateways"),
		resourceids.UserSpecifiedSegment("gatewayName", "gatewayValue"),
		resourceids.StaticSegment("vpnConnections", "vpnConnections", "vpnConnections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionValue"),
	}
}

// String returns a human-readable description of this V P N Connection ID
func (id VPNConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Gateway Name: %q", id.GatewayName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("V P N Connection (%s)", strings.Join(components, "\n"))
}
