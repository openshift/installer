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
	"fmt"
	"net"
	"reflect"
	"regexp"
	"strconv"
	"strings"
	"time"

	apierrors "k8s.io/apimachinery/pkg/api/errors"
	"k8s.io/apimachinery/pkg/runtime"
	"k8s.io/apimachinery/pkg/util/validation/field"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/client"
	"sigs.k8s.io/controller-runtime/pkg/webhook/admission"

	"sigs.k8s.io/cluster-api-provider-azure/util/versions"
	webhookutils "sigs.k8s.io/cluster-api-provider-azure/util/webhook"
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
func (mw *azureManagedControlPlaneWebhook) Default(_ context.Context, obj runtime.Object) error {
	m, ok := obj.(*AzureManagedControlPlane)
	if !ok {
		return apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}

	m.Spec.Version = setDefaultVersion(m.Spec.Version)
	m.Spec.SKU = setDefaultSku(m.Spec.SKU)
	m.Spec.FleetsMember = setDefaultFleetsMember(m.Spec.FleetsMember, m.Labels)

	if err := m.setDefaultSSHPublicKey(); err != nil {
		ctrl.Log.WithName("AzureManagedControlPlaneWebHookLogger").Error(err, "setDefaultSSHPublicKey failed")
	}

	m.setDefaultResourceGroupName()
	m.setDefaultNodeResourceGroupName()
	m.setDefaultVirtualNetwork()
	m.setDefaultSubnet()
	m.setDefaultOIDCIssuerProfile()
	m.setDefaultDNSPrefix()
	m.setDefaultAKSExtensions()

	return nil
}

// +kubebuilder:webhook:verbs=create;update,path=/validate-infrastructure-cluster-x-k8s-io-v1beta1-azuremanagedcontrolplane,mutating=false,failurePolicy=fail,groups=infrastructure.cluster.x-k8s.io,resources=azuremanagedcontrolplanes,versions=v1beta1,name=validation.azuremanagedcontrolplanes.infrastructure.cluster.x-k8s.io,sideEffects=None,admissionReviewVersions=v1;v1beta1

// ValidateCreate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateCreate(_ context.Context, obj runtime.Object) (admission.Warnings, error) {
	m, ok := obj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}

	return nil, m.Validate(mw.Client)
}

