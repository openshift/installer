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

package v1beta1

import (
	"context"
	"fmt"
	"reflect"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/webhook"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
)

// SetupWebhookWithManager sets up and registers the webhook with the manager.
func (c *AzureCluster) SetupWebhookWithManager(mgr ctrl.Manager) error {
	w := new(AzureClusterWebhook)
	return ctrl.NewWebhookManagedBy(mgr).
		For(c).
		WithValidator(w).
		WithDefaulter(w).
		Complete()
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azurecluster,mutating=false,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,versions=v1beta1,name=validation.azurecluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1
// +kubebuilder:webhook:verbs=create;update,path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azurecluster,mutating=true,failurePolicy=fail,matchPolicy=Equivalent,groups=infrastructure.cluster.x-k8s.io,resources=azureclusters,versions=v1beta1,name=default.azurecluster.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// AzureClusterWebhook implements validating and mutating webhook for AzureCluster.
type AzureClusterWebhook struct{}

var _ webhook.CustomValidator = &AzureClusterWebhook{}
var _ webhook.CustomDefaulter = &AzureClusterWebhook{}

// Default implements webhook.CustomDefaulter so a webhook will be registered for the type.
func (*AzureClusterWebhook) Default(_ context.Context, obj runtime.Object) error {
	c, ok := obj.(*AzureCluster)
	if !ok {
		return fmt.Errorf("expected an AzureCluster object but got %T", c)
	}

	c.setDefaults()
	return nil
}

// ValidateCreate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	c, ok := obj.(*AzureCluster)
	if !ok {
		return nil, fmt.Errorf("expected an AzureCluster object but got %T", c)
	}

	return c.validateCluster(nil)
}

// ValidateUpdate implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateUpdate(_ context.Context, oldRaw, newObj runtime.Object) (admission.Warnings, error) {
	c, ok := newObj.(*AzureCluster)
	if !ok {
		return nil, fmt.Errorf("expected an AzureCluster object but got %T", c)
	}

	var allErrs field.ErrorList
	old := oldRaw.(*AzureCluster)

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
		old.setAzureEnvironmentDefault()

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

	allErrs = append(allErrs, c.validateSubnetUpdate(old)...)

	if len(allErrs) == 0 {
		return c.validateCluster(old)
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureClusterKind).GroupKind(), c.Name, allErrs)
}

// validateSubnetUpdate validates a ClusterSpec.NetworkSpec.Subnets for immutability.
func (c *AzureCluster) validateSubnetUpdate(old *AzureCluster) field.ErrorList {
	var allErrs field.ErrorList

	oldSubnetMap := make(map[string]SubnetSpec, len(old.Spec.NetworkSpec.Subnets))
	oldSubnetIndex := make(map[string]int, len(old.Spec.NetworkSpec.Subnets))
	for i, subnet := range old.Spec.NetworkSpec.Subnets {
		oldSubnetMap[subnet.Name] = subnet
		oldSubnetIndex[subnet.Name] = i
	}
	for i, subnet := range c.Spec.NetworkSpec.Subnets {
		if oldSubnet, ok := oldSubnetMap[subnet.Name]; ok {
			// Verify the CIDR blocks haven't changed for an owned Vnet.
			// A non-owned Vnet's CIDR block can change based on what's
			// defined in the spec vs what's been loaded from Azure directly.
			// This technically allows the cidr block to be modified in the brief
			// moments before the Vnet is created (because the tags haven't been
			// set yet) but once the Vnet has been created it becomes immutable.
			if old.Spec.NetworkSpec.Vnet.Tags.HasOwned(old.Name) && !reflect.DeepEqual(subnet.CIDRBlocks, oldSubnet.CIDRBlocks) {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("CIDRBlocks"),
						c.Spec.NetworkSpec.Subnets[i].CIDRBlocks, "field is immutable"),
				)
			}
			if subnet.RouteTable.Name != oldSubnet.RouteTable.Name {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("RouteTable").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].RouteTable.Name, "field is immutable"),
				)
			}
			if (subnet.NatGateway.Name != oldSubnet.NatGateway.Name) && (oldSubnet.NatGateway.Name != "") {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("NatGateway").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].NatGateway.Name, "field is immutable"),
				)
			}
			if subnet.SecurityGroup.Name != oldSubnet.SecurityGroup.Name {
				allErrs = append(allErrs,
					field.Invalid(field.NewPath("spec", "networkSpec", "subnets").Index(oldSubnetIndex[subnet.Name]).Child("SecurityGroup").Child("Name"),
						c.Spec.NetworkSpec.Subnets[i].SecurityGroup.Name, "field is immutable"),
				)
			}
		}
	}

	return allErrs
}

// ValidateDelete implements webhook.CustomValidator so a webhook will be registered for the type.
func (*AzureClusterWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}
