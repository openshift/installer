package streamingjobs

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = StreamingJobId{}

// StreamingJobId is a struct representing the Resource ID for a Streaming Job
type StreamingJobId struct {
	SubscriptionId    string
	ResourceGroupName string
	StreamingJobName  string
}

// NewStreamingJobID returns a new StreamingJobId struct
func NewStreamingJobID(subscriptionId string, resourceGroupName string, streamingJobName string) StreamingJobId {
	return StreamingJobId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		StreamingJobName:  streamingJobName,
	}
}

// ParseStreamingJobID parses 'input' into a StreamingJobId
func ParseStreamingJobID(input string) (*StreamingJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingJobId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseStreamingJobIDInsensitively parses 'input' case-insensitively into a StreamingJobId
// note: this method should only be used for API response data and not user input
func ParseStreamingJobIDInsensitively(input string) (*StreamingJobId, error) {
	parser := resourceids.NewParserFromResourceIdType(StreamingJobId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StreamingJobId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.StreamingJobName, ok = parsed.Parsed["streamingJobName"]; !ok {
		return nil, fmt.Errorf("the segment 'streamingJobName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateStreamingJobID checks that 'input' can be parsed as a Streaming Job ID
func ValidateStreamingJobID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStreamingJobID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Streaming Job ID
func (id StreamingJobId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.StreamAnalytics/streamingJobs/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.StreamingJobName)
}

// Segments returns a slice of Resource ID Segments which comprise this Streaming Job ID
func (id StreamingJobId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftStreamAnalytics", "Microsoft.StreamAnalytics", "Microsoft.StreamAnalytics"),
		resourceids.StaticSegment("staticStreamingJobs", "streamingJobs", "streamingJobs"),
		resourceids.UserSpecifiedSegment("streamingJobName", "streamingJobValue"),
	}
}

// String returns a human-readable description of this Streaming Job ID
func (id StreamingJobId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Streaming Job Name: %q", id.StreamingJobName),
	}
	return fmt.Sprintf("Streaming Job (%s)", strings.Join(components, "\n"))
}
