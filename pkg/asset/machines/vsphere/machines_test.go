package vsphere

import (
	"net/netip"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	"github.com/stretchr/testify/assert"
	"sigs.k8s.io/yaml"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/conversion"
	"github.com/openshift/installer/pkg/types/vsphere"
)

const testClusterID = "test"

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
    vcenters:
    - server: your.vcenter.example.com
      username: username
      password: password
      datacenters: 
      - dc1
      - dc2
      - dc3
      - dc4
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
        datastore: /dc1/datastore/datastore1
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
        datastore: /dc2/datastore/datastore2
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
        datastore: /dc3/datastore/datastore3
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
        datastore: /dc4/datastore/datastore4
    hosts:
    - role: bootstrap
      networkDevice:
        ipaddrs:
        - 192.168.101.240/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
    - role: control-plane
      failureDomain: deployzone-us-east-1a
      networkDevice:
        ipAddrs:
        - 192.168.101.241/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
    - role: control-plane
      failureDomain: deployzone-us-east-2a
      networkDevice:
        ipAddrs:
        - 192.168.101.242/24
        gateway: 192.168.101.1
        nameservers:
        - 192.168.101.2
    - role: control-plane
      failureDomain: deployzone-us-east-3a
      networkDevice:
        ipAddrs:
        - 192.168.101.243/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
    - role: compute
      failureDomain: deployzone-us-east-1a
      networkDevice:
        ipAddrs:
        - 192.168.101.244/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
    - role: compute
      failureDomain: deployzone-us-east-2a
      networkDevice:
        ipAddrs:
        - 192.168.101.245/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
    - role: compute
      failureDomain: deployzone-us-east-3a
      networkDevice:
        ipAddrs:
        - 192.168.101.246/24
        gateway: 192.168.101.1        
        nameservers:
        - 192.168.101.2
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

var machinePoolSingleZones = types.MachinePool{
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

func assertOnUnexpectedErrorState(t *testing.T, expectedError string, err error) {
	t.Helper()
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
	clusterID := testClusterID
	installConfig, err := parseInstallConfig()
	if err != nil {
		assert.Errorf(t, err, "unable to parse sample install config")
		return
	}
	defaultClusterResourcePool, err := parseInstallConfig()
	if err != nil {
		t.Error(err)
		return
	}
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
			testCase:                   "zones distributed among control plane machines(zone count matches machine count)",
			machinePool:                &machinePoolValidZones,
			maxAllowedWorkspaceMatches: 1,
			installConfig:              installConfig,
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
			testCase:                   "all masters in single zone / pool",
			machinePool:                &machinePoolSingleZones,
			maxAllowedWorkspaceMatches: 3,
			installConfig:              installConfig,
			workspaces: []machineapi.Workspace{
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc1",
					Folder:       "/dc1/vm/folder1",
					Datastore:    "/dc1/datastore/datastore1",
					ResourcePool: "/dc1/host/c1/Resources/rp1",
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
			data, err := Machines(clusterID, tc.installConfig, tc.machinePool, "", "", "")
			assertOnUnexpectedErrorState(t, tc.expectedError, err)

			if len(tc.workspaces) > 0 {
				var matchCountByIndex []int
				for range tc.workspaces {
					matchCountByIndex = append(matchCountByIndex, 0)
				}

				for _, machine := range data.Machines {
					// check if expected workspaces are returned
					machineWorkspace := machine.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec).Workspace
					for idx, workspace := range tc.workspaces {
						if cmp.Equal(workspace, *machineWorkspace) {
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

				if data.ControlPlaneMachineSet != nil {
					// Make sure FDs equal same quantity as config
					fds := data.ControlPlaneMachineSet.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains
					if len(fds.VSphere) != len(tc.workspaces) {
						t.Errorf("machine workspace count %d does not equal number of failure domains [count: %d] in CPMS", len(tc.workspaces), len(fds.VSphere))
					}
				}
			}
		})
	}
}

func TestHostsToMachines(t *testing.T) {
	clusterID := testClusterID
	installConfig, err := parseInstallConfig()
	if err != nil {
		assert.Errorf(t, err, "unable to parse sample install config")
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
		installConfig *types.InstallConfig
		role          string
		machineCount  int
	}{
		{
			testCase:      "Static IP - ControlPlane",
			machinePool:   &machinePoolValidZones,
			installConfig: installConfig,
			role:          "master",
			machineCount:  3,
		},
		{
			testCase:      "Static IP - Compute",
			machinePool:   &machinePoolValidZones,
			installConfig: installConfig,
			role:          "worker",
			machineCount:  3,
		},
	}

	for _, tc := range testCases {
		t.Run(tc.testCase, func(t *testing.T) {
			data, err := Machines(clusterID, tc.installConfig, tc.machinePool, "", tc.role, "")
			assertOnUnexpectedErrorState(t, tc.expectedError, err)

			// Check machine counts
			if len(data.Machines) != tc.machineCount {
				t.Errorf("machine count (%v) did not match expected (%v).", len(data.Machines), tc.machineCount)
			}

			// Check Claim counts
			if len(data.IPClaims) != tc.machineCount {
				t.Errorf("ip address claim count (%v) did not match expected (%v).", len(data.IPClaims), tc.machineCount)
			}

			// Check ip address counts
			if len(data.IPAddresses) != tc.machineCount {
				t.Errorf("ip address count (%v) did not match expected (%v).", len(data.IPAddresses), tc.machineCount)
			}

			// Verify static IP was set on all machines
			for index, machine := range data.Machines {
				provider, success := machine.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)
				if !success {
					t.Errorf("Unable to convert vshere machine provider spec.")
				}

				if len(provider.Network.Devices) == 1 {
					// Check IP
					if provider.Network.Devices[0].AddressesFromPools == nil {
						t.Errorf("AddressesFromPools is not set: %v", machine)
					}

					// Check nameserver
					if provider.Network.Devices[0].Nameservers == nil || provider.Network.Devices[0].Nameservers[0] == "" {
						t.Errorf("Nameserver is not set: %v", machine)
					}

					gateway := data.IPAddresses[index].Spec.Gateway
					ip, err := netip.ParseAddr(gateway)
					if err != nil {
						t.Error(err)
					}
					targetGateway := "192.168.101.1"
					if ip.Is6() {
						targetGateway = "2001:0db8:85a3:0000:0000:8a2e:0370:7334"
					}

					// Check gateway
					if gateway != targetGateway {
						t.Errorf("Gateway is incorrect: %v", gateway)
					}
				} else {
					t.Errorf("Devices not set for machine: %v", machine)
				}
			}
		})
	}
}
