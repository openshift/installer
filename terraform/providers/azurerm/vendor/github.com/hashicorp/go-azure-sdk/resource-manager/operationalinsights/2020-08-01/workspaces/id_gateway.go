package workspaces

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = GatewayId{}

// GatewayId is a struct representing the Resource ID for a Gateway
type GatewayId struct {
	SubscriptionId    string
	ResourceGroupName string
	WorkspaceName     string
	GatewayId         string
}

// NewGatewayID returns a new GatewayId struct
func NewGatewayID(subscriptionId string, resourceGroupName string, workspaceName string, gatewayId string) GatewayId {
	return GatewayId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		WorkspaceName:     workspaceName,
		GatewayId:         gatewayId,
	}
}

// ParseGatewayID parses 'input' into a GatewayId
func ParseGatewayID(input string) (*GatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(GatewayId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.GatewayId, ok = parsed.Parsed["gatewayId"]; !ok {
		return nil, fmt.Errorf("the segment 'gatewayId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseGatewayIDInsensitively parses 'input' case-insensitively into a GatewayId
// note: this method should only be used for API response data and not user input
func ParseGatewayIDInsensitively(input string) (*GatewayId, error) {
	parser := resourceids.NewParserFromResourceIdType(GatewayId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := GatewayId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.WorkspaceName, ok = parsed.Parsed["workspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'workspaceName' was not found in the resource id %q", input)
	}

	if id.GatewayId, ok = parsed.Parsed["gatewayId"]; !ok {
		return nil, fmt.Errorf("the segment 'gatewayId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateGatewayID checks that 'input' can be parsed as a Gateway ID
func ValidateGatewayID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseGatewayID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Gateway ID
func (id GatewayId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.OperationalInsights/workspaces/%s/gateways/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.WorkspaceName, id.GatewayId)
}

// Segments returns a slice of Resource ID Segments which comprise this Gateway ID
func (id GatewayId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftOperationalInsights", "Microsoft.OperationalInsights", "Microsoft.OperationalInsights"),
		resourceids.StaticSegment("staticWorkspaces", "workspaces", "workspaces"),
		resourceids.UserSpecifiedSegment("workspaceName", "workspaceValue"),
		resourceids.StaticSegment("staticGateways", "gateways", "gateways"),
		resourceids.UserSpecifiedSegment("gatewayId", "gatewayIdValue"),
	}
}

// String returns a human-readable description of this Gateway ID
func (id GatewayId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Workspace Name: %q", id.WorkspaceName),
		fmt.Sprintf("Gateway: %q", id.GatewayId),
	}
	return fmt.Sprintf("Gateway (%s)", strings.Join(components, "\n"))
}
