package dedicatedhosts

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = HostId{}

// HostId is a struct representing the Resource ID for a Host
type HostId struct {
	SubscriptionId    string
	ResourceGroupName string
	HostGroupName     string
	HostName          string
}

// NewHostID returns a new HostId struct
func NewHostID(subscriptionId string, resourceGroupName string, hostGroupName string, hostName string) HostId {
	return HostId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		HostGroupName:     hostGroupName,
		HostName:          hostName,
	}
}

// ParseHostID parses 'input' into a HostId
func ParseHostID(input string) (*HostId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.HostGroupName, ok = parsed.Parsed["hostGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostGroupName' was not found in the resource id %q", input)
	}

	if id.HostName, ok = parsed.Parsed["hostName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseHostIDInsensitively parses 'input' case-insensitively into a HostId
// note: this method should only be used for API response data and not user input
func ParseHostIDInsensitively(input string) (*HostId, error) {
	parser := resourceids.NewParserFromResourceIdType(HostId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := HostId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.HostGroupName, ok = parsed.Parsed["hostGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostGroupName' was not found in the resource id %q", input)
	}

	if id.HostName, ok = parsed.Parsed["hostName"]; !ok {
		return nil, fmt.Errorf("the segment 'hostName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateHostID checks that 'input' can be parsed as a Host ID
func ValidateHostID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseHostID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Host ID
func (id HostId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/hostGroups/%s/hosts/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.HostGroupName, id.HostName)
}

// Segments returns a slice of Resource ID Segments which comprise this Host ID
func (id HostId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCompute", "Microsoft.Compute", "Microsoft.Compute"),
		resourceids.StaticSegment("staticHostGroups", "hostGroups", "hostGroups"),
		resourceids.UserSpecifiedSegment("hostGroupName", "hostGroupValue"),
		resourceids.StaticSegment("staticHosts", "hosts", "hosts"),
		resourceids.UserSpecifiedSegment("hostName", "hostValue"),
	}
}

// String returns a human-readable description of this Host ID
func (id HostId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Host Group Name: %q", id.HostGroupName),
		fmt.Sprintf("Host Name: %q", id.HostName),
	}
	return fmt.Sprintf("Host (%s)", strings.Join(components, "\n"))
}
