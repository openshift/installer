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
	"errors"
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	kerrors "k8s.io/apimachinery/pkg/util/errors"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	"sigs.k8s.io/cluster-api-provider-azure/feature"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capifeature "sigs.k8s.io/cluster-api/feature"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"
)

var (
	kubeSemver                 = regexp.MustCompile(`^v(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)\.(0|[1-9][0-9]*)([-0-9a-zA-Z_\.+]*)?$`)
	rMaxNodeProvisionTime      = regexp.MustCompile(`^(\d+)m$`)
	rScaleDownTime             = regexp.MustCompile(`^(\d+)m$`)
	rScaleDownDelayAfterDelete = regexp.MustCompile(`^(\d+)s$`)
	rScanInterval              = regexp.MustCompile(`^(\d+)s$`)
)

// SetupAzureManagedControlPlaneWebhookWithManager sets up and registers the webhook with the manager.
func SetupAzureManagedControlPlaneWebhookWithManager(mgr ctrl.Manager) error {
	mw := &azureManagedControlPlaneWebhook{Client: mgr.GetClient()}
	return ctrl.NewWebhookManagedBy(mgr).
		For(&AzureManagedControlPlane{}).
		WithDefaulter(mw).
		WithValidator(mw).
		Complete()
}

// +kubebuilder:webhook:path=/mutate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplane,mutating=true,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,verbs=create;update,versions=v1beta1,name=default.azuremanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// azureManagedControlPlaneWebhook implements a validating and defaulting webhook for AzureManagedControlPlane.
type azureManagedControlPlaneWebhook struct {
	Client client.Client
}

