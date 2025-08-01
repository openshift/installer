// Package nutanix implements tests for Nutanix MachineSets generation.
package nutanix

import (
	"testing"

	"k8s.io/utils/ptr"

	v1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// provider1 is a function variable to allow assignment for testing.
var provider1 = func(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string, failureDomain *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
	// default implementation of provider
	return nil, nil
}

// mockProvider is a mock function that returns a mock provider configuration.
func mockProvider(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage, userDataSecret string, fd *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
	// Mock return value matching the expected type
	return &v1.NutanixMachineProviderConfig{}, nil
}

// setProviderForTest reassigns the provider function to the mock for testing.
func setProviderForTest() {
	provider1 = mockProvider
}

// clearProviderForTest resets the provider function to its original implementation.
func clearProviderForTest() {
	provider1 = func(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string, failureDomain *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
		// original provider logic placeholder
		return nil, nil
	}
}

// TestMachineSetsReplicaNil tests MachineSets generation when replica counts and failure domains vary.
func TestMachineSetsReplicaNil(t *testing.T) {
	setProviderForTest()
	defer clearProviderForTest()

	// Nutanix InstallConfig with failure domains defined.
	ic := &types.InstallConfig{
		Platform: types.Platform{
			Nutanix: &nutanix.Platform{
				FailureDomains: []nutanix.FailureDomain{
					{
						Name: "region-dc1-zone-undefined",
						PrismElement: nutanix.PrismElement{
							Name: "PE1",
							UUID: "123456-2345-3454-4455",
						},
						SubnetUUIDs: []string{"1234-3453432-342343-43", "23434234-324343-3434-2434"},
					},
					{
						Name: "region-dc2-zone-undefined",
						PrismElement: nutanix.PrismElement{
							Name: "PE2",
							UUID: "123456-2345-3454-4456",
						},
						SubnetUUIDs: []string{"1234-3453432-342343-432", "23434234-324343-3434-24345"},
					},
				},
				PrismElements: []nutanix.PrismElement{
					{
						Name: "PE1",
						UUID: "1234-3453-345-3333-4565",
					},
				},
				SubnetUUIDs: []string{"123456-2345-3454-4455"},
			},
		},
	}

	// Define machine pools with different failure domain configurations.
	machinePools := []*types.MachinePool{
		{
			Name: "worker",
			Platform: types.MachinePoolPlatform{
				Nutanix: &nutanix.MachinePool{
					FailureDomains: []string{"region-dc1-zone-undefined", "region-dc2-zone-undefined"},
				},
			},
			Replicas: ptr.To(int64(3)),
		},
		{
			Name: "worker",
			Platform: types.MachinePoolPlatform{
				Nutanix: &nutanix.MachinePool{
					FailureDomains: []string{},
				},
			},
			Replicas: ptr.To(int64(3)),
		},
		{
			Name: "worker",
			Platform: types.MachinePoolPlatform{
				Nutanix: &nutanix.MachinePool{
					FailureDomains: []string{"region-dc1-zone-undefined", "region-dc2-zone-undefined"},
				},
			},
			// Replicas is nil here
		},
	}

	// Test MachineSets creation for each machine pool.
	for _, mp := range machinePools {
		machinesets, err := MachineSets("test", ic, mp, "fake-img", "worker", "userdata")
		if err != nil {
			t.Error(err)
			return
		}

		expectedNum := len(mp.Platform.Nutanix.FailureDomains)
		if expectedNum == 0 {
			expectedNum = 1
		}

		if len(machinesets) != expectedNum {
			t.Errorf("unexpected number of machine sets:\nwant: %v\ngot:  %v", expectedNum, len(machinesets))
		}
	}
}
