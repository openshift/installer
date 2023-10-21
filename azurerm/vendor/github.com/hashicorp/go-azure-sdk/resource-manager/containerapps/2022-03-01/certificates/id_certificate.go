package certificates

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = CertificateId{}

// CertificateId is a struct representing the Resource ID for a Certificate
type CertificateId struct {
	SubscriptionId         string
	ResourceGroupName      string
	ManagedEnvironmentName string
	CertificateName        string
}

// NewCertificateID returns a new CertificateId struct
func NewCertificateID(subscriptionId string, resourceGroupName string, managedEnvironmentName string, certificateName string) CertificateId {
	return CertificateId{
		SubscriptionId:         subscriptionId,
		ResourceGroupName:      resourceGroupName,
		ManagedEnvironmentName: managedEnvironmentName,
		CertificateName:        certificateName,
	}
}

// ParseCertificateID parses 'input' into a CertificateId
func ParseCertificateID(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CertificateId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, fmt.Errorf("the segment 'certificateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseCertificateIDInsensitively parses 'input' case-insensitively into a CertificateId
// note: this method should only be used for API response data and not user input
func ParseCertificateIDInsensitively(input string) (*CertificateId, error) {
	parser := resourceids.NewParserFromResourceIdType(CertificateId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := CertificateId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.ManagedEnvironmentName, ok = parsed.Parsed["managedEnvironmentName"]; !ok {
		return nil, fmt.Errorf("the segment 'managedEnvironmentName' was not found in the resource id %q", input)
	}

	if id.CertificateName, ok = parsed.Parsed["certificateName"]; !ok {
		return nil, fmt.Errorf("the segment 'certificateName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateCertificateID checks that 'input' can be parsed as a Certificate ID
func ValidateCertificateID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseCertificateID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Certificate ID
func (id CertificateId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.App/managedEnvironments/%s/certificates/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.ManagedEnvironmentName, id.CertificateName)
}

// Segments returns a slice of Resource ID Segments which comprise this Certificate ID
func (id CertificateId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftApp", "Microsoft.App", "Microsoft.App"),
		resourceids.StaticSegment("staticManagedEnvironments", "managedEnvironments", "managedEnvironments"),
		resourceids.UserSpecifiedSegment("managedEnvironmentName", "managedEnvironmentValue"),
		resourceids.StaticSegment("staticCertificates", "certificates", "certificates"),
		resourceids.UserSpecifiedSegment("certificateName", "certificateValue"),
	}
}

// String returns a human-readable description of this Certificate ID
func (id CertificateId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Managed Environment Name: %q", id.ManagedEnvironmentName),
		fmt.Sprintf("Certificate Name: %q", id.CertificateName),
	}
	return fmt.Sprintf("Certificate (%s)", strings.Join(components, "\n"))
}
