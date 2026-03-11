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
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1beta2 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachinetemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachinetemplates,verbs=create;update,versions=v1beta2,name=mibmvpcmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpcmachinetemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcmachinetemplates,versions=v1beta2,name=vibmvpcmachinetemplate.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMVPCMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1beta2.IBMVPCMachineTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMVPCMachineTemplate implements a validation and defaulting webhook for IBMVPCMachineTemplate.
type IBMVPCMachineTemplate struct{}

var _ webhook.CustomDefaulter = &IBMVPCMachineTemplate{}
var _ webhook.CustomValidator = &IBMVPCMachineTemplate{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) Default(_ context.Context, obj runtime.Object) error {
	objValue, ok := obj.(*infrav1beta2.IBMVPCMachineTemplate)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a IBMVPCMachineTemplate but got a %T", obj))
	}
	defaultIBMVPCMachineSpec(&objValue.Spec.Template.Spec)
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	objValue, ok := obj.(*infrav1beta2.IBMVPCMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a IBMVPCMachineTemplate but got a %T", obj))
	}
	var allErrs field.ErrorList
	allErrs = append(allErrs, validateIBMVPCMachineBootVolume(objValue.Spec.Template.Spec)...)

	return nil, aggregateObjErrors(objValue.GroupVersionKind().GroupKind(), objValue.Name, allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateUpdate(_ context.Context, _, _ runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMVPCMachineTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
