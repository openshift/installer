package nutanix

import (
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/utils/ptr"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// TestMachinesWorkerGPUConfig verifies that when a worker MachinePool includes GPU
// configuration, the generated Machines have the GPU settings propagated to the
// NutanixMachineProviderConfig. Covers Polarion workitem OCP-75546.
func TestMachinesWorkerGPUConfig(t *testing.T) {
	clusterID := "test-cluster"
	role := "worker"
	osImage := "rhcos-nutanix"
	userDataSecret := "worker-user-data"

	platform := &nutanix.Platform{
		PrismElements: []nutanix.PrismElement{
			{
				Name: "PE1",
				UUID: "pe-uuid-1111-2222-3333",
			},
		},
		SubnetUUIDs: []string{"subnet-uuid-aaaa-bbbb"},
	}

	tests := []struct {
		name         string
		gpus         []machinev1.NutanixGPU
		replicas     int64
		expectedGPUs int
	}{
		{
			name: "worker with single GPU by name",
			gpus: []machinev1.NutanixGPU{
				{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("Tesla-T4")},
			},
			replicas:     2,
			expectedGPUs: 1,
		},
		{
			name: "worker with single GPU by deviceID",
			gpus: []machinev1.NutanixGPU{
				{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(4318))},
			},
			replicas:     2,
			expectedGPUs: 1,
		},
		{
			name: "worker with multiple GPUs",
			gpus: []machinev1.NutanixGPU{
				{Type: machinev1.NutanixGPUIdentifierName, Name: ptr.To("Tesla-T4")},
				{Type: machinev1.NutanixGPUIdentifierDeviceID, DeviceID: ptr.To(int32(4318))},
			},
			replicas:     2,
			expectedGPUs: 2,
		},
		{
			name:         "worker without GPUs",
			gpus:         nil,
			replicas:     2,
			expectedGPUs: 0,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			mpool := &nutanix.MachinePool{
				NumCPUs:           4,
				NumCoresPerSocket: 2,
				MemoryMiB:         16384,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 120,
				},
				GPUs: tc.gpus,
			}

			pool := &types.MachinePool{
				Name:     role,
				Replicas: ptr.To(tc.replicas),
				Platform: types.MachinePoolPlatform{
					Nutanix: mpool,
				},
			}

			ic := &types.InstallConfig{
				Platform: types.Platform{
					Nutanix: platform,
				},
			}

			machines, _, err := Machines(clusterID, ic, pool, osImage, role, userDataSecret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Len(t, machines, int(tc.replicas))

			for _, machine := range machines {
				providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig)
				if !ok {
					t.Fatal("failed to cast provider spec to NutanixMachineProviderConfig")
				}
				assert.Len(t, providerSpec.GPUs, tc.expectedGPUs, "GPU count mismatch in provider config")

				for i, gpu := range providerSpec.GPUs {
					assert.Equal(t, tc.gpus[i].Type, gpu.Type, "GPU type mismatch")
					switch gpu.Type {
					case machinev1.NutanixGPUIdentifierName:
						assert.Equal(t, tc.gpus[i].Name, gpu.Name, "GPU name mismatch")
					case machinev1.NutanixGPUIdentifierDeviceID:
						assert.Equal(t, tc.gpus[i].DeviceID, gpu.DeviceID, "GPU deviceID mismatch")
					}
				}
			}
		})
	}
}

