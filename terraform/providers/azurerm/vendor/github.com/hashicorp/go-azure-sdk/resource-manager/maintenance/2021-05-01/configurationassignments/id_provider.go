package configurationassignments

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProviderId{}

// ProviderId is a struct representing the Resource ID for a Provider
type ProviderId struct {
	SubscriptionId    string
	ResourceGroupName string
	ProviderName      string
	ResourceType      string
	ResourceName      string
}

// NewProviderID returns a new ProviderId struct
func NewProviderID(subscriptionId string, resourceGroupName string, providerName string, resourceType string, resourceName string) ProviderId {
	return ProviderId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ProviderName:      providerName,
		ResourceType:      resourceType,
		ResourceName:      resourceName,
	}
}

// ParseProviderID parses 'input' into a ProviderId
func ParseProviderID(input string) (*ProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProviderIDInsensitively parses 'input' case-insensitively into a ProviderId
// note: this method should only be used for API response data and not user input
func ParseProviderIDInsensitively(input string) (*ProviderId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProviderId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProviderId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProviderName, ok = parsed.Parsed["providerName"]; !ok {
		return nil, fmt.Errorf("the segment 'providerName' was not found in the resource id %q", input)
	}

	if id.ResourceType, ok = parsed.Parsed["resourceType"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceType' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProviderID checks that 'input' can be parsed as a Provider ID
func ValidateProviderID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProviderID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provider ID
func (id ProviderId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/%s/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProviderName, id.ResourceType, id.ResourceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provider ID
func (id ProviderId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.UserSpecifiedSegment("providerName", "providerValue"),
		resourceids.UserSpecifiedSegment("resourceType", "resourceTypeValue"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
	}
}

// String returns a human-readable description of this Provider ID
func (id ProviderId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provider Name: %q", id.ProviderName),
		fmt.Sprintf("Resource Type: %q", id.ResourceType),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
	}
	return fmt.Sprintf("Provider (%s)", strings.Join(components, "\n"))
}