// Default implements webhook.Defaulter so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) Default(ctx context.Context, obj runtime.Object) error {
	m, ok := obj.(*AzureManagedControlPlane)
	if !ok {
		return apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}
	if m.Spec.NetworkPlugin == nil {
		networkPlugin := "azure"
		m.Spec.NetworkPlugin = &networkPlugin
	}
	if m.Spec.LoadBalancerSKU == nil {
		loadBalancerSKU := "Standard"
		m.Spec.LoadBalancerSKU = &loadBalancerSKU
	}

	if m.Spec.Version != "" && !strings.HasPrefix(m.Spec.Version, "v") {
		normalizedVersion := "v" + m.Spec.Version
		m.Spec.Version = normalizedVersion
	}

	if m.Spec.Identity == nil {
		m.Spec.Identity = &Identity{
			Type: ManagedControlPlaneIdentityTypeSystemAssigned,
		}
	}

	if err := m.setDefaultSSHPublicKey(); err != nil {
		ctrl.Log.WithName("AzureManagedControlPlaneWebHookLogger").Error(err, "setDefaultSSHPublicKey failed")
	}

	// PaidManagedControlPlaneTier has been replaced with StandardManagedControlPlaneTier since v2023-02-01.
	if m.Spec.SKU != nil && m.Spec.SKU.Tier == PaidManagedControlPlaneTier {
		m.Spec.SKU.Tier = StandardManagedControlPlaneTier
		ctrl.Log.WithName("AzureManagedControlPlaneWebHookLogger").Info("Paid SKU tier is deprecated and has been replaced by Standard")
	}

	m.setDefaultNodeResourceGroupName()
	m.setDefaultVirtualNetwork()
	m.setDefaultSubnet()
	m.setDefaultSku()
	m.setDefaultAutoScalerProfile()
	m.setDefaultOIDCIssuerProfile()

	return nil
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplane,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,versions=v1beta1,name=validation.azuremanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateCreate(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	m, ok := obj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}
	// NOTE: AzureManagedControlPlane relies upon MachinePools, which is behind a feature gate flag.
	// The webhook must prevent creating new objects in case the feature flag is disabled.
	if !feature.Gates.Enabled(capifeature.MachinePool) {
		return nil, field.Forbidden(
			field.NewPath("spec"),
			"can be set only if the Cluster API 'MachinePool' feature flag is enabled",
		)
	}

	return nil, m.Validate(mw.Client)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateUpdate(ctx context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldObj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}
	m, ok := newObj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SubscriptionID"),
		old.Spec.SubscriptionID,
		m.Spec.SubscriptionID); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "ResourceGroupName"),
		old.Spec.ResourceGroupName,
		m.Spec.ResourceGroupName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "NodeResourceGroupName"),
		old.Spec.NodeResourceGroupName,
		m.Spec.NodeResourceGroupName); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "Location"),
		old.Spec.Location,
		m.Spec.Location); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "SSHPublicKey"),
		old.Spec.SSHPublicKey,
		m.Spec.SSHPublicKey); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "DNSServiceIP"),
		old.Spec.DNSServiceIP,
		m.Spec.DNSServiceIP); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "NetworkPlugin"),
		old.Spec.NetworkPlugin,
		m.Spec.NetworkPlugin); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "NetworkPolicy"),
		old.Spec.NetworkPolicy,
		m.Spec.NetworkPolicy); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "LoadBalancerSKU"),
		old.Spec.LoadBalancerSKU,
		m.Spec.LoadBalancerSKU); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "HTTPProxyConfig"),
		old.Spec.HTTPProxyConfig,
		m.Spec.HTTPProxyConfig); err != nil {
		allErrs = append(allErrs, err)
	}

	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "AzureEnvironment"),
		old.Spec.AzureEnvironment,
		m.Spec.AzureEnvironment); err != nil {
		allErrs = append(allErrs, err)
	}

	if old.Spec.AADProfile != nil {
		if m.Spec.AADProfile == nil {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "AADProfile"),
					m.Spec.AADProfile,
					"field cannot be nil, cannot disable AADProfile"))
		} else {
			if !m.Spec.AADProfile.Managed && old.Spec.AADProfile.Managed {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("Spec", "AADProfile.Managed"),
						m.Spec.AADProfile.Managed,
						"cannot set AADProfile.Managed to false"))
			}
			if len(m.Spec.AADProfile.AdminGroupObjectIDs) == 0 {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("Spec", "AADProfile.AdminGroupObjectIDs"),
						m.Spec.AADProfile.AdminGroupObjectIDs,
						"length of AADProfile.AdminGroupObjectIDs cannot be zero"))
			}
		}
	}

	// Consider removing this once moves out of preview
	// Updating outboundType after cluster creation (PREVIEW)
	// https://learn.microsoft.com/en-us/azure/aks/egress-outboundtype#updating-outboundtype-after-cluster-creation-preview
	if err := webhookutils.ValidateImmutable(
		field.NewPath("Spec", "OutboundType"),
		old.Spec.OutboundType,
		m.Spec.OutboundType); err != nil {
		allErrs = append(allErrs, err)
	}

	if errs := m.validateVirtualNetworkUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateAPIServerAccessProfileUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateNetworkPluginModeUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateOIDCIssuerProfileUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, m.Validate(mw.Client)
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind("AzureManagedControlPlane").GroupKind(), m.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateDelete(ctx context.Context, obj runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Validate the Azure Managed Control Plane and return an aggregate error.
func (m *AzureManagedControlPlane) Validate(cli client.Client) error {
	validators := []func(client client.Client) error{
		m.validateName,
		m.validateVersion,
		m.validateSSHKey,
		m.validateLoadBalancerProfile,
		m.validateAPIServerAccessProfile,
		m.validateManagedClusterNetwork,
		m.validateAutoScalerProfile,
		m.validateIdentity,
		m.validateNetworkPluginMode,
	}

	var errs []error
	for _, validator := range validators {
		if err := validator(cli); err != nil {
			errs = append(errs, err)
		}
	}

	return kerrors.NewAggregate(errs)
}

// validateVersion validates the Kubernetes version.
func (m *AzureManagedControlPlane) validateVersion(_ client.Client) error {
	if !kubeSemver.MatchString(m.Spec.Version) {
		return errors.New("must be a valid semantic version")
	}

	return nil
}

// validateSSHKey validates an SSHKey.
func (m *AzureManagedControlPlane) validateSSHKey(_ client.Client) error {
	if sshKey := m.Spec.SSHPublicKey; sshKey != nil && *sshKey != "" {
		if errs := ValidateSSHKey(*sshKey, field.NewPath("sshKey")); len(errs) > 0 {
			return kerrors.NewAggregate(errs.ToAggregate().Errors())
		}
	}

	return nil
}

// validateLoadBalancerProfile validates a LoadBalancerProfile.
func (m *AzureManagedControlPlane) validateLoadBalancerProfile(_ client.Client) error {
	if m.Spec.LoadBalancerProfile != nil {
		var errs []error
		var allErrs field.ErrorList
		numOutboundIPTypes := 0

		if m.Spec.LoadBalancerProfile.ManagedOutboundIPs != nil {
			if *m.Spec.LoadBalancerProfile.ManagedOutboundIPs < 1 || *m.Spec.LoadBalancerProfile.ManagedOutboundIPs > 100 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "LoadBalancerProfile", "ManagedOutboundIPs"), *m.Spec.LoadBalancerProfile.ManagedOutboundIPs, "value should be in between 1 and 100"))
			}
		}

		if m.Spec.LoadBalancerProfile.AllocatedOutboundPorts != nil {
			if *m.Spec.LoadBalancerProfile.AllocatedOutboundPorts < 0 || *m.Spec.LoadBalancerProfile.AllocatedOutboundPorts > 64000 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "LoadBalancerProfile", "AllocatedOutboundPorts"), *m.Spec.LoadBalancerProfile.AllocatedOutboundPorts, "value should be in between 0 and 64000"))
			}
		}

		if m.Spec.LoadBalancerProfile.IdleTimeoutInMinutes != nil {
			if *m.Spec.LoadBalancerProfile.IdleTimeoutInMinutes < 4 || *m.Spec.LoadBalancerProfile.IdleTimeoutInMinutes > 120 {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "LoadBalancerProfile", "IdleTimeoutInMinutes"), *m.Spec.LoadBalancerProfile.IdleTimeoutInMinutes, "value should be in between 4 and 120"))
			}
		}

		if m.Spec.LoadBalancerProfile.ManagedOutboundIPs != nil {
			numOutboundIPTypes++
		}
		if len(m.Spec.LoadBalancerProfile.OutboundIPPrefixes) > 0 {
			numOutboundIPTypes++
		}
		if len(m.Spec.LoadBalancerProfile.OutboundIPs) > 0 {
			numOutboundIPTypes++
		}
		if numOutboundIPTypes > 1 {
			errs = append(errs, errors.New("load balancer profile must specify at most one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs"))
		}

		if len(allErrs) > 0 {
			agg := kerrors.NewAggregate(allErrs.ToAggregate().Errors())
			errs = append(errs, agg)
		}

		return kerrors.NewAggregate(errs)
	}

	return nil
}

