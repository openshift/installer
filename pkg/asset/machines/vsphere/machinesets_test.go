package vsphere

import (
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
	vsphereapis "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var machineComputePoolValidZones = types.MachinePool{
	Name:     "worker",
	Replicas: &machinePoolReplicas,
	Platform: types.MachinePoolPlatform{
		VSphere: &vsphere.MachinePool{
			NumCPUs:           4,
			NumCoresPerSocket: 2,
			MemoryMiB:         16384,
			OSDisk: vsphere.OSDisk{
				DiskSizeGB: 60,
			},
			Zones: []string{
				"region-dc1-zone-c1",
				"region-dc1-zone-c2",
				"region-dc1-zone-c3",
			},
		},
	},
	Hyperthreading: "true",
	Architecture:   types.ArchitectureAMD64,
}

var machineComputePoolReducedZones = types.MachinePool{
	Name:     "worker",
	Replicas: &machinePoolReplicas,
	Platform: types.MachinePoolPlatform{
		VSphere: &vsphere.MachinePool{
			NumCPUs:           4,
			NumCoresPerSocket: 2,
			MemoryMiB:         16384,
			OSDisk: vsphere.OSDisk{
				DiskSizeGB: 60,
			},
			Zones: []string{
				"region-dc1-zone-c1",
				"region-dc1-zone-c2",
			},
		},
	},
	Hyperthreading: "true",
	Architecture:   types.ArchitectureAMD64,
}

var machineComputePoolUndefinedZones = types.MachinePool{
	Name:     "worker",
	Replicas: &machinePoolReplicas,
	Platform: types.MachinePoolPlatform{
		VSphere: &vsphere.MachinePool{
			NumCPUs:           4,
			NumCoresPerSocket: 2,
			MemoryMiB:         16384,
			OSDisk: vsphere.OSDisk{
				DiskSizeGB: 60,
			},
			Zones: []string{
				"region-dc1-zone-undefined",
			},
		},
	},
	Hyperthreading: "true",
	Architecture:   types.ArchitectureAMD64,
}

var machineComputePoolNoZones = types.MachinePool{
	Name:     "worker",
	Replicas: &machinePoolReplicas,
	Platform: types.MachinePoolPlatform{
		VSphere: &vsphere.MachinePool{
			NumCPUs:           4,
			NumCoresPerSocket: 2,
			MemoryMiB:         16384,
			OSDisk: vsphere.OSDisk{
				DiskSizeGB: 60,
			},
		},
	},
	Hyperthreading: "true",
	Architecture:   types.ArchitectureAMD64,
}

func TestConfigMachinesets(t *testing.T) {
	clusterID := "test"
	installConfig, err := parseInstallConfig()
	testCases := []struct {
		testCase      string
		expectedError string
		machinePool   *types.MachinePool
		workspaces    []vsphereapis.Workspace
		maxReplicas   int
		minReplicas   int
	}{
		{
			testCase:    "machinepool with no defined zones should use legacy configuration",
			machinePool: &machineComputePoolNoZones,
			minReplicas: 3,
			maxReplicas: 3,
			workspaces: []vsphereapis.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "datacenter",
					Folder:       "/datacenter/vm/test",
					Datastore:    "datastore",
					ResourcePool: "/datacenter/host//Resources",
				},
			},
		},
		{
			testCase:    "zones distributed among compute machinesets(zone count matches machineset count)",
			machinePool: &machineComputePoolValidZones,
			maxReplicas: 1,
			minReplicas: 1,
			workspaces: []vsphereapis.Workspace{
				{
					Server:       "your.first-vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/test",
					Datastore:    "DC1_DS1",
					ResourcePool: "/dc1/host/DC1_C1/Resources",
				},
				{
					Server:       "your.first-vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/test",
					Datastore:    "DC1_DS2",
					ResourcePool: "/dc1/host/DC1_C2/Resources",
				},
				{
					Server:       "your.first-vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/test",
					Datastore:    "DC1_DS3",
					ResourcePool: "/dc1/host/DC1_C3/Resources",
				},
			},
		},
		{
			testCase:      "undefined zone in machinepool results in error",
			machinePool:   &machineComputePoolUndefinedZones,
			expectedError: "zone [region-dc1-zone-undefined] specified by machinepool is not defined",
		},
		{
			testCase:    "zones distributed among compute machinesets(zone less than compute machinesets count)",
			machinePool: &machineComputePoolReducedZones,
			maxReplicas: 2,
			minReplicas: 1,
			workspaces: []vsphereapis.Workspace{
				{
					Server:       "your.first-vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/test",
					Datastore:    "DC1_DS1",
					ResourcePool: "/dc1/host/DC1_C1/Resources",
				},
				{
					Server:       "your.first-vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/test",
					Datastore:    "DC1_DS2",
					ResourcePool: "/dc1/host/DC1_C2/Resources",
				},
			},
		},
	}

	for _, tc := range testCases {
		if err != nil {
			assert.Errorf(t, err, "unable to parse sample install config")
		}
		t.Run(tc.testCase, func(t *testing.T) {

			machineSets, err := MachineSets(clusterID, installConfig, tc.machinePool, "", "", "")
			assertOnUnexpectedErrorState(tc.expectedError, err, t)

			if len(tc.workspaces) > 0 {
				var matchCountByIndex []int
				for range tc.workspaces {
					matchCountByIndex = append(matchCountByIndex, 0)
				}
				for _, machineSet := range machineSets {
					replicas := int(*machineSet.Spec.Replicas)
					if replicas > tc.maxReplicas {
						t.Errorf("machineset has too many replicas[max: %d]", tc.maxReplicas)
					} else if replicas < tc.minReplicas {
						t.Errorf("machineset has too few replicas[max: %d]", tc.minReplicas)
					}
					machineWorkspace := machineSet.Spec.Template.Spec.ProviderSpec.Value.Object.(*vsphereapis.VSphereMachineProviderSpec).Workspace
					for idx, workspace := range tc.workspaces {
						if reflect.DeepEqual(workspace, *machineWorkspace) {
							matchCountByIndex[idx]++
						}
					}
				}
				for _, count := range matchCountByIndex {
					if count > 1 {
						t.Errorf("zone was applied to machineset too many times[max: 1]")
					}
				}
			}
		})
	}
}
