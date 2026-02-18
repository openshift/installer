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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var gcpMachinePoolLog = logf.Log.WithName("gcpmachinepool-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *GCPMachinePool) SetupWebhookWithManager(mgr ctrl.Manager) error {
	validator := new(gcpMachinePoolWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(r).
		WithValidator(validator).
		Complete()
}

type gcpMachinePoolWebhook struct{}

//+kubebuilder:webhook:verbs=update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmachinepool,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmachinepools,versions=v1beta1,name=validation.gcpmachinepool.infrastructure.cluster.x-k8s.io,admissionReviewVersions=v1

var _ webhook.CustomValidator = &gcpMachinePoolWebhook{}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpMachinePoolWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*GCPMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPMachinePool but got a %T", obj))
	}

	gcpMachinePoolLog.Info("Validating GCPMachinePool create", "name", r.Name)

	// Add custom validation logic upon creation if needed.

	return nil, nil
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*gcpMachinePoolWebhook) ValidateUpdate(_ context.Context, _, newObj runtime.Object) (admission.Warnings, error) {
	r, ok := newObj.(*GCPMachinePool)
	if !ok {
		return nil, fmt.Errorf("expected a GCPMachinePool object but got %T", r)
	}

	gcpMachinePoolLog.Info("Validating GCPMachinePool update", "name", r.Name)

	// Add custom validation logic upon update if needed.

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*gcpMachinePoolWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*GCPMachinePool)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a GCPMachinePool but got a %T", obj))
	}

	gcpMachinePoolLog.Info("Validating GCPMachinePool delete", "name", r.Name)

	// Add custom validation logic upon deletion if needed.

	return nil, nil
}
