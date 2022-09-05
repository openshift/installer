package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualHubBGPConnectionId{}

// VirtualHubBGPConnectionId is a struct representing the Resource ID for a Virtual Hub B G P Connection
type VirtualHubBGPConnectionId struct {
	SubscriptionId string
	ResourceGroup  string
	HubName        string
	ConnectionName string
}

// NewVirtualHubBGPConnectionID returns a new VirtualHubBGPConnectionId struct
func NewVirtualHubBGPConnectionID(subscriptionId string, resourceGroup string, hubName string, connectionName string) VirtualHubBGPConnectionId {
	return VirtualHubBGPConnectionId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		HubName:        hubName,
		ConnectionName: connectionName,
	}
}

// ParseVirtualHubBGPConnectionID parses 'input' into a VirtualHubBGPConnectionId
func ParseVirtualHubBGPConnectionID(input string) (*VirtualHubBGPConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubBGPConnectionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubBGPConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.HubName, ok = parsed.Parsed["hubName"]; !ok {
		return nil, fmt.Errorf("the segment 'hubName' was not found in the resource id %q", input)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVirtualHubBGPConnectionIDInsensitively parses 'input' case-insensitively into a VirtualHubBGPConnectionId
// note: this method should only be used for API response data and not user input
func ParseVirtualHubBGPConnectionIDInsensitively(input string) (*VirtualHubBGPConnectionId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubBGPConnectionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubBGPConnectionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.HubName, ok = parsed.Parsed["hubName"]; !ok {
		return nil, fmt.Errorf("the segment 'hubName' was not found in the resource id %q", input)
	}

	if id.ConnectionName, ok = parsed.Parsed["connectionName"]; !ok {
		return nil, fmt.Errorf("the segment 'connectionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVirtualHubBGPConnectionID checks that 'input' can be parsed as a Virtual Hub B G P Connection ID
func ValidateVirtualHubBGPConnectionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualHubBGPConnectionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Hub B G P Connection ID
func (id VirtualHubBGPConnectionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/bgpConnections/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.HubName, id.ConnectionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Hub B G P Connection ID
func (id VirtualHubBGPConnectionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("virtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("hubName", "hubValue"),
		resourceids.StaticSegment("bgpConnections", "bgpConnections", "bgpConnections"),
		resourceids.UserSpecifiedSegment("connectionName", "connectionValue"),
	}
}

// String returns a human-readable description of this Virtual Hub B G P Connection ID
func (id VirtualHubBGPConnectionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Hub Name: %q", id.HubName),
		fmt.Sprintf("Connection Name: %q", id.ConnectionName),
	}
	return fmt.Sprintf("Virtual Hub B G P Connection (%s)", strings.Join(components, "\n"))
}
