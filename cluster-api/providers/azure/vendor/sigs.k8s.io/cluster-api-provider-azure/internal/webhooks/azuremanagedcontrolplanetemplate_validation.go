/*
Copyright The Kubernetes Authors.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

    http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

package webhooks

import (
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/util/versions"
)

// Validate the Azure Managed Control Plane Template and return an aggregate error.
func validateAzureManagedControlPlaneTemplate(mcp *infrav1.AzureManagedControlPlaneTemplate, cli client.Client) error {
	var allErrs field.ErrorList

	allErrs = append(allErrs, validateVersion(
		mcp.Spec.Template.Spec.Version,
		field.NewPath("spec").Child("template").Child("spec").Child("version"))...)

	allErrs = append(allErrs, validateLoadBalancerProfile(
		mcp.Spec.Template.Spec.LoadBalancerProfile,
		field.NewPath("spec").Child("template").Child("spec").Child("loadBalancerProfile"))...)

	allErrs = append(allErrs, validateManagedClusterNetwork(
		cli,
		mcp.Labels,
		mcp.Namespace,
		mcp.Spec.Template.Spec.DNSServiceIP,
		mcp.Spec.Template.Spec.VirtualNetwork.Subnet,
		field.NewPath("spec").Child("template").Child("spec"))...)

	allErrs = append(allErrs, validateName(mcp.Name, field.NewPath("name"))...)

	allErrs = append(allErrs, validateAutoScalerProfile(mcp.Spec.Template.Spec.AutoScalerProfile, field.NewPath("spec").Child("template").Child("spec").Child("autoScalerProfile"))...)

	allErrs = append(allErrs, validateAKSExtensions(mcp.Spec.Template.Spec.Extensions, field.NewPath("spec").Child("extensions"))...)

	allErrs = append(allErrs, validateAzureManagedControlPlaneClassSpecSecurityProfile(&mcp.Spec.Template.Spec.AzureManagedControlPlaneClassSpec)...)

	allErrs = append(allErrs, validateNetworkPolicy(mcp.Spec.Template.Spec.NetworkPolicy, mcp.Spec.Template.Spec.NetworkDataplane, field.NewPath("spec").Child("template").Child("spec").Child("networkPolicy"))...)

	allErrs = append(allErrs, validateNetworkDataplane(mcp.Spec.Template.Spec.NetworkDataplane, mcp.Spec.Template.Spec.NetworkPolicy, mcp.Spec.Template.Spec.NetworkPluginMode, field.NewPath("spec").Child("template").Child("spec").Child("networkDataplane"))...)

	allErrs = append(allErrs, validateAPIServerAccessProfile(mcp.Spec.Template.Spec.APIServerAccessProfile, field.NewPath("spec").Child("template").Child("spec").Child("apiServerAccessProfile"))...)

	allErrs = append(allErrs, validateAMCPVirtualNetwork(mcp.Spec.Template.Spec.VirtualNetwork, field.NewPath("spec").Child("template").Child("spec").Child("virtualNetwork"))...)

	return allErrs.ToAggregate()
}

// validateK8sVersionUpdate validates K8s version.
func validateAzureManagedControlPlaneTemplateK8sVersionUpdate(mcp *infrav1.AzureManagedControlPlaneTemplate, old *infrav1.AzureManagedControlPlaneTemplate) field.ErrorList {
	var allErrs field.ErrorList
	if hv := versions.GetHigherK8sVersion(mcp.Spec.Template.Spec.Version, old.Spec.Template.Spec.Version); hv != mcp.Spec.Template.Spec.Version {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "version"),
			mcp.Spec.Template.Spec.Version, "field version cannot be downgraded"),
		)
	}
	return allErrs
}

// validateVirtualNetworkTemplateUpdate validates update to VirtualNetworkTemplate.
func validateAzureManagedControlPlaneTemplateVirtualNetworkTemplateUpdate(mcp *infrav1.AzureManagedControlPlaneTemplate, old *infrav1.AzureManagedControlPlaneTemplate) field.ErrorList {
	var allErrs field.ErrorList

	if old.Spec.Template.Spec.VirtualNetwork.CIDRBlock != mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "template", "spec", "virtualNetwork", "cidrBlock"),
				mcp.Spec.Template.Spec.VirtualNetwork.CIDRBlock,
				"Virtual Network CIDRBlock is immutable"))
	}

	if old.Spec.Template.Spec.VirtualNetwork.Subnet.Name != mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "template", "spec", "virtualNetwork", "subnet", "name"),
				mcp.Spec.Template.Spec.VirtualNetwork.Subnet.Name,
				"Subnet Name is immutable"))
	}

	// NOTE: This only works because we force the user to set the CIDRBlock for both the
	// managed and unmanaged Vnets. If we ever update the subnet cidr based on what's
	// actually set in the subnet, and it is different from what's in the Spec, for
	// unmanaged Vnets like we do with the AzureCluster this logic will break.
	if old.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock != mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "template", "spec", "virtualNetwork", "subnet", "cidrBlock"),
				mcp.Spec.Template.Spec.VirtualNetwork.Subnet.CIDRBlock,
				"Subnet CIDRBlock is immutable"))
	}

	if errs := validateAzureManagedControlPlaneClassSpecSecurityProfileUpdate(&mcp.Spec.Template.Spec.AzureManagedControlPlaneClassSpec, &old.Spec.Template.Spec.AzureManagedControlPlaneClassSpec); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	return allErrs
}

// validateAPIServerAccessProfileTemplateUpdate validates update to APIServerAccessProfileTemplate.
func validateAzureManagedControlPlaneTemplateAPIServerAccessProfileTemplateUpdate(mcp *infrav1.AzureManagedControlPlaneTemplate, old *infrav1.AzureManagedControlPlaneTemplate) field.ErrorList {
	var allErrs field.ErrorList

	newAPIServerAccessProfileNormalized := &infrav1.APIServerAccessProfile{}
	oldAPIServerAccessProfileNormalized := &infrav1.APIServerAccessProfile{}
	if mcp.Spec.Template.Spec.APIServerAccessProfile != nil {
		newAPIServerAccessProfileNormalized = &infrav1.APIServerAccessProfile{
			APIServerAccessProfileClassSpec: infrav1.APIServerAccessProfileClassSpec{
				EnablePrivateCluster:           mcp.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateCluster,
				PrivateDNSZone:                 mcp.Spec.Template.Spec.APIServerAccessProfile.PrivateDNSZone,
				EnablePrivateClusterPublicFQDN: mcp.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
			},
		}
	}
	if old.Spec.Template.Spec.APIServerAccessProfile != nil {
		oldAPIServerAccessProfileNormalized = &infrav1.APIServerAccessProfile{
			APIServerAccessProfileClassSpec: infrav1.APIServerAccessProfileClassSpec{
				EnablePrivateCluster:           old.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateCluster,
				PrivateDNSZone:                 old.Spec.Template.Spec.APIServerAccessProfile.PrivateDNSZone,
				EnablePrivateClusterPublicFQDN: old.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
			},
		}
	}

	if !reflect.DeepEqual(newAPIServerAccessProfileNormalized, oldAPIServerAccessProfileNormalized) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "template", "spec", "apiServerAccessProfile"),
				mcp.Spec.Template.Spec.APIServerAccessProfile, "fields are immutable"),
		)
	}

	return allErrs
}
