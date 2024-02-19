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

package v1alpha7

import (
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const openStackClusterTemplateImmutableMsg = "OpenStackClusterTemplate spec.template.spec field is immutable. Please create new resource instead."

func (r *OpenStackClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1alpha7-openstackclustertemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclustertemplates,versions=v1alpha7,name=default.openstackclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha7-openstackclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclustertemplates,versions=v1alpha7,name=validation.openstackclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var (
	_ webhook.Defaulter = &OpenStackClusterTemplate{}
	_ webhook.Validator = &OpenStackClusterTemplate{}
)

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *OpenStackClusterTemplate) Default() {
	if r.Spec.Template.Spec.IdentityRef != nil && r.Spec.Template.Spec.IdentityRef.Kind == "" {
		r.Spec.Template.Spec.IdentityRef.Kind = defaultIdentityRefKind
	}
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *OpenStackClusterTemplate) ValidateCreate() (admission.Warnings, error) {
	var allErrs field.ErrorList

	if r.Spec.Template.Spec.IdentityRef != nil && r.Spec.Template.Spec.IdentityRef.Kind != defaultIdentityRefKind {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "identityRef", "kind"), "must be a Secret"))
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *OpenStackClusterTemplate) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldRaw.(*OpenStackClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an OpenStackClusterTemplate but got a %T", oldRaw))
	}

	if !reflect.DeepEqual(r.Spec.Template.Spec, old.Spec.Template.Spec) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("OpenStackClusterTemplate", "spec", "template", "spec"), r, openStackClusterTemplateImmutableMsg),
		)
	}

	return aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *OpenStackClusterTemplate) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
