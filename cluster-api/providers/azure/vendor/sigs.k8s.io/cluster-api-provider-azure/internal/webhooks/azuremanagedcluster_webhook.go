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

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (w *AzureManagedClusterWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureManagedCluster{}).
		WithValidator(w).
		Complete()
}

// AzureManagedClusterWebhook implements a validating webhook for AzureManagedCluster.
type AzureManagedClusterWebhook struct{}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedclusters,versions=v1beta1,name=validation.azuremanagedclusters.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ admission.Validator[*infrav1.AzureManagedCluster] = &AzureManagedClusterWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (w *AzureManagedClusterWebhook) ValidateCreate(_ context.Context, _ *infrav1.AzureManagedCluster) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (w *AzureManagedClusterWebhook) ValidateUpdate(_ context.Context, _, _ *infrav1.AzureManagedCluster) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (w *AzureManagedClusterWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureManagedCluster) (admission.Warnings, error) {
	return nil, nil
}
