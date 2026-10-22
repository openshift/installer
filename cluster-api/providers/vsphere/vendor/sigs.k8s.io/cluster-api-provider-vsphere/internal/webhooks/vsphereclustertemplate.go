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

package webhooks

import (
	"context"
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/conversion"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-vsphereclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vsphereclustertemplates,versions=v1beta2,name=validation.vsphereclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereClusterTemplate implements a validation webhook for VSphereClusterTemplate.
type VSphereClusterTemplate struct{}

var _ admission.Validator[*infrav1.VSphereClusterTemplate] = &VSphereClusterTemplate{}

func (webhook *VSphereClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.VSphereClusterTemplate{}).
		WithValidator(webhook).
		WithConverter(conversion.VSphereClusterTemplate).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateCreate(_ context.Context, _ *infrav1.VSphereClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateUpdate(_ context.Context, oldTyped, newTyped *infrav1.VSphereClusterTemplate) (admission.Warnings, error) {
	if !reflect.DeepEqual(newTyped.Spec.Template.Spec, oldTyped.Spec.Template.Spec) {
		return nil, field.Forbidden(field.NewPath("spec", "template", "spec"), "VSphereClusterTemplate spec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereClusterTemplate) ValidateDelete(_ context.Context, _ *infrav1.VSphereClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}
