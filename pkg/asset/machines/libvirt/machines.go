// Package libvirt generates Machine objects for libvirt.
package libvirt

import (
	"fmt"

	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1beta1"
	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/libvirt"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != libvirt.Name {
		return nil, fmt.Errorf("non-Libvirt configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != libvirt.Name {
		return nil, fmt.Errorf("non-Libvirt machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Libvirt

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(clusterID, config.Networking.MachineNetwork[0].CIDR.String(), platform, userDataSecret)
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

func provider(clusterID string, networkInterfaceAddress string, platform *libvirt.Platform, userDataSecret string) *libvirtprovider.LibvirtMachineProviderConfig {
	return &libvirtprovider.LibvirtMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "libvirtproviderconfig.openshift.io/v1beta1",
			Kind:       "LibvirtMachineProviderConfig",
		},
		DomainMemory: 8192,
		DomainVcpu:   4,
		Ignition: &libvirtprovider.Ignition{
			UserDataSecret: userDataSecret,
		},
		Volume: &libvirtprovider.Volume{
			PoolName:     clusterID,
			BaseVolumeID: fmt.Sprintf("%s-base", clusterID),
		},
		NetworkInterfaceName:    clusterID,
		NetworkInterfaceAddress: networkInterfaceAddress,
		Autostart:               false,
		URI:                     platform.URI,
	}
}
