package validation

import (
	"fmt"
	"net"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *alibabacloud.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	allErrs = append(allErrs, validateMetadataServerIP(n)...)

	return allErrs
}

func validateMetadataServerIP(n *types.Networking) field.ErrorList {
	allErrs := field.ErrorList{}
	metadataServerIP := "100.100.100.200"
	networkPath := field.NewPath("networking")
	machineNetworkPath := networkPath.Child("machineNetwork")
	clusterNetworkPath := networkPath.Child("clusterNetwork")
	serviceNetworkPath := networkPath.Child("serviceNetwork")

	for _, network := range n.MachineNetwork {
		if network.CIDR.Contains(net.ParseIP(metadataServerIP)) {
			allErrs = append(allErrs, field.Invalid(machineNetworkPath, network.CIDR.String(), fmt.Errorf("contains %s which is reserved for the metadata service", metadataServerIP).Error()))
		}
	}
	for _, network := range n.ClusterNetwork {
		if network.CIDR.Contains(net.ParseIP(metadataServerIP)) {
			allErrs = append(allErrs, field.Invalid(clusterNetworkPath, network.CIDR.String(), fmt.Errorf("contains %s which is reserved for the metadata service", metadataServerIP).Error()))
		}
	}
	for _, network := range n.ServiceNetwork {
		if network.Contains(net.ParseIP(metadataServerIP)) {
			allErrs = append(allErrs, field.Invalid(serviceNetworkPath, network.String(), fmt.Errorf("contains %s which is reserved for the metadata service", metadataServerIP).Error()))
		}
	}

	return allErrs
}
