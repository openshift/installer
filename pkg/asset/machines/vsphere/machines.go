// Package vsphere generates Machine objects for vsphere.
package vsphere

import (
	"fmt"
	"regexp"

	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/pkg/errors"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) ([]machineapi.Machine, error) {
	if configPlatform := config.Platform.Name(); configPlatform != vsphere.Name {
		return nil, fmt.Errorf("non vsphere configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != vsphere.Name {
		return nil, fmt.Errorf("non-VSphere machine-pool: %q", poolPlatform)
	}

	var machines []machineapi.Machine

	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere

	azs := mpool.Zones
	numOfAZs := len(azs)
	definedZones := make(map[string]*vsphere.Platform)

	if numOfAZs > 0 {
		zones, err := getDefinedZones(platform)
		if err != nil {
			return machines, err
		}
		definedZones = zones
	}

	replicas := int64(1)
	if pool.Replicas != nil {
		replicas = *pool.Replicas
	}

	for idx := int64(0); idx < replicas; idx++ {
		var failureDomain *vsphere.FailureDomain
		if numOfAZs > 0 {
			desiredZone := mpool.Zones[int(idx)%numOfAZs]
			if _, exists := definedZones[desiredZone]; !exists {
				return nil, errors.Errorf("zone [%s] specified by machinepool is not defined", desiredZone)
			}
			platform = definedZones[desiredZone]
			for idx, knownFailureDomain := range platform.FailureDomains {
				if knownFailureDomain.Name == desiredZone {
					failureDomain = &platform.FailureDomains[idx]
				}
			}
		}

		machineLabels := map[string]string{
			"machine.openshift.io/cluster-api-cluster":      clusterID,
			"machine.openshift.io/cluster-api-machine-role": role,
			"machine.openshift.io/cluster-api-machine-type": role,
		}

		osImageForZone := osImage

		if failureDomain != nil {
			osImageForZone = fmt.Sprintf("%s-%s-%s", osImage, failureDomain.Region, failureDomain.Zone)
		}

		provider, err := provider(clusterID, platform, mpool, osImageForZone, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}

		machine := machineapi.Machine{
			TypeMeta: metav1.TypeMeta{
				APIVersion: "machine.openshift.io/v1beta1",
				Kind:       "Machine",
			},
			ObjectMeta: metav1.ObjectMeta{
				Namespace: "openshift-machine-api",
				Name:      fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx),
				Labels:    machineLabels,
			},
			Spec: machineapi.MachineSpec{
				ProviderSpec: machineapi.ProviderSpec{
					Value: &runtime.RawExtension{Object: provider},
				},
				// we don't need to set Versions, because we control those via operators.
			},
		}
		machines = append(machines, machine)
	}
	return machines, nil
}

func provider(clusterID string, platform *vsphere.Platform, mpool *vsphere.MachinePool, osImage string, userDataSecret string) (*machineapi.VSphereMachineProviderSpec, error) {
	folder := fmt.Sprintf("/%s/vm/%s", platform.Datacenter, clusterID)

	resourcePool := fmt.Sprintf("/%s/host/%s/Resources", platform.Datacenter, platform.Cluster)
	resourcePoolPrefix := "^\\/(.*?)\\/host\\/(.*?)"
	hasFullPath, _ := regexp.MatchString(resourcePoolPrefix, platform.Cluster)
	if hasFullPath {
		resourcePool = fmt.Sprintf("%s/Resources", platform.Cluster)
	}

	if platform.Folder != "" {
		folder = platform.Folder
	}
	if platform.ResourcePool != "" {
		resourcePool = platform.ResourcePool
	}

	return &machineapi.VSphereMachineProviderSpec{
		TypeMeta: metav1.TypeMeta{
			APIVersion: machineapi.SchemeGroupVersion.String(),
			Kind:       "VSphereMachineProviderSpec",
		},
		UserDataSecret:    &corev1.LocalObjectReference{Name: userDataSecret},
		CredentialsSecret: &corev1.LocalObjectReference{Name: "vsphere-cloud-credentials"},
		Template:          osImage,
		Network: machineapi.NetworkSpec{
			Devices: []machineapi.NetworkDeviceSpec{
				{
					NetworkName: platform.Network,
				},
			},
		},
		Workspace: &machineapi.Workspace{
			Server:       platform.VCenter,
			Datacenter:   platform.Datacenter,
			Datastore:    platform.DefaultDatastore,
			Folder:       folder,
			ResourcePool: resourcePool,
		},
		NumCPUs:           mpool.NumCPUs,
		NumCoresPerSocket: mpool.NumCoresPerSocket,
		MemoryMiB:         mpool.MemoryMiB,
		DiskGiB:           mpool.OSDisk.DiskSizeGB,
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
}
