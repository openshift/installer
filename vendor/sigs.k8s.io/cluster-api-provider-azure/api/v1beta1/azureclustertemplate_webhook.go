/*
Copyright 2022 The Kubernetes Authors.

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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// AzureClusterTemplateImmutableMsg is the message used for errors on fields that are immutable.
const AzureClusterTemplateImmutableMsg = "AzureClusterTemplate spec.template.spec field is immutable. Please create new resource instead. ref doc: https://cluster-api.sigs.k8s.io/tasks/experimental-features/cluster-class/change-clusterclass.html"

// SetupWebhookWithManager will set up the webhook to be managed by the specified manager.
func (c *AzureClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(azureClusterTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azureclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclustertemplates,versions=v1beta1,name=validation.azureclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azureclustertemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclustertemplates,versions=v1beta1,name=default.azureclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureClusterTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &azureClusterTemplateWebhook{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*azureClusterTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	c, ok := obj.(*AzureClusterTemplate)
	if !ok {
		return fmt.Errorf("expected an AzureClusterTemplate object but got %T", c)
	}

	c.setDefaults()
	return nil
}

var _ webhook.CustomValidator = &azureClusterTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	c, ok := obj.(*AzureClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AzureClusterTemplate object but got %T", c)
	}

	return c.validateClusterTemplate()
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterTemplateWebhook) ValidateUpdate(_ context.Context, oldRaw, newObj runtime.Object) (admission.Warnings, error) {
	c, ok := newObj.(*AzureClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AzureClusterTemplate object but got %T", c)
	}

	var allErrs field.ErrorList
	old := oldRaw.(*AzureClusterTemplate)
	if !reflect.DeepEqual(c.Spec.Template.Spec, old.Spec.Template.Spec) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("AzureClusterTemplate", "spec", "template", "spec"), c, AzureClusterTemplateImmutableMsg),
		)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureClusterTemplateKind).GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*azureClusterTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