// validateAPIServerAccessProfile validates an APIServerAccessProfile.
func (m *AzureManagedControlPlane) validateAPIServerAccessProfile(_ client.Client) error {
	if m.Spec.APIServerAccessProfile != nil {
		var allErrs field.ErrorList
		for _, ipRange := range m.Spec.APIServerAccessProfile.AuthorizedIPRanges {
			if _, _, err := net.ParseCIDR(ipRange); err != nil {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "APIServerAccessProfile", "AuthorizedIPRanges"), ipRange, "invalid CIDR format"))
			}
		}
		if len(allErrs) > 0 {
			return kerrors.NewAggregate(allErrs.ToAggregate().Errors())
		}
	}
	return nil
}

// validateManagedClusterNetwork validates the Cluster network values.
func (m *AzureManagedControlPlane) validateManagedClusterNetwork(cli client.Client) error {
	ctx := context.Background()

	// Fetch the Cluster.
	clusterName, ok := m.Labels[clusterv1.ClusterNameLabel]
	if !ok {
		return nil
	}

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: m.Namespace,
		Name:      clusterName,
	}

	if err := cli.Get(ctx, key, ownerCluster); err != nil {
		return err
	}

	var (
		allErrs     field.ErrorList
		serviceCIDR string
	)

	if clusterNetwork := ownerCluster.Spec.ClusterNetwork; clusterNetwork != nil {
		if clusterNetwork.Services != nil {
			// A user may provide zero or one CIDR blocks. If they provide an empty array,
			// we ignore it and use the default. AKS doesn't support > 1 Service/Pod CIDR.
			if len(clusterNetwork.Services.CIDRBlocks) > 1 {
				allErrs = append(allErrs, field.TooMany(field.NewPath("Cluster", "Spec", "ClusterNetwork", "Services", "CIDRBlocks"), len(clusterNetwork.Services.CIDRBlocks), 1))
			}
			if len(clusterNetwork.Services.CIDRBlocks) == 1 {
				serviceCIDR = clusterNetwork.Services.CIDRBlocks[0]
			}
		}
		if clusterNetwork.Pods != nil {
			// A user may provide zero or one CIDR blocks. If they provide an empty array,
			// we ignore it and use the default. AKS doesn't support > 1 Service/Pod CIDR.
			if len(clusterNetwork.Pods.CIDRBlocks) > 1 {
				allErrs = append(allErrs, field.TooMany(field.NewPath("Cluster", "Spec", "ClusterNetwork", "Pods", "CIDRBlocks"), len(clusterNetwork.Pods.CIDRBlocks), 1))
			}
		}
	}

	if m.Spec.DNSServiceIP != nil {
		if serviceCIDR == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("Cluster", "Spec", "ClusterNetwork", "Services", "CIDRBlocks"), "service CIDR must be specified if specifying DNSServiceIP"))
		}
		_, cidr, err := net.ParseCIDR(serviceCIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "Spec", "ClusterNetwork", "Services", "CIDRBlocks"), serviceCIDR, fmt.Sprintf("failed to parse cluster service cidr: %v", err)))
		}

		dnsIP := net.ParseIP(*m.Spec.DNSServiceIP)
		if dnsIP == nil { // dnsIP will be nil if the string is not a valid IP
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "DNSServiceIP"), *m.Spec.DNSServiceIP, "must be a valid IP address"))
		}

		if dnsIP != nil && !cidr.Contains(dnsIP) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "Spec", "ClusterNetwork", "Services", "CIDRBlocks"), serviceCIDR, "DNSServiceIP must reside within the associated cluster serviceCIDR"))
		}

		// AKS only supports .10 as the last octet for the DNSServiceIP.
		// Refer to: https://learn.microsoft.com/en-us/azure/aks/configure-kubenet#create-an-aks-cluster-with-system-assigned-managed-identities
		targetSuffix := ".10"
		if dnsIP != nil && !strings.HasSuffix(dnsIP.String(), targetSuffix) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "DNSServiceIP"), *m.Spec.DNSServiceIP, fmt.Sprintf("must end with %q", targetSuffix)))
		}
	}

	if errs := validatePrivateEndpoints(m.Spec.VirtualNetwork.Subnet.PrivateEndpoints, []string{m.Spec.VirtualNetwork.Subnet.CIDRBlock}, field.NewPath("Spec", "VirtualNetwork.Subnet.PrivateEndpoints")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) > 0 {
		return kerrors.NewAggregate(allErrs.ToAggregate().Errors())
	}
	return nil
}

