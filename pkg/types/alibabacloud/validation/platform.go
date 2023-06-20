package validation

import (
	"fmt"
	"net"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/alibabacloud"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *alibabacloud.Platform, n *types.Networking, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	logrus.Warn("Alibaba Cloud is deprecated and will be removed in a future OpenShift version. Please reach out to your Red Hat Support or Technical Account Manager for more information.")

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	if p.PrivateZoneID != "" {
		if p.VpcID == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("vpcID"), "when using existing privatezones, an existing VPC must be used"))
		}
	}

	if len(p.VSwitchIDs) > 0 {
		allErrs = append(allErrs, validateVSwitches(p, fldPath)...)
	}

	allErrs = append(allErrs, validateMetadataServerIP(n)...)

	return allErrs
}

func validateVSwitches(p *alibabacloud.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.VpcID == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("vpcID"), "when using existing VSwitches, an existing VPC must be used"))
	}

	vswitchIDs := map[string]bool{}
	for idx, vswitchID := range p.VSwitchIDs {
		if vswitchIDs[vswitchID] {
			allErrs = append(allErrs, field.Duplicate(fldPath.Child("vswitchIDs").Index(idx), vswitchID))
		} else {
			vswitchIDs[vswitchID] = true
		}
	}
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
