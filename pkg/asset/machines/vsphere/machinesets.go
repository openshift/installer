// Package vsphere generates Machine objects for vsphere.package vsphere
package vsphere

import (
	"fmt"

	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func getMachineSetWithPlatform(
	clusterID string,
	name string,
	mpool *vsphere.MachinePool,
	osImage string,
	failureDomain vsphere.FailureDomain,
	vcenter *vsphere.VCenter,
	replicas int32,
	role,
	userDataSecret string) (*machineapi.MachineSet, error) {
	provider, err := provider(clusterID, vcenter, failureDomain, mpool, osImage, userDataSecret)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create provider")
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
			Replicas: &replicas,
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
	return mset, nil
}

func getVCenterFromServerName(server string, platformSpec *vsphere.Platform) (*vsphere.VCenter, error) {
	for _, vCenter := range platformSpec.VCenters {
		if vCenter.Server == server {
			return &vCenter, nil
		}
	}
	return nil, errors.Errorf("unable to find vCenter %s", server)
}

func getDefinedZonesFromTopology(p *vsphere.Platform) (map[string]vsphere.FailureDomain, error) {
	zones := make(map[string]vsphere.FailureDomain)
	for _, failureDomain := range p.FailureDomains {
		zones[failureDomain.Name] = failureDomain
	}
	return zones, nil
}

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != vsphere.Name {
		return nil, fmt.Errorf("non vsphere configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != vsphere.Name {
		return nil, fmt.Errorf("non-vsphere machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere
	// The machinepool has no zones defined, there are FailureDomains
	// This is a vSphere zonal installation. Generate machinepool zone
	// list.
	if len(mpool.Zones) == 0 {
		for _, fd := range config.VSphere.FailureDomains {
			mpool.Zones = append(mpool.Zones, fd.Name)
		}
	}
	azs := mpool.Zones
	total := 0
	if pool.Replicas != nil {
		total = int(*pool.Replicas)
	}
	numOfAZs := len(azs)
	machinesets := make([]*machineapi.MachineSet, 0, numOfAZs)

	zones, err := getDefinedZonesFromTopology(platform)

	if err != nil {
		return machinesets, err
	}
	for idx := range azs {
		replicas := int32(total / numOfAZs)
		if idx < total%numOfAZs {
			replicas++
		}
		desiredZone := azs[idx]
		if _, exists := zones[desiredZone]; !exists {
			return nil, errors.Errorf("zone [%s] specified by machinepool is not defined", desiredZone)
		}
		name := fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx)

		failureDomain := zones[desiredZone]

		vcenter, err := getVCenterFromServerName(failureDomain.Server, platform)
		if err != nil {
			return nil, err
		}

		osImageForZone := failureDomain.Topology.Template
		if failureDomain.Topology.Template == "" {
			osImageForZone = fmt.Sprintf("%s-%s-%s", osImage, failureDomain.Region, failureDomain.Zone)
		}
		machineset, err := getMachineSetWithPlatform(
			clusterID,
			name,
			mpool,
			osImageForZone,
			failureDomain,
			vcenter,
			replicas,
			role,
			userDataSecret)
		if err != nil {
			return machinesets, err
		}
		machinesets = append(machinesets, machineset)
	}

	return machinesets, nil
}