// validateAPIServerAccessProfileUpdate validates update to APIServerAccessProfile.
func (m *AzureManagedControlPlane) validateAPIServerAccessProfileUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	newAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	oldAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	if m.Spec.APIServerAccessProfile != nil {
		newAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			EnablePrivateCluster:           m.Spec.APIServerAccessProfile.EnablePrivateCluster,
			PrivateDNSZone:                 m.Spec.APIServerAccessProfile.PrivateDNSZone,
			EnablePrivateClusterPublicFQDN: m.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
		}
	}
	if old.Spec.APIServerAccessProfile != nil {
		oldAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			EnablePrivateCluster:           old.Spec.APIServerAccessProfile.EnablePrivateCluster,
			PrivateDNSZone:                 old.Spec.APIServerAccessProfile.PrivateDNSZone,
			EnablePrivateClusterPublicFQDN: old.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
		}
	}

	if !reflect.DeepEqual(newAPIServerAccessProfileNormalized, oldAPIServerAccessProfileNormalized) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("Spec", "APIServerAccessProfile"),
				m.Spec.APIServerAccessProfile, "fields (except for AuthorizedIPRanges) are immutable"),
		)
	}

	return allErrs
}

// validateVirtualNetworkUpdate validates update to VirtualNetwork.
func (m *AzureManagedControlPlane) validateVirtualNetworkUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if old.Spec.VirtualNetwork.Name != m.Spec.VirtualNetwork.Name {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "VirtualNetwork.Name"),
				m.Spec.VirtualNetwork.Name,
				"Virtual Network Name is immutable"))
	}

	if old.Spec.VirtualNetwork.CIDRBlock != m.Spec.VirtualNetwork.CIDRBlock {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "VirtualNetwork.CIDRBlock"),
				m.Spec.VirtualNetwork.CIDRBlock,
				"Virtual Network CIDRBlock is immutable"))
	}

	if old.Spec.VirtualNetwork.Subnet.Name != m.Spec.VirtualNetwork.Subnet.Name {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "VirtualNetwork.Subnet.Name"),
				m.Spec.VirtualNetwork.Subnet.Name,
				"Subnet Name is immutable"))
	}

	// NOTE: This only works because we force the user to set the CIDRBlock for both the
	// managed and unmanaged Vnets. If we ever update the subnet cidr based on what's
	// actually set in the subnet, and it is different from what's in the Spec, for
	// unmanaged Vnets like we do with the AzureCluster this logic will break.
	if old.Spec.VirtualNetwork.Subnet.CIDRBlock != m.Spec.VirtualNetwork.Subnet.CIDRBlock {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "VirtualNetwork.Subnet.CIDRBlock"),
				m.Spec.VirtualNetwork.Subnet.CIDRBlock,
				"Subnet CIDRBlock is immutable"))
	}

	if old.Spec.VirtualNetwork.ResourceGroup != m.Spec.VirtualNetwork.ResourceGroup {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("Spec", "VirtualNetwork.ResourceGroup"),
				m.Spec.VirtualNetwork.ResourceGroup,
				"Virtual Network Resource Group is immutable"))
	}
	return allErrs
}

