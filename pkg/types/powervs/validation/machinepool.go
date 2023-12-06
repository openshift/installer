package validation

import (
	"fmt"
	"regexp"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/intstr"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/powervs"
)

var validSMTLevels = sets.New[string]("1", "2", "3", "4", "5", "6", "7", "8", "on", "off")

// ValidateMachinePool checks that the specified machine pool is valid.
func ValidateMachinePool(p *powervs.MachinePool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// Validate VolumeIDs
	volumes := make(map[string]bool)
	for i, volumeID := range p.VolumeIDs {
		_, err := uuid.Parse(volumeID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("volumeIDs").Index(i), volumeID, "volume ID must be a valid UUID"))
		}
		if _, ok := volumes[volumeID]; ok {
			allErrs = append(allErrs, field.Duplicate(fldPath.Child("volumeIDs").Index(i), volumeID))
			continue
		}
		volumes[volumeID] = true
	}

	// Validate ProcType
	if p.ProcType != "" {
		procTypes := sets.NewString("Shared", "Dedicated", "Capped")
		if !procTypes.Has(string(p.ProcType)) {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("procType"), p.ProcType, procTypes.List()))
		}
	}

	// Validate PowerVS Memory and Processors max and min limits
	// More details about PowerVS limits.
	// https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-pricing-virtual-server
	// Validate Processors
	var processors float64
	var err error
	switch p.Processors.Type {
	case intstr.Int:
		processors = float64(p.Processors.IntVal)
	case intstr.String:
		processors, err = strconv.ParseFloat(p.Processors.StrVal, 64)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), p.Processors.StrVal, "processors must be a valid floating point number"))
		}
		if err == nil && !regexp.MustCompile(`^(^$|[0]*|25[0]*|5[0]*|75[0]*)$`).MatchString(strings.Split(p.Processors.StrVal, ".")[1]) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), processors, "processors must be in increments of .25"))
		}
	}
	if processors != 0 {
		switch p.ProcType {
		case "Shared", "Capped", "":
			if processors < 0.5 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), processors, "minimum number of processors must be .5 cores for capped or shared ProcType"))
			}
		case "Dedicated":
			if processors < 1 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), processors, "minimum number of processors must be from 1 core for Dedicated ProcType"))
			}
		}
	}

	// Validate SMTLevel
	if p.SMTLevel != "" && !validSMTLevels.Has(p.SMTLevel) {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("smtLevel"), p.SMTLevel, fmt.Sprintf("Valid SMT Levels are %s", sets.List(validSMTLevels))))
	}

	// Validate SysType
	if p.SysType != "" {
		if !powervs.AllKnownSysTypes().Has(p.SysType) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("sysType"), p.SysType, "unknown system type specified"))
		}
	}

	// Validate for Maximum Memory and Processors limits
	if p.MemoryGiB != 0 {
		if p.MemoryGiB < 4 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.MemoryGiB, "memory must be an integer number of GB that is at least 4"))
		}
	}

	// Validate Mimimum Memory limit of machinepool.
	s922TypeRegex := regexp.MustCompile(`^s922(-.*|)$`)
	e980TypeRegex := regexp.MustCompile(`^e980$`)

	switch {
	case e980TypeRegex.MatchString(p.SysType):
		if p.MemoryGiB > 15307 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.MemoryGiB, "maximum memory limit for the e980 SysType is 15307GiB"))
		}

		if processors > 143 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), processors, "maximum processors limit for e980 SysType is 143 cores"))
		}
	case s922TypeRegex.MatchString(p.SysType):
		if p.MemoryGiB > 942 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("memory"), p.MemoryGiB, "maximum memory limit for the s922 SysType is 942GiB"))
		}

		if processors > 15 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("processors"), processors, "maximum processors limit for s922 SysType is 15 cores"))
		}
	}
	return allErrs
}