// ValidateUpdate implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateUpdate(_ context.Context, oldObj, newObj runtime.Object) (admission.Warnings, error) {
	var allErrs field.ErrorList
	old, ok := oldObj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}
	m, ok := newObj.(*AzureManagedControlPlane)
	if !ok {
		return nil, apierrors.NewBadRequest("expected an AzureManagedControlPlane")
	}

	immutableFields := []struct {
		path *field.Path
		old  interface{}
		new  interface{}
	}{
		{field.NewPath("spec", "subscriptionID"), old.Spec.SubscriptionID, m.Spec.SubscriptionID},
		{field.NewPath("spec", "resourceGroupName"), old.Spec.ResourceGroupName, m.Spec.ResourceGroupName},
		{field.NewPath("spec", "nodeResourceGroupName"), old.Spec.NodeResourceGroupName, m.Spec.NodeResourceGroupName},
		{field.NewPath("spec", "location"), old.Spec.Location, m.Spec.Location},
		{field.NewPath("spec", "sshPublicKey"), old.Spec.SSHPublicKey, m.Spec.SSHPublicKey},
		{field.NewPath("spec", "dnsServiceIP"), old.Spec.DNSServiceIP, m.Spec.DNSServiceIP},
		{field.NewPath("spec", "networkPlugin"), old.Spec.NetworkPlugin, m.Spec.NetworkPlugin},
		{field.NewPath("spec", "networkPolicy"), old.Spec.NetworkPolicy, m.Spec.NetworkPolicy},
		{field.NewPath("spec", "networkDataplane"), old.Spec.NetworkDataplane, m.Spec.NetworkDataplane},
		{field.NewPath("spec", "loadBalancerSKU"), old.Spec.LoadBalancerSKU, m.Spec.LoadBalancerSKU},
		{field.NewPath("spec", "httpProxyConfig"), old.Spec.HTTPProxyConfig, m.Spec.HTTPProxyConfig},
		{field.NewPath("spec", "azureEnvironment"), old.Spec.AzureEnvironment, m.Spec.AzureEnvironment},
	}

	for _, f := range immutableFields {
		if err := webhookutils.ValidateImmutable(f.path, f.old, f.new); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	// This nil check is only to streamline tests from having to define this correctly in every test case.
	// Normally, the defaulting webhooks will always set the new DNSPrefix so users can never entirely unset it.
	if m.Spec.DNSPrefix != nil {
		// Pre-1.12 versions of CAPZ do not set this field while 1.12+ defaults it, so emulate the current
		// defaulting here to avoid unrelated updates from failing this immutability check due to the
		// nil -> non-nil transition.
		oldDNSPrefix := old.Spec.DNSPrefix
		if oldDNSPrefix == nil {
			oldDNSPrefix = ptr.To(old.Name)
		}
		if err := webhookutils.ValidateImmutable(
			field.NewPath("spec", "dnsPrefix"),
			oldDNSPrefix,
			m.Spec.DNSPrefix,
		); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	// Consider removing this once moves out of preview
	// Updating outboundType after cluster creation (PREVIEW)
	// https://learn.microsoft.com/en-us/azure/aks/egress-outboundtype#updating-outboundtype-after-cluster-creation-preview
	if err := webhookutils.ValidateImmutable(
		field.NewPath("spec", "outboundType"),
		old.Spec.OutboundType,
		m.Spec.OutboundType); err != nil {
		allErrs = append(allErrs, err)
	}

	if errs := m.validateVirtualNetworkUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateAddonProfilesUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateAPIServerAccessProfileUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateNetworkPluginModeUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateAADProfileUpdateAndLocalAccounts(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateAutoUpgradeProfile(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateK8sVersionUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateOIDCIssuerProfileUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.validateFleetsMemberUpdate(old); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateAKSExtensionsUpdate(old.Spec.Extensions, m.Spec.Extensions); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := m.Spec.AzureManagedControlPlaneClassSpec.validateSecurityProfileUpdate(&old.Spec.AzureManagedControlPlaneClassSpec); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if len(allErrs) == 0 {
		return nil, m.Validate(mw.Client)
	}

	return nil, apierrors.NewInvalid(GroupVersion.WithKind(AzureManagedControlPlaneKind).GroupKind(), m.Name, allErrs)
}

// ValidateDelete implements webhook.Validator so a webhook will be registered for the type.
func (mw *azureManagedControlPlaneWebhook) ValidateDelete(_ context.Context, _ runtime.Object) (admission.Warnings, error) {
	return nil, nil
}

// Validate the Azure Managed Control Plane and return an aggregate error.
func (m *AzureManagedControlPlane) Validate(cli client.Client) error {
	var allErrs field.ErrorList
	validators := []func(client client.Client) field.ErrorList{
		m.validateSSHKey,
		m.validateIdentity,
		m.validateNetworkPluginMode,
		m.validateDNSPrefix,
		m.validateDisableLocalAccounts,
	}
	for _, validator := range validators {
		if err := validator(cli); err != nil {
			allErrs = append(allErrs, err...)
		}
	}

	allErrs = append(allErrs, validateVersion(
		m.Spec.Version,
		field.NewPath("spec").Child("version"))...)

	allErrs = append(allErrs, validateLoadBalancerProfile(
		m.Spec.LoadBalancerProfile,
		field.NewPath("spec").Child("loadBalancerProfile"))...)

	allErrs = append(allErrs, validateManagedClusterNetwork(
		cli,
		m.Labels,
		m.Namespace,
		m.Spec.DNSServiceIP,
		m.Spec.VirtualNetwork.Subnet,
		field.NewPath("spec"))...)

	allErrs = append(allErrs, validateName(m.Name, field.NewPath("name"))...)

	allErrs = append(allErrs, validateAutoScalerProfile(m.Spec.AutoScalerProfile, field.NewPath("spec").Child("autoScalerProfile"))...)

	allErrs = append(allErrs, validateAKSExtensions(m.Spec.Extensions, field.NewPath("spec").Child("aksExtensions"))...)

	allErrs = append(allErrs, m.Spec.AzureManagedControlPlaneClassSpec.validateSecurityProfile()...)

	allErrs = append(allErrs, validateNetworkPolicy(m.Spec.NetworkPolicy, m.Spec.NetworkDataplane, field.NewPath("spec").Child("networkPolicy"))...)

	allErrs = append(allErrs, validateNetworkDataplane(m.Spec.NetworkDataplane, m.Spec.NetworkPolicy, m.Spec.NetworkPluginMode, field.NewPath("spec").Child("networkDataplane"))...)

	allErrs = append(allErrs, validateAPIServerAccessProfile(m.Spec.APIServerAccessProfile, field.NewPath("spec").Child("apiServerAccessProfile"))...)

	allErrs = append(allErrs, validateAMCPVirtualNetwork(m.Spec.VirtualNetwork, field.NewPath("spec").Child("virtualNetwork"))...)

	allErrs = append(allErrs, validateFleetsMember(m.Spec.FleetsMember, field.NewPath("spec").Child("fleetsMember"))...)

	return allErrs.ToAggregate()
}

func (m *AzureManagedControlPlane) validateDNSPrefix(_ client.Client) field.ErrorList {
	if m.Spec.DNSPrefix == nil {
		return nil
	}

	// Regex pattern for DNS prefix validation
	// 1. Between 1 and 54 characters long: {1,54}
	// 2. Alphanumerics and hyphens: [a-zA-Z0-9-]
	// 3. Start and end with alphanumeric: ^[a-zA-Z0-9].*[a-zA-Z0-9]$
	pattern := `^[a-zA-Z0-9][a-zA-Z0-9-]{0,52}[a-zA-Z0-9]$`
	regex := regexp.MustCompile(pattern)
	if regex.MatchString(ptr.Deref(m.Spec.DNSPrefix, "")) {
		return nil
	}
	allErrs := field.ErrorList{
		field.Invalid(field.NewPath("spec", "dnsPrefix"), *m.Spec.DNSPrefix, "DNSPrefix is invalid, does not match regex: "+pattern),
	}
	return allErrs
}

// validateSecurityProfile validates SecurityProfile.
func (m *AzureManagedControlPlaneClassSpec) validateSecurityProfile() field.ErrorList {
	allErrs := field.ErrorList{}
	if err := m.validateAzureKeyVaultKms(); err != nil {
		allErrs = append(allErrs, err...)
	}
	if err := m.validateWorkloadIdentity(); err != nil {
		allErrs = append(allErrs, err...)
	}
	return allErrs
}

// validateAzureKeyVaultKms validates AzureKeyVaultKms.
func (m *AzureManagedControlPlaneClassSpec) validateAzureKeyVaultKms() field.ErrorList {
	if m.SecurityProfile != nil && m.SecurityProfile.AzureKeyVaultKms != nil {
		if !m.isUserManagedIdentityEnabled() {
			allErrs := field.ErrorList{
				field.Invalid(field.NewPath("spec", "securityProfile", "azureKeyVaultKms", "keyVaultResourceID"),
					m.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID,
					"Spec.SecurityProfile.AzureKeyVaultKms can be set only when Spec.Identity.Type is UserAssigned"),
			}
			return allErrs
		}
		keyVaultNetworkAccess := ptr.Deref(m.SecurityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess, KeyVaultNetworkAccessTypesPublic)
		keyVaultResourceID := ptr.Deref(m.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID, "")
		if keyVaultNetworkAccess == KeyVaultNetworkAccessTypesPrivate && keyVaultResourceID == "" {
			allErrs := field.ErrorList{
				field.Invalid(field.NewPath("spec", "securityProfile", "azureKeyVaultKms", "keyVaultResourceID"),
					m.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID,
					"Spec.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID cannot be empty when Spec.SecurityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess is Private"),
			}
			return allErrs
		}
		if keyVaultNetworkAccess == KeyVaultNetworkAccessTypesPublic && keyVaultResourceID != "" {
			allErrs := field.ErrorList{
				field.Invalid(field.NewPath("spec", "securityProfile", "azureKeyVaultKms", "keyVaultResourceID"), m.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID,
					"Spec.SecurityProfile.AzureKeyVaultKms.KeyVaultResourceID should be empty when Spec.SecurityProfile.AzureKeyVaultKms.KeyVaultNetworkAccess is Public"),
			}
			return allErrs
		}
	}
	return nil
}

// validateWorkloadIdentity validates WorkloadIdentity.
func (m *AzureManagedControlPlaneClassSpec) validateWorkloadIdentity() field.ErrorList {
	if m.SecurityProfile != nil && m.SecurityProfile.WorkloadIdentity != nil && !m.isOIDCEnabled() {
		allErrs := field.ErrorList{
			field.Invalid(field.NewPath("spec", "securityProfile", "workloadIdentity"), m.SecurityProfile.WorkloadIdentity,
				"Spec.SecurityProfile.WorkloadIdentity cannot be enabled when Spec.OIDCIssuerProfile is disabled"),
		}
		return allErrs
	}
	return nil
}

// validateDisableLocalAccounts disabling local accounts for AAD based clusters.
func (m *AzureManagedControlPlane) validateDisableLocalAccounts(_ client.Client) field.ErrorList {
	if m.Spec.DisableLocalAccounts != nil && m.Spec.AADProfile == nil {
		return field.ErrorList{
			field.Invalid(field.NewPath("spec", "disableLocalAccounts"), *m.Spec.DisableLocalAccounts, "DisableLocalAccounts should be set only for AAD enabled clusters"),
		}
	}
	return nil
}

// validateVersion validates the Kubernetes version.
func validateVersion(version string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if !kubeSemver.MatchString(version) {
		allErrs = append(allErrs, field.Invalid(fldPath, version, "must be a valid semantic version"))
	}

	return allErrs
}

// validateSSHKey validates an SSHKey.
func (m *AzureManagedControlPlane) validateSSHKey(_ client.Client) field.ErrorList {
	if sshKey := m.Spec.SSHPublicKey; sshKey != nil && *sshKey != "" {
		if errs := ValidateSSHKey(*sshKey, field.NewPath("sshKey")); len(errs) > 0 {
			return errs
		}
	}

	return nil
}

// validateLoadBalancerProfile validates a LoadBalancerProfile.
func validateLoadBalancerProfile(loadBalancerProfile *LoadBalancerProfile, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if loadBalancerProfile != nil {
		numOutboundIPTypes := 0

		if loadBalancerProfile.ManagedOutboundIPs != nil {
			if *loadBalancerProfile.ManagedOutboundIPs < 1 || *loadBalancerProfile.ManagedOutboundIPs > 100 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("ManagedOutboundIPs"), *loadBalancerProfile.ManagedOutboundIPs, "value should be in between 1 and 100"))
			}
		}

		if loadBalancerProfile.AllocatedOutboundPorts != nil {
			if *loadBalancerProfile.AllocatedOutboundPorts < 0 || *loadBalancerProfile.AllocatedOutboundPorts > 64000 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("AllocatedOutboundPorts"), *loadBalancerProfile.AllocatedOutboundPorts, "value should be in between 0 and 64000"))
			}
		}

		if loadBalancerProfile.IdleTimeoutInMinutes != nil {
			if *loadBalancerProfile.IdleTimeoutInMinutes < 4 || *loadBalancerProfile.IdleTimeoutInMinutes > 120 {
				allErrs = append(allErrs, field.Invalid(fldPath.Child("IdleTimeoutInMinutes"), *loadBalancerProfile.IdleTimeoutInMinutes, "value should be in between 4 and 120"))
			}
		}

		if loadBalancerProfile.ManagedOutboundIPs != nil {
			numOutboundIPTypes++
		}
		if len(loadBalancerProfile.OutboundIPPrefixes) > 0 {
			numOutboundIPTypes++
		}
		if len(loadBalancerProfile.OutboundIPs) > 0 {
			numOutboundIPTypes++
		}
		if numOutboundIPTypes > 1 {
			allErrs = append(allErrs, field.Forbidden(fldPath, "load balancer profile must specify at most one of ManagedOutboundIPs, OutboundIPPrefixes and OutboundIPs"))
		}
	}

	return allErrs
}

