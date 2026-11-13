// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	gcptypes "github.com/openshift/installer/pkg/types/gcp"
)

func TestConfigMasters(t *testing.T) {
	clusterID := "test"
	testCases := []struct {
		testCase            string
		publishingStrategy  types.PublishingStrategy
		expectedTargetPools []string
	}{
		{
			testCase:           "External",
			publishingStrategy: types.ExternalPublishingStrategy,
			expectedTargetPools: []string{
				fmt.Sprintf("%s-api", clusterID),
			},
		},
		{
			testCase:            "Internal",
			publishingStrategy:  types.InternalPublishingStrategy,
			expectedTargetPools: nil,
		},
	}

	for _, tc := range testCases {
		machines := []machineapi.Machine{
			{
				Spec: machineapi.MachineSpec{
					ProviderSpec: machineapi.ProviderSpec{
						Value: &runtime.RawExtension{Object: &machineapi.GCPMachineProviderSpec{}},
					},
				},
			},
			{
				Spec: machineapi.MachineSpec{
					ProviderSpec: machineapi.ProviderSpec{
						Value: &runtime.RawExtension{Object: &machineapi.GCPMachineProviderSpec{}},
					},
				},
			},
		}
		controlPlaneMachineSet := &machinev1.ControlPlaneMachineSet{
			Spec: machinev1.ControlPlaneMachineSetSpec{
				Template: machinev1.ControlPlaneMachineSetTemplate{
					OpenShiftMachineV1Beta1Machine: &machinev1.OpenShiftMachineV1Beta1MachineTemplate{
						Spec: machineapi.MachineSpec{
							ProviderSpec: machineapi.ProviderSpec{
								Value: &runtime.RawExtension{
									Object: &machineapi.GCPMachineProviderSpec{},
								},
							},
						},
					},
				},
			},
		}
		t.Run(tc.testCase, func(t *testing.T) {
			err := ConfigMasters(machines, controlPlaneMachineSet, clusterID, tc.publishingStrategy)
			assert.NoError(t, err)
			for _, machine := range machines {
				providerSpec := machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
				assert.Equal(t, providerSpec.TargetPools, tc.expectedTargetPools)
			}
		})
	}
}

func TestMachinesBootImages(t *testing.T) {
	const defaultImage = "projects/rhcos/global/images/default-image"
	cases := []struct {
		name             string
		customImage      *gcptypes.OSImage
		wantMachineImage string
		wantTemplate     string
	}{
		{
			name:             "default image",
			wantMachineImage: defaultImage,
		},
		{
			name: "custom image",
			customImage: &gcptypes.OSImage{
				Project: "custom-project",
				Name:    "custom-image",
			},
			wantMachineImage: "projects/custom-project/global/images/custom-image",
			wantTemplate:     "projects/custom-project/global/images/custom-image",
		},
	}

	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			config, pool := gcpImageFixture("master")
			pool.Platform.GCP.OSImage = tc.customImage

			machines, cpms, err := Machines("cluster", config, pool, defaultImage, "master", "master-user-data")
			if !assert.NoError(t, err) || !assert.Len(t, machines, 2) || !assert.NotNil(t, cpms) {
				return
			}

			cpmsProvider := cpms.Spec.Template.OpenShiftMachineV1Beta1Machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
			assert.Equal(t, tc.wantTemplate, cpmsProvider.Disks[0].Image)
			for _, machine := range machines {
				provider := machine.Spec.ProviderSpec.Value.Object.(*machineapi.GCPMachineProviderSpec)
				assert.Equal(t, tc.wantMachineImage, provider.Disks[0].Image)
			}
		})
	}
}

func gcpImageFixture(role string) (*types.InstallConfig, *types.MachinePool) {
	config := &types.InstallConfig{
		Platform: types.Platform{GCP: &gcptypes.Platform{ProjectID: "test-project", Region: "us-central1"}},
	}
	pool := &types.MachinePool{
		Name:     role,
		Replicas: ptr.To(int64(2)),
		Platform: types.MachinePoolPlatform{GCP: &gcptypes.MachinePool{
			Zones:        []string{"us-central1-a", "us-central1-b"},
			InstanceType: "n2-standard-4",
		}},
	}
	return config, pool
}
