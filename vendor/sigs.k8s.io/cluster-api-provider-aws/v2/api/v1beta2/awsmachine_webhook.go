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
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-aws/v2/feature"
)

// log is for logging in this package.
var log = ctrl.Log.WithName("awsmachine-resource")

func (r *AWSMachine) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,versions=v1beta2,name=validation.awsmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsmachine,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=awsmachines,versions=v1beta2,name=mawsmachine.kb.io,name=mutation.awsmachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var (
	_ webhook.Validator = &AWSMachine{}
	_ webhook.Defaulter = &AWSMachine{}
)

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachine) ValidateCreate() (admission.Warnings, error) {
	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateCloudInitSecret()...)
	allErrs = append(allErrs, r.validateIgnitionAndCloudInit()...)
	allErrs = append(allErrs, r.validateRootVolume()...)
	allErrs = append(allErrs, r.validateNonRootVolumes()...)
	allErrs = append(allErrs, r.validateSSHKeyName()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachine) ValidateUpdate(old runtime.Object) (admission.Warnings, error) {
	newAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new AWSMachine to unstructured object")),
		})
	}
	oldAWSMachine, err := runtime.DefaultUnstructuredConverter.ToUnstructured(old)
	if err != nil {
		return nil, apierrors.NewInvalid(GroupVersion.WithKind("AWSMachine").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old AWSMachine to unstructured object")),
		})
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.validateCloudInitSecret()...)
	allErrs = append(allErrs, r.validateAdditionalSecurityGroups()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	newAWSMachineSpec := newAWSMachine["spec"].(map[string]interface{})
	oldAWSMachineSpec := oldAWSMachine["spec"].(map[string]interface{})

	// allow changes to providerID
	delete(oldAWSMachineSpec, "providerID")
	delete(newAWSMachineSpec, "providerID")

	// allow changes to instanceID
	delete(oldAWSMachineSpec, "instanceID")
	delete(newAWSMachineSpec, "instanceID")

	// allow changes to additionalTags
	delete(oldAWSMachineSpec, "additionalTags")
	delete(newAWSMachineSpec, "additionalTags")

	// allow changes to additionalSecurityGroups
	delete(oldAWSMachineSpec, "additionalSecurityGroups")
	delete(newAWSMachineSpec, "additionalSecurityGroups")

	// allow changes to secretPrefix, secretCount, and secureSecretsBackend
	if cloudInit, ok := oldAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretPrefix")
		delete(cloudInit, "secretCount")
		delete(cloudInit, "secureSecretsBackend")
	}

	if cloudInit, ok := newAWSMachineSpec["cloudInit"].(map[string]interface{}); ok {
		delete(cloudInit, "secretPrefix")
		delete(cloudInit, "secretCount")
		delete(cloudInit, "secureSecretsBackend")
	}

	// allow changes to enableResourceNameDNSAAAARecord and enableResourceNameDNSARecord
	if privateDNSName, ok := oldAWSMachineSpec["privateDnsName"].(map[string]interface{}); ok {
		delete(privateDNSName, "enableResourceNameDnsAAAARecord")
		delete(privateDNSName, "enableResourceNameDnsARecord")
	}

	if privateDNSName, ok := newAWSMachineSpec["privateDnsName"].(map[string]interface{}); ok {
		delete(privateDNSName, "enableResourceNameDnsAAAARecord")
		delete(privateDNSName, "enableResourceNameDnsARecord")
	}

	if !cmp.Equal(oldAWSMachineSpec, newAWSMachineSpec) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec"), "cannot be modified"))
	}

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

func (r *AWSMachine) validateCloudInitSecret() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.CloudInit.InsecureSkipSecretsManager {
		if r.Spec.CloudInit.SecretPrefix != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretPrefix"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if r.Spec.CloudInit.SecretCount != 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
		if r.Spec.CloudInit.SecureSecretsBackend != "" {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secureSecretsBackend"), "cannot be set if spec.cloudInit.insecureSkipSecretsManager is true"))
		}
	}

	if (r.Spec.CloudInit.SecretPrefix != "") != (r.Spec.CloudInit.SecretCount != 0) {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit", "secretCount"), "must be set together with spec.CloudInit.SecretPrefix"))
	}

	return allErrs
}

