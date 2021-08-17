package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types/aws"
)

// tagRegex is used to check that the keys and values of a tag contain only valid characters.
var tagRegex = regexp.MustCompile(`^[0-9A-Za-z_.:/=+-@]*$`)

// kubernetesNamespaceRegex is used to check that a tag key is not in the kubernetes.io namespace.
var kubernetesNamespaceRegex = regexp.MustCompile(`^([^/]*\.)?kubernetes.io/`)

// openshiftNamespaceRegex is used to check that a tag key is not in the openshift.io namespace.
var openshiftNamespaceRegex = regexp.MustCompile(`^([^/]*\.)?openshift.io/`)

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *aws.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	if p.HostedZone != "" {
		if len(p.Subnets) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostedZone"), p.HostedZone, "may not use an existing hosted zone when not using existing subnets"))
		}
	}

	allErrs = append(allErrs, validateServiceEndpoints(p.ServiceEndpoints, fldPath.Child("serviceEndpoints"))...)
	allErrs = append(allErrs, validateUserTags(p.UserTags, p.ExperimentalPropagateUserTag, fldPath.Child("userTags"))...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}
	return allErrs
}

func validateUserTags(tags map[string]string, propagatingTags bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(tags) == 0 {
		return allErrs
	}
	if len(tags) > 8 {
		allErrs = append(allErrs, field.Invalid(fldPath, len(tags), "number of user tags cannot be more than 8"))
	}
	for key, value := range tags {
		if strings.EqualFold(key, "Name") {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(key), tags[key], "Name key is not allowed for user defined tags"))
		}
		if propagatingTags {
			if err := validateTag(key, value); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath.Key(key), value, err.Error()))
			}
		} else {
			if strings.HasPrefix(key, "kubernetes.io/cluster/") {
				allErrs = append(allErrs, field.Invalid(fldPath.Key(key), tags[key], "Keys with prefix 'kubernetes.io/cluster/' are not allowed for user defined tags"))
			}
		}
	}
	return allErrs
}

// validateTag checks the following things to ensure that the tag is acceptable as an additional tag.
// * The key and value contain only valid characters.
// * The key is not empty and at most 128 characters.
// * The value is not empty and at most 256 characters. Note that, while many AWS services accept empty tag values,
//   the additional tags may be applied to resources in services that do not accept empty tag values. Consequently,
//   OpenShift cannot accept empty tag values.
// * The key is not in the kubernetes.io namespace.
// * The key is not in the openshift.io namespace.
func validateTag(key, value string) error {
	if !tagRegex.MatchString(key) {
		return fmt.Errorf("key contains invalid characters")
	}
	if !tagRegex.MatchString(value) {
		return fmt.Errorf("value contains invalid characters")
	}
	if len(key) == 0 {
		return fmt.Errorf("key is empty")
	}
	if len(key) > 128 {
		return fmt.Errorf("key is too long")
	}
	if len(value) == 0 {
		return fmt.Errorf("value is empty")
	}
	if len(value) > 256 {
		return fmt.Errorf("value is too long")
	}
	if kubernetesNamespaceRegex.MatchString(key) {
		return fmt.Errorf("key is in the kubernetes.io namespace")
	}
	if openshiftNamespaceRegex.MatchString(key) {
		return fmt.Errorf("key is in the openshift.io namespace")
	}
	return nil
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