// TestMachinesMultiSubnetConfig verifies that when a platform or failure domain
// has multiple SubnetUUIDs (multi-NIC), the generated Machines have all subnets
// propagated to the NutanixMachineProviderConfig. Covers Polarion workitem OCP-78555.
func TestMachinesMultiSubnetConfig(t *testing.T) {
	clusterID := "test-cluster"
	role := "worker"
	osImage := "rhcos-nutanix"
	userDataSecret := "worker-user-data"

	tests := []struct {
		name            string
		subnetUUIDs     []string
		failureDomains  []nutanix.FailureDomain
		mpoolFDs        []string
		expectedSubnets int
	}{
		{
			name:            "single subnet (default)",
			subnetUUIDs:     []string{"subnet-uuid-1"},
			expectedSubnets: 1,
		},
		{
			name:            "multiple subnets for multi-NIC",
			subnetUUIDs:     []string{"subnet-uuid-1", "subnet-uuid-2"},
			expectedSubnets: 2,
		},
		{
			name:            "three subnets for multi-NIC",
			subnetUUIDs:     []string{"subnet-uuid-1", "subnet-uuid-2", "subnet-uuid-3"},
			expectedSubnets: 3,
		},
		{
			name:        "multiple subnets from failure domain",
			subnetUUIDs: []string{"subnet-uuid-platform"},
			failureDomains: []nutanix.FailureDomain{
				{
					Name: "fd-multi-nic",
					PrismElement: nutanix.PrismElement{
						Name: "PE1",
						UUID: "pe-uuid-1111-2222-3333",
					},
					SubnetUUIDs: []string{"subnet-uuid-fd-1", "subnet-uuid-fd-2"},
				},
			},
			mpoolFDs:        []string{"fd-multi-nic"},
			expectedSubnets: 2,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			platform := &nutanix.Platform{
				PrismElements: []nutanix.PrismElement{
					{
						Name: "PE1",
						UUID: "pe-uuid-1111-2222-3333",
					},
				},
				SubnetUUIDs: tc.subnetUUIDs,
			}
			if tc.failureDomains != nil {
				platform.FailureDomains = tc.failureDomains
			}

			mpool := &nutanix.MachinePool{
				NumCPUs:           4,
				NumCoresPerSocket: 2,
				MemoryMiB:         16384,
				OSDisk: nutanix.OSDisk{
					DiskSizeGiB: 120,
				},
				FailureDomains: tc.mpoolFDs,
			}

			pool := &types.MachinePool{
				Name:     role,
				Replicas: ptr.To(int64(2)),
				Platform: types.MachinePoolPlatform{
					Nutanix: mpool,
				},
			}

			ic := &types.InstallConfig{
				Platform: types.Platform{
					Nutanix: platform,
				},
			}

			machines, _, err := Machines(clusterID, ic, pool, osImage, role, userDataSecret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Len(t, machines, 2)

			for _, machine := range machines {
				providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig)
				if !ok {
					t.Fatal("failed to cast provider spec to NutanixMachineProviderConfig")
				}
				assert.Len(t, providerSpec.Subnets, tc.expectedSubnets, "subnet count mismatch in provider config")

				for _, subnet := range providerSpec.Subnets {
					assert.Equal(t, machinev1.NutanixIdentifierUUID, subnet.Type, "subnet identifier type should be UUID")
					assert.NotNil(t, subnet.UUID, "subnet UUID must not be nil")
				}
			}
		})
	}
}

