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
var _ = ctrl.Log.WithName("awsclusterstaticidentity-resource")

func (r *AWSClusterStaticIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsClusterStaticIdentityWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterstaticidentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterstaticidentities,versions=v1beta2,name=validation.awsclusterstaticidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsclusterstaticidentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterstaticidentities,versions=v1beta2,name=default.awsclusterstaticidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsClusterStaticIdentityWebhook struct{}

var (
	_ webhook.CustomValidator = &awsClusterStaticIdentityWebhook{}
	_ webhook.CustomDefaulter = &awsClusterStaticIdentityWebhook{}
)

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsClusterStaticIdentityWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSClusterStaticIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterStaticIdentity object but got %T", r)
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

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsClusterStaticIdentityWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*awsClusterStaticIdentityWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSClusterStaticIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterStaticIdentity object but got %T", r)
	}

	oldP, ok := oldObj.(*AWSClusterStaticIdentity)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an AWSClusterStaticIdentity but got a %T", oldObj))
	}

	if oldP.Spec.SecretRef != r.Spec.SecretRef {
		return nil, field.Invalid(field.NewPath("spec", "secretRef"),
			r.Spec.SecretRef, "field cannot be updated")
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

// Default should return the default AWSClusterStaticIdentity.
func (*awsClusterStaticIdentityWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSClusterStaticIdentity)
	if !ok {
		return fmt.Errorf("expected an AWSClusterStaticIdentity object but got %T", r)
	}
	SetDefaults_Labels(&r.ObjectMeta)
	return nil
}