func validateAMCPVirtualNetwork(virtualNetwork ManagedControlPlaneVirtualNetwork, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	// VirtualNetwork and the CIDR blocks get defaulted in the defaulting webhook, so we can assume they are always set.
	if !reflect.DeepEqual(virtualNetwork, ManagedControlPlaneVirtualNetwork{}) {
		_, parentNet, vnetErr := net.ParseCIDR(virtualNetwork.CIDRBlock)
		if vnetErr != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("CIDRBlock"), virtualNetwork.CIDRBlock, "pre-existing virtual networks CIDR block is invalid"))
		}
		subnetIP, _, subnetErr := net.ParseCIDR(virtualNetwork.Subnet.CIDRBlock)
		if subnetErr != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("Subnet", "CIDRBlock"), virtualNetwork.CIDRBlock, "pre-existing subnets CIDR block is invalid"))
		}
		if vnetErr == nil && subnetErr == nil && !parentNet.Contains(subnetIP) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("CIDRBlock"), virtualNetwork.CIDRBlock, "pre-existing virtual networks CIDR block should contain the subnet CIDR block"))
		}
	}
	return allErrs
}

func validateFleetsMember(fleetsMember *FleetsMember, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if fleetsMember != nil && fleetsMember.Name != "" {
		match, _ := regexp.MatchString(`^[a-z0-9]([-a-z0-9]*[a-z0-9])?$`, fleetsMember.Name)
		if !match {
			allErrs = append(allErrs,
				field.Invalid(
					fldPath.Child("Name"),
					fleetsMember.Name,
					"Name must match ^[a-z0-9]([-a-z0-9]*[a-z0-9])?$",
				),
			)
		}
	}

	return allErrs
}

