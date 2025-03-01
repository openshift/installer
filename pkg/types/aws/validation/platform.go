package validation

import (
	"fmt"
	"net/url"
	"regexp"
	"strings"

	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/sets"
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
func ValidatePlatform(p *aws.Platform, publish types.PublishingStrategy, cm types.CredentialsMode, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}

	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region must be specified"))
	}

	if p.HostedZone != "" {
		if len(p.VPC.Subnets) == 0 {
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

	allErrs = append(allErrs, validateSubnets(p.VPC.Subnets, publish, fldPath.Child("vpc", "subnets"))...)
	allErrs = append(allErrs, validateServiceEndpoints(p.ServiceEndpoints, fldPath.Child("serviceEndpoints"))...)
	allErrs = append(allErrs, validateUserTags(p.UserTags, p.PropagateUserTag, fldPath.Child("userTags"))...)

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

func validateSubnets(subnets []aws.Subnet, publish types.PublishingStrategy, fldPath *field.Path) field.ErrorList {
	if len(subnets) == 0 {
		return nil
	}

	allErrs := field.ErrorList{}

	// Either all subnets must be assigned roles (manual role selection)
	// or none of the subnets should have roles assigned (automatic role selection).
	for _, subnet := range subnets {
		if (len(subnet.Roles) > 0) != (len(subnets[0].Roles) > 0) {
			allErrs = append(allErrs, field.Forbidden(fldPath, "either all subnets must be assigned roles or none of the subnets should have roles assigned"))
			break
		}
	}

	supportedRoles := sets.New(
		aws.ClusterNodeSubnetRole,
		aws.EdgeNodeSubnetRole,
		aws.BootstrapNodeSubnetRole,
		aws.IngressControllerLBSubnetRole,
		aws.ControlPlaneExternalLBSubnetRole,
		aws.ControlPlaneInternalLBSubnetRole,
	)
	// A mapping of a role to its assigned subnets.
	subnetsForRole := make(map[aws.SubnetRoleType][]aws.Subnet)
	// An indicator whether none of the subnets are assigned roles.
	autoRoleSelection := true
	// A tracker to check duplicate subnet IDs.
	subnetTracker := make(map[string]int)

	for snIdx, subnet := range subnets {
		subnetFldPath := fldPath.Index(snIdx)

		// Subnet ID must be unique.
		if _, ok := subnetTracker[string(subnet.ID)]; ok {
			allErrs = append(allErrs, field.Duplicate(subnetFldPath.Child("id"), subnet.ID))
		} else {
			subnetTracker[string(subnet.ID)] = snIdx
		}

		// A set of role types assigned to a subnet
		// for quick search later.
		subnetRoleTypes := sets.New[aws.SubnetRoleType]()
		// A tracker to check duplicate subnet roles.
		roleTracker := make(map[string]int)

		for rIdx, role := range subnet.Roles {
			if !supportedRoles.Has(role.Type) {
				allErrs = append(allErrs, field.NotSupported(subnetFldPath.Child("roles").Index(rIdx).Child("type"), role.Type, sets.List(supportedRoles)))
				continue // Role type is unsupported. No further validation.
			}

			// Subnet roles must not contain duplicates.
			if _, ok := roleTracker[string(role.Type)]; ok {
				allErrs = append(allErrs, field.Duplicate(subnetFldPath.Child("roles").Index(rIdx).Child("type"), role.Type))
			} else {
				roleTracker[string(role.Type)] = rIdx
			}

			subnetRoleTypes.Insert(role.Type)
		}

		// Role ControlPlaneExternalLB and ControlPlaneInternalLB must not be both specified.
		// on the same subnet.
		if subnetRoleTypes.HasAll(aws.ControlPlaneExternalLBSubnetRole, aws.ControlPlaneInternalLBSubnetRole) {
			allErrs = append(allErrs, field.Forbidden(subnetFldPath.Child("roles"), "must not have both ControlPlaneExternalLB and ControlPlaneInternalLB role"))
		}

		// EdgeNode cannot be combined with any other roles.
		if subnetRoleTypes.Has(aws.EdgeNodeSubnetRole) && subnetRoleTypes.Len() > 1 {
			allErrs = append(allErrs, field.Forbidden(subnetFldPath.Child("roles"), "must not combine EdgeNode role with any other roles"))
		}

		autoRoleSelection = autoRoleSelection && len(subnetRoleTypes) == 0
		for rType := range subnetRoleTypes {
			subnetsForRole[rType] = append(subnetsForRole[rType], subnet)
		}
	}

	if !autoRoleSelection {
		// The IngressController's API only allows 10 subnets.
		if len(subnetsForRole[aws.IngressControllerLBSubnetRole]) > 10 {
			allErrs = append(allErrs, field.Forbidden(fldPath, "must not include more than 10 subnets with the IngressControllerLB role"))
		}

		// If the cluster is private, ControlPlaneExternalLB role is not allowed
		// as only an internal control plane load balancer will be created.
		if publish == types.InternalPublishingStrategy && len(subnetsForRole[aws.ControlPlaneExternalLBSubnetRole]) > 0 {
			allErrs = append(allErrs, field.Forbidden(fldPath, "must not include subnets with the ControlPlaneExternalLBSubnetRole role in a private cluster"))
		}

		// ClusterNode, IngressControllerLB, ControlPlaneExternalLB, and ControlPlaneInternalLB
		// must be assigned to at least 1 subnet.
		missingRoles := sets.New[aws.SubnetRoleType]()
		for rType := range supportedRoles {
			switch rType {
			// EdgeNode role is optional.
			case aws.EdgeNodeSubnetRole:
			// If the cluster is internal, ControlPlaneExternalLB role is not required.
			case aws.ControlPlaneExternalLBSubnetRole:
				if publish == types.InternalPublishingStrategy {
					continue
				}
				fallthrough
			default:
				if len(subnetsForRole[rType]) == 0 {
					missingRoles.Insert(rType)
				}
			}
		}

		if missingRoles.Len() > 0 {
			allErrs = append(allErrs, field.Invalid(fldPath, subnets, fmt.Sprintf("roles %v must be assigned to at least 1 subnet", sets.List(missingRoles))))
		}
	}

	return allErrs
}
