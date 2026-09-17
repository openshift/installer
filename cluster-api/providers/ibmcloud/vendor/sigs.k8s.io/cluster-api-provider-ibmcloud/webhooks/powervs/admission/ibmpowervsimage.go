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

package powervs

import (
	"context"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

// Ensure IBMPowerVSImage implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMPowerVSImage] = &IBMPowerVSImage{}
	_ admission.Defaulter[*infrav1.IBMPowerVSImage] = &IBMPowerVSImage{}
)

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsimage,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,versions=v1beta3,name=mibmpowervsimage.kb.io,sideEffects=None,admissionReviewVersions=v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsimage,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsimages,versions=v1beta3,name=vibmpowervsimage.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *IBMPowerVSImage) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMPowerVSImage{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSImage implements a validation and defaulting webhook for IBMPowerVSImage.
type IBMPowerVSImage struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSImage) Default(_ context.Context, _ *infrav1.IBMPowerVSImage) error {
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateCreate(_ context.Context, _ *infrav1.IBMPowerVSImage) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateUpdate(_ context.Context, _, _ *infrav1.IBMPowerVSImage) (warnings admission.Warnings, err error) {
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSImage) ValidateDelete(_ context.Context, _ *infrav1.IBMPowerVSImage) (admission.Warnings, error) {
	return nil, nil
}
