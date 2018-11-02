// Package libvirt generates Machine objects for libvirt.
package libvirt

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != types.PlatformNameLibvirt {
		return nil, fmt.Errorf("non-Libvirt configuration: %q", configPlatform)
	}
	// FIXME: empty is a valid case for Libvirt as we don't use it.
	if poolPlatform := pool.Platform.Name(); poolPlatform != "" && poolPlatform != types.PlatformNameLibvirt {
		return nil, fmt.Errorf("non-Libvirt machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.Libvirt
	// FIXME: libvirt actuator does not support any options from machinepool.
	// mpool := pool.Platform.Libvirt

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	provider := provider(platform, pool.Name)
	name := fmt.Sprintf("%s-%s-%d", clustername, pool.Name, 0)
	mset := clusterapi.MachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cluster.k8s.io/v1alpha1",
			Kind:       "MachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-cluster-api",
			Name:      name,
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster":      clustername,
				"sigs.k8s.io/cluster-api-machine-role": role,
				"sigs.k8s.io/cluster-api-machine-type": role,
			},
		},
		Spec: clusterapi.MachineSetSpec{
			Replicas: pointer.Int32Ptr(int32(total)),
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"sigs.k8s.io/cluster-api-machineset": name,
					"sigs.k8s.io/cluster-api-cluster":    clustername,
				},
			},
			Template: clusterapi.MachineTemplateSpec{
				ObjectMeta: metav1.ObjectMeta{
					Labels: map[string]string{
						"sigs.k8s.io/cluster-api-machineset":   name,
						"sigs.k8s.io/cluster-api-cluster":      clustername,
						"sigs.k8s.io/cluster-api-machine-role": role,
						"sigs.k8s.io/cluster-api-machine-type": role,
					},
				},
				Spec: clusterapi.MachineSpec{
					ProviderConfig: clusterapi.ProviderConfig{
						Value: &runtime.RawExtension{Object: provider},
					},
					// we don't need to set Versions, because we control those via cluster operators.
				},
			},
		},
	}

	return []clusterapi.MachineSet{mset}, nil
}
