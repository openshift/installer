package ibmcloud

import (
	"testing"

	"github.com/stretchr/testify/assert"

	ibmcloudprovider "github.com/openshift/machine-api-provider-ibmcloud/pkg/apis/ibmcloudprovider/v1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ibmcloud"
)

func testInstallConfig(pool *ibmcloud.MachinePool) (*types.InstallConfig, *types.MachinePool) {
	config := &types.InstallConfig{
		Platform: types.Platform{
			IBMCloud: &ibmcloud.Platform{
				Region: "us-east",
			},
		},
	}
	replicas := int64(1)
	machinePool := &types.MachinePool{
		Name:     "master",
		Replicas: &replicas,
		Platform: types.MachinePoolPlatform{
			IBMCloud: pool,
		},
	}
	return config, machinePool
}

// TestMachinesBootVolume verifies that boot volume options configured on the
// machine pool reach the generated MAPI provider spec. The CAPI machines are
// built from this same provider spec, but GenerateMachines cannot be covered
// here because it looks up the machine SSH key through the IBM Cloud API.
func TestMachinesBootVolume(t *testing.T) {
	cases := []struct {
		name     string
		pool     *ibmcloud.MachinePool
		expected ibmcloudprovider.IBMCloudMachineBootVolume
	}{
		{
			name: "no boot volume configured",
			pool: &ibmcloud.MachinePool{
				InstanceType: "bx2-4x16",
				Zones:        []string{"us-east-1"},
			},
			expected: ibmcloudprovider.IBMCloudMachineBootVolume{},
		},
		{
			name: "encryption key only",
			pool: &ibmcloud.MachinePool{
				InstanceType: "bx2-4x16",
				Zones:        []string{"us-east-1"},
				BootVolume: &ibmcloud.BootVolume{
					EncryptionKey: "crn:v1:bluemix:public:kms:us-east:a/accountid:instanceid:key:keyid",
				},
			},
			expected: ibmcloudprovider.IBMCloudMachineBootVolume{
				EncryptionKey: "crn:v1:bluemix:public:kms:us-east:a/accountid:instanceid:key:keyid",
			},
		},
		{
			name: "sdp profile with iops and bandwidth",
			pool: &ibmcloud.MachinePool{
				InstanceType: "bx2-4x16",
				Zones:        []string{"us-east-1"},
				BootVolume: &ibmcloud.BootVolume{
					Profile:   "sdp",
					SizeGiB:   500,
					IOPS:      16000,
					Bandwidth: 2000,
				},
			},
			expected: ibmcloudprovider.IBMCloudMachineBootVolume{
				Profile:   "sdp",
				SizeGiB:   500,
				Iops:      16000,
				Bandwidth: 2000,
			},
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config, pool := testInstallConfig(tc.pool)
			machines, err := Machines("test-cluster", config, map[string]string{"us-east-1": "subnet-1"}, pool, "master", "master-user-data")
			if !assert.NoError(t, err) || !assert.Len(t, machines, 1) {
				return
			}

			providerSpec, ok := machines[0].Spec.ProviderSpec.Value.Object.(*ibmcloudprovider.IBMCloudMachineProviderSpec)
			if !assert.True(t, ok, "expected an IBMCloudMachineProviderSpec") {
				return
			}
			assert.Equal(t, tc.expected, providerSpec.BootVolume)
		})
	}
}
