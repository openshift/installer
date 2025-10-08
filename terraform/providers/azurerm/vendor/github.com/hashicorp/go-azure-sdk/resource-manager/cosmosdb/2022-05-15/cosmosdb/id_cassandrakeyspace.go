package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CassandraKeyspaceId{}

// CassandraKeyspaceId is a struct representing the Resource ID for a Cassandra Keyspace
type CassandraKeyspaceId struct {
	SubscriptionId    string
	ResourceGroupName string
	AccountName       string
	KeyspaceName      string
}

// NewCassandraKeyspaceID returns a new CassandraKeyspaceId struct
func NewCassandraKeyspaceID(subscriptionId string, resourceGroupName string, accountName string, keyspaceName string) CassandraKeyspaceId {
	return CassandraKeyspaceId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		AccountName:       accountName,
		KeyspaceName:      keyspaceName,
	}
}

// ParseCassandraKeyspaceID parses 'input' into a CassandraKeyspaceId
func ParseCassandraKeyspaceID(input string) (*CassandraKeyspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraKeyspaceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraKeyspaceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.KeyspaceName, ok = parsed.Parsed["keyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyspaceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCassandraKeyspaceIDInsensitively parses 'input' case-insensitively into a CassandraKeyspaceId
// note: this method should only be used for API response data and not user input
func ParseCassandraKeyspaceIDInsensitively(input string) (*CassandraKeyspaceId, error) {
	parser := resourceids.NewParserFromResourceIdType(CassandraKeyspaceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CassandraKeyspaceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.AccountName, ok = parsed.Parsed["accountName"]; !ok {
		return nil, fmt.Errorf("the segment 'accountName' was not found in the resource id %q", input)
	}

	if id.KeyspaceName, ok = parsed.Parsed["keyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyspaceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCassandraKeyspaceID checks that 'input' can be parsed as a Cassandra Keyspace ID
func ValidateCassandraKeyspaceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCassandraKeyspaceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Cassandra Keyspace ID
func (id CassandraKeyspaceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/cassandraKeyspaces/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.AccountName, id.KeyspaceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Cassandra Keyspace ID
func (id CassandraKeyspaceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("accountName", "accountValue"),
		resourceids.StaticSegment("staticCassandraKeyspaces", "cassandraKeyspaces", "cassandraKeyspaces"),
		resourceids.UserSpecifiedSegment("keyspaceName", "keyspaceValue"),
	}
}

// String returns a human-readable description of this Cassandra Keyspace ID
func (id CassandraKeyspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Account Name: %q", id.AccountName),
		fmt.Sprintf("Keyspace Name: %q", id.KeyspaceName),
	}
	return fmt.Sprintf("Cassandra Keyspace (%s)", strings.Join(components, "\n"))
}
