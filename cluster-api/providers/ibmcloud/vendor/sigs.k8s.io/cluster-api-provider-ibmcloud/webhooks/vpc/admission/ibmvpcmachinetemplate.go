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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta2"
)

// Ensure IBMVPCMachineTemplate implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMVPCMachineTemplate] = &IBMVPCMachineTemplate{}
	_ admission.Defaulter[*infrav1.IBMVPCMachineTemplate] = &IBMVPCMachineTemplate{}
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachinetemplates,verbs=create;update,versions=v1beta2,name=mibmvpcmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachinetemplates,versions=v1beta2,name=vibmvpcmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMVPCMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMVPCMachineTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMVPCMachineTemplate implements a validation and defaulting webhook for IBMVPCMachineTemplate.
type IBMVPCMachineTemplate struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) Default(_ context.Context, obj *infrav1.IBMVPCMachineTemplate) error {
	defaultIBMVPCMachineSpec(&obj.Spec.Template.Spec)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateCreate(_ context.Context, obj *infrav1.IBMVPCMachineTemplate) (admission.Warnings, error) {
	allErrs := validateIBMVPCMachineVolume(obj.Spec.Template.Spec)
	return nil, aggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateUpdate(_ context.Context, _, _ *infrav1.IBMVPCMachineTemplate) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateDelete(_ context.Context, _ *infrav1.IBMVPCMachineTemplate) (admission.Warnings, error) {
	return nil, nil
}
