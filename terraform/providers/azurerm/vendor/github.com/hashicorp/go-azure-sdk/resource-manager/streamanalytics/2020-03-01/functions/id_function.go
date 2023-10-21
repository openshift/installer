package functions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = FunctionId{}

// FunctionId is a struct representing the Resource ID for a Function
type FunctionId struct {
	SubscriptionId    string
	ResourceGroupName string
	StreamingJobName  string
	FunctionName      string
}

// NewFunctionID returns a new FunctionId struct
func NewFunctionID(subscriptionId string, resourceGroupName string, streamingJobName string, functionName string) FunctionId {
	return FunctionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StreamingJobName:  streamingJobName,
		FunctionName:      functionName,
	}
}

// ParseFunctionID parses 'input' into a FunctionId
func ParseFunctionID(input string) (*FunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(FunctionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FunctionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	if id.FunctionName, ok = parsed.Parsed["functionName"]; !ok {
		return nil, fmt.Errorf("the segment 'functionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseFunctionIDInsensitively parses 'input' case-insensitively into a FunctionId
// note: this method should only be used for API response data and not user input
func ParseFunctionIDInsensitively(input string) (*FunctionId, error) {
	parser := resourceids.NewParserFromResourceIdType(FunctionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := FunctionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	if id.FunctionName, ok = parsed.Parsed["functionName"]; !ok {
		return nil, fmt.Errorf("the segment 'functionName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateFunctionID checks that 'input' can be parsed as a Function ID
func ValidateFunctionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseFunctionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Function ID
func (id FunctionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s/functions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName, id.FunctionName)
}

// Segments returns a slice of Resource ID Segments which comprise this Function ID
func (id FunctionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStreamAnalytics", "Microsoft.StreamAnalytics", "Microsoft.StreamAnalytics"),
		resourceids.StaticSegment("staticStreamingJobs", "streamingJobs", "streamingJobs"),
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobValue"),
		resourceids.StaticSegment("staticFunctions", "functions", "functions"),
		resourceids.UserSpecifiedSegment("functionName", "functionValue"),
	}
}

// String returns a human-readable description of this Function ID
func (id FunctionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Streaming Job Name: %q", id.StreamingJobName),
		fmt.Sprintf("Function Name: %q", id.FunctionName),
	}
	return fmt.Sprintf("Function (%s)", strings.Join(components, "\n"))
}