// validateAPIServerAccessProfile validates an APIServerAccessProfile.
func validateAPIServerAccessProfile(apiServerAccessProfile *APIServerAccessProfile, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if apiServerAccessProfile != nil {
		for _, ipRange := range apiServerAccessProfile.AuthorizedIPRanges {
			if _, _, err := net.ParseCIDR(ipRange); err != nil {
				allErrs = append(allErrs, field.Invalid(fldPath, ipRange, "invalid CIDR format"))
			}
		}

		// privateDNSZone should either be "System" or "None" or the private dns zone name should be in either of these
		// formats: 'private.<location>.azmk8s.io,privatelink.<location>.azmk8s.io,[a-zA-Z0-9-]{1,32}.private.<location>.azmk8s.io,
		// [a-zA-Z0-9-]{1,32}.privatelink.<location>.azmk8s.io'. The validation below follows the guidelines mentioned at
		// https://learn.microsoft.com/azure/aks/private-clusters?tabs=azure-portal#configure-a-private-dns-zone.
		// Performing a lower case comparison to avoid case sensitivity.
		if apiServerAccessProfile.PrivateDNSZone != nil {
			privateDNSZone := strings.ToLower(ptr.Deref(apiServerAccessProfile.PrivateDNSZone, ""))
			if !strings.EqualFold(strings.ToLower(privateDNSZone), "system") &&
				!strings.EqualFold(strings.ToLower(privateDNSZone), "none") {
				// Extract substring starting from "privatednszones/"
				startIndex := strings.Index(strings.ToLower(privateDNSZone), "privatednszones/")
				if startIndex == -1 {
					allErrs = append(allErrs, field.Invalid(fldPath, privateDNSZone, "invalid private DNS zone"))
					return allErrs
				}

				// Private DNS Zones can only be used by private clusters.
				if !ptr.Deref(apiServerAccessProfile.EnablePrivateCluster, false) {
					allErrs = append(allErrs, field.Invalid(fldPath, apiServerAccessProfile.EnablePrivateCluster, "Private Cluster should be enabled to use PrivateDNSZone"))
					return allErrs
				}

				extractedPrivateDNSZone := privateDNSZone[startIndex+len("privatednszones/"):]

				patternWithLocation := `^(privatelink|private)\.[a-zA-Z0-9]+\.(azmk8s\.io)$`
				locationRegex := regexp.MustCompile(patternWithLocation)
				patternWithSubzone := `^[a-zA-Z0-9-]{1,32}\.(privatelink|private)\.[a-zA-Z0-9]+\.(azmk8s\.io)$`
				subzoneRegex := regexp.MustCompile(patternWithSubzone)

				// check if privateDNSZone is a valid resource ID
				if !locationRegex.MatchString(extractedPrivateDNSZone) && !subzoneRegex.MatchString(extractedPrivateDNSZone) {
					allErrs = append(allErrs, field.Invalid(fldPath, privateDNSZone, "invalid privateDnsZone resource ID. Each label the private dns zone name should be in either of these formats: 'private.<location>.azmk8s.io,privatelink.<location>.azmk8s.io,[a-zA-Z0-9-]{1,32}.private.<location>.azmk8s.io,[a-zA-Z0-9-]{1,32}.privatelink.<location>.azmk8s.io'"))
				}
			}
		}
	}
	return allErrs
}

// validateManagedClusterNetwork validates the Cluster network values.
func validateManagedClusterNetwork(cli client.Client, labels map[string]string, namespace string, dnsServiceIP *string, subnet ManagedControlPlaneSubnet, fldPath *field.Path) field.ErrorList {
	var (
		allErrs     field.ErrorList
		serviceCIDR string
	)

	ctx := context.Background()

	// Fetch the Cluster.
	clusterName, ok := labels[clusterv1.ClusterNameLabel]
	if !ok {
		return nil
	}

	ownerCluster := &clusterv1.Cluster{}
	key := client.ObjectKey{
		Namespace: namespace,
		Name:      clusterName,
	}

	if err := cli.Get(ctx, key, ownerCluster); err != nil {
		allErrs = append(allErrs, field.InternalError(field.NewPath("Cluster", "spec", "clusterNetwork"), err))
		return allErrs
	}

	if clusterNetwork := ownerCluster.Spec.ClusterNetwork; clusterNetwork != nil {
		if clusterNetwork.Services != nil {
			// A user may provide zero or one CIDR blocks. If they provide an empty array,
			// we ignore it and use the default. AKS doesn't support > 1 Service/Pod CIDR.
			if len(clusterNetwork.Services.CIDRBlocks) > 1 {
				allErrs = append(allErrs, field.TooMany(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "cidrBlocks"), len(clusterNetwork.Services.CIDRBlocks), 1))
			}
			if len(clusterNetwork.Services.CIDRBlocks) == 1 {
				serviceCIDR = clusterNetwork.Services.CIDRBlocks[0]
			}
		}
		if clusterNetwork.Pods != nil {
			// A user may provide zero or one CIDR blocks. If they provide an empty array,
			// we ignore it and use the default. AKS doesn't support > 1 Service/Pod CIDR.
			if len(clusterNetwork.Pods.CIDRBlocks) > 1 {
				allErrs = append(allErrs, field.TooMany(field.NewPath("Cluster", "spec", "clusterNetwork", "pods", "cidrBlocks"), len(clusterNetwork.Pods.CIDRBlocks), 1))
			}
		}
	}

	if dnsServiceIP != nil {
		if serviceCIDR == "" {
			allErrs = append(allErrs, field.Required(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "cidrBlocks"), "service CIDR must be specified if specifying DNSServiceIP"))
		}
		_, cidr, err := net.ParseCIDR(serviceCIDR)
		if err != nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "cidrBlocks"), serviceCIDR, fmt.Sprintf("failed to parse cluster service cidr: %v", err)))
		}

		dnsIP := net.ParseIP(*dnsServiceIP)
		if dnsIP == nil { // dnsIP will be nil if the string is not a valid IP
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "dnsServiceIP"), *dnsServiceIP, "must be a valid IP address"))
		}

		if dnsIP != nil && !cidr.Contains(dnsIP) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "cidrBlocks"), serviceCIDR, "DNSServiceIP must reside within the associated cluster serviceCIDR"))
		}

		// AKS only supports .10 as the last octet for the DNSServiceIP.
		// Refer to: https://learn.microsoft.com/en-us/azure/aks/configure-kubenet#create-an-aks-cluster-with-system-assigned-managed-identities
		targetSuffix := ".10"
		if dnsIP != nil && !strings.HasSuffix(dnsIP.String(), targetSuffix) {
			allErrs = append(allErrs, field.Invalid(field.NewPath("Cluster", "spec", "clusterNetwork", "services", "dnsServiceIP"), *dnsServiceIP, fmt.Sprintf("must end with %q", targetSuffix)))
		}
	}

	if errs := validatePrivateEndpoints(subnet.PrivateEndpoints, []string{subnet.CIDRBlock}, fldPath.Child("VirtualNetwork.Subnet.PrivateEndpoints")); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	return allErrs
}

