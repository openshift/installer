/*
Copyright 2026 The Kubernetes Authors.

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

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
	"sigs.k8s.io/cluster-api/util/topology"
)

func (w *AWSMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.AWSMachineTemplate{}).
		WithValidator(w).
		Complete()
}

// AWSMachineTemplate implements a custom validation webhook for AWSMachineTemplate.
// Note: we use a custom validator to access the request context for SSA of AWSMachineTemplate.
// +kubebuilder:object:generate=false
type AWSMachineTemplate struct{}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachinetemplates,versions=v1beta2,name=validation.awsmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomValidator = &AWSMachineTemplate{}

func (w *AWSMachineTemplate) validateRootVolume(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec
	if spec.RootVolume == nil {
		return allErrs
	}

	if infrav1.VolumeTypesProvisioned.Has(string(spec.RootVolume.Type)) && spec.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if spec.RootVolume.Throughput != nil {
		if spec.RootVolume.Type != infrav1.VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.throughput"), "throughput is valid only for type 'gp3'"))
		}
		if *spec.RootVolume.Throughput < 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.rootVolume.throughput"), "throughput must be nonnegative"))
		}
	}

	if spec.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func (w *AWSMachineTemplate) validateNonRootVolumes(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec

	for _, volume := range spec.NonRootVolumes {
		if infrav1.VolumeTypesProvisioned.Has(string(volume.Type)) && volume.IOPS == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.iops"), "iops required if type is 'io1' or 'io2'"))
		}

		if volume.Throughput != nil {
			if volume.Type != infrav1.VolumeTypeGP3 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.throughput"), "throughput is valid only for type 'gp3'"))
			}
			if *volume.Throughput < 0 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.throughput"), "throughput must be nonnegative"))
			}
		}

		if volume.DeviceName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.nonRootVolumes.deviceName"), "non root volume should have device name"))
		}
	}

	return allErrs
}

func (w *AWSMachineTemplate) validateAdditionalSecurityGroups(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec

	for _, additionalSecurityGroup := range spec.AdditionalSecurityGroups {
		if len(additionalSecurityGroup.Filters) > 0 && additionalSecurityGroup.ID != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "additionalSecurityGroups"), "only one of ID or Filters may be specified, specifying both is forbidden"))
		}
	}
	return allErrs
}

func (w *AWSMachineTemplate) validateCloudInitSecret(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec
	if spec.CloudInit.InsecureSkipSecretsManager {
		if spec.CloudInit.SecretPrefix != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretPrefix"), "cannot be set if spec.template.spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if spec.CloudInit.SecretCount != 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretCount"), "cannot be set if spec.template.spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if spec.CloudInit.SecureSecretsBackend != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secureSecretsBackend"), "cannot be set if spec.template.spec.cloudInit.insecureSkipSecretsManager is true"))
		}
	}

	if (spec.CloudInit.SecretPrefix != "") != (spec.CloudInit.SecretCount != 0) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretCount"), "must be set together with spec.template.spec.CloudInit.SecretPrefix"))
	}

	return allErrs
}

func (w *AWSMachineTemplate) cloudInitConfigured(r *infrav1.AWSMachineTemplate) bool {
	spec := r.Spec.Template.Spec
	configured := false

	configured = configured || spec.CloudInit.SecretPrefix != ""
	configured = configured || spec.CloudInit.SecretCount != 0
	configured = configured || spec.CloudInit.SecureSecretsBackend != ""
	configured = configured || spec.CloudInit.InsecureSkipSecretsManager

	return configured
}

func (w *AWSMachineTemplate) ignitionEnabled(r *infrav1.AWSMachineTemplate) bool {
	return r.Spec.Template.Spec.Ignition != nil
}

func (w *AWSMachineTemplate) validateIgnitionAndCloudInit(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) && w.ignitionEnabled(r) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "ignition"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	if w.ignitionEnabled(r) && w.cloudInitConfigured(r) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit"),
			"cannot be set if spec.template.spec.ignition is set"))
	}

	return allErrs
}

func (w *AWSMachineTemplate) validateHostAllocation(r *infrav1.AWSMachineTemplate) field.ErrorList {
	var allErrs field.ErrorList

	spec := r.Spec.Template.Spec

	// Check if both hostID and dynamicHostAllocation are specified
	hasHostID := spec.HostID != nil && len(*spec.HostID) > 0
	hasDynamicHostAllocation := spec.DynamicHostAllocation != nil

	if hasHostID && hasDynamicHostAllocation {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.template.spec.hostID"), "hostID and dynamicHostAllocation are mutually exclusive"), field.Forbidden(field.NewPath("spec.template.spec.dynamicHostAllocation"), "hostID and dynamicHostAllocation are mutually exclusive"))
	}

	// HostID, HostAffinity, and DynamicHostAllocation can only be set when Tenancy is "host"
	if hasHostID && spec.Tenancy != hostTenancy {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.template.spec.hostID"), "hostID can only be set when tenancy is 'host'"))
	}

	if spec.HostAffinity != nil && *spec.HostAffinity == hostAffinity && spec.Tenancy != hostTenancy {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.template.spec.hostAffinity"), "hostAffinity can only be set to 'host' when tenancy is 'host'"))
	}

	if hasDynamicHostAllocation && spec.Tenancy != hostTenancy {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.template.spec.dynamicHostAllocation"), "dynamicHostAllocation can only be set when tenancy is 'host'"))
	}

	// When hostAffinity is "host", either hostID or dynamicHostAllocation must be specified
	if spec.HostAffinity != nil && *spec.HostAffinity == hostAffinity && !hasHostID && !hasDynamicHostAllocation {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.template.spec.hostID"), "hostID or dynamicHostAllocation must be set when hostAffinity is 'host'"))
	}

	// DHA needs to have hostAffinity set to "host" to make sure it does not drift off its allocated host when the instance is restarted, otherwise there will be a host not in use still allocated.
	if hasDynamicHostAllocation && (spec.HostAffinity == nil || *spec.HostAffinity != hostAffinity) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.template.spec.dynamicHostAllocation"), "dynamicHostAllocation can only be set when hostAffinity is 'host'"))
	}

	return allErrs
}

func (w *AWSMachineTemplate) validateSSHKeyName(r *infrav1.AWSMachineTemplate) field.ErrorList {
	return validateSSHKeyName(r.Spec.Template.Spec.SSHKeyName)
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (w *AWSMachineTemplate) ValidateCreate(_ context.Context, raw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	obj, ok := raw.(*infrav1.AWSMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a AWSMachineTemplate but got a %T", raw))
	}

	spec := obj.Spec.Template.Spec

	if spec.CloudInit.SecretPrefix != "" {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "cloudInit", "secretPrefix"), "cannot be set in templates"))
	}

	if spec.CloudInit.SecretCount != 0 {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "cannot be set in templates"))
	}

	if spec.ProviderID != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "template", "spec", "providerID"), "cannot be set in templates"))
	}

	allErrs = append(allErrs, w.validateCloudInitSecret(obj)...)
	allErrs = append(allErrs, w.validateIgnitionAndCloudInit(obj)...)
	allErrs = append(allErrs, w.validateRootVolume(obj)...)
	allErrs = append(allErrs, w.validateNonRootVolumes(obj)...)
	allErrs = append(allErrs, w.validateSSHKeyName(obj)...)
	allErrs = append(allErrs, w.validateAdditionalSecurityGroups(obj)...)
	allErrs = append(allErrs, obj.Spec.Template.Spec.AdditionalTags.Validate()...)
	allErrs = append(allErrs, w.validateHostAllocation(obj)...)

	return nil, aggregateObjErrors(obj.GroupVersionKind().GroupKind(), obj.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (w *AWSMachineTemplate) ValidateUpdate(ctx context.Context, oldRaw runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	newAWSMachineTemplate, ok := newRaw.(*infrav1.AWSMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a AWSMachineTemplate but got a %T", newRaw))
	}
	oldAWSMachineTemplate, ok := oldRaw.(*infrav1.AWSMachineTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a AWSMachineTemplate but got a %T", oldRaw))
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	var allErrs field.ErrorList

	if !topology.IsDryRunRequest(req, newAWSMachineTemplate) && !cmp.Equal(newAWSMachineTemplate.Spec, oldAWSMachineTemplate.Spec) {
		if oldAWSMachineTemplate.Spec.Template.Spec.InstanceMetadataOptions == nil {
			oldAWSMachineTemplate.Spec.Template.Spec.InstanceMetadataOptions = newAWSMachineTemplate.Spec.Template.Spec.InstanceMetadataOptions
		}

		if !cmp.Equal(newAWSMachineTemplate.Spec.Template.Spec, oldAWSMachineTemplate.Spec.Template.Spec) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "template", "spec"), newAWSMachineTemplate, "AWSMachineTemplate.Spec is immutable"),
			)
		}
	}

	return nil, aggregateObjErrors(newAWSMachineTemplate.GroupVersionKind().GroupKind(), newAWSMachineTemplate.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (w *AWSMachineTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
