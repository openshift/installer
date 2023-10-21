package redis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = LinkedServerId{}

// LinkedServerId is a struct representing the Resource ID for a Linked Server
type LinkedServerId struct {
	SubscriptionId    string
	ResourceGroupName string
	RedisName         string
	LinkedServerName  string
}

// NewLinkedServerID returns a new LinkedServerId struct
func NewLinkedServerID(subscriptionId string, resourceGroupName string, redisName string, linkedServerName string) LinkedServerId {
	return LinkedServerId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		RedisName:         redisName,
		LinkedServerName:  linkedServerName,
	}
}

// ParseLinkedServerID parses 'input' into a LinkedServerId
func ParseLinkedServerID(input string) (*LinkedServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LinkedServerId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LinkedServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, fmt.Errorf("the segment 'redisName' was not found in the resource id %q", input)
	}

	if id.LinkedServerName, ok = parsed.Parsed["linkedServerName"]; !ok {
		return nil, fmt.Errorf("the segment 'linkedServerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseLinkedServerIDInsensitively parses 'input' case-insensitively into a LinkedServerId
// note: this method should only be used for API response data and not user input
func ParseLinkedServerIDInsensitively(input string) (*LinkedServerId, error) {
	parser := resourceids.NewParserFromResourceIdType(LinkedServerId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := LinkedServerId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.RedisName, ok = parsed.Parsed["redisName"]; !ok {
		return nil, fmt.Errorf("the segment 'redisName' was not found in the resource id %q", input)
	}

	if id.LinkedServerName, ok = parsed.Parsed["linkedServerName"]; !ok {
		return nil, fmt.Errorf("the segment 'linkedServerName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateLinkedServerID checks that 'input' can be parsed as a Linked Server ID
func ValidateLinkedServerID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseLinkedServerID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Linked Server ID
func (id LinkedServerId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Cache/redis/%s/linkedServers/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.RedisName, id.LinkedServerName)
}

// Segments returns a slice of Resource ID Segments which comprise this Linked Server ID
func (id LinkedServerId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftCache", "Microsoft.Cache", "Microsoft.Cache"),
		resourceids.StaticSegment("staticRedis", "redis", "redis"),
		resourceids.UserSpecifiedSegment("redisName", "redisValue"),
		resourceids.StaticSegment("staticLinkedServers", "linkedServers", "linkedServers"),
		resourceids.UserSpecifiedSegment("linkedServerName", "linkedServerValue"),
	}
}

// String returns a human-readable description of this Linked Server ID
func (id LinkedServerId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Redis Name: %q", id.RedisName),
		fmt.Sprintf("Linked Server Name: %q", id.LinkedServerName),
	}
	return fmt.Sprintf("Linked Server (%s)", strings.Join(components, "\n"))
}
