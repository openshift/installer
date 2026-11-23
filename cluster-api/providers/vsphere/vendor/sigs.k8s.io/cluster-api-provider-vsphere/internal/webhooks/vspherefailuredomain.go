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

package webhooks

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/conversion"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-vspherefailuredomain,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspherefailuredomains,versions=v1beta2,name=validation.vspherefailuredomain.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereFailureDomain implements a validation and defaulting webhook for VSphereFailureDomain.
type VSphereFailureDomain struct{}

var _ admission.Validator[*infrav1.VSphereFailureDomain] = &VSphereFailureDomain{}

func (webhook *VSphereFailureDomain) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.VSphereFailureDomain{}).
		WithValidator(webhook).
		WithConverter(conversion.VSphereFailureDomain).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereFailureDomain) ValidateCreate(_ context.Context, obj *infrav1.VSphereFailureDomain) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if obj.Spec.Topology.ComputeCluster == "" && obj.Spec.Topology.Hosts.IsDefined() {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), "cannot be empty if Hosts is not empty"))
	}

	if obj.Spec.Region.Type == infrav1.HostGroupFailureDomain {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Region", "Type"), fmt.Sprintf("region's Failure Domain type cannot be %s", obj.Spec.Region.Type)))
	}

	if obj.Spec.Zone.Type == infrav1.HostGroupFailureDomain && !obj.Spec.Topology.Hosts.IsDefined() {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "Hosts"), fmt.Sprintf("cannot be nil if zone's Failure Domain type is %s", obj.Spec.Zone.Type)))
	}

	if obj.Spec.Region.Type == infrav1.ComputeClusterFailureDomain && obj.Spec.Topology.ComputeCluster == "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), fmt.Sprintf("cannot be nil if region's Failure Domain type is %s", obj.Spec.Region.Type)))
	}

	if obj.Spec.Zone.Type == infrav1.ComputeClusterFailureDomain && obj.Spec.Topology.ComputeCluster == "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "Topology", "ComputeCluster"), fmt.Sprintf("cannot be nil if zone's Failure Domain type is %s", obj.Spec.Zone.Type)))
	}

	if len(obj.Spec.Topology.NetworkConfigurations) != 0 && len(obj.Spec.Topology.Networks) != 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "topology", "networks"), "cannot be set if spec.topology.networkConfigurations is already set"))
	}

	for i, networkConfig := range obj.Spec.Topology.NetworkConfigurations {
		if networkConfig.NetworkName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec", "topology", "networkConfigurations").Index(i).Child("networkName"), "cannot be empty"))
		}
	}

	return nil, AggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereFailureDomain) ValidateUpdate(_ context.Context, oldTyped, newTyped *infrav1.VSphereFailureDomain) (admission.Warnings, error) {
	if !reflect.DeepEqual(newTyped.Spec, oldTyped.Spec) {
		return nil, field.Forbidden(field.NewPath("spec"), "VSphereFailureDomainSpec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereFailureDomain) ValidateDelete(_ context.Context, _ *infrav1.VSphereFailureDomain) (admission.Warnings, error) {
	return nil, nil
}
