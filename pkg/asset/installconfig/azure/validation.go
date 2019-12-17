package azure

import (
	"context"
	"fmt"
	"net"

	aznetwork "github.com/Azure/azure-sdk-for-go/services/network/mgmt/2018-12-01/network"
	aztypes "github.com/openshift/installer/pkg/types/azure"

	"github.com/openshift/installer/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate executes platform-specific validation.
func Validate(client API, ic *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	allErrs = append(allErrs, validateNetworks(client, ic.Azure, ic.Networking.MachineNetwork, field.NewPath("platform").Child("azure"))...)
	return allErrs.ToAggregate()
}

// validateNetworks checks that the user-provided VNet and subnets are valid.
func validateNetworks(client API, p *aztypes.Platform, machineNetworks []types.MachineNetworkEntry, fieldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.VirtualNetwork != "" {
		_, err := client.GetVirtualNetwork(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("virtualNetwork"), p.VirtualNetwork, err.Error()))
		}

		computeSubnet, err := client.GetComputeSubnet(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork, p.ComputeSubnet)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("computeSubnet"), p.ComputeSubnet, "failed to retrieve compute subnet"))
		}

		allErrs = append(allErrs, validateSubnet(client, fieldPath.Child("computeSubnet"), computeSubnet, p.ComputeSubnet, machineNetworks)...)

		controlPlaneSubnet, err := client.GetControlPlaneSubnet(context.TODO(), p.NetworkResourceGroupName, p.VirtualNetwork, p.ControlPlaneSubnet)
		if err != nil {
			return append(allErrs, field.Invalid(fieldPath.Child("controlPlaneSubnet"), p.ControlPlaneSubnet, "failed to retrieve control plane subnet"))
		}

		allErrs = append(allErrs, validateSubnet(client, fieldPath.Child("controlPlaneSubnet"), controlPlaneSubnet, p.ControlPlaneSubnet, machineNetworks)...)
	}

	return allErrs
}

// validateSubnet checks that the subnet is in the same network as the machine CIDR
func validateSubnet(client API, fieldPath *field.Path, subnet *aznetwork.Subnet, subnetName string, networks []types.MachineNetworkEntry) field.ErrorList {
	allErrs := field.ErrorList{}

	subnetIP, _, err := net.ParseCIDR(*subnet.AddressPrefix)
	if err != nil {
		return append(allErrs, field.Invalid(fieldPath, subnetName, "unable to parse subnet CIDR"))
	}

	allErrs = append(allErrs, validateMachineNetworksContainIP(fieldPath, networks, *subnet.Name, subnetIP)...)
	return allErrs
}

func validateMachineNetworksContainIP(fldPath *field.Path, networks []types.MachineNetworkEntry, subnetName string, ip net.IP) field.ErrorList {
	for _, network := range networks {
		if network.CIDR.Contains(ip) {
			return nil
		}
	}
	return field.ErrorList{field.Invalid(fldPath, subnetName, fmt.Sprintf("subnet %s address prefix is outside of the specified machine networks", ip))}
}
