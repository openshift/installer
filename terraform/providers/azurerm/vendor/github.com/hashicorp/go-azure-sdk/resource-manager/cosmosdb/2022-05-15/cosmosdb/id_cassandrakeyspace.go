package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CassandraKeyspaceId{}

// CassandraKeyspaceId is a struct representing the Resource ID for a Cassandra Keyspace
type CassandraKeyspaceId struct {
	SubscriptionId        string
	ResourceGroupName     string
	DatabaseAccountName   string
	CassandraKeyspaceName string
}

// NewCassandraKeyspaceID returns a new CassandraKeyspaceId struct
func NewCassandraKeyspaceID(subscriptionId string, resourceGroupName string, databaseAccountName string, cassandraKeyspaceName string) CassandraKeyspaceId {
	return CassandraKeyspaceId{
		SubscriptionId:        subscriptionId,
		ResourceGroupName:     resourceGroupName,
		DatabaseAccountName:   databaseAccountName,
		CassandraKeyspaceName: cassandraKeyspaceName,
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

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseAccountName' was not found in the resource id %q", input)
	}

	if id.CassandraKeyspaceName, ok = parsed.Parsed["cassandraKeyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'cassandraKeyspaceName' was not found in the resource id %q", input)
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

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseAccountName' was not found in the resource id %q", input)
	}

	if id.CassandraKeyspaceName, ok = parsed.Parsed["cassandraKeyspaceName"]; !ok {
		return nil, fmt.Errorf("the segment 'cassandraKeyspaceName' was not found in the resource id %q", input)
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
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.CassandraKeyspaceName)
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
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticCassandraKeyspaces", "cassandraKeyspaces", "cassandraKeyspaces"),
		resourceids.UserSpecifiedSegment("cassandraKeyspaceName", "cassandraKeyspaceValue"),
	}
}

// String returns a human-readable description of this Cassandra Keyspace ID
func (id CassandraKeyspaceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Cassandra Keyspace Name: %q", id.CassandraKeyspaceName),
	}
	return fmt.Sprintf("Cassandra Keyspace (%s)", strings.Join(components, "\n"))
}
