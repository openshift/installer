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
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var ibmpowervsmachinetemplatelog = logf.Log.WithName("ibmpowervsmachinetemplate-resource")

func (r *IBMPowerVSMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,verbs=create;update,versions=v1beta2,name=mibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &IBMPowerVSMachineTemplate{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) Default() {
	ibmpowervsmachinetemplatelog.Info("default", "name", r.Name)
	defaultIBMPowerVSMachineSpec(&r.Spec.Template.Spec)
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,versions=v1beta2,name=vibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Validator = &IBMPowerVSMachineTemplate{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateCreate() (admission.Warnings, error) {
	ibmpowervsmachinetemplatelog.Info("validate create", "name", r.Name)
	return r.validateIBMPowerVSMachineTemplate()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	ibmpowervsmachinetemplatelog.Info("validate update", "name", r.Name)
	return r.validateIBMPowerVSMachineTemplate()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateDelete() (admission.Warnings, error) {
	ibmpowervsmachinetemplatelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

func (r *IBMPowerVSMachineTemplate) validateIBMPowerVSMachineTemplate() (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := r.validateIBMPowerVSMachineTemplateNetwork(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateIBMPowerVSMachineTemplateImage(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateIBMPowerVSMachineTemplateMemory(); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := r.validateIBMPowerVSMachineTemplateProcessors(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMPowerVSMachineTemplate"},
		r.Name, allErrs)
}

func (r *IBMPowerVSMachineTemplate) validateIBMPowerVSMachineTemplateNetwork() *field.Error {
	if res, err := validateIBMPowerVSNetworkReference(r.Spec.Template.Spec.Network); !res {
		return err
	}
	return nil
}

func (r *IBMPowerVSMachineTemplate) validateIBMPowerVSMachineTemplateImage() *field.Error {
	mt := r.Spec.Template

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

func (r *IBMPowerVSMachineTemplate) validateIBMPowerVSMachineTemplateMemory() *field.Error {
	if res := validateIBMPowerVSMemoryValues(r.Spec.Template.Spec.MemoryGiB); !res {
		return field.Invalid(field.NewPath("spec", "template", "spec", "memoryGiB"), r.Spec.Template.Spec.MemoryGiB, "Invalid Memory value - must be a positive integer no lesser than 2")
	}
	return nil
}

func (r *IBMPowerVSMachineTemplate) validateIBMPowerVSMachineTemplateProcessors() *field.Error {
	if res := validateIBMPowerVSProcessorValues(r.Spec.Template.Spec.Processors); !res {
		return field.Invalid(field.NewPath("spec", "template", "spec", "processors"), r.Spec.Template.Spec.Processors, "Invalid Processors value - must be non-empty and positive floating-point number no lesser than 0.25")
	}
	return nil
}
