// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	corev1 "k8s.io/api/core/v1"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(config *types.InstallConfig, pool *types.MachinePool, role string) ([]clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	mpool := pool.Platform.AWS
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machinesets []clusterapi.MachineSet
	for idx, az := range azs {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		name := fmt.Sprintf("%s-%s-%s", clustername, pool.Name, az)
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
				Replicas: &replicas,
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
							ValueFrom: &clusterapi.ProviderConfigSource{
								MachineClass: &clusterapi.MachineClassRef{
									ObjectReference: &corev1.ObjectReference{
										Kind:      "MachineClass",
										Name:      fmt.Sprintf("%s-%s-%s", clustername, role, az),
										Namespace: "openshift-cluster-api",
									},
								},
							},
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
