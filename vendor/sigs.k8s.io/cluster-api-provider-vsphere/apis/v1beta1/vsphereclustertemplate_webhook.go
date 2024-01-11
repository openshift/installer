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
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *VSphereClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-vsphereclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vsphereclustertemplates,versions=v1beta1,name=validation.vsphereclustertemplate.infrastructure.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Validator = &VSphereClusterTemplate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereClusterTemplate) ValidateCreate() (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereClusterTemplate) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	old := oldRaw.(*VSphereClusterTemplate) //nolint:forcetypeassert
	if !reflect.DeepEqual(r.Spec.Template.Spec, old.Spec.Template.Spec) {
		return nil, field.Forbidden(field.NewPath("spec", "template", "spec"), "VSphereClusterTemplate spec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereClusterTemplate) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
