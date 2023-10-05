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

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/pointer"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	"sigs.k8s.io/cluster-api-provider-aws/pkg/eks"
)

const (
	maxNodegroupNameLength = 64
)

// log is for logging in this package.
var mmpLog = logf.Log.WithName("awsmanagedmachinepool-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSManagedMachinePool.
func (r *AWSManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-awsmanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta1,name=validation.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-awsmanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta1,name=default.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Defaulter = &AWSManagedMachinePool{}
var _ webhook.Validator = &AWSManagedMachinePool{}

func (r *AWSManagedMachinePool) validateScaling() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.Scaling != nil { // nolint:nestif
		minField := field.NewPath("spec", "scaling", "minSize")
		maxField := field.NewPath("spec", "scaling", "maxSize")
		min := r.Spec.Scaling.MinSize
		max := r.Spec.Scaling.MaxSize
		if min != nil {
			if *min < 0 {
				allErrs = append(allErrs, field.Invalid(minField, *min, "must be greater or equal zero"))
			}
			if max != nil && *max < *min {
				allErrs = append(allErrs, field.Invalid(maxField, *max, fmt.Sprintf("must be greater than field %s", minField.String())))
			}
		}
		if max != nil && *max < 0 {
			allErrs = append(allErrs, field.Invalid(maxField, *max, "must be greater than zero"))
		}
	}
	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func (r *AWSManagedMachinePool) validateNodegroupUpdateConfig() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.UpdateConfig != nil {
		nodegroupUpdateConfigField := field.NewPath("spec", "updateConfig")

		if r.Spec.UpdateConfig.MaxUnavailable == nil && r.Spec.UpdateConfig.MaxUnavailablePercentage == nil {
			allErrs = append(allErrs, field.Invalid(nodegroupUpdateConfigField, r.Spec.UpdateConfig, "must specify one of maxUnavailable or maxUnavailablePercentage when using nodegroup updateconfig"))
		}

		if r.Spec.UpdateConfig.MaxUnavailable != nil && r.Spec.UpdateConfig.MaxUnavailablePercentage != nil {
			allErrs = append(allErrs, field.Invalid(nodegroupUpdateConfigField, r.Spec.UpdateConfig, "cannot specify both maxUnavailable and maxUnavailablePercentage"))
		}
	}

	if len(allErrs) == 0 {
		return nil
	}
	return allErrs
}

func (r *AWSManagedMachinePool) validateRemoteAccess() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.RemoteAccess == nil {
		return allErrs
	}
	remoteAccessPath := field.NewPath("spec", "remoteAccess")
	sourceSecurityGroups := r.Spec.RemoteAccess.SourceSecurityGroups

	if public := r.Spec.RemoteAccess.Public; public && len(sourceSecurityGroups) > 0 {
		allErrs = append(
			allErrs,
			field.Invalid(remoteAccessPath.Child("sourceSecurityGroups"), sourceSecurityGroups, "must be empty if public is set"),
		)
	}

	return allErrs
}

// ValidateCreate will do any extra validation when creating a AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ValidateCreate() error {
	mmpLog.Info("AWSManagedMachinePool validate create", "name", r.Name)

	var allErrs field.ErrorList

	if r.Spec.EKSNodegroupName == "" {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.eksNodegroupName"), "eksNodegroupName is required"))
	}
	if errs := r.validateScaling(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateRemoteAccess(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateNodegroupUpdateConfig(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ValidateUpdate(old runtime.Object) error {
	mmpLog.Info("AWSManagedMachinePool validate update", "name", r.Name)
	oldPool, ok := old.(*AWSManagedMachinePool)
	if !ok {
		return apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedMachinePool").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old AWSManagedMachinePool to object")),
		})
	}

	var allErrs field.ErrorList
	allErrs = append(allErrs, r.validateImmutable(oldPool)...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if errs := r.validateScaling(); errs != nil || len(errs) == 0 {
		allErrs = append(allErrs, errs...)
	}
	if errs := r.validateNodegroupUpdateConfig(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil
	}

	return apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete allows you to add any extra validation when deleting.
func (r *AWSManagedMachinePool) ValidateDelete() error {
	mmpLog.Info("AWSManagedMachinePool validate delete", "name", r.Name)

	return nil
}

func (r *AWSManagedMachinePool) validateImmutable(old *AWSManagedMachinePool) field.ErrorList {
	var allErrs field.ErrorList

	appendErrorIfMutated := func(old, update interface{}, name string) {
		if !cmp.Equal(old, update) {
			allErrs = append(
				allErrs,
				field.Invalid(field.NewPath("spec", name), update, "field is immutable"),
			)
		}
	}
	appendErrorIfSetAndMutated := func(old, update interface{}, name string) {
		if !reflect.ValueOf(old).IsZero() && !cmp.Equal(old, update) {
			allErrs = append(
				allErrs,
				field.Invalid(field.NewPath("spec", name), update, "field is immutable"),
			)
		}
	}

	if old.Spec.EKSNodegroupName != "" {
		appendErrorIfMutated(old.Spec.EKSNodegroupName, r.Spec.EKSNodegroupName, "eksNodegroupName")
	}
	appendErrorIfMutated(old.Spec.SubnetIDs, r.Spec.SubnetIDs, "subnetIDs")
	appendErrorIfSetAndMutated(old.Spec.RoleName, r.Spec.RoleName, "roleName")
	appendErrorIfMutated(old.Spec.DiskSize, r.Spec.DiskSize, "diskSize")
	appendErrorIfMutated(old.Spec.AMIType, r.Spec.AMIType, "amiType")
	appendErrorIfMutated(old.Spec.RemoteAccess, r.Spec.RemoteAccess, "remoteAccess")
	appendErrorIfSetAndMutated(old.Spec.CapacityType, r.Spec.CapacityType, "capacityType")

	return allErrs
}

// Default will set default values for the AWSManagedMachinePool.
func (r *AWSManagedMachinePool) Default() {
	mmpLog.Info("AWSManagedMachinePool setting defaults", "name", r.Name)

	if r.Spec.EKSNodegroupName == "" {
		mmpLog.Info("EKSNodegroupName is empty, generating name")
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxNodegroupNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
			return
		}

		mmpLog.Info("Generated EKSNodegroupName", "nodegroup-name", name)
		r.Spec.EKSNodegroupName = name
	}

	if r.Spec.UpdateConfig == nil {
		r.Spec.UpdateConfig = &UpdateConfig{
			MaxUnavailable: pointer.Int(1),
		}
	}
}
