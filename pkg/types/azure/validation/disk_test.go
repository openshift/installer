package validation

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/azure"
)

var (
	subscriptionID        = "aF675309-bE11-cD22-aF33-bE3606808909"
	resourceGroup         = "Test-res.o(ur)Ce_gRoup"
	diskEncryptionSetName = "teSt-encrypTion_Set"
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
			expected:  fmt.Sprintf(`diskEncryptionSet.diskEncryptionSet: Invalid value: azure.DiskEncryptionSet{SubscriptionID:"%s", ResourceGroup:".+", Name:"%s"}: disk encryption sets are not supported on this platform`, subscriptionID, diskEncryptionSetName),
		},
		{
			name: "invalid disk encryption set (invalid subscription ID)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.OSDisk.DiskEncryptionSet.SubscriptionID = "invalid"
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  `subscriptionID: Invalid value: "invalid": invalid subscription ID format`,
		},
		{
			name: "invalid disk encryption set (invalid resource group)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.OSDisk.DiskEncryptionSet.ResourceGroup = ""
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  `resourceGroup: Invalid value: "": invalid resource group format`,
		},
		{
			name: "invalid disk encryption set (invalid name)",
			pool: func() *azure.MachinePool {
				p := validDiskEncryptionMachinePool()
				p.OSDisk.DiskEncryptionSet.Name = ""
				return p
			}(),
			cloudName: azure.PublicCloud,
			expected:  `diskEncryptionSetName: Invalid value: "": invalid name format`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateDiskEncryption(tc.pool, tc.cloudName, field.NewPath("test-path")).ToAggregate()
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
			expected:  `encryptionAtHost: Invalid value: true: encryption at host is not supported on this platform`,
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateEncryptionAtHost(tc.pool, tc.cloudName, field.NewPath("test-path")).ToAggregate()
			if tc.expected == "" {
				assert.NoError(t, err)
			} else {
				assert.Regexp(t, tc.expected, err)
			}
		})
	}
}
