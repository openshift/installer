package vsphere

import (
	"fmt"
	"testing"

	"github.com/google/go-cmp/cmp"
	"github.com/stretchr/testify/assert"

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

var machineComputePoolSingleZone = types.MachinePool{
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
			testCase:      "all workers in single zone",
			machinePool:   &machineComputePoolSingleZone,
			maxReplicas:   3,
			minReplicas:   3,
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
			},
		},
		{
			testCase:      "undefined zone in machinepool results in error",
			machinePool:   &machineComputePoolUndefinedZones,
			installConfig: installConfig,
			expectedError: "zone [region-dc1-zone-undefined] specified by machinepool is not defined",
		},
		{
			testCase: "no zones specified uses all failure domains",
			machinePool: &types.MachinePool{
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
			},
			maxReplicas:   1,
			minReplicas:   0,
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
				{
					Server:       "your.vcenter.example.com",
					Datacenter:   "dc4",
					Folder:       "/dc4/vm/folder4",
					Datastore:    "/dc4/datastore/datastore4",
					ResourcePool: "/dc4/host/c4/Resources/rp4",
				},
			},
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
			machineSets, err := MachineSets(clusterID, tc.installConfig, tc.machinePool, "", "")
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

func TestMachineSetsStaticIP(t *testing.T) {
	clusterID := "test"
	installConfig, err := parseInstallConfig()
	if err != nil {
		t.Fatalf("unable to parse sample install config. %s", err.Error())
	}

	machineSets, err := MachineSets(clusterID, installConfig, &machineComputePoolValidZones, "worker", "")
	if err != nil {
		t.Fatalf("unexpected error: %s", err.Error())
	}

	if len(machineSets) == 0 {
		t.Fatal("expected at least one machineset")
	}

	for _, machineSet := range machineSets {
		provider := machineSet.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)

		if len(provider.Network.Devices) == 0 {
			t.Fatalf("machineset %s: expected at least one network device", machineSet.Name)
		}

		for idx, device := range provider.Network.Devices {
			if len(device.AddressesFromPools) == 0 {
				t.Errorf("machineset %s: device %d missing AddressesFromPools for static IP", machineSet.Name, idx)
				continue
			}
			pool := device.AddressesFromPools[0]
			expectedName := fmt.Sprintf("default-%d", idx)
			if pool.Group != "installer.openshift.io" {
				t.Errorf("machineset %s: device %d AddressesFromPools group = %s, want installer.openshift.io", machineSet.Name, idx, pool.Group)
			}
			if pool.Name != expectedName {
				t.Errorf("machineset %s: device %d AddressesFromPools name = %s, want %s", machineSet.Name, idx, pool.Name, expectedName)
			}
			if pool.Resource != "IPPool" {
				t.Errorf("machineset %s: device %d AddressesFromPools resource = %s, want IPPool", machineSet.Name, idx, pool.Resource)
			}

			if len(device.Nameservers) == 0 {
				t.Errorf("machineset %s: device %d missing Nameservers for static IP", machineSet.Name, idx)
			}
		}
	}
}

func TestMachineSetsMultiVCenter(t *testing.T) {
	clusterID := "test"

	config, err := parseInstallConfig()
	if err != nil {
		t.Fatalf("unable to parse sample install config: %v", err)
	}
	config.VSphere.Hosts = nil

	// Add a second vCenter
	config.VSphere.VCenters = append(config.VSphere.VCenters, vsphere.VCenter{
		Server:      "second-vcenter.example.com",
		Port:        443,
		Username:    "user2",
		Password:    "pass2",
		Datacenters: []string{"dc3"},
	})

	// Point the third failure domain to the second vCenter
	for i := range config.VSphere.FailureDomains {
		if config.VSphere.FailureDomains[i].Name == "deployzone-us-east-3a" {
			config.VSphere.FailureDomains[i].Server = "second-vcenter.example.com"
		}
	}

	expectedServers := map[string]string{
		"deployzone-us-east-1a": "your.vcenter.example.com",
		"deployzone-us-east-2a": "your.vcenter.example.com",
		"deployzone-us-east-3a": "second-vcenter.example.com",
	}

	machineSets, err := MachineSets(clusterID, config, &machineComputePoolValidZones, "worker", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Len(t, machineSets, len(expectedServers), "should produce one machineset per zone")

	for _, ms := range machineSets {
		providerSpec := ms.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)
		zone := providerSpec.Workspace.Datacenter

		// Map datacenter back to zone name for lookup
		dcToZone := map[string]string{
			"dc1": "deployzone-us-east-1a",
			"dc2": "deployzone-us-east-2a",
			"dc3": "deployzone-us-east-3a",
		}
		zoneName, ok := dcToZone[zone]
		if !ok {
			t.Fatalf("machineset %s: unexpected datacenter %q not in dcToZone map", ms.Name, zone)
		}
		expectedServer, ok := expectedServers[zoneName]
		if !ok {
			t.Fatalf("machineset %s: zone %q not in expectedServers map", ms.Name, zoneName)
		}

		assert.Equal(t, expectedServer, providerSpec.Workspace.Server,
			"machineset %s (zone %s): workspace server should match failure domain's vCenter", ms.Name, zoneName)
	}
}

