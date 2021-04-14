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
	mpool := pool.Platform.EquinixMetal

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	var machines []machineapi.Machine
	for idx := int64(0); idx < total; idx++ {
		provider := provider(clusterID, platform, mpool.Plan, mpool.CustomData, role, userDataSecret, osImage)

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

func provider(clusterID string, platform *equinixmetal.Platform, plan, customData, role, userDataSecret, osImage string) *equinixprovider.EquinixMetalMachineProviderConfig {
	// TOOD(displague) This IPXE script url contains the kernel and initrd
	// parameters needed to an official RHCOS image from the official mirror.
	// Equinix Metal devices can be created with userdata values of #!ipxe" to
	// avoid the need for a hosted IPXE script, but userdata must be a valid
	// Ignition Config for IPI purposes. I am actively seeking a EM feature to
	// permit ipxescripturl to support data urls or to offer an ipxescript
	// (content, not url) device creation field. Notably, this static script
	// should be dynamic based on the osImage parameter which may reflect a
	// differnet version and architecture than this gist offers.
	ipxeScriptURL := "https://gist.githubusercontent.com/displague/5282172449a83c7b83821f8f8333a072/raw/f7300a5ab652e923dddacb5c9f206864c4c2aceb/rhcos.ipxe"
	_ = osImage

	spec := equinixprovider.EquinixMetalMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: equinixprovider.SchemeGroupVersion.String(),
			Kind:       "EquinixMetalMachineProviderConfig",
		},
		CustomData: customData,
		// Facility:      platform.Facility,
		Metro:         platform.Metro,
		OS:            "custom_ipxe",
		ProjectID:     platform.ProjectID,
		IPXEScriptURL: ipxeScriptURL,
		BillingCycle:  "hourly",
		MachineType:   plan,
		Tags:          []string{"openshift-ipi", fmt.Sprintf("%s-%s", clusterID, role)},
		// TODO(displague) ssh keys will need to be defined in the project
		// SshKeys:      []string{},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "equinixmetal-credentials"},
	}
	return &spec
}
