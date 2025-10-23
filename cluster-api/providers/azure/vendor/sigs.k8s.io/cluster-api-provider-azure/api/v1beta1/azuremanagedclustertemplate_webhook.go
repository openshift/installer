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

package v1beta1

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *AzureManagedClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(azureManagedClusterTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		Complete()
}

// +kubebuilder:webhook:verbs=update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedclustertemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedclustertemplates,versions=v1beta1,name=validation.azuremanagedclustertemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureManagedClusterTemplateWebhook struct{}

var _ webhook.CustomValidator = &azureManagedClusterTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *azureManagedClusterTemplateWebhook) ValidateCreate(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *azureManagedClusterTemplateWebhook) ValidateUpdate(_ context.Context, _, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *azureManagedClusterTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