// validateAutoUpgradeProfile validates auto upgrade profile.
func (m *AzureManagedControlPlane) validateAutoUpgradeProfile(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if old.Spec.AutoUpgradeProfile != nil {
		if old.Spec.AutoUpgradeProfile.UpgradeChannel != nil && (m.Spec.AutoUpgradeProfile == nil || m.Spec.AutoUpgradeProfile.UpgradeChannel == nil) {
			// Prevent AutoUpgradeProfile.UpgradeChannel to be set to nil.
			// Unsetting the field is not allowed.
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("Spec", "AutoUpgradeProfile", "UpgradeChannel"),
					old.Spec.AutoUpgradeProfile.UpgradeChannel,
					"field cannot be set to nil, to disable auto upgrades set the channel to none."))
		}
	}
	return allErrs
}

// validateK8sVersionUpdate validates K8s version.
func (m *AzureManagedControlPlane) validateK8sVersionUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if hv := versions.GetHigherK8sVersion(m.Spec.Version, old.Spec.Version); hv != m.Spec.Version {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "version"),
			m.Spec.Version, "field version cannot be downgraded"),
		)
	}

	if old.Status.AutoUpgradeVersion != "" && m.Spec.Version != old.Spec.Version {
		if hv := versions.GetHigherK8sVersion(m.Spec.Version, old.Status.AutoUpgradeVersion); hv != m.Spec.Version {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "version"),
				m.Spec.Version, "version is auto-upgraded to "+old.Status.AutoUpgradeVersion+", cannot be downgraded"),
			)
		}
	}
	return allErrs
}

// validateAPIServerAccessProfileUpdate validates update to APIServerAccessProfile.
func (m *AzureManagedControlPlane) validateAPIServerAccessProfileUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	newAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	oldAPIServerAccessProfileNormalized := &APIServerAccessProfile{}
	if m.Spec.APIServerAccessProfile != nil {
		newAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			APIServerAccessProfileClassSpec: APIServerAccessProfileClassSpec{
				EnablePrivateCluster:           m.Spec.APIServerAccessProfile.EnablePrivateCluster,
				PrivateDNSZone:                 m.Spec.APIServerAccessProfile.PrivateDNSZone,
				EnablePrivateClusterPublicFQDN: m.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
			},
		}
	}
	if old.Spec.APIServerAccessProfile != nil {
		oldAPIServerAccessProfileNormalized = &APIServerAccessProfile{
			APIServerAccessProfileClassSpec: APIServerAccessProfileClassSpec{
				EnablePrivateCluster:           old.Spec.APIServerAccessProfile.EnablePrivateCluster,
				PrivateDNSZone:                 old.Spec.APIServerAccessProfile.PrivateDNSZone,
				EnablePrivateClusterPublicFQDN: old.Spec.APIServerAccessProfile.EnablePrivateClusterPublicFQDN,
			},
		}
	}

	if !reflect.DeepEqual(newAPIServerAccessProfileNormalized, oldAPIServerAccessProfileNormalized) {
		allErrs = append(allErrs,
			field.Invalid(field.NewPath("spec", "apiServerAccessProfile"),
				m.Spec.APIServerAccessProfile, "fields (except for AuthorizedIPRanges) are immutable"),
		)
	}

	return allErrs
}

// validateAddonProfilesUpdate validates update to AddonProfiles.
func (m *AzureManagedControlPlane) validateAddonProfilesUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	newAddonProfileMap := map[string]struct{}{}
	if len(old.Spec.AddonProfiles) != 0 {
		for _, addonProfile := range m.Spec.AddonProfiles {
			newAddonProfileMap[addonProfile.Name] = struct{}{}
		}
		for i, addonProfile := range old.Spec.AddonProfiles {
			if _, ok := newAddonProfileMap[addonProfile.Name]; !ok {
				allErrs = append(allErrs, field.Invalid(
					field.NewPath("spec", "addonProfiles"),
					m.Spec.AddonProfiles,
					fmt.Sprintf("cannot remove addonProfile %s, To disable this AddonProfile, update Spec.AddonProfiles[%v].Enabled to false", addonProfile.Name, i)))
			}
		}
	}
	return allErrs
}

// validateVirtualNetworkUpdate validates update to VirtualNetwork.
func (m *AzureManagedControlPlane) validateVirtualNetworkUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if old.Spec.VirtualNetwork.Name != m.Spec.VirtualNetwork.Name {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "virtualNetwork", "name"),
				m.Spec.VirtualNetwork.Name,
				"Virtual Network Name is immutable"))
	}

	if old.Spec.VirtualNetwork.CIDRBlock != m.Spec.VirtualNetwork.CIDRBlock {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "virtualNetwork", "cidrBlock"),
				m.Spec.VirtualNetwork.CIDRBlock,
				"Virtual Network CIDRBlock is immutable"))
	}

	if old.Spec.VirtualNetwork.Subnet.Name != m.Spec.VirtualNetwork.Subnet.Name {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "virtualNetwork", "subnet", "name"),
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
				field.NewPath("spec", "virtualNetwork", "subnet", "cidrBlock"),
				m.Spec.VirtualNetwork.Subnet.CIDRBlock,
				"Subnet CIDRBlock is immutable"))
	}

	if old.Spec.VirtualNetwork.ResourceGroup != m.Spec.VirtualNetwork.ResourceGroup {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "virtualNetwork", "resourceGroup"),
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
		allErrs = append(allErrs, field.Forbidden(field.NewPath("spec", "networkPluginMode"), fmt.Sprintf("%q NetworkPluginMode cannot be enabled when NetworkPolicy is set", NetworkPluginModeOverlay)))
	}

	return allErrs
}

// validateAADProfileUpdateAndLocalAccounts validates updates for AADProfile.
func (m *AzureManagedControlPlane) validateAADProfileUpdateAndLocalAccounts(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList
	if old.Spec.AADProfile != nil {
		if m.Spec.AADProfile == nil {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "aadProfile"),
					m.Spec.AADProfile,
					"field cannot be nil, cannot disable AADProfile"))
		} else {
			if !m.Spec.AADProfile.Managed && old.Spec.AADProfile.Managed {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("spec", "aadProfile", "managed"),
						m.Spec.AADProfile.Managed,
						"cannot set AADProfile.Managed to false"))
			}
			if len(m.Spec.AADProfile.AdminGroupObjectIDs) == 0 {
				allErrs = append(allErrs,
					field.Invalid(
						field.NewPath("spec", "aadProfile", "adminGroupObjectIDs"),
						m.Spec.AADProfile.AdminGroupObjectIDs,
						"length of AADProfile.AdminGroupObjectIDs cannot be zero"))
			}
		}
	}

	if old.Spec.DisableLocalAccounts == nil &&
		m.Spec.DisableLocalAccounts != nil &&
		m.Spec.AADProfile == nil {
		allErrs = append(allErrs,
			field.Invalid(
				field.NewPath("spec", "disableLocalAccounts"),
				m.Spec.DisableLocalAccounts,
				"DisableLocalAccounts can be set only for AAD enabled clusters"))
	}

	if old.Spec.DisableLocalAccounts != nil {
		// Prevent DisableLocalAccounts modification if it was already set to some value
		if err := webhookutils.ValidateImmutable(
			field.NewPath("spec", "disableLocalAccounts"),
			m.Spec.DisableLocalAccounts,
			old.Spec.DisableLocalAccounts,
		); err != nil {
			allErrs = append(allErrs, err)
		}
	}

	return allErrs
}

