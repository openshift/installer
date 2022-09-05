package recordsets

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ZoneId{}

// ZoneId is a struct representing the Resource ID for a Zone
type ZoneId struct {
	SubscriptionId    string
	ResourceGroupName string
	ZoneName          string
	RecordType        RecordType
}

// NewZoneID returns a new ZoneId struct
func NewZoneID(subscriptionId string, resourceGroupName string, zoneName string, recordType RecordType) ZoneId {
	return ZoneId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ZoneName:          zoneName,
		RecordType:        recordType,
	}
}

// ParseZoneID parses 'input' into a ZoneId
func ParseZoneID(input string) (*ZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(ZoneId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ZoneName, ok = parsed.Parsed["zoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'zoneName' was not found in the resource id %q", input)
	}

	if v, ok := parsed.Parsed["recordType"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'recordType' was not found in the resource id %q", input)
		}

		recordType, err := parseRecordType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.RecordType = *recordType
	}

	return &id, nil
}

// ParseZoneIDInsensitively parses 'input' case-insensitively into a ZoneId
// note: this method should only be used for API response data and not user input
func ParseZoneIDInsensitively(input string) (*ZoneId, error) {
	parser := resourceids.NewParserFromResourceIdType(ZoneId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ZoneId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ZoneName, ok = parsed.Parsed["zoneName"]; !ok {
		return nil, fmt.Errorf("the segment 'zoneName' was not found in the resource id %q", input)
	}

	if v, ok := parsed.Parsed["recordType"]; true {
		if !ok {
			return nil, fmt.Errorf("the segment 'recordType' was not found in the resource id %q", input)
		}

		recordType, err := parseRecordType(v)
		if err != nil {
			return nil, fmt.Errorf("parsing %q: %+v", v, err)
		}
		id.RecordType = *recordType
	}

	return &id, nil
}

// ValidateZoneID checks that 'input' can be parsed as a Zone ID
func ValidateZoneID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseZoneID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Zone ID
func (id ZoneId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Network/dnsZones/%s/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ZoneName, string(id.RecordType))
}

// Segments returns a slice of Resource ID Segments which comprise this Zone ID
func (id ZoneId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftNetwork", "Microsoft.Network", "Microsoft.Network"),
		resourceids.StaticSegment("staticDnsZones", "dnsZones", "dnsZones"),
		resourceids.UserSpecifiedSegment("zoneName", "zoneValue"),
		resourceids.ConstantSegment("recordType", PossibleValuesForRecordType(), "A"),
	}
}

// String returns a human-readable description of this Zone ID
func (id ZoneId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Zone Name: %q", id.ZoneName),
		fmt.Sprintf("Record Type: %q", string(id.RecordType)),
	}
	return fmt.Sprintf("Zone (%s)", strings.Join(components, "\n"))
}
