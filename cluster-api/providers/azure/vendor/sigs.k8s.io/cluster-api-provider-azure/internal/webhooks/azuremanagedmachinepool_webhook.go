/*
Copyright The Kubernetes Authors.

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
	"regexp"

	"github.com/pkg/errors"
	apierrors "k8s.io/apimachinery/pkg/api/errors"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

var validNodePublicPrefixID = regexp.MustCompile(`(?i)^/?subscriptions/[0-9a-f]{8}-([0-9a-f]{4}-){3}[0-9a-f]{12}/resourcegroups/[^/]+/providers/microsoft\.network/publicipprefixes/[^/]+$`)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (mw *AzureManagedMachinePoolWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	mw.client = mgr.GetClient()

	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureManagedMachinePool{}).
		WithDefaulter(mw).
		WithValidator(mw).
		Complete()
}

//+kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,verbs=create;update,versions=v1beta1,name=default.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureManagedMachinePoolWebhook implements a validating and defaulting webhook for AzureManagedMachinePool.
type AzureManagedMachinePoolWebhook struct {
	client client.Client
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mw *AzureManagedMachinePoolWebhook) Default(_ context.Context, m *infrav1.AzureManagedMachinePool) error {
	if m.Labels == nil {
		m.Labels = make(map[string]string)
	}
	m.Labels[infrav1.LabelAgentPoolMode] = m.Spec.Mode

	if m.Spec.Name == nil || *m.Spec.Name == "" {
		m.Spec.Name = &m.Name
	}

	if m.Spec.OSType == nil {
		m.Spec.OSType = ptr.To(infrav1.DefaultOSType)
	}

	return nil
}

//+kubebuilder:webhook:verbs=create;update;delete,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedmachinepool,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedmachinepools,versions=v1beta1,name=validation.azuremanagedmachinepools.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedMachinePoolWebhook) ValidateCreate(_ context.Context, m *infrav1.AzureManagedMachinePool) (admission.Warnings, error) {
	var errs []error

	errs = append(errs, validateMaxPods(
		m.Spec.MaxPods,
		field.NewPath("spec", "maxPods")))

	errs = append(errs, validateOSType(
		m.Spec.Mode,
		m.Spec.OSType,
		field.NewPath("spec", "osType")))

	errs = append(errs, validateMPName(
		m.Name,
		m.Spec.Name,
		m.Spec.OSType,
		field.NewPath("spec", "name")))

	errs = append(errs, validateNodeLabels(
		m.Spec.NodeLabels,
		field.NewPath("spec", "nodeLabels")))

	errs = append(errs, validateNodePublicIPPrefixID(
		m.Spec.NodePublicIPPrefixID,
		field.NewPath("spec", "nodePublicIPPrefixID")))

	errs = append(errs, validateEnableNodePublicIP(
		m.Spec.EnableNodePublicIP,
		m.Spec.NodePublicIPPrefixID,
		field.NewPath("spec", "enableNodePublicIP")))

	errs = append(errs, validateKubeletConfig(
		m.Spec.KubeletConfig,
		field.NewPath("spec", "kubeletConfig")))

	errs = append(errs, validateLinuxOSConfig(
		m.Spec.LinuxOSConfig,
		m.Spec.KubeletConfig,
		field.NewPath("spec", "linuxOSConfig")))

	errs = append(errs, validateMPSubnetName(
		m.Spec.SubnetName,
		field.NewPath("spec", "subnetName")))

	return nil, kerrors.NewAggregate(errs)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedMachinePoolWebhook) ValidateUpdate(_ context.Context, old, m *infrav1.AzureManagedMachinePool) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "name"),
		old.Spec.Name,
		m.Spec.Name); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := validateNodeLabels(m.Spec.NodeLabels, field.NewPath("spec", "nodeLabels")); err != nil {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "nodeLabels"),
				m.Spec.NodeLabels,
				err.Error()))
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "osType"),
		old.Spec.OSType,
		m.Spec.OSType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "sku"),
		old.Spec.SKU,
		m.Spec.SKU); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "osDiskSizeGB"),
		old.Spec.OSDiskSizeGB,
		m.Spec.OSDiskSizeGB); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "subnetName"),
		old.Spec.SubnetName,
		m.Spec.SubnetName); err != nil && old.Spec.SubnetName != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "enableFIPS"),
		old.Spec.EnableFIPS,
		m.Spec.EnableFIPS); err != nil && old.Spec.EnableFIPS != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "enableEncryptionAtHost"),
		old.Spec.EnableEncryptionAtHost,
		m.Spec.EnableEncryptionAtHost); err != nil && old.Spec.EnableEncryptionAtHost != nil {
		allErrs = append(allErrs, err)
	}

	if !webhookutils.EnsureStringSlicesAreEquivalent(m.Spec.AvailabilityZones, old.Spec.AvailabilityZones) {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "availabilityZones"),
				m.Spec.AvailabilityZones,
				"field is immutable"))
	}

	if m.Spec.Mode != string(infrav1.NodePoolModeSystem) && old.Spec.Mode == string(infrav1.NodePoolModeSystem) {
		// validate for last system node pool
		if err := validateLastSystemNodePool(mw.client, m.Labels, m.Namespace, m.Annotations); err != nil {
			allErrs = append(allErrs, field.Forbidden(
				field.NewPath("spec", "mode"),
				"Cannot change node pool mode to User, you must have at least one System node pool in your cluster"))
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "maxPods"),
		old.Spec.MaxPods,
		m.Spec.MaxPods); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "osDiskType"),
		old.Spec.OsDiskType,
		m.Spec.OsDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "scaleSetPriority"),
		old.Spec.ScaleSetPriority,
		m.Spec.ScaleSetPriority); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "enableUltraSSD"),
		old.Spec.EnableUltraSSD,
		m.Spec.EnableUltraSSD); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "enableNodePublicIP"),
		old.Spec.EnableNodePublicIP,
		m.Spec.EnableNodePublicIP); err != nil {
		allErrs = append(allErrs, err)
	}
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "nodePublicIPPrefixID"),
		old.Spec.NodePublicIPPrefixID,
		m.Spec.NodePublicIPPrefixID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "kubeletConfig"),
		old.Spec.KubeletConfig,
		m.Spec.KubeletConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "kubeletDiskType"),
		old.Spec.KubeletDiskType,
		m.Spec.KubeletDiskType); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "linuxOSConfig"),
		old.Spec.LinuxOSConfig,
		m.Spec.LinuxOSConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if len(allErrs) != 0 {
		return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureManagedMachinePoolKind).GroupKind(), m.Name, allErrs)
	}

	return nil, nil
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *AzureManagedMachinePoolWebhook) ValidateDelete(_ context.Context, m *infrav1.AzureManagedMachinePool) (admission.Warnings, error) {
	if m.Spec.Mode != string(infrav1.NodePoolModeSystem) {
		return nil, nil
	}

	return nil, errors.Wrapf(validateLastSystemNodePool(mw.client, m.Labels, m.Namespace, m.Annotations), "if the delete is triggered via owner MachinePool please refer to trouble shooting section in https://capz.sigs.k8s.io/topics/managedcluster.html")
}
