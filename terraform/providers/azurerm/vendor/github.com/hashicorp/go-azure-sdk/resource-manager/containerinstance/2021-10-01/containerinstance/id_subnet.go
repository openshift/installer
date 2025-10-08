package containerinstance

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SubnetId{}

// SubnetId is a struct representing the Resource ID for a Subnet
type SubnetId struct {
	SubscriptionId     string
	ResourceGroupName  string
	VirtualNetworkName string
	SubnetName         string
}

// NewSubnetID returns a new SubnetId struct
func NewSubnetID(subscriptionId string, resourceGroupName string, virtualNetworkName string, subnetName string) SubnetId {
	return SubnetId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		VirtualNetworkName: virtualNetworkName,
		SubnetName:         subnetName,
	}
}

// ParseSubnetID parses 'input' into a SubnetId
func ParseSubnetID(input string) (*SubnetId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubnetId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubnetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkName, ok = parsed.Parsed["virtualNetworkName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkName' was not found in the resource id %q", input)
	}

	if id.SubnetName, ok = parsed.Parsed["subnetName"]; !ok {
		return nil, fmt.Errorf("the segment 'subnetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSubnetIDInsensitively parses 'input' case-insensitively into a SubnetId
// note: this method should only be used for API response data and not user input
func ParseSubnetIDInsensitively(input string) (*SubnetId, error) {
	parser := resourceids.NewParserFromResourceIdType(SubnetId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SubnetId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VirtualNetworkName, ok = parsed.Parsed["virtualNetworkName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualNetworkName' was not found in the resource id %q", input)
	}

	if id.SubnetName, ok = parsed.Parsed["subnetName"]; !ok {
		return nil, fmt.Errorf("the segment 'subnetName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSubnetID checks that 'input' can be parsed as a Subnet ID
func ValidateSubnetID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSubnetID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Subnet ID
func (id SubnetId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualNetworks/%s/subnets/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VirtualNetworkName, id.SubnetName)
}

// Segments returns a slice of Resource ID Segments which comprise this Subnet ID
func (id SubnetId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticVirtualNetworks", "virtualNetworks", "virtualNetworks"),
		resourceids.UserSpecifiedSegment("virtualNetworkName", "virtualNetworkValue"),
		resourceids.StaticSegment("staticSubnets", "subnets", "subnets"),
		resourceids.UserSpecifiedSegment("subnetName", "subnetValue"),
	}
}

// String returns a human-readable description of this Subnet ID
func (id SubnetId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Virtual Network Name: %q", id.VirtualNetworkName),
		fmt.Sprintf("Subnet Name: %q", id.SubnetName),
	}
	return fmt.Sprintf("Subnet (%s)", strings.Join(components, "\n"))
}
