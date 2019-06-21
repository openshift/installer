package validation

import (
	"sort"

	"github.com/openshift/installer/pkg/types/gcp"
	"k8s.io/apimachinery/pkg/util/validation/field"
)

var (
	// Regions is a map of known GCP regions. The key of the map is
	// the short name of the region. The value of the map is the long
	// name of the region.
	Regions = map[string]string{
		"northamerica-northeast1": "Montréal",
		"us-central":              "Iowa",
		"us-west2":                "Los Angeles",
		"us-east1":                "South Carolina",
		"us-east4":                "Northern Virginia",
		"southamerica-east1":      "São Paulo",
		"europe-west":             "Belgium",
		"europe-west2":            "London",
		"europe-west3":            "Frankfurt",
		"europe-west6":            "Zürich",
		"asia-northeast1":         "Tokyo",
		"asia-northeast2":         "Osaka",
		"asia-east2":              "Hong Kong",
		"asia-south1":             "Mumbai",
		"australia-southeast1":    "Sydney",
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
func ValidatePlatform(p *gcp.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if _, ok := Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, validRegionValues))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}
