package backuppolicies

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = BackupPolicyId{}

// BackupPolicyId is a struct representing the Resource ID for a Backup Policy
type BackupPolicyId struct {
	SubscriptionId    string
	ResourceGroupName string
	BackupVaultName   string
	BackupPolicyName  string
}

// NewBackupPolicyID returns a new BackupPolicyId struct
func NewBackupPolicyID(subscriptionId string, resourceGroupName string, backupVaultName string, backupPolicyName string) BackupPolicyId {
	return BackupPolicyId{
		SubscriptionId:    subscriptionId,
		ResourceGroupName: resourceGroupName,
		BackupVaultName:   backupVaultName,
		BackupPolicyName:  backupPolicyName,
	}
}

// ParseBackupPolicyID parses 'input' into a BackupPolicyId
func ParseBackupPolicyID(input string) (*BackupPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupPolicyId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.BackupVaultName, ok = parsed.Parsed["backupVaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupVaultName' was not found in the resource id %q", input)
	}

	if id.BackupPolicyName, ok = parsed.Parsed["backupPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseBackupPolicyIDInsensitively parses 'input' case-insensitively into a BackupPolicyId
// note: this method should only be used for API response data and not user input
func ParseBackupPolicyIDInsensitively(input string) (*BackupPolicyId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupPolicyId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupPolicyId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.BackupVaultName, ok = parsed.Parsed["backupVaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupVaultName' was not found in the resource id %q", input)
	}

	if id.BackupPolicyName, ok = parsed.Parsed["backupPolicyName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupPolicyName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateBackupPolicyID checks that 'input' can be parsed as a Backup Policy ID
func ValidateBackupPolicyID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupPolicyID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Policy ID
func (id BackupPolicyId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/backupPolicies/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.BackupVaultName, id.BackupPolicyName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Policy ID
func (id BackupPolicyId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("backupVaultName", "backupVaultValue"),
		resourceids.StaticSegment("staticBackupPolicies", "backupPolicies", "backupPolicies"),
		resourceids.UserSpecifiedSegment("backupPolicyName", "backupPolicyValue"),
	}
}

// String returns a human-readable description of this Backup Policy ID
func (id BackupPolicyId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Backup Vault Name: %q", id.BackupVaultName),
		fmt.Sprintf("Backup Policy Name: %q", id.BackupPolicyName),
	}
	return fmt.Sprintf("Backup Policy (%s)", strings.Join(components, "\n"))
}
