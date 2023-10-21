package publicmaintenanceconfigurations

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = PublicMaintenanceConfigurationId{}

// PublicMaintenanceConfigurationId is a struct representing the Resource ID for a Public Maintenance Configuration
type PublicMaintenanceConfigurationId struct {
	SubscriptionId                     string
	PublicMaintenanceConfigurationName string
}

// NewPublicMaintenanceConfigurationID returns a new PublicMaintenanceConfigurationId struct
func NewPublicMaintenanceConfigurationID(subscriptionId string, publicMaintenanceConfigurationName string) PublicMaintenanceConfigurationId {
	return PublicMaintenanceConfigurationId{
		SubscriptionId:                     subscriptionId,
		PublicMaintenanceConfigurationName: publicMaintenanceConfigurationName,
	}
}

// ParsePublicMaintenanceConfigurationID parses 'input' into a PublicMaintenanceConfigurationId
func ParsePublicMaintenanceConfigurationID(input string) (*PublicMaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublicMaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublicMaintenanceConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.PublicMaintenanceConfigurationName, ok = parsed.Parsed["publicMaintenanceConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'publicMaintenanceConfigurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParsePublicMaintenanceConfigurationIDInsensitively parses 'input' case-insensitively into a PublicMaintenanceConfigurationId
// note: this method should only be used for API response data and not user input
func ParsePublicMaintenanceConfigurationIDInsensitively(input string) (*PublicMaintenanceConfigurationId, error) {
	parser := resourceids.NewParserFromResourceIdType(PublicMaintenanceConfigurationId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := PublicMaintenanceConfigurationId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.PublicMaintenanceConfigurationName, ok = parsed.Parsed["publicMaintenanceConfigurationName"]; !ok {
		return nil, fmt.Errorf("the segment 'publicMaintenanceConfigurationName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidatePublicMaintenanceConfigurationID checks that 'input' can be parsed as a Public Maintenance Configuration ID
func ValidatePublicMaintenanceConfigurationID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParsePublicMaintenanceConfigurationID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Public Maintenance Configuration ID
func (id PublicMaintenanceConfigurationId) ID() string {
	fmtString := "/subscriptions/%s/providers/Microsoft.Maintenance/publicMaintenanceConfigurations/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.PublicMaintenanceConfigurationName)
}

// Segments returns a slice of Resource ID Segments which comprise this Public Maintenance Configuration ID
func (id PublicMaintenanceConfigurationId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftMaintenance", "Microsoft.Maintenance", "Microsoft.Maintenance"),
		resourceids.StaticSegment("staticPublicMaintenanceConfigurations", "publicMaintenanceConfigurations", "publicMaintenanceConfigurations"),
		resourceids.UserSpecifiedSegment("publicMaintenanceConfigurationName", "publicMaintenanceConfigurationValue"),
	}
}

// String returns a human-readable description of this Public Maintenance Configuration ID
func (id PublicMaintenanceConfigurationId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Public Maintenance Configuration Name: %q", id.PublicMaintenanceConfigurationName),
	}
	return fmt.Sprintf("Public Maintenance Configuration (%s)", strings.Join(components, "\n"))
}
