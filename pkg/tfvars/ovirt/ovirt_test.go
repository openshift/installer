package ovirt

import (
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/cluster-api-provider-ovirt/pkg/apis/ovirtprovider/v1beta1"
	"github.com/openshift/installer/pkg/types/ovirt"
)

func defaultMachineSpec() *v1beta1.OvirtMachineProviderSpec {
	return &v1beta1.OvirtMachineProviderSpec{
		InstanceTypeId: "",
		CPU: &v1beta1.CPU{
			Sockets: 1,
			Cores:   8,
			Threads: 1,
		},
		MemoryMB:            16000,
		OSDisk:              &v1beta1.Disk{SizeGB: 31},
		VMType:              "high_performance",
		AffinityGroupsNames: []string{"clusterName-xxxxx-controlplane"},
		AutoPinningPolicy:   "none",
	}
}

var defaultTerraformOvirtVarsJSON = `{
  "ovirt_url": "https://ovirt-engine.com",
  "ovirt_username": "admin",
  "ovirt_password": "test",
  "ovirt_cafile": "ca-file-content",
  "ovirt_ca_bundle": "",
  "ovirt_insecure": false,
  "ovirt_cluster_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "ovirt_storage_domain_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "ovirt_network_name": "ovirt-network",
  "ovirt_vnic_profile_id": "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
  "ovirt_affinity_groups": [
    {
      "description": "AffinityGroup for spreading each control plane machines to a different host",
      "enforcing": true,
      "name": "clusterName-xxxxx-controlplane",
      "priority": 5
    },
    {
      "description": "AffinityGroup for spreading each compute machine to a different host",
      "enforcing": true,
      "name": "clusterName-xxxxx-compute",
      "priority": 3
    }
  ],
  "ovirt_base_image_name": "some-base-image",
  "ovirt_master_instance_type_id": "",
  "ovirt_master_vm_type": "high_performance",
  "ovirt_master_memory": 16000,
  "ovirt_master_cores": 8,
  "ovirt_master_sockets": 1,
  "ovirt_master_threads": 1,
  "ovirt_master_os_disk_gb": 31,
  "ovirt_master_affinity_groups": [
    "clusterName-xxxxx-controlplane"
  ],
  "ovirt_master_auto_pinning_policy": "none",
  "ovirt_master_hugepages": 0,
  "ovirt_master_clone": null,
  "ovirt_master_sparse": null,
  "ovirt_master_format": ""
}`

func TestSetPlatformDefaults(t *testing.T) {
	cases := []struct {
		name            string
		auth            Auth
		clusterID       string
		storageDomainID string
		networkName     string
		vnicProfileID   string
		baseImage       string
		infraID         string
		masterSpec      *v1beta1.OvirtMachineProviderSpec
		affinityGroups  []ovirt.AffinityGroup
		masterCount     int
		expected        []byte
	}{
		{
			name: "default",
			auth: Auth{
				URL:      "https://ovirt-engine.com",
				Username: "admin",
				Password: "test",
				Cafile:   "ca-file-content",
			},
			clusterID:       "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			storageDomainID: "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			networkName:     "ovirt-network",
			vnicProfileID:   "xxxxxxxx-xxxx-xxxx-xxxx-xxxxxxxxxxxx",
			baseImage:       "some-base-image",
			infraID:         "clusterName-xxxxx",
			masterSpec:      defaultMachineSpec(),
			affinityGroups: []ovirt.AffinityGroup{{
				Name:        "controlplane",
				Priority:    5,
				Description: "AffinityGroup for spreading each control plane machines to a different host",
				Enforcing:   true,
			}, {
				Name:        "compute",
				Priority:    3,
				Description: "AffinityGroup for spreading each compute machine to a different host",
				Enforcing:   true,
			}},
			masterCount: 3,
			expected:    []byte(defaultTerraformOvirtVarsJSON),
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			tfVar, err := TFVars(
				tc.auth, tc.clusterID, tc.storageDomainID, tc.networkName,
				tc.vnicProfileID, tc.baseImage, tc.infraID, tc.masterSpec,
				tc.affinityGroups)
			if err != nil {
				t.Fatalf("failed during test case %s: %v", tc.name, err)
			}
			assert.JSONEq(t, string(tc.expected), string(tfVar), "unexpected ovirt-specific Terraform variables file")
		})
	}
}
