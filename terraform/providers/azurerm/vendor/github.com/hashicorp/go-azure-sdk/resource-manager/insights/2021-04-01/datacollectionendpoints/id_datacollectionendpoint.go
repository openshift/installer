package datacollectionendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DataCollectionEndpointId{}

// DataCollectionEndpointId is a struct representing the Resource ID for a Data Collection Endpoint
type DataCollectionEndpointId struct {
	SubscriptionId             string
	ResourceGroupName          string
	DataCollectionEndpointName string
}

// NewDataCollectionEndpointID returns a new DataCollectionEndpointId struct
func NewDataCollectionEndpointID(subscriptionId string, resourceGroupName string, dataCollectionEndpointName string) DataCollectionEndpointId {
	return DataCollectionEndpointId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		DataCollectionEndpointName: dataCollectionEndpointName,
	}
}

// ParseDataCollectionEndpointID parses 'input' into a DataCollectionEndpointId
func ParseDataCollectionEndpointID(input string) (*DataCollectionEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataCollectionEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataCollectionEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DataCollectionEndpointName, ok = parsed.Parsed["dataCollectionEndpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataCollectionEndpointName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDataCollectionEndpointIDInsensitively parses 'input' case-insensitively into a DataCollectionEndpointId
// note: this method should only be used for API response data and not user input
func ParseDataCollectionEndpointIDInsensitively(input string) (*DataCollectionEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataCollectionEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataCollectionEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DataCollectionEndpointName, ok = parsed.Parsed["dataCollectionEndpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataCollectionEndpointName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDataCollectionEndpointID checks that 'input' can be parsed as a Data Collection Endpoint ID
func ValidateDataCollectionEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataCollectionEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Collection Endpoint ID
func (id DataCollectionEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/dataCollectionEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DataCollectionEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Collection Endpoint ID
func (id DataCollectionEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDataCollectionEndpoints", "dataCollectionEndpoints", "dataCollectionEndpoints"),
		resourceids.UserSpecifiedSegment("dataCollectionEndpointName", "dataCollectionEndpointValue"),
	}
}

// String returns a human-readable description of this Data Collection Endpoint ID
func (id DataCollectionEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Data Collection Endpoint Name: %q", id.DataCollectionEndpointName),
	}
	return fmt.Sprintf("Data Collection Endpoint (%s)", strings.Join(components, "\n"))
}
