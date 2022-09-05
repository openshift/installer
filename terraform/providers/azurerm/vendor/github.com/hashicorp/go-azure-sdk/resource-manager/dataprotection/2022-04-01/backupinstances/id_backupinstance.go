package backupinstances

import (
	"fmt"
	"strings"

	"github.com/hashicorp/go-azure-helpers/resourcemanager/resourceids"
)

var _ resourceids.ResourceId = BackupInstanceId{}

// BackupInstanceId is a struct representing the Resource ID for a Backup Instance
type BackupInstanceId struct {
	SubscriptionId     string
	ResourceGroupName  string
	VaultName          string
	BackupInstanceName string
}

// NewBackupInstanceID returns a new BackupInstanceId struct
func NewBackupInstanceID(subscriptionId string, resourceGroupName string, vaultName string, backupInstanceName string) BackupInstanceId {
	return BackupInstanceId{
		SubscriptionId:     subscriptionId,
		ResourceGroupName:  resourceGroupName,
		VaultName:          vaultName,
		BackupInstanceName: backupInstanceName,
	}
}

// ParseBackupInstanceID parses 'input' into a BackupInstanceId
func ParseBackupInstanceID(input string) (*BackupInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupInstanceId{})
	parsed, err := parser.Parse(input, false)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.BackupInstanceName, ok = parsed.Parsed["backupInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupInstanceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ParseBackupInstanceIDInsensitively parses 'input' case-insensitively into a BackupInstanceId
// note: this method should only be used for API response data and not user input
func ParseBackupInstanceIDInsensitively(input string) (*BackupInstanceId, error) {
	parser := resourceids.NewParserFromResourceIdType(BackupInstanceId{})
	parsed, err := parser.Parse(input, true)
	if err != nil {
		return nil, fmt.Errorf("parsing %q: %+v", input, err)
	}

	var ok bool
	id := BackupInstanceId{}

	if id.SubscriptionId, ok = parsed.Parsed["subscriptionId"]; !ok {
		return nil, fmt.Errorf("the segment 'subscriptionId' was not found in the resource id %q", input)
	}

	if id.ResourceGroupName, ok = parsed.Parsed["resourceGroupName"]; !ok {
		return nil, fmt.Errorf("the segment 'resourceGroupName' was not found in the resource id %q", input)
	}

	if id.VaultName, ok = parsed.Parsed["vaultName"]; !ok {
		return nil, fmt.Errorf("the segment 'vaultName' was not found in the resource id %q", input)
	}

	if id.BackupInstanceName, ok = parsed.Parsed["backupInstanceName"]; !ok {
		return nil, fmt.Errorf("the segment 'backupInstanceName' was not found in the resource id %q", input)
	}

	return &id, nil
}

// ValidateBackupInstanceID checks that 'input' can be parsed as a Backup Instance ID
func ValidateBackupInstanceID(input interface{}, key string) (warnings []string, errors []error) {
	v, ok := input.(string)
	if !ok {
		errors = append(errors, fmt.Errorf("expected %q to be a string", key))
		return
	}

	if _, err := ParseBackupInstanceID(v); err != nil {
		errors = append(errors, err)
	}

	return
}

// ID returns the formatted Backup Instance ID
func (id BackupInstanceId) ID() string {
	fmtString := "/subscriptions/%s/resourceGroups/%s/providers/Microsoft.DataProtection/backupVaults/%s/backupInstances/%s"
	return fmt.Sprintf(fmtString, id.SubscriptionId, id.ResourceGroupName, id.VaultName, id.BackupInstanceName)
}

// Segments returns a slice of Resource ID Segments which comprise this Backup Instance ID
func (id BackupInstanceId) Segments() []resourceids.Segment {
	return []resourceids.Segment{
		resourceids.StaticSegment("staticSubscriptions", "subscriptions", "subscriptions"),
		resourceids.SubscriptionIdSegment("subscriptionId", "12345678-1234-9876-4563-123456789012"),
		resourceids.StaticSegment("staticResourceGroups", "resourceGroups", "resourceGroups"),
		resourceids.ResourceGroupSegment("resourceGroupName", "example-resource-group"),
		resourceids.StaticSegment("staticProviders", "providers", "providers"),
		resourceids.ResourceProviderSegment("staticMicrosoftDataProtection", "Microsoft.DataProtection", "Microsoft.DataProtection"),
		resourceids.StaticSegment("staticBackupVaults", "backupVaults", "backupVaults"),
		resourceids.UserSpecifiedSegment("vaultName", "vaultValue"),
		resourceids.StaticSegment("staticBackupInstances", "backupInstances", "backupInstances"),
		resourceids.UserSpecifiedSegment("backupInstanceName", "backupInstanceValue"),
	}
}

// String returns a human-readable description of this Backup Instance ID
func (id BackupInstanceId) String() string {
	components := []string{
		fmt.Sprintf("Subscription: %q", id.SubscriptionId),
		fmt.Sprintf("Resource Group Name: %q", id.ResourceGroupName),
		fmt.Sprintf("Vault Name: %q", id.VaultName),
		fmt.Sprintf("Backup Instance Name: %q", id.BackupInstanceName),
	}
	return fmt.Sprintf("Backup Instance (%s)", strings.Join(components, "\n"))
}
