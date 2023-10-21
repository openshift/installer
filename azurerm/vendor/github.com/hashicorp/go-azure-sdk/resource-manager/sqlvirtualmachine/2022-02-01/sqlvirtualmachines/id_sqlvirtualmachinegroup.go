package sqlvirtualmachines

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = SqlVirtualMachineGroupId{}

// SqlVirtualMachineGroupId is a struct representing the Resource ID for a Sql Virtual Machine Group
type SqlVirtualMachineGroupId struct {
	SubscriptionId             string
	ResourceGroupName          string
	SqlVirtualMachineGroupName string
}

// NewSqlVirtualMachineGroupID returns a new SqlVirtualMachineGroupId struct
func NewSqlVirtualMachineGroupID(subscriptionId string, resourceGroupName string, sqlVirtualMachineGroupName string) SqlVirtualMachineGroupId {
	return SqlVirtualMachineGroupId{
		SubscriptionId:             subscriptionId,
		ResourceGroupName:          resourceGroupName,
		SqlVirtualMachineGroupName: sqlVirtualMachineGroupName,
	}
}

// ParseSqlVirtualMachineGroupID parses 'input' into a SqlVirtualMachineGroupId
func ParseSqlVirtualMachineGroupID(input string) (*SqlVirtualMachineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlVirtualMachineGroupId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlVirtualMachineGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SqlVirtualMachineGroupName, ok = parsed.Parsed["sqlVirtualMachineGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'sqlVirtualMachineGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseSqlVirtualMachineGroupIDInsensitively parses 'input' case-insensitively into a SqlVirtualMachineGroupId
// note: this method should only be used for API response data and not user input
func ParseSqlVirtualMachineGroupIDInsensitively(input string) (*SqlVirtualMachineGroupId, error) {
	parser := resourceids.NewParserFromResourceIdType(SqlVirtualMachineGroupId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := SqlVirtualMachineGroupId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.SqlVirtualMachineGroupName, ok = parsed.Parsed["sqlVirtualMachineGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'sqlVirtualMachineGroupName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateSqlVirtualMachineGroupID checks that 'input' can be parsed as a Sql Virtual Machine Group ID
func ValidateSqlVirtualMachineGroupID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseSqlVirtualMachineGroupID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.SqlVirtualMachine/sqlVirtualMachineGroups/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.SqlVirtualMachineGroupName)
}

// Segments returns a slice of Resource ID Segments which comprise this Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftSqlVirtualMachine", "Microsoft.SqlVirtualMachine", "Microsoft.SqlVirtualMachine"),
		resourceids.StaticSegment("staticSqlVirtualMachineGroups", "sqlVirtualMachineGroups", "sqlVirtualMachineGroups"),
		resourceids.UserSpecifiedSegment("sqlVirtualMachineGroupName", "sqlVirtualMachineGroupValue"),
	}
}

// String returns a human-readable description of this Sql Virtual Machine Group ID
func (id SqlVirtualMachineGroupId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Sql Virtual Machine Group Name: %q", id.SqlVirtualMachineGroupName),
	}
	return fmt.Sprintf("Sql Virtual Machine Group (%s)", strings.Join(components, "\n"))
}