// validateSecurityProfileUpdate validates a SecurityProfile update.
func (m *AzureManagedControlPlaneClassSpec) validateSecurityProfileUpdate(old *AzureManagedControlPlaneClassSpec) field.ErrorList {
	var allErrs field.ErrorList
	if old.SecurityProfile != nil {
		if errAzureKeyVaultKms := m.validateAzureKeyVaultKmsUpdate(old); errAzureKeyVaultKms != nil {
			allErrs = append(allErrs, errAzureKeyVaultKms...)
		}
		if errWorkloadIdentity := m.validateWorkloadIdentityUpdate(old); errWorkloadIdentity != nil {
			allErrs = append(allErrs, errWorkloadIdentity...)
		}
		if errWorkloadIdentity := m.validateImageCleanerUpdate(old); errWorkloadIdentity != nil {
			allErrs = append(allErrs, errWorkloadIdentity...)
		}
		if errWorkloadIdentity := m.validateDefender(old); errWorkloadIdentity != nil {
			allErrs = append(allErrs, errWorkloadIdentity...)
		}
	}
	return allErrs
}

// validateAzureKeyVaultKmsUpdate validates AzureKeyVaultKmsUpdate profile.
func (m *AzureManagedControlPlaneClassSpec) validateAzureKeyVaultKmsUpdate(old *AzureManagedControlPlaneClassSpec) field.ErrorList {
	var allErrs field.ErrorList
	if old.SecurityProfile.AzureKeyVaultKms != nil {
		if m.SecurityProfile == nil || m.SecurityProfile.AzureKeyVaultKms == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "securityProfile", "azureKeyVaultKms"),
				nil, "cannot unset Spec.SecurityProfile.AzureKeyVaultKms profile to disable the profile please set Spec.SecurityProfile.AzureKeyVaultKms.Enabled to false"))
			return allErrs
		}
	}
	return allErrs
}

// validateWorkloadIdentityUpdate validates WorkloadIdentityUpdate profile.
func (m *AzureManagedControlPlaneClassSpec) validateWorkloadIdentityUpdate(old *AzureManagedControlPlaneClassSpec) field.ErrorList {
	var allErrs field.ErrorList
	if old.SecurityProfile.WorkloadIdentity != nil {
		if m.SecurityProfile == nil || m.SecurityProfile.WorkloadIdentity == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "securityProfile", "workloadIdentity"),
				nil, "cannot unset Spec.SecurityProfile.WorkloadIdentity, to disable workloadIdentity please set Spec.SecurityProfile.WorkloadIdentity.Enabled to false"))
		}
	}
	return allErrs
}

// validateImageCleanerUpdate validates ImageCleanerUpdate profile.
func (m *AzureManagedControlPlaneClassSpec) validateImageCleanerUpdate(old *AzureManagedControlPlaneClassSpec) field.ErrorList {
	var allErrs field.ErrorList
	if old.SecurityProfile.ImageCleaner != nil {
		if m.SecurityProfile == nil || m.SecurityProfile.ImageCleaner == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "securityProfile", "imageCleaner"),
				nil, "cannot unset Spec.SecurityProfile.ImageCleaner, to disable imageCleaner please set Spec.SecurityProfile.ImageCleaner.Enabled to false"))
		}
	}
	return allErrs
}

// validateDefender validates defender profile.
func (m *AzureManagedControlPlaneClassSpec) validateDefender(old *AzureManagedControlPlaneClassSpec) field.ErrorList {
	var allErrs field.ErrorList
	if old.SecurityProfile.Defender != nil {
		if m.SecurityProfile == nil || m.SecurityProfile.Defender == nil {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "securityProfile", "defender"),
				nil, "cannot unset Spec.SecurityProfile.Defender, to disable defender please set Spec.SecurityProfile.Defender.SecurityMonitoring.Enabled to false"))
		}
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
					field.NewPath("spec", "oidcIssuerProfile", "enabled"),
					"cannot be disabled",
				),
			)
		}
	}
	return allErrs
}

// validateFleetsMemberUpdate validates a FleetsMember.
func (m *AzureManagedControlPlane) validateFleetsMemberUpdate(old *AzureManagedControlPlane) field.ErrorList {
	var allErrs field.ErrorList

	if old.Spec.FleetsMember == nil || m.Spec.FleetsMember == nil {
		return allErrs
	}
	if old.Spec.FleetsMember.Name != "" && old.Spec.FleetsMember.Name != m.Spec.FleetsMember.Name {
		allErrs = append(allErrs,
			field.Forbidden(
				field.NewPath("spec", "fleetsMember", "name"),
				"Name is immutable",
			),
		)
	}

	return allErrs
}

