package validation

import (
	"fmt"
	configv1 "github.com/openshift/api/config/v1"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
)

// tagRegex is used to check that the keys and values of a tag contain only valid characters.
var tagRegex = regexp.MustCompile(`^[0-9A-Za-z_.:/=+-@\p{Z}]*$`)

// kubernetesNamespaceRegex is used to check that a tag key is not in the kubernetes.io namespace.
var kubernetesNamespaceRegex = regexp.MustCompile(`^([^/]*\.)?kubernetes.io/`)

// openshiftNamespaceRegex is used to check that a tag key is not in the openshift.io namespace.
var openshiftNamespaceRegex = regexp.MustCompile(`^([^/]*\.)?openshift.io/`)

// userTagLimit is defined in openshift/api
// https://github.com/openshift/api/blob/1265e99256880f8679d1b74561c0bc7932067c43/config/v1/types_infrastructure.go#L370-L376
const userTagLimit = 25

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *aws.Platform, cm types.CredentialsMode, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	if p.HostedZone != "" {
		if len(p.Subnets) == 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostedZone"), p.HostedZone, "may not use an existing hosted zone when not using existing subnets"))
		}
	}

	if p.HostedZoneRole != "" {
		if p.HostedZone == "" {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("hostedZoneRole"), p.HostedZoneRole, "may not specify a role to assume for hosted zone operations without also specifying a hosted zone"))
		}

		if cm != types.ManualCredentialsMode && cm != types.PassthroughCredentialsMode {
			errMsg := "when specifying a hostedZoneRole, either Passthrough or Manual credential mode must be specified"
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("credentialsMode"), errMsg))
		}
	}

	allErrs = append(allErrs, validateServiceEndpoints(p.ServiceEndpoints, fldPath.Child("serviceEndpoints"))...)
	allErrs = append(allErrs, validateUserTags(p.UserTags, p.PropagateUserTag, fldPath.Child("userTags"))...)

	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p, p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
	}

	allErrs = append(allErrs, checkForEIPAllocationsWithNLB(p, fldPath.Child("aws", "lbType"))...)

	if p != nil && p.EIPAllocations != nil && len(p.EIPAllocations.IngressNetworkLoadBalancer) > 0 {
		var allocationIDs []string
		for _, allocationID := range p.EIPAllocations.IngressNetworkLoadBalancer {
			allocationIDs = append(allocationIDs, string(allocationID))
		}

		// Validate EIP allocation IDs
		allErrs = append(allErrs, validateEIPAllocationsList(allocationIDs, fldPath.Child("aws", "eipAllocations", "ingressNetworkLoadBalancer"))...)
	}
	return allErrs
}

// validateEIPAllocations validates an array of EIP allocation IDs based on given rules.
func validateEIPAllocationsList(allocationIDs []string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	// Rule 1: Check for duplicates
	seen := make(map[string]bool)
	for _, id := range allocationIDs {
		if seen[id] {
			allErrs = append(allErrs, field.Invalid(fldPath, id, "cannot have duplicate EIP Allocation IDs"))
		}
		seen[id] = true
	}

	// Rule 2: Check the max length of the array
	if len(allocationIDs) > 10 {
		allErrs = append(allErrs, field.TooMany(fldPath, len(allocationIDs), 10))
	}

	// Regular expression to validate EIP allocation format
	eipRegex := regexp.MustCompile(`^eipalloc-[a-fA-F0-9]{17}$`)

	// Validate each allocation ID
	for _, id := range allocationIDs {
		if len(id) != 26 {
			// Rule 3: Check for minimum and maximum length
			allErrs = append(allErrs, field.Invalid(fldPath, id, "invalid EIP allocation ID length"))
		} else if len(id) < 9 || id[:9] != "eipalloc-" {
			// Rule 4: Value shall start with "eipalloc-"
			allErrs = append(allErrs, field.Invalid(fldPath, id, "eipAllocations should start with 'eipalloc-'"))
		} else if !eipRegex.MatchString(id) {
			// Rule 5: Check if the value matches the regex
			allErrs = append(allErrs, field.Invalid(fldPath, id, "eipAllocations must be 'eipalloc-' followed by exactly 17 hexadecimal characters (0-9, a-f, A-F)"))
		}
	}
	return allErrs
}

// Checks if the type was set to NLB when eipAllocations are provided to the install-config.
func checkForEIPAllocationsWithNLB(aws *aws.Platform, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if aws.EIPAllocations == nil || len(aws.EIPAllocations.IngressNetworkLoadBalancer) == 0 {
		// EIPAllocations not provided
		return nil
	}
	if len(aws.LBType) == 0 || (len(aws.LBType) > 0 && aws.LBType != configv1.NLB) {
		allErrs = append(allErrs, field.Required(fldPath, "lbType NLB must be specified"))
	}
	return allErrs
}

func validateUserTags(tags map[string]string, propagatingTags bool, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(tags) == 0 {
		return allErrs
	}
	if len(tags) > 8 {
		logrus.Warnf("Due to a limit of 10 tags on S3 Bucket Objects, only the first eight lexicographically sorted tags will be applied to the bootstrap ignition object, which is a temporary resource only used during installation")
	}

	if len(tags) > userTagLimit {
		allErrs = append(allErrs, field.TooMany(fldPath, len(tags), userTagLimit))
	}
	for key, value := range tags {
		if err := validateTag(key, value); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(key), value, err.Error()))
		}
	}
	return allErrs
}

// validateTag checks the following things to ensure that the tag is acceptable as an additional tag.
//   - The key and value contain only valid characters.
//   - The key is not empty and at most 128 characters.
//   - The value is not empty and at most 256 characters. Note that, while many AWS services accept empty tag values,
//     the additional tags may be applied to resources in services that do not accept empty tag values. Consequently,
//     OpenShift cannot accept empty tag values.
//   - The key is not in the kubernetes.io namespace.
//   - The key is not in the openshift.io namespace.
func validateTag(key, value string) error {
	if strings.EqualFold(key, "Name") {
		return fmt.Errorf("\"Name\" key is not allowed for user defined tags")
	}
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
