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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime/schema"
	"k8s.io/apimachinery/pkg/util/validation/field"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/vpc/v1beta2"
)

// Ensure IBMVPCCluster implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMVPCCluster] = &IBMVPCCluster{}
	_ admission.Defaulter[*infrav1.IBMVPCCluster] = &IBMVPCCluster{}
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpccluster,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters,verbs=create;update,versions=v1beta2,name=mibmvpccluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmvpccluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmvpcclusters,versions=v1beta2,name=vibmvpccluster.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMVPCCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMVPCCluster{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMVPCCluster implements a validation and defaulting webhook for IBMVPCCluster.
type IBMVPCCluster struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMVPCCluster) Default(_ context.Context, _ *infrav1.IBMVPCCluster) error {
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateCreate(_ context.Context, obj *infrav1.IBMVPCCluster) (admission.Warnings, error) {
	return validateIBMVPCCluster(obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateUpdate(_ context.Context, _, newObj *infrav1.IBMVPCCluster) (warnings admission.Warnings, err error) {
	return validateIBMVPCCluster(newObj)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMVPCCluster) ValidateDelete(_ context.Context, _ *infrav1.IBMVPCCluster) (admission.Warnings, error) {
	return nil, nil
}

func validateIBMVPCCluster(vpcCluster *infrav1.IBMVPCCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := validateIBMVPCClusterControlPlane(vpcCluster); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		schema.GroupKind{Group: "infrastructure.cluster.x-k8s.io", Kind: "IBMVPCCluster"},
		vpcCluster.Name, allErrs)
}

func validateIBMVPCClusterControlPlane(vpcCluster *infrav1.IBMVPCCluster) *field.Error {
	if vpcCluster.Spec.ControlPlaneEndpoint.Host == "" && vpcCluster.Spec.ControlPlaneLoadBalancer == nil {
		return field.Invalid(field.NewPath(""), "", "One of - ControlPlaneEndpoint or ControlPlaneLoadBalancer must be specified")
	}
	return nil
}
