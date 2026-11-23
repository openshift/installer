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

package powervs

import (
	"context"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

// Ensure IBMPowerVSMachineTemplate implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMPowerVSMachineTemplate] = &IBMPowerVSMachineTemplate{}
	_ admission.Defaulter[*infrav1.IBMPowerVSMachineTemplate] = &IBMPowerVSMachineTemplate{}
)

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,versions=v1beta3,name=mibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachinetemplates,versions=v1beta3,name=vibmpowervsmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *IBMPowerVSMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMPowerVSMachineTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSMachineTemplate implements a validation and defaulting webhook for IBMPowerVSMachineTemplate.
type IBMPowerVSMachineTemplate struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) Default(_ context.Context, obj *infrav1.IBMPowerVSMachineTemplate) error {
	defaultIBMPowerVSMachineSpec(&obj.Spec.Template.Spec)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateCreate(_ context.Context, obj *infrav1.IBMPowerVSMachineTemplate) (admission.Warnings, error) {
	return validateIBMPowerVSMachineTemplate(obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateUpdate(_ context.Context, _, newObj *infrav1.IBMPowerVSMachineTemplate) (warnings admission.Warnings, err error) {
	return validateIBMPowerVSMachineTemplate(newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachineTemplate) ValidateDelete(_ context.Context, _ *infrav1.IBMPowerVSMachineTemplate) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSMachineTemplate(machineTemplate *infrav1.IBMPowerVSMachineTemplate) (admission.Warnings, error) {
	// Network: CRD CEL on ResourceIdentifier enforces exactly-one of ID/Name.
	// Image:   CRD CEL on IBMPowerVSMachineImage enforces type↔reference/import mutual exclusion;
	//          Enum marker rejects invalid type; MinLength=1 on ImageReference.Name rejects empty import name.
	// Memory:  +kubebuilder:validation:Minimum=2 on MemoryGiB enforces the minimum at the CRD level.
	// Processors: intstr.IntOrString has no CRD-expressible minimum, so the webhook is the right layer.
	if err := validateIBMPowerVSMachineTemplateProcessors(machineTemplate); err != nil {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: infrastructureGroup, Kind: "IBMPowerVSMachineTemplate"},
			machineTemplate.Name, field.ErrorList{err})
	}
	return nil, nil
}

func validateIBMPowerVSMachineTemplateProcessors(machineTemplate *infrav1.IBMPowerVSMachineTemplate) *field.Error {
	spec := machineTemplate.Spec.Template.Spec
	err := validateIBMPowerVSProcessorValues(spec.ProcessorType, spec.Processors)
	if err == nil {
		return nil
	}
	// Re-point the field path to the template's spec location.
	return field.Invalid(field.NewPath("spec", "template", "spec", "processors"), spec.Processors, err.Detail)
}
