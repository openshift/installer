// Package equinixmetal generates Machine objects for equinixmetal.
package equinixmetal

import (
	"fmt"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/pointer"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/equinixmetal"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role,
	userDataSecret string) ([]*machineapi.MachineSet, error) {

	if configPlatform := config.Platform.Name(); configPlatform != equinixmetal.Name {
		return nil, fmt.Errorf("non-equinixmetal configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != "" && poolPlatform != equinixmetal.Name {
		return nil, fmt.Errorf("non-equinixmetal machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.EquinixMetal
	mpool := pool.Platform.EquinixMetal
	plan := mpool.Plan
	customData := mpool.CustomData
	total := int64(1)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}

	provider := provider(clusterID, platform, plan, customData, role, userDataSecret, osImage)
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
				ObjectMeta: metav1.ObjectMeta{
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
