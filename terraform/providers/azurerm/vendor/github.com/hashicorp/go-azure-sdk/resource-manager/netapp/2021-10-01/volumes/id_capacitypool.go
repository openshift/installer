package volumes

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CapacityPoolId{}

// CapacityPoolId is a struct representing the Resource ID for a Capacity Pool
type CapacityPoolId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	PoolName          string
}

// NewCapacityPoolID returns a new CapacityPoolId struct
func NewCapacityPoolID(subscriptionId string, resourceGroupName string, accountName string, poolName string) CapacityPoolId {
	return CapacityPoolId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		PoolName:          poolName,
	}
}

// ParseCapacityPoolID parses 'input' into a CapacityPoolId
func ParseCapacityPoolID(input string) (*CapacityPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityPoolId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.PoolName, ok = parsed.Parsed["poolName"]; !ok {
		return nil, fmt.Errorf("the segment 'poolName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCapacityPoolIDInsensitively parses 'input' case-insensitively into a CapacityPoolId
// note: this method should only be used for API response data and not user input
func ParseCapacityPoolIDInsensitively(input string) (*CapacityPoolId, error) {
	parser := resourceids.NewParserFromResourceIdType(CapacityPoolId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CapacityPoolId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.PoolName, ok = parsed.Parsed["poolName"]; !ok {
		return nil, fmt.Errorf("the segment 'poolName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCapacityPoolID checks that 'input' can be parsed as a Capacity Pool ID
func ValidateCapacityPoolID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCapacityPoolID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Capacity Pool ID
func (id CapacityPoolId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.NetApp/netAppAccounts/%s/capacityPools/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.PoolName)
}

// Segments returns a slice of Resource ID Segments which comprise this Capacity Pool ID
func (id CapacityPoolId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetApp", "Microsoft.NetApp", "Microsoft.NetApp"),
		resourceids.StaticSegment("staticNetAppAccounts", "netAppAccounts", "netAppAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticCapacityPools", "capacityPools", "capacityPools"),
		resourceids.UserSpecifiedSegment("poolName", "poolValue"),
	}
}

// String returns a human-readable description of this Capacity Pool ID
func (id CapacityPoolId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Pool Name: %q", id.PoolName),
	}
	return fmt.Sprintf("Capacity Pool (%s)", strings.Join(components, "\n"))
}
