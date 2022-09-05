package validation

import (
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

	allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, "default", fldPath.Child("defaultMachinePlatform"))...)

	return allErrs
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
