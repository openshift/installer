/*
Copyright 2025 The Kubernetes Authors.

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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-azure/feature"
)

// SetupAzureASOManagedControlPlaneWebhookWithManager sets up and registers the webhook with the manager.
func SetupAzureASOManagedControlPlaneWebhookWithManager(mgr ctrl.Manager) error {
	azureASOManagedControlPlaneWebhook := &azureASOManagedControlPlaneWebhook{}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AzureASOManagedControlPlane{}).
		WithValidator(azureASOManagedControlPlaneWebhook).
		Complete()
}

// azureASOManagedControlPlaneWebhook implements a validating webhook for AzureASOManagedControlPlane.
type azureASOManagedControlPlaneWebhook struct {
}

// +kubebuilder:webhook:verbs=create,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azureasomanagedcontrolplane,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azureasomanagedcontrolplanes,versions=v1beta1,name=validation.azureasomanagedcontrolplane.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (ampw *azureASOManagedControlPlaneWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	_, ok := obj.(*AzureASOManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureASOManagedControlPlane")
	}
	if !feature.Gates.Enabled(feature.ASOAPI) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			fmt.Sprintf("can be set only if the %s feature flag is enabled", feature.ASOAPI),
		)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (ampw *azureASOManagedControlPlaneWebhook) ValidateUpdate(_ context.Context, _, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (ampw *azureASOManagedControlPlaneWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
