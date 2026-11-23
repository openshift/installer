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
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	infrav1 "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"
	apiinternal "sigs.k8s.io/cluster-api-provider-azure/internal/api/v1beta1"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (w *AzureClusterWebhook) SetupWebhookWithManager(mgr ctrl.Manager) error {
	return ctrl.NewWebhookManagedBy(mgr, &infrav1.AzureCluster{}).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azurecluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,versions=v1beta1,name=validation.azurecluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azurecluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,versions=v1beta1,name=default.azurecluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureClusterWebhook implements validating and mutating webhook for AzureCluster.
type AzureClusterWebhook struct{}

var _ admission.Validator[*infrav1.AzureCluster] = &AzureClusterWebhook{}
var _ admission.Defaulter[*infrav1.AzureCluster] = &AzureClusterWebhook{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*AzureClusterWebhook) Default(_ context.Context, c *infrav1.AzureCluster) error {
	apiinternal.SetDefaultsAzureCluster(c)
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateCreate(_ context.Context, c *infrav1.AzureCluster) (admission.Warnings, error) {
	return validateAzureCluster(c, nil)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateUpdate(_ context.Context, old, c *infrav1.AzureCluster) (admission.Warnings, error) {
	var allErrs field.ErrorList

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "resourceGroup"),
		old.Spec.ResourceGroup,
		c.Spec.ResourceGroup); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "subscriptionID"),
		old.Spec.SubscriptionID,
		c.Spec.SubscriptionID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "location"),
		old.Spec.Location,
		c.Spec.Location); err != nil {
		allErrs = append(allErrs, err)
	}

	if old.Spec.ControlPlaneEndpoint.Host != "" && c.Spec.ControlPlaneEndpoint.Host != old.Spec.ControlPlaneEndpoint.Host {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneEndpoint", "host"),
				c.Spec.ControlPlaneEndpoint.Host, "field is immutable"),
		)
	}

	if old.Spec.ControlPlaneEndpoint.Port != 0 && c.Spec.ControlPlaneEndpoint.Port != old.Spec.ControlPlaneEndpoint.Port {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "controlPlaneEndpoint", "port"),
				c.Spec.ControlPlaneEndpoint.Port, "field is immutable"),
		)
	}

	if !reflect.DeepEqual(c.Spec.AzureEnvironment, old.Spec.AzureEnvironment) {
		// The equality failure could be because of default mismatch between v1alpha3 and v1beta1. This happens because
		// the new object `r` will have run through the default webhooks but the old object `old` would not have so.
		// This means if the old object was in v1alpha3, it would not get the new defaults set in v1beta1 resulting
		// in object inequality. To workaround this, we set the v1beta1 defaults here so that the old object also gets
		// the new defaults.
		apiinternal.SetDefaultAzureClusterAzureEnvironment(old)

		// if it's still not equal, return error.
		if !reflect.DeepEqual(c.Spec.AzureEnvironment, old.Spec.AzureEnvironment) {
			allErrs = append(allErrs,
				field.Invalid(field.NewPath("spec", "azureEnvironment"),
					c.Spec.AzureEnvironment, "field is immutable"),
			)
		}
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "networkSpec", "privateDNSZoneName"),
		old.Spec.NetworkSpec.PrivateDNSZoneName,
		c.Spec.NetworkSpec.PrivateDNSZoneName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "networkSpec", "privateDNSZoneResourceGroup"),
		old.Spec.NetworkSpec.PrivateDNSZoneResourceGroup,
		c.Spec.NetworkSpec.PrivateDNSZoneResourceGroup); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "networkSpec", "privateDNSZone"),
		old.Spec.NetworkSpec.PrivateDNSZone,
		c.Spec.NetworkSpec.PrivateDNSZone); err != nil {
		allErrs = append(allErrs, err)
	}

	// Allow enabling azure bastion but avoid disabling it.
	if old.Spec.BastionSpec.AzureBastion != nil && !reflect.DeepEqual(old.Spec.BastionSpec.AzureBastion, c.Spec.BastionSpec.AzureBastion) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "bastionSpec", "azureBastion"),
				c.Spec.BastionSpec.AzureBastion, "azure bastion cannot be removed from a cluster"),
		)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "networkSpec", "controlPlaneOutboundLB"),
		old.Spec.NetworkSpec.ControlPlaneOutboundLB,
		c.Spec.NetworkSpec.ControlPlaneOutboundLB); err != nil {
		allErrs = append(allErrs, err)
	}

	// Validate availability zones are immutable for load balancers
	if c.Spec.NetworkSpec.APIServerLB != nil && old.Spec.NetworkSpec.APIServerLB != nil {
		if !webhookutils.EnsureStringSlicesAreEquivalent(
			c.Spec.NetworkSpec.APIServerLB.AvailabilityZones,
			old.Spec.NetworkSpec.APIServerLB.AvailabilityZones) {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "networkSpec", "apiServerLB", "availabilityZones"),
					c.Spec.NetworkSpec.APIServerLB.AvailabilityZones,
					"field is immutable"))
		}
	}

	if c.Spec.NetworkSpec.NodeOutboundLB != nil && old.Spec.NetworkSpec.NodeOutboundLB != nil {
		if !webhookutils.EnsureStringSlicesAreEquivalent(
			c.Spec.NetworkSpec.NodeOutboundLB.AvailabilityZones,
			old.Spec.NetworkSpec.NodeOutboundLB.AvailabilityZones) {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "networkSpec", "nodeOutboundLB", "availabilityZones"),
					c.Spec.NetworkSpec.NodeOutboundLB.AvailabilityZones,
					"field is immutable"))
		}
	}

	if c.Spec.NetworkSpec.ControlPlaneOutboundLB != nil && old.Spec.NetworkSpec.ControlPlaneOutboundLB != nil {
		if !webhookutils.EnsureStringSlicesAreEquivalent(
			c.Spec.NetworkSpec.ControlPlaneOutboundLB.AvailabilityZones,
			old.Spec.NetworkSpec.ControlPlaneOutboundLB.AvailabilityZones) {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "networkSpec", "controlPlaneOutboundLB", "availabilityZones"),
					c.Spec.NetworkSpec.ControlPlaneOutboundLB.AvailabilityZones,
					"field is immutable"))
		}
	}

	allErrs = append(allErrs, validateAzureClusterSubnetUpdate(c, old)...)

	if len(allErrs) == 0 {
		return validateAzureCluster(c, old)
	}

	return nil, apierrors.NewInvalid(infrav1.GroupVersion.WithKind(infrav1.AzureClusterKind).GroupKind(), c.Name, allErrs)
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateDelete(_ context.Context, _ *infrav1.AzureCluster) (admission.Warnings, error) {
	return nil, nil
}
