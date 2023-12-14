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
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api/util/topology"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// OpenStackMachineTemplateImmutableMsg ...
const OpenStackMachineTemplateImmutableMsg = "OpenStackMachineTemplate spec.template.spec field is immutable. Please create a new resource instead. Ref doc: https://cluster-api.sigs.k8s.io/tasks/change-machine-template.html"

// +kubebuilder:object:generate=false
type OpenStackMachineTemplateWebhook struct{}

func (r *OpenStackMachineTemplateWebhook) SetupWebhookWithManager(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(&OpenStackMachineTemplate{}).
		WithValidator(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha7-openstackmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackmachinetemplates,versions=v1alpha7,name=validation.openstackmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.CustomValidator = &OpenStackMachineTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *OpenStackMachineTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	openStackMachineTemplate, ok := obj.(*OpenStackMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an OpenStackMachineTemplate but got a %T", obj))
	}

	var allErrs field.ErrorList

	if openStackMachineTemplate.Spec.Template.Spec.ProviderID != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "providerID"), "cannot be set in templates"))
	}

	return aggregateObjErrors(openStackMachineTemplate.GroupVersionKind().GroupKind(), openStackMachineTemplate.Name, allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *OpenStackMachineTemplateWebhook) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldRaw.(*OpenStackMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an OpenStackMachineTemplate but got a %T", oldRaw))
	}

	newObj, ok := newRaw.(*OpenStackMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected an OpenStackMachineTemplate but got a %T", oldRaw))
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	if !topology.ShouldSkipImmutabilityChecks(req, newObj) &&
		!reflect.DeepEqual(newObj.Spec.Template.Spec, old.Spec.Template.Spec) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "template", "spec"), r, OpenStackMachineTemplateImmutableMsg),
		)
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *OpenStackMachineTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
