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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var _ = ctrl.Log.WithName("awsclusterroleidentity-resource")

func (r *AWSClusterRoleIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsClusterRoleIdentityWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterroleidentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterroleidentities,versions=v1beta2,name=validation.awsclusterroleidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterroleidentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterroleidentities,versions=v1beta2,name=default.awsclusterroleidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsClusterRoleIdentityWebhook struct{}

var (
	_ webhook.CustomValidator = &awsClusterRoleIdentityWebhook{}
	_ webhook.CustomDefaulter = &awsClusterRoleIdentityWebhook{}
)

// ValidateCreate will do any extra validation when creating an AWSClusterRoleIdentity.
func (*awsClusterRoleIdentityWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSClusterRoleIdentity)
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
func (*awsClusterRoleIdentityWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate will do any extra validation when updating an AWSClusterRoleIdentity.
func (*awsClusterRoleIdentityWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSClusterRoleIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterRoleIdentity object but got %T", r)
	}

	oldP, ok := oldObj.(*AWSClusterRoleIdentity)
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
func (*awsClusterRoleIdentityWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSClusterRoleIdentity)
	if !ok {
		return fmt.Errorf("expected an AWSClusterRoleIdentity object but got %T", r)
	}
	SetDefaults_Labels(&r.ObjectMeta)
	return nil
}
