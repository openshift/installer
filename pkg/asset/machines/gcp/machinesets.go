// Package gcp generates Machine objects for gcp.
package gcp

import (
	"fmt"
	"strings"

	machineapi "github.com/openshift/cluster-api/pkg/apis/machine/v1beta1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/pkg/errors"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != gcp.Name {
		return nil, fmt.Errorf("non-GCP machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.GCP
	mpool := pool.Platform.GCP
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machinesets []*machineapi.MachineSet
	for idx, az := range azs {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		provider, err := provider(clusterID, platform, mpool, osImage, idx, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		name := fmt.Sprintf("%s-%s-%s", clusterID, pool.Name, strings.TrimPrefix(az, fmt.Sprintf("%s-", platform.Region)))
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
				Replicas: &replicas,
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
		machinesets = append(machinesets, mset)
	}

	return machinesets, nil
}
