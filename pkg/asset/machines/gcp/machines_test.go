// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"
	"testing"

	"github.com/stretchr/testify/assert"
	"k8s.io/apimachinery/pkg/runtime"

	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
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
