// Package kubevirt generates Machine objects for kubevirt.
package kubevirt

import (
	"fmt"

	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/kubevirt"

	kubevirtprovider "github.com/openshift/cluster-api-provider-kubevirt/pkg/apis/kubevirtprovider/v1alpha1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != kubevirt.Name {
		return nil, fmt.Errorf("non-kubevirt configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != kubevirt.Name {
		return nil, fmt.Errorf("non-kubevirt machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Kubevirt

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(clusterID, platform, pool, userDataSecret, config)
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

func provider(clusterID string, platform *kubevirt.Platform, pool *types.MachinePool, userDataSecret string, config *types.InstallConfig) *kubevirtprovider.KubevirtMachineProviderSpec {
	interfaceBindingMethod := "InterfaceBridge"
	if config.Kubevirt.InterfaceBindingMethod != "" {
		interfaceBindingMethod = "Interface" + config.Kubevirt.InterfaceBindingMethod
	}
	spec := kubevirtprovider.KubevirtMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "kubevirtproviderconfig.openshift.io/v1alpha1",
			Kind:       "KubevirtMachineProviderSpec",
		},
		SourcePvcName:              fmt.Sprintf("%s-source-pvc", clusterID),
		RequestedMemory:            pool.Platform.Kubevirt.Memory,
		RequestedCPU:               pool.Platform.Kubevirt.CPU,
		RequestedStorage:           pool.Platform.Kubevirt.StorageSize,
		StorageClassName:           platform.StorageClass,
		IgnitionSecretName:         userDataSecret,
		NetworkName:                platform.NetworkName,
		InterfaceBindingMethod:     interfaceBindingMethod,
		PersistentVolumeAccessMode: platform.PersistentVolumeAccessMode,
	}
	return &spec
}
