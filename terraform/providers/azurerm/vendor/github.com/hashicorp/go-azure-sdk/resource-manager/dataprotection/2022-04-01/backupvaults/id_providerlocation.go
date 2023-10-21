package backupvaults

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProviderLocationId{}

// ProviderLocationId is a struct representing the Resource ID for a Provider Location
type ProviderLocationId struct {
	SubscriptionId    string
	ResourceGroupName string
	LocationName      string
}

// NewProviderLocationID returns a new ProviderLocationId struct
func NewProviderLocationID(subscriptionId string, resourceGroupName string, locationName string) ProviderLocationId {
	return ProviderLocationId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		LocationName:      locationName,
	}
}

// ParseProviderLocationID parses 'input' into a ProviderLocationId
func ParseProviderLocationID(input string) (*ProviderLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLocationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, fmt.Errorf("the segment 'locationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProviderLocationIDInsensitively parses 'input' case-insensitively into a ProviderLocationId
// note: this method should only be used for API response data and not user input
func ParseProviderLocationIDInsensitively(input string) (*ProviderLocationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderLocationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderLocationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.LocationName, ok = parsed.Parsed["locationName"]; !ok {
		return nil, fmt.Errorf("the segment 'locationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProviderLocationID checks that 'input' can be parsed as a Provider Location ID
func ValidateProviderLocationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderLocationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider Location ID
func (id ProviderLocationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/locations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.LocationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider Location ID
func (id ProviderLocationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticLocations", "locations", "locations"),
		resourceids.UserSpecifiedSegment("locationName", "locationValue"),
	}
}

// String returns a human-readable description of this Provider Location ID
func (id ProviderLocationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Location Name: %q", id.LocationName),
	}
	return fmt.Sprintf("Provider Location (%s)", strings.Join(components, "\n"))
}
