package cosmosdb

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = MongodbDatabaseId{}

// MongodbDatabaseId is a struct representing the Resource ID for a Mongodb Database
type MongodbDatabaseId struct {
	SubscriptionId      string
	ResourceGroupName   string
	DatabaseAccountName string
	MongodbDatabaseName string
}

// NewMongodbDatabaseID returns a new MongodbDatabaseId struct
func NewMongodbDatabaseID(subscriptionId string, resourceGroupName string, databaseAccountName string, mongodbDatabaseName string) MongodbDatabaseId {
	return MongodbDatabaseId{
		SubscriptionId:      subscriptionId,
		ResourceGroupName:   resourceGroupName,
		DatabaseAccountName: databaseAccountName,
		MongodbDatabaseName: mongodbDatabaseName,
	}
}

// ParseMongodbDatabaseID parses 'input' into a MongodbDatabaseId
func ParseMongodbDatabaseID(input string) (*MongodbDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbDatabaseId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseAccountName' was not found in the resource id %q", input)
	}

	if id.MongodbDatabaseName, ok = parsed.Parsed["mongodbDatabaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'mongodbDatabaseName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseMongodbDatabaseIDInsensitively parses 'input' case-insensitively into a MongodbDatabaseId
// note: this method should only be used for API response data and not user input
func ParseMongodbDatabaseIDInsensitively(input string) (*MongodbDatabaseId, error) {
	parser := resourceids.NewParserFromResourceIdType(MongodbDatabaseId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := MongodbDatabaseId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.DatabaseAccountName, ok = parsed.Parsed["databaseAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'databaseAccountName' was not found in the resource id %q", input)
	}

	if id.MongodbDatabaseName, ok = parsed.Parsed["mongodbDatabaseName"]; !ok {
		return nil, fmt.Errorf("the segment 'mongodbDatabaseName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateMongodbDatabaseID checks that 'input' can be parsed as a Mongodb Database ID
func ValidateMongodbDatabaseID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseMongodbDatabaseID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Mongodb Database ID
func (id MongodbDatabaseId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DocumentDB/databaseAccounts/%s/mongodbDatabases/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.DatabaseAccountName, id.MongodbDatabaseName)
}

// Segments returns a slice of Resource ID Segments which comprise this Mongodb Database ID
func (id MongodbDatabaseId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDocumentDB", "Microsoft.DocumentDB", "Microsoft.DocumentDB"),
		resourceids.StaticSegment("staticDatabaseAccounts", "databaseAccounts", "databaseAccounts"),
		resourceids.UserSpecifiedSegment("databaseAccountName", "databaseAccountValue"),
		resourceids.StaticSegment("staticMongodbDatabases", "mongodbDatabases", "mongodbDatabases"),
		resourceids.UserSpecifiedSegment("mongodbDatabaseName", "mongodbDatabaseValue"),
	}
}

// String returns a human-readable description of this Mongodb Database ID
func (id MongodbDatabaseId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Database Account Name: %q", id.DatabaseAccountName),
		fmt.Sprintf("Mongodb Database Name: %q", id.MongodbDatabaseName),
	}
	return fmt.Sprintf("Mongodb Database (%s)", strings.Join(components, "\n"))
}
