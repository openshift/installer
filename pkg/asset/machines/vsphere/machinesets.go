//Package vsphere generates Machine objects for vsphere.package vsphere
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

func getDefinedZones(platformSpec *vsphere.Platform) (map[string]*vsphere.Platform, error) {
	zones := make(map[string]*vsphere.Platform)
	if len(platformSpec.VCenters) > 0 {
		for _, vcenter := range platformSpec.VCenters {
			var vcPlatform = vsphere.Platform{
				VCenter:  vcenter.Server,
				Username: vcenter.User,
				Password: vcenter.Password,
			}
			for _, region := range vcenter.Regions {
				regionPlatform := vcPlatform
				regionPlatform.Datacenter = region.Datacenter
				if len(region.Zones) == 0 {
					return zones, errors.Errorf("region[%s] has no defined zones", region.Name)
				}
				for _, zone := range region.Zones {
					if _, exists := zones[zone.Name]; exists {
						return zones, errors.Errorf("zones with duplicate name[%s] defined", zone.Name)
					}
					zonePlatform := regionPlatform
					zonePlatform.Cluster = zone.Cluster
					zonePlatform.DefaultDatastore = zone.Datastore
					zonePlatform.Network = zone.Network
					zones[zone.Name] = &zonePlatform
				}
			}
		}
	}
	return zones, nil
}

// MachineSets returns a list of machinesets for a machinepool.
func MachineSets(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]*machineapi.MachineSet, error) {
	if configPlatform := config.Platform.Name(); configPlatform != vsphere.Name {
		return nil, fmt.Errorf("non vsphere configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != vsphere.Name {
		return nil, fmt.Errorf("non-VSphere machine-pool: %q", poolPlatform)
	}
	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere
	azs := mpool.Zones
	total := int32(0)
	if pool.Replicas != nil {
		total = int32(*pool.Replicas)
	}
	numOfAZs := int32(len(azs))
	var machinesets []*machineapi.MachineSet
	if numOfAZs > 0 {
		zones, err := getDefinedZones(platform)
		if err != nil {
			return machinesets, err
		}
		for idx := range azs {
			replicas := int32(total / numOfAZs)
			if int32(idx) < total%numOfAZs {
				replicas++
			}
			desiredZone := azs[idx]
			if _, exists := zones[desiredZone]; !exists {
				return nil, errors.Errorf("zone [%s] specified by machinepool is not defined", desiredZone)
			}
			name := fmt.Sprintf("%s-%s-%s", clusterID, pool.Name, desiredZone)
			machineset, err := getMachineSetWithPlatform(
				clusterID,
				name,
				mpool,
				osImage,
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
			total,
			role,
			userDataSecret)
		if err != nil {
			return machinesets, err
		}
		machinesets = append(machinesets, machineset)
	}
	return machinesets, nil
}
