// Package libvirt generates Machine objects for libvirt.
package libvirt

import (
	"fmt"

	libvirtprovider "github.com/openshift/cluster-api-provider-libvirt/pkg/apis/libvirtproviderconfig/v1alpha1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/libvirt"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]clusterapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != libvirt.Name {
		return nil, fmt.Errorf("non-Libvirt configuration: %q", configPlatform)
	}
	// FIXME: empty is a valid case for Libvirt as we don't use it.
	if poolPlatform := pool.Platform.Name(); poolPlatform != "" && poolPlatform != libvirt.Name {
		return nil, fmt.Errorf("non-Libvirt machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.Libvirt
	// FIXME: libvirt actuator does not support any options from machinepool.
	// mpool := pool.Platform.Libvirt

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(clustername, config.Networking.MachineCIDR.String(), platform, userDataSecret)
	var machines []clusterapi.Machine
	for idx := int64(0); idx < total; idx++ {
		machine := clusterapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.k8s.io/v1alpha1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-cluster-api",
				Name:      fmt.Sprintf("%s-%s-%d", clustername, pool.Name, idx),
				Labels: map[string]string{
					"sigs.k8s.io/cluster-api-cluster":      clustername,
					"sigs.k8s.io/cluster-api-machine-role": role,
					"sigs.k8s.io/cluster-api-machine-type": role,
				},
			},
			Spec: clusterapi.MachineSpec{
				ProviderSpec: clusterapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via cluster operators.
			},
		}
		machines = append(machines, machine)
	}

	return machines, nil
}

func provider(clusterName string, networkInterfaceAddress string, platform *libvirt.Platform, userDataSecret string) *libvirtprovider.LibvirtMachineProviderConfig {
	return &libvirtprovider.LibvirtMachineProviderConfig{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "libvirtproviderconfig.k8s.io/v1alpha1",
			Kind:       "LibvirtMachineProviderConfig",
		},
		DomainMemory: 2048,
		DomainVcpu:   2,
		Ignition: &libvirtprovider.Ignition{
			UserDataSecret: userDataSecret,
		},
		Volume: &libvirtprovider.Volume{
			PoolName:     "default",
			BaseVolumeID: fmt.Sprintf("/var/lib/libvirt/images/%s-base", clusterName),
		},
		NetworkInterfaceName:    clusterName,
		NetworkInterfaceAddress: networkInterfaceAddress,
		Autostart:               false,
		URI:                     platform.URI,
	}
}
