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

package webhooks

import (
	"context"
	"reflect"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/apis/meta/v1/unstructured"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	infrav1 "sigs.k8s.io/cluster-api-provider-gcp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmachinetemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpmachinetemplates,versions=v1beta1,name=validation.gcpmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmachinetemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=gcpmachinetemplates,versions=v1beta1,name=default.gcpmachinetemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

func (r *GCPMachineTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.GCPMachineTemplate{}).
		WithValidator(r).
		WithDefaulter(r).
		Complete()
}

// GCPMachineTemplate implements a validating and defaulting webhook for GCPMachineTemplate.
type GCPMachineTemplate struct{}

var (
	_ admission.Validator[*infrav1.GCPMachineTemplate] = &GCPMachineTemplate{}
	_ admission.Defaulter[*infrav1.GCPMachineTemplate] = &GCPMachineTemplate{}
)

func (*GCPMachineTemplate) ValidateCreate(_ context.Context, r *infrav1.GCPMachineTemplate) (admission.Warnings, error) {
	clusterlog.Info("validate create", "name", r.Name)

	return nil, validateConfidentialCompute(r.Spec.Template.Spec)
}

func (*GCPMachineTemplate) ValidateUpdate(_ context.Context, oldObj, r *infrav1.GCPMachineTemplate) (admission.Warnings, error) {
	newGCPMachineTemplate, err := runtime.DefaultUnstructuredConverter.ToUnstructured(r)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachineTemplate").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert new GCPMachineTemplate to unstructured object")),
		})
	}
	oldGCPMachineTemplate, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObj)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachineTemplate").GroupKind(), r.Name, field.ErrorList{
			field.InternalError(nil, errors.Wrap(err, "failed to convert old GCPMachineTemplate to unstructured object")),
		})
	}

	newGCPMachineTemplateSpec := newGCPMachineTemplate["spec"].(map[string]interface{})
	oldGCPMachineTemplateSpec := oldGCPMachineTemplate["spec"].(map[string]interface{})

	// allow changes to providerID, additionalLabels and additionalNetworkTags
	for _, fieldName := range []string{"providerID", "additionalLabels", "additionalNetworkTags"} {
		unstructured.RemoveNestedField(oldGCPMachineTemplateSpec, "template", "spec", fieldName)
		unstructured.RemoveNestedField(newGCPMachineTemplateSpec, "template", "spec", fieldName)
	}

	if !reflect.DeepEqual(oldGCPMachineTemplateSpec, newGCPMachineTemplateSpec) {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind("GCPMachineTemplate").GroupKind(), r.Name, field.ErrorList{
			field.Forbidden(field.NewPath("spec"), "cannot be modified"),
		})
	}

	return nil, nil
}

func (*GCPMachineTemplate) ValidateDelete(_ context.Context, _ *infrav1.GCPMachineTemplate) (admission.Warnings, error) {
	return nil, nil
}

func (*GCPMachineTemplate) Default(_ context.Context, _ *infrav1.GCPMachineTemplate) error {
	return nil
}
