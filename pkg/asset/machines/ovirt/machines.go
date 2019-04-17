// Package ovirt generates Machine objects for ovirt.
package ovirt

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"

	ovirtprovider "github.com/ovirt/cluster-api-provider-ovirt/pkg/apis/ovirtclusterproviderconfig/v1beta1"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != ovirt.Name {
		return nil, fmt.Errorf("non-ovirt configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != ovirt.Name {
		return nil, fmt.Errorf("non-ovirt machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Ovirt

	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	provider := provider(platform, userDataSecret)
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

func provider(platform *ovirt.Platform, userDataSecret string) *ovirtprovider.OvirtMachineProviderSpec {
	return &ovirtprovider.OvirtMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ovirtproviderconfig.openshift.io/v1beta1",
			Kind:       "OvirtMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "ovirt-credentials"},
		ClusterId:         platform.ClusterID,
		TemplateId:        platform.TemplateID,
	}
}
