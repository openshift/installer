package validation

import (
	"fmt"
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
		allErrs = append(allErrs, ValidateMachinePool(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
		allErrs = append(allErrs, ValidateDefaultDiskType(p.DefaultMachinePlatform, fldPath.Child("defaultMachinePlatform"))...)
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

	switch cloud := p.CloudName; cloud {
	case azure.StackCloud:
		allErrs = append(allErrs, validateAzureStack(p, fldPath)...)
	default:
		if p.ARMEndpoint != "" {
			allErrs = append(allErrs, field.Required(fldPath.Child("armEndpoint"), fmt.Sprintf("ARM endpoint must not be set when the cloud name is %s", cloud)))
		}
	}

	return allErrs
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
