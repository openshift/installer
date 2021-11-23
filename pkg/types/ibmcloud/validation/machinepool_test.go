package validation

import (
	"testing"

	"github.com/openshift/installer/pkg/types/ibmcloud"
	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/util/validation/field"
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
