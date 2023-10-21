package managedenvironmentsstorages

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = StorageId{}

// StorageId is a struct representing the Resource ID for a Storage
type StorageId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	StorageName            string
}

// NewStorageID returns a new StorageId struct
func NewStorageID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, storageName string) StorageId {
	return StorageId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		StorageName:            storageName,
	}
}

// ParseStorageID parses 'input' into a StorageId
func ParseStorageID(input string) (*StorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	if id.StorageName, ok = parsed.Parsed["storageName"]; !ok {
		return nil, fmt.Errorf("the segment 'storageName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseStorageIDInsensitively parses 'input' case-insensitively into a StorageId
// note: this method should only be used for API response data and not user input
func ParseStorageIDInsensitively(input string) (*StorageId, error) {
	parser := resourceids.NewParserFromResourceIdType(StorageId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := StorageId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	if id.StorageName, ok = parsed.Parsed["storageName"]; !ok {
		return nil, fmt.Errorf("the segment 'storageName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateStorageID checks that 'input' can be parsed as a Storage ID
func ValidateStorageID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseStorageID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Storage ID
func (id StorageId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/storages/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.StorageName)
}

// Segments returns a slice of Resource ID Segments which comprise this Storage ID
func (id StorageId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentValue"),
		resourceids.StaticSegment("staticStorages", "storages", "storages"),
		resourceids.UserSpecifiedSegment("storageName", "storageValue"),
	}
}

// String returns a human-readable description of this Storage ID
func (id StorageId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Storage Name: %q", id.StorageName),
	}
	return fmt.Sprintf("Storage (%s)", strings.Join(components, "\n"))
}
