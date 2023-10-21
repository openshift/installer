package updateruns

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = UpdateId{}

// UpdateId is a struct representing the Resource ID for a Update
type UpdateId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	UpdateName        string
}

// NewUpdateID returns a new UpdateId struct
func NewUpdateID(subscriptionId string, resourceGroupName string, clusterName string, updateName string) UpdateId {
	return UpdateId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		UpdateName:        updateName,
	}
}

// ParseUpdateID parses 'input' into a UpdateId
func ParseUpdateID(input string) (*UpdateId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.UpdateName, ok = parsed.Parsed["updateName"]; !ok {
		return nil, fmt.Errorf("the segment 'updateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseUpdateIDInsensitively parses 'input' case-insensitively into a UpdateId
// note: this method should only be used for API response data and not user input
func ParseUpdateIDInsensitively(input string) (*UpdateId, error) {
	parser := resourceids.NewParserFromResourceIdType(UpdateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := UpdateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.UpdateName, ok = parsed.Parsed["updateName"]; !ok {
		return nil, fmt.Errorf("the segment 'updateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateUpdateID checks that 'input' can be parsed as a Update ID
func ValidateUpdateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseUpdateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Update ID
func (id UpdateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/updates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.UpdateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Update ID
func (id UpdateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticUpdates", "updates", "updates"),
		resourceids.UserSpecifiedSegment("updateName", "updateValue"),
	}
}

// String returns a human-readable description of this Update ID
func (id UpdateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Update Name: %q", id.UpdateName),
	}
	return fmt.Sprintf("Update (%s)", strings.Join(components, "\n"))
}
