package containerappsrevisions

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ContainerAppId{}

// ContainerAppId is a struct representing the Resource ID for a Container App
type ContainerAppId struct {
	SubscriptionId    string
	ResourceGroupName string
	ContainerAppName  string
}

// NewContainerAppID returns a new ContainerAppId struct
func NewContainerAppID(subscriptionId string, resourceGroupName string, containerAppName string) ContainerAppId {
	return ContainerAppId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ContainerAppName:  containerAppName,
	}
}

// ParseContainerAppID parses 'input' into a ContainerAppId
func ParseContainerAppID(input string) (*ContainerAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContainerAppId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContainerAppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ContainerAppName, ok = parsed.Parsed["containerAppName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerAppName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseContainerAppIDInsensitively parses 'input' case-insensitively into a ContainerAppId
// note: this method should only be used for API response data and not user input
func ParseContainerAppIDInsensitively(input string) (*ContainerAppId, error) {
	parser := resourceids.NewParserFromResourceIdType(ContainerAppId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ContainerAppId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ContainerAppName, ok = parsed.Parsed["containerAppName"]; !ok {
		return nil, fmt.Errorf("the segment 'containerAppName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateContainerAppID checks that 'input' can be parsed as a Container App ID
func ValidateContainerAppID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseContainerAppID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Container App ID
func (id ContainerAppId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/containerApps/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ContainerAppName)
}

// Segments returns a slice of Resource ID Segments which comprise this Container App ID
func (id ContainerAppId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticContainerApps", "containerApps", "containerApps"),
		resourceids.UserSpecifiedSegment("containerAppName", "containerAppValue"),
	}
}

// String returns a human-readable description of this Container App ID
func (id ContainerAppId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Container App Name: %q", id.ContainerAppName),
	}
	return fmt.Sprintf("Container App (%s)", strings.Join(components, "\n"))
}
