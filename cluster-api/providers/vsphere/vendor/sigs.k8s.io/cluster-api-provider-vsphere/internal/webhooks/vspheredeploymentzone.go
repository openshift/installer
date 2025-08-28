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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-vspheredeploymentzone,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=vspheredeploymentzones,versions=v1beta1,name=default.vspheredeploymentzone.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereDeploymentZoneWebhook implements a validation and defaulting webhook for VSphereDeploymentZone.
type VSphereDeploymentZoneWebhook struct{}

var _ webhook.CustomDefaulter = &VSphereDeploymentZoneWebhook{}

func (webhook *VSphereDeploymentZoneWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&infrav1.VSphereDeploymentZone{}).
		WithDefaulter(webhook).
		Complete()
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (webhook *VSphereDeploymentZoneWebhook) Default(_ context.Context, obj runtime.Object) error {
	typedObj, ok := obj.(*infrav1.VSphereDeploymentZone)
	if !ok {
		return apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereDeploymentZone but got a %T", obj))
	}
	if typedObj.Spec.ControlPlane == nil {
		typedObj.Spec.ControlPlane = ptr.To(true)
	}

	return nil
}
