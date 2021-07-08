package validation

import (
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/ovirt"
	"github.com/openshift/installer/pkg/validate"
)

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *ovirt.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.CPU != nil {
		if p.CPU.Cores <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("cores"), p.CPU.Cores, "CPU cores must be positive"))
		}
		if p.CPU.Sockets <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("sockets"), p.CPU.Sockets, "CPU sockets must be positive"))
		}
	}

	if p.MemoryMB < 0 {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("memoryMB"), p.MemoryMB, "Memory value must be nonnegative"))
	}

	if p.VMType != "" && !ValidVMType(p.VMType) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("vmType"), p.VMType, fmt.Sprintf("VM type must be one of %s", supportedVMTypes())))
	}

	if p.InstanceTypeID != "" {
		if p.CPU != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("instanceTypeID"), p.InstanceTypeID, "mixing instanceTypeID and CPU is not supported"))
		}
		if p.MemoryMB > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("instanceTypeID"), p.InstanceTypeID, "mixing instanceTypeID and Memory  is not supported"))
		}
		if err := validate.UUID(p.InstanceTypeID); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("instanceTypeID"), p.InstanceTypeID, err.Error()))
		}
	}

	if p.OSDisk != nil {
		if p.OSDisk.SizeGB <= 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("sizeGB"), p.OSDisk.SizeGB, "disk size must be positive"))
		}
	}

	if p.AutoPinningPolicy != "" && !ValidAutoPinningPolicy(p.AutoPinningPolicy) {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("autoPinningPolicy"), p.AutoPinningPolicy,
			[]string{string(ovirt.AutoPinningNone), string(ovirt.AutoPinningResizeAndPin)}))
	}

	if p.Hugepages > 0 {
		if p.Hugepages != 2048 && p.Hugepages != 1048576 {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("hugepages"), p.Hugepages,
				[]string{string(ovirt.Hugepages2MB), string(ovirt.Hugepages1GB)}))
		}
	}

	return allErrs
}

// ValidVMType returns true if the vmType is supported.
func ValidVMType(vmType ovirt.VMType) bool {
	for _, v := range supportedVMTypes() {
		if vmType == v {
			return true
		}
	}
	return false
}

// supportedVMTypes returns a slice of all supported VMTypes.
func supportedVMTypes() []ovirt.VMType {
	return []ovirt.VMType{
		ovirt.VMTypeDesktop,
		ovirt.VMTypeServer,
		ovirt.VMTypeHighPerformance,
	}
}

// ValidAutoPinningPolicy returns true if the AutoPinningPolicy is supported.
func ValidAutoPinningPolicy(autoPinningPolicy ovirt.AutoPinningPolicy) bool {
	for _, v := range supportedAutoPinningPolicies() {
		if autoPinningPolicy == v {
			return true
		}
	}
	return false
}

// supportedAutoPinningPolicies returns a slice of all supported AutoPinningPolicy.
func supportedAutoPinningPolicies() []ovirt.AutoPinningPolicy {
	return []ovirt.AutoPinningPolicy{
		ovirt.AutoPinningNone,
		ovirt.AutoPinningResizeAndPin,
	}
}
