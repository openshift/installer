package validation

import (
	"sort"

	"github.com/openshift/installer/pkg/types/azure"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	// Regions is a map of the known azure regions. The key of the map is
	// the short name of the region. The value of the map is the long
	// name of the region.
	Regions = map[string]string{
		"eastus": "East US",
		"westus": "West US",
	}

	validRegionValues = func() []string {
		validValues := make([]string, len(Regions))
		i := 0
		for r := range Regions {
			validValues[i] = r
			i++
		}
		sort.Strings(validValues)
		return validValues
	}()
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *azure.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if _, ok := Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, validRegionValues))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
