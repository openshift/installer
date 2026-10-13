/*
Copyright 2021 The Kubernetes Authors.

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
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (c *AzureClusterIdentity) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(azureClusterIdentityWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		WithValidator(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azureclusteridentity,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusteridentities,versions=v1beta1,name=validation.azureclusteridentity.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureClusterIdentityWebhook struct{}

var _ webhook.CustomValidator = &azureClusterIdentityWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterIdentityWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	c, ok := obj.(*AzureClusterIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AzureClusterIdentity object but got %T", c)
	}

	return c.validateClusterIdentity()
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterIdentityWebhook) ValidateUpdate(_ context.Context, oldRaw, newObj runtime.Object) (admission.Warnings, error) {
	c, ok := newObj.(*AzureClusterIdentity)
	if !ok {
		return nil, fmt.Errorf("expected an AzureClusterIdentity object but got %T", c)
	}

	var allErrs field.ErrorList
	old := oldRaw.(*AzureClusterIdentity)
	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "Type"),
		old.Spec.Type,
		c.Spec.Type); err != nil {
		allErrs = append(allErrs, err)
	}
	if len(allErrs) == 0 {
		return c.validateClusterIdentity()
	}
	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureClusterIdentityKind).GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterIdentityWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
