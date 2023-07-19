// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	"github.com/gophercloud/utils/openstack/clientconfig"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	clusterapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

const maxInt32 int64 = int64(^uint32(0)) >> 1

// MachineSets returns the MachineSets encoded by the given machine-pool. The
// number of returned MachineSets, while being capped to the number of
// replicas, depends on the variable-length fields in the machine-pool. Each
// MachineSet generates a set of identical Machines; to encode for Machines
// spread on, say, three availability zones, three MachineSets must be
// produced. Note that for each variable-length field (currently: Compute
// availability zones, Storage availability zones and Root volume types), when
// more than one is specified, values of identical index are grouped in the
// same MachineSet.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string, clientOpts *clientconfig.ClientOpts) ([]*clusterapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}
	if pool.Replicas == nil || *pool.Replicas < 1 {
		return nil, nil
	}
	platform := config.Platform.OpenStack
	mpool := pool.Platform.OpenStack
	trunkSupport, err := checkNetworkExtensionAvailability(platform.Cloud, "trunk", clientOpts)
	if err != nil {
		return nil, err
	}

	failureDomains := failureDomainsFromSpec(*mpool)
	numberOfFailureDomains := int64(len(failureDomains))

	machinesets := make([]*clusterapi.MachineSet, len(failureDomains))
	for idx := range machinesets {
		var replicaNumber int32
		{
			replicas := *pool.Replicas / numberOfFailureDomains
			if int64(idx) < *pool.Replicas%numberOfFailureDomains {
				replicas++
			}
			if replicas > maxInt32 {
				return nil, fmt.Errorf("the number of requested worker replicas (%d) is too high. Each MachineSet can hold %d replicas; the install-config encodes for %d MachineSets, which gives us a replica number of %d", *pool.Replicas, maxInt32, numberOfFailureDomains, replicas)
			}
			replicaNumber = int32(replicas)
		}

		providerSpec := generateProviderSpec(
			clusterID,
			platform,
			mpool,
			osImage,
			role,
			userDataSecret,
			trunkSupport,
			failureDomains[idx],
		)

		// Set unique name for the machineset
		name := fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx)

		machinesets[idx] = &clusterapi.MachineSet{
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
				Replicas: &replicaNumber,
				Selector: metav1.LabelSelector{
					MatchLabels: map[string]string{
						"machine.openshift.io/cluster-api-machineset": name,
						"machine.openshift.io/cluster-api-cluster":    clusterID,
					},
				},
				Template: clusterapi.MachineTemplateSpec{
					ObjectMeta: clusterapi.ObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-machineset":   name,
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: clusterapi.MachineSpec{
						ProviderSpec: clusterapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: providerSpec},
						},
						// we don't need to set Versions, because we control those via cluster operators.
					},
				},
			},
		}
	}

	return machinesets, nil
}
