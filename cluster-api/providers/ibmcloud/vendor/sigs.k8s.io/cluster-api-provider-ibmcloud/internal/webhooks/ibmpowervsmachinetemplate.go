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

package webhooks

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,verbs=create;update,versions=v1beta2,name=mibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,versions=v1beta2,name=vibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMPowerVSMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1beta2.IBMPowerVSMachineTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSMachineTemplate implements a validation and defaulting webhook for IBMPowerVSMachineTemplate.
type IBMPowerVSMachineTemplate struct{}

var _ webhook.CustomDefaulter = &IBMPowerVSMachineTemplate{}
var _ webhook.CustomValidator = &IBMPowerVSMachineTemplate{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) Default(_ context.Context, obj runtime.Object) error {
	objValue, ok := obj.(*infrav1beta2.IBMPowerVSMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSMachineTemplate but got a %T", obj))
	}
	defaultIBMPowerVSMachineSpec(&objValue.Spec.Template.Spec)
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	objValue, ok := obj.(*infrav1beta2.IBMPowerVSMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSMachineTemplate but got a %T", obj))
	}
	return validateIBMPowerVSMachineTemplate(objValue)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (warnings admission.Warnings, err error) {
	objValue, ok := newObj.(*infrav1beta2.IBMPowerVSMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMPowerVSMachineTemplate but got a %T", newObj))
	}
	return validateIBMPowerVSMachineTemplate(objValue)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSMachineTemplate(machineTemplate *infrav1beta2.IBMPowerVSMachineTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := validateIBMPowerVSMachineTemplateNetwork(machineTemplate); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateIBMPowerVSMachineTemplateImage(machineTemplate); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateIBMPowerVSMachineTemplateMemory(machineTemplate); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := validateIBMPowerVSMachineTemplateProcessors(machineTemplate); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMPowerVSMachineTemplate"},
		machineTemplate.Name, allErrs)
}

func validateIBMPowerVSMachineTemplateNetwork(machineTemplate *infrav1beta2.IBMPowerVSMachineTemplate) *field.Error {
	if res, err := validateIBMPowerVSNetworkReference(machineTemplate.Spec.Template.Spec.Network); !res {
		return err
	}
	return nil
}

func validateIBMPowerVSMachineTemplateImage(machineTemplate *infrav1beta2.IBMPowerVSMachineTemplate) *field.Error {
	mt := machineTemplate.Spec.Template

	if mt.Spec.Image == nil && mt.Spec.ImageRef == nil {
		return field.Invalid(field.NewPath(""), "", "One of - Image or ImageRef must be specified")
	}

	if mt.Spec.Image != nil && mt.Spec.ImageRef != nil {
		return field.Invalid(field.NewPath(""), "", "Only one of - Image or ImageRef maybe be specified")
	}

	if mt.Spec.Image != nil {
		if res, err := validateIBMPowerVSResourceReference(*mt.Spec.Image, "Image"); !res {
			return err
		}
	}

	return nil
}

func validateIBMPowerVSMachineTemplateMemory(machineTemplate *infrav1beta2.IBMPowerVSMachineTemplate) *field.Error {
	if res := validateIBMPowerVSMemoryValues(machineTemplate.Spec.Template.Spec.MemoryGiB); !res {
		return field.Invalid(field.NewPath("spec", "template", "spec", "memoryGiB"), machineTemplate.Spec.Template.Spec.MemoryGiB, "Invalid Memory value - must be a positive integer no lesser than 2")
	}
	return nil
}

func validateIBMPowerVSMachineTemplateProcessors(machineTemplate *infrav1beta2.IBMPowerVSMachineTemplate) *field.Error {
	if res := validateIBMPowerVSProcessorValues(machineTemplate.Spec.Template.Spec.Processors); !res {
		return field.Invalid(field.NewPath("spec", "template", "spec", "processors"), machineTemplate.Spec.Template.Spec.Processors, "Invalid Processors value - must be non-empty and positive floating-point number no lesser than 0.25")
	}
	return nil
}
