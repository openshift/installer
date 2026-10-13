/*
Copyright 2023 The Kubernetes Authors.

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
	"context"
	"fmt"
	"net"

	"github.com/blang/semver"
	kmsArnRegexpValidator "github.com/openshift-online/ocm-common/pkg/resource/validations"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the webhooks for the ROSAControlPlane.
func (r *ROSAControlPlane) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(rosaControlPlaneWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-controlplane-cluster-x-k8s-io-v1beta2-rosacontrolplane,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,versions=v1beta2,name=validation.rosacontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-controlplane-cluster-x-k8s-io-v1beta2-rosacontrolplane,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,versions=v1beta2,name=default.rosacontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type rosaControlPlaneWebhook struct{}

var _ webhook.CustomDefaulter = &rosaControlPlaneWebhook{}
var _ webhook.CustomValidator = &rosaControlPlaneWebhook{}

// ValidateCreate implements admission.Validator.
func (*rosaControlPlaneWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	r, ok := obj.(*ROSAControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected an ROSAControlPlane object but got %T", r)
	}

	var allErrs field.ErrorList

	if err := r.validateVersion(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateEtcdEncryptionKMSArn(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateExternalAuthProviders(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateClusterRegistryConfig(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateRosaRoleConfig(); err != nil {
		allErrs = append(allErrs, err)
	}

	allErrs = append(allErrs, r.validateNetwork()...)
	allErrs = append(allErrs, r.Spec.AdditionalTags.Validate()...)

	if err := r.validateROSANetwork(); err != nil {
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

func (r *ROSAControlPlane) validateClusterRegistryConfig() *field.Error {
	if r.Spec.ClusterRegistryConfig != nil {
		if r.Spec.ClusterRegistryConfig.RegistrySources != nil {
			if len(r.Spec.ClusterRegistryConfig.RegistrySources.AllowedRegistries) > 0 && len(r.Spec.ClusterRegistryConfig.RegistrySources.BlockedRegistries) > 0 {
				return field.Invalid(field.NewPath("spec.clusterRegistryConfig.registrySources"), r.Spec.ClusterRegistryConfig.RegistrySources, "allowedRegistries and blockedRegistries are mutually exclusive fields")
			}
		}
	}

	return nil
}

// ValidateUpdate implements admission.Validator.
func (*rosaControlPlaneWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	r, ok := newObj.(*ROSAControlPlane)
	if !ok {
		return nil, fmt.Errorf("expected an ROSAControlPlane object but got %T", r)
	}

	var allErrs field.ErrorList

	if err := r.validateVersion(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateEtcdEncryptionKMSArn(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateRosaRoleConfig(); err != nil {
		allErrs = append(allErrs, err)
	}

	allErrs = append(allErrs, r.validateNetwork()...)
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

// ValidateDelete implements admission.Validator.
func (*rosaControlPlaneWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

func (r *ROSAControlPlane) validateVersion() *field.Error {
	_, err := semver.Parse(r.Spec.Version)
	if err != nil {
		return field.Invalid(field.NewPath("spec.version"), r.Spec.Version, "must be a valid semantic version")
	}

	return nil
}

func (r *ROSAControlPlane) validateNetwork() field.ErrorList {
	var allErrs field.ErrorList
	if r.Spec.Network == nil {
		return allErrs
	}

	rootPath := field.NewPath("spec", "network")

	if r.Spec.Network.MachineCIDR != "" {
		_, _, err := net.ParseCIDR(r.Spec.Network.MachineCIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(rootPath.Child("machineCIDR"), r.Spec.Network.MachineCIDR, "must be valid CIDR block"))
		}
	}

	if r.Spec.Network.PodCIDR != "" {
		_, _, err := net.ParseCIDR(r.Spec.Network.PodCIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(rootPath.Child("podCIDR"), r.Spec.Network.PodCIDR, "must be valid CIDR block"))
		}
	}

	if r.Spec.Network.ServiceCIDR != "" {
		_, _, err := net.ParseCIDR(r.Spec.Network.ServiceCIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(rootPath.Child("serviceCIDR"), r.Spec.Network.ServiceCIDR, "must be valid CIDR block"))
		}
	}

	return allErrs
}

func (r *ROSAControlPlane) validateEtcdEncryptionKMSArn() *field.Error {
	err := kmsArnRegexpValidator.ValidateKMSKeyARN(&r.Spec.EtcdEncryptionKMSARN)
	if err != nil {
		return field.Invalid(field.NewPath("spec.etcdEncryptionKMSARN"), r.Spec.EtcdEncryptionKMSARN, err.Error())
	}

	return nil
}

func (r *ROSAControlPlane) validateExternalAuthProviders() *field.Error {
	if !r.Spec.EnableExternalAuthProviders && len(r.Spec.ExternalAuthProviders) > 0 {
		return field.Invalid(field.NewPath("spec.ExternalAuthProviders"), r.Spec.ExternalAuthProviders,
			"can only be set if spec.EnableExternalAuthProviders is set to 'True'")
	}

	return nil
}

func (r *ROSAControlPlane) validateRosaRoleConfig() *field.Error {
	hasRoleFields := r.Spec.OIDCID != "" || r.Spec.InstallerRoleARN != "" || r.Spec.SupportRoleARN != "" || r.Spec.WorkerRoleARN != "" ||
		r.Spec.RolesRef.IngressARN != "" || r.Spec.RolesRef.ImageRegistryARN != "" || r.Spec.RolesRef.StorageARN != "" ||
		r.Spec.RolesRef.NetworkARN != "" || r.Spec.RolesRef.KubeCloudControllerARN != "" || r.Spec.RolesRef.NodePoolManagementARN != "" ||
		r.Spec.RolesRef.ControlPlaneOperatorARN != "" || r.Spec.RolesRef.KMSProviderARN != ""

	if r.Spec.RosaRoleConfigRef != nil {
		if hasRoleFields {
			return field.Invalid(field.NewPath("spec.rosaRoleConfigRef"), r.Spec.RosaRoleConfigRef, "RosaRoleConfigRef and role fields such as installerRoleARN, supportRoleARN, workerRoleARN, rolesRef and oidcID are mutually exclusive")
		}
		return nil
	}

	if r.Spec.OIDCID == "" {
		return field.Invalid(field.NewPath("spec.oidcID"), r.Spec.OIDCID, "must be specified")
	}
	if r.Spec.InstallerRoleARN == "" {
		return field.Invalid(field.NewPath("spec.installerRoleARN"), r.Spec.InstallerRoleARN, "must be specified")
	}
	if r.Spec.SupportRoleARN == "" {
		return field.Invalid(field.NewPath("spec.supportRoleARN"), r.Spec.SupportRoleARN, "must be specified")
	}
	if r.Spec.WorkerRoleARN == "" {
		return field.Invalid(field.NewPath("spec.workerRoleARN"), r.Spec.WorkerRoleARN, "must be specified")
	}
	if r.Spec.RolesRef.IngressARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.ingressARN"), r.Spec.RolesRef.IngressARN, "must be specified")
	}
	if r.Spec.RolesRef.ImageRegistryARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.imageRegistryARN"), r.Spec.RolesRef.ImageRegistryARN, "must be specified")
	}
	if r.Spec.RolesRef.StorageARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.storageARN"), r.Spec.RolesRef.StorageARN, "must be specified")
	}
	if r.Spec.RolesRef.NetworkARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.networkARN"), r.Spec.RolesRef.NetworkARN, "must be specified")
	}
	if r.Spec.RolesRef.KubeCloudControllerARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.kubeCloudControllerARN"), r.Spec.RolesRef.KubeCloudControllerARN, "must be specified")
	}
	if r.Spec.RolesRef.NodePoolManagementARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.nodePoolManagementARN"), r.Spec.RolesRef.NodePoolManagementARN, "must be specified")
	}
	if r.Spec.RolesRef.ControlPlaneOperatorARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.controlPlaneOperatorARN"), r.Spec.RolesRef.ControlPlaneOperatorARN, "must be specified")
	}
	if r.Spec.RolesRef.KMSProviderARN == "" {
		return field.Invalid(field.NewPath("spec.rolesRef.kmsProviderARN"), r.Spec.RolesRef.KMSProviderARN, "must be specified")
	}
	return nil
}

func (r *ROSAControlPlane) validateROSANetwork() *field.Error {
	if r.Spec.ROSANetworkRef != nil {
		if r.Spec.Subnets != nil {
			return field.Forbidden(field.NewPath("spec.rosaNetworkRef"), "spec.subnets and spec.rosaNetworkRef are mutually exclusive")
		}
		if r.Spec.AvailabilityZones != nil {
			return field.Forbidden(field.NewPath("spec.rosaNetworkRef"), "spec.availabilityZones and spec.rosaNetworkRef are mutually exclusive")
		}
	}

	if r.Spec.ROSANetworkRef == nil && r.Spec.Subnets == nil {
		return field.Required(field.NewPath("spec.subnets"), "spec.subnets cannot be empty when spec.rosaNetworkRef is unspecified")
	}

	if r.Spec.ROSANetworkRef == nil && r.Spec.AvailabilityZones == nil {
		return field.Required(field.NewPath("spec.availabilityZones"), "spec.availabilityZones cannot be empty when spec.rosaNetworkRef is unspecified")
	}

	return nil
}

// Default implements admission.Defaulter.
func (*rosaControlPlaneWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*ROSAControlPlane)
	if !ok {
		return fmt.Errorf("expected an ROSAControlPlane object but got %T", r)
	}

	SetObjectDefaults_ROSAControlPlane(r)
	return nil
}