// validateNetworkPluginModeUpdate validates update to NetworkPluginMode.
func (m *AzureManagedControlPlane) validateNetworkPluginModeUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	if ptr.Deref(old.Spec.NetworkPluginMode, "") != NetworkPluginModeOverlay &&
		ptr.Deref(m.Spec.NetworkPluginMode, "") == NetworkPluginModeOverlay &&
		old.Spec.NetworkPolicy != nil {
		allErrs = append(allErrs, field.Forbidden(field.NewPath("Spec", "NetworkPluginMode"), fmt.Sprintf("%q NetworkPluginMode cannot be enabled when NetworkPolicy is set", NetworkPluginModeOverlay)))
	}

	return allErrs
}

// validateOIDCIssuerProfile validates an OIDCIssuerProfile.
func (m *AzureManagedControlPlane) validateOIDCIssuerProfileUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	if m.Spec.OIDCIssuerProfile != nil && old.Spec.OIDCIssuerProfile != nil {
		if m.Spec.OIDCIssuerProfile.Enabled != nil && old.Spec.OIDCIssuerProfile.Enabled != nil &&
			!*m.Spec.OIDCIssuerProfile.Enabled && *old.Spec.OIDCIssuerProfile.Enabled {
			allErrs = append(allErrs,
				field.Forbidden(
					field.NewPath("Spec", "OIDCIssuerProfile", "Enabled"),
					"cannot be disabled",
				),
			)
		}
	}

	return allErrs
}

