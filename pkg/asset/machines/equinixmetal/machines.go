// Package equinixmetal generates Machine objects for equinixmetal.
package equinixmetal

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/equinixmetal"

	equinixprovider "github.com/openshift/cluster-api-provider-equinix-metal/pkg/apis/equinixmetal/v1beta1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != equinixmetal.Name {
		return nil, fmt.Errorf("non-equinixmetal configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != equinixmetal.Name {
		return nil, fmt.Errorf("non-equinixmetal machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.EquinixMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(platform, pool, userDataSecret, osImage)
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels: map[string]string{
					"machine.openshift.io/cluster-api-cluster":      clusterID,
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
				},
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via cluster operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(platform *equinixmetal.Platform, pool *types.MachinePool, userDataSecret string, osImage string) *equinixprovider.EquinixMetalMachineProviderConfig {
	spec := equinixprovider.EquinixMetalMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "equinixprovider.k8s.io/v1alpha1",
			Kind:       "EquinixMetalMachineProviderSpec",
		},
		// TODO(displague) which role? doesn't matter - not present in latest providerconfig
		// Roles:     []equinixprovider.MachineRole{equinixprovider.MasterRole, equinixprovider.NodeRole},
		// TODO(displague) should Facility be offered as a slice? EM-API can pick from a list
		Facility:  platform.FacilityCode,
		OS:        "custom_ipxe",
		ProjectID: platform.ProjectID,
		// TODO(displague) IPXE / osImage / platform.BootstrapOSImage / platform.ClusterOSImage?
		IPXEScriptURL: osImage,
		BillingCycle:  "hourly",
		MachineType:   "t1.small.x86", // TODO(displague) must provide a type
		Tags:          []string{"openshift-ipi", "wip", "TODO"},
		// SshKeys:      []string{},
		UserDataSecret: &corev1.LocalObjectReference{Name: userDataSecret},
		// CredentialsSecret: &corev1.LocalObjectReference{Name: "equinixmetal-credentials"},
	}
	return &spec
}