// validateAKSExtensionsUpdate validates update to AKS extensions.
func validateAKSExtensionsUpdate(old []AKSExtension, current []AKSExtension) field.ErrorList {
	var allErrs field.ErrorList

	oldAKSExtensionsMap := make(map[string]AKSExtension, len(old))
	oldAKSExtensionsIndex := make(map[string]int, len(old))
	for i, extension := range old {
		oldAKSExtensionsMap[extension.Name] = extension
		oldAKSExtensionsIndex[extension.Name] = i
	}
	for i, extension := range current {
		oldExtension, ok := oldAKSExtensionsMap[extension.Name]
		if !ok {
			continue
		}
		if extension.Name != oldExtension.Name {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "name"),
					extension.Name,
					"field is immutable",
				),
			)
		}
		if (oldExtension.ExtensionType != nil && extension.ExtensionType != nil) && *extension.ExtensionType != *oldExtension.ExtensionType {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "extensionType"),
					extension.ExtensionType,
					"field is immutable",
				),
			)
		}
		if (extension.Plan != nil && oldExtension.Plan != nil) && *extension.Plan != *oldExtension.Plan {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "plan"),
					extension.Plan,
					"field is immutable",
				),
			)
		}
		if extension.Scope != oldExtension.Scope {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "scope"),
					extension.Scope,
					"field is immutable",
				),
			)
		}
		if (extension.ReleaseTrain != nil && oldExtension.ReleaseTrain != nil) && *extension.ReleaseTrain != *oldExtension.ReleaseTrain {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "releaseTrain"),
					extension.ReleaseTrain,
					"field is immutable",
				),
			)
		}
		if (extension.Version != nil && oldExtension.Version != nil) && *extension.Version != *oldExtension.Version {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "version"),
					extension.Version,
					"field is immutable",
				),
			)
		}
		if extension.Identity != oldExtension.Identity {
			allErrs = append(allErrs,
				field.Invalid(
					field.NewPath("spec", "extensions", fmt.Sprintf("[%d]", i), "identity"),
					extension.Identity,
					"field is immutable",
				),
			)
		}
	}

	return allErrs
}

func validateName(name string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if lName := strings.ToLower(name); strings.Contains(lName, "microsoft") ||
		strings.Contains(lName, "windows") {
		allErrs = append(allErrs, field.Invalid(fldPath.Child("Name"), name,
			"cluster name is invalid because 'MICROSOFT' and 'WINDOWS' can't be used as either a whole word or a substring in the name"))
	}

	return allErrs
}

// validateAKSExtensions validates the AKS extensions.
func validateAKSExtensions(extensions []AKSExtension, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	for _, extension := range extensions {
		if extension.Version != nil && (extension.AutoUpgradeMinorVersion == nil || (extension.AutoUpgradeMinorVersion != nil && *extension.AutoUpgradeMinorVersion)) {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("Version"), "Version must not be given if AutoUpgradeMinorVersion is true (or not provided, as it is true by default)"))
		}
		if extension.AutoUpgradeMinorVersion == ptr.To(false) && extension.ReleaseTrain != nil {
			allErrs = append(allErrs, field.Forbidden(fldPath.Child("ReleaseTrain"), "ReleaseTrain must not be given if AutoUpgradeMinorVersion is false"))
		}
		if extension.Scope != nil {
			if extension.Scope.ScopeType == ExtensionScopeCluster {
				if extension.Scope.ReleaseNamespace == "" {
					allErrs = append(allErrs, field.Required(fldPath.Child("Scope", "ReleaseNamespace"), "ReleaseNamespace must be provided if Scope is Cluster"))
				}
				if extension.Scope.TargetNamespace != "" {
					allErrs = append(allErrs, field.Forbidden(fldPath.Child("Scope", "TargetNamespace"), "TargetNamespace can only be given if Scope is Namespace"))
				}
			} else if extension.Scope.ScopeType == ExtensionScopeNamespace {
				if extension.Scope.TargetNamespace == "" {
					allErrs = append(allErrs, field.Required(fldPath.Child("Scope", "TargetNamespace"), "TargetNamespace must be provided if Scope is Namespace"))
				}
				if extension.Scope.ReleaseNamespace != "" {
					allErrs = append(allErrs, field.Forbidden(fldPath.Child("Scope", "ReleaseNamespace"), "ReleaseNamespace can only be given if Scope is Cluster"))
				}
			}
		}
	}

	return allErrs
}

// validateNetworkPolicy validates the networkPolicy.
func validateNetworkPolicy(networkPolicy *string, networkDataplane *NetworkDataplaneType, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if networkPolicy == nil {
		return nil
	}

	if *networkPolicy == "cilium" && networkDataplane != nil && *networkDataplane != NetworkDataplaneTypeCilium {
		allErrs = append(allErrs, field.Invalid(fldPath, networkPolicy, "cilium network policy can only be used with cilium network dataplane"))
	}

	return allErrs
}

// validateNetworkDataplane validates the NetworkDataplane.
func validateNetworkDataplane(networkDataplane *NetworkDataplaneType, networkPolicy *string, networkPluginMode *NetworkPluginMode, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if networkDataplane == nil {
		return nil
	}

	if *networkDataplane == NetworkDataplaneTypeCilium && (networkPluginMode == nil || *networkPluginMode != NetworkPluginModeOverlay) {
		allErrs = append(allErrs, field.Invalid(fldPath, networkDataplane, "cilium network dataplane can only be used with overlay network plugin mode"))
	}
	if *networkDataplane == NetworkDataplaneTypeCilium && (networkPolicy == nil || *networkPolicy != "cilium") {
		allErrs = append(allErrs, field.Invalid(fldPath, networkDataplane, "cilium dataplane requires network policy cilium."))
	}

	return allErrs
}

// validateAutoScalerProfile validates an AutoScalerProfile.
func validateAutoScalerProfile(autoScalerProfile *AutoScalerProfile, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList

	if autoScalerProfile == nil {
		return nil
	}

	if errs := validateIntegerStringGreaterThanZero(autoScalerProfile.MaxEmptyBulkDelete, fldPath, "MaxEmptyBulkDelete"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateIntegerStringGreaterThanZero(autoScalerProfile.MaxGracefulTerminationSec, fldPath, "MaxGracefulTerminationSec"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateMaxNodeProvisionTime(autoScalerProfile.MaxNodeProvisionTime, fldPath); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if autoScalerProfile.MaxTotalUnreadyPercentage != nil {
		val, err := strconv.Atoi(*autoScalerProfile.MaxTotalUnreadyPercentage)
		if err != nil || val < 0 || val > 100 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "autoscalerProfile", "maxTotalUnreadyPercentage"), autoScalerProfile.MaxTotalUnreadyPercentage, "invalid value"))
		}
	}

	if errs := validateNewPodScaleUpDelay(autoScalerProfile.NewPodScaleUpDelay, fldPath); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateIntegerStringGreaterThanZero(autoScalerProfile.OkTotalUnreadyCount, fldPath, "okTotalUnreadyCount"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScanInterval(autoScalerProfile.ScanInterval, fldPath); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScaleDownTime(autoScalerProfile.ScaleDownDelayAfterAdd, fldPath, "scaleDownDelayAfterAdd"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScaleDownDelayAfterDelete(autoScalerProfile.ScaleDownDelayAfterDelete, fldPath); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScaleDownTime(autoScalerProfile.ScaleDownDelayAfterFailure, fldPath, "scaleDownDelayAfterFailure"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScaleDownTime(autoScalerProfile.ScaleDownUnneededTime, fldPath, "scaleDownUnneededTime"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if errs := validateScaleDownTime(autoScalerProfile.ScaleDownUnreadyTime, fldPath, "scaleDownUnreadyTime"); len(errs) > 0 {
		allErrs = append(allErrs, errs...)
	}

	if autoScalerProfile.ScaleDownUtilizationThreshold != nil {
		val, err := strconv.ParseFloat(*autoScalerProfile.ScaleDownUtilizationThreshold, 32)
		if err != nil || val < 0 || val > 1 {
			allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "autoscalerProfile", "scaleDownUtilizationThreshold"), autoScalerProfile.ScaleDownUtilizationThreshold, "invalid value"))
		}
	}

	return allErrs
}

