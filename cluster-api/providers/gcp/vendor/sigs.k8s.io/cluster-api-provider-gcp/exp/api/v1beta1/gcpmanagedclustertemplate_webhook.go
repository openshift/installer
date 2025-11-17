/*
Copyright 2025 The Kubernetes Authors.

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
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var gmctlog = logf.Log.WithName("gcpclustertemplate-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *GCPManagedClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mctw := new(gcpManagedClusterTemplateWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(mctw).
		Complete()
}

type gcpManagedClusterTemplateWebhook struct{}

//+kubebuilder:webhook:verbs=update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedclustertemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedclustertemplates,versions=v1beta1,name=vgcpmanagedclustertemplate.kb.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &gcpManagedClusterTemplateWebhook{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedClusterTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*GCPManagedClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPManagedClusterTemplate but got a %T", obj))
	}

	gmctlog.Info("Validating GCPManagedClusterTemplate create", "name", r.Name)

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedClusterTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	old, ok := oldObj.(*GCPManagedClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPManagedClusterTemplate but got a %T", oldObj))
	}

	r, ok := newObj.(*GCPManagedClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPManagedClusterTemplate but got a %T", newObj))
	}

	gmctlog.Info("Validating GCPManagedClusterTemplate update", "name", r.Name)

	if !reflect.DeepEqual(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("GCPManagedClusterTemplate.Spec is immutable")
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*gcpManagedClusterTemplateWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*GCPManagedClusterTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPManagedClusterTemplate but got a %T", obj))
	}

	gmctlog.Info("Validint GCPManagedClusterTemplate delete", "name", r.Name)

	return nil, nil
}
