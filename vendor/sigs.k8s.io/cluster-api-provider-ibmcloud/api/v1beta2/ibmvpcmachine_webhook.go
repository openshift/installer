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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var ibmvpcmachinelog = logf.Log.WithName("ibmvpcmachine-resource")

func (r *IBMVPCMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,verbs=create;update,versions=v1beta2,name=mibmvpcmachine.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &IBMVPCMachine{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMVPCMachine) Default() {
	ibmvpcmachinelog.Info("default", "name", r.Name)
	defaultIBMVPCMachineSpec(&r.Spec)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachine,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachines,versions=v1beta2,name=vibmvpcmachine.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Validator = &IBMVPCMachine{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateCreate() (admission.Warnings, error) {
	ibmvpcmachinelog.Info("validate create", "name", r.Name)

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateIBMVPCMachineBootVolume()...)

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	ibmvpcmachinelog.Info("validate update", "name", r.Name)
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachine) ValidateDelete() (admission.Warnings, error) {
	ibmvpcmachinelog.Info("validate delete", "name", r.Name)
	return nil, nil
}

func (r *IBMVPCMachine) validateIBMVPCMachineBootVolume() field.ErrorList {
	return validateBootVolume(r.Spec)
}
