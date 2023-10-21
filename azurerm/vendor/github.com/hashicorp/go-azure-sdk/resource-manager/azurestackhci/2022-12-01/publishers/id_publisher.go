package publishers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PublisherId{}

// PublisherId is a struct representing the Resource ID for a Publisher
type PublisherId struct {
	SubscriptionId    string
	ResourceGroupName string
	ClusterName       string
	PublisherName     string
}

// NewPublisherID returns a new PublisherId struct
func NewPublisherID(subscriptionId string, resourceGroupName string, clusterName string, publisherName string) PublisherId {
	return PublisherId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ClusterName:       clusterName,
		PublisherName:     publisherName,
	}
}

// ParsePublisherID parses 'input' into a PublisherId
func ParsePublisherID(input string) (*PublisherId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublisherId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublisherId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.PublisherName, ok = parsed.Parsed["publisherName"]; !ok {
		return nil, fmt.Errorf("the segment 'publisherName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParsePublisherIDInsensitively parses 'input' case-insensitively into a PublisherId
// note: this method should only be used for API response data and not user input
func ParsePublisherIDInsensitively(input string) (*PublisherId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublisherId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublisherId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.PublisherName, ok = parsed.Parsed["publisherName"]; !ok {
		return nil, fmt.Errorf("the segment 'publisherName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidatePublisherID checks that 'input' can be parsed as a Publisher ID
func ValidatePublisherID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePublisherID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Publisher ID
func (id PublisherId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AzureStackHCI/clusters/%s/publishers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ClusterName, id.PublisherName)
}

// Segments returns a slice of Resource ID Segments which comprise this Publisher ID
func (id PublisherId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAzureStackHCI", "Microsoft.AzureStackHCI", "Microsoft.AzureStackHCI"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticPublishers", "publishers", "publishers"),
		resourceids.UserSpecifiedSegment("publisherName", "publisherValue"),
	}
}

// String returns a human-readable description of this Publisher ID
func (id PublisherId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Publisher Name: %q", id.PublisherName),
	}
	return fmt.Sprintf("Publisher (%s)", strings.Join(components, "\n"))
}