// validateMaxNodeProvisionTime validates update to AutoscalerProfile.MaxNodeProvisionTime.
func validateMaxNodeProvisionTime(maxNodeProvisionTime *string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(maxNodeProvisionTime, "") != "" {
		if !rMaxNodeProvisionTime.MatchString(ptr.Deref(maxNodeProvisionTime, "")) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("MaxNodeProvisionTime"), maxNodeProvisionTime, "invalid value"))
		}
	}
	return allErrs
}

// validateScanInterval validates update to AutoscalerProfile.ScanInterval.
func validateScanInterval(scanInterval *string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(scanInterval, "") != "" {
		if !rScanInterval.MatchString(ptr.Deref(scanInterval, "")) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ScanInterval"), scanInterval, "invalid value"))
		}
	}
	return allErrs
}

// validateNewPodScaleUpDelay validates update to AutoscalerProfile.NewPodScaleUpDelay.
func validateNewPodScaleUpDelay(newPodScaleUpDelay *string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(newPodScaleUpDelay, "") != "" {
		_, err := time.ParseDuration(ptr.Deref(newPodScaleUpDelay, ""))
		if err != nil {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("NewPodScaleUpDelay"), newPodScaleUpDelay, "invalid value"))
		}
	}
	return allErrs
}

// validateScaleDownDelayAfterDelete validates update to AutoscalerProfile.ScaleDownDelayAfterDelete value.
func validateScaleDownDelayAfterDelete(scaleDownDelayAfterDelete *string, fldPath *field.Path) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(scaleDownDelayAfterDelete, "") != "" {
		if !rScaleDownDelayAfterDelete.MatchString(ptr.Deref(scaleDownDelayAfterDelete, "")) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child("ScaleDownDelayAfterDelete"), ptr.Deref(scaleDownDelayAfterDelete, ""), "invalid value"))
		}
	}
	return allErrs
}

// validateScaleDownTime validates update to AutoscalerProfile.ScaleDown* values.
func validateScaleDownTime(scaleDownValue *string, fldPath *field.Path, fieldName string) field.ErrorList {
	var allErrs field.ErrorList
	if ptr.Deref(scaleDownValue, "") != "" {
		if !rScaleDownTime.MatchString(ptr.Deref(scaleDownValue, "")) {
			allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldName), ptr.Deref(scaleDownValue, ""), "invalid value"))
		}
	}
	return allErrs
}

// validateIntegerStringGreaterThanZero validates that a string value is an integer greater than zero.
func validateIntegerStringGreaterThanZero(input *string, fldPath *field.Path, fieldName string) field.ErrorList {
	var allErrs field.ErrorList

	if input != nil {
		val, err := strconv.Atoi(*input)
		if err != nil || val < 0 {
			allErrs = append(allErrs, field.Invalid(fldPath.Child(fieldName), input, "invalid value"))
		}
	}

	return allErrs
}

// validateIdentity validates an Identity.
func (m *AzureManagedControlPlane) validateIdentity(_ client.Client) field.ErrorList {
	var allErrs field.ErrorList

	if m.Spec.Identity != nil {
		if m.Spec.Identity.Type == ManagedControlPlaneIdentityTypeUserAssigned {
			if m.Spec.Identity.UserAssignedIdentityResourceID == "" {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "identity", "userAssignedIdentityResourceID"), m.Spec.Identity.UserAssignedIdentityResourceID, "cannot be empty if Identity.Type is UserAssigned"))
			}
		} else {
			if m.Spec.Identity.UserAssignedIdentityResourceID != "" {
				allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "identity", "userAssignedIdentityResourceID"), m.Spec.Identity.UserAssignedIdentityResourceID, "should be empty if Identity.Type is SystemAssigned"))
			}
		}
	}

	if len(allErrs) > 0 {
		return allErrs
	}

	return nil
}

// validateNetworkPluginMode validates a NetworkPluginMode.
func (m *AzureManagedControlPlane) validateNetworkPluginMode(_ client.Client) field.ErrorList {
	var allErrs field.ErrorList

	const kubenet = "kubenet"
	if ptr.Deref(m.Spec.NetworkPluginMode, "") == NetworkPluginModeOverlay &&
		ptr.Deref(m.Spec.NetworkPlugin, "") == kubenet {
		allErrs = append(allErrs, field.Invalid(field.NewPath("spec", "networkPluginMode"), m.Spec.NetworkPluginMode, fmt.Sprintf("cannot be set to %q when NetworkPlugin is %q", NetworkPluginModeOverlay, kubenet)))
	}

	if len(allErrs) > 0 {
		return allErrs
	}

	return nil
}

// isOIDCEnabled return true if OIDC issuer is enabled.
func (m *AzureManagedControlPlaneClassSpec) isOIDCEnabled() bool {
	if m.OIDCIssuerProfile == nil {
		return false
	}
	if m.OIDCIssuerProfile.Enabled == nil {
		return false
	}
	return *m.OIDCIssuerProfile.Enabled
}

// isUserManagedIdentityEnabled checks if user assigned identity is set.
func (m *AzureManagedControlPlaneClassSpec) isUserManagedIdentityEnabled() bool {
	if m.Identity == nil {
		return false
	}
	if m.Identity.Type != ManagedControlPlaneIdentityTypeUserAssigned {
		return false
	}
	return true
}
