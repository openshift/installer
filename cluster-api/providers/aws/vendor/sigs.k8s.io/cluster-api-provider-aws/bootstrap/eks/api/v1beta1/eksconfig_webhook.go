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
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
)

// SetupWebhookWithManager will setup the webhooks for the EKSConfig.
func (r *EKSConfig) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-bootstrap-cluster-x-k8s-io-v1beta1-eksconfig,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=bootstrap.cluster.x-k8s.io,resources=eksconfig,versions=v1beta1,name=validation.eksconfigs.bootstrap.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-bootstrap-cluster-x-k8s-io-v1beta1-eksconfig,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=bootstrap.cluster.x-k8s.io,resources=eksconfig,versions=v1beta1,name=default.eksconfigs.bootstrap.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

var _ webhook.Defaulter = &EKSConfig{}
var _ webhook.Validator = &EKSConfig{}

// ValidateCreate will do any extra validation when creating a EKSConfig.
func (r *EKSConfig) ValidateCreate() error {
	return nil
}

// ValidateUpdate will do any extra validation when updating a EKSConfig.
func (r *EKSConfig) ValidateUpdate(old runtime.Object) error {
	return nil
}

// ValidateDelete allows you to add any extra validation when deleting.
func (r *EKSConfig) ValidateDelete() error {
	return nil
}

// Default will set default values for the EKSConfig.
func (r *EKSConfig) Default() {
}
