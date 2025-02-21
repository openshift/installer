package vsphere

import (
	"context"
	"encoding/json"
	"errors"
	"fmt"
	"strings"

	"github.com/sirupsen/logrus"
	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	capv "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	machinev1 "github.com/openshift/api/machine/v1beta1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig/vsphere"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/utils"
)

const (
	masterRole = "master"
)

// ProviderSpecFromRawExtension unmarshals the JSON-encoded spec.
func ProviderSpecFromRawExtension(rawExtension *runtime.RawExtension) (*machinev1.VSphereMachineProviderSpec, error) {
	if rawExtension == nil {
		return &machinev1.VSphereMachineProviderSpec{}, nil
	}

	spec := new(machinev1.VSphereMachineProviderSpec)
	if err := json.Unmarshal(rawExtension.Raw, &spec); err != nil {
		return nil, fmt.Errorf("error unmarshalling providerSpec: %w", err)
	}

	return spec, nil
}

func getNetworkInventoryPath(vcenterContext vsphere.VCenterContext, networkName string, providerSpec *machinev1.VSphereMachineProviderSpec) (string, error) {
	// if networkName is a path, we'll assume that a full path was provided by the admin
	if strings.Contains(networkName, "/") {
		return networkName, nil
	}

	// else, we'll dereference the network name to a full path using the resource pool
	for _, clusterNetworkMap := range vcenterContext.ClusterNetworkMap {
		if _, networkInContext := clusterNetworkMap.NetworkNames[networkName]; !networkInContext {
			continue
		}

		for _, resourcePool := range clusterNetworkMap.ResourcePools {
			if resourcePool.InventoryPath == providerSpec.Workspace.ResourcePool {
				return clusterNetworkMap.NetworkNames[networkName], nil
			}
		}

		// This is a case for UPI (terraform or powercli) the resource pool will not exist
		// prior to running openshift-install create manifests.
		// This also will keep backward compatibility as this was not required to CAPI implementation.
		if strings.Contains(providerSpec.Workspace.ResourcePool, clusterNetworkMap.Cluster) {
			logrus.Debugf("using cluster %s as selector for network device path %s", clusterNetworkMap.Cluster, networkName)
			return clusterNetworkMap.NetworkNames[networkName], nil
		}
	}
	return "", fmt.Errorf("unable to find network %s in resource pool %s", networkName, providerSpec.Workspace.ResourcePool)
}

