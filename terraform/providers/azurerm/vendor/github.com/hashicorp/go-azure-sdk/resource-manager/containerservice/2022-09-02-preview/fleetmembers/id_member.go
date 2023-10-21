package fleetmembers

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = MemberId{}

// MemberId is a struct representing the Resource ID for a Member
type MemberId struct {
	SubscriptionId    string
	ResourceGroupName string
	FleetName         string
	MemberName        string
}

// NewMemberID returns a new MemberId struct
func NewMemberID(subscriptionId string, resourceGroupName string, fleetName string, memberName string) MemberId {
	return MemberId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		FleetName:         fleetName,
		MemberName:        memberName,
	}
}

// ParseMemberID parses 'input' into a MemberId
func ParseMemberID(input string) (*MemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(MemberId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MemberId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FleetName, ok = parsed.Parsed["fleetName"]; !ok {
		return nil, fmt.Errorf("the segment 'fleetName' was not found in the resource id %q", input)
	}

	if id.MemberName, ok = parsed.Parsed["memberName"]; !ok {
		return nil, fmt.Errorf("the segment 'memberName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseMemberIDInsensitively parses 'input' case-insensitively into a MemberId
// note: this method should only be used for API response data and not user input
func ParseMemberIDInsensitively(input string) (*MemberId, error) {
	parser := resourceids.NewParserFromResourceIdType(MemberId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MemberId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.FleetName, ok = parsed.Parsed["fleetName"]; !ok {
		return nil, fmt.Errorf("the segment 'fleetName' was not found in the resource id %q", input)
	}

	if id.MemberName, ok = parsed.Parsed["memberName"]; !ok {
		return nil, fmt.Errorf("the segment 'memberName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateMemberID checks that 'input' can be parsed as a Member ID
func ValidateMemberID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMemberID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Member ID
func (id MemberId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.ContainerService/fleets/%s/members/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.FleetName, id.MemberName)
}

// Segments returns a slice of Resource ID Segments which comprise this Member ID
func (id MemberId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftContainerService", "Microsoft.ContainerService", "Microsoft.ContainerService"),
		resourceids.StaticSegment("staticFleets", "fleets", "fleets"),
		resourceids.UserSpecifiedSegment("fleetName", "fleetValue"),
		resourceids.StaticSegment("staticMembers", "members", "members"),
		resourceids.UserSpecifiedSegment("memberName", "memberValue"),
	}
}

// String returns a human-readable description of this Member ID
func (id MemberId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Fleet Name: %q", id.FleetName),
		fmt.Sprintf("Member Name: %q", id.MemberName),
	}
	return fmt.Sprintf("Member (%s)", strings.Join(components, "\n"))
}
