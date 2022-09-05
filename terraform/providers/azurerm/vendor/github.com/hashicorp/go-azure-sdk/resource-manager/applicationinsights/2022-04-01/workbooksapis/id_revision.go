package workbooksapis

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = RevisionId{}

// RevisionId is a struct representing the Resource ID for a Revision
type RevisionId struct {
	SubscriptionId    string
	ResourceGroupName string
	ResourceName      string
	RevisionId        string
}

// NewRevisionID returns a new RevisionId struct
func NewRevisionID(subscriptionId string, resourceGroupName string, resourceName string, revisionId string) RevisionId {
	return RevisionId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ResourceName:      resourceName,
		RevisionId:        revisionId,
	}
}

// ParseRevisionID parses 'input' into a RevisionId
func ParseRevisionID(input string) (*RevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RevisionId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RevisionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.RevisionId, ok = parsed.Parsed["revisionId"]; !ok {
		return nil, fmt.Errorf("the segment 'revisionId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseRevisionIDInsensitively parses 'input' case-insensitively into a RevisionId
// note: this method should only be used for API response data and not user input
func ParseRevisionIDInsensitively(input string) (*RevisionId, error) {
	parser := resourceids.NewParserFromResourceIdType(RevisionId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := RevisionId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ResourceName, ok = parsed.Parsed["resourceName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceName' was not found in the resource id %q", input)
	}

	if id.RevisionId, ok = parsed.Parsed["revisionId"]; !ok {
		return nil, fmt.Errorf("the segment 'revisionId' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateRevisionID checks that 'input' can be parsed as a Revision ID
func ValidateRevisionID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseRevisionID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Revision ID
func (id RevisionId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Insights/workbooks/%s/revisions/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ResourceName, id.RevisionId)
}

// Segments returns a slice of Resource ID Segments which comprise this Revision ID
func (id RevisionId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftInsights", "Microsoft.Insights", "Microsoft.Insights"),
		resourceids.StaticSegment("staticWorkbooks", "workbooks", "workbooks"),
		resourceids.UserSpecifiedSegment("resourceName", "resourceValue"),
		resourceids.StaticSegment("staticRevisions", "revisions", "revisions"),
		resourceids.UserSpecifiedSegment("revisionId", "revisionIdValue"),
	}
}

// String returns a human-readable description of this Revision ID
func (id RevisionId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Resource Name: %q", id.ResourceName),
		fmt.Sprintf("Revision: %q", id.RevisionId),
	}
	return fmt.Sprintf("Revision (%s)", strings.Join(components, "\n"))
}
