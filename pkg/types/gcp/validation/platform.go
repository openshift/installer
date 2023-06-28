package validation

import (
	"os"
	"sort"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/validate"
)

var (
	// Regions is a map of known GCP regions. The key of the map is
	// the short name of the region. The value of the map is the long
	// name of the region.
	Regions = map[string]string{
		// List from: https://cloud.google.com/compute/docs/regions-zones/
		"asia-east1":              "Changhua County, Taiwan",
		"asia-east2":              "Hong Kong",
		"asia-northeast1":         "Tokyo, Japan",
		"asia-northeast2":         "Osaka, Japan",
		"asia-northeast3":         "Seoul, South Korea",
		"asia-south1":             "Mumbai, India",
		"asia-south2":             "Delhi, India",
		"asia-southeast1":         "Jurong West, Singapore",
		"asia-southeast2":         "Jakarta, Indonesia",
		"australia-southeast1":    "Sydney, Australia",
		"australia-southeast2":    "Melbourne, Australia",
		"europe-central2":         "Warsaw, Poland",
		"europe-north1":           "Hamina, Finland",
		"europe-west1":            "St. Ghislain, Belgium",
		"europe-west2":            "London, England, UK",
		"europe-west3":            "Frankfurt, Germany",
		"europe-west4":            "Eemshaven, Netherlands",
		"europe-west6":            "Zürich, Switzerland",
		"europe-west8":            "Milan, Italy",
		"europe-west9":            "Paris, France",
		"europe-west12":           "Turin, Italy",
		"europe-southwest1":       "Madrid, Spain",
		"me-central1":             "Doha, Qatar, Middle East",
		"me-west1":                "Tel Aviv, Israel",
		"northamerica-northeast1": "Montréal, Québec, Canada",
		"northamerica-northeast2": "Toronto, Ontario, Canada",
		"southamerica-east1":      "São Paulo, Brazil",
		"southamerica-west1":      "Santiago, Chile",
		"us-central1":             "Council Bluffs, Iowa, USA",
		"us-east1":                "Moncks Corner, South Carolina, USA",
		"us-east4":                "Ashburn, Virginia, USA",
		"us-east5":                "Columbus, Ohio, USA",
		"us-south1":               "Dallas, Texas, USA",
		"us-west1":                "The Dalles, Oregon, USA",
		"us-west2":                "Los Angeles, California, USA",
		"us-west3":                "Salt Lake City, Utah, USA",
		"us-west4":                "Las Vegas, Nevada, USA",
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
func ValidatePlatform(p *gcp.Platform, fldPath *field.Path, ic *types.InstallConfig) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "must provide a region"))
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
		allErrs = append(allErrs, ValidateDefaultDiskType(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	if p.NetworkProjectID != "" {
		if p.Network == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("network"), "must provide a network when a networkProjectID is specified"))
		}
		if ic.CredentialsMode != types.ManualCredentialsMode && ic.CredentialsMode != types.PassthroughCredentialsMode {
			allErrs = append(allErrs, field.NotSupported(fldPath.Child("credentialsMode"),
				ic.CredentialsMode, []string{string(types.ManualCredentialsMode), string(types.PassthroughCredentialsMode)}))
		}
	}
	if p.Network != "" {
		if p.ComputeSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("computeSubnet"), "must provide a compute subnet when a network is specified"))
		}
		if p.ControlPlaneSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("controlPlaneSubnet"), "must provide a control plane subnet when a network is specified"))
		}
	}

	if (p.ComputeSubnet != "" || p.ControlPlaneSubnet != "") && p.Network == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("network"), "must provide a VPC network when supplying subnets"))
	}

	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); ok && oi != "" && len(p.Licenses) > 0 {
		allErrs = append(allErrs, field.Forbidden(fldPath.Child("licenses"), "the use of custom image licenses is forbidden if an OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE is specified"))
	}

	for i, license := range p.Licenses {
		if validate.URIWithProtocol(license, "https") != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("licenses").Index(i), license, "licenses must be URLs (https) only"))
		}
	}

	return allErrs
}
