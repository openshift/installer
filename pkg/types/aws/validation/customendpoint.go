package validation

import (
	"sort"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

var (
	supportedServices = map[string]bool{"ec2": true, "elb": true, "iam": true, "route53": true, "s3": true}

	validServices = func() []string {
		validValues := make([]string, len(supportedServices))
		i := 0
		for r := range supportedServices {
			validValues[i] = r
			i++
		}
		sort.Strings(validValues)
		return validValues
	}()
)

// ValidateCustomEndpoints validates list of Custom Endpoints and appends the list of errors
func ValidateCustomEndpoints(endpoints *[]aws.CustomEndpoint, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	for _, endpoint := range *endpoints {
		allErrs = append(allErrs, ValidateCustomEndpoint(&endpoint, fldPath.Child("CustomEndpoint"))...)
	}
	return allErrs
}

// ValidateCustomEndpoint validates a Custom Endpoint object checking in supported service list etc.
// and appends the list of errors
func ValidateCustomEndpoint(endpoint *aws.CustomEndpoint, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if endpoint.Service == "" {
		// Checks for empty custom endpoint struct
		return allErrs
	}
	if val, ok := supportedServices[endpoint.Service]; !ok || !val {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("Service"), endpoint.Service, validServices))
	}
	if len(endpoint.URL) == 0 {
		allErrs = append(allErrs, field.Required(fldPath.Child("URL"), "Service Override URL cannot be empty"))
	}
	return allErrs
}
