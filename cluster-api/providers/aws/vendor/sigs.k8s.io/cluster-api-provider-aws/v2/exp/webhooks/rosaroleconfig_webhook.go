package webhooks

import (
	"context"
	"fmt"

	"github.com/blang/semver"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	runtime "k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	expinfrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/exp/api/v1beta2"
)

// ROSARoleConfig implements a custom validation webhook for ROSARoleConfig.
type ROSARoleConfig struct{}

// SetupWebhookWithManager will setup the webhooks for the ROSARoleConfig.
func (w *ROSARoleConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&expinfrav1.ROSARoleConfig{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=validation.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-rosaroleconfig,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=rosaroleconfigs,versions=v1beta2,name=default.rosaroleconfig.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.CustomDefaulter = &ROSARoleConfig{}
var _ webhook.CustomValidator = &ROSARoleConfig{}

// ValidateCreate implements admission.Validator.
func (w *ROSARoleConfig) ValidateCreate(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	roleConfig, ok := obj.(*expinfrav1.ROSARoleConfig)
	if !ok {
		return nil, fmt.Errorf("expected an ROSARoleConfig object but got %T", roleConfig)
	}

	var allErrs field.ErrorList
	if roleConfig.Spec.OidcProviderType == expinfrav1.Managed && roleConfig.Spec.OperatorRoleConfig.OIDCID != "" {
		err := field.Invalid(field.NewPath("spec.operatorRoleConfig.oidcId"), roleConfig.Spec.OperatorRoleConfig.OIDCID, "cannot be set with Managed oidc provider type")
		allErrs = append(allErrs, err)
	} else if roleConfig.Spec.OidcProviderType == expinfrav1.Unmanaged && roleConfig.Spec.OperatorRoleConfig.OIDCID == "" {
		err := field.Invalid(field.NewPath("spec.operatorRoleConfig.oidcId"), roleConfig.Spec.OperatorRoleConfig.OIDCID, "must set operatorRoleConfig.oidcId with UnManaged oidc provider type")
		allErrs = append(allErrs, err)
	}

	_, vErr := semver.Parse(roleConfig.Spec.AccountRoleConfig.Version)
	if vErr != nil {
		err := field.Invalid(field.NewPath("spec.accountRoleConfig.version"), roleConfig.Spec.AccountRoleConfig.Version, "must be a valid semantic version")
		allErrs = append(allErrs, err)
	}

	if len(allErrs) > 0 {
		return nil, apierrors.NewInvalid(
			roleConfig.GroupVersionKind().GroupKind(),
			roleConfig.Name,
			allErrs)
	}

	return nil, nil
}

// ValidateUpdate implements admission.Validator.
func (w *ROSARoleConfig) ValidateUpdate(ctx context.Context, old runtime.Object, updated runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements admission.Validator.
func (w *ROSARoleConfig) ValidateDelete(ctx context.Context, obj runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// Default implements admission.Defaulter.
func (w *ROSARoleConfig) Default(ctx context.Context, obj runtime.Object) error {
	return nil
}
