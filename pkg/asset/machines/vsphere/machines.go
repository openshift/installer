// Package vsphere generates Machine objects for vsphere.
package vsphere

import (
	"fmt"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"

	machineapi "github.com/openshift/api/machine/v1beta1"
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

	var failureDomain vsphere.FailureDomain
	var machines []machineapi.Machine
	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere
	replicas := int64(1)

	numOfZones := len(mpool.Zones)

	zones, err := getDefinedZonesFromTopology(platform)
	if err != nil {
		return machines, err
	}

	if pool.Replicas != nil {
		replicas = *pool.Replicas
	}

	// Create hosts to populate from.  Copying so we can remove without changing original
	// and only put in the ones that match the role.
	var hosts []*vsphere.Host
	if config.Platform.VSphere.Hosts != nil {
		for _, host := range config.Platform.VSphere.Hosts {
			if (host.IsCompute() && role == "worker") || (host.IsControlPlane() && role == "master") {
				logrus.Debugf("Adding host for static ip assignment: %v - %v", host.FailureDomain, host.NetworkDevice.IPAddrs[0])
				hosts = append(hosts, host)
			}
		}
	}

	for idx := int64(0); idx < replicas; idx++ {
		logrus.Debugf("Creating %v machine %v", role, idx)
		var host *vsphere.Host
		desiredZone := mpool.Zones[int(idx)%numOfZones]
		if hosts != nil && int(idx) < len(hosts) {
			host = hosts[idx]
			if host.FailureDomain != "" {
				desiredZone = host.FailureDomain
			}
		}
		logrus.Debugf("Desired zone: %v", desiredZone)

		if _, exists := zones[desiredZone]; !exists {
			return nil, errors.Errorf("zone [%s] specified by machinepool is not defined", desiredZone)
		}

		failureDomain = zones[desiredZone]

		machineLabels := map[string]string{
			"machine.openshift.io/cluster-api-cluster":      clusterID,
			"machine.openshift.io/cluster-api-machine-role": role,
			"machine.openshift.io/cluster-api-machine-type": role,
		}

		osImageForZone := failureDomain.Topology.Template
		if failureDomain.Topology.Template == "" {
			osImageForZone = fmt.Sprintf("%s-%s-%s", osImage, failureDomain.Region, failureDomain.Zone)
		}

		vcenter, err := getVCenterFromServerName(failureDomain.Server, platform)
		if err != nil {
			return nil, errors.Wrap(err, "unable to find vCenter in failure domains")
		}
		provider, err := provider(clusterID, vcenter, failureDomain, mpool, osImageForZone, userDataSecret)
		if err != nil {
			return nil, errors.Wrap(err, "failed to create provider")
		}

		// Apply static IP if configured
		applyNetworkConfig(host, provider)

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

// applyNetworkConfig this function will apply the static ip configuration to the networkDevice
// field in the provider spec.  The function will use the desired zone to determine which config
// to apply and then remove that host config from the hosts array.
func applyNetworkConfig(host *vsphere.Host, provider *machineapi.VSphereMachineProviderSpec) {
	if host != nil {
		networkDevice := host.NetworkDevice
		if networkDevice != nil {
			provider.Network.Devices[0].IPAddrs = networkDevice.IPAddrs
			provider.Network.Devices[0].Nameservers = networkDevice.Nameservers
			provider.Network.Devices[0].Gateway = networkDevice.Gateway
		}
	}
}

func provider(clusterID string, vcenter *vsphere.VCenter, failureDomain vsphere.FailureDomain, mpool *vsphere.MachinePool, osImage string, userDataSecret string) (*machineapi.VSphereMachineProviderSpec, error) {
	networkDeviceSpec := make([]machineapi.NetworkDeviceSpec, len(failureDomain.Topology.Networks))

	// If failureDomain.Topology.Folder is empty this will be used
	folder := fmt.Sprintf("/%s/vm/%s", failureDomain.Topology.Datacenter, clusterID)

	// If failureDomain.Topology.ResourcePool is empty this will be used
	// computeCluster is required to be a path
	resourcePool := fmt.Sprintf("%s/Resources", failureDomain.Topology.ComputeCluster)

	if failureDomain.Topology.Folder != "" {
		folder = failureDomain.Topology.Folder
	}
	if failureDomain.Topology.ResourcePool != "" {
		resourcePool = failureDomain.Topology.ResourcePool
	}

	for i, network := range failureDomain.Topology.Networks {
		networkDeviceSpec[i] = machineapi.NetworkDeviceSpec{NetworkName: network}
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
			Devices: networkDeviceSpec,
		},
		Workspace: &machineapi.Workspace{
			Server:       vcenter.Server,
			Datacenter:   failureDomain.Topology.Datacenter,
			Datastore:    failureDomain.Topology.Datastore,
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
