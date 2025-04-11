/*
Copyright 2024 The Kubernetes Authors.

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

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api/util/topology"
	"sigs.k8s.io/controller-runtime/pkg/builder"
	"sigs.k8s.io/controller-runtime/pkg/manager"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1alpha1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1alpha1"
	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
)

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1alpha1-openstackserver,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=openstackservers,versions=v1alpha1,name=validation.openstackserver.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1beta1

func SetupOpenStackServerWebhook(mgr manager.Manager) error {
	return builder.WebhookManagedBy(mgr).
		For(&infrav1alpha1.OpenStackServer{}).
		WithValidator(&openStackServerWebhook{}).
		Complete()
}

type openStackServerWebhook struct{}

// Compile-time assertion that openStackServerWebhook implements webhook.CustomValidator.
var _ webhook.CustomValidator = &openStackServerWebhook{}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackServerWebhook) ValidateCreate(_ context.Context, objRaw runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	newObj, err := castToOpenStackServer(objRaw)
	if err != nil {
		return nil, err
	}

	if newObj.Spec.RootVolume != nil && newObj.Spec.AdditionalBlockDevices != nil {
		for _, device := range newObj.Spec.AdditionalBlockDevices {
			if device.Name == "root" {
				allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "additionalBlockDevices"), "cannot contain a device named \"root\" when rootVolume is set"))
			}
		}
	}

	for _, port := range newObj.Spec.Ports {
		if ptr.Deref(port.DisablePortSecurity, false) && len(port.SecurityGroups) > 0 {
			allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "ports"), "cannot have security groups when DisablePortSecurity is set to true"))
		}
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackServerWebhook) ValidateUpdate(ctx context.Context, oldObjRaw, newObjRaw runtime.Object) (admission.Warnings, error) {
	oldObj, err := castToOpenStackServer(oldObjRaw)
	if err != nil {
		return nil, err
	}

	newObj, err := castToOpenStackServer(newObjRaw)
	if err != nil {
		return nil, err
	}

	req, err := admission.RequestFromContext(ctx)
	if err != nil {
		return nil, apierrors.NewBadRequest(fmt.Sprintf("expected a admission.Request inside context: %v", err))
	}

	newOpenStackServer, err := runtime.DefaultUnstructuredConverter.ToUnstructured(newObj)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.SchemeGroupVersion.WithKind("OpenStackServer").GroupKind(), newObj.Name, field.ErrorList{
			field.InternalError(nil, fmt.Errorf("failed to convert new OpenStackServer to unstructured object: %w", err)),
		})
	}
	oldOpenStackServer, err := runtime.DefaultUnstructuredConverter.ToUnstructured(oldObjRaw)
	if err != nil {
		return nil, apierrors.NewInvalid(infrav1.SchemeGroupVersion.WithKind("OpenStackServer").GroupKind(), newObj.Name, field.ErrorList{
			field.InternalError(nil, fmt.Errorf("failed to convert old OpenStackServer to unstructured object: %w", err)),
		})
	}

	var allErrs field.ErrorList

	newOpenStackServerSpec := newOpenStackServer["spec"].(map[string]interface{})
	oldOpenStackServerSpec := oldOpenStackServer["spec"].(map[string]interface{})

	// allow changes to identifyRef
	delete(oldOpenStackServerSpec, "identityRef")
	delete(newOpenStackServerSpec, "identityRef")

	if !topology.ShouldSkipImmutabilityChecks(req, newObj) &&
		!reflect.DeepEqual(newObj.Spec, oldObj.Spec) {
		allErrs = append(allErrs,
			field.Forbidden(field.NewPath("spec"), "OpenStackServer spec field is immutable. Please create a new resource instead."),
		)
	}

	return aggregateObjErrors(newObj.GroupVersionKind().GroupKind(), newObj.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*openStackServerWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

func castToOpenStackServer(obj runtime.Object) (*infrav1alpha1.OpenStackServer, error) {
	cast, ok := obj.(*infrav1alpha1.OpenStackServer)
	if !ok {
		return nil, fmt.Errorf("expected an OpenStackServer but got a %T", obj)
	}
	return cast, nil
}
