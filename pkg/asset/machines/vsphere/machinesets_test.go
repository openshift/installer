package vsphere

import (
	"testing"

	"github.com/google/go-cmp/cmp"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
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
				"deployzone-us-east-1a",
				"deployzone-us-east-2a",
				"deployzone-us-east-3a",
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
				"deployzone-us-east-1a",
				"deployzone-us-east-2a",
			},
		},
	},
	Hyperthreading: "true",
	Architecture:   types.ArchitectureAMD64,
}

var machineComputePoolReducedReplicas = types.MachinePool{
	Name:     "worker",
	Replicas: &machinePoolMinReplicas,
	Platform: types.MachinePoolPlatform{
		VSphere: &vsphere.MachinePool{
			NumCPUs:           4,
			NumCoresPerSocket: 2,
			MemoryMiB:         16384,
			OSDisk: vsphere.OSDisk{
				DiskSizeGB: 60,
			},
			Zones: []string{
				"deployzone-us-east-1a",
				"deployzone-us-east-2a",
				"deployzone-us-east-3a",
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

func TestConfigMachinesets(t *testing.T) {
	clusterID := "test"
	osImage := "test-cluster-xyzxyz-rhcos"
	installConfig, err := parseInstallConfig()

	if err != nil {
		t.Errorf("unable to parse sample install config. %s", err.Error())
		return
	}

	defaultClusterResourcePool, err := parseInstallConfig()
	if err != nil {
		t.Error(err)
		return
	}
	defaultClusterResourcePool.VSphere.FailureDomains[0].Topology.ResourcePool = ""

	testCases := []struct {
		testCase      string
		expectedError string
		machinePool   *types.MachinePool
		workspaces    []machineapi.Workspace
		maxReplicas   int
		minReplicas   int
		totalReplicas int
		installConfig *types.InstallConfig
	}{
		{
			testCase:      "zones distributed among compute machinesets(zone count matches machineset count)",
			machinePool:   &machineComputePoolValidZones,
			maxReplicas:   1,
			minReplicas:   1,
			totalReplicas: 3,
			installConfig: installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "/dc1/datastore/datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "/dc2/datastore/datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc3",
					Folder:       "/dc3/vm/folder3",
					Datastore:    "/dc3/datastore/datastore3",
					ResourcePool: "/dc3/host/c3/Resources/rp3",
				},
			},
		},
		{
			testCase:      "undefined zone in machinepool results in error",
			machinePool:   &machineComputePoolUndefinedZones,
			installConfig: installConfig,
			expectedError: "zone [region-dc1-zone-undefined] specified by machinepool is not defined",
		},
		{
			testCase:      "zones distributed among compute machinesets(zone count less than compute machinesets count)",
			machinePool:   &machineComputePoolReducedZones,
			maxReplicas:   2,
			minReplicas:   1,
			totalReplicas: 3,
			installConfig: installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "/dc1/datastore/datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "/dc2/datastore/datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
			},
		},
		{
			testCase:      "zones distributed among compute machinesets(machinesets count less than zone count)",
			machinePool:   &machineComputePoolReducedReplicas,
			maxReplicas:   1,
			minReplicas:   0,
			totalReplicas: 2,
			installConfig: installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "/dc1/datastore/datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "/dc2/datastore/datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
			},
		},
		{
			testCase:      "full path to cluster resource pool when no pool provided via placement constraint",
			machinePool:   &machinePoolValidZones,
			maxReplicas:   1,
			minReplicas:   0,
			totalReplicas: 3,
			installConfig: defaultClusterResourcePool,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "/dc1/datastore/datastore1",
					ResourcePool: "/dc1/host/c1/Resources",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "/dc2/datastore/datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc3",
					Folder:       "/dc3/vm/folder3",
					Datastore:    "/dc3/datastore/datastore3",
					ResourcePool: "/dc3/host/c3/Resources/rp3",
				},
			},
		},
	}
	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			machineSets, err := MachineSets(clusterID, tc.installConfig, tc.machinePool, osImage, "", "")
			assertOnUnexpectedErrorState(t, tc.expectedError, err)

			if len(tc.workspaces) > 0 {
				var matchCountByIndex []int
				for range tc.workspaces {
					matchCountByIndex = append(matchCountByIndex, 0)
				}
				replicaCount := 0
				for _, machineSet := range machineSets {
					// check if replica counts match expected
					replicas := int(*machineSet.Spec.Replicas)
					replicaCount += replicas
					if replicas > tc.maxReplicas {
						t.Errorf("machineset %s has too many replicas[max: %d] found %d", machineSet.Name, tc.maxReplicas, replicas)
					} else if replicas < tc.minReplicas {
						t.Errorf("machineset %s has too few replicas[max: %d] found %d", machineSet.Name, tc.maxReplicas, replicas)
					}

					// check if expected workspaces are returned
					machineWorkspace := machineSet.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec).Workspace
					for idx, workspace := range tc.workspaces {
						if cmp.Equal(workspace, *machineWorkspace) {
							matchCountByIndex[idx]++
						}
					}
				}
				if tc.totalReplicas != 0 && tc.totalReplicas != replicaCount {
					t.Errorf("expected replica count %d but %d replicas configured", tc.totalReplicas, replicaCount)
				}
				for _, count := range matchCountByIndex {
					if count > 1 {
						t.Errorf("expected machineset workspace encountered too many times[max: 1]")
					}
					if count == 0 {
						t.Errorf("expected machineset workspace was not applied to a zone")
					}
				}
			}
		})
	}
}
