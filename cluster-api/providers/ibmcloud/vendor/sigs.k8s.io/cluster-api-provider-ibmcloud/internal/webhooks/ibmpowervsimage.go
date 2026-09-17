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

package webhooks

import (
	"context"

	"k8s.io/apimachinery/pkg/runtime"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/v1beta2"
)

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsimage,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,verbs=create;update,versions=v1beta2,name=mibmpowervsimage.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-ibmpowervsimage,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,versions=v1beta2,name=vibmpowervsimage.kb.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

func (r *IBMPowerVSImage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.IBMPowerVSImage{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSImage implements a validation and defaulting webhook for IBMPowerVSImage.
type IBMPowerVSImage struct{}

var _ webhook.CustomDefaulter = &IBMPowerVSImage{}
var _ webhook.CustomValidator = &IBMPowerVSImage{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (r *IBMPowerVSImage) Default(_ context.Context, _ runtime.Object) error {
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateCreate(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateUpdate(_ context.Context, _, _ runtime.Object) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
