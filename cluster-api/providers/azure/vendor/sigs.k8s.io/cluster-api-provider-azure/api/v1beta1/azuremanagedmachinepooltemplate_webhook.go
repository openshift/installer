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

package v1beta1

import (
	"context"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupAzureManagedMachinePoolTemplateWebhookWithManager will set up the webhook to be managed by the specified manager.
func SetupAzureManagedMachinePoolTemplateWebhookWithManager(mgr ctrl.Manager) error {
	mpw := &azureManagedMachinePoolTemplateWebhook{Client: mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AzureManagedMachinePoolTemplate{}).
		WithDefaulter(mpw).
		WithValidator(mpw).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepooltemplate,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepooltemplates,verbs=create;update,versions=v1beta1,name=default.azuremanagedmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

type azureManagedMachinePoolTemplateWebhook struct {
	Client client.Client
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mpw *azureManagedMachinePoolTemplateWebhook) Default(_ context.Context, obj runtime.Object) error {
	mp, ok := obj.(*AzureManagedMachinePoolTemplate)
	if !ok {
		return apierrors.NewBadRequest("expected an AzureManagedMachinePoolTemplate")
	}
	if mp.Labels == nil {
		mp.Labels = make(map[string]string)
	}
	mp.Labels[LabelAgentPoolMode] = mp.Spec.Template.Spec.Mode

	if mp.Spec.Template.Spec.Name == nil || *mp.Spec.Template.Spec.Name == "" {
		mp.Spec.Template.Spec.Name = &mp.Name
	}

	setDefault[*string](&mp.Spec.Template.Spec.OSType, ptr.To(DefaultOSType))

	return nil
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepooltemplate,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepooltemplates,versions=v1beta1,name=validation.azuremanagedmachinepooltemplates.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mpw *azureManagedMachinePoolTemplateWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	mp, ok := obj.(*AzureManagedMachinePoolTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePoolTemplate")
	}

	var errs []error

	errs = append(errs, validateMaxPods(
		mp.Spec.Template.Spec.MaxPods,
		field.NewPath("spec", "template", "spec", "maxPods")))

	errs = append(errs, validateOSType(
		mp.Spec.Template.Spec.Mode,
		mp.Spec.Template.Spec.OSType,
		field.NewPath("spec", "template", "spec", "osType")))

	errs = append(errs, validateMPName(
		mp.Name,
		mp.Spec.Template.Spec.Name,
		mp.Spec.Template.Spec.OSType,
		field.NewPath("spec", "template", "spec", "name")))

	errs = append(errs, validateNodeLabels(
		mp.Spec.Template.Spec.NodeLabels,
		field.NewPath("spec", "template", "spec", "nodeLabels")))

	errs = append(errs, validateNodePublicIPPrefixID(
		mp.Spec.Template.Spec.NodePublicIPPrefixID,
		field.NewPath("spec", "template", "spec", "nodePublicIPPrefixID")))

	errs = append(errs, validateEnableNodePublicIP(
		mp.Spec.Template.Spec.EnableNodePublicIP,
		mp.Spec.Template.Spec.NodePublicIPPrefixID,
		field.NewPath("spec", "template", "spec", "enableNodePublicIP")))

	errs = append(errs, validateKubeletConfig(
		mp.Spec.Template.Spec.KubeletConfig,
		field.NewPath("spec", "template", "spec", "kubeletConfig")))

	errs = append(errs, validateLinuxOSConfig(
		mp.Spec.Template.Spec.LinuxOSConfig,
		mp.Spec.Template.Spec.KubeletConfig,
		field.NewPath("spec", "template", "spec", "linuxOSConfig")))

	return nil, kerrors.NewAggregate(errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mpw *azureManagedMachinePoolTemplateWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldObj.(*AzureManagedMachinePoolTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePoolTemplate")
	}
	mp, ok := newObj.(*AzureManagedMachinePoolTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePoolTemplate")
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "name"),
		old.Spec.Template.Spec.Name,
		mp.Spec.Template.Spec.Name); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := validateNodeLabels(mp.Spec.Template.Spec.NodeLabels, field.NewPath("spec", "template", "spec", "nodeLabels")); err != nil {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "template", "spec", "nodeLabels"),
				mp.Spec.Template.Spec.NodeLabels,
				err.Error()))
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "osType"),
		old.Spec.Template.Spec.OSType,
		mp.Spec.Template.Spec.OSType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "sku"),
		old.Spec.Template.Spec.SKU,
		mp.Spec.Template.Spec.SKU); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "osDiskSizeGB"),
		old.Spec.Template.Spec.OSDiskSizeGB,
		mp.Spec.Template.Spec.OSDiskSizeGB); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "subnetName"),
		old.Spec.Template.Spec.SubnetName,
		mp.Spec.Template.Spec.SubnetName); err != nil && old.Spec.Template.Spec.SubnetName != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "enableFIPS"),
		old.Spec.Template.Spec.EnableFIPS,
		mp.Spec.Template.Spec.EnableFIPS); err != nil && old.Spec.Template.Spec.EnableFIPS != nil {
		allErrs = append(allErrs, err)
	}

	if !webhookutils.EnsureStringSlicesAreEquivalent(mp.Spec.Template.Spec.AvailabilityZones, old.Spec.Template.Spec.AvailabilityZones) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "template", "spec", "availabilityZones"),
				mp.Spec.Template.Spec.AvailabilityZones,
				"field is immutable"))
	}

	if mp.Spec.Template.Spec.Mode != string(NodePoolModeSystem) && old.Spec.Template.Spec.Mode == string(NodePoolModeSystem) {
		// validate for last system node pool
		if err := validateLastSystemNodePool(mpw.Client, mp.Spec.Template.Spec.NodeLabels, mp.Namespace, mp.Annotations); err != nil {
			allErrs = append(allErrs, field.Forbidden(
				field.NewPath("spec", "template", "spec", "mode"),
				"Cannot change node pool mode to User, you must have at least one System node pool in your cluster"))
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "maxPods"),
		old.Spec.Template.Spec.MaxPods,
		mp.Spec.Template.Spec.MaxPods); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "osDiskType"),
		old.Spec.Template.Spec.OsDiskType,
		mp.Spec.Template.Spec.OsDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "scaleSetPriority"),
		old.Spec.Template.Spec.ScaleSetPriority,
		mp.Spec.Template.Spec.ScaleSetPriority); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "enableUltraSSD"),
		old.Spec.Template.Spec.EnableUltraSSD,
		mp.Spec.Template.Spec.EnableUltraSSD); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "enableNodePublicIP"),
		old.Spec.Template.Spec.EnableNodePublicIP,
		mp.Spec.Template.Spec.EnableNodePublicIP); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "nodePublicIPPrefixID"),
		old.Spec.Template.Spec.NodePublicIPPrefixID,
		mp.Spec.Template.Spec.NodePublicIPPrefixID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "kubeletConfig"),
		old.Spec.Template.Spec.KubeletConfig,
		mp.Spec.Template.Spec.KubeletConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "kubeletDiskType"),
		old.Spec.Template.Spec.KubeletDiskType,
		mp.Spec.Template.Spec.KubeletDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "template", "spec", "linuxOSConfig"),
		old.Spec.Template.Spec.LinuxOSConfig,
		mp.Spec.Template.Spec.LinuxOSConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) == 0 {
		return nil, nil
	}
	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureManagedMachinePoolTemplateKind).GroupKind(), mp.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mpw *azureManagedMachinePoolTemplateWebhook) ValidateDelete(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	mp, ok := obj.(*AzureManagedMachinePoolTemplate)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedMachinePoolTemplate")
	}
	if mp.Spec.Template.Spec.Mode != string(NodePoolModeSystem) {
		return nil, nil
	}

	return nil, errors.Wrapf(validateLastSystemNodePool(mpw.Client, mp.Spec.Template.Spec.NodeLabels, mp.Namespace, mp.Annotations), "if the delete is triggered via owner MachinePool please refer to trouble shooting section in https://capz.sigs.k8s.io/topics/managedcluster.html")
}
