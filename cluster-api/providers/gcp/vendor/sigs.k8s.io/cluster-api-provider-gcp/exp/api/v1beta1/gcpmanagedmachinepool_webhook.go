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

package v1beta1

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	webhookutils "sigs.k8s.io/cluster-api-provider-gcp/util/webhook"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	maxNodePoolNameLength = 40
)

// log is for logging in this package.
var gcpmanagedmachinepoollog = logf.Log.WithName("gcpmanagedmachinepooltemplate-resource")

func (r *GCPManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(gcpManagedMachinePoolWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedmachinepool,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedmachinepools,verbs=create;update,versions=v1beta1,name=mgcpmanagedmachinepool.kb.io,admissionReviewVersions=v1

type gcpManagedMachinePoolWebhook struct{}

var _ webhook.CustomDefaulter = &gcpManagedMachinePoolWebhook{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (*gcpManagedMachinePoolWebhook) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedmachinepool,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedmachinepools,verbs=create;update,versions=v1beta1,name=vgcpmanagedmachinepool.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &gcpManagedMachinePoolWebhook{}

func validateNodePoolName(name string, fldPath *field.Path) *field.Error {
	if len(name) > maxNodePoolNameLength {
		return (field.Invalid(
			fldPath,
			name,
			fmt.Sprintf("node pool name cannot have more than %d characters", maxNodePoolNameLength),
		))
	}

	return nil
}

func validateScaling(scaling *NodePoolAutoScaling, minField, maxField, locationPolicyField *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// cannot specify autoscaling config if autoscaling is disabled
	if scaling.EnableAutoscaling != nil && !*scaling.EnableAutoscaling {
		if scaling.MinCount != nil {
			allErrs = append(allErrs, field.Forbidden(minField, "minCount cannot be specified when autoscaling is disabled"))
		}
		if scaling.MaxCount != nil {
			allErrs = append(allErrs, field.Forbidden(maxField, "maxCount cannot be specified when autoscaling is disabled"))
		}
		if scaling.LocationPolicy != nil {
			allErrs = append(allErrs, field.Forbidden(locationPolicyField, "locationPolicy cannot be specified when autoscaling is disabled"))
		}
	}

	if scaling.MinCount != nil {
		// validates min >= 0
		if *scaling.MinCount < 0 {
			allErrs = append(allErrs, field.Invalid(minField, *scaling.MinCount, "must be greater or equal zero"))
		}
		// validates min <= max
		if scaling.MaxCount != nil && *scaling.MaxCount < *scaling.MinCount {
			allErrs = append(allErrs, field.Invalid(maxField, *scaling.MaxCount, "must be greater than field "+minField.String()))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}

	return allErrs
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedMachinePoolWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*GCPManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an GCPManagedMachinePool object but got %T", r)
	}

	gcpmanagedmachinepoollog.Info("Validating GCPManagedMachinePool create", "name", r.Name)

	var allErrs field.ErrorList

	if err := validateNodePoolName(
		r.Spec.NodePoolName,
		field.NewPath("spec", "NodePoolName")); err != nil {
		allErrs = append(allErrs, err)
	}

	if r.Spec.Scaling != nil {
		if errs := validateScaling(
			r.Spec.Scaling,
			field.NewPath("spec", "scaling", "minCount"),
			field.NewPath("spec", "scaling", "maxCount"),
			field.NewPath("spec", "scaling", "locationPolicy"),
		); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "diskSizeGb"),
		r.Spec.DiskSizeGb,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "localSsdCount"),
		r.Spec.LocalSsdCount,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "maxPodsPerNode"),
		r.Spec.MaxPodsPerNode,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedMachinePoolWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	old, ok := oldObj.(*GCPManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an GCPManagedMachinePool object but got %T", old)
	}

	r, ok := newObj.(*GCPManagedMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an GCPManagedMachinePool object but got %T", r)
	}

	gcpmanagedmachinepoollog.Info("Validating GCPManagedMachinePool update", "name", r.Name)

	var allErrs field.ErrorList

	if r.Spec.Scaling != nil {
		if errs := validateScaling(
			r.Spec.Scaling,
			field.NewPath("spec", "scaling", "minCount"),
			field.NewPath("spec", "scaling", "maxCount"),
			field.NewPath("spec", "scaling", "locationPolicy"),
		); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "instanceType"),
		old.Spec.InstanceType,
		r.Spec.InstanceType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodePoolName"),
		old.Spec.NodePoolName,
		r.Spec.NodePoolName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "machineType"),
		old.Spec.MachineType,
		r.Spec.MachineType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "diskSizeGb"),
		old.Spec.DiskSizeGb,
		r.Spec.DiskSizeGb); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "diskType"),
		old.Spec.DiskType,
		r.Spec.DiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "localSsdCount"),
		old.Spec.LocalSsdCount,
		r.Spec.LocalSsdCount); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "management"),
		old.Spec.Management,
		r.Spec.Management); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "maxPodsPerNode"),
		old.Spec.MaxPodsPerNode,
		r.Spec.MaxPodsPerNode); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodeNetwork", "podRangeName"),
		old.Spec.NodeNetwork.PodRangeName,
		r.Spec.NodeNetwork.PodRangeName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodeNetwork", "createPodRange"),
		old.Spec.NodeNetwork.CreatePodRange,
		r.Spec.NodeNetwork.CreatePodRange); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodeNetwork", "podRangeCidrBlock"),
		old.Spec.NodeNetwork.PodRangeCidrBlock,
		r.Spec.NodeNetwork.PodRangeCidrBlock); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodeSecurity"),
		old.Spec.NodeSecurity,
		r.Spec.NodeSecurity); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind("GCPManagedMachinePool").GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedMachinePoolWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
