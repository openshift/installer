package validation

import (
	"fmt"
	"os"
	"regexp"
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

	// userLabelKeyRegex is for verifying that the label key contains only allowed characters.
	userLabelKeyRegex = regexp.MustCompile(`^[a-z][0-9a-z_-]{1,62}$`)

	// userLabelValueRegex is for verifying that the label value contains only allowed characters.
	userLabelValueRegex = regexp.MustCompile(`^[0-9a-z_-]{1,63}$`)

	// userLabelKeyPrefixRegex is for verifying that the label key does not contain restricted prefixes.
	userLabelKeyPrefixRegex = regexp.MustCompile(`^(?i)(kubernetes\-io|openshift\-io)`)

	// userTagKeyRegex is for verifying that the tag key contains only allowed characters.
	userTagKeyRegex = regexp.MustCompile(`^[a-zA-Z0-9]([0-9A-Za-z_.-]{0,61}[a-zA-Z0-9])?$`)

	// userTagValueRegex is for verifying that the tag value contains only allowed characters.
	userTagValueRegex = regexp.MustCompile(`^[a-zA-Z0-9]([0-9A-Za-z_.@%=+:,*#&()\[\]{}\-\s]{0,61}[a-zA-Z0-9])?$`)

	// userTagParentIDRegex is for verifying that the tag parentID contains only allowed characters.
	userTagParentIDRegex = regexp.MustCompile(`(^[1-9][0-9]{0,31}$)|(^[a-z][a-z0-9-]{4,28}[a-z0-9]$)`)
)

// maxUserLabelLimit is the maximum userLabels that can be configured as defined in openshift/api.
// https://github.com/openshift/api/commit/ae73a19d05c35068af16c9aeff375d0b7c936a8a#diff-07b264a49084976b670fb699badaca1795027d6ea732a99226a5388104f6174fR592-R602
const maxUserLabelLimit = 32

// maxUserTagLimit is the maximum userTags that can be configured as defined in openshift/api.
// https://github.com/openshift/api/commit/ae73a19d05c35068af16c9aeff375d0b7c936a8a#diff-07b264a49084976b670fb699badaca1795027d6ea732a99226a5388104f6174fR604-R613
const maxUserTagLimit = 50

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

	// check if configured userLabels are valid.
	allErrs = append(allErrs, validateUserLabels(p.UserLabels, fldPath.Child("userLabels"))...)

	// check if configured userTags are valid.
	allErrs = append(allErrs, validateUserTags(p.UserTags, fldPath.Child("userTags"))...)

	return allErrs
}

// validateUserLabels verifies if configured number of UserLabels is not more than
// allowed limit and the label keys and values are valid.
func validateUserLabels(labels []gcp.UserLabel, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(labels) == 0 {
		return allErrs
	}

	if len(labels) > maxUserLabelLimit {
		allErrs = append(allErrs, field.TooMany(fldPath, len(labels), maxUserLabelLimit))
	}

	for _, label := range labels {
		if err := validateLabel(label.Key, label.Value); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(label.Key), label.Value, err.Error()))
		}
	}
	return allErrs
}

// validateLabel checks the following to ensure that the label configured is acceptable.
//   - The key and value contain only allowed characters.
//   - The key is not empty and at most 63 characters and starts with a lowercase letter.
//   - The value is not empty and at most 63 characters.
//   - The key and value must contain only lowercase letters, numeric characters,
//     underscores, and dashes.
//   - The key cannot be Name or have kubernetes.io, openshift.io prefixes.
func validateLabel(key, value string) error {
	if !userLabelKeyRegex.MatchString(key) {
		return fmt.Errorf("label key is invalid or contains invalid characters. Label key can have a maximum of 63 characters and cannot be empty. Label key must begin with a lowercase letter, and must contain only lowercase letters, numeric characters, and the following special characters `_-`")
	}
	if !userLabelValueRegex.MatchString(value) {
		return fmt.Errorf("label value is invalid or contains invalid characters. Label value can have a maximum of 63 characters and cannot be empty. Value must contain only lowercase letters, numeric characters, and the following special characters `_-`")
	}
	if userLabelKeyPrefixRegex.MatchString(key) {
		return fmt.Errorf("label key contains restricted prefix. Label key cannot have `kubernetes-io`, `openshift-io` prefixes")
	}
	return nil
}

// validateUserTags verifies if configured number of UserTags is not more than
// allowed limit and the tag keys and values are valid.
func validateUserTags(tags []gcp.UserTag, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(tags) == 0 {
		return allErrs
	}

	if len(tags) > maxUserTagLimit {
		allErrs = append(allErrs, field.TooMany(fldPath, len(tags), maxUserTagLimit))
	}

	for _, tag := range tags {
		if err := validateTag(tag.ParentID, tag.Key, tag.Value); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(tag.Key), tag.Value, err.Error()))
		}
	}

	return allErrs
}

// validateTag checks the following to ensure that the tag configured is acceptable. Though
// the criteria is for tag resources to pre-exist, tags will be validated to catch the
// error much earlier.
//   - The key and value contain only allowed characters.
//   - The key and value is not empty and can have at most 63 characters.
//   - The ParentID can be either OrganizationID or ProjectID.
//   - OrganizationID must consist of decimal numbers, and cannot have leading zeroes.
//   - ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers,
//     and hyphens, and must start with a letter, and cannot end with a hyphen.
func validateTag(parentID, key, value string) error {
	if !userTagParentIDRegex.MatchString(parentID) {
		return fmt.Errorf("tag parentID is invalid or contains invalid characters. ParentID can have a maximum of 32 characters and cannot be empty. ParentID can be either OrganizationID or ProjectID. OrganizationID must consist of decimal numbers, and cannot have leading zeroes and ProjectID must be 6 to 30 characters in length, can only contain lowercase letters, numbers, and hyphens, and must start with a letter, and cannot end with a hyphen")
	}
	if !userTagKeyRegex.MatchString(key) {
		return fmt.Errorf("tag key is invalid or contains invalid characters. Tag key can have a maximum of 63 characters and cannot be empty. Tag key must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `._-`")
	}
	if !userTagValueRegex.MatchString(value) {
		return fmt.Errorf("tag value is invalid or contains invalid characters. Tag value can have a maximum of 63 characters and cannot be empty. Tag value must begin and end with an alphanumeric character, and must contain only uppercase, lowercase alphanumeric characters, and the following special characters `_-.@%%=+:,*#&(){}[]` and spaces")
	}
	return nil
}
