package validation

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ibmcloud"
	ibmclouddefaults "github.com/openshift/installer/pkg/types/ibmcloud/defaults"
)

var (
	validType            = "valid-type"
	validZones           = []string{"us-east-1", "us-east-2"}
	validEncryptionKey   = "crn:v1:bluemix:public:kms:global:a/accountid:service:key:keyid"
	invalidEncryptionKey = "v1:bluemix:kms:global:a/accountid:service:key:keyid"
)

func TestValidateMachinePool(t *testing.T) {
	platform := &ibmcloud.Platform{Region: "us-east"}
	cases := []struct {
		name        string
		machinepool *ibmcloud.MachinePool
		valid       bool
	}{
		{
			name:        "minimal",
			machinepool: &ibmcloud.MachinePool{},
			valid:       true,
		},
		{
			name: "valid type",
			machinepool: &ibmcloud.MachinePool{
				InstanceType: validType,
			},
			valid: true,
		},
		{
			name: "valid zones",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
			},
			valid: true,
		},
		{
			name: "valid bootVolume",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					EncryptionKey: validEncryptionKey,
				},
			},
			valid: true,
		},
		{
			name: "invalid bootVolume",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					EncryptionKey: invalidEncryptionKey,
				},
			},
			valid: false,
		},
		{
			name: "valid sdp bootVolume",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					Profile:   "sdp",
					SizeGiB:   500,
					IOPS:      3000,
					Bandwidth: 1000,
				},
			},
			valid: true,
		},
		{
			name: "valid custom bootVolume with iops",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					Profile: "custom",
					SizeGiB: 100,
					IOPS:    3000,
				},
			},
			valid: true,
		},
		{
			name: "invalid bootVolume profile",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					Profile: "not-a-real-profile",
				},
			},
			valid: false,
		},
		{
			name: "invalid iops on a profile that does not support it",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					Profile: "general-purpose",
					IOPS:    3000,
				},
			},
			valid: false,
		},
		{
			name: "invalid bandwidth on a non-sdp profile",
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{
					Profile:   "custom",
					IOPS:      3000,
					Bandwidth: 1000,
				},
			},
			valid: false,
		},
		{
			name: "valid dedicatedHosts 1",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Profile: validType,
					},
					{
						Profile: validType,
					},
				},
				InstanceType: validType,
			},
			valid: true,
		},
		{
			name: "valid dedicatedHosts 2",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
					{
						Name: "name",
					},
				},
				InstanceType: validType,
			},
			valid: true,
		},
		{
			name: "valid dedicatedHosts 3",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
					{
						Profile: validType,
					},
				},
				InstanceType: validType,
			},
			valid: true,
		},
		{
			name: "invalid dedicatedHosts 1",
			machinepool: &ibmcloud.MachinePool{
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
					{
						Profile: validType,
					},
				},
				InstanceType: validType,
			},
			valid: false,
		},
		{
			name: "invalid dedicatedHosts 2",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
				},
				InstanceType: validType,
			},
			valid: false,
		},
		{
			name: "invalid dedicatedHosts 3",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
					{
						Profile: "invalid",
					},
				},
				InstanceType: validType,
			},
			valid: false,
		},
		{
			name: "invalid dedicatedHosts 4",
			machinepool: &ibmcloud.MachinePool{
				Zones: validZones,
				DedicatedHosts: []ibmcloud.DedicatedHost{
					{
						Name: "name",
					},
					{
						Profile: validType,
					},
				},
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			err := ValidateMachinePool(platform, tc.machinepool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}

func TestValidateMachinePoolInheritedBootVolume(t *testing.T) {
	cases := []struct {
		name                   string
		defaultMachinePlatform *ibmcloud.MachinePool
		machinepool            *ibmcloud.MachinePool
		valid                  bool
	}{
		{
			name: "valid iops on a profile inherited from defaultMachinePlatform",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "sdp", SizeGiB: 250},
			},
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{IOPS: 12000},
			},
			valid: true,
		},
		{
			name: "valid bandwidth on a profile inherited from defaultMachinePlatform",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "sdp"},
			},
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Bandwidth: 2000},
			},
			valid: true,
		},
		{
			name: "invalid inherited iops on a profile that does not support it",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "custom", IOPS: 3000},
			},
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "general-purpose"},
			},
			valid: false,
		},
		{
			name: "invalid iops inherited from a profile-less default on a pool without a profile",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{IOPS: 3000},
			},
			machinepool: &ibmcloud.MachinePool{},
			valid:       false,
		},
		{
			name: "valid iops from a profile-less default completed by the pool's sdp profile",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{IOPS: 3000},
			},
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "sdp"},
			},
			valid: true,
		},
		{
			name: "invalid inherited bandwidth on a non-sdp profile",
			defaultMachinePlatform: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "sdp", Bandwidth: 2000},
			},
			machinepool: &ibmcloud.MachinePool{
				BootVolume: &ibmcloud.BootVolume{Profile: "custom"},
			},
			valid: false,
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			platform := &ibmcloud.Platform{Region: "us-east", DefaultMachinePlatform: tc.defaultMachinePlatform}
			// Pools are validated after defaulting merges defaultMachinePlatform into them.
			ibmclouddefaults.Apply(platform.DefaultMachinePlatform, tc.machinepool)
			err := ValidateMachinePool(platform, tc.machinepool, field.NewPath("test-path")).ToAggregate()
			if tc.valid {
				assert.NoError(t, err)
			} else {
				assert.Error(t, err)
			}
		})
	}
}
