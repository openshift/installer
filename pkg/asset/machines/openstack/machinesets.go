// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	clusterapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string, trunkSupport bool) ([]*clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.OpenStack
	mpool := pool.Platform.OpenStack

	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}

	numOfAZs := int32(len(mpool.Zones))
	// TODO(flaper87): Add support for availability zones
	var machinesets []*clusterapi.MachineSet

	for idx, az := range mpool.Zones {
		replicas := int32(total / numOfAZs)
		if int32(idx) < total%numOfAZs {
			replicas++
		}
		provider, err := generateProvider(clusterID, platform, mpool, osImage, az, role, userDataSecret, trunkSupport)
		if err != nil {
			return nil, err
		}

		// TODO(flaper87): Implement AZ support sometime soon
		//name := fmt.Sprintf("%s-%s-%s", clustername, pool.Name, az)
		name := fmt.Sprintf("%s-%s", clusterID, pool.Name)
		mset := &clusterapi.MachineSet{
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
			Spec: clusterapi.MachineSetSpec{
				Replicas: &total,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Template: clusterapi.MachineTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
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
	}

	return machinesets, nil
}
