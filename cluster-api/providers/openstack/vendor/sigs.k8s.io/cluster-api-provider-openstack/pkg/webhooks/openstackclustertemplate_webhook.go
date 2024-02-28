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

package webhooks

import (
	"context"
	"fmt"
	"reflect"

	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-openstackclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclustertemplates,versions=v1beta1,name=validation.openstackclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

func SetupOpenStackClusterTemplateWebhook(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.OpenStackClusterTemplate{}).
		WithValidator(&openStackClusterTemplateWebhook{}).
		Complete()
}

type openStackClusterTemplateWebhook struct{}

// Compile-time assertion that openStackClusterTemplateWebhook implements webhook.CustomValidator.
var _ webhook.CustomValidator = &openStackClusterTemplateWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateCreate(_ context.Context, objRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	newObj, err := castToOpenStackClusterTemplate(objRaw)
	if err != nil {
		return nil, err
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateUpdate(_ context.Context, oldObjRaw, newObjRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	oldObj, err := castToOpenStackClusterTemplate(oldObjRaw)
	if err != nil {
		return nil, err
	}
	newObj, err := castToOpenStackClusterTemplate(newObjRaw)
	if err != nil {
		return nil, err
	}

	if !reflect.DeepEqual(newObj.Spec.Template.Spec, oldObj.Spec.Template.Spec) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("OpenStackClusterTemplate", "spec", "template", "spec"), newObj, "OpenStackClusterTemplate spec.template.spec field is immutable. Please create new resource instead."),
		)
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func castToOpenStackClusterTemplate(obj runtime.Object) (*infrav1.OpenStackClusterTemplate, error) {
	cast, ok := obj.(*infrav1.OpenStackClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an OpenStackClusterTemplate but got a %T", obj)
	}
	return cast, nil
}
