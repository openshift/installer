// Package powervs generates Machine objects for IBM Power VS.
package powervs

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/powervs"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != powervs.Name {
		return nil, fmt.Errorf("non-powerVS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != powervs.Name {
		return nil, fmt.Errorf("non-powerVS machine-pool: %q", poolPlatform)
	}

	platform := config.Platform.PowerVS
	mpool := pool.Platform.PowerVS
	var network string
	image := fmt.Sprintf("rhcos-%s", clusterID)

	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}
	var machinesets []*machineapi.MachineSet
	provider, err := provider(clusterID, platform, mpool, userDataSecret, image, network)
	if err != nil {
		return nil, fmt.Errorf("failed to create provider: %w", err)
	}

	name := fmt.Sprintf("%s-%s", clusterID, pool.Name)
	mset := &machineapi.MachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1beta1",
			Kind:       "MachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      name,
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machineapi.MachineSetSpec{
			Replicas: &total,
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
				},
			},
		},
	}
	machinesets = append(machinesets, mset)

	return machinesets, nil
}
