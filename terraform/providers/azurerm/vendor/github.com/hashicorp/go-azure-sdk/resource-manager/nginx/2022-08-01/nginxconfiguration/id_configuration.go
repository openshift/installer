package nginxconfiguration

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ConfigurationId{}

// ConfigurationId is a struct representing the Resource ID for a Configuration
type ConfigurationId struct {
	SubscriptionId      string
	ResourceGroupName   string
	NginxDeploymentName string
	ConfigurationName   string
}

// NewConfigurationID returns a new ConfigurationId struct
func NewConfigurationID(subscriptionId string, resourceGroupName string, nginxDeploymentName string, configurationName string) ConfigurationId {
	return ConfigurationId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		NginxDeploymentName: nginxDeploymentName,
		ConfigurationName:   configurationName,
	}
}

// ParseConfigurationID parses 'input' into a ConfigurationId
func ParseConfigurationID(input string) (*ConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NginxDeploymentName, ok = parsed.Parsed["nginxDeploymentName"]; !ok {
		return nil, fmt.Errorf("the segment 'nginxDeploymentName' was not found in the resource id %q", input)
	}

	if id.ConfigurationName, ok = parsed.Parsed["configurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'configurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseConfigurationIDInsensitively parses 'input' case-insensitively into a ConfigurationId
// note: this method should only be used for API response data and not user input
func ParseConfigurationIDInsensitively(input string) (*ConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(ConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.NginxDeploymentName, ok = parsed.Parsed["nginxDeploymentName"]; !ok {
		return nil, fmt.Errorf("the segment 'nginxDeploymentName' was not found in the resource id %q", input)
	}

	if id.ConfigurationName, ok = parsed.Parsed["configurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'configurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateConfigurationID checks that 'input' can be parsed as a Configuration ID
func ValidateConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Configuration ID
func (id ConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Nginx.NginxPlus/nginxDeployments/%s/configurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.NginxDeploymentName, id.ConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Configuration ID
func (id ConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticNginxNginxPlus", "Nginx.NginxPlus", "Nginx.NginxPlus"),
		resourceids.StaticSegment("staticNginxDeployments", "nginxDeployments", "nginxDeployments"),
		resourceids.UserSpecifiedSegment("nginxDeploymentName", "nginxDeploymentValue"),
		resourceids.StaticSegment("staticConfigurations", "configurations", "configurations"),
		resourceids.UserSpecifiedSegment("configurationName", "configurationValue"),
	}
}

// String returns a human-readable description of this Configuration ID
func (id ConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Nginx Deployment Name: %q", id.NginxDeploymentName),
		fmt.Sprintf("Configuration Name: %q", id.ConfigurationName),
	}
	return fmt.Sprintf("Configuration (%s)", strings.Join(components, "\n"))
}
