package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProvisioningServiceId{}

// ProvisioningServiceId is a struct representing the Resource ID for a Provisioning Service
type ProvisioningServiceId struct {
	SubscriptionId          string
	ResourceGroupName       string
	ProvisioningServiceName string
}

// NewProvisioningServiceID returns a new ProvisioningServiceId struct
func NewProvisioningServiceID(subscriptionId string, resourceGroupName string, provisioningServiceName string) ProvisioningServiceId {
	return ProvisioningServiceId{
		SubscriptionId:          subscriptionId,
		ResourceGroupName:       resourceGroupName,
		ProvisioningServiceName: provisioningServiceName,
	}
}

// ParseProvisioningServiceID parses 'input' into a ProvisioningServiceId
func ParseProvisioningServiceID(input string) (*ProvisioningServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProvisioningServiceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProvisioningServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProvisioningServiceName, ok = parsed.Parsed["provisioningServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'provisioningServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProvisioningServiceIDInsensitively parses 'input' case-insensitively into a ProvisioningServiceId
// note: this method should only be used for API response data and not user input
func ParseProvisioningServiceIDInsensitively(input string) (*ProvisioningServiceId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProvisioningServiceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProvisioningServiceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ProvisioningServiceName, ok = parsed.Parsed["provisioningServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'provisioningServiceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProvisioningServiceID checks that 'input' can be parsed as a Provisioning Service ID
func ValidateProvisioningServiceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProvisioningServiceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Provisioning Service ID
func (id ProvisioningServiceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Devices/provisioningServices/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ProvisioningServiceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Provisioning Service ID
func (id ProvisioningServiceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevices", "Microsoft.Devices", "Microsoft.Devices"),
		resourceids.StaticSegment("staticProvisioningServices", "provisioningServices", "provisioningServices"),
		resourceids.UserSpecifiedSegment("provisioningServiceName", "provisioningServiceValue"),
	}
}

// String returns a human-readable description of this Provisioning Service ID
func (id ProvisioningServiceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Provisioning Service Name: %q", id.ProvisioningServiceName),
	}
	return fmt.Sprintf("Provisioning Service (%s)", strings.Join(components, "\n"))
}
