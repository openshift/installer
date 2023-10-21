package integrationaccountmaps

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = MapId{}

// MapId is a struct representing the Resource ID for a Map
type MapId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	MapName                string
}

// NewMapID returns a new MapId struct
func NewMapID(subscriptionId string, resourceGroupName string, integrationAccountName string, mapName string) MapId {
	return MapId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		MapName:                mapName,
	}
}

// ParseMapID parses 'input' into a MapId
func ParseMapID(input string) (*MapId, error) {
	parser := resourceids.NewParserFromResourceIdType(MapId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'integrationAccountName' was not found in the resource id %q", input)
	}

	if id.MapName, ok = parsed.Parsed["mapName"]; !ok {
		return nil, fmt.Errorf("the segment 'mapName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseMapIDInsensitively parses 'input' case-insensitively into a MapId
// note: this method should only be used for API response data and not user input
func ParseMapIDInsensitively(input string) (*MapId, error) {
	parser := resourceids.NewParserFromResourceIdType(MapId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MapId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'integrationAccountName' was not found in the resource id %q", input)
	}

	if id.MapName, ok = parsed.Parsed["mapName"]; !ok {
		return nil, fmt.Errorf("the segment 'mapName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateMapID checks that 'input' can be parsed as a Map ID
func ValidateMapID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMapID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Map ID
func (id MapId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/maps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.MapName)
}

// Segments returns a slice of Resource ID Segments which comprise this Map ID
func (id MapId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountValue"),
		resourceids.StaticSegment("staticMaps", "maps", "maps"),
		resourceids.UserSpecifiedSegment("mapName", "mapValue"),
	}
}

// String returns a human-readable description of this Map ID
func (id MapId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Map Name: %q", id.MapName),
	}
	return fmt.Sprintf("Map (%s)", strings.Join(components, "\n"))
}
