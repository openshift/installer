package inputs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = InputId{}

// InputId is a struct representing the Resource ID for a Input
type InputId struct {
	SubscriptionId    string
	ResourceGroupName string
	StreamingJobName  string
	InputName         string
}

// NewInputID returns a new InputId struct
func NewInputID(subscriptionId string, resourceGroupName string, streamingJobName string, inputName string) InputId {
	return InputId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StreamingJobName:  streamingJobName,
		InputName:         inputName,
	}
}

// ParseInputID parses 'input' into a InputId
func ParseInputID(input string) (*InputId, error) {
	parser := resourceids.NewParserFromResourceIdType(InputId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InputId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	if id.InputName, ok = parsed.Parsed["inputName"]; !ok {
		return nil, fmt.Errorf("the segment 'inputName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseInputIDInsensitively parses 'input' case-insensitively into a InputId
// note: this method should only be used for API response data and not user input
func ParseInputIDInsensitively(input string) (*InputId, error) {
	parser := resourceids.NewParserFromResourceIdType(InputId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := InputId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	if id.InputName, ok = parsed.Parsed["inputName"]; !ok {
		return nil, fmt.Errorf("the segment 'inputName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateInputID checks that 'input' can be parsed as a Input ID
func ValidateInputID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseInputID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Input ID
func (id InputId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s/inputs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName, id.InputName)
}

// Segments returns a slice of Resource ID Segments which comprise this Input ID
func (id InputId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStreamAnalytics", "Microsoft.StreamAnalytics", "Microsoft.StreamAnalytics"),
		resourceids.StaticSegment("staticStreamingJobs", "streamingJobs", "streamingJobs"),
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobValue"),
		resourceids.StaticSegment("staticInputs", "inputs", "inputs"),
		resourceids.UserSpecifiedSegment("inputName", "inputValue"),
	}
}

// String returns a human-readable description of this Input ID
func (id InputId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Streaming Job Name: %q", id.StreamingJobName),
		fmt.Sprintf("Input Name: %q", id.InputName),
	}
	return fmt.Sprintf("Input (%s)", strings.Join(components, "\n"))
}
