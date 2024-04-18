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
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	capifeature "sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *AzureManagedCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcluster,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedclusters,versions=v1beta1,name=validation.azuremanagedclusters.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

var _ webhook.Validator = &AzureManagedCluster{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *AzureManagedCluster) ValidateCreate() (admission.Warnings, error) {
	// NOTE: AzureManagedCluster relies upon MachinePools, which is behind a feature gate flag.
	// The webhook must prevent creating new objects in case the feature flag is disabled.
	if !feature.Gates.Enabled(capifeature.MachinePool) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			"can be set only if the Cluster API 'MachinePool' feature flag is enabled",
		)
	}
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *AzureManagedCluster) ValidateUpdate(oldRaw runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *AzureManagedCluster) ValidateDelete() (admission.Warnings, error) {
	return nil, nil
}
