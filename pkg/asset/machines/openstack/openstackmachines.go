// Package openstack generates Machine objects for openstack.
package openstack

import (
	"fmt"

	v1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"
	capo "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha7"
	capi "sigs.k8s.io/cluster-api/api/v1beta1"

	machinev1 "github.com/openshift/api/machine/v1"
	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/manifests/capiutils"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
)

// GenerateMachines returns manifests and runtime objects to provision the control plane (including bootstrap, if applicable) nodes using CAPI.
func GenerateMachines(clusterID string, config *types.InstallConfig, pool *types.MachinePool, osImage, role string, trunkSupport bool) ([]*asset.RuntimeFile, error) {
	if configPlatform := config.Platform.Name(); configPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack configuration: %q", configPlatform)
	}
	if poolPlatform := pool.Platform.Name(); poolPlatform != openstack.Name {
		return nil, fmt.Errorf("non-OpenStack machine-pool: %q", poolPlatform)
	}

	mpool := pool.Platform.OpenStack

	total := int64(1)
	if role == "master" && pool.Replicas != nil {
		total = *pool.Replicas
	}

	var result []*asset.RuntimeFile
	failureDomains := failureDomainsFromSpec(*mpool)
	for idx := int64(0); idx < total; idx++ {
		failureDomain := failureDomains[uint(idx)%uint(len(failureDomains))]
		machineSpec, err := generateMachineSpec(
			clusterID,
			config.Platform.OpenStack,
			mpool,
			osImage,
			role,
			trunkSupport,
			failureDomain,
		)
		if err != nil {
			return nil, err
		}

		machineName := fmt.Sprintf("%s-%s-%d", clusterID, pool.Name, idx)
		machineLabels := map[string]string{
			"cluster.x-k8s.io/control-plane": "",
		}
		if role == "bootstrap" {
			machineName = capiutils.GenerateBoostrapMachineName(clusterID)
			machineLabels = map[string]string{
				"cluster.x-k8s.io/control-plane": "",
				"install.openshift.io/bootstrap": "",
			}
		}
		openStackMachine := &capo.OpenStackMachine{
			ObjectMeta: metav1.ObjectMeta{
				Name:   machineName,
				Labels: machineLabels,
			},
			Spec: *machineSpec,
		}
		openStackMachine.SetGroupVersionKind(capo.GroupVersion.WithKind("OpenStackMachine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_inframachine_%s.yaml", openStackMachine.Name)},
			Object: openStackMachine,
		})

		// The instanceSpec used to create the server uses the failureDomain from CAPI Machine
		// defined bellow. This field must match a Key on FailureDomains stored in the cluster.
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/blob/main/controllers/openstackmachine_controller.go#L472
		// TODO (maysa): test this
		machine := &capi.Machine{
			ObjectMeta: metav1.ObjectMeta{
				Name: openStackMachine.Name,
				Labels: map[string]string{
					"cluster.x-k8s.io/control-plane": "",
				},
			},
			Spec: capi.MachineSpec{
				ClusterName: clusterID,
				Bootstrap: capi.Bootstrap{
					DataSecretName: ptr.To(fmt.Sprintf("%s-%s", clusterID, role)),
				},
				InfrastructureRef: v1.ObjectReference{
					APIVersion: capo.GroupVersion.String(),
					Kind:       "OpenStackMachine",
					Name:       openStackMachine.Name,
				},
				FailureDomain: &failureDomain.AvailabilityZone,
			},
		}
		machine.SetGroupVersionKind(capi.GroupVersion.WithKind("Machine"))

		result = append(result, &asset.RuntimeFile{
			File:   asset.File{Filename: fmt.Sprintf("10_machine_%s.yaml", machine.Name)},
			Object: machine,
		})
	}
	return result, nil
}

