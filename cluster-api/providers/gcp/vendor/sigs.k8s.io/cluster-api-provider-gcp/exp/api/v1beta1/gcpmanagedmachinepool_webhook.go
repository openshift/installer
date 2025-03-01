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
	"fmt"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

const (
	maxNodePoolNameLength = 40
)

// log is for logging in this package.
var gcpmanagedmachinepoollog = logf.Log.WithName("gcpmanagedmachinepool-resource")

func (r *GCPManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedmachinepool,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedmachinepools,verbs=create;update,versions=v1beta1,name=mgcpmanagedmachinepool.kb.io,admissionReviewVersions=v1

var _ webhook.Defaulter = &GCPManagedMachinePool{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *GCPManagedMachinePool) Default() {
	gcpmanagedmachinepoollog.Info("default", "name", r.Name)
}

//+kubebuilder:webhook:path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedmachinepool,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedmachinepools,verbs=create;update,versions=v1beta1,name=vgcpmanagedmachinepool.kb.io,admissionReviewVersions=v1

var _ webhook.Validator = &GCPManagedMachinePool{}

// validateSpec validates that the GCPManagedMachinePool spec is valid.
func (r *GCPManagedMachinePool) validateSpec() field.ErrorList {
	var allErrs field.ErrorList

	// Validate node pool name
	if len(r.Spec.NodePoolName) > maxNodePoolNameLength {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "NodePoolName"),
				r.Spec.NodePoolName, fmt.Sprintf("node pool name cannot have more than %d characters", maxNodePoolNameLength)),
		)
	}

	if errs := r.validateScaling(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := r.validateNonNegative(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

// validateScaling validates that the GCPManagedMachinePool autoscaling spec is valid.
func (r *GCPManagedMachinePool) validateScaling() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.Scaling != nil {
		minField := field.NewPath("spec", "scaling", "minCount")
		maxField := field.NewPath("spec", "scaling", "maxCount")
		locationPolicyField := field.NewPath("spec", "scaling", "locationPolicy")

		minCount := r.Spec.Scaling.MinCount
		maxCount := r.Spec.Scaling.MaxCount
		locationPolicy := r.Spec.Scaling.LocationPolicy

		// cannot specify autoscaling config if autoscaling is disabled
		if r.Spec.Scaling.EnableAutoscaling != nil && !*r.Spec.Scaling.EnableAutoscaling {
			if minCount != nil {
				allErrs = append(allErrs, field.Forbidden(minField, "minCount cannot be specified when autoscaling is disabled"))
			}
			if maxCount != nil {
				allErrs = append(allErrs, field.Forbidden(maxField, "maxCount cannot be specified when autoscaling is disabled"))
			}
			if locationPolicy != nil {
				allErrs = append(allErrs, field.Forbidden(locationPolicyField, "locationPolicy cannot be specified when autoscaling is disabled"))
			}
		}

		if minCount != nil {
			// validates min >= 0
			if *minCount < 0 {
				allErrs = append(allErrs, field.Invalid(minField, *minCount, "must be greater or equal zero"))
			}
			// validates min <= max
			if maxCount != nil && *maxCount < *minCount {
				allErrs = append(allErrs, field.Invalid(maxField, *maxCount, "must be greater than field "+minField.String()))
			}
		}
	}
	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func appendErrorIfNegative[T int32 | int64](value *T, name string, errs *field.ErrorList) {
	if value != nil && *value < 0 {
		*errs = append(*errs, field.Invalid(field.NewPath("spec", name), *value, "must be non-negative"))
	}
}

// validateNonNegative validates that non-negative GCPManagedMachinePool spec fields are not negative.
func (r *GCPManagedMachinePool) validateNonNegative() field.ErrorList {
	var allErrs field.ErrorList

	appendErrorIfNegative(r.Spec.DiskSizeGb, "diskSizeGb", &allErrs)
	appendErrorIfNegative(r.Spec.MaxPodsPerNode, "maxPodsPerNode", &allErrs)
	appendErrorIfNegative(r.Spec.LocalSsdCount, "localSsdCount", &allErrs)

	return allErrs
}

func appendErrorIfMutated(old, update interface{}, name string, errs *field.ErrorList) {
	if !cmp.Equal(old, update) {
		*errs = append(
			*errs,
			field.Invalid(field.NewPath("spec", name), update, "field is immutable"),
		)
	}
}

// validateImmutable validates that immutable GCPManagedMachinePool spec fields are not mutated.
func (r *GCPManagedMachinePool) validateImmutable(old *GCPManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList

	appendErrorIfMutated(old.Spec.InstanceType, r.Spec.InstanceType, "instanceType", &allErrs)
	appendErrorIfMutated(old.Spec.NodePoolName, r.Spec.NodePoolName, "nodePoolName", &allErrs)
	appendErrorIfMutated(old.Spec.MachineType, r.Spec.MachineType, "machineType", &allErrs)
	appendErrorIfMutated(old.Spec.DiskSizeGb, r.Spec.DiskSizeGb, "diskSizeGb", &allErrs)
	appendErrorIfMutated(old.Spec.DiskType, r.Spec.DiskType, "diskType", &allErrs)
	appendErrorIfMutated(old.Spec.LocalSsdCount, r.Spec.LocalSsdCount, "localSsdCount", &allErrs)
	appendErrorIfMutated(old.Spec.Management, r.Spec.Management, "management", &allErrs)
	appendErrorIfMutated(old.Spec.MaxPodsPerNode, r.Spec.MaxPodsPerNode, "maxPodsPerNode", &allErrs)
	appendErrorIfMutated(old.Spec.NodeNetwork.PodRangeName, r.Spec.NodeNetwork.PodRangeName, "podRangeName", &allErrs)
	appendErrorIfMutated(old.Spec.NodeNetwork.CreatePodRange, r.Spec.NodeNetwork.CreatePodRange, "createPodRange", &allErrs)
	appendErrorIfMutated(old.Spec.NodeNetwork.PodRangeCidrBlock, r.Spec.NodeNetwork.PodRangeCidrBlock, "podRangeCidrBlock", &allErrs)
	appendErrorIfMutated(old.Spec.NodeSecurity, r.Spec.NodeSecurity, "nodeSecurity", &allErrs)

	return allErrs
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedMachinePool) ValidateCreate() (admission.Warnings, error) {
	gcpmanagedmachinepoollog.Info("validate create", "name", r.Name)
	var allErrs field.ErrorList

	if errs := r.validateSpec(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind("GCPManagedMachinePool").GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedMachinePool) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	gcpmanagedmachinepoollog.Info("validate update", "name", r.Name)
	var allErrs field.ErrorList
	old := oldRaw.(*GCPManagedMachinePool)

	if errs := r.validateImmutable(old); errs != nil {
		allErrs = append(allErrs, errs...)
	}

	if errs := r.validateSpec(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind("GCPManagedMachinePool").GroupKind(), r.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *GCPManagedMachinePool) ValidateDelete() (admission.Warnings, error) {
	gcpmanagedmachinepoollog.Info("validate delete", "name", r.Name)

	return nil, nil
}
