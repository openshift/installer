package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ManagedEnvironmentId{}

// ManagedEnvironmentId is a struct representing the Resource ID for a Managed Environment
type ManagedEnvironmentId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
}

// NewManagedEnvironmentID returns a new ManagedEnvironmentId struct
func NewManagedEnvironmentID(subscriptionId string, resourceGroupName string, managedEnvironmentName string) ManagedEnvironmentId {
	return ManagedEnvironmentId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
	}
}

// ParseManagedEnvironmentID parses 'input' into a ManagedEnvironmentId
func ParseManagedEnvironmentID(input string) (*ManagedEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedEnvironmentId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedEnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseManagedEnvironmentIDInsensitively parses 'input' case-insensitively into a ManagedEnvironmentId
// note: this method should only be used for API response data and not user input
func ParseManagedEnvironmentIDInsensitively(input string) (*ManagedEnvironmentId, error) {
	parser := resourceids.NewParserFromResourceIdType(ManagedEnvironmentId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ManagedEnvironmentId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateManagedEnvironmentID checks that 'input' can be parsed as a Managed Environment ID
func ValidateManagedEnvironmentID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseManagedEnvironmentID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Managed Environment ID
func (id ManagedEnvironmentId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName)
}

// Segments returns a slice of Resource ID Segments which comprise this Managed Environment ID
func (id ManagedEnvironmentId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentValue"),
	}
}

// String returns a human-readable description of this Managed Environment ID
func (id ManagedEnvironmentId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
	}
	return fmt.Sprintf("Managed Environment (%s)", strings.Join(components, "\n"))
}
