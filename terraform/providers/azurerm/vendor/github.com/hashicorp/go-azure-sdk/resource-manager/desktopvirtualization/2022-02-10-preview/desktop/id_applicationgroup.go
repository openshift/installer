package desktop

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ApplicationGroupId{}

// ApplicationGroupId is a struct representing the Resource ID for a Application Group
type ApplicationGroupId struct {
	SubscriptionId       string
	ResourceGroupName    string
	ApplicationGroupName string
}

// NewApplicationGroupID returns a new ApplicationGroupId struct
func NewApplicationGroupID(subscriptionId string, resourceGroupName string, applicationGroupName string) ApplicationGroupId {
	return ApplicationGroupId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		ApplicationGroupName: applicationGroupName,
	}
}

// ParseApplicationGroupID parses 'input' into a ApplicationGroupId
func ParseApplicationGroupID(input string) (*ApplicationGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseApplicationGroupIDInsensitively parses 'input' case-insensitively into a ApplicationGroupId
// note: this method should only be used for API response data and not user input
func ParseApplicationGroupIDInsensitively(input string) (*ApplicationGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(ApplicationGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ApplicationGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ApplicationGroupName, ok = parsed.Parsed["applicationGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'applicationGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateApplicationGroupID checks that 'input' can be parsed as a Application Group ID
func ValidateApplicationGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseApplicationGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Application Group ID
func (id ApplicationGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DesktopVirtualization/applicationGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ApplicationGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Application Group ID
func (id ApplicationGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDesktopVirtualization", "Microsoft.DesktopVirtualization", "Microsoft.DesktopVirtualization"),
		resourceids.StaticSegment("staticApplicationGroups", "applicationGroups", "applicationGroups"),
		resourceids.UserSpecifiedSegment("applicationGroupName", "applicationGroupValue"),
	}
}

// String returns a human-readable description of this Application Group ID
func (id ApplicationGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Application Group Name: %q", id.ApplicationGroupName),
	}
	return fmt.Sprintf("Application Group (%s)", strings.Join(components, "\n"))
}
