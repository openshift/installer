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

	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
	"sigs.k8s.io/cluster-api-provider-vsphere/internal/webhooks/conversion"
)

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta2-vspheredeploymentzone,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones,versions=v1beta2,name=default.vspheredeploymentzone.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1

// VSphereDeploymentZone implements a defaulting webhook for VSphereDeploymentZone.
type VSphereDeploymentZone struct{}

var _ admission.Defaulter[*infrav1.VSphereDeploymentZone] = &VSphereDeploymentZone{}

func (webhook *VSphereDeploymentZone) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.VSphereDeploymentZone{}).
		WithDefaulter(webhook).
		WithConverter(conversion.VSphereDeploymentZone).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereDeploymentZone) Default(_ context.Context, typedObj *infrav1.VSphereDeploymentZone) error {
	if typedObj.Spec.ControlPlane == nil {
		typedObj.Spec.ControlPlane = ptr.To(true)
	}
	return nil
}
