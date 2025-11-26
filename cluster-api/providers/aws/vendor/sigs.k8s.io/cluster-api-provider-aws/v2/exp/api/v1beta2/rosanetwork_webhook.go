package v1beta2

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager will setup the webhooks for the ROSANetwork.
func (r *ROSANetwork) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(rosaNetworkWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosanetwork,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,versions=v1beta2,name=validation.rosanetwork.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosanetwork,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,versions=v1beta2,name=default.rosanetwork.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type rosaNetworkWebhook struct{}

var _ webhook.CustomDefaulter = &rosaNetworkWebhook{}
var _ webhook.CustomValidator = &rosaNetworkWebhook{}

// ValidateCreate implements admission.Validator.
func (r *rosaNetworkWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	rosaNet, ok := obj.(*ROSANetwork)
	if !ok {
		return nil, fmt.Errorf("expected an ROSANetwork object but got %T", rosaNet)
	}

	var allErrs field.ErrorList
	if rosaNet.Spec.AvailabilityZoneCount == 0 && len(rosaNet.Spec.AvailabilityZones) == 0 {
		err := field.Invalid(field.NewPath("spec.AvailabilityZones"), rosaNet.Spec.AvailabilityZones, "Either AvailabilityZones OR AvailabilityZoneCount must be set.")
		allErrs = append(allErrs, err)
	}
	if rosaNet.Spec.AvailabilityZoneCount != 0 && len(rosaNet.Spec.AvailabilityZones) > 0 {
		err := field.Invalid(field.NewPath("spec.AvailabilityZones"), rosaNet.Spec.AvailabilityZones, "Either AvailabilityZones OR AvailabilityZoneCount can be set.")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(
			rosaNet.GroupVersionKind().GroupKind(),
			rosaNet.Name,
			allErrs)
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator.
func (r *rosaNetworkWebhook) ValidateUpdate(ctx context.Context, old runtime.Object, updated runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements admission.Validator.
func (r *rosaNetworkWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// Default implements admission.Defaulter.
func (r *rosaNetworkWebhook) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}
