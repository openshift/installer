package validation

import (
	"fmt"
	"testing"

	"github.com/openshift/installer/pkg/types/azure"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	subscriptionID        = "08675309-1111-2222-3333-303606808909"
	resourceGroup         = "test-resource-group"
	diskEncryptionSetName = "test-encryption-set"
	diskEncryptionSetID   = fmt.Sprintf("/subscriptions/%s/resourceGroups/%s/providers/Microsoft.Compute/diskEncryptionSets/%s",
		subscriptionID, resourceGroup, diskEncryptionSetName)
)

func validDiskEncryptionMachinePool() *azure.MachinePool {
	return &azure.MachinePool{
		OSDisk: azure.OSDisk{
			DiskType: "Premium_LRS",
			DiskEncryptionSet: &azure.DiskEncryptionSet{
				SubscriptionID: subscriptionID,
				ResourceGroup:  resourceGroup,
				Name:           diskEncryptionSetName,
			},
		},
	}
}

func TestValidateDiskEncryption(t *testing.T) {
	cases := []struct {
		name      string
		pool      *azure.MachinePool
		cloudName azure.CloudEnvironment
		expected  string
	}{
		{
			name:      "valid disk encryption set",
			pool:      validDiskEncryptionMachinePool(),
			cloudName: azure.PublicCloud,
			expected:  "",
		},
		{
			name:      "invalid disk encryption set (platform is stack cloud)",
			pool:      validDiskEncryptionMachinePool(),
			cloudName: azure.StackCloud,
			expected:  fmt.Sprintf(`azure.osDisk.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:"%s", Name:"%s"}: disk encryption sets are not supported on this platform`, subscriptionID, resourceGroup, diskEncryptionSetName),
		},
		{
			name: "invalid disk encryption set (missing name)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.OSDisk.DiskEncryptionSet.Name = ""
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  `azure.osDisk.diskEncryptionSet.name: Required value: name is required when specifying a diskEncryptionSet`,
		},
		{
			name: "invalid disk encryption set (missing resource group)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.OSDisk.DiskEncryptionSet.ResourceGroup = ""
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  `azure.osDisk.diskEncryptionSet.resourceGroup: Required value: resourceGroup is required when specifying a diskEncryptionSet`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDiskEncryption(tc.pool, tc.cloudName, field.NewPath("azure").Child("osDisk", "diskEncryptionSet")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}

func TestValidateEncryptionAtHost(t *testing.T) {
	cases := []struct {
		name      string
		pool      *azure.MachinePool
		cloudName azure.CloudEnvironment
		expected  string
	}{
		{
			name: "valid disk encryption at host",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.EncryptionAtHost = true
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  "",
		},
		{
			name: "invalid disk encryption at host (platform is stack cloud)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.EncryptionAtHost = true
				return p
			}(),
			cloudName: azure.StackCloud,
			expected:  `azure.osDisk.encryptionAtHost: Invalid value: true: encryption at host is not supported on this platform`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEncryptionAtHost(tc.pool, tc.cloudName, field.NewPath("azure").Child("osDisk", "encryptionAtHost")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