// TestMachinesControlPlaneMachineSet verifies that the Machines function generates
// a ControlPlaneMachineSet (CPMS) with the correct metadata, state, replicas, labels,
// and failure domain configuration. Covers Polarion workitem OCP-66571.
func TestMachinesControlPlaneMachineSet(t *testing.T) {
	clusterID := "test-cluster"
	role := "master"
	osImage := "rhcos-nutanix"
	userDataSecret := "master-user-data"

	basePlatform := &nutanix.Platform{
		PrismElements: []nutanix.PrismElement{
			{
				Name: "PE1",
				UUID: "pe-uuid-1111-2222-3333",
			},
		},
		SubnetUUIDs: []string{"subnet-uuid-aaaa-bbbb"},
	}

	baseMachinePool := &nutanix.MachinePool{
		NumCPUs:           4,
		NumCoresPerSocket: 2,
		MemoryMiB:         16384,
		OSDisk: nutanix.OSDisk{
			DiskSizeGiB: 120,
		},
	}

	tests := []struct {
		name                string
		replicas            int64
		failureDomains      []nutanix.FailureDomain
		mpoolFailureDomains []string
		expectFDCount       int
	}{
		{
			name:          "CPMS is generated with Active state and correct replicas without failure domains",
			replicas:      3,
			expectFDCount: 0,
		},
		{
			name:     "CPMS includes single failure domain reference for single-FD cluster",
			replicas: 3,
			failureDomains: []nutanix.FailureDomain{
				{
					Name: "fd-single",
					PrismElement: nutanix.PrismElement{
						Name: "PE1",
						UUID: "pe-uuid-1111-2222-3333",
					},
					SubnetUUIDs: []string{"subnet-uuid-fd-single"},
				},
			},
			mpoolFailureDomains: []string{"fd-single"},
			expectFDCount:       1,
		},
		{
			name:     "CPMS includes failure domain references when failure domains are configured",
			replicas: 3,
			failureDomains: []nutanix.FailureDomain{
				{
					Name: "fd-1",
					PrismElement: nutanix.PrismElement{
						Name: "PE1",
						UUID: "pe-uuid-1111-2222-3333",
					},
					SubnetUUIDs: []string{"subnet-uuid-fd1"},
				},
				{
					Name: "fd-2",
					PrismElement: nutanix.PrismElement{
						Name: "PE2",
						UUID: "pe-uuid-4444-5555-6666",
					},
					SubnetUUIDs: []string{"subnet-uuid-fd2"},
				},
				{
					Name: "fd-3",
					PrismElement: nutanix.PrismElement{
						Name: "PE3",
						UUID: "pe-uuid-7777-8888-9999",
					},
					SubnetUUIDs: []string{"subnet-uuid-fd3"},
				},
			},
			mpoolFailureDomains: []string{"fd-1", "fd-2", "fd-3"},
			expectFDCount:       3,
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			platform := basePlatform.DeepCopy()
			if tc.failureDomains != nil {
				platform.FailureDomains = tc.failureDomains
			}

			mpool := &nutanix.MachinePool{
				NumCPUs:           baseMachinePool.NumCPUs,
				NumCoresPerSocket: baseMachinePool.NumCoresPerSocket,
				MemoryMiB:         baseMachinePool.MemoryMiB,
				OSDisk:            baseMachinePool.OSDisk,
				FailureDomains:    tc.mpoolFailureDomains,
			}

			pool := &types.MachinePool{
				Name:     role,
				Replicas: ptr.To(tc.replicas),
				Platform: types.MachinePoolPlatform{
					Nutanix: mpool,
				},
			}

			ic := &types.InstallConfig{
				Platform: types.Platform{
					Nutanix: platform,
				},
			}

			machines, cpms, err := Machines(clusterID, ic, pool, osImage, role, userDataSecret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			// Verify machines were generated
			assert.Len(t, machines, int(tc.replicas))

			// Verify CPMS is not nil
			if cpms == nil {
				t.Fatal("ControlPlaneMachineSet must be generated")
			}

			// Verify CPMS state is Active
			assert.Equal(t, machinev1.ControlPlaneMachineSetStateActive, cpms.Spec.State)

			// Verify CPMS metadata
			assert.Equal(t, "cluster", cpms.Name)
			assert.Equal(t, "openshift-machine-api", cpms.Namespace)
			assert.Equal(t, "ControlPlaneMachineSet", cpms.Kind)
			assert.Equal(t, "machine.openshift.io/v1", cpms.APIVersion)

			// Verify replicas
			if cpms.Spec.Replicas == nil {
				t.Fatal("CPMS replicas must not be nil")
			}
			assert.Equal(t, int32(tc.replicas), *cpms.Spec.Replicas)

			// Verify labels
			assert.Equal(t, clusterID, cpms.Labels["machine.openshift.io/cluster-api-cluster"])

			// Verify selector labels
			assert.Equal(t, role, cpms.Spec.Selector.MatchLabels["machine.openshift.io/cluster-api-machine-role"])
			assert.Equal(t, role, cpms.Spec.Selector.MatchLabels["machine.openshift.io/cluster-api-machine-type"])
			assert.Equal(t, clusterID, cpms.Spec.Selector.MatchLabels["machine.openshift.io/cluster-api-cluster"])

			// Verify template machine type
			assert.Equal(t, machinev1.OpenShiftMachineV1Beta1MachineType, cpms.Spec.Template.MachineType)

			// Verify template has machine spec with provider
			if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine == nil {
				t.Fatal("CPMS template OpenShiftMachineV1Beta1Machine must not be nil")
			}
			templateLabels := cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.ObjectMeta.Labels
			assert.Equal(t, clusterID, templateLabels["machine.openshift.io/cluster-api-cluster"])
			assert.Equal(t, role, templateLabels["machine.openshift.io/cluster-api-machine-role"])
			assert.Equal(t, role, templateLabels["machine.openshift.io/cluster-api-machine-type"])

			// Verify provider spec is set in template
			assert.NotNil(t, cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value)

			// Verify failure domains
			if tc.expectFDCount > 0 {
				if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains == nil {
					t.Fatal("CPMS failure domains must not be nil when failure domains are configured")
				}
				fds := cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains
				assert.Len(t, fds.Nutanix, tc.expectFDCount)
				for i, fdRef := range fds.Nutanix {
					assert.Equal(t, tc.mpoolFailureDomains[i], fdRef.Name)
				}
			} else if cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains != nil {
				assert.Empty(t, cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.FailureDomains.Nutanix)
			}
		})
	}
}

