package datacollectionruleassociations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DataCollectionRuleId{}

// DataCollectionRuleId is a struct representing the Resource ID for a Data Collection Rule
type DataCollectionRuleId struct {
	SubscriptionId         string
	ResourceGroupName      string
	DataCollectionRuleName string
}

// NewDataCollectionRuleID returns a new DataCollectionRuleId struct
func NewDataCollectionRuleID(subscriptionId string, resourceGroupName string, dataCollectionRuleName string) DataCollectionRuleId {
	return DataCollectionRuleId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		DataCollectionRuleName: dataCollectionRuleName,
	}
}

// ParseDataCollectionRuleID parses 'input' into a DataCollectionRuleId
func ParseDataCollectionRuleID(input string) (*DataCollectionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataCollectionRuleId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataCollectionRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DataCollectionRuleName, ok = parsed.Parsed["dataCollectionRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataCollectionRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDataCollectionRuleIDInsensitively parses 'input' case-insensitively into a DataCollectionRuleId
// note: this method should only be used for API response data and not user input
func ParseDataCollectionRuleIDInsensitively(input string) (*DataCollectionRuleId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataCollectionRuleId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataCollectionRuleId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DataCollectionRuleName, ok = parsed.Parsed["dataCollectionRuleName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataCollectionRuleName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDataCollectionRuleID checks that 'input' can be parsed as a Data Collection Rule ID
func ValidateDataCollectionRuleID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataCollectionRuleID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Collection Rule ID
func (id DataCollectionRuleId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/dataCollectionRules/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DataCollectionRuleName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Collection Rule ID
func (id DataCollectionRuleId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticDataCollectionRules", "dataCollectionRules", "dataCollectionRules"),
		resourceids.UserSpecifiedSegment("dataCollectionRuleName", "dataCollectionRuleValue"),
	}
}

// String returns a human-readable description of this Data Collection Rule ID
func (id DataCollectionRuleId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Data Collection Rule Name: %q", id.DataCollectionRuleName),
	}
	return fmt.Sprintf("Data Collection Rule (%s)", strings.Join(components, "\n"))
}
