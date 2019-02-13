// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/pkg/errors"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.OpenStack
	mpool := pool.Platform.OpenStack

	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}

	// TODO(flaper87): Add support for availability zones
	var machinesets []clusterapi.MachineSet
	az := ""
	provider, err := provider(clusterID, clustername, platform, mpool, osImage, az, role, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
	}
	// TODO(flaper87): Implement AZ support sometime soon
	//name := fmt.Sprintf("%s-%s-%s", clustername, pool.Name, az)
	name := fmt.Sprintf("%s-%s", clustername, pool.Name)
	mset := clusterapi.MachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "cluster.k8s.io/v1alpha1",
			Kind:       "MachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      name,
			Labels: map[string]string{
				"sigs.k8s.io/cluster-api-cluster":      clustername,
				"sigs.k8s.io/cluster-api-machine-role": role,
				"sigs.k8s.io/cluster-api-machine-type": role,
			},
		},
		Spec: clusterapi.MachineSetSpec{
			Replicas: &total,
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
					ProviderSpec: clusterapi.ProviderSpec{
						Value: &runtime.RawExtension{Object: provider},
					},
					// we don't need to set Versions, because we control those via cluster operators.
				},
			},
		},
	}
	machinesets = append(machinesets, mset)

	return machinesets, nil
}
