package webhooks

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

// ROSANetwork implements a custom validation webhook for ROSANetwork.
type ROSANetwork struct{}

// SetupWebhookWithManager will setup the webhooks for the ROSANetwork.
func (w *ROSANetwork) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&expinfrav1.ROSANetwork{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosanetwork,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,versions=v1beta2,name=validation.rosanetwork.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosanetwork,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosanetworks,versions=v1beta2,name=default.rosanetwork.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &ROSANetwork{}
var _ webhook.CustomValidator = &ROSANetwork{}

// ValidateCreate implements admission.Validator.
func (w *ROSANetwork) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	rosaNet, ok := obj.(*expinfrav1.ROSANetwork)
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
func (w *ROSANetwork) ValidateUpdate(ctx context.Context, old runtime.Object, updated runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements admission.Validator.
func (w *ROSANetwork) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// Default implements admission.Defaulter.
func (w *ROSANetwork) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}
