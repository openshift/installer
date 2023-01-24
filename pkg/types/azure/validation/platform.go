package validation

import (
	"fmt"
	"regexp"
	"sort"

	"k8s.io/apimachinery/pkg/util/validation/field"

	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

var (
	validCloudNames = map[azure.CloudEnvironment]bool{
		azure.PublicCloud:       true,
		azure.USGovernmentCloud: true,
		azure.ChinaCloud:        true,
		azure.GermanCloud:       true,
		azure.StackCloud:        true,
	}

	validCloudNameValues = func() []string {
		v := make([]string, 0, len(validCloudNames))
		for n := range validCloudNames {
			v = append(v, string(n))
		}
		return v
	}()
)

var (
	// tagKeyRegex is for verifying that the tag key contains only allowed characters.
	tagKeyRegex = regexp.MustCompile(`^[a-zA-Z][0-9A-Za-z_.=+\-@]{1,127}$`)

	// tagValueRegex is for verifying that the tag value contains only allowed characters.
	tagValueRegex = regexp.MustCompile(`^[0-9A-Za-z_.=+\-@]{1,256}$`)

	// tagKeyPrefixRegex is for verifying that the tag value does not contain restricted prefixes.
	tagKeyPrefixRegex = regexp.MustCompile(`^(?i)(name$|kubernetes\.io|openshift\.io|microsoft|azure|windows)`)
)

// maxUserTagLimit is the maximum userTags that can be configured as defined in openshift/api.
// https://github.com/openshift/api/blob/068483260288d83eac56053e202761b1702d46f5/config/v1/types_infrastructure.go#L482-L488
const maxUserTagLimit = 10

// ValidatePlatform checks that the specified platform is valid.
func ValidatePlatform(p *azure.Platform, publish types.PublishingStrategy, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if p.Region == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("region"), "region should be set to one of the supported Azure regions"))
	}
	if !p.IsARO() && publish == types.ExternalPublishingStrategy {
		if p.BaseDomainResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("baseDomainResourceGroupName"), "baseDomainResourceGroupName is the resource group name where the azure dns zone is deployed"))
		}
	}
	if p.DefaultMachinePlatform != nil {
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, "", p, fldPath.Child("defaultMachinePlatform"))...)
	}
	if p.VirtualNetwork != "" {
		if p.ComputeSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("computeSubnet"), "must provide a compute subnet when a virtual network is specified"))
		}
		if p.ControlPlaneSubnet == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("controlPlaneSubnet"), "must provide a control plane subnet when a virtual network is specified"))
		}
		if p.NetworkResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("networkResourceGroupName"), "must provide a network resource group when a virtual network is specified"))
		}
	}
	if (p.ComputeSubnet != "" || p.ControlPlaneSubnet != "") && (p.VirtualNetwork == "" || p.NetworkResourceGroupName == "") {
		if p.VirtualNetwork == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("virtualNetwork"), "must provide a virtual network when supplying subnets"))
		}
		if p.NetworkResourceGroupName == "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("networkResourceGroupName"), "must provide a network resource group when supplying subnets"))
		}
	}
	if !validCloudNames[p.CloudName] {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("cloudName"), p.CloudName, validCloudNameValues))
	}

	if _, ok := validOutboundTypes[p.OutboundType]; !ok {
		allErrs = append(allErrs, field.NotSupported(fldPath.Child("outboundType"), p.OutboundType, validOutboundTypeValues))
	}
	if p.OutboundType == azure.UserDefinedRoutingOutboundType && p.VirtualNetwork == "" {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("outboundType"), p.OutboundType, fmt.Sprintf("%s is only allowed when installing to pre-existing network", azure.UserDefinedRoutingOutboundType)))
	}

	// check if configured userTags are valid.
	allErrs = append(allErrs, validateUserTags(p.UserTags, fldPath.Child("userTags"))...)

	switch cloud := p.CloudName; cloud {
	case azure.StackCloud:
		allErrs = append(allErrs, validateAzureStack(p, fldPath)...)
	default:
		if p.ARMEndpoint != "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("armEndpoint"), fmt.Sprintf("ARM endpoint must not be set when the cloud name is %s", cloud)))
		}
		if p.ClusterOSImage != "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("clusterOSImage"), fmt.Sprintf("clusterOSImage must not be set when the cloud name is %s", cloud)))
		}
	}

	return allErrs
}

// validateUserTags verifies if configured number of UserTags is not more than
// allowed limit and the tag keys and values are valid.
func validateUserTags(tags map[string]string, fldPath *field.Path) field.ErrorList {
	allErrs := field.ErrorList{}
	if len(tags) == 0 {
		return allErrs
	}

	if len(tags) > maxUserTagLimit {
		allErrs = append(allErrs, field.TooMany(fldPath, len(tags), maxUserTagLimit))
	}

	for key, value := range tags {
		if err := validateTag(key, value); err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Key(key), value, err.Error()))
		}
	}
	return allErrs
}

// validateTag checks the following to ensure that the tag configured is acceptable.
//   - The key and value contain only allowed characters.
//   - The key is not empty and at most 128 characters and starts with an alphabet.
//   - The value is not empty and at most 256 characters.
//     Note: Although azure allows empty value, the tags may be applied to resources
//     in services that do not accept empty tag values. Consequently, OpenShift cannot
//     accept empty tag values.
//   - The key cannot be Name or have kubernetes.io, openshift.io, microsoft, azure,
//     windows prefixes.
func validateTag(key, value string) error {
	if !tagKeyRegex.MatchString(key) {
		return fmt.Errorf("key is invalid or contains invalid characters")
	}
	if !tagValueRegex.MatchString(value) {
		return fmt.Errorf("value is invalid or contains invalid characters")
	}
	if tagKeyPrefixRegex.MatchString(key) {
		return fmt.Errorf("key contains restricted prefix")
	}
	return nil
}

var (
	validOutboundTypes = map[azure.OutboundType]struct{}{
		azure.LoadbalancerOutboundType:       {},
		azure.UserDefinedRoutingOutboundType: {},
	}

	validOutboundTypeValues = func() []string {
		v := make([]string, 0, len(validOutboundTypes))
		for m := range validOutboundTypes {
			v = append(v, string(m))
		}
		sort.Strings(v)
		return v
	}()
)

func validateAzureStack(p *azure.Platform, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if p.ARMEndpoint == "" {
		allErrs = append(allErrs, field.Required(fldPath.Child("armEndpoint"), "ARM endpoint must be set when installing on Azure Stack"))
	}
	if p.OutboundType == azure.UserDefinedRoutingOutboundType {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("outboundType"), p.OutboundType, "Azure Stack does not support user-defined routing"))
	}
	return allErrs
}
