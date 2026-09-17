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

// Ensure IBMPowerVSMachine implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMPowerVSMachine] = &IBMPowerVSMachine{}
	_ admission.Defaulter[*infrav1.IBMPowerVSMachine] = &IBMPowerVSMachine{}
)

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines,versions=v1beta3,name=mibmpowervsmachine.kb.io,sideEffects=None,admissionReviewVersions=v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsmachine,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsmachines,versions=v1beta3,name=vibmpowervsmachine.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *IBMPowerVSMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMPowerVSMachine{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSMachine implements a validation and defaulting webhook for IBMPowerVSMachine.
type IBMPowerVSMachine struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSMachine) Default(_ context.Context, obj *infrav1.IBMPowerVSMachine) error {
	defaultIBMPowerVSMachineSpec(&obj.Spec)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachine) ValidateCreate(_ context.Context, obj *infrav1.IBMPowerVSMachine) (admission.Warnings, error) {
	return validateIBMPowerVSMachine(obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachine) ValidateUpdate(_ context.Context, _, newObj *infrav1.IBMPowerVSMachine) (warnings admission.Warnings, err error) {
	return validateIBMPowerVSMachine(newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSMachine) ValidateDelete(_ context.Context, _ *infrav1.IBMPowerVSMachine) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMPowerVSMachine(machine *infrav1.IBMPowerVSMachine) (admission.Warnings, error) {
	// Network: CRD CEL on ResourceIdentifier enforces exactly-one of ID/Name.
	// Image:   CRD CEL on IBMPowerVSMachineImage enforces type↔reference/import mutual exclusion;
	//          Enum marker rejects invalid type; MinLength=1 on ImageReference.Name rejects empty import name.
	// Memory:  +kubebuilder:validation:Minimum=2 on MemoryGiB enforces the minimum at the CRD level.
	// Processors: intstr.IntOrString has no CRD-expressible minimum, so the webhook is the right layer.
	if err := validateIBMPowerVSMachineProcessors(machine); err != nil {
		return nil, apierrors.NewInvalid(
			schema.GroupKind{Group: infrastructureGroup, Kind: "IBMPowerVSMachine"},
			machine.Name, field.ErrorList{err})
	}
	return nil, nil
}

func validateIBMPowerVSMachineProcessors(machine *infrav1.IBMPowerVSMachine) *field.Error {
	return validateIBMPowerVSProcessorValues(machine.Spec.ProcessorType, machine.Spec.Processors)
}