func generateMachineSpec(clusterID string, platform *openstack.Platform, mpool *openstack.MachinePool, osImage string, role string, trunkSupport bool, failureDomain machinev1.OpenStackFailureDomain) (*capo.OpenStackMachineSpec, error) {
	port := capo.PortOpts{}

	addressPairs := []capo.AddressPair{}
	for _, apiVIP := range platform.APIVIPs {
		addressPairs = append(addressPairs, capo.AddressPair{IPAddress: apiVIP})
	}
	for _, ingressVIP := range platform.IngressVIPs {
		addressPairs = append(addressPairs, capo.AddressPair{IPAddress: ingressVIP})
	}

	if platform.ControlPlanePort != nil {
		port.Network = &capo.NetworkFilter{
			Name: platform.ControlPlanePort.Network.Name,
			ID:   platform.ControlPlanePort.Network.ID,
		}

		var fixedIPs []capo.FixedIP
		for _, fixedIP := range platform.ControlPlanePort.FixedIPs {
			fixedIPs = append(fixedIPs, capo.FixedIP{
				Subnet: &capo.SubnetFilter{
					ID:   fixedIP.Subnet.ID,
					Name: fixedIP.Subnet.Name,
				}})
		}
		port.FixedIPs = fixedIPs
		if len(addressPairs) > 0 {
			port.AllowedAddressPairs = addressPairs
		}
	} else {
		port = capo.PortOpts{
			FixedIPs: []capo.FixedIP{
				{
					Subnet: &capo.SubnetFilter{
						// NOTE(mandre) the format of the subnet name changes when letting CAPI create it.
						// So solely rely on tags for now.
						Tags: fmt.Sprintf("openshiftClusterID=%s", clusterID),
					},
				},
			},
		}
		if len(addressPairs) > 0 {
			port.AllowedAddressPairs = addressPairs
		}
	}

	additionalPorts := make([]capo.PortOpts, 0, len(mpool.AdditionalNetworkIDs))
	for _, networkID := range mpool.AdditionalNetworkIDs {
		additionalPorts = append(additionalPorts, capo.PortOpts{
			Network: &capo.NetworkFilter{
				ID: networkID,
			},
		})
	}

	securityGroups := []capo.SecurityGroupFilter{
		{
			// Bootstrap and Master share the same security group
			Name: fmt.Sprintf("%s-master", clusterID),
		},
	}

	for _, securityGroup := range mpool.AdditionalSecurityGroupIDs {
		securityGroups = append(securityGroups, capo.SecurityGroupFilter{ID: securityGroup})
	}

	// FIXME: Uncomment when the server group rework merged
	// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/pull/1779
	// serverGroupName := clusterID + "-" + role
	spec := capo.OpenStackMachineSpec{
		CloudName: CloudName,
		Flavor:    mpool.FlavorName,
		IdentityRef: &capo.OpenStackIdentityReference{
			Kind: "Secret",
			Name: clusterID + "-cloud-config",
		},
		// FIXME(stephenfin): We probably want a FIP for bootstrap?
		// TODO: This is an image name. Migrate to a filter with Name when API v1alpha8 is released.
		Image:          osImage,
		Ports:          append([]capo.PortOpts{port}, additionalPorts...),
		SecurityGroups: securityGroups,
		// FIXME: Uncomment when the server group rework merged
		// https://github.com/kubernetes-sigs/cluster-api-provider-openstack/pull/1779
		//ServerGroup: *capo.ServerGroupFilter{
		//	"Name": serverGroupName,
		// },
		ServerMetadata: map[string]string{
			"Name":               fmt.Sprintf("%s-%s", clusterID, role),
			"openshiftClusterID": clusterID,
		},
		Trunk: trunkSupport,
		Tags: []string{
			fmt.Sprintf("openshiftClusterID=%s", clusterID),
		},
	}

	if mpool.RootVolume != nil {
		spec.RootVolume = &capo.RootVolume{
			Size:             mpool.RootVolume.Size,
			VolumeType:       failureDomain.RootVolume.VolumeType,
			AvailabilityZone: failureDomain.RootVolume.AvailabilityZone,
		}
	}

	return &spec, nil
}
