// Package aws generates Machine objects for aws.
package aws

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/intstr"
	"sigs.k8s.io/cluster-api/pkg/apis/cluster/common"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/pkg/errors"
)

// MachineDeployments returns a list of machineDeployments for a machinepool.
func MachineDeployments(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]clusterapi.MachineDeployment, error) {
	if configPlatform := config.Platform.Name(); configPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != aws.Name {
		return nil, fmt.Errorf("non-AWS machine-pool: %q", poolPlatform)
	}
	clustername := config.ObjectMeta.Name
	platform := config.Platform.AWS
	mpool := pool.Platform.AWS
	azs := mpool.Zones

	total := int64(0)
	if pool.Replicas != nil {
		total = *pool.Replicas
	}
	numOfAZs := int64(len(azs))
	var machineDeployments []clusterapi.MachineDeployment
	for idx, az := range azs {
		replicas := int32(total / numOfAZs)
		if int64(idx) < total%numOfAZs {
			replicas++
		}

		provider, err := provider(clusterID, clustername, platform, mpool, osImage, idx, role, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}
		name := fmt.Sprintf("%s-%s-%s", clustername, pool.Name, az)
		// Rolling machines one by one
		maxSurge := intstr.FromInt(1)
		maxUnavailable := intstr.FromInt(0)
		minReadySeconds := int32(0)
		mdep := clusterapi.MachineDeployment{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "cluster.k8s.io/v1alpha1",
				Kind:       "MachineDeployment",
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
			Spec: clusterapi.MachineDeploymentSpec{
				MinReadySeconds: &minReadySeconds,
				Strategy: &clusterapi.MachineDeploymentStrategy{
					Type: common.RollingUpdateMachineDeploymentStrategyType,
					RollingUpdate: &clusterapi.MachineRollingUpdateDeployment{
						MaxSurge:       &maxSurge,
						MaxUnavailable: &maxUnavailable,
					},
				},
				Replicas: &replicas,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"sigs.k8s.io/cluster-api-machinedeployment": name,
						"sigs.k8s.io/cluster-api-cluster":           clustername,
					},
				},
				Template: clusterapi.MachineTemplateSpec{
					ObjectMeta: metav1.ObjectMeta{
						Labels: map[string]string{
							"sigs.k8s.io/cluster-api-machinedeployment": name,
							"sigs.k8s.io/cluster-api-cluster":           clustername,
							"sigs.k8s.io/cluster-api-machine-role":      role,
							"sigs.k8s.io/cluster-api-machine-type":      role,
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
		machineDeployments = append(machineDeployments, mdep)
	}

	return machineDeployments, nil
}
