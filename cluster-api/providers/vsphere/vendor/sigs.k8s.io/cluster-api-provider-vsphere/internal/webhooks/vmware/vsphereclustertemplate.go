/*
Copyright 2026 The Kubernetes Authors.

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

// Package vmware is the package for webhooks of vmware resources.
package vmware

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/api/supervisor/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/vmware/conversion"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta2-vsphereclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclustertemplates,versions=v1beta2,name=validation.vsphereclustertemplate.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereClusterTemplate implements a validation webhook for VSphereClusterTemplate.
type VSphereClusterTemplate struct {
	NetworkProvider string
}

var _ admission.Validator[*vmwarev1.VSphereClusterTemplate] = &VSphereClusterTemplate{}

// SetupWebhookWithManager sets up the webhook with the Manager.
func (webhook *VSphereClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &vmwarev1.VSphereClusterTemplate{}).
		WithValidator(webhook).
		WithConverter(conversion.VSphereClusterTemplate).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateCreate(_ context.Context, obj *vmwarev1.VSphereClusterTemplate) (admission.Warnings, error) {
	return webhook.validate(obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateUpdate(_ context.Context, _, newObj *vmwarev1.VSphereClusterTemplate) (admission.Warnings, error) {
	return webhook.validate(newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateDelete(_ context.Context, _ *vmwarev1.VSphereClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}

// validateClusterTemplateNetwork validates the network configuration of the VSphereClusterTemplate.
func (webhook *VSphereClusterTemplate) validateClusterTemplateNetwork(template *vmwarev1.VSphereClusterTemplate) field.ErrorList {
	var allErrs field.ErrorList

	if !feature.Gates.Enabled(feature.MultiNetworks) && template.Spec.Template.Spec.Network.NSXVPC.CreateSubnetSet != nil {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "template", "spec", "network", "nsxVPC", "createSubnetSet"),
			"createSubnetSet can only be set when MultiNetworks feature gate is enabled",
		))
	}
	if template.Spec.Template.Spec.Network.NSXVPC.IsDefined() && webhook.NetworkProvider != manager.NSXVPCNetworkProvider {
		allErrs = append(allErrs, field.Forbidden(
			field.NewPath("spec", "template", "spec", "network", "nsxVPC"),
			"nsxVPC can only be set when network provider is NSX-VPC",
		))
	}

	return allErrs
}

// validate aggregates all validations for the VSphereClusterTemplate.
func (webhook *VSphereClusterTemplate) validate(template *vmwarev1.VSphereClusterTemplate) (admission.Warnings, error) {
	allErrs := webhook.validateClusterTemplateNetwork(template)
	allErrs = append(allErrs, validateFailureDomainsControlPlaneSelector(
		template.Spec.Template.Spec.FailureDomains.ControlPlane.Selector,
		field.NewPath("spec", "template", "spec", "failureDomains", "controlPlane", "selector"),
	)...)

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(template.GroupVersionKind().GroupKind(), template.Name, allErrs)
	}

	return nil, nil
}
