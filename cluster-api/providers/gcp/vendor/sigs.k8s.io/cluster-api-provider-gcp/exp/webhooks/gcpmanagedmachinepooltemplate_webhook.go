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
	"fmt"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	expinfrav1 "sigs.k8s.io/cluster-api-provider-gcp/exp/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	logf "sigs.k8s.io/controller-runtime/pkg/log"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "sigs.k8s.io/cluster-api-provider-gcp/util/webhook"
)

// log is for logging in this package.
var gmmplog = logf.Log.WithName("gcpmanagedmachinepool-resource")

// SetupGCPManagedMachinePoolTemplateWebhookWithManager sets up and registers the webhook with the manager.
func SetupGCPManagedMachinePoolTemplateWebhookWithManager(mgr ctrl.Manager) error {
	mmptw := &GCPManagedMachinePoolTemplate{Client: mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&expinfrav1.GCPManagedMachinePoolTemplate{}).
		WithDefaulter(mmptw).
		WithValidator(mmptw).
		Complete()
}

// GCPManagedMachinePoolTemplate implements a validating and defaulting webhook for GCPManagedMachinePoolTemplate.
type GCPManagedMachinePoolTemplate struct {
	Client client.Client
}

//+kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-gcpmanagedmachinepooltemplate,mutating=true,failurePolicy=fail,sideEffects=None,groups=infrastructure.cluster.x-k8s.io,resources=gcpmanagedmachinepooltemplates,versions=v1beta1,name=mgcpmanagedmachinepooltemplate.kb.io,admissionReviewVersions=v1

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (*GCPManagedMachinePoolTemplate) Default(_ context.Context, obj runtime.Object) error {
	r, ok := obj.(*expinfrav1.GCPManagedMachinePoolTemplate)
	if !ok {
		return fmt.Errorf("expected a GCPManagedMachinePoolTemplate object but got %T", r)
	}

	gmmplog.Info("default", "name", r.Name)

	return nil
}

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (*GCPManagedMachinePoolTemplate) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	r, ok := obj.(*expinfrav1.GCPManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected an GCPManagedMachinePoolTemplate object but got %T", r)
	}

	gmmplog.Info("Validating GCPManagedMachinePoolTemplate create", "name", r.Name)

	var allErrs field.ErrorList

	if err := validateNodePoolName(
		r.Spec.Template.Spec.NodePoolName,
		field.NewPath("spec", "NodePoolName")); err != nil {
		allErrs = append(allErrs, err)
	}

	if r.Spec.Template.Spec.Scaling != nil {
		if errs := validateScaling(
			r.Spec.Template.Spec.Scaling,
			field.NewPath("spec", "scaling", "minCount"),
			field.NewPath("spec", "scaling", "maxCount"),
			field.NewPath("spec", "scaling", "locationPolicy"),
		); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "diskSizeGb"),
		r.Spec.Template.Spec.DiskSizeGb,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "localSsdCount"),
		r.Spec.Template.Spec.LocalSsdCount,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateNonNegative(
		field.NewPath("spec", "template", "spec", "maxPodsPerNode"),
		r.Spec.Template.Spec.MaxPodsPerNode,
	); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}

	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (*GCPManagedMachinePoolTemplate) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	old, ok := oldObj.(*expinfrav1.GCPManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a GCPManagedMachinePoolTemplate object but got %T", old)
	}

	r, ok := newObj.(*expinfrav1.GCPManagedMachinePoolTemplate)
	if !ok {
		return nil, fmt.Errorf("expected a GCPManagedMachinePoolTemplate object but got %T", r)
	}

	gcpmanagedmachinepoollog.Info("Validating GCPManagedMachinePoolTemplate update", "name", r.Name)

	var allErrs field.ErrorList

	if r.Spec.Template.Spec.Scaling != nil {
		if errs := validateScaling(
			r.Spec.Template.Spec.Scaling,
			field.NewPath("spec", "scaling", "minCount"),
			field.NewPath("spec", "scaling", "maxCount"),
			field.NewPath("spec", "scaling", "locationPolicy"),
		); len(errs) > 0 {
			allErrs = append(allErrs, errs...)
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "instanceType"),
		old.Spec.Template.Spec.InstanceType,
		r.Spec.Template.Spec.InstanceType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodePoolName"),
		old.Spec.Template.Spec.NodePoolName,
		r.Spec.Template.Spec.NodePoolName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "machineType"),
		old.Spec.Template.Spec.MachineType,
		r.Spec.Template.Spec.MachineType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "diskSizeGb"),
		old.Spec.Template.Spec.DiskSizeGb,
		r.Spec.Template.Spec.DiskSizeGb); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "diskType"),
		old.Spec.Template.Spec.DiskType,
		r.Spec.Template.Spec.DiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "localSsdCount"),
		old.Spec.Template.Spec.LocalSsdCount,
		r.Spec.Template.Spec.LocalSsdCount); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "management"),
		old.Spec.Template.Spec.Management,
		r.Spec.Template.Spec.Management); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "maxPodsPerNode"),
		old.Spec.Template.Spec.MaxPodsPerNode,
		r.Spec.Template.Spec.MaxPodsPerNode); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodeNetwork", "podRangeName"),
		old.Spec.Template.Spec.NodeNetwork.PodRangeName,
		r.Spec.Template.Spec.NodeNetwork.PodRangeName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodeNetwork", "createPodRange"),
		old.Spec.Template.Spec.NodeNetwork.CreatePodRange,
		r.Spec.Template.Spec.NodeNetwork.CreatePodRange); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodeNetwork", "podRangeCidrBlock"),
		old.Spec.Template.Spec.NodeNetwork.PodRangeCidrBlock,
		r.Spec.Template.Spec.NodeNetwork.PodRangeCidrBlock); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodeSecurity"),
		old.Spec.Template.Spec.NodeSecurity,
		r.Spec.Template.Spec.NodeSecurity); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(
		r.GroupVersionKind().GroupKind(),
		r.Name,
		allErrs,
	)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (*GCPManagedMachinePoolTemplate) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
