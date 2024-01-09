/*
Copyright 2021 The Kubernetes Authors.

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

package v1beta1

import (
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *VSphereFailureDomain) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-vspherefailuredomain,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspherefailuredomains,versions=v1beta1,name=validation.vspherefailuredomain.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-vspherefailuredomain,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspherefailuredomains,verbs=create;update,versions=v1beta1,name=default.vspherefailuredomain.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Validator = &VSphereFailureDomain{}

var _ webhook.Defaulter = &VSphereFailureDomain{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereFailureDomain) ValidateCreate() (admission.Warnings, error) {
	var allErrs field.ErrorList

	if r.Spec.Topology.ComputeCluster == nil && r.Spec.Topology.Hosts != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), "cannot be empty if Hosts is not empty"))
	}

	if r.Spec.Region.Type == HostGroupFailureDomain {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Region", "Type"), fmt.Sprintf("region's Failure Domain type cannot be %s", r.Spec.Region.Type)))
	}

	if r.Spec.Zone.Type == HostGroupFailureDomain && r.Spec.Topology.Hosts == nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "Hosts"), fmt.Sprintf("cannot be nil if zone's Failure Domain type is %s", r.Spec.Zone.Type)))
	}

	if r.Spec.Region.Type == ComputeClusterFailureDomain && r.Spec.Topology.ComputeCluster == nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), fmt.Sprintf("cannot be nil if region's Failure Domain type is %s", r.Spec.Region.Type)))
	}

	if r.Spec.Zone.Type == ComputeClusterFailureDomain && r.Spec.Topology.ComputeCluster == nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), fmt.Sprintf("cannot be nil if zone's Failure Domain type is %s", r.Spec.Zone.Type)))
	}

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereFailureDomain) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	oldVSphereFailureDomain, ok := old.(*VSphereFailureDomain)
	if !ok || !reflect.DeepEqual(r.Spec, oldVSphereFailureDomain.Spec) {
		return nil, field.Forbidden(field.NewPath("spec"), "VSphereFailureDomainSpec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *VSphereFailureDomain) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *VSphereFailureDomain) Default() {
	if r.Spec.Zone.AutoConfigure == nil {
		r.Spec.Zone.AutoConfigure = pointer.Bool(false)
	}

	if r.Spec.Region.AutoConfigure == nil {
		r.Spec.Region.AutoConfigure = pointer.Bool(false)
	}
}
