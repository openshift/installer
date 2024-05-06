// Package vsphere generates Machine objects for vsphere.
package vsphere

import (
	"fmt"
	"path"
	"strconv"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	ipamv1 "sigs.k8s.io/cluster-api/exp/ipam/api/v1beta1"

	v1 "github.com/openshift/api/config/v1"
	machinev1 "github.com/openshift/api/machine/v1"
	machineapi "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// MachineData contains all result output from the Machines() function.
type MachineData struct {
	Machines               []machineapi.Machine
	ControlPlaneMachineSet *machinev1.ControlPlaneMachineSet
	IPClaims               []ipamv1.IPAddressClaim
	IPAddresses            []ipamv1.IPAddress
}

// Machines returns a list of machines for a machinepool.
func Machines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role, userDataSecret string) (*MachineData, error) {
	data := &MachineData{}
	if configPlatform := config.Platform.Name(); configPlatform != vsphere.Name {
		return data, fmt.Errorf("non vsphere configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != vsphere.Name {
		return data, fmt.Errorf("non-VSphere machine-pool: %q", poolPlatform)
	}

	var failureDomain vsphere.FailureDomain
	platform := config.Platform.VSphere
	mpool := pool.Platform.VSphere
	replicas := int32(1)

	numOfZones := len(mpool.Zones)

	zones, err := getDefinedZonesFromTopology(platform)
	if err != nil {
		return data, err
	}

	if pool.Replicas != nil {
		replicas = int32(*pool.Replicas)
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

	failureDomains := []machinev1.VSphereFailureDomain{}

	vsphereMachineProvider := &machineapi.VSphereMachineProviderSpec{}

	for idx := int32(0); idx < replicas; idx++ {
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
			return data, errors.Errorf("zone [%s] specified by machinepool is not defined", desiredZone)
		}

		failureDomain = zones[desiredZone]

		machineLabels := map[string]string{
			"machine.openshift.io/cluster-api-cluster":      clusterID,
			"machine.openshift.io/cluster-api-machine-role": role,
			"machine.openshift.io/cluster-api-machine-type": role,
		}

		if !hasFailureDomain(failureDomains, failureDomain.Name) {
			failureDomains = append(failureDomains, machinev1.VSphereFailureDomain{
				Name: failureDomain.Name,
			})
		}

		osImageForZone := failureDomain.Topology.Template
		if failureDomain.Topology.Template == "" {
			osImageForZone = fmt.Sprintf("%s-%s-%s", osImage, failureDomain.Region, failureDomain.Zone)
		}

		vcenter, err := getVCenterFromServerName(failureDomain.Server, platform)
		if err != nil {
			return data, errors.Wrap(err, "unable to find vCenter in failure domains")
		}
		provider, err := provider(clusterID, vcenter, failureDomain, mpool, osImageForZone, userDataSecret)
		if err != nil {
			return data, errors.Wrap(err, "failed to create provider")
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

		// Apply static IP if configured
		claim, address, err := applyNetworkConfig(host, provider, machine)
		if err != nil {
			return data, err
		} else if claim != nil && address != nil {
			data.IPClaims = append(data.IPClaims, claim...)
			data.IPAddresses = append(data.IPAddresses, address...)
		}
		data.Machines = append(data.Machines, machine)

		vsphereMachineProvider = provider.DeepCopy()
	}

	// when multiple zones are defined, network and workspace are derived from the topology
	origProv := vsphereMachineProvider.DeepCopy()
	if len(failureDomains) >= 1 {
		vsphereMachineProvider.Network = machineapi.NetworkSpec{}
		vsphereMachineProvider.Workspace = &machineapi.Workspace{}
		vsphereMachineProvider.Template = ""
	}

	// Only set AddressesFromPools and Nameservers if AddressesFromPools is > 0, else revert to
	// the older static IP manifest way.
	if len(hosts) > 0 {
		if len(origProv.Network.Devices[0].AddressesFromPools) > 0 {
			vsphereMachineProvider.Network.Devices = []machineapi.NetworkDeviceSpec{
				{
					AddressesFromPools: origProv.Network.Devices[0].AddressesFromPools,
					Nameservers:        origProv.Network.Devices[0].Nameservers,
				},
			}
		} else {
			// Older static IP config, lets remove network since it'll come from FD
			vsphereMachineProvider.Network = machineapi.NetworkSpec{}
		}
	}

	data.ControlPlaneMachineSet = &machinev1.ControlPlaneMachineSet{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "machine.openshift.io/v1",
			Kind:       "ControlPlaneMachineSet",
		},
		ObjectMeta: metav1.ObjectMeta{
			Namespace: "openshift-machine-api",
			Name:      "cluster",
			Labels: map[string]string{
				"machine.openshift.io/cluster-api-cluster": clusterID,
			},
		},
		Spec: machinev1.ControlPlaneMachineSetSpec{
			Replicas: &replicas,
			State:    machinev1.ControlPlaneMachineSetStateActive,
			Selector: metav1.LabelSelector{
				MatchLabels: map[string]string{
					"machine.openshift.io/cluster-api-machine-role": role,
					"machine.openshift.io/cluster-api-machine-type": role,
					"machine.openshift.io/cluster-api-cluster":      clusterID,
				},
			},
			Template: machinev1.ControlPlaneMachineSetTemplate{
				MachineType: machinev1.OpenShiftMachineV1Beta1MachineType,
				OpenShiftMachineV1Beta1Machine: &machinev1.OpenShiftMachineV1Beta1MachineTemplate{
					FailureDomains: &machinev1.FailureDomains{
						Platform: v1.VSpherePlatformType,
						VSphere:  failureDomains,
					},
					ObjectMeta: machinev1.ControlPlaneMachineSetTemplateObjectMeta{
						Labels: map[string]string{
							"machine.openshift.io/cluster-api-cluster":      clusterID,
							"machine.openshift.io/cluster-api-machine-role": role,
							"machine.openshift.io/cluster-api-machine-type": role,
						},
					},
					Spec: machineapi.MachineSpec{
						ProviderSpec: machineapi.ProviderSpec{
							Value: &runtime.RawExtension{Object: vsphereMachineProvider},
						},
					},
				},
			},
		},
	}

	return data, nil
}

// applyNetworkConfig this function will apply the static ip configuration to the networkDevice
// field in the provider spec.  The function will use the desired zone to determine which config
// to apply and then remove that host config from the hosts array.
func applyNetworkConfig(host *vsphere.Host, provider *machineapi.VSphereMachineProviderSpec, machine machineapi.Machine) ([]ipamv1.IPAddressClaim, []ipamv1.IPAddress, error) {
	var ipClaims []ipamv1.IPAddressClaim
	var ipAddrs []ipamv1.IPAddress
	if host != nil {
		networkDevice := host.NetworkDevice
		if networkDevice != nil {
			for idx, address := range networkDevice.IPAddrs {
				provider.Network.Devices[0].Nameservers = networkDevice.Nameservers
				provider.Network.Devices[0].AddressesFromPools = append(provider.Network.Devices[0].AddressesFromPools, machineapi.AddressesFromPool{
					Group:    "installer.openshift.io",
					Name:     fmt.Sprintf("default-%d", idx),
					Resource: "IPPool",
				},
				)

				// Generate the capi networking objects
				slashIndex := strings.Index(address, "/")
				ipAddress := address[0:slashIndex]
				prefix, err := strconv.Atoi(address[slashIndex+1:])
				if err != nil {
					return nil, nil, errors.Wrap(err, "unable to determine address prefix")
				}
				ipClaim, ipAddr := generateCapiNetwork(machine.Name, ipAddress, networkDevice.Gateway, prefix, 0, idx)
				ipClaims = append(ipClaims, *ipClaim)
				ipAddrs = append(ipAddrs, *ipAddr)
			}
		}
	}

	return ipClaims, ipAddrs, nil
}

// generateCapiNetwork this function will create IPAddressClaim and IPAddress for the specified information.
func generateCapiNetwork(machineName, ipAddress, gateway string, prefix, deviceIndex, ipIndex int) (*ipamv1.IPAddressClaim, *ipamv1.IPAddress) {
	// Generate PoolRef
	apigroup := "installer.openshift.io"
	poolRef := corev1.TypedLocalObjectReference{
		APIGroup: &apigroup,
		Kind:     "IPPool",
		Name:     fmt.Sprintf("default-%d", ipIndex),
	}

	// Generate IPAddressClaim
	ipclaim := &ipamv1.IPAddressClaim{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ipam.cluster.x-k8s.io/v1beta1",
			Kind:       "IPAddressClaim",
		},
		ObjectMeta: metav1.ObjectMeta{
			Finalizers: []string{
				machineapi.IPClaimProtectionFinalizer,
			},
			Name:      fmt.Sprintf("%s-claim-%d-%d", machineName, deviceIndex, ipIndex),
			Namespace: "openshift-machine-api",
		},
		Spec: ipamv1.IPAddressClaimSpec{
			PoolRef: poolRef,
		},
	}

	// Populate IPAddress info
	ipaddr := &ipamv1.IPAddress{
		TypeMeta: metav1.TypeMeta{
			APIVersion: "ipam.cluster.x-k8s.io/v1beta1",
			Kind:       "IPAddress",
		},
		ObjectMeta: metav1.ObjectMeta{
			Name:      fmt.Sprintf("%s-claim-%d-%d", machineName, deviceIndex, ipIndex),
			Namespace: "openshift-machine-api",
		},
		Spec: ipamv1.IPAddressSpec{
			Address: ipAddress,
			ClaimRef: corev1.LocalObjectReference{
				Name: ipclaim.Name,
			},
			Gateway: gateway,
			PoolRef: poolRef,
			Prefix:  prefix,
		},
	}

	ipclaim.Status = ipamv1.IPAddressClaimStatus{
		AddressRef: corev1.LocalObjectReference{
			Name: ipaddr.Name,
		},
	}

	return ipclaim, ipaddr
}

