package manifests

import (
	"testing"

	"github.com/ghodss/yaml"
	"github.com/stretchr/testify/assert"

	storagev1 "k8s.io/api/storage/v1"

	azuretypes "github.com/openshift/installer/pkg/types/azure"
)

// TestStorageClassConfig tests that a valid storage class manifest is created.
func TestStorageClassConfig(t *testing.T) {
	createStorageClassConfig := func() *storagev1.StorageClass {
		machinePool := &azuretypes.MachinePool{
			OSDisk: azuretypes.OSDisk{
				DiskType: "Premium_LRS",
				DiskEncryptionSet: &azuretypes.DiskEncryptionSet{
					SubscriptionID: "08675309-1111-2222-3333-303606808909",
					ResourceGroup:  "test-resource-group",
					Name:           "test-encryption-set",
				},
			},
		}
		return azureStorageClass(machinePool, "master")
	}

	expectedConfig := createStorageClassConfig()
	expectedYaml := `allowVolumeExpansion: true
apiVersion: storage.k8s.io/v1
kind: StorageClass
metadata:
  creationTimestamp: null
  name: master-managed-premium-encrypted-cmk
parameters:
  diskencryptionsetid: /subscriptions/08675309-1111-2222-3333-303606808909/resourceGroups/test-resource-group/providers/Microsoft.Compute/diskEncryptionSets/test-encryption-set
  kind: Managed
  resourcegroup: test-resource-group
  skuname: Premium_LRS
provisioner: disk.csi.azure.com
reclaimPolicy: Delete
volumeBindingMode: WaitForFirstConsumer
`
	yaml, err := yaml.Marshal(expectedConfig)
	assert.NoError(t, err, "failed to create storage class config")
	assert.Equal(t, expectedYaml, string(yaml), "unexpected storage class config")
}
