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
	"reflect"

	"k8s.io/apimachinery/pkg/util/validation/field"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-openstackclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackclustertemplates,versions=v1beta2,name=validation.openstackclustertemplate.v1beta2.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

func SetupOpenStackClusterTemplateWebhook(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr, &infrav1.OpenStackClusterTemplate{}).
		WithValidator(&openStackClusterTemplateWebhook{}).
		Complete()
}

type openStackClusterTemplateWebhook struct{}

var _ admission.Validator[*infrav1.OpenStackClusterTemplate] = &openStackClusterTemplateWebhook{}

// ValidateCreate implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateCreate(_ context.Context, newObj *infrav1.OpenStackClusterTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateUpdate implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj *infrav1.OpenStackClusterTemplate) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if !reflect.DeepEqual(newObj.Spec.Template.Spec, oldObj.Spec.Template.Spec) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("OpenStackClusterTemplate", "spec", "template", "spec"), newObj, "OpenStackClusterTemplate spec.template.spec field is immutable. Please create new resource instead."),
		)
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements admission.Validator so a webhook will be registered for the type.
func (*openStackClusterTemplateWebhook) ValidateDelete(_ context.Context, _ *infrav1.OpenStackClusterTemplate) (admission.Warnings, error) {
	return nil, nil
}
