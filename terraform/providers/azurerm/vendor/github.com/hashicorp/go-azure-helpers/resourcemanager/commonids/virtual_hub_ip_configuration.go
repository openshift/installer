package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = VirtualHubIPConfigurationId{}

// VirtualHubIPConfigurationId is a struct representing the Resource ID for a Virtual Hub I P Configuration
type VirtualHubIPConfigurationId struct {
	SubscriptionId string
	ResourceGroup  string
	VirtualHubName string
	IpConfigName   string
}

// NewVirtualHubIPConfigurationID returns a new VirtualHubIPConfigurationId struct
func NewVirtualHubIPConfigurationID(subscriptionId string, resourceGroup string, virtualHubName string, ipConfigName string) VirtualHubIPConfigurationId {
	return VirtualHubIPConfigurationId{
		SubscriptionId: subscriptionId,
		ResourceGroup:  resourceGroup,
		VirtualHubName: virtualHubName,
		IpConfigName:   ipConfigName,
	}
}

// ParseVirtualHubIPConfigurationID parses 'input' into a VirtualHubIPConfigurationId
func ParseVirtualHubIPConfigurationID(input string) (*VirtualHubIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubIPConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubIPConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.VirtualHubName, ok = parsed.Parsed["virtualHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualHubName' was not found in the resource id %q", input)
	}

	if id.IpConfigName, ok = parsed.Parsed["ipConfigName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipConfigName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseVirtualHubIPConfigurationIDInsensitively parses 'input' case-insensitively into a VirtualHubIPConfigurationId
// note: this method should only be used for API response data and not user input
func ParseVirtualHubIPConfigurationIDInsensitively(input string) (*VirtualHubIPConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(VirtualHubIPConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := VirtualHubIPConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.VirtualHubName, ok = parsed.Parsed["virtualHubName"]; !ok {
		return nil, fmt.Errorf("the segment 'virtualHubName' was not found in the resource id %q", input)
	}

	if id.IpConfigName, ok = parsed.Parsed["ipConfigName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipConfigName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateVirtualHubIPConfigurationID checks that 'input' can be parsed as a Virtual Hub I P Configuration ID
func ValidateVirtualHubIPConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseVirtualHubIPConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Virtual Hub I P Configuration ID
func (id VirtualHubIPConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/virtualHubs/%s/ipConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.VirtualHubName, id.IpConfigName)
}

// Segments returns a slice of Resource ID Segments which comprise this Virtual Hub I P Configuration ID
func (id VirtualHubIPConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("virtualHubs", "virtualHubs", "virtualHubs"),
		resourceids.UserSpecifiedSegment("virtualHubName", "virtualHubValue"),
		resourceids.StaticSegment("ipConfigurations", "ipConfigurations", "ipConfigurations"),
		resourceids.UserSpecifiedSegment("ipConfigName", "ipConfigValue"),
	}
}

// String returns a human-readable description of this Virtual Hub I P Configuration ID
func (id VirtualHubIPConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Virtual Hub Name: %q", id.VirtualHubName),
		fmt.Sprintf("Ip Config Name: %q", id.IpConfigName),
	}
	return fmt.Sprintf("Virtual Hub I P Configuration (%s)", strings.Join(components, "\n"))
}
