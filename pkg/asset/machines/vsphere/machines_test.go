package vsphere

import (
	"github.com/ghodss/yaml"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/vsphere"
	vsphereapis "github.com/openshift/machine-api-operator/pkg/apis/vsphereprovider/v1beta1"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"reflect"
	"testing"
)

var installConfigSample = `
apiVersion: v1
baseDomain: example.com
controlPlane:
  name: master
  platform:
    vsphere:
      zones:
      - region-dc1-zone-c1
      - region-dc1-zone-c2
      - region-dc1-zone-c3
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
      - region-dc1-zone-c1
      - region-dc1-zone-c2
      - region-dc1-zone-c3
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
    - server: your.first-vcenter.example.com
      user: username
      password: password
      regions:
      - name: region-dc1
        datacenter: dc1
        vcenter: your.first-vcenter.example.com
        zones:
        - name: region-dc1-zone-c1
          cluster: DC1_C1
          resourcepool: /DC1/host/DC1_C1/Resources
          network: "vm network"
          datastore: "DC1_DS1"
        - name: region-dc1-zone-c2
          cluster: DC1_C2
          resourcepool: /DC1/host/DC1_C2/Resources
          network: "vm network"
          datastore: "DC1_DS2"
        - name: region-dc1-zone-c3
          cluster: DC1_C3
          resourcepool: /DC1/host/DC1_C3/Resources
          network: "vm network"
          datastore: "DC1_DS3"
    - name: your.second-vcenter.example.com
      port: 8443
      user: username
      password: password
      regions: []
pullSecret:
sshKey:`

var machinePoolReplicas = int64(3)

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
				"region-dc1-zone-c1",
				"region-dc1-zone-c2",
				"region-dc1-zone-c3",
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
				"region-dc1-zone-c1",
				"region-dc1-zone-c2",
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
	testCases := []struct {
		testCase                   string
		expectedError              string
		machinePool                *types.MachinePool
		workspaces                 []vsphereapis.Workspace
		maxAllowedWorkspaceMatches int
		minAllowedWorkspaceMatches int
	}{
		{
			testCase:                   "machinepool with no defined zones should use legacy configuration",
			machinePool:                &machinePoolNoZones,
			minAllowedWorkspaceMatches: 3,
			maxAllowedWorkspaceMatches: 3,
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
			testCase:                   "zones distributed among control plane machines(zone count matches machine count)",
			machinePool:                &machinePoolValidZones,
			maxAllowedWorkspaceMatches: 1,
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
			machinePool:   &machinePoolUndefinedZones,
			expectedError: "zone [region-dc1-zone-undefined] specified by machinepool is not defined",
		},
		{
			testCase:                   "zones distributed among control plane machines(zone less than machine count)",
			machinePool:                &machinePoolReducedZones,
			maxAllowedWorkspaceMatches: 2,
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
			machines, err := Machines(clusterID, installConfig, tc.machinePool, "", "", "")
			assertOnUnexpectedErrorState(tc.expectedError, err, t)

			if len(tc.workspaces) > 0 {
				var matchCountByIndex []int
				for range tc.workspaces {
					matchCountByIndex = append(matchCountByIndex, 0)
				}

				for _, machine := range machines {
					machineWorkspace := machine.Spec.ProviderSpec.Value.Object.(*vsphereapis.VSphereMachineProviderSpec).Workspace
					for idx, workspace := range tc.workspaces {
						if reflect.DeepEqual(workspace, *machineWorkspace) {
							matchCountByIndex[idx]++
						}
					}
				}
				for _, count := range matchCountByIndex {
					if count > tc.maxAllowedWorkspaceMatches {
						t.Errorf("zone was applied to machine too many times[max: %d]", tc.maxAllowedWorkspaceMatches)
					}
					if count < tc.minAllowedWorkspaceMatches {
						t.Errorf("zone was applied to machine too few times[min: %d]", tc.minAllowedWorkspaceMatches)
					}
				}
			}
		})
	}
}
