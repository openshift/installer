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
var ibmvpcclusterlog = logf.Log.WithName("ibmvpccluster-resource")

func (r *IBMVPCCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpccluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters,verbs=create;update,versions=v1beta2,name=mibmvpccluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &IBMVPCCluster{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMVPCCluster) Default() {
	ibmvpcclusterlog.Info("default", "name", r.Name)
}

// TODO(user): change verbs to "verbs=create;update;delete" if you want to enable deletion validation.
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpccluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters,versions=v1beta2,name=vibmvpccluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Validator = &IBMVPCCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateCreate() (admission.Warnings, error) {
	ibmvpcclusterlog.Info("validate create", "name", r.Name)
	return r.validateIBMVPCCluster()
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateUpdate(_ runtime.Object) (admission.Warnings, error) {
	ibmvpcclusterlog.Info("validate update", "name", r.Name)
	return r.validateIBMVPCCluster()
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateDelete() (admission.Warnings, error) {
	ibmvpcclusterlog.Info("validate delete", "name", r.Name)
	return nil, nil
}

func (r *IBMVPCCluster) validateIBMVPCCluster() (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := r.validateIBMVPCClusterControlPlane(); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMVPCCluster"},
		r.Name, allErrs)
}

func (r *IBMVPCCluster) validateIBMVPCClusterControlPlane() *field.Error {
	if r.Spec.ControlPlaneEndpoint.Host == "" && r.Spec.ControlPlaneLoadBalancer == nil {
		return field.Invalid(field.NewPath(""), "", "One of - ControlPlaneEndpoint or ControlPlaneLoadBalancer must be specified")
	}
	return nil
}
