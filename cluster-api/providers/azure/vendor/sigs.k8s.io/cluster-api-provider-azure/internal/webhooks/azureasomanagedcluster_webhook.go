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
	"fmt"

	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (ampw *AzureASOManagedClusterWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureASOManagedCluster{}).
		WithValidator(ampw).
		Complete()
}

// AzureASOManagedClusterWebhook implements a validating and defaulting webhook for AzureASOManagedCluster.
type AzureASOManagedClusterWebhook struct {
}

// +kubebuilder:webhook:verbs=create,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azureasomanagedcluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedclusters,versions=v1beta1,name=validation.azureasomanagedcluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (ampw *AzureASOManagedClusterWebhook) ValidateCreate(_ context.Context, _ *infrav1.AzureASOManagedCluster) (admission.Warnings, error) {
	if !feature.Gates.Enabled(feature.ASOAPI) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			fmt.Sprintf("can be set only if the %s feature flag is enabled", feature.ASOAPI),
		)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (ampw *AzureASOManagedClusterWebhook) ValidateUpdate(_ context.Context, _, _ *infrav1.AzureASOManagedCluster) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (ampw *AzureASOManagedClusterWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureASOManagedCluster) (admission.Warnings, error) {
	return nil, nil
}
