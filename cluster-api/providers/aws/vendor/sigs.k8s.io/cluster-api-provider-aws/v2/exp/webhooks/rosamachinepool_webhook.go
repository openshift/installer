package webhooks

import (
	"context"
	"fmt"

	"github.com/blang/semver"
	"github.com/google/go-cmp/cmp"
	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

// ROSAMachinePool implements a custom validation webhook for ROSAMachinePool.
type ROSAMachinePool struct{}

// SetupWebhookWithManager will setup the webhooks for the ROSAMachinePool.
func (w *ROSAMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&expinfrav1.ROSAMachinePool{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosamachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools,versions=v1beta2,name=validation.rosamachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosamachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosamachinepools,versions=v1beta2,name=default.rosamachinepool.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &ROSAMachinePool{}
var _ webhook.CustomValidator = &ROSAMachinePool{}

// ValidateCreate implements admission.Validator.
func (w *ROSAMachinePool) ValidateCreate(_ context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	r, ok := obj.(*expinfrav1.ROSAMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an ROSAMachinePool object but got %T", r)
	}

	var allErrs field.ErrorList

	if err := w.validateVersion(r); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := w.validateNodeDrainGracePeriod(r); err != nil {
		allErrs = append(allErrs, err)
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

// ValidateUpdate implements admission.Validator.
func (w *ROSAMachinePool) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (warnings admission.Warnings, err error) {
	r, ok := newObj.(*expinfrav1.ROSAMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected an ROSAMachinePool object but got %T", r)
	}

	oldPool, ok := oldObj.(*expinfrav1.ROSAMachinePool)
	if !ok {
		return nil, apierrors.NewInvalid(expinfrav1.GroupVersion.WithKind("ROSAMachinePool").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.New("failed to convert old ROSAMachinePool to object")),
		})
	}

	var allErrs field.ErrorList
	if err := w.validateVersion(r); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := w.validateNodeDrainGracePeriod(r); err != nil {
		allErrs = append(allErrs, err)
	}

	allErrs = append(allErrs, validateImmutable(oldPool.Spec.AdditionalSecurityGroups, r.Spec.AdditionalSecurityGroups, "additionalSecurityGroups")...)
	allErrs = append(allErrs, validateImmutable(oldPool.Spec.AdditionalTags, r.Spec.AdditionalTags, "additionalTags")...)

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
func (w *ROSAMachinePool) ValidateDelete(_ context.Context, _ runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

func (w *ROSAMachinePool) validateVersion(r *expinfrav1.ROSAMachinePool) *field.Error {
	if r.Spec.Version == "" {
		return nil
	}
	_, err := semver.Parse(r.Spec.Version)
	if err != nil {
		return field.Invalid(field.NewPath("spec.version"), r.Spec.Version, "must be a valid semantic version")
	}

	return nil
}

func (w *ROSAMachinePool) validateNodeDrainGracePeriod(r *expinfrav1.ROSAMachinePool) *field.Error {
	if r.Spec.NodeDrainGracePeriod == nil {
		return nil
	}

	if r.Spec.NodeDrainGracePeriod.Minutes() > 10080 {
		return field.Invalid(field.NewPath("spec.nodeDrainGracePeriod"), r.Spec.NodeDrainGracePeriod,
			"max supported duration is 1 week (10080m|168h)")
	}

	return nil
}

func validateImmutable(old, updated interface{}, name string) field.ErrorList {
	var allErrs field.ErrorList

	if !cmp.Equal(old, updated) {
		allErrs = append(
			allErrs,
			field.Invalid(field.NewPath("spec", name), updated, "field is immutable"),
		)
	}

	return allErrs
}

// Default implements admission.Defaulter.
func (w *ROSAMachinePool) Default(ctx context.Context, obj runtime.Object) error {
	r, ok := obj.(*expinfrav1.ROSAMachinePool)
	if !ok {
		return fmt.Errorf("expected an ROSAMachinePool object but got %T", r)
	}

	r.Default()
	return nil
}
