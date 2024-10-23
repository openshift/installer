package v1beta2

import (
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
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-controlplane-cluster-x-k8s-io-v1beta2-rosacontrolplane,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,versions=v1beta2,name=validation.rosacontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-controlplane-cluster-x-k8s-io-v1beta2-rosacontrolplane,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=controlplane.cluster.x-k8s.io,resources=rosacontrolplanes,versions=v1beta2,name=default.rosacontrolplanes.controlplane.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Defaulter = &ROSAControlPlane{}
var _ webhook.Validator = &ROSAControlPlane{}

// ValidateCreate implements admission.Validator.
func (r *ROSAControlPlane) ValidateCreate() (warnings admission.Warnings, err error) {
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
func (r *ROSAControlPlane) ValidateUpdate(old runtime.Object) (warnings admission.Warnings, err error) {
	var allErrs field.ErrorList

	if err := r.validateVersion(); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := r.validateEtcdEncryptionKMSArn(); err != nil {
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
func (r *ROSAControlPlane) ValidateDelete() (warnings admission.Warnings, err error) {
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

// Default implements admission.Defaulter.
func (r *ROSAControlPlane) Default() {
	SetObjectDefaults_ROSAControlPlane(r)
}
