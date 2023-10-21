package fleetmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FleetId{}

// FleetId is a struct representing the Resource ID for a Fleet
type FleetId struct {
	SubscriptionId    string
	ResourceGroupName string
	FleetName         string
}

// NewFleetID returns a new FleetId struct
func NewFleetID(subscriptionId string, resourceGroupName string, fleetName string) FleetId {
	return FleetId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FleetName:         fleetName,
	}
}

// ParseFleetID parses 'input' into a FleetId
func ParseFleetID(input string) (*FleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(FleetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FleetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FleetName, ok = parsed.Parsed["fleetName"]; !ok {
		return nil, fmt.Errorf("the segment 'fleetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseFleetIDInsensitively parses 'input' case-insensitively into a FleetId
// note: this method should only be used for API response data and not user input
func ParseFleetIDInsensitively(input string) (*FleetId, error) {
	parser := resourceids.NewParserFromResourceIdType(FleetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FleetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FleetName, ok = parsed.Parsed["fleetName"]; !ok {
		return nil, fmt.Errorf("the segment 'fleetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateFleetID checks that 'input' can be parsed as a Fleet ID
func ValidateFleetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFleetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Fleet ID
func (id FleetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Fleet ID
func (id FleetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetValue"),
	}
}

// String returns a human-readable description of this Fleet ID
func (id FleetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
	}
	return fmt.Sprintf("Fleet (%s)", strings.Join(components, "\n"))
}
