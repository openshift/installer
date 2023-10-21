package administrators

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FlexibleServerId{}

// FlexibleServerId is a struct representing the Resource ID for a Flexible Server
type FlexibleServerId struct {
	SubscriptionId     string
	ResourceGroupName  string
	FlexibleServerName string
}

// NewFlexibleServerID returns a new FlexibleServerId struct
func NewFlexibleServerID(subscriptionId string, resourceGroupName string, flexibleServerName string) FlexibleServerId {
	return FlexibleServerId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		FlexibleServerName: flexibleServerName,
	}
}

// ParseFlexibleServerID parses 'input' into a FlexibleServerId
func ParseFlexibleServerID(input string) (*FlexibleServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FlexibleServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FlexibleServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FlexibleServerName, ok = parsed.Parsed["flexibleServerName"]; !ok {
		return nil, fmt.Errorf("the segment 'flexibleServerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseFlexibleServerIDInsensitively parses 'input' case-insensitively into a FlexibleServerId
// note: this method should only be used for API response data and not user input
func ParseFlexibleServerIDInsensitively(input string) (*FlexibleServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(FlexibleServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FlexibleServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FlexibleServerName, ok = parsed.Parsed["flexibleServerName"]; !ok {
		return nil, fmt.Errorf("the segment 'flexibleServerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateFlexibleServerID checks that 'input' can be parsed as a Flexible Server ID
func ValidateFlexibleServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFlexibleServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Flexible Server ID
func (id FlexibleServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/flexibleServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FlexibleServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Flexible Server ID
func (id FlexibleServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticFlexibleServers", "flexibleServers", "flexibleServers"),
		resourceids.UserSpecifiedSegment("flexibleServerName", "flexibleServerValue"),
	}
}

// String returns a human-readable description of this Flexible Server ID
func (id FlexibleServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Flexible Server Name: %q", id.FlexibleServerName),
	}
	return fmt.Sprintf("Flexible Server (%s)", strings.Join(components, "\n"))
}
