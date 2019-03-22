package azure

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
	clusterapi "sigs.k8s.io/cluster-api/pkg/apis/cluster/v1alpha1"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != azure.Name {
		return nil, fmt.Errorf("non-azure configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != azure.Name {
		return nil, fmt.Errorf("non-azure machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.Azure
	mpool := pool.Platform.Azure

	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}

	// TODO(flaper87): Add support for availability zones
	var machinesets []*clusterapi.MachineSet
	az := ""
	provider, err := provider(clusterID, platform, mpool, osImage, az, role, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
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

	return machinesets, nil
}
