package vsphere

import (
	"reflect"
	"testing"

	"github.com/ghodss/yaml"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/vsphere"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
)

var installConfigSample = `
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    vsphere:
      zones:
      - deployzone-us-east-1a
      - deployzone-us-east-2a
      - deployzone-us-east-3a
      cpus: 8
      coresPerSocket: 2
      memoryMB: 24576
      osDisk:
        diskSizeGB: 512
  replicas: 3
compute:
- name: worker
  platform:
    vsphere:
      zones:
      - deployzone-us-east-1a
      - deployzone-us-east-2a
      - deployzone-us-east-3a
      cpus: 8
      coresPerSocket: 2
      memoryMB: 24576
      osDisk:
        diskSizeGB: 512
  replicas: 3
metadata:
  name: test-cluster
platform:
  vSphere:
    vCenter: your.vcenter.example.com
    username: username
    password: password
    datacenter: datacenter
    defaultDatastore: datastore
    network: portgroup
    vcenters:
    - server: your.vcenter.example.com
      username: username
      password: password
      datacenters: 
      - datacenter
    failureDomains:
    - name: deployzone-us-east-1a
      server: your.vcenter.example.com
      region: us-east
      zone: us-east-1a
      topology:
        resourcePool: /dc1/host/c1/Resources/rp1
        folder: /dc1/vm/folder1
        datacenter: dc1
        computeCluster: /dc1/host/c1
        networks:
        - network1
        datastore: datastore1
    - name: deployzone-us-east-2a
      server: your.vcenter.example.com
      region: us-east
      zone: us-east-2a
      topology:
        resourcePool: /dc2/host/c2/Resources/rp2
        folder: /dc2/vm/folder2
        datacenter: dc2
        computeCluster: /dc2/host/c2
        networks:
        - network1
        datastore: datastore2
    - name: deployzone-us-east-3a
      server: your.vcenter.example.com
      region: us-east
      zone: us-east-3a
      topology:
        resourcePool: /dc3/host/c3/Resources/rp3
        folder: /dc3/vm/folder3
        datacenter: dc3
        computeCluster: /dc3/host/c3
        networks:
        - network1
        datastore: datastore3
    - name: deployzone-us-east-4a
      server: your.vcenter.example.com 
      region: us-east
      zone: us-east-4a
      topology:
        resourcePool: /dc4/host/c4/Resources/rp4
        folder: /dc4/vm/folder4
        datacenter: dc4
        computeCluster: /dc4/host/c4
        networks:
        - network1
        datastore: datastore4
pullSecret:
sshKey:`

var machinePoolReplicas = int64(3)
var machinePoolMinReplicas = int64(2)

var machinePoolValidZones = types.MachinePool{
	Name:     "master",
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

var machinePoolReducedZones = types.MachinePool{
	Name:     "master",
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

var machinePoolUndefinedZones = types.MachinePool{
	Name:     "master",
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

var machinePoolNoZones = types.MachinePool{
	Name:     "master",
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

func parseInstallConfig() (*types.InstallConfig, error) {
	config := &types.InstallConfig{}
	if err := yaml.Unmarshal([]byte(installConfigSample), config); err != nil {
		return nil, errors.Wrapf(err, "failed to unmarshal sample install config")
	}

	// Upconvert any deprecated fields
	if err := conversion.ConvertInstallConfig(config); err != nil {
		return nil, errors.Wrap(err, "failed to upconvert install config")
	}
	return config, nil
}

func assertOnUnexpectedErrorState(expectedError string, err error, t *testing.T) {
	if expectedError == "" && err != nil {
		t.Errorf("unexpected error encountered: %s", err.Error())
	} else if expectedError != "" {
		if err != nil {
			if expectedError != err.Error() {
				t.Errorf("expected error does not match intended error. should be: %s was: %s", expectedError, err.Error())
			}
		} else {
			t.Errorf("expected error but error did not occur")
		}
	}
}

func TestConfigMasters(t *testing.T) {
	clusterID := "test"
	installConfig, err := parseInstallConfig()
	if err != nil {
		assert.Errorf(t, err, "unable to parse sample install config")
		return
	}
	defaultClusterResourcePool, err := parseInstallConfig()
	defaultClusterResourcePool.VSphere.FailureDomains[0].Topology.ResourcePool = ""

	testCases := []struct {
		testCase                   string
		expectedError              string
		machinePool                *types.MachinePool
		workspaces                 []machineapi.Workspace
		installConfig              *types.InstallConfig
		maxAllowedWorkspaceMatches int
		minAllowedWorkspaceMatches int
	}{
		{
			testCase:                   "machinepool with no defined zones should use legacy configuration",
			machinePool:                &machinePoolNoZones,
			minAllowedWorkspaceMatches: 3,
			maxAllowedWorkspaceMatches: 3,
			installConfig:              installConfig,
			workspaces: []machineapi.Workspace{
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
			testCase:                   "zones distributed among control plane machines(zone count matches machine count)",
			machinePool:                &machinePoolValidZones,
			maxAllowedWorkspaceMatches: 1,
			installConfig:              installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc3",
					Folder:       "/dc3/vm/folder3",
					Datastore:    "datastore3",
					ResourcePool: "/dc3/host/c3/Resources/rp3",
				},
			},
		},
		{
			testCase:      "undefined zone in machinepool results in error",
			machinePool:   &machinePoolUndefinedZones,
			installConfig: installConfig,
			expectedError: "zone [region-dc1-zone-undefined] specified by machinepool is not defined",
		},
		{
			testCase:                   "zones distributed among control plane machines(zone count is less than machine count)",
			machinePool:                &machinePoolReducedZones,
			maxAllowedWorkspaceMatches: 2,
			minAllowedWorkspaceMatches: 1,
			installConfig:              installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
			},
		},
		{
			testCase:                   "full path to cluster resource pool when no pool provided via placement constraint",
			machinePool:                &machinePoolValidZones,
			maxAllowedWorkspaceMatches: 1,
			minAllowedWorkspaceMatches: 1,
			installConfig:              defaultClusterResourcePool,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "datastore1",
					ResourcePool: "/dc1/host/c1/Resources",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc2",
					Folder:       "/dc2/vm/folder2",
					Datastore:    "datastore2",
					ResourcePool: "/dc2/host/c2/Resources/rp2",
				},
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc3",
					Folder:       "/dc3/vm/folder3",
					Datastore:    "datastore3",
					ResourcePool: "/dc3/host/c3/Resources/rp3",
				},
			},
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			machines, err := Machines(clusterID, tc.installConfig, tc.machinePool, "", "", "")
			assertOnUnexpectedErrorState(tc.expectedError, err, t)

			if len(tc.workspaces) > 0 {
				var matchCountByIndex []int
				for range tc.workspaces {
					matchCountByIndex = append(matchCountByIndex, 0)
				}

				for _, machine := range machines {
					// check if expected workspaces are returned
					machineWorkspace := machine.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec).Workspace
					for idx, workspace := range tc.workspaces {
						if reflect.DeepEqual(workspace, *machineWorkspace) {
							matchCountByIndex[idx]++
						}
					}
				}
				for _, count := range matchCountByIndex {
					if count > tc.maxAllowedWorkspaceMatches {
						t.Errorf("machine workspace was enountered too many times[max: %d]", tc.maxAllowedWorkspaceMatches)
					}
					if count < tc.minAllowedWorkspaceMatches {
						t.Errorf("machine workspace was enountered too few times[min: %d]", tc.minAllowedWorkspaceMatches)
					}
				}
			}
		})
	}
}