// GenerateMachines returns a list of capi machines.
func GenerateMachines(ctx context.Context, clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage string, role string, metadata *vsphere.Metadata) ([]*asset.RuntimeFile, error) {
	data, err := Machines(clusterID, config, pool, osImage, role, "")
	if err != nil {
		return nil, fmt.Errorf("unable to retrieve machines: %w", err)
	}
	machines := data.Machines

	capvMachines := make([]*capv.VSphereMachine, 0, len(machines))
	result := make([]*asset.RuntimeFile, 0, len(machines))
	staticIP := false

	for mIndex, machine := range machines {
		providerSpec, ok := machine.Spec.ProviderSpec.Value.Object.(*machinev1.VSphereMachineProviderSpec)
		if !ok {
			return nil, errors.New("unable to convert ProviderSpec to VSphereMachineProviderSpec")
		}

		vcenterContext := metadata.VCenterContexts[providerSpec.Workspace.Server]
		resourcePool := providerSpec.Workspace.ResourcePool

		customVMXKeys := map[string]string{
			"guestinfo.hostname": machine.Name,
			"guestinfo.domain":   strings.TrimSuffix(config.ClusterDomain(), "."),
			"stealclock.enable":  "TRUE",
		}

		capvNetworkDevices := []capv.NetworkDeviceSpec{}
		for _, networkDevice := range providerSpec.Network.Devices {
			networkName, err := getNetworkInventoryPath(vcenterContext, networkDevice.NetworkName, providerSpec)
			if err != nil {
				return nil, fmt.Errorf("unable to get network inventory path: %w", err)
			}
			deviceSpec := capv.NetworkDeviceSpec{
				NetworkName: networkName,
				DHCP4:       true,
			}

			// Static IP configured.  Add kargs.
			if len(networkDevice.AddressesFromPools) > 0 {
				staticIP = true
				kargs, err := utils.ConstructNetworkKargsFromMachine(data.IPClaims, data.IPAddresses, &machines[mIndex], networkDevice)
				if err != nil {
					return nil, fmt.Errorf("unable to get static ip config for machine %v: %w", machine.Name, err)
				}
				customVMXKeys["guestinfo.afterburn.initrd.network-kargs"] = kargs
			}
			capvNetworkDevices = append(capvNetworkDevices, deviceSpec)
		}

		vsphereMachine := &capv.VSphereMachine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiutils.Namespace,
				Name:      machine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},

			Spec: capv.VSphereMachineSpec{
				VirtualMachineCloneSpec: capv.VirtualMachineCloneSpec{
					CloneMode:     capv.FullClone,
					CustomVMXKeys: customVMXKeys,
					Network: capv.NetworkSpec{
						Devices: capvNetworkDevices,
					},
					Folder:            providerSpec.Workspace.Folder,
					Template:          providerSpec.Template,
					Datacenter:        providerSpec.Workspace.Datacenter,
					Server:            providerSpec.Workspace.Server,
					NumCPUs:           providerSpec.NumCPUs,
					NumCoresPerSocket: providerSpec.NumCoresPerSocket,
					MemoryMiB:         providerSpec.MemoryMiB,
					DiskGiB:           providerSpec.DiskGiB,
					Datastore:         providerSpec.Workspace.Datastore,
					ResourcePool:      resourcePool,
				},
			},
		}

		// only set failureDomainName if VMGroup is defined as vm-host group
		// is the only scenario we create vspherefailuredomainspec and vspheredeploymentzone
		if providerSpec.Workspace.VMGroup != "" {
			if failureDomainName, ok := data.MachineFailureDomain[machine.Name]; ok {
				vsphereMachine.Spec.FailureDomain = &failureDomainName
			} else {
				return nil, fmt.Errorf("unable to find failure domain for machine %s", machine.Name)
			}
		}

		// If we have additional disks to add to VM, lets iterate through them and add to CAPV machine
		if len(providerSpec.DataDisks) > 0 {
			dataDisks := []capv.VSphereDisk{}
			for _, disk := range providerSpec.DataDisks {
				newDisk := capv.VSphereDisk{
					Name:    disk.Name,
					SizeGiB: disk.SizeGiB,
				}
				dataDisks = append(dataDisks, newDisk)
			}
			vsphereMachine.Spec.DataDisks = dataDisks
		}

		vsphereMachine.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereMachine"))
		capvMachines = append(capvMachines, vsphereMachine)

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", vsphereMachine.Name)},
			Object: vsphereMachine,
		})

		// Need to determine the infrastructure ref since there may be multi vcenters.
		clusterName := clusterID
		for index, vcenter := range config.Platform.VSphere.VCenters {
			if vcenter.Server == providerSpec.Workspace.Server {
				clusterName = fmt.Sprintf("%v-%d", clusterID, index)
				break
			}
		}

		// Create capi machine for vspheremachine
		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Namespace: capiutils.Namespace,
				Name:      vsphereMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterName,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capv.GroupVersion.String(),
					Kind:       "VSphereMachine",
					Name:       vsphereMachine.Name,
				},
			},
		}
		machine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
			Object: machine,
		})
	}

	// as part of provisioning control plane nodes, we need to create a bootstrap node as well
	if role == masterRole {
		customVMXKeys := map[string]string{}

		// If we detected static IP for masters, lets apply to bootstrap as well.
		if staticIP {
			kargs, err := utils.ConstructKargsForBootstrap(config)
			if err != nil {
				return nil, fmt.Errorf("unable to get static ip config for bootstrap: %w", err)
			}
			customVMXKeys["guestinfo.afterburn.initrd.network-kargs"] = kargs
		}

		bootstrapSpec := capvMachines[0].Spec
		bootstrapSpec.CustomVMXKeys = customVMXKeys
		bootstrapVSphereMachine := &capv.VSphereMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name: fmt.Sprintf("%s-bootstrap", clusterID),
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: bootstrapSpec,
		}
		bootstrapVSphereMachine.SetGroupVersionKind(capv.GroupVersion.WithKind("VSphereMachine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", bootstrapVSphereMachine.Name)},
			Object: bootstrapVSphereMachine,
		})

		// Need to determine the infrastructure ref since there may be multi vcenters.
		clusterName := clusterID
		for index, vcenter := range config.Platform.VSphere.VCenters {
			if vcenter.Server == bootstrapSpec.Server {
				clusterName = fmt.Sprintf("%v-%d", clusterID, index)
				break
			}
		}

		bootstrapMachine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: bootstrapVSphereMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterName,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-bootstrap", clusterID)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capv.GroupVersion.String(),
					Kind:       "VSphereMachine",
					Name:       bootstrapVSphereMachine.Name,
				},
			},
		}
		bootstrapMachine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))
		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", bootstrapVSphereMachine.Name)},
			Object: bootstrapMachine,
		})
	}
	return result, nil
}
