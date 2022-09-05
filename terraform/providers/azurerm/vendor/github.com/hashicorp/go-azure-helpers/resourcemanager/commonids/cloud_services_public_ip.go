package commonids

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CloudServicesPublicIPAddressId{}

// CloudServicesPublicIPAddressId is a struct representing the Resource ID for a Cloud Services Public I P Address
type CloudServicesPublicIPAddressId struct {
	SubscriptionId       string
	ResourceGroup        string
	CloudServiceName     string
	RoleInstanceName     string
	NetworkInterfaceName string
	IpConfigurationName  string
	PublicIPAddressName  string
}

// NewCloudServicesPublicIPAddressID returns a new CloudServicesPublicIPAddressId struct
func NewCloudServicesPublicIPAddressID(subscriptionId string, resourceGroup string, cloudServiceName string, roleInstanceName string, networkInterfaceName string, ipConfigurationName string, publicIPAddressName string) CloudServicesPublicIPAddressId {
	return CloudServicesPublicIPAddressId{
		SubscriptionId:       subscriptionId,
		ResourceGroup:        resourceGroup,
		CloudServiceName:     cloudServiceName,
		RoleInstanceName:     roleInstanceName,
		NetworkInterfaceName: networkInterfaceName,
		IpConfigurationName:  ipConfigurationName,
		PublicIPAddressName:  publicIPAddressName,
	}
}

// ParseCloudServicesPublicIPAddressID parses 'input' into a CloudServicesPublicIPAddressId
func ParseCloudServicesPublicIPAddressID(input string) (*CloudServicesPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(CloudServicesPublicIPAddressId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CloudServicesPublicIPAddressId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'cloudServiceName' was not found in the resource id %q", input)
	}

	if id.RoleInstanceName, ok = parsed.Parsed["roleInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'roleInstanceName' was not found in the resource id %q", input)
	}

	if id.NetworkInterfaceName, ok = parsed.Parsed["networkInterfaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'networkInterfaceName' was not found in the resource id %q", input)
	}

	if id.IpConfigurationName, ok = parsed.Parsed["ipConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipConfigurationName' was not found in the resource id %q", input)
	}

	if id.PublicIPAddressName, ok = parsed.Parsed["publicIPAddressName"]; !ok {
		return nil, fmt.Errorf("the segment 'publicIPAddressName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCloudServicesPublicIPAddressIDInsensitively parses 'input' case-insensitively into a CloudServicesPublicIPAddressId
// note: this method should only be used for API response data and not user input
func ParseCloudServicesPublicIPAddressIDInsensitively(input string) (*CloudServicesPublicIPAddressId, error) {
	parser := resourceids.NewParserFromResourceIdType(CloudServicesPublicIPAddressId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CloudServicesPublicIPAddressId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroup, ok = parsed.Parsed["resourceGroup"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroup' was not found in the resource id %q", input)
	}

	if id.CloudServiceName, ok = parsed.Parsed["cloudServiceName"]; !ok {
		return nil, fmt.Errorf("the segment 'cloudServiceName' was not found in the resource id %q", input)
	}

	if id.RoleInstanceName, ok = parsed.Parsed["roleInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'roleInstanceName' was not found in the resource id %q", input)
	}

	if id.NetworkInterfaceName, ok = parsed.Parsed["networkInterfaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'networkInterfaceName' was not found in the resource id %q", input)
	}

	if id.IpConfigurationName, ok = parsed.Parsed["ipConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'ipConfigurationName' was not found in the resource id %q", input)
	}

	if id.PublicIPAddressName, ok = parsed.Parsed["publicIPAddressName"]; !ok {
		return nil, fmt.Errorf("the segment 'publicIPAddressName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCloudServicesPublicIPAddressID checks that 'input' can be parsed as a Cloud Services Public I P Address ID
func ValidateCloudServicesPublicIPAddressID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCloudServicesPublicIPAddressID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cloud Services Public I P Address ID
func (id CloudServicesPublicIPAddressId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/cloudServices/%s/roleInstances/%s/networkInterfaces/%s/ipConfigurations/%s/publicIPAddresses/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroup, id.CloudServiceName, id.RoleInstanceName, id.NetworkInterfaceName, id.IpConfigurationName, id.PublicIPAddressName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cloud Services Public I P Address ID
func (id CloudServicesPublicIPAddressId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("subscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("resourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroup", "example-resource-group"),
		resourceids.StaticSegment("providers", "providers", "providers"),
		resourceids.ResourceProviderSegment("resourceProvider", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("cloudServices", "cloudServices", "cloudServices"),
		resourceids.UserSpecifiedSegment("cloudServiceName", "cloudServiceValue"),
		resourceids.StaticSegment("roleInstances", "roleInstances", "roleInstances"),
		resourceids.UserSpecifiedSegment("roleInstanceName", "roleInstanceValue"),
		resourceids.StaticSegment("networkInterfaces", "networkInterfaces", "networkInterfaces"),
		resourceids.UserSpecifiedSegment("networkInterfaceName", "networkInterfaceValue"),
		resourceids.StaticSegment("ipConfigurations", "ipConfigurations", "ipConfigurations"),
		resourceids.UserSpecifiedSegment("ipConfigurationName", "ipConfigurationValue"),
		resourceids.StaticSegment("publicIPAddresses", "publicIPAddresses", "publicIPAddresses"),
		resourceids.UserSpecifiedSegment("publicIPAddressName", "publicIPAddressValue"),
	}
}

// String returns a human-readable description of this Cloud Services Public I P Address ID
func (id CloudServicesPublicIPAddressId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group: %q", id.ResourceGroup),
		fmt.Sprintf("Cloud Service Name: %q", id.CloudServiceName),
		fmt.Sprintf("Role Instance Name: %q", id.RoleInstanceName),
		fmt.Sprintf("Network Interface Name: %q", id.NetworkInterfaceName),
		fmt.Sprintf("Ip Configuration Name: %q", id.IpConfigurationName),
		fmt.Sprintf("Public I P Address Name: %q", id.PublicIPAddressName),
	}
	return fmt.Sprintf("Cloud Services Public I P Address (%s)", strings.Join(components, "\n"))
}
