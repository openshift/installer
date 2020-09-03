package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"sort"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *aws.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	allErrs = append(allErrs, validateServiceEndpoints(p.ServiceEndpoints, fldPath.Child("serviceEndpoints"))...)
	allErrs = append(allErrs, validateUserTags(p.UserTags, fldPath.Child("userTags"))...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}

func validateUserTags(tags map[string]string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(tags) == 0 {
		return allErrs
	}
	keys := make([]string, 0, len(tags))
	for k := range tags {
		keys = append(keys, k)
	}
	sort.Strings(keys)
	for _, key := range keys {
		if strings.EqualFold(key, "Name") {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(key), tags[key], "Name key is not allowed for user defined tags"))
		}
		if strings.HasPrefix(key, "kubernetes.io/cluster/") {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(key), tags[key], "Keys with prefix 'kubernetes.io/cluster/' are not allowed for user defined tags"))
		}
	}
	return allErrs
}

func validateServiceEndpoints(endpoints []aws.ServiceEndpoint, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	tracker := map[string]int{}
	for idx, e := range endpoints {
		fldp := fldPath.Index(idx)
		if eidx, ok := tracker[e.Name]; ok {
			allErrs = append(allErrs, field.Invalid(fldp.Child("name"), e.Name, fmt.Sprintf("duplicate service endpoint not allowed for %s, service endpoint already defined at %s", e.Name, fldPath.Index(eidx))))
		} else {
			tracker[e.Name] = idx
		}

		if err := validateServiceURL(e.URL); err != nil {
			allErrs = append(allErrs, field.Invalid(fldp.Child("url"), e.URL, err.Error()))
		}
	}
	return allErrs
}

var schemeRE = regexp.MustCompile("^([^:]+)://")

func validateServiceURL(uri string) error {
	endpoint := uri
	if !schemeRE.MatchString(endpoint) {
		scheme := "https"
		endpoint = fmt.Sprintf("%s://%s", scheme, endpoint)
	}

	u, err := url.Parse(endpoint)
	if err != nil {
		return err
	}
	if u.Hostname() == "" {
		return fmt.Errorf("host cannot be empty, empty host provided")
	}
	if s := u.Scheme; s != "https" {
		return fmt.Errorf("invalid scheme %s, only https allowed", s)
	}
	if r := u.RequestURI(); r != "/" {
		return fmt.Errorf("no path or request parameters must be provided, %q was provided", r)
	}

	return nil
}