func TestMachineSetsHostGroupVMGroup(t *testing.T) {
	clusterID := "test"

	config, err := parseInstallConfig()
	if err != nil {
		t.Fatalf("unable to parse sample install config: %v", err)
	}
	config.VSphere.Hosts = nil

	// Configure failure domains with HostGroup zone type
	for i := range config.VSphere.FailureDomains {
		config.VSphere.FailureDomains[i].RegionType = vsphere.ComputeClusterFailureDomain
		config.VSphere.FailureDomains[i].ZoneType = vsphere.HostGroupFailureDomain
		config.VSphere.FailureDomains[i].Topology.HostGroup = fmt.Sprintf("host-group-%d", i)
	}

	machineSets, err := MachineSets(clusterID, config, &machineComputePoolValidZones, "worker", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Len(t, machineSets, 3, "should produce one machineset per zone")

	// Map datacenter back to failure domain name to verify VMGroup
	dcToFD := map[string]string{
		"dc1": "deployzone-us-east-1a",
		"dc2": "deployzone-us-east-2a",
		"dc3": "deployzone-us-east-3a",
	}

	for _, ms := range machineSets {
		providerSpec := ms.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)
		fdName, ok := dcToFD[providerSpec.Workspace.Datacenter]
		if !ok {
			t.Fatalf("machineset %s: unexpected datacenter %q not in dcToFD map", ms.Name, providerSpec.Workspace.Datacenter)
		}

		expectedVMGroup := fmt.Sprintf("%s-%s", clusterID, fdName)
		assert.Equal(t, expectedVMGroup, providerSpec.Workspace.VMGroup,
			"machineset %s (fd %s): workspace VMGroup should be <clusterID>-<failureDomainName>", ms.Name, fdName)
	}
}

func TestMachineSetsMultiZoneTemplate(t *testing.T) {
	clusterID := "test"
	userTemplate := "/dc1/vm/rhcos-414.92.202306141028-0-vmware.x86_64.ova"

	config, err := parseInstallConfig()
	if err != nil {
		t.Fatalf("unable to parse sample install config: %v", err)
	}
	config.VSphere.Hosts = nil

	// Set template only on the first failure domain; leave others empty
	config.VSphere.FailureDomains[0].Topology.Template = userTemplate

	expectedTemplates := map[string]string{
		"dc1": userTemplate,
		"dc2": fmt.Sprintf("%s-rhcos-%s", clusterID, "deployzone-us-east-2a"),
		"dc3": fmt.Sprintf("%s-rhcos-%s", clusterID, "deployzone-us-east-3a"),
	}

	machineSets, err := MachineSets(clusterID, config, &machineComputePoolValidZones, "worker", "")
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}

	assert.Len(t, machineSets, 3, "should produce one machineset per zone")

	for _, ms := range machineSets {
		providerSpec := ms.Spec.Template.Spec.ProviderSpec.Value.Object.(*machineapi.VSphereMachineProviderSpec)
		dc := providerSpec.Workspace.Datacenter
		expectedTemplate := expectedTemplates[dc]

		assert.Equal(t, expectedTemplate, providerSpec.Template,
			"machineset %s (dc %s): template should match expected value", ms.Name, dc)
	}
}
