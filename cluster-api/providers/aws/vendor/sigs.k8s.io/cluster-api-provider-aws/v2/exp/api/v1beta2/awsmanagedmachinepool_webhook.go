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
	"fmt"
	"reflect"

	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/klog/v2"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/eks"
)

const (
	maxNodegroupNameLength = 64
)

// log is for logging in this package.
var mmpLog = ctrl.Log.WithName("awsmanagedmachinepool-resource")

// SetupWebhookWithManager will setup the webhooks for the AWSManagedMachinePool.
func (r *AWSManagedMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=validation.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmanagedmachinepools,versions=v1beta2,name=default.awsmanagedmachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &AWSManagedMachinePool{}
var _ webhook.Validator = &AWSManagedMachinePool{}

func (r *AWSManagedMachinePool) validateScaling() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.Scaling != nil { //nolint:nestif
		minField := field.NewPath("spec", "scaling", "minSize")
		maxField := field.NewPath("spec", "scaling", "maxSize")
		minSize := r.Spec.Scaling.MinSize
		maxSize := r.Spec.Scaling.MaxSize
		if minSize != nil {
			if *minSize < 0 {
				allErrs = append(allErrs, field.Invalid(minField, *minSize, "must be greater or equal zero"))
			}
			if maxSize != nil && *maxSize < *minSize {
				allErrs = append(allErrs, field.Invalid(maxField, *maxSize, fmt.Sprintf("must be greater than field %s", minField.String())))
			}
		}
		if maxSize != nil && *maxSize < 0 {
			allErrs = append(allErrs, field.Invalid(maxField, *maxSize, "must be greater than zero"))
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

func (r *AWSManagedMachinePool) validateLaunchTemplate() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.AWSLaunchTemplate == nil {
		return allErrs
	}

	if r.Spec.InstanceType != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "InstanceType"), r.Spec.InstanceType, "InstanceType cannot be specified when LaunchTemplate is specified"))
	}
	if r.Spec.DiskSize != nil {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "DiskSize"), r.Spec.DiskSize, "DiskSize cannot be specified when LaunchTemplate is specified"))
	}

	if r.Spec.AWSLaunchTemplate.IamInstanceProfile != "" {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "AWSLaunchTemplate", "IamInstanceProfile"), r.Spec.AWSLaunchTemplate.IamInstanceProfile, "IAM instance profile in launch template is prohibited in EKS managed node group"))
	}

	return allErrs
}

// ValidateCreate will do any extra validation when creating a AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ValidateCreate() (admission.Warnings, error) {
	mmpLog.Info("AWSManagedMachinePool validate create", "managed-machine-pool", klog.KObj(r))

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
	if errs := r.validateLaunchTemplate(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate will do any extra validation when updating a AWSManagedMachinePool.
func (r *AWSManagedMachinePool) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	mmpLog.Info("AWSManagedMachinePool validate update", "managed-machine-pool", klog.KObj(r))
	oldPool, ok := old.(*AWSManagedMachinePool)
	if !ok {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSManagedMachinePool").GroupKind(), r.Name, field.ErrorList{
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
	if errs := r.validateLaunchTemplate(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
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

// ValidateDelete allows you to add any extra validation when deleting.
func (r *AWSManagedMachinePool) ValidateDelete() (admission.Warnings, error) {
	mmpLog.Info("AWSManagedMachinePool validate delete", "managed-machine-pool", klog.KObj(r))

	return nil, nil
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
	appendErrorIfMutated(old.Spec.AvailabilityZones, r.Spec.AvailabilityZones, "availabilityZones")
	appendErrorIfMutated(old.Spec.AvailabilityZoneSubnetType, r.Spec.AvailabilityZoneSubnetType, "availabilityZoneSubnetType")
	if (old.Spec.AWSLaunchTemplate != nil && r.Spec.AWSLaunchTemplate == nil) ||
		(old.Spec.AWSLaunchTemplate == nil && r.Spec.AWSLaunchTemplate != nil) {
		allErrs = append(
			allErrs,
			field.Invalid(field.NewPath("spec", "AWSLaunchTemplate"), old.Spec.AWSLaunchTemplate, "field is immutable"),
		)
	}
	if old.Spec.AWSLaunchTemplate != nil && r.Spec.AWSLaunchTemplate != nil {
		appendErrorIfMutated(old.Spec.AWSLaunchTemplate.Name, r.Spec.AWSLaunchTemplate.Name, "awsLaunchTemplate.name")
	}

	return allErrs
}

// Default will set default values for the AWSManagedMachinePool.
func (r *AWSManagedMachinePool) Default() {
	mmpLog.Info("AWSManagedMachinePool setting defaults", "managed-machine-pool", klog.KObj(r))

	if r.Spec.EKSNodegroupName == "" {
		mmpLog.Info("EKSNodegroupName is empty, generating name")
		name, err := eks.GenerateEKSName(r.Name, r.Namespace, maxNodegroupNameLength)
		if err != nil {
			mmpLog.Error(err, "failed to create EKS nodegroup name")
			return
		}

		mmpLog.Info("Generated EKSNodegroupName", "nodegroup", klog.KRef(r.Namespace, name))
		r.Spec.EKSNodegroupName = name
	}

	if r.Spec.UpdateConfig == nil {
		r.Spec.UpdateConfig = &UpdateConfig{
			MaxUnavailable: ptr.To[int](1),
		}
	}
}
