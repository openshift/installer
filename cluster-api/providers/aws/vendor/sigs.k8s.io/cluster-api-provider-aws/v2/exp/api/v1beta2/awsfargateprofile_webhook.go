/*
Copyright 2022 The Kubernetes Authors.

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

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	maxProfileNameLength = 100
	maxIAMRoleNameLength = 64
)

// SetupWebhookWithManager will setup the webhooks for the AWSFargateProfile.
func (r *AWSFargateProfile) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsFargateProfileWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsfargateprofile,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,versions=v1beta2,name=default.awsfargateprofile.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsfargateprofile,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsfargateprofiles,versions=v1beta2,name=validation.awsfargateprofile.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsFargateProfileWebhook struct{}

var _ webhook.CustomDefaulter = &awsFargateProfileWebhook{}
var _ webhook.CustomValidator = &awsFargateProfileWebhook{}

// Default will set default values for the AWSFargateProfile.
func (*awsFargateProfileWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSFargateProfile)
	if !ok {
		return fmt.Errorf("expected an AWSFargateProfile object but got %T", r)
	}

	if r.Labels == nil {
		r.Labels = make(map[string]string)
	}
	r.Labels[clusterv1.ClusterNameLabel] = r.Spec.ClusterName

	if r.Spec.ProfileName == "" {
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxProfileNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
			return nil
		}

		r.Spec.ProfileName = name
	}
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*awsFargateProfileWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSFargateProfile)
	if !ok {
		return nil, fmt.Errorf("expected an AWSFargateProfile object but got %T", r)
	}

	gv := r.GroupVersionKind().GroupKind()
	old, ok := oldObj.(*AWSFargateProfile)
	if !ok {
		return nil, apierrors.NewInvalid(gv, r.Name, field.ErrorList{
			field.InternalError(nil, errors.Errorf("failed to convert old %s to object", gv.Kind)),
		})
	}

	var allErrs field.ErrorList

	// Spec is immutable, but if the new RoleName is the generated one(or default if EnableIAM is disabled) and
	// the old RoleName is nil, then ignore checking that field.
	if old.Spec.RoleName == "" {
		roleName, err := eks.GenerateEKSName(
			"fargate",
			fmt.Sprintf("%s-%s", r.Spec.ClusterName, r.Spec.ProfileName),
			maxIAMRoleNameLength,
		)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS fargate role name")
		}

		if r.Spec.RoleName == roleName || r.Spec.RoleName == DefaultEKSFargateRole {
			r.Spec.RoleName = ""
		}
	}

	// Spec is immutable, but if the new ProfileName is the defaulted one and
	// the old ProfileName is nil, then ignore checking that field.
	if old.Spec.ProfileName == "" {
		name, err := eks.GenerateEKSName(old.Name, old.Namespace, maxProfileNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
		}
		if r.Spec.ProfileName == name {
			r.Spec.ProfileName = ""
		}
	}

	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	// remove additionalTags from equal check since they are mutable
	old.Spec.AdditionalTags = nil
	r.Spec.AdditionalTags = nil

	if !cmp.Equal(old.Spec, r.Spec) {
		allErrs = append(
			allErrs,
			field.Invalid(field.NewPath("spec"), r.Spec, "is immutable"),
		)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		gv,
		r.Name,
		allErrs,
	)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*awsFargateProfileWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSFargateProfile)
	if !ok {
		return nil, fmt.Errorf("expected an AWSFargateProfile object but got %T", r)
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)
	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*awsFargateProfileWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
