// Package nutanix generates Machine objects for nutanix.package nutanix
package nutanix

import (
	"fmt"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/nutanix"
)

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != nutanix.Name {
		return nil, fmt.Errorf("non nutanix configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != nutanix.Name {
		return nil, fmt.Errorf("non-nutanix machine-pool: %q", poolPlatform)
	}

	platform := config.Platform.Nutanix
	mpool := pool.Platform.Nutanix
	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}

	var machinesets []*machineapi.MachineSet
	numOfFDs := int32(len(mpool.FailureDomains))
	numOfMachineSets := numOfFDs
	if numOfMachineSets == 0 {
		numOfMachineSets = 1
	} else if pool.Replicas != nil && numOfMachineSets > total {
		numOfMachineSets = total
	}

	fdName2ReplicasMap := make(map[string]*int32, numOfMachineSets)
	fdName2FDsMap := make(map[string]*nutanix.FailureDomain, numOfFDs)

	if numOfFDs == 0 {
		if pool.Replicas != nil {
			fdName2ReplicasMap[""] = ptr.To(int32(*pool.Replicas))
		} else {
			fdName2ReplicasMap[""] = nil
		}
	} else {
		// When failure domains is configured for the workers, evenly distribute
		// the machineset replicas to the failure domains, based on order.
		for _, fdName := range mpool.FailureDomains {
			fd, err := platform.GetFailureDomainByName(fdName)
			if err != nil {
				return nil, err
			}
			fdName2FDsMap[fdName] = fd
		}

		if pool.Replicas != nil {
			for i := int32(0); i < total; i++ {
				idx := i % numOfFDs
				fdName := mpool.FailureDomains[idx]
				replica := int32(1)
				if ra, ok := fdName2ReplicasMap[fdName]; ok && ra != nil {
					replica = *ra + 1
				}
				fdName2ReplicasMap[fdName] = ptr.To(replica)
			}
		} else {
			for _, fdName := range mpool.FailureDomains {
				fdName2ReplicasMap[fdName] = nil
			}
		}
	}

	var idx int32
	for fdName, replicaPtr := range fdName2ReplicasMap {
		name := fmt.Sprintf("%s-%s", clusterID, pool.Name)

		var failureDomain *nutanix.FailureDomain
		if fdName != "" {
			failureDomain = fdName2FDsMap[fdName]
			name = fmt.Sprintf("%s-%v", name, idx)
		}

		provider, err := provider(clusterID, platform, mpool, osImage, userDataSecret, failureDomain)
		if err != nil {
			return nil, fmt.Errorf("failed to create provider: %w", err)
		}

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
				Replicas: replicaPtr,
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
		machinesets = append(machinesets, mset)
		idx++
	}

	return machinesets, nil
}
