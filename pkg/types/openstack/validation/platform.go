package validation

import (
	"errors"
	"net"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateClusterName(c.ObjectMeta.Name)...)

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	err := validateVIP(p.APIVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("apiVIP"), p.APIVIP, err.Error()))
	}

	err = validateVIP(p.IngressVIP, n)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("ingressVIP"), p.IngressVIP, err.Error()))
	}

	allErrs = append(allErrs, validateDefaultMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)

	return allErrs
}

// validateVIP is a convenience function for validating VIP port and usage
func validateVIP(vip string, n *types.Networking) error {
	if vip != "" {
		if err := validate.IP(vip); err != nil {
			return err
		}

		if !n.MachineNetwork[0].CIDR.Contains(net.ParseIP(vip)) {
			return errors.New("IP is not in the machineNetwork")
		}
	}
	return nil
}

func validateClusterName(name string) (allErrs field.ErrorList) {
	if len(name) > 14 {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), name, "cluster name is too long, please restrict it to 14 characters"))
	}

	if strings.Contains(name, ".") {
		allErrs = append(allErrs, field.Invalid(field.NewPath("metadata", "name"), name, "cluster name can't contain \".\" character"))
	}

	return
}