func provider(clusterID string, vcenter *vsphere.VCenter, failureDomain vsphere.FailureDomain, mpool *vsphere.MachinePool, osImage string, userDataSecret string) (*machineapi.VSphereMachineProviderSpec, error) {
	networkDeviceSpec := make([]machineapi.NetworkDeviceSpec, len(failureDomain.Topology.Networks))

	// If failureDomain.Topology.Folder is empty this will be used
	folder := path.Clean(fmt.Sprintf("/%s/vm/%s", failureDomain.Topology.Datacenter, clusterID))

	// If failureDomain.Topology.ResourcePool is empty this will be used
	// computeCluster is required to be a path
	resourcePool := path.Clean(fmt.Sprintf("%s/Resources", failureDomain.Topology.ComputeCluster))

	if failureDomain.Topology.Folder != "" {
		folder = failureDomain.Topology.Folder
	}
	if failureDomain.Topology.ResourcePool != "" {
		resourcePool = failureDomain.Topology.ResourcePool
	}

	resourcePool = path.Clean(resourcePool)

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
		TagIDs:            failureDomain.Topology.TagIDs,
		NumCPUs:           mpool.NumCPUs,
		NumCoresPerSocket: mpool.NumCoresPerSocket,
		MemoryMiB:         mpool.MemoryMiB,
		DiskGiB:           mpool.OSDisk.DiskSizeGB,
	}, nil
}

// ConfigMasters sets the PublicIP flag and assigns a set of load balancers to the given machines
func ConfigMasters(machines []machineapi.Machine, clusterID string) {
}

func hasFailureDomain(failureDomains []machinev1.VSphereFailureDomain, failureDomain string) bool {
	for _, fd := range failureDomains {
		if fd.Name == failureDomain {
			return true
		}
	}
	return false
}