func (m *AzureManagedControlPlane) validateName(_ client.Client) error {
	if lName := strings.ToLower(m.Name); strings.Contains(lName, "microsoft") ||
		strings.Contains(lName, "windows") {
		return field.Invalid(field.NewPath("Name"), m.Name,
			"cluster name is invalid because 'MICROSOFT' and 'WINDOWS' can't be used as either a whole word or a substring in the name")
	}

	return nil
}

// validateAutoScalerProfile validates an AutoScalerProfile.
func (m *AzureManagedControlPlane) validateAutoScalerProfile(_ client.Client) error {
	var allErrs field.ErrorList

	if m.Spec.AutoScalerProfile == nil {
		return nil
	}

	if errs := m.validateIntegerStringGreaterThanZero(m.Spec.AutoScalerProfile.MaxEmptyBulkDelete, "MaxEmptyBulkDelete"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateIntegerStringGreaterThanZero(m.Spec.AutoScalerProfile.MaxGracefulTerminationSec, "MaxGracefulTerminationSec"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateMaxNodeProvisionTime(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if m.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage != nil {
		val, err := strconv.Atoi(*m.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage)
		if err != nil || val < 0 || val > 100 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "MaxTotalUnreadyPercentage"), m.Spec.AutoScalerProfile.MaxTotalUnreadyPercentage, "invalid value"))
		}
	}

	if errs := m.validateNewPodScaleUpDelay(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateIntegerStringGreaterThanZero(m.Spec.AutoScalerProfile.OkTotalUnreadyCount, "OkTotalUnreadyCount"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScanInterval(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScaleDownTime(m.Spec.AutoScalerProfile.ScaleDownDelayAfterAdd, "ScaleDownDelayAfterAdd"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScaleDownDelayAfterDelete(); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScaleDownTime(m.Spec.AutoScalerProfile.ScaleDownDelayAfterFailure, "ScaleDownDelayAfterFailure"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScaleDownTime(m.Spec.AutoScalerProfile.ScaleDownUnneededTime, "ScaleDownUnneededTime"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateScaleDownTime(m.Spec.AutoScalerProfile.ScaleDownUnreadyTime, "ScaleDownUnreadyTime"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if m.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold != nil {
		val, err := strconv.ParseFloat(*m.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold, 32)
		if err != nil || val < 0 || val > 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "ScaleDownUtilizationThreshold"), m.Spec.AutoScalerProfile.ScaleDownUtilizationThreshold, "invalid value"))
		}
	}

	if len(allErrs) > 0 {
		return kerrors.NewAggregate(allErrs.ToAggregate().Errors())
	}

	return nil
}

// validateMaxNodeProvisionTime validates update to AutoscalerProfile.MaxNodeProvisionTime.
func (m *AzureManagedControlPlane) validateMaxNodeProvisionTime() field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(m.Spec.AutoScalerProfile.MaxNodeProvisionTime, "") != "" {
		if !rMaxNodeProvisionTime.MatchString(ptr.Deref(m.Spec.AutoScalerProfile.MaxNodeProvisionTime, "")) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "MaxNodeProvisionTime"), m.Spec.AutoScalerProfile.MaxNodeProvisionTime, "invalid value"))
		}
	}
	return allErrs
}

// validateScanInterval validates update to AutoscalerProfile.ScanInterval.
func (m *AzureManagedControlPlane) validateScanInterval() field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(m.Spec.AutoScalerProfile.ScanInterval, "") != "" {
		if !rScanInterval.MatchString(ptr.Deref(m.Spec.AutoScalerProfile.ScanInterval, "")) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "ScanInterval"), m.Spec.AutoScalerProfile.ScanInterval, "invalid value"))
		}
	}
	return allErrs
}

