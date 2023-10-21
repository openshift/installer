package integrationaccountassemblies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = AssemblyId{}

// AssemblyId is a struct representing the Resource ID for a Assembly
type AssemblyId struct {
	SubscriptionId         string
	ResourceGroupName      string
	IntegrationAccountName string
	AssemblyName           string
}

// NewAssemblyID returns a new AssemblyId struct
func NewAssemblyID(subscriptionId string, resourceGroupName string, integrationAccountName string, assemblyName string) AssemblyId {
	return AssemblyId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		IntegrationAccountName: integrationAccountName,
		AssemblyName:           assemblyName,
	}
}

// ParseAssemblyID parses 'input' into a AssemblyId
func ParseAssemblyID(input string) (*AssemblyId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssemblyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssemblyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'integrationAccountName' was not found in the resource id %q", input)
	}

	if id.AssemblyName, ok = parsed.Parsed["assemblyName"]; !ok {
		return nil, fmt.Errorf("the segment 'assemblyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseAssemblyIDInsensitively parses 'input' case-insensitively into a AssemblyId
// note: this method should only be used for API response data and not user input
func ParseAssemblyIDInsensitively(input string) (*AssemblyId, error) {
	parser := resourceids.NewParserFromResourceIdType(AssemblyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := AssemblyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.IntegrationAccountName, ok = parsed.Parsed["integrationAccountName"]; !ok {
		return nil, fmt.Errorf("the segment 'integrationAccountName' was not found in the resource id %q", input)
	}

	if id.AssemblyName, ok = parsed.Parsed["assemblyName"]; !ok {
		return nil, fmt.Errorf("the segment 'assemblyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateAssemblyID checks that 'input' can be parsed as a Assembly ID
func ValidateAssemblyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseAssemblyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Assembly ID
func (id AssemblyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Logic/integrationAccounts/%s/assemblies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.IntegrationAccountName, id.AssemblyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Assembly ID
func (id AssemblyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftLogic", "Microsoft.Logic", "Microsoft.Logic"),
		resourceids.StaticSegment("staticIntegrationAccounts", "integrationAccounts", "integrationAccounts"),
		resourceids.UserSpecifiedSegment("integrationAccountName", "integrationAccountValue"),
		resourceids.StaticSegment("staticAssemblies", "assemblies", "assemblies"),
		resourceids.UserSpecifiedSegment("assemblyName", "assemblyValue"),
	}
}

// String returns a human-readable description of this Assembly ID
func (id AssemblyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Integration Account Name: %q", id.IntegrationAccountName),
		fmt.Sprintf("Assembly Name: %q", id.AssemblyName),
	}
	return fmt.Sprintf("Assembly (%s)", strings.Join(components, "\n"))
}
