package serverkeys

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = KeyId{}

// KeyId is a struct representing the Resource ID for a Key
type KeyId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServerName        string
	KeyName           string
}

// NewKeyID returns a new KeyId struct
func NewKeyID(subscriptionId string, resourceGroupName string, serverName string, keyName string) KeyId {
	return KeyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServerName:        serverName,
		KeyName:           keyName,
	}
}

// ParseKeyID parses 'input' into a KeyId
func ParseKeyID(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, fmt.Errorf("the segment 'serverName' was not found in the resource id %q", input)
	}

	if id.KeyName, ok = parsed.Parsed["keyName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseKeyIDInsensitively parses 'input' case-insensitively into a KeyId
// note: this method should only be used for API response data and not user input
func ParseKeyIDInsensitively(input string) (*KeyId, error) {
	parser := resourceids.NewParserFromResourceIdType(KeyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := KeyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServerName, ok = parsed.Parsed["serverName"]; !ok {
		return nil, fmt.Errorf("the segment 'serverName' was not found in the resource id %q", input)
	}

	if id.KeyName, ok = parsed.Parsed["keyName"]; !ok {
		return nil, fmt.Errorf("the segment 'keyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateKeyID checks that 'input' can be parsed as a Key ID
func ValidateKeyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseKeyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Key ID
func (id KeyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DBforPostgreSQL/servers/%s/keys/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServerName, id.KeyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Key ID
func (id KeyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDBforPostgreSQL", "Microsoft.DBforPostgreSQL", "Microsoft.DBforPostgreSQL"),
		resourceids.StaticSegment("staticServers", "servers", "servers"),
		resourceids.UserSpecifiedSegment("serverName", "serverValue"),
		resourceids.StaticSegment("staticKeys", "keys", "keys"),
		resourceids.UserSpecifiedSegment("keyName", "keyValue"),
	}
}

// String returns a human-readable description of this Key ID
func (id KeyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Server Name: %q", id.ServerName),
		fmt.Sprintf("Key Name: %q", id.KeyName),
	}
	return fmt.Sprintf("Key (%s)", strings.Join(components, "\n"))
}