// validateNewPodScaleUpDelay validates update to AutoscalerProfile.NewPodScaleUpDelay.
func (m *AzureManagedControlPlane) validateNewPodScaleUpDelay() field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(m.Spec.AutoScalerProfile.NewPodScaleUpDelay, "") != "" {
		_, err := time.ParseDuration(ptr.Deref(m.Spec.AutoScalerProfile.NewPodScaleUpDelay, ""))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "NewPodScaleUpDelay"), m.Spec.AutoScalerProfile.NewPodScaleUpDelay, "invalid value"))
		}
	}
	return allErrs
}

// validateScaleDownDelayAfterDelete validates update to AutoscalerProfile.ScaleDownDelayAfterDelete value.
func (m *AzureManagedControlPlane) validateScaleDownDelayAfterDelete() field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(m.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete, "") != "" {
		if !rScaleDownDelayAfterDelete.MatchString(ptr.Deref(m.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete, "")) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", "ScaleDownDelayAfterDelete"), ptr.Deref(m.Spec.AutoScalerProfile.ScaleDownDelayAfterDelete, ""), "invalid value"))
		}
	}
	return allErrs
}

// validateScaleDownTime validates update to AutoscalerProfile.ScaleDown* values.
func (m *AzureManagedControlPlane) validateScaleDownTime(scaleDownValue *string, fieldName string) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(scaleDownValue, "") != "" {
		if !rScaleDownTime.MatchString(ptr.Deref(scaleDownValue, "")) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", fieldName), ptr.Deref(scaleDownValue, ""), "invalid value"))
		}
	}
	return allErrs
}

// validateIntegerStringGreaterThanZero validates that a string value is an integer greater than zero.
func (m *AzureManagedControlPlane) validateIntegerStringGreaterThanZero(input *string, fieldName string) field.ErrorList {
	var allErrs field.ErrorList

	if input != nil {
		val, err := strconv.Atoi(*input)
		if err != nil || val < 0 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "AutoscalerProfile", fieldName), input, "invalid value"))
		}
	}

	return allErrs
}

// validateIdentity validates an Identity.
func (m *AzureManagedControlPlane) validateIdentity(_ client.Client) error {
	var allErrs field.ErrorList

	if m.Spec.Identity != nil {
		if m.Spec.Identity.Type == ManagedControlPlaneIdentityTypeUserAssigned {
			if m.Spec.Identity.UserAssignedIdentityResourceID == "" {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "Identity", "UserAssignedIdentityResourceID"), m.Spec.Identity.UserAssignedIdentityResourceID, "cannot be empty if Identity.Type is UserAssigned"))
			}
		} else {
			if m.Spec.Identity.UserAssignedIdentityResourceID != "" {
				allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "Identity", "UserAssignedIdentityResourceID"), m.Spec.Identity.UserAssignedIdentityResourceID, "should be empty if Identity.Type is SystemAssigned"))
			}
		}
	}

	if len(allErrs) > 0 {
		return kerrors.NewAggregate(allErrs.ToAggregate().Errors())
	}

	return nil
}

// validateNetworkPluginMode validates a NetworkPluginMode.
func (m *AzureManagedControlPlane) validateNetworkPluginMode(_ client.Client) error {
	var allErrs field.ErrorList

	const kubenet = "kubenet"
	if ptr.Deref(m.Spec.NetworkPluginMode, "") == NetworkPluginModeOverlay &&
		ptr.Deref(m.Spec.NetworkPlugin, "") == kubenet {
		allErrs = append(allErrs, field.Invalid(field.NewPath("Spec", "NetworkPluginMode"), m.Spec.NetworkPluginMode, fmt.Sprintf("cannot be set to %q when NetworkPlugin is %q", NetworkPluginModeOverlay, kubenet)))
	}

	if len(allErrs) > 0 {
		return kerrors.NewAggregate(allErrs.ToAggregate().Errors())
	}

	return nil
}
