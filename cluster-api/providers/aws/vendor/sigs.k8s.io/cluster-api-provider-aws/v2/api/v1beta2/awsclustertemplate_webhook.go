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

package v1beta2

import (
	"context"
	"fmt"

	"github.com/google/go-cmp/cmp"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

func (r *AWSClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(awsClusterTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta2-awsclustertemplate,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustertemplates,versions=v1beta2,name=validation.awsclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-awsclustertemplate,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=awsclustertemplates,versions=v1beta2,name=default.awsclustertemplate.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type awsClusterTemplateWebhook struct{}

var _ webhook.CustomDefaulter = &awsClusterTemplateWebhook{}
var _ webhook.CustomValidator = &awsClusterTemplateWebhook{}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (*awsClusterTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*AWSClusterTemplate)
	if !ok {
		return fmt.Errorf("expected an AWSClusterTemplate object but got %T", r)
	}

	SetObjectDefaults_AWSClusterTemplate(r)
	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*awsClusterTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*AWSClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterTemplate object but got %T", r)
	}

	var allErrs field.ErrorList

	allErrs = append(allErrs, r.Spec.Template.Spec.Bastion.Validate()...)
	allErrs = append(allErrs, validateSSHKeyName(r.Spec.Template.Spec.SSHKeyName)...)

	return nil, aggregateObjErrors(r.GroupVersionKind().GroupKind(), r.Name, allErrs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*awsClusterTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*AWSClusterTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an AWSClusterTemplate object but got %T", r)
	}

	old := oldObj.(*AWSClusterTemplate)

	if !cmp.Equal(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("AWSClusterTemplate.Spec is immutable")
	}
	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*awsClusterTemplateWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
