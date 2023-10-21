package datastores

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = DataStoreId{}

// DataStoreId is a struct representing the Resource ID for a Data Store
type DataStoreId struct {
	SubscriptionId    string
	ResourceGroupName string
	PrivateCloudName  string
	ClusterName       string
	DataStoreName     string
}

// NewDataStoreID returns a new DataStoreId struct
func NewDataStoreID(subscriptionId string, resourceGroupName string, privateCloudName string, clusterName string, dataStoreName string) DataStoreId {
	return DataStoreId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		PrivateCloudName:  privateCloudName,
		ClusterName:       clusterName,
		DataStoreName:     dataStoreName,
	}
}

// ParseDataStoreID parses 'input' into a DataStoreId
func ParseDataStoreID(input string) (*DataStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataStoreId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateCloudName, ok = parsed.Parsed["privateCloudName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateCloudName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.DataStoreName, ok = parsed.Parsed["dataStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseDataStoreIDInsensitively parses 'input' case-insensitively into a DataStoreId
// note: this method should only be used for API response data and not user input
func ParseDataStoreIDInsensitively(input string) (*DataStoreId, error) {
	parser := resourceids.NewParserFromResourceIdType(DataStoreId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := DataStoreId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.PrivateCloudName, ok = parsed.Parsed["privateCloudName"]; !ok {
		return nil, fmt.Errorf("the segment 'privateCloudName' was not found in the resource id %q", input)
	}

	if id.ClusterName, ok = parsed.Parsed["clusterName"]; !ok {
		return nil, fmt.Errorf("the segment 'clusterName' was not found in the resource id %q", input)
	}

	if id.DataStoreName, ok = parsed.Parsed["dataStoreName"]; !ok {
		return nil, fmt.Errorf("the segment 'dataStoreName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateDataStoreID checks that 'input' can be parsed as a Data Store ID
func ValidateDataStoreID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseDataStoreID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Data Store ID
func (id DataStoreId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.AVS/privateClouds/%s/clusters/%s/dataStores/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.PrivateCloudName, id.ClusterName, id.DataStoreName)
}

// Segments returns a slice of Resource ID Segments which comprise this Data Store ID
func (id DataStoreId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftAVS", "Microsoft.AVS", "Microsoft.AVS"),
		resourceids.StaticSegment("staticPrivateClouds", "privateClouds", "privateClouds"),
		resourceids.UserSpecifiedSegment("privateCloudName", "privateCloudValue"),
		resourceids.StaticSegment("staticClusters", "clusters", "clusters"),
		resourceids.UserSpecifiedSegment("clusterName", "clusterValue"),
		resourceids.StaticSegment("staticDataStores", "dataStores", "dataStores"),
		resourceids.UserSpecifiedSegment("dataStoreName", "dataStoreValue"),
	}
}

// String returns a human-readable description of this Data Store ID
func (id DataStoreId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Private Cloud Name: %q", id.PrivateCloudName),
		fmt.Sprintf("Cluster Name: %q", id.ClusterName),
		fmt.Sprintf("Data Store Name: %q", id.DataStoreName),
	}
	return fmt.Sprintf("Data Store (%s)", strings.Join(components, "\n"))
}
