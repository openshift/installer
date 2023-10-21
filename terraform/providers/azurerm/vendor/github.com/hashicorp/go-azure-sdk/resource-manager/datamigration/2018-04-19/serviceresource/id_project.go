package serviceresource

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = ProjectId{}

// ProjectId is a struct representing the Resource ID for a Project
type ProjectId struct {
	SubscriptionId    string
	ResourceGroupName string
	ServiceName       string
	ProjectName       string
}

// NewProjectID returns a new ProjectId struct
func NewProjectID(subscriptionId string, resourceGroupName string, serviceName string, projectName string) ProjectId {
	return ProjectId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		ServiceName:       serviceName,
		ProjectName:       projectName,
	}
}

// ParseProjectID parses 'input' into a ProjectId
func ParseProjectID(input string) (*ProjectId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProjectId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProjectId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, fmt.Errorf("the segment 'serviceName' was not found in the resource id %q", input)
	}

	if id.ProjectName, ok = parsed.Parsed["projectName"]; !ok {
		return nil, fmt.Errorf("the segment 'projectName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseProjectIDInsensitively parses 'input' case-insensitively into a ProjectId
// note: this method should only be used for API response data and not user input
func ParseProjectIDInsensitively(input string) (*ProjectId, error) {
	parser := resourceids.NewParserFromResourceIdType(ProjectId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := ProjectId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ServiceName, ok = parsed.Parsed["serviceName"]; !ok {
		return nil, fmt.Errorf("the segment 'serviceName' was not found in the resource id %q", input)
	}

	if id.ProjectName, ok = parsed.Parsed["projectName"]; !ok {
		return nil, fmt.Errorf("the segment 'projectName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateProjectID checks that 'input' can be parsed as a Project ID
func ValidateProjectID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseProjectID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Project ID
func (id ProjectId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataMigration/services/%s/projects/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ServiceName, id.ProjectName)
}

// Segments returns a slice of Resource ID Segments which comprise this Project ID
func (id ProjectId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.UserSpecifiedSegment("resourceGroupName", "resourceGroupValue"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataMigration", "Microsoft.DataMigration", "Microsoft.DataMigration"),
		resourceids.StaticSegment("staticServices", "services", "services"),
		resourceids.UserSpecifiedSegment("serviceName", "serviceValue"),
		resourceids.StaticSegment("staticProjects", "projects", "projects"),
		resourceids.UserSpecifiedSegment("projectName", "projectValue"),
	}
}

// String returns a human-readable description of this Project ID
func (id ProjectId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Service Name: %q", id.ServiceName),
		fmt.Sprintf("Project Name: %q", id.ProjectName),
	}
	return fmt.Sprintf("Project (%s)", strings.Join(components, "\n"))
}
