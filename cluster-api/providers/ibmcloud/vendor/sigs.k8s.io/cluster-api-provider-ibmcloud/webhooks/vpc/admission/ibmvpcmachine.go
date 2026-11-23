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

package vpc

import (
	"context"

	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta2"
)

// Ensure IBMVPCMachine implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMVPCMachine] = &IBMVPCMachine{}
	_ admission.Defaulter[*infrav1.IBMVPCMachine] = &IBMVPCMachine{}
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,verbs=create;update,versions=v1beta2,name=mibmvpcmachine.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachine,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,versions=v1beta2,name=vibmvpcmachine.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMVPCMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMVPCMachine{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMVPCMachine implements a validation and defaulting webhook for IBMVPCMachine.
type IBMVPCMachine struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMVPCMachine) Default(_ context.Context, obj *infrav1.IBMVPCMachine) error {
	defaultIBMVPCMachineSpec(&obj.Spec)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateCreate(_ context.Context, obj *infrav1.IBMVPCMachine) (admission.Warnings, error) {
	allErrs := validateIBMVPCMachineVolume(obj.Spec)
	return nil, aggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateUpdate(_ context.Context, _, _ *infrav1.IBMVPCMachine) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateDelete(_ context.Context, _ *infrav1.IBMVPCMachine) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMVPCMachineVolume(spec infrav1.IBMVPCMachineSpec) field.ErrorList {
	return validateVolumes(spec)
}
