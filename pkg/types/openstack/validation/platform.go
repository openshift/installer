package validation

import (
	"errors"
	"strings"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/validate"

	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *openstack.Platform, n *types.Networking, fldPath *field.Path, c *types.InstallConfig) field.ErrorList {
	var allErrs field.ErrorList

	for _, ip := range p.ExternalDNS {
		if err := validate.IP(ip); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("externalDNS"), p.ExternalDNS, err.Error()))
		}
	}

	err := validateFailureDomainsPlatform(p, c)
	if err != nil {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("failureDomains"), p.FailureDomains, err.Error()))
	}

	allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, "default", fldPath.Child("defaultMachinePlatform"))...)

	return allErrs
}

func validateFailureDomainsPlatform(p *openstack.Platform, c *types.InstallConfig) error {
	for _, domain := range p.FailureDomains {
		if domain.Name == "" {
			return errors.New(("must specify a failure domain name"))
		}
		if domain.Subnet == "" {
			return errors.New(("must specify a failure domain subnet"))
		}
	}

	controlPlaneFailureDomains := c.ControlPlane.Platform.OpenStack.FailureDomainNames != nil
	computeFailureDomains := c.Compute[0].Platform.OpenStack.FailureDomainNames != nil
	if controlPlaneFailureDomains != computeFailureDomains {
		return errors.New("must specify failure domains for both control plane and compute")
	}

	if len(p.FailureDomains) > 0 {
		if p.MachinesSubnet == "" {
			return errors.New(("must specify a machinesSubnet when failure domains are specified"))
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
