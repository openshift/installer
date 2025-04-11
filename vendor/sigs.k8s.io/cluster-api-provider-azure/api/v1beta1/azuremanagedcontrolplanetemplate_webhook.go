/*
Copyright 2023 The Kubernetes Authors.

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

package v1beta1

import (
	"context"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	"sigs.k8s.io/cluster-api-provider-azure/util/versions"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
	capifeature "sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupAzureManagedControlPlaneTemplateWebhookWithManager will set up the webhook to be managed by the specified manager.
func SetupAzureManagedControlPlaneTemplateWebhookWithManager(mgr ctrl.Manager) error {
	mcpw := &azureManagedControlPlaneTemplateWebhook{Client: mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AzureManagedControlPlaneTemplate{}).
		WithDefaulter(mcpw).
		WithValidator(mcpw).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplanetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanetemplates,versions=v1beta1,name=validation.azuremanagedcontrolplanetemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplanetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanetemplates,versions=v1beta1,name=default.azuremanagedcontrolplanetemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureManagedControlPlaneTemplateWebhook struct {
	Client client.Client
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mcpw *azureManagedControlPlaneTemplateWebhook) Default(ctx context.Context, obj runtime.Object) error {
	mcp, ok := obj.(*AzureManagedControlPlaneTemplate)
	if !ok {
		return apierrors.NewBadRequest("expected an AzureManagedControlPlaneTemplate")
	}
	mcp.setDefaults()
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *azureManagedControlPlaneTemplateWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	mcp, ok := obj.(*AzureManagedControlPlaneTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlaneTemplate")
	}
	// NOTE: AzureManagedControlPlaneTemplate relies upon MachinePools, which is behind a feature gate flag.
	// The webhook must prevent creating new objects in case the feature flag is disabled.
	if !feature.Gates.Enabled(capifeature.MachinePool) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			"can be set only if the Cluster API 'MachinePool' feature flag is enabled",
		)
	}

	return nil, mcp.validateManagedControlPlaneTemplate(mcpw.Client)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *azureManagedControlPlaneTemplateWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldObj.(*AzureManagedControlPlaneTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlaneTemplate")
	}
	mcp, ok := newObj.(*AzureManagedControlPlaneTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlaneTemplate")
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "subscriptionID"),
		old.Spec.Template.Spec.SubscriptionID,
		mcp.Spec.Template.Spec.SubscriptionID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "location"),
		old.Spec.Template.Spec.Location,
		mcp.Spec.Template.Spec.Location); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "dnsServiceIP"),
		old.Spec.Template.Spec.DNSServiceIP,
		mcp.Spec.Template.Spec.DNSServiceIP); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "networkPlugin"),
		old.Spec.Template.Spec.NetworkPlugin,
		mcp.Spec.Template.Spec.NetworkPlugin); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "networkPolicy"),
		old.Spec.Template.Spec.NetworkPolicy,
		mcp.Spec.Template.Spec.NetworkPolicy); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "networkDataplane"),
		old.Spec.Template.Spec.NetworkDataplane,
		mcp.Spec.Template.Spec.NetworkDataplane); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "loadBalancerSKU"),
		old.Spec.Template.Spec.LoadBalancerSKU,
		mcp.Spec.Template.Spec.LoadBalancerSKU); err != nil {
		allErrs = append(allErrs, err)
	}

	if old.Spec.Template.Spec.AADProfile != nil {
		if mcp.Spec.Template.Spec.AADProfile == nil {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "template", "spec", "aadProfile"),
					mcp.Spec.Template.Spec.AADProfile,
					"field cannot be nil, cannot disable AADProfile"))
		} else {
			if !mcp.Spec.Template.Spec.AADProfile.Managed && old.Spec.Template.Spec.AADProfile.Managed {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("spec", "template", "spec", "aadProfile", "managed"),
						mcp.Spec.Template.Spec.AADProfile.Managed,
						"cannot set AADProfile.Managed to false"))
			}
			if len(mcp.Spec.Template.Spec.AADProfile.AdminGroupObjectIDs) == 0 {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("spec", "template", "spec", "aadProfile", "adminGroupObjectIDs"),
						mcp.Spec.Template.Spec.AADProfile.AdminGroupObjectIDs,
						"length of AADProfile.AdminGroupObjectIDs cannot be zero"))
			}
		}
	}

	// Consider removing this once moves out of preview
	// Updating outboundType after cluster creation (PREVIEW)
	// https://learn.microsoft.com/en-us/azure/aks/egress-outboundtype#updating-outboundtype-after-cluster-creation-preview
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "outboundType"),
		old.Spec.Template.Spec.OutboundType,
		mcp.Spec.Template.Spec.OutboundType); err != nil {
		allErrs = append(allErrs, err)
	}

	if errs := mcp.validateVirtualNetworkTemplateUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := mcp.validateAPIServerAccessProfileTemplateUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAKSExtensionsUpdate(old.Spec.Template.Spec.Extensions, mcp.Spec.Template.Spec.Extensions); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := mcp.validateK8sVersionUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, mcp.validateManagedControlPlaneTemplate(mcpw.Client)
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureManagedControlPlaneTemplateKind).GroupKind(), mcp.Name, allErrs)
}

// Validate the Azure Managed Control Plane Template and return an aggregate error.
func (mcp *AzureManagedControlPlaneTemplate) validateManagedControlPlaneTemplate(cli client.Client) error {
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

	allErrs = append(allErrs, mcp.Spec.Template.Spec.AzureManagedControlPlaneClassSpec.validateSecurityProfile()...)

	allErrs = append(allErrs, validateNetworkPolicy(mcp.Spec.Template.Spec.NetworkPolicy, mcp.Spec.Template.Spec.NetworkDataplane, field.NewPath("spec").Child("template").Child("spec").Child("networkPolicy"))...)

	allErrs = append(allErrs, validateNetworkDataplane(mcp.Spec.Template.Spec.NetworkDataplane, mcp.Spec.Template.Spec.NetworkPolicy, mcp.Spec.Template.Spec.NetworkPluginMode, field.NewPath("spec").Child("template").Child("spec").Child("networkDataplane"))...)

	allErrs = append(allErrs, validateAPIServerAccessProfile(mcp.Spec.Template.Spec.APIServerAccessProfile, field.NewPath("spec").Child("template").Child("spec").Child("apiServerAccessProfile"))...)

	allErrs = append(allErrs, validateAMCPVirtualNetwork(mcp.Spec.Template.Spec.VirtualNetwork, field.NewPath("spec").Child("template").Child("spec").Child("virtualNetwork"))...)

	return allErrs.ToAggregate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *azureManagedControlPlaneTemplateWebhook) ValidateDelete(ctx context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// validateK8sVersionUpdate validates K8s version.
func (mcp *AzureManagedControlPlaneTemplate) validateK8sVersionUpdate(old *AzureManagedControlPlaneTemplate) field.ErrorList {
	var allErrs field.ErrorList
	if hv := versions.GetHigherK8sVersion(mcp.Spec.Template.Spec.Version, old.Spec.Template.Spec.Version); hv != mcp.Spec.Template.Spec.Version {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "template", "spec", "version"),
			mcp.Spec.Template.Spec.Version, "field version cannot be downgraded"),
		)
	}
	return allErrs
}

// validateVirtualNetworkTemplateUpdate validates update to VirtualNetworkTemplate.
func (mcp *AzureManagedControlPlaneTemplate) validateVirtualNetworkTemplateUpdate(old *AzureManagedControlPlaneTemplate) field.ErrorList {
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

	if errs := mcp.Spec.Template.Spec.AzureManagedControlPlaneClassSpec.validateSecurityProfileUpdate(&old.Spec.Template.Spec.AzureManagedControlPlaneClassSpec); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	return allErrs
}

// validateAPIServerAccessProfileTemplateUpdate validates update to APIServerAccessProfileTemplate.
func (mcp *AzureManagedControlPlaneTemplate) validateAPIServerAccessProfileTemplateUpdate(old *AzureManagedControlPlaneTemplate) field.ErrorList {
	var allErrs field.ErrorList

	newAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	oldAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	if mcp.Spec.Template.Spec.APIServerAccessProfile != nil {
		newAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			APIServerAccessProfileClassSpec: APIServerAccessProfileClassSpec{
				EnablePrivateCluster:           mcp.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateCluster,
				PrivateDNSZone:                 mcp.Spec.Template.Spec.APIServerAccessProfile.PrivateDNSZone,
				EnablePrivateClusterPublicFQDN: mcp.Spec.Template.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
			},
		}
	}
	if old.Spec.Template.Spec.APIServerAccessProfile != nil {
		oldAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			APIServerAccessProfileClassSpec: APIServerAccessProfileClassSpec{
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
