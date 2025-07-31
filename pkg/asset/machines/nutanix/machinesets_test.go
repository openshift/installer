package nutanix

import (
	"testing"

	v1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
	"k8s.io/utils/ptr"
)

// provider1 is a function variable to allow assignment for testing
var provider1 = func(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string, failureDomain *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
	// default implementation of provider
	return nil, nil
}

// mockProvider is a mock function that returns a mock provider configuration
func mockProvider(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage, userDataSecret string, fd *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
	// Mock return value matching the expected type
	return &v1.NutanixMachineProviderConfig{}, nil
}

// setProviderForTest reassigns the provider function for testing
func setProviderForTest() {
	provider1 = mockProvider
}

// clearProviderForTest resets the provider to its original state
func clearProviderForTest() {
	provider1 = func(clusterID string, platform *nutanix.Platform, mpool *nutanix.MachinePool, osImage string, userDataSecret string, failureDomain *nutanix.FailureDomain) (*v1.NutanixMachineProviderConfig, error) {
		// original provider logic
		return nil, nil
	}
}

// TestMachineSetsReplicaNil tests the creation of machine sets when replica count is nil
func TestMachineSetsReplicaNil(t *testing.T) {
	setProviderForTest()
	defer clearProviderForTest()

	// Nutanix InstallConfig with no failure domains (simulate a real-world structure)
	ic := &types.InstallConfig{
		Platform: types.Platform{
			Nutanix: &nutanix.Platform{
				// purposely *no* FailureDomains to match the test case
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

	// Machine pool references to failure domains
	var machinePoolList []types.MachinePool
	machinePool1 := &types.MachinePool{
		Name: "worker",
		Platform: types.MachinePoolPlatform{
			Nutanix: &nutanix.MachinePool{
				FailureDomains: []string{"region-dc1-zone-undefined", "region-dc2-zone-undefined"},
			},
		},
		Replicas: ptr.To(int64(3)),
	}
	machinePool2 := &types.MachinePool{
		Name: "worker",
		Platform: types.MachinePoolPlatform{
			Nutanix: &nutanix.MachinePool{
				FailureDomains: []string{},
			},
		},
		Replicas: ptr.To(int64(3)),
	}
	machinePool3 := &types.MachinePool{
		Name: "worker",
		Platform: types.MachinePoolPlatform{
			Nutanix: &nutanix.MachinePool{
				FailureDomains: []string{"region-dc1-zone-undefined", "region-dc2-zone-undefined"},
			},
		},
	}

	// Add machine pools to list
	machinePoolList = append(machinePoolList, *machinePool1)
	machinePoolList = append(machinePoolList, *machinePool2)
	machinePoolList = append(machinePoolList, *machinePool3)

	// Iterate over machine pool list and check machine set creation
	for i := range machinePoolList {
		machinesets, err := MachineSets("test", ic, &machinePoolList[i], "fake-img", "worker", "userdata")
		if err != nil {
			t.Error(err)
			return
		}

		// Ensure that the number of machine sets matches the number of failure domains
		expectedNumMachineSets := len(machinePoolList[i].Platform.Nutanix.FailureDomains)
		if expectedNumMachineSets == 0 {
			expectedNumMachineSets = 1
		}
		if len(machinesets) != expectedNumMachineSets {
			t.Errorf("unexpected number of machine sets:\nwant: %v\ngot:  %v", expectedNumMachineSets, len(machinesets))
		}
	}
}
