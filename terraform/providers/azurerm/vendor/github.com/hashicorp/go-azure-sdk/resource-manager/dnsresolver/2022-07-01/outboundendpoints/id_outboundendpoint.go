package outboundendpoints

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = OutboundEndpointId{}

// OutboundEndpointId is a struct representing the Resource ID for a Outbound Endpoint
type OutboundEndpointId struct {
	SubscriptionId       string
	ResourceGroupName    string
	DnsResolverName      string
	OutboundEndpointName string
}

// NewOutboundEndpointID returns a new OutboundEndpointId struct
func NewOutboundEndpointID(subscriptionId string, resourceGroupName string, dnsResolverName string, outboundEndpointName string) OutboundEndpointId {
	return OutboundEndpointId{
		SubscriptionId:       subscriptionId,
		ResourceGroupName:    resourceGroupName,
		DnsResolverName:      dnsResolverName,
		OutboundEndpointName: outboundEndpointName,
	}
}

// ParseOutboundEndpointID parses 'input' into a OutboundEndpointId
func ParseOutboundEndpointID(input string) (*OutboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutboundEndpointId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutboundEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DnsResolverName, ok = parsed.Parsed["dnsResolverName"]; !ok {
		return nil, fmt.Errorf("the segment 'dnsResolverName' was not found in the resource id %q", input)
	}

	if id.OutboundEndpointName, ok = parsed.Parsed["outboundEndpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'outboundEndpointName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseOutboundEndpointIDInsensitively parses 'input' case-insensitively into a OutboundEndpointId
// note: this method should only be used for API response data and not user input
func ParseOutboundEndpointIDInsensitively(input string) (*OutboundEndpointId, error) {
	parser := resourceids.NewParserFromResourceIdType(OutboundEndpointId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := OutboundEndpointId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DnsResolverName, ok = parsed.Parsed["dnsResolverName"]; !ok {
		return nil, fmt.Errorf("the segment 'dnsResolverName' was not found in the resource id %q", input)
	}

	if id.OutboundEndpointName, ok = parsed.Parsed["outboundEndpointName"]; !ok {
		return nil, fmt.Errorf("the segment 'outboundEndpointName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateOutboundEndpointID checks that 'input' can be parsed as a Outbound Endpoint ID
func ValidateOutboundEndpointID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseOutboundEndpointID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Outbound Endpoint ID
func (id OutboundEndpointId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsResolvers/%s/outboundEndpoints/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DnsResolverName, id.OutboundEndpointName)
}

// Segments returns a slice of Resource ID Segments which comprise this Outbound Endpoint ID
func (id OutboundEndpointId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsResolvers", "dnsResolvers", "dnsResolvers"),
		resourceids.UserSpecifiedSegment("dnsResolverName", "dnsResolverValue"),
		resourceids.StaticSegment("staticOutboundEndpoints", "outboundEndpoints", "outboundEndpoints"),
		resourceids.UserSpecifiedSegment("outboundEndpointName", "outboundEndpointValue"),
	}
}

// String returns a human-readable description of this Outbound Endpoint ID
func (id OutboundEndpointId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Dns Resolver Name: %q", id.DnsResolverName),
		fmt.Sprintf("Outbound Endpoint Name: %q", id.OutboundEndpointName),
	}
	return fmt.Sprintf("Outbound Endpoint (%s)", strings.Join(components, "\n"))
}
