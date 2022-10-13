// Package vsphere generates Machine objects for vsphere.package vsphere
package vsphere

import (
	"fmt"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/pkg/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

func getMachineSetWithPlatform(
	clusterID string,
	name string,
	mpool *vsphere.MachinePool,
	osImage string,
	platform *vsphere.Platform,
	replicas int32,
	role,
	userDataSecret string) (*machineapi.MachineSet, error) {
	provider, err := provider(clusterID, platform, mpool, osImage, userDataSecret)
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

func getFailureDomain(domainName string, platformSpec *vsphere.Platform) (*vsphere.FailureDomain, error) {
	for _, failureDomain := range platformSpec.FailureDomains {
		if failureDomain.Name == domainName {
			return &failureDomain, nil
		}
	}
	return nil, errors.Errorf("%s is not a defined failure domain", domainName)
}

func getfailureDomain(failureDomainName string, platformSpec *vsphere.Platform) (*vsphere.FailureDomain, error) {
	for _, failureDomain := range platformSpec.FailureDomains {
		if failureDomain.Name == failureDomainName {
			return &failureDomain, nil
		}
	}
	return nil, errors.Errorf("%s is not a defined deployment zone", failureDomainName)
}

// getDefinedZones retrieves zones and associated platform specs that are appropriate to the machine role
func getDefinedZones(platformSpec *vsphere.Platform) (map[string]*vsphere.Platform, error) {
	zones := make(map[string]*vsphere.Platform)

	for _, failureDomain := range platformSpec.FailureDomains {
		vCenter, err := getVCenterFromServerName(failureDomain.Server, platformSpec)
		if err != nil {
			return nil, err
		}
		failureDomain, err := getFailureDomain(failureDomain.Name, platformSpec)
		if err != nil {
			return nil, err
		}
		var vcPlatform = vsphere.Platform{
			VCenter:          vCenter.Server,
			Username:         vCenter.Username,
			Password:         vCenter.Password,
			Datacenter:       failureDomain.Topology.Datacenter,
			DefaultDatastore: failureDomain.Topology.Datastore,
			Folder:           failureDomain.Topology.Folder,
			Cluster:          failureDomain.Topology.ComputeCluster,
			ResourcePool:     failureDomain.Topology.ResourcePool,
			APIVIPs:          platformSpec.APIVIPs,
			IngressVIPs:      platformSpec.IngressVIPs,
			Network:          failureDomain.Topology.Networks[0],
			DiskType:         platformSpec.DiskType,
			FailureDomains:   platformSpec.FailureDomains,
		}
		zones[failureDomain.Name] = &vcPlatform
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
	azs := mpool.Zones
	total := 0
	if pool.Replicas != nil {
		total = int(*pool.Replicas)
	}
	numOfAZs := len(azs)
	var machinesets []*machineapi.MachineSet
	if numOfAZs > 0 {
		zones, err := getDefinedZones(platform)
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

			failureDomain, err := getFailureDomain(desiredZone, platform)
			if err != nil {
				return nil, err
			}

			osImageForZone := fmt.Sprintf("%s-%s-%s", osImage, failureDomain.Region, failureDomain.Zone)
			machineset, err := getMachineSetWithPlatform(
				clusterID,
				name,
				mpool,
				osImageForZone,
				zones[desiredZone],
				replicas,
				role,
				userDataSecret)
			if err != nil {
				return machinesets, err
			}
			machinesets = append(machinesets, machineset)
		}
	} else {
		name := fmt.Sprintf("%s-%s", clusterID, pool.Name)
		machineset, err := getMachineSetWithPlatform(
			clusterID,
			name,
			mpool,
			osImage,
			platform,
			int32(total),
			role,
			userDataSecret)
		if err != nil {
			return machinesets, err
		}
		machinesets = append(machinesets, machineset)
	}
	return machinesets, nil
}