func (r *AWSMachine) cloudInitConfigured() bool {
	configured := false

	configured = configured || r.Spec.CloudInit.SecretPrefix != ""
	configured = configured || r.Spec.CloudInit.SecretCount != 0
	configured = configured || r.Spec.CloudInit.SecureSecretsBackend != ""
	configured = configured || r.Spec.CloudInit.InsecureSkipSecretsManager

	return configured
}

func (r *AWSMachine) ignitionEnabled() bool {
	return r.Spec.Ignition != nil
}

func (r *AWSMachine) validateIgnitionAndCloudInit() field.ErrorList {
	var allErrs field.ErrorList

	// Feature gate is not enabled but ignition is enabled then send a forbidden error.
	if !feature.Gates.Enabled(feature.BootstrapFormatIgnition) && r.ignitionEnabled() {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ignition"),
			"can be set only if the BootstrapFormatIgnition feature gate is enabled"))
	}

	if r.ignitionEnabled() && r.cloudInitConfigured() {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "cloudInit"), "cannot be set if spec.ignition is set"))
	}

	return allErrs
}

func (r *AWSMachine) validateRootVolume() field.ErrorList {
	var allErrs field.ErrorList

	if r.Spec.RootVolume == nil {
		return allErrs
	}

	if VolumeTypesProvisioned.Has(string(r.Spec.RootVolume.Type)) && r.Spec.RootVolume.IOPS == 0 {
		allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.iops"), "iops required if type is 'io1' or 'io2'"))
	}

	if r.Spec.RootVolume.Throughput != nil {
		if r.Spec.RootVolume.Type != VolumeTypeGP3 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.throughput"), "throughput is valid only for type 'gp3'"))
		}
		if *r.Spec.RootVolume.Throughput < 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.rootVolume.throughput"), "throughput must be nonnegative"))
		}
	}

	if r.Spec.RootVolume.DeviceName != "" {
		log.Info("root volume shouldn't have a device name (this can be ignored if performing a `clusterctl move`)")
	}

	return allErrs
}

func (r *AWSMachine) validateNonRootVolumes() field.ErrorList {
	var allErrs field.ErrorList

	for _, volume := range r.Spec.NonRootVolumes {
		if VolumeTypesProvisioned.Has(string(volume.Type)) && volume.IOPS == 0 {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.iops"), "iops required if type is 'io1' or 'io2'"))
		}

		if volume.Throughput != nil {
			if volume.Type != VolumeTypeGP3 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.throughput"), "throughput is valid only for type 'gp3'"))
			}
			if *volume.Throughput < 0 {
				allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.throughput"), "throughput must be nonnegative"))
			}
		}

		if volume.DeviceName == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("spec.nonRootVolumes.deviceName"), "non root volume should have device name"))
		}
	}

	return allErrs
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AWSMachine) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhook.Defaulter such that an empty CloudInit will be defined with a default
// SecureSecretsBackend as SecretBackendSecretsManager iff InsecureSkipSecretsManager is unset.
func (r *AWSMachine) Default() {
	if !r.Spec.CloudInit.InsecureSkipSecretsManager && r.Spec.CloudInit.SecureSecretsBackend == "" && !r.ignitionEnabled() {
		r.Spec.CloudInit.SecureSecretsBackend = SecretBackendSecretsManager
	}

	if r.ignitionEnabled() && r.Spec.Ignition.Version == "" {
		if r.Spec.Ignition == nil {
			r.Spec.Ignition = &Ignition{}
		}

		r.Spec.Ignition.Version = DefaultIgnitionVersion
	}
}

func (r *AWSMachine) validateAdditionalSecurityGroups() field.ErrorList {
	var allErrs field.ErrorList

	for _, additionalSecurityGroup := range r.Spec.AdditionalSecurityGroups {
		if len(additionalSecurityGroup.Filters) > 0 && additionalSecurityGroup.ID != nil {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec.additionalSecurityGroups"), "only one of ID or Filters may be specified, specifying both is forbidden"))
		}
	}
	return allErrs
}

func (r *AWSMachine) validateSSHKeyName() field.ErrorList {
	return validateSSHKeyName(r.Spec.SSHKeyName)
}
