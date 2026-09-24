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
	"context"

	"github.com/go-logr/logr"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	apiinternal "sigs.k8s.io/cluster-api-provider-azure/internal/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager will set up the webhook to be managed by the specified manager.
func (mcpw *AzureManagedControlPlaneTemplateWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mcpw.client = mgr.GetClient()
	mcpw.logger = mgr.GetLogger().WithName("AzureManagedControlPlaneTemplate")

	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureManagedControlPlaneTemplate{}).
		WithDefaulter(mcpw).
		WithValidator(mcpw).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplanetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanetemplates,versions=v1beta1,name=validation.azuremanagedcontrolplanetemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplanetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanetemplates,versions=v1beta1,name=default.azuremanagedcontrolplanetemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureManagedControlPlaneTemplateWebhook implements a validating and defaulting webhook for AzureManagedControlPlaneTemplate.
type AzureManagedControlPlaneTemplateWebhook struct {
	client client.Client
	logger logr.Logger
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mcpw *AzureManagedControlPlaneTemplateWebhook) Default(_ context.Context, mcp *infrav1.AzureManagedControlPlaneTemplate) error {
	apiinternal.SetDefaultsAzureManagedControlPlaneTemplate(mcpw.logger, mcp)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *AzureManagedControlPlaneTemplateWebhook) ValidateCreate(_ context.Context, mcp *infrav1.AzureManagedControlPlaneTemplate) (admission.Warnings, error) {
	return nil, validateAzureManagedControlPlaneTemplate(mcp, mcpw.client)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *AzureManagedControlPlaneTemplateWebhook) ValidateUpdate(_ context.Context, old, mcp *infrav1.AzureManagedControlPlaneTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList
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

	if errs := validateAzureManagedControlPlaneTemplateVirtualNetworkTemplateUpdate(mcp, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAzureManagedControlPlaneTemplateAPIServerAccessProfileTemplateUpdate(mcp, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAKSExtensionsUpdate(old.Spec.Template.Spec.Extensions, mcp.Spec.Template.Spec.Extensions); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := validateAzureManagedControlPlaneTemplateK8sVersionUpdate(mcp, old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, validateAzureManagedControlPlaneTemplate(mcp, mcpw.client)
	}

	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureManagedControlPlaneTemplateKind).GroupKind(), mcp.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mcpw *AzureManagedControlPlaneTemplateWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureManagedControlPlaneTemplate) (admission.Warnings, error) {
	return nil, nil
}
