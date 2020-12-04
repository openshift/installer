// Package ovirt generates Machine objects for ovirt.
package ovirt

import (
	"fmt"

	machineapi "github.com/openshift/machine-api-operator/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/ovirt"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role,
	userDataSecret string) ([]*machineapi.MachineSet, error) {

	if configPlatform := config.Platform.Name(); configPlatform != ovirt.Name {
		return nil, fmt.Errorf("non-ovirt configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != "" && poolPlatform != ovirt.Name {
		return nil, fmt.Errorf("non-ovirt machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Ovirt

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	provider := provider(platform, pool, userDataSecret, osImage)
	name := fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, 0)
	mset := &machineapi.MachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "MachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      name,
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster":      clusterID,
				"machine.openshift.io/cluster-api-machine-role": role,
				"machine.openshift.io/cluster-api-machine-type": role,
			},
		},
		Spec: machineapi.MachineSetSpec{
			Replicas: pointer.Int32Ptr(int32(total)),
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machineset": name,
					"machine.openshift.io/cluster-api-cluster":    clusterID,
				},
			},
			Template: machineapi.MachineTemplateSpec{
				ObjectMeta: machineapi.ObjectMeta{
					Labels: map[string]string{
						"machine.openshift.io/cluster-api-machineset":   name,
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
			},
		},
	}

	return []*machineapi.MachineSet{mset}, nil
}
