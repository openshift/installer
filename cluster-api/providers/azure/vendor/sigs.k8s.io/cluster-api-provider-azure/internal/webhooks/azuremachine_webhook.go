/*
Copyright The Kubernetes Authors.

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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	apiinternal "sigs.k8s.io/cluster-api-provider-azure/internal/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (mw *AzureMachineWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mw.client = mgr.GetClient()

	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureMachine{}).
		WithDefaulter(mw).
		WithValidator(mw).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremachine,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremachines,versions=v1beta1,name=validation.azuremachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremachine,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremachines,versions=v1beta1,name=default.azuremachine.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureMachineWebhook implements a validating and defaulting webhook for AzureMachines.
type AzureMachineWebhook struct {
	client client.Client
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureMachineWebhook) ValidateCreate(_ context.Context, m *infrav1.AzureMachine) (admission.Warnings, error) {
	spec := m.Spec

	allErrs := validateAzureMachineSpec(spec)

	roleAssignmentName := ""
	if spec.SystemAssignedIdentityRole != nil {
		roleAssignmentName = spec.SystemAssignedIdentityRole.Name
	}

	if errs := ValidateSystemAssignedIdentity(spec.Identity, "", roleAssignmentName, field.NewPath("roleAssignmentName")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureMachineKind).GroupKind(), m.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureMachineWebhook) ValidateUpdate(_ context.Context, old, m *infrav1.AzureMachine) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "image"),
		old.Spec.Image,
		m.Spec.Image); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "identity"),
		old.Spec.Identity,
		m.Spec.Identity); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "systemAssignedIdentityRole"),
		old.Spec.SystemAssignedIdentityRole,
		m.Spec.SystemAssignedIdentityRole); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "userAssignedIdentities"),
		old.Spec.UserAssignedIdentities,
		m.Spec.UserAssignedIdentities); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "roleAssignmentName"),
		old.Spec.RoleAssignmentName,             //nolint:staticcheck
		m.Spec.RoleAssignmentName); err != nil { //nolint:staticcheck
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "osDisk"),
		old.Spec.OSDisk,
		m.Spec.OSDisk); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "dataDisks"),
		old.Spec.DataDisks,
		m.Spec.DataDisks); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "sshPublicKey"),
		old.Spec.SSHPublicKey,
		m.Spec.SSHPublicKey); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "allocatePublicIP"),
		old.Spec.AllocatePublicIP,
		m.Spec.AllocatePublicIP); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "enableIPForwarding"),
		old.Spec.EnableIPForwarding,
		m.Spec.EnableIPForwarding); err != nil {
		allErrs = append(allErrs, err)
	}

	// Spec.AcceleratedNetworking can only be reset to nil and no other changes apart from that
	// is accepted if the field is set.
	// Ref issue #3518
	if err := webhookutils.ValidateZeroTransition(
		field.NewPath("Spec", "AcceleratedNetworking"),
		old.Spec.AcceleratedNetworking,             //nolint:staticcheck
		m.Spec.AcceleratedNetworking); err != nil { //nolint:staticcheck
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SpotVMOptions"),
		old.Spec.SpotVMOptions,
		m.Spec.SpotVMOptions); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SecurityProfile"),
		old.Spec.SecurityProfile,
		m.Spec.SecurityProfile); err != nil {
		allErrs = append(allErrs, err)
	}

	if old.Spec.Diagnostics != nil {
		if err := webhookutils.ValidateImmutable(
			field.NewPath("spec", "diagnostics"),
			old.Spec.Diagnostics,
			m.Spec.Diagnostics); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	if !reflect.DeepEqual(m.Spec.NetworkInterfaces, old.Spec.NetworkInterfaces) {
		// The defaulting webhook may have migrated values from the old SubnetName field to the new NetworkInterfaces format.
		apiinternal.SetDefaultAzureMachineSpecNetworkInterfaces(&old.Spec)

		// The reconciler will populate the SubnetName on the first interface if the user left it blank.
		if old.Spec.NetworkInterfaces[0].SubnetName == "" && m.Spec.NetworkInterfaces[0].SubnetName != "" {
			old.Spec.NetworkInterfaces[0].SubnetName = m.Spec.NetworkInterfaces[0].SubnetName
		}

		// Enforce immutability for all other changes to NetworkInterfaces.
		if !reflect.DeepEqual(m.Spec.NetworkInterfaces, old.Spec.NetworkInterfaces) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "networkInterfaces"),
					m.Spec.NetworkInterfaces, "field is immutable"),
			)
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "capacityReservationGroupID"),
		old.Spec.CapacityReservationGroupID,
		m.Spec.CapacityReservationGroupID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "disableExtensionOperations"),
		old.Spec.DisableExtensionOperations,
		m.Spec.DisableExtensionOperations); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "disableVMBootstrapExtension"),
		old.Spec.DisableVMBootstrapExtension,
		m.Spec.DisableVMBootstrapExtension); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureMachineKind).GroupKind(), m.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureMachineWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureMachine) (admission.Warnings, error) {
	return nil, nil
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mw *AzureMachineWebhook) Default(_ context.Context, m *infrav1.AzureMachine) error {
	return apiinternal.SetDefaultsAzureMachine(m, mw.client)
}
