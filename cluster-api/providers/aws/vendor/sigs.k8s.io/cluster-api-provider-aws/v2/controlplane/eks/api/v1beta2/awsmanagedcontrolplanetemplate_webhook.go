/*
Copyright 2025 The Kubernetes Authors.

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

package v1beta2

import (
	"context"
	"fmt"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// templateLog is used for logging in this package.
var templateLog = ctrl.Log.WithName("awsmanagedcontrolplanetemplate-resource")

// SetupWebhookWithManager sets up the webhook with the Manager.
func (r *AWSManagedControlPlaneTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsManagedControlPlaneTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-controlplane-cluster-x-k8s-io-v1beta2-awsmanagedcontrolplanetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanetemplates,versions=v1beta2,name=validation.awsmanagedcontrolplanetemplates.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-controlplane-cluster-x-k8s-io-v1beta2-awsmanagedcontrolplanetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=awsmanagedcontrolplanetemplates,versions=v1beta2,name=default.awsmanagedcontrolplanetemplates.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsManagedControlPlaneTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &awsManagedControlPlaneTemplateWebhook{}
var _ webhook.CustomValidator = &awsManagedControlPlaneTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsManagedControlPlaneTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedControlPlaneTemplate object but got %T", r)
	}

	templateLog.Info("Validating AWSManagedControlPlaneTemplate create", "name", r.Name)

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateEKSVersion(nil)...)
	allErrs = append(allErrs, r.Spec.Template.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateRestrictPrivateSubnets()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.Template.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validateNetwork()...)
	allErrs = append(allErrs, r.validatePrivateDNSHostnameTypeOnLaunch()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsManagedControlPlaneTemplateWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSManagedControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedControlPlaneTemplate object but got %T", r)
	}

	templateLog.Info("AWSManagedControlPlaneTemplate validate update", "control-plane-template", klog.KObj(r))

	oldAWSManagedControlplaneTemplate, ok := oldObj.(*AWSManagedControlPlaneTemplate)
	if !ok {
		return nil, apierrors.NewInvalid(
			GroupVersion.WithKind("AWSManagedControlPlaneTemplate").GroupKind(),
			r.Name,
			field.ErrorList{field.InternalError(nil, errors.New("failed to convert old AWSManagedControlPlaneTemplate to object"))},
		)
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateEKSVersion(oldAWSManagedControlplaneTemplate)...)
	allErrs = append(allErrs, r.Spec.Template.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, r.validateIAMAuthConfig()...)
	allErrs = append(allErrs, r.validateSecondaryCIDR()...)
	allErrs = append(allErrs, r.validateEKSAddons()...)
	allErrs = append(allErrs, r.validateDisableVPCCNI()...)
	allErrs = append(allErrs, r.validateRestrictPrivateSubnets()...)
	allErrs = append(allErrs, r.validateKubeProxy()...)
	allErrs = append(allErrs, r.Spec.Template.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, r.validatePrivateDNSHostnameTypeOnLaunch()...)

	if r.Spec.Template.Spec.Region != oldAWSManagedControlplaneTemplate.Spec.Template.Spec.Region {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "region"), r.Spec.Template.Spec.Region, "field is immutable"),
		)
	}

	// If encryptionConfig is already set, do not allow removal of it.
	if oldAWSManagedControlplaneTemplate.Spec.Template.Spec.EncryptionConfig != nil && r.Spec.Template.Spec.EncryptionConfig == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "encryptionConfig"), r.Spec.Template.Spec.EncryptionConfig, "disabling EKS encryption is not allowed after it has been enabled"),
		)
	}

	// If encryptionConfig is already set, do not allow change in provider
	if r.Spec.Template.Spec.EncryptionConfig != nil &&
		r.Spec.Template.Spec.EncryptionConfig.Provider != nil &&
		oldAWSManagedControlplaneTemplate.Spec.Template.Spec.EncryptionConfig != nil &&
		oldAWSManagedControlplaneTemplate.Spec.Template.Spec.EncryptionConfig.Provider != nil &&
		*r.Spec.Template.Spec.EncryptionConfig.Provider != *oldAWSManagedControlplaneTemplate.Spec.Template.Spec.EncryptionConfig.Provider {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "encryptionConfig", "provider"), r.Spec.Template.Spec.EncryptionConfig.Provider, "changing EKS encryption is not allowed after it has been enabled"),
		)
	}

	// If a identityRef is already set, do not allow removal of it.
	if oldAWSManagedControlplaneTemplate.Spec.Template.Spec.IdentityRef != nil && r.Spec.Template.Spec.IdentityRef == nil {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "identityRef"),
				r.Spec.Template.Spec.IdentityRef, "field cannot be set to nil"),
		)
	}

	if oldAWSManagedControlplaneTemplate.Spec.Template.Spec.NetworkSpec.VPC.IsIPv6Enabled() != r.Spec.Template.Spec.NetworkSpec.VPC.IsIPv6Enabled() {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "network", "vpc", "enableIPv6"), r.Spec.Template.Spec.NetworkSpec.VPC.IsIPv6Enabled(), "changing IP family is not allowed after it has been set"))
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsManagedControlPlaneTemplateWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSManagedControlPlaneTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSManagedControlPlaneTemplate object but got %T", r)
	}

	templateLog.Info("Validating AWSManagedControlPlaneTemplate delete", "name", r.Name)
	// No validation logic on deletion.
	return nil, nil
}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*awsManagedControlPlaneTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSManagedControlPlaneTemplate)
	if !ok {
		return fmt.Errorf("expected an AWSManagedControlPlaneTemplate object but got %T", r)
	}

	templateLog.Info("AWSManagedControlPlaneTemplate setting defaults", "control-plane", klog.KObj(r))

	if r.Spec.Template.Spec.IdentityRef == nil {
		r.Spec.Template.Spec.IdentityRef = &infrav1.AWSIdentityReference{
			Kind: infrav1.ControllerIdentityKind,
			Name: infrav1.AWSClusterControllerIdentityName,
		}
	}

	infrav1.SetDefaults_Bastion(&r.Spec.Template.Spec.Bastion)
	infrav1.SetDefaults_NetworkSpec(&r.Spec.Template.Spec.NetworkSpec)

	return nil
}

func (r *AWSManagedControlPlaneTemplate) validateEKSVersion(old *AWSManagedControlPlaneTemplate) field.ErrorList {
	path := field.NewPath("spec.template.spec.version")
	var oldVersion *string
	if old != nil {
		oldVersion = old.Spec.Template.Spec.Version
	}
	return validateEKSVersion(r.Spec.Template.Spec.Version, oldVersion, r.Spec.Template.Spec.NetworkSpec, path)
}

func (r *AWSManagedControlPlaneTemplate) validateIAMAuthConfig() field.ErrorList {
	return validateIAMAuthConfig(r.Spec.Template.Spec.IAMAuthenticatorConfig, field.NewPath("spec.template.spec.iamAuthenticatorConfig"))
}

func (r *AWSManagedControlPlaneTemplate) validateSecondaryCIDR() field.ErrorList {
	return validateSecondaryCIDR(r.Spec.Template.Spec.SecondaryCidrBlock, field.NewPath("spec", "template", "spec", "secondaryCidrBlock"))
}

func (r *AWSManagedControlPlaneTemplate) validateEKSAddons() field.ErrorList {
	return validateEKSAddons(r.Spec.Template.Spec.Version, r.Spec.Template.Spec.NetworkSpec, r.Spec.Template.Spec.Addons, field.NewPath("spec.template.spec"))
}

func (r *AWSManagedControlPlaneTemplate) validateDisableVPCCNI() field.ErrorList {
	return validateDisableVPCCNI(r.Spec.Template.Spec.VpcCni, r.Spec.Template.Spec.Addons, field.NewPath("spec.template.spec"))
}

func (r *AWSManagedControlPlaneTemplate) validateRestrictPrivateSubnets() field.ErrorList {
	return validateRestrictPrivateSubnets(r.Spec.Template.Spec.RestrictPrivateSubnets, r.Spec.Template.Spec.NetworkSpec, "", field.NewPath("spec.template.spec"))
}

func (r *AWSManagedControlPlaneTemplate) validateKubeProxy() field.ErrorList {
	return validateKubeProxy(r.Spec.Template.Spec.KubeProxy, r.Spec.Template.Spec.Addons, field.NewPath("spec.template.spec"))
}

func (r *AWSManagedControlPlaneTemplate) validateNetwork() field.ErrorList {
	return validateNetwork("AWSManagedControlPlaneTemplate", r.Spec.Template.Spec.NetworkSpec, r.Spec.Template.Spec.SecondaryCidrBlock, field.NewPath("spec.template.spec"))
}

func (r *AWSManagedControlPlaneTemplate) validatePrivateDNSHostnameTypeOnLaunch() field.ErrorList {
	return validatePrivateDNSHostnameTypeOnLaunch(r.Spec.Template.Spec.NetworkSpec, field.NewPath("spec.template.spec"))
}
