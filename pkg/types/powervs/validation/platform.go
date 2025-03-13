package validation

import (
	"fmt"
	"net/url"
	"regexp"

	"github.com/google/uuid"
	"k8s.io/apimachinery/pkg/util/sets"
	"k8s.io/apimachinery/pkg/util/validation/field"

	configv1 "github.com/openshift/api/config/v1"
	"github.com/openshift/installer/pkg/types/powervs"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *powervs.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	// validate Zone
	if p.Zone == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("zone"), "zone must be specified"))
		// Region checking is nonsense if Zone is invalid
		return allErrs
	} else if ok := powervs.ValidateZone(p.Zone); !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("zone"), p.Zone, powervs.ZoneNames()))
		// Region checking is nonsense if Zone is invalid
		return allErrs
	}

	// validate Region
	if p.Region == "" {
		p.Region = powervs.RegionFromZone(p.Zone)
	}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region not findable from specified zone"))
	} else if _, ok := powervs.Regions[p.Region]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("region"), p.Region, powervs.RegionShortNames()))
	}

	// validate DefaultMachinePlatform
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	// validate ServiceInstanceGUID
	if p.ServiceInstanceGUID != "" {
		_, err := uuid.Parse(p.ServiceInstanceGUID)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ServiceInstanceGUID"), p.ServiceInstanceGUID, "ServiceInstanceGUID must be a valid UUID"))
		}
	}
	if p.ServiceEndpoints != nil {
		allErrs = append(allErrs, validateServiceEndpoints(p.ServiceEndpoints, fldPath.Child("serviceEndpoints"))...)
	}

	return allErrs
}

// validateServiceEndpoints checks that the specified ServiceEndpoints are valid.
func validateServiceEndpoints(endpoints []configv1.PowerVSServiceEndpoint, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	knownEndpoints := sets.New[string]()
	for index, endpoint := range endpoints {
		fldp := fldPath.Index(index)
		if knownEndpoints.Has(endpoint.Name) {
			allErrs = append(allErrs, field.Duplicate(fldp.Child("name"), endpoint.Name))
		}
		knownEndpoints.Insert(endpoint.Name)

		if err := validateServiceURL(endpoint.URL); err != nil {
			allErrs = append(allErrs, field.Invalid(fldp.Child("url"), endpoint.URL, err.Error()))
		}
	}
	return allErrs
}

// schemeRE is used to check whether a string starts with a scheme (URI format).
var schemeRE = regexp.MustCompile("^([^:]+)://")

// versionPath is the regexp for a trailing API version in URL path ('/v1', '/v22/', etc.)
var versionPath = regexp.MustCompile(`(/v\d+[/]{0,1})$`)

// validateServiceURL checks that a string meets certain URI expectations.
func validateServiceURL(uri string) error {
	endpoint := uri
	httpsScheme := "https"

	// determine if the endpoint (uri) starts with an URI scheme
	// add 'https' scheme if not
	if !schemeRE.MatchString(endpoint) {
		endpoint = fmt.Sprintf("%s://%s", httpsScheme, endpoint)
	}

	// verify the endpoint meets the following criteria
	// 1. contains a hostname
	// 2. uses 'https' scheme
	// 3. contains no path or request parameters, except API version paths ('/v1')
	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("empty hostname provided, it cannot be empty")
	}
	// check the scheme in case one was provided and is not 'https' (we didn't set it above)
	if s := u.Scheme; s != httpsScheme {
		return fmt.Errorf("invalid scheme %s, only https is allowed", s)
	}
	// check that the path is empty ('/'), or only contains an API version ('/v1'), by using regexp to replace the API version and should result in empty string
	if r := u.RequestURI(); r != "/" && versionPath.ReplaceAllString(r, "") != "" {
		return fmt.Errorf("no path or request parameters can be provided, %q was provided", r)
	}

	return nil
}
