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

package webhooks

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterroleidentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterroleidentities,versions=v1beta2,name=validation.awsclusterroleidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterroleidentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterroleidentities,versions=v1beta2,name=default.awsclusterroleidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AWSClusterRoleIdentity implements a validating and defaulting webhook for AWSClusterRoleIdentity.
type AWSClusterRoleIdentity struct{}

var (
	_ webhook.CustomValidator = &AWSClusterRoleIdentity{}
	_ webhook.CustomDefaulter = &AWSClusterRoleIdentity{}
)

func (w *AWSClusterRoleIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.AWSClusterRoleIdentity{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// ValidateCreate will do any extra validation when creating an AWSClusterRoleIdentity.
func (*AWSClusterRoleIdentity) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*infrav1.AWSClusterRoleIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterRoleIdentity object but got %T", r)
	}

	if r.Spec.SourceIdentityRef == nil {
		return nil, field.Invalid(field.NewPath("spec", "sourceIdentityRef"),
			r.Spec.SourceIdentityRef, "field cannot be set to nil")
	}

	// Validate selector parses as Selector
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return nil, field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil, nil
}

// ValidateDelete allows you to add any extra validation when deleting an AWSClusterRoleIdentity.
func (*AWSClusterRoleIdentity) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate will do any extra validation when updating an AWSClusterRoleIdentity.
func (*AWSClusterRoleIdentity) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*infrav1.AWSClusterRoleIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterRoleIdentity object but got %T", r)
	}

	oldP, ok := oldObj.(*infrav1.AWSClusterRoleIdentity)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an AWSClusterRoleIdentity but got a %T", oldObj))
	}

	// If a SourceIdentityRef is set, do not allow removal of it.
	if oldP.Spec.SourceIdentityRef != nil && r.Spec.SourceIdentityRef == nil {
		return nil, field.Invalid(field.NewPath("spec", "sourceIdentityRef"),
			r.Spec.SourceIdentityRef, "field cannot be set to nil")
	}

	// Validate selector parses as Selector
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return nil, field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil, nil
}

// Default will set default values for the AWSClusterRoleIdentity.
func (*AWSClusterRoleIdentity) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*infrav1.AWSClusterRoleIdentity)
	if !ok {
		return fmt.Errorf("expected an AWSClusterRoleIdentity object but got %T", r)
	}
	infrav1.SetDefaults_Labels(&r.ObjectMeta)
	return nil
}
