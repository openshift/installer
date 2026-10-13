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
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var _ = ctrl.Log.WithName("awsclustercontrolleridentity-resource")

func (r *AWSClusterControllerIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsClusterControllerIdentityWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsclustercontrolleridentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustercontrolleridentities,versions=v1beta2,name=validation.awsclustercontrolleridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsclustercontrolleridentity,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustercontrolleridentities,versions=v1beta2,name=default.awsclustercontrolleridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsClusterControllerIdentityWebhook struct{}

var (
	_ webhook.CustomValidator = &awsClusterControllerIdentityWebhook{}
	_ webhook.CustomDefaulter = &awsClusterControllerIdentityWebhook{}
)

// ValidateCreate will do any extra validation when creating an AWSClusterControllerIdentity.
func (*awsClusterControllerIdentityWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSClusterControllerIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterControllerIdentity object but got %T", r)
	}

	// Ensures AWSClusterControllerIdentity being singleton by only allowing "default" as name
	if r.Name != AWSClusterControllerIdentityName {
		return nil, field.Invalid(field.NewPath("name"),
			r.Name, "AWSClusterControllerIdentity is a singleton and only acceptable name is default")
	}

	// Validate selector parses as Selector if AllowedNameSpaces is populated
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return nil, field.Invalid(field.NewPath("spec", "allowedNamespaces", "selector"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil, nil
}

// ValidateDelete allows you to add any extra validation when deleting an AWSClusterControllerIdentity.
func (*awsClusterControllerIdentityWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate will do any extra validation when updating an AWSClusterControllerIdentity.
func (*awsClusterControllerIdentityWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSClusterControllerIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterControllerIdentity object but got %T", r)
	}

	oldP, ok := oldObj.(*AWSClusterControllerIdentity)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an AWSClusterControllerIdentity but got a %T", oldObj))
	}

	if !cmp.Equal(r.Spec, oldP.Spec) {
		return nil, errors.New("AWSClusterControllerIdentity is immutable")
	}

	if r.Name != oldP.Name {
		return nil, field.Invalid(field.NewPath("name"),
			r.Name, "AWSClusterControllerIdentity is a singleton and only acceptable name is default")
	}

	// Validate selector parses as Selector if AllowedNameSpaces is not nil
	if r.Spec.AllowedNamespaces != nil {
		_, err := metav1.LabelSelectorAsSelector(&r.Spec.AllowedNamespaces.Selector)
		if err != nil {
			return nil, field.Invalid(field.NewPath("spec", "allowedNamespaces", "selectors"), r.Spec.AllowedNamespaces.Selector, err.Error())
		}
	}

	return nil, nil
}

// Default will set default values for the AWSClusterControllerIdentity.
func (*awsClusterControllerIdentityWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSClusterControllerIdentity)
	if !ok {
		return fmt.Errorf("expected an AWSClusterControllerIdentity object but got %T", r)
	}

	SetDefaults_Labels(&r.ObjectMeta)
	return nil
}
