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

package webhooks

import (
	"context"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

// log is for logging in this package.
var gmctlog = logf.Log.WithName("gcpclustertemplate-resource")

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (r *GCPManagedClusterTemplate) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &expinfrav1.GCPManagedClusterTemplate{}).
		WithValidator(r).
		Complete()
}

// GCPManagedClusterTemplate implements a validating webhook for GCPManagedClusterTemplate.
type GCPManagedClusterTemplate struct{}

//+kubebuilder:webhook:verbs=update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedclustertemplate,mutating=false,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedclustertemplates,versions=v1beta1,name=vgcpmanagedclustertemplate.kb.io,admissionReviewVersions=v1

var _ admission.Validator[*expinfrav1.GCPManagedClusterTemplate] = &GCPManagedClusterTemplate{}

func (*GCPManagedClusterTemplate) ValidateCreate(_ context.Context, r *expinfrav1.GCPManagedClusterTemplate) (admission.Warnings, error) {
	gmctlog.Info("Validating GCPManagedClusterTemplate create", "name", r.Name)

	return nil, nil
}

func (*GCPManagedClusterTemplate) ValidateUpdate(_ context.Context, old, r *expinfrav1.GCPManagedClusterTemplate) (admission.Warnings, error) {
	gmctlog.Info("Validating GCPManagedClusterTemplate update", "name", r.Name)

	if !reflect.DeepEqual(r.Spec, old.Spec) {
		return nil, apierrors.NewBadRequest("GCPManagedClusterTemplate.Spec is immutable")
	}

	return nil, nil
}

func (*GCPManagedClusterTemplate) ValidateDelete(_ context.Context, r *expinfrav1.GCPManagedClusterTemplate) (admission.Warnings, error) {
	gmctlog.Info("Validint GCPManagedClusterTemplate delete", "name", r.Name)

	return nil, nil
}
