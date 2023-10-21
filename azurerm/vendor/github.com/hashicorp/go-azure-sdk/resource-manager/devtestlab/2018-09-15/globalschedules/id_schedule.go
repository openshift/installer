package globalschedules

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScheduleId{}

// ScheduleId is a struct representing the Resource ID for a Schedule
type ScheduleId struct {
	SubscriptionId    string
	ResourceGroupName string
	ScheduleName      string
}

// NewScheduleID returns a new ScheduleId struct
func NewScheduleID(subscriptionId string, resourceGroupName string, scheduleName string) ScheduleId {
	return ScheduleId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ScheduleName:      scheduleName,
	}
}

// ParseScheduleID parses 'input' into a ScheduleId
func ParseScheduleID(input string) (*ScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ScheduleName, ok = parsed.Parsed["scheduleName"]; !ok {
		return nil, fmt.Errorf("the segment 'scheduleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScheduleIDInsensitively parses 'input' case-insensitively into a ScheduleId
// note: this method should only be used for API response data and not user input
func ParseScheduleIDInsensitively(input string) (*ScheduleId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ScheduleName, ok = parsed.Parsed["scheduleName"]; !ok {
		return nil, fmt.Errorf("the segment 'scheduleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScheduleID checks that 'input' can be parsed as a Schedule ID
func ValidateScheduleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScheduleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Schedule ID
func (id ScheduleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DevTestLab/schedules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ScheduleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Schedule ID
func (id ScheduleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDevTestLab", "Microsoft.DevTestLab", "Microsoft.DevTestLab"),
		resourceids.StaticSegment("staticSchedules", "schedules", "schedules"),
		resourceids.UserSpecifiedSegment("scheduleName", "scheduleValue"),
	}
}

// String returns a human-readable description of this Schedule ID
func (id ScheduleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Schedule Name: %q", id.ScheduleName),
	}
	return fmt.Sprintf("Schedule (%s)", strings.Join(components, "\n"))
}
