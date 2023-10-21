package scheduledactions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ScheduledActionId{}

// ScheduledActionId is a struct representing the Resource ID for a Scheduled Action
type ScheduledActionId struct {
	ScheduledActionName string
}

// NewScheduledActionID returns a new ScheduledActionId struct
func NewScheduledActionID(scheduledActionName string) ScheduledActionId {
	return ScheduledActionId{
		ScheduledActionName: scheduledActionName,
	}
}

// ParseScheduledActionID parses 'input' into a ScheduledActionId
func ParseScheduledActionID(input string) (*ScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduledActionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduledActionId{}

	if id.ScheduledActionName, ok = parsed.Parsed["scheduledActionName"]; !ok {
		return nil, fmt.Errorf("the segment 'scheduledActionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseScheduledActionIDInsensitively parses 'input' case-insensitively into a ScheduledActionId
// note: this method should only be used for API response data and not user input
func ParseScheduledActionIDInsensitively(input string) (*ScheduledActionId, error) {
	parser := resourceids.NewParserFromResourceIdType(ScheduledActionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ScheduledActionId{}

	if id.ScheduledActionName, ok = parsed.Parsed["scheduledActionName"]; !ok {
		return nil, fmt.Errorf("the segment 'scheduledActionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateScheduledActionID checks that 'input' can be parsed as a Scheduled Action ID
func ValidateScheduledActionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseScheduledActionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Scheduled Action ID
func (id ScheduledActionId) ID() string {
	fmtString := "/providers/Microsoft.CostManagement/scheduledActions/%s"
	return fmt.Sprintf(fmtString, id.ScheduledActionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Scheduled Action ID
func (id ScheduledActionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCostManagement", "Microsoft.CostManagement", "Microsoft.CostManagement"),
		resourceids.StaticSegment("staticScheduledActions", "scheduledActions", "scheduledActions"),
		resourceids.UserSpecifiedSegment("scheduledActionName", "scheduledActionValue"),
	}
}

// String returns a human-readable description of this Scheduled Action ID
func (id ScheduledActionId) String() string {
	components := []string{
		fmt.Sprintf("Scheduled Action Name: %q", id.ScheduledActionName),
	}
	return fmt.Sprintf("Scheduled Action (%s)", strings.Join(components, "\n"))
}
