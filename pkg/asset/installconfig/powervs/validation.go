package powervs

import (
	"context"
	"fmt"

	"github.com/openshift/installer/pkg/types"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

// Validate executes platform specific validation/
func Validate(config *types.InstallConfig) error {
	allErrs := field.ErrorList{}

	if config.Platform.PowerVS == nil {
		allErrs = append(allErrs, field.Required(field.NewPath("platform", "powervs"), "Power VS Validation requires a Power VS platform configuration."))
	} else {
		if config.ControlPlane != nil {
			fldPath := field.NewPath("controlPlane")
			allErrs = append(allErrs, validateMachinePool(fldPath, config.ControlPlane)...)
		}
		for idx, compute := range config.Compute {
			fldPath := field.NewPath("compute").Index(idx)
			allErrs = append(allErrs, validateMachinePool(fldPath, &compute)...)
		}
	}
	return allErrs.ToAggregate()
}

func validateMachinePool(fldPath *field.Path, machinePool *types.MachinePool) field.ErrorList {
	allErrs := field.ErrorList{}
	if machinePool.Architecture != "ppc64le" {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("architecture"), machinePool.Architecture, []string{"ppc64le"}))
	}
	return allErrs
}

// ValidatePreExistingPublicDNS ensure no pre-existing DNS record exists in the CIS
// DNS zone for cluster's Kubernetes API.
func ValidatePreExistingPublicDNS(client API, ic *types.InstallConfig, metadata *Metadata) error {
	allErrs := field.ErrorList{}

	// Get CIS CRN
	crn, err := metadata.CISInstanceCRN(context.TODO())
	if err != nil {
		return err
	}

	// Get CIS zone ID by name
	zoneID, err := client.GetDNSZoneIDByName(context.TODO(), ic.BaseDomain)
	if err != nil {
		return append(allErrs, field.InternalError(field.NewPath("baseDomain"), err)).ToAggregate()
	}

	// Search for existing records
	recordNames := [...]string{fmt.Sprintf("api.%s", ic.ClusterDomain()), fmt.Sprintf("api-int.%s", ic.ClusterDomain())}
	for _, recordName := range recordNames {
		records, err := client.GetDNSRecordsByName(context.TODO(), crn, zoneID, recordName)
		if err != nil {
			allErrs = append(allErrs, field.InternalError(field.NewPath("baseDomain"), err))
		}

		// DNS record exists
		if len(records) != 0 {
			allErrs = append(allErrs, field.Duplicate(field.NewPath("baseDomain"), fmt.Sprintf("record %s already exists in CIS zone (%s) and might be in use by another cluster, please remove it to continue", recordName, zoneID)))
		}
	}

	return allErrs.ToAggregate()
}
