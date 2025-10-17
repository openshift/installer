// Package nutanix generates Machine objects for Nutanix.
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

// MachineSets returns a list of machine sets for a given machine pool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	// Check if the platform is Nutanix.
	if configPlatform := config.Platform.Name(); configPlatform != nutanix.Name {
		return nil, fmt.Errorf("non-nutanix configuration: %q", configPlatform)
	}

	// Check if the machine pool platform is Nutanix.
	if poolPlatform := pool.Platform.Name(); poolPlatform != nutanix.Name {
		return nil, fmt.Errorf("non-nutanix machine-pool: %q", poolPlatform)
	}

	platform := config.Platform.Nutanix
	mpool := pool.Platform.Nutanix
	var total int32

	// Get total number of replicas.
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}

	var machinesets []*machineapi.MachineSet
	numOfFDs := int32(len(mpool.FailureDomains))
	numOfMachineSets := numOfFDs

	// Adjust number of machine sets if necessary.
	if numOfMachineSets == 0 {
		numOfMachineSets = 1
	} else if pool.Replicas != nil && numOfMachineSets > total {
		numOfMachineSets = total
	}

	fdName2ReplicasMap := make(map[string]*int32, numOfMachineSets)
	fdName2FDsMap := make(map[string]*nutanix.FailureDomain, numOfFDs)

	// If no failure domains, assign replicas to a single machine set.
	if numOfFDs == 0 {
		if pool.Replicas != nil {
			fdName2ReplicasMap[""] = ptr.To(int32(*pool.Replicas))
		} else {
			fdName2ReplicasMap[""] = nil
		}
	} else {
		// Distribute replicas to failure domains based on order.
		for _, fdName := range mpool.FailureDomains {
			fd, err := platform.GetFailureDomainByName(fdName)
			if err != nil {
				return nil, err
			}
			fdName2FDsMap[fdName] = fd
		}

		// Distribute replicas evenly across failure domains.
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

	// Create MachineSets based on failure domains and replicas.
	var idx int32
	for fdName, replicaPtr := range fdName2ReplicasMap {
		name := fmt.Sprintf("%s-%s", clusterID, pool.Name)

		var failureDomain *nutanix.FailureDomain
		if fdName != "" {
			failureDomain = fdName2FDsMap[fdName]
			name = fmt.Sprintf("%s-%d", name, idx)
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
						// Versions are managed separately via cluster operators.
					},
				},
			},
		}
		machinesets = append(machinesets, mset)
		idx++
	}

	return machinesets, nil
}
