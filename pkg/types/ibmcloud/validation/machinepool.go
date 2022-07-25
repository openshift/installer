package validation

import (
	"fmt"
	"strings"

	"github.com/IBM-Cloud/bluemix-go/crn"
	"github.com/openshift/installer/pkg/types/ibmcloud"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// ValidateMachinePool validates the MachinePool.
func ValidateMachinePool(platform *ibmcloud.Platform, mp *ibmcloud.MachinePool, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for i, zone := range mp.Zones {
		if !strings.HasPrefix(zone, platform.Region) {
			allErrs = append(allErrs, field.Invalid(path.Child("zones").Index(i), zone, fmt.Sprintf("zone not in configured region (%s)", platform.Region)))
		}
	}

	if mp.DedicatedHosts != nil {
		allErrs = append(allErrs, validateDedicatedHosts(mp.DedicatedHosts, mp.InstanceType, mp.Zones, path.Child("dedicatedHosts"))...)

		if mp.InstanceType == "" {
			allErrs = append(allErrs, field.Invalid(path.Child("type"), mp.InstanceType, "type is required, default type not supported for dedicated hosts"))
		}
	}

	if mp.BootVolume != nil {
		allErrs = append(allErrs, validateBootVolume(mp.BootVolume, path.Child("bootVolume"))...)
	}
	return allErrs
}

func validateBootVolume(bv *ibmcloud.BootVolume, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if bv.EncryptionKey != "" {
		_, parseErr := crn.Parse(bv.EncryptionKey)
		if parseErr != nil {
			allErrs = append(allErrs, field.Invalid(path.Child("encryptionKey"), bv.EncryptionKey, "encryptionKey is not a valid IBM CRN"))
		}
	}
	return allErrs
}

func validateDedicatedHosts(dhosts []ibmcloud.DedicatedHost, itype string, zones []string, path *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Length of dedicated hosts must match platform zones
	if len(dhosts) != len(zones) {
		allErrs = append(allErrs, field.Invalid(path, dhosts, fmt.Sprintf("number of dedicated hosts does not match list of zones (%s)", zones)))
	}

	for i, dhost := range dhosts {
		// Dedicated host name or profile is required
		if dhost.Name == "" && dhost.Profile == "" {
			allErrs = append(allErrs, field.Invalid(path.Index(i), dhost.Profile, "name or profile must be set"))
		}

		// Instance type must be in the same profile family as dedicated host
		if dhost.Profile != "" && itype != "" {
			if strings.Split(dhost.Profile, "-")[0] != strings.Split(itype, "-")[0] {
				allErrs = append(allErrs, field.Invalid(path.Index(i).Child("profile"), dhost.Profile, fmt.Sprintf("profile does not support expected instance type (%s)", itype)))
			}
		}
	}

	return allErrs
}
