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

package powervs

import (
	"context"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"

	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

// Ensure IBMPowerVSClusterTemplate implements the typed webhook interfaces.
var (
	_ admission.Validator[*infrav1.IBMPowerVSClusterTemplate] = &IBMPowerVSClusterTemplate{}
	_ admission.Defaulter[*infrav1.IBMPowerVSClusterTemplate] = &IBMPowerVSClusterTemplate{}
)

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsclustertemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclustertemplates,versions=v1beta3,name=mibmpowervsclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1
//+kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta3-ibmpowervsclustertemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=ibmpowervsclustertemplates,versions=v1beta3,name=vibmpowervsclustertemplate.kb.io,sideEffects=None,admissionReviewVersions=v1

func (r *IBMPowerVSClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.IBMPowerVSClusterTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// IBMPowerVSClusterTemplate implements a validation and defaulting webhook for IBMPowerVSClusterTemplate.
type IBMPowerVSClusterTemplate struct{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (r *IBMPowerVSClusterTemplate) Default(_ context.Context, _ *infrav1.IBMPowerVSClusterTemplate) error {
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSClusterTemplate) ValidateCreate(_ context.Context, _ *infrav1.IBMPowerVSClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSClusterTemplate) ValidateUpdate(_ context.Context, oldObj, newObj *infrav1.IBMPowerVSClusterTemplate) (warnings admission.Warnings, err error) {
	if !reflect.DeepEqual(newObj.Spec, oldObj.Spec) {
		return nil, apierrors.NewBadRequest("IBMPowerVSClusterTemplate.Spec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (r *IBMPowerVSClusterTemplate) ValidateDelete(_ context.Context, _ *infrav1.IBMPowerVSClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}
