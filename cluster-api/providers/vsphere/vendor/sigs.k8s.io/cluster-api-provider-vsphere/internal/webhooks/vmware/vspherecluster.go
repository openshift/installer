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

// Package vmware is the package for webhooks of vmware resources.
package vmware

import (
	"context"
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	vmwarev1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/vmware/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/feature"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/manager"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-vmware-infrastructure-cluster-x-k8s-io-v1beta1-vspherecluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=vmware.infrastructure.cluster.x-k8s.io,resources=vsphereclusters,versions=v1beta1,name=validation.vspherecluster.vmware.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

// VSphereCluster implements a validation and defaulting webhook for VSphereCluster.
type VSphereCluster struct {
	// NetworkProvider is the network provider used by Supervisor based clusters
	NetworkProvider string
}

var _ webhook.CustomValidator = &VSphereCluster{}

func (webhook *VSphereCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr).
		For(&vmwarev1.VSphereCluster{}).
		WithValidator(webhook).
		Complete()
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereCluster) ValidateCreate(_ context.Context, objRaw runtime.Object) (admission.Warnings, error) {
	obj, ok := objRaw.(*vmwarev1.VSphereCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereCluster but got a %T", objRaw))
	}
	return webhook.validateClusterNetwork(obj)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereCluster) ValidateUpdate(_ context.Context, _ runtime.Object, newRaw runtime.Object) (admission.Warnings, error) {
	newTyped, ok := newRaw.(*vmwarev1.VSphereCluster)
	if !ok {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a VSphereCluster but got a %T", newRaw))
	}

	return webhook.validateClusterNetwork(newTyped)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (webhook *VSphereCluster) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func (webhook *VSphereCluster) validateClusterNetwork(cluster *vmwarev1.VSphereCluster) (admission.Warnings, error) {
	if !feature.Gates.Enabled(feature.MultiNetworks) && cluster.Spec.Network.NSXVPC.CreateSubnetSet != nil {
		return nil, apierrors.NewInvalid(cluster.GroupVersionKind().GroupKind(), cluster.Name, field.ErrorList{
			field.Forbidden(field.NewPath("spec", "network", "nsxVPC", "createSubnetSet"), "createSubnetSet can only be set when MultiNetworks feature gate is enabled"),
		})
	}
	if cluster.Spec.Network.NSXVPC.IsDefined() && webhook.NetworkProvider != manager.NSXVPCNetworkProvider {
		return nil, apierrors.NewInvalid(cluster.GroupVersionKind().GroupKind(), cluster.Name, field.ErrorList{
			field.Forbidden(field.NewPath("spec", "network", "nsxVPC"), "nsxVPC can only be set when network provider is NSX-VPC"),
		})
	}
	return nil, nil
}