// TestMachinesPreloadedOSImageName verifies that when a preloaded OS image name
// is passed as the osImage parameter, the generated Machines reference that image
// name in their NutanixMachineProviderConfig. Covers Polarion workitem OCP-78659.
func TestMachinesPreloadedOSImageName(t *testing.T) {
	clusterID := "test-cluster"
	userDataSecret := "master-user-data"

	platform := &nutanix.Platform{
		PrismElements: []nutanix.PrismElement{
			{
				Name: "PE1",
				UUID: "pe-uuid-1111-2222-3333",
			},
		},
		SubnetUUIDs: []string{"subnet-uuid-aaaa-bbbb"},
	}

	mpool := &nutanix.MachinePool{
		NumCPUs:           8,
		NumCoresPerSocket: 2,
		MemoryMiB:         16384,
		OSDisk: nutanix.OSDisk{
			DiskSizeGiB: 120,
		},
	}

	tests := []struct {
		name          string
		osImage       string
		role          string
		expectedImage string
	}{
		{
			name:          "preloaded OS image name used for master machines",
			osImage:       "shared-rhcos-4.16",
			role:          "master",
			expectedImage: "shared-rhcos-4.16",
		},
		{
			name:          "preloaded OS image name used for worker machines",
			osImage:       "my-preloaded-rhcos",
			role:          "worker",
			expectedImage: "my-preloaded-rhcos",
		},
	}

	for _, tc := range tests {
		t.Run(tc.name, func(t *testing.T) {
			pool := &types.MachinePool{
				Name:     tc.role,
				Replicas: ptr.To(int64(3)),
				Platform: types.MachinePoolPlatform{
					Nutanix: mpool,
				},
			}

			ic := &types.InstallConfig{
				Platform: types.Platform{
					Nutanix: platform,
				},
			}

			machines, _, err := Machines(clusterID, ic, pool, tc.osImage, tc.role, userDataSecret)
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}

			assert.Len(t, machines, 3)

			for _, machine := range machines {
				providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.NutanixMachineProviderConfig)
				if !ok {
					t.Fatal("failed to cast provider spec to NutanixMachineProviderConfig")
				}
				assert.Equal(t, machinev1.NutanixIdentifierName, providerSpec.Image.Type, "image identifier type should be Name")
				assert.NotNil(t, providerSpec.Image.Name, "image name must not be nil")
				assert.Equal(t, tc.expectedImage, *providerSpec.Image.Name, "image name should match the preloaded OS image name")
			}
		})
	}
}
