/*
Copyright 2021 The Kubernetes Authors.

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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// log is for logging in this package.
var _ = logf.Log.WithName("awsclusterstaticidentity-resource")

func (r *AWSClusterStaticIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-awsclusterstaticidentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterstaticidentities,versions=v1beta1,name=validation.awsclusterstaticidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-awsclusterstaticidentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclusterstaticidentities,versions=v1beta1,name=default.awsclusterstaticidentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var (
	_ webhook.Validator = &AWSClusterStaticIdentity{}
	_ webhook.Defaulter = &AWSClusterStaticIdentity{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSClusterStaticIdentity) ValidateCreate() error {
	// Validate selector parses as Selector
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSClusterStaticIdentity) ValidateDelete() error {
	return nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSClusterStaticIdentity) ValidateUpdate(old runtime.Object) error {
	oldP, ok := old.(*AWSClusterStaticIdentity)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected an AWSClusterStaticIdentity but got a %T", old))
	}

	if oldP.Spec.SecretRef != r.Spec.SecretRef {
		return field.Invalid(field.NewPath("spec", "secretRef"),
			r.Spec.SecretRef, "field cannot be updated")
	}

	// Validate selector parses as Selector
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil
}

// Default should return the default AWSClusterStaticIdentity.
func (r *AWSClusterStaticIdentity) Default() {
	SetDefaults_Labels(&r.ObjectMeta)
}
