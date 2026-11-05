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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (w *AzureClusterIdentityWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureClusterIdentity{}).
		WithValidator(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azureclusteridentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusteridentities,versions=v1beta1,name=validation.azureclusteridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureClusterIdentityWebhook implements a validating webhook for AzureClusterIdentity.
type AzureClusterIdentityWebhook struct{}

var _ admission.Validator[*infrav1.AzureClusterIdentity] = &AzureClusterIdentityWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterIdentityWebhook) ValidateCreate(_ context.Context, c *infrav1.AzureClusterIdentity) (admission.Warnings, error) {
	return validateAzureClusterIdentity(c)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterIdentityWebhook) ValidateUpdate(_ context.Context, old, c *infrav1.AzureClusterIdentity) (admission.Warnings, error) {
	var allErrs field.ErrorList
	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "Type"),
		old.Spec.Type,
		c.Spec.Type); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return validateAzureClusterIdentity(c)
	}
	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureClusterIdentityKind).GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterIdentityWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureClusterIdentity) (admission.Warnings, error) {
	return nil, nil
}
