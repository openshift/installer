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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

const (
	// ManagedClusterFinalizer allows Reconcile to clean up Azure resources associated with the AzureManagedControlPlane before
	// removing it from the apiserver.
	ManagedClusterFinalizer = "azuremanagedcontrolplane.infrastructure.cluster.x-k8s.io"

	// PrivateDNSZoneModeSystem represents mode System for azuremanagedcontrolplane.
	PrivateDNSZoneModeSystem string = "System"

	// PrivateDNSZoneModeNone represents mode None for azuremanagedcontrolplane.
	PrivateDNSZoneModeNone string = "None"
)

// ManagedControlPlaneOutboundType enumerates the values for the managed control plane OutboundType.
type ManagedControlPlaneOutboundType string

const (
	// ManagedControlPlaneOutboundTypeLoadBalancer ...
	ManagedControlPlaneOutboundTypeLoadBalancer ManagedControlPlaneOutboundType = "loadBalancer"
	// ManagedControlPlaneOutboundTypeManagedNATGateway ...
	ManagedControlPlaneOutboundTypeManagedNATGateway ManagedControlPlaneOutboundType = "managedNATGateway"
	// ManagedControlPlaneOutboundTypeUserAssignedNATGateway ...
	ManagedControlPlaneOutboundTypeUserAssignedNATGateway ManagedControlPlaneOutboundType = "userAssignedNATGateway"
	// ManagedControlPlaneOutboundTypeUserDefinedRouting ...
	ManagedControlPlaneOutboundTypeUserDefinedRouting ManagedControlPlaneOutboundType = "userDefinedRouting"
)

// ManagedControlPlaneIdentityType enumerates the values for managed control plane identity type.
type ManagedControlPlaneIdentityType string

const (
	// ManagedControlPlaneIdentityTypeSystemAssigned Use an implicitly created system-assigned managed identity to manage
	// cluster resources. Components in the control plane such as kube-controller-manager will use the
	// system-assigned managed identity to manipulate Azure resources.
	ManagedControlPlaneIdentityTypeSystemAssigned ManagedControlPlaneIdentityType = ManagedControlPlaneIdentityType(VMIdentitySystemAssigned)
	// ManagedControlPlaneIdentityTypeUserAssigned Use a user-assigned identity to manage cluster resources.
	// Components in the control plane such as kube-controller-manager will use the specified user-assigned
	// managed identity to manipulate Azure resources.
	ManagedControlPlaneIdentityTypeUserAssigned ManagedControlPlaneIdentityType = ManagedControlPlaneIdentityType(VMIdentityUserAssigned)
)

// NetworkPluginMode is the mode the network plugin should use.
type NetworkPluginMode string

const (
	// NetworkPluginModeOverlay is used with networkPlugin=azure, pods are given IPs from the PodCIDR address space but use Azure
	// Routing Domains rather than Kubenet's method of route tables.
	// See also [AKS doc].
	//
	// [AKS doc]: https://aka.ms/aks/azure-cni-overlay
	NetworkPluginModeOverlay NetworkPluginMode = "overlay"
)

// AzureManagedControlPlaneSpec defines the desired state of AzureManagedControlPlane.
type AzureManagedControlPlaneSpec struct {
	// Version defines the desired Kubernetes version.
	// +kubebuilder:validation:MinLength:=2
	Version string `json:"version"`

	// ResourceGroupName is the name of the Azure resource group for this AKS Cluster.
	// Immutable.
	ResourceGroupName string `json:"resourceGroupName"`

	// NodeResourceGroupName is the name of the resource group
	// containing cluster IaaS resources. Will be populated to default
	// in webhook.
	// Immutable.
	// +optional
	NodeResourceGroupName string `json:"nodeResourceGroupName,omitempty"`

	// VirtualNetwork describes the vnet for the AKS cluster. Will be created if it does not exist.
	// Immutable except for `subnet`.
	// +optional
	VirtualNetwork ManagedControlPlaneVirtualNetwork `json:"virtualNetwork,omitempty"`

	// SubscriptionID is the GUID of the Azure subscription to hold this cluster.
	// Immutable.
	// +optional
	SubscriptionID string `json:"subscriptionID,omitempty"`

	// Location is a string matching one of the canonical Azure region names. Examples: "westus2", "eastus".
	// Immutable.
	Location string `json:"location"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// Immutable, populated by the AKS API at create.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// AdditionalTags is an optional set of tags to add to Azure resources managed by the Azure provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// NetworkPlugin used for building Kubernetes network.
	// Allowed values are "azure", "kubenet".
	// Immutable.
	// +kubebuilder:validation:Enum=azure;kubenet
	// +optional
	NetworkPlugin *string `json:"networkPlugin,omitempty"`

	// NetworkPluginMode is the mode the network plugin should use.
	// Allowed value is "overlay".
	// +kubebuilder:validation:Enum=overlay
	// +optional
	NetworkPluginMode *NetworkPluginMode `json:"networkPluginMode,omitempty"`

	// NetworkPolicy used for building Kubernetes network.
	// Allowed values are "azure", "calico".
	// Immutable.
	// +kubebuilder:validation:Enum=azure;calico
	// +optional
	NetworkPolicy *string `json:"networkPolicy,omitempty"`

	// Outbound configuration used by Nodes.
	// Immutable.
	// +kubebuilder:validation:Enum=loadBalancer;managedNATGateway;userAssignedNATGateway;userDefinedRouting
	// +optional
	OutboundType *ManagedControlPlaneOutboundType `json:"outboundType,omitempty"`

	// SSHPublicKey is a string literal containing an ssh public key base64 encoded.
	// Use empty string to autogenerate new key. Use null value to not set key.
	// Immutable.
	// +optional
	SSHPublicKey *string `json:"sshPublicKey,omitempty"`

	// DNSServiceIP is an IP address assigned to the Kubernetes DNS service.
	// It must be within the Kubernetes service address range specified in serviceCidr.
	// Immutable.
	// +optional
	DNSServiceIP *string `json:"dnsServiceIP,omitempty"`

	// LoadBalancerSKU is the SKU of the loadBalancer to be provisioned.
	// Immutable.
	// +kubebuilder:validation:Enum=Basic;Standard
	// +optional
	LoadBalancerSKU *string `json:"loadBalancerSKU,omitempty"`

	// IdentityRef is a reference to a AzureClusterIdentity to be used when reconciling this cluster
	// +optional
	IdentityRef *corev1.ObjectReference `json:"identityRef,omitempty"`

	// AadProfile is Azure Active Directory configuration to integrate with AKS for aad authentication.
	// +optional
	AADProfile *AADProfile `json:"aadProfile,omitempty"`

	// AddonProfiles are the profiles of managed cluster add-on.
	// +optional
	AddonProfiles []AddonProfile `json:"addonProfiles,omitempty"`

	// SKU is the SKU of the AKS to be provisioned.
	// +optional
	SKU *AKSSku `json:"sku,omitempty"`

	// LoadBalancerProfile is the profile of the cluster load balancer.
	// +optional
	LoadBalancerProfile *LoadBalancerProfile `json:"loadBalancerProfile,omitempty"`

	// APIServerAccessProfile is the access profile for AKS API server.
	// Immutable except for `authorizedIPRanges`.
	// +optional
	APIServerAccessProfile *APIServerAccessProfile `json:"apiServerAccessProfile,omitempty"`

	// AutoscalerProfile is the parameters to be applied to the cluster-autoscaler when enabled
	// +optional
	AutoScalerProfile *AutoScalerProfile `json:"autoscalerProfile,omitempty"`

	// AzureEnvironment is the name of the AzureCloud to be used.
	// The default value that would be used by most users is "AzurePublicCloud", other values are:
	// - ChinaCloud: "AzureChinaCloud"
	// - PublicCloud: "AzurePublicCloud"
	// - USGovernmentCloud: "AzureUSGovernmentCloud"
	// +optional
	AzureEnvironment string `json:"azureEnvironment,omitempty"`

	// Identity configuration used by the AKS control plane.
	// +optional
	Identity *Identity `json:"identity,omitempty"`

	// KubeletUserAssignedIdentity is the user-assigned identity for kubelet.
	// For authentication with Azure Container Registry.
	// +optional
	KubeletUserAssignedIdentity string `json:"kubeletUserAssignedIdentity,omitempty"`

	// HTTPProxyConfig is the HTTP proxy configuration for the cluster.
	// Immutable.
	// +optional
	HTTPProxyConfig *HTTPProxyConfig `json:"httpProxyConfig,omitempty"`

	// OIDCIssuerProfile is the OIDC issuer profile of the Managed Cluster.
	// +optional
	OIDCIssuerProfile *OIDCIssuerProfile `json:"oidcIssuerProfile,omitempty"`

	// DNSPrefix allows the user to customize dns prefix.
	// Immutable.
	// +optional
	DNSPrefix *string `json:"dnsPrefix,omitempty"`
	// DisableLocalAccounts disables getting static credentials for this cluster when set. Expected to only be used for AAD clusters.
	// +optional
	DisableLocalAccounts *bool `json:"disableLocalAccounts,omitempty"`
}

// HTTPProxyConfig is the HTTP proxy configuration for the cluster.
type HTTPProxyConfig struct {
	// HTTPProxy is the HTTP proxy server endpoint to use.
	// +optional
	HTTPProxy *string `json:"httpProxy,omitempty"`

	// HTTPSProxy is the HTTPS proxy server endpoint to use.
	// +optional
	HTTPSProxy *string `json:"httpsProxy,omitempty"`

	// NoProxy indicates the endpoints that should not go through proxy.
	// +optional
	NoProxy []string `json:"noProxy,omitempty"`

	// TrustedCA is the alternative CA cert to use for connecting to proxy servers.
	// +optional
	TrustedCA *string `json:"trustedCa,omitempty"`
}

// AADProfile - AAD integration managed by AKS.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/managed-aad
type AADProfile struct {
	// Managed - Whether to enable managed AAD.
	// +kubebuilder:validation:Required
	Managed bool `json:"managed"`

	// AdminGroupObjectIDs - AAD group object IDs that will have admin role of the cluster.
	// +kubebuilder:validation:Required
	AdminGroupObjectIDs []string `json:"adminGroupObjectIDs"`
}

// AddonProfile represents a managed cluster add-on.
type AddonProfile struct {
	// Name - The name of the managed cluster add-on.
	Name string `json:"name"`

	// Config - Key-value pairs for configuring the add-on.
	// +optional
	Config map[string]string `json:"config,omitempty"`

	// Enabled - Whether the add-on is enabled or not.
	Enabled bool `json:"enabled"`
}

// AzureManagedControlPlaneSkuTier - Tier of a managed cluster SKU.
// +kubebuilder:validation:Enum=Free;Paid;Standard
type AzureManagedControlPlaneSkuTier string

const (
	// FreeManagedControlPlaneTier is the free tier of AKS without corresponding SLAs.
	FreeManagedControlPlaneTier AzureManagedControlPlaneSkuTier = "Free"
	// PaidManagedControlPlaneTier is the paid tier of AKS with corresponding SLAs.
	// Deprecated. It has been replaced with StandardManagedControlPlaneTier.
	PaidManagedControlPlaneTier AzureManagedControlPlaneSkuTier = "Paid"
	// StandardManagedControlPlaneTier is the standard tier of AKS with corresponding SLAs.
	StandardManagedControlPlaneTier AzureManagedControlPlaneSkuTier = "Standard"
)

// AKSSku - AKS SKU.
type AKSSku struct {
	// Tier - Tier of an AKS cluster.
	Tier AzureManagedControlPlaneSkuTier `json:"tier"`
}

// LoadBalancerProfile - Profile of the cluster load balancer.
// At most one of `managedOutboundIPs`, `outboundIPPrefixes`, or `outboundIPs` may be specified.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/load-balancer-standard
type LoadBalancerProfile struct {
	// ManagedOutboundIPs - Desired managed outbound IPs for the cluster load balancer.
	// +optional
	ManagedOutboundIPs *int32 `json:"managedOutboundIPs,omitempty"`

	// OutboundIPPrefixes - Desired outbound IP Prefix resources for the cluster load balancer.
	// +optional
	OutboundIPPrefixes []string `json:"outboundIPPrefixes,omitempty"`

	// OutboundIPs - Desired outbound IP resources for the cluster load balancer.
	// +optional
	OutboundIPs []string `json:"outboundIPs,omitempty"`

	// AllocatedOutboundPorts - Desired number of allocated SNAT ports per VM. Allowed values must be in the range of 0 to 64000 (inclusive). The default value is 0 which results in Azure dynamically allocating ports.
	// +optional
	AllocatedOutboundPorts *int32 `json:"allocatedOutboundPorts,omitempty"`

	// IdleTimeoutInMinutes - Desired outbound flow idle timeout in minutes. Allowed values must be in the range of 4 to 120 (inclusive). The default value is 30 minutes.
	// +optional
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
}

// APIServerAccessProfile tunes the accessibility of the cluster's control plane.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/api-server-authorized-ip-ranges
type APIServerAccessProfile struct {
	// AuthorizedIPRanges - Authorized IP Ranges to kubernetes API server.
	// +optional
	AuthorizedIPRanges []string `json:"authorizedIPRanges,omitempty"`
	// EnablePrivateCluster - Whether to create the cluster as a private cluster or not.
	// +optional
	EnablePrivateCluster *bool `json:"enablePrivateCluster,omitempty"`
	// PrivateDNSZone - Private dns zone mode for private cluster.
	// +kubebuilder:validation:Enum=System;None
	// +optional
	PrivateDNSZone *string `json:"privateDNSZone,omitempty"`
	// EnablePrivateClusterPublicFQDN - Whether to create additional public FQDN for private cluster or not.
	// +optional
	EnablePrivateClusterPublicFQDN *bool `json:"enablePrivateClusterPublicFQDN,omitempty"`
}

// ManagedControlPlaneVirtualNetwork describes a virtual network required to provision AKS clusters.
type ManagedControlPlaneVirtualNetwork struct {
	Name      string `json:"name"`
	CIDRBlock string `json:"cidrBlock"`
	// Immutable except for `serviceEndpoints`.
	// +optional
	Subnet ManagedControlPlaneSubnet `json:"subnet,omitempty"`
	// ResourceGroup is the name of the Azure resource group for the VNet and Subnet.
	// +optional
	ResourceGroup string `json:"resourceGroup,omitempty"`
}

// ManagedControlPlaneSubnet describes a subnet for an AKS cluster.
type ManagedControlPlaneSubnet struct {
	Name      string `json:"name"`
	CIDRBlock string `json:"cidrBlock"`

	// ServiceEndpoints is a slice of Virtual Network service endpoints to enable for the subnets.
	// +optional
	ServiceEndpoints ServiceEndpoints `json:"serviceEndpoints,omitempty"`

	// PrivateEndpoints is a slice of Virtual Network private endpoints to create for the subnets.
	// +optional
	PrivateEndpoints PrivateEndpoints `json:"privateEndpoints,omitempty"`
}

// AzureManagedControlPlaneStatus defines the observed state of AzureManagedControlPlane.
type AzureManagedControlPlaneStatus struct {
	// Ready is true when the provider resource is ready.
	// +optional
	Ready bool `json:"ready,omitempty"`

	// Initialized is true when the control plane is available for initial contact.
	// This may occur before the control plane is fully ready.
	// In the AzureManagedControlPlane implementation, these are identical.
	// +optional
	Initialized bool `json:"initialized,omitempty"`

	// Conditions defines current service state of the AzureManagedControlPlane.
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`

	// LongRunningOperationStates saves the states for Azure long-running operations so they can be continued on the
	// next reconciliation loop.
	// +optional
	LongRunningOperationStates Futures `json:"longRunningOperationStates,omitempty"`

	// OIDCIssuerProfile is the OIDC issuer profile of the Managed Cluster.
	// +optional
	OIDCIssuerProfile *OIDCIssuerProfileStatus `json:"oidcIssuerProfile,omitempty"`
}

// OIDCIssuerProfileStatus is the OIDC issuer profile of the Managed Cluster.
type OIDCIssuerProfileStatus struct {
	// IssuerURL is the OIDC issuer url of the Managed Cluster.
	// +optional
	IssuerURL *string `json:"issuerURL,omitempty"`
}

// AutoScalerProfile parameters to be applied to the cluster-autoscaler.
// See also [AKS doc], [K8s doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/cluster-autoscaler#use-the-cluster-autoscaler-profile
// [K8s doc]: https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-the-parameters-to-ca
type AutoScalerProfile struct {
	// BalanceSimilarNodeGroups - Valid values are 'true' and 'false'. The default is false.
	// +kubebuilder:validation:Enum="true";"false"
	// +optional
	BalanceSimilarNodeGroups *BalanceSimilarNodeGroups `json:"balanceSimilarNodeGroups,omitempty"`
	// Expander - If not specified, the default is 'random'. See [expanders](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-expanders) for more information.
	// +kubebuilder:validation:Enum=least-waste;most-pods;priority;random
	// +optional
	Expander *Expander `json:"expander,omitempty"`
	// MaxEmptyBulkDelete - The default is 10.
	// +optional
	MaxEmptyBulkDelete *string `json:"maxEmptyBulkDelete,omitempty"`
	// MaxGracefulTerminationSec - The default is 600.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +optional
	MaxGracefulTerminationSec *string `json:"maxGracefulTerminationSec,omitempty"`
	// MaxNodeProvisionTime - The default is '15m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +optional
	MaxNodeProvisionTime *string `json:"maxNodeProvisionTime,omitempty"`
	// MaxTotalUnreadyPercentage - The default is 45. The maximum is 100 and the minimum is 0.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +kubebuilder:validation:MaxLength=3
	// +kubebuilder:validation:MinLength=1
	// +optional
	MaxTotalUnreadyPercentage *string `json:"maxTotalUnreadyPercentage,omitempty"`
	// NewPodScaleUpDelay - For scenarios like burst/batch scale where you don't want CA to act before the kubernetes scheduler could schedule all the pods, you can tell CA to ignore unscheduled pods before they're a certain age. The default is '0s'. Values must be an integer followed by a unit ('s' for seconds, 'm' for minutes, 'h' for hours, etc).
	// +optional
	NewPodScaleUpDelay *string `json:"newPodScaleUpDelay,omitempty"`
	// OkTotalUnreadyCount - This must be an integer. The default is 3.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +optional
	OkTotalUnreadyCount *string `json:"okTotalUnreadyCount,omitempty"`
	// ScanInterval - How often cluster is reevaluated for scale up or down. The default is '10s'.
	// +kubebuilder:validation:Pattern=`^(\d+)s$`
	// +optional
	ScanInterval *string `json:"scanInterval,omitempty"`
	// ScaleDownDelayAfterAdd - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +optional
	ScaleDownDelayAfterAdd *string `json:"scaleDownDelayAfterAdd,omitempty"`
	// ScaleDownDelayAfterDelete - The default is the scan-interval. Values must be an integer followed by an 's'. No unit of time other than seconds (s) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)s$`
	// +optional
	ScaleDownDelayAfterDelete *string `json:"scaleDownDelayAfterDelete,omitempty"`
	// ScaleDownDelayAfterFailure - The default is '3m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +optional
	ScaleDownDelayAfterFailure *string `json:"scaleDownDelayAfterFailure,omitempty"`
	// ScaleDownUnneededTime - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +optional
	ScaleDownUnneededTime *string `json:"scaleDownUnneededTime,omitempty"`
	// ScaleDownUnreadyTime - The default is '20m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +optional
	ScaleDownUnreadyTime *string `json:"scaleDownUnreadyTime,omitempty"`
	// ScaleDownUtilizationThreshold - The default is '0.5'.
	// +optional
	ScaleDownUtilizationThreshold *string `json:"scaleDownUtilizationThreshold,omitempty"`
	// SkipNodesWithLocalStorage - The default is false.
	// +kubebuilder:validation:Enum="true";"false"
	// +optional
	SkipNodesWithLocalStorage *SkipNodesWithLocalStorage `json:"skipNodesWithLocalStorage,omitempty"`
	// SkipNodesWithSystemPods - The default is true.
	// +kubebuilder:validation:Enum="true";"false"
	// +optional
	SkipNodesWithSystemPods *SkipNodesWithSystemPods `json:"skipNodesWithSystemPods,omitempty"`
}

// BalanceSimilarNodeGroups enumerates the values for BalanceSimilarNodeGroups.
type BalanceSimilarNodeGroups string

const (
	// BalanceSimilarNodeGroupsTrue ...
	BalanceSimilarNodeGroupsTrue BalanceSimilarNodeGroups = "true"
	// BalanceSimilarNodeGroupsFalse ...
	BalanceSimilarNodeGroupsFalse BalanceSimilarNodeGroups = "false"
)

// SkipNodesWithLocalStorage enumerates the values for SkipNodesWithLocalStorage.
type SkipNodesWithLocalStorage string

const (
	// SkipNodesWithLocalStorageTrue ...
	SkipNodesWithLocalStorageTrue SkipNodesWithLocalStorage = "true"
	// SkipNodesWithLocalStorageFalse ...
	SkipNodesWithLocalStorageFalse SkipNodesWithLocalStorage = "false"
)

// SkipNodesWithSystemPods enumerates the values for SkipNodesWithSystemPods.
type SkipNodesWithSystemPods string

const (
	// SkipNodesWithSystemPodsTrue ...
	SkipNodesWithSystemPodsTrue SkipNodesWithSystemPods = "true"
	// SkipNodesWithSystemPodsFalse ...
	SkipNodesWithSystemPodsFalse SkipNodesWithSystemPods = "false"
)

// Expander enumerates the values for Expander.
type Expander string

const (
	// ExpanderLeastWaste ...
	ExpanderLeastWaste Expander = "least-waste"
	// ExpanderMostPods ...
	ExpanderMostPods Expander = "most-pods"
	// ExpanderPriority ...
	ExpanderPriority Expander = "priority"
	// ExpanderRandom ...
	ExpanderRandom Expander = "random"
)

// Identity represents the Identity configuration for an AKS control plane.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/en-us/azure/aks/use-managed-identity
type Identity struct {
	// Type - The Identity type to use.
	// +kubebuilder:validation:Enum=SystemAssigned;UserAssigned
	// +optional
	Type ManagedControlPlaneIdentityType `json:"type,omitempty"`

	// UserAssignedIdentityResourceID - Identity ARM resource ID when using user-assigned identity.
	// +optional
	UserAssignedIdentityResourceID string `json:"userAssignedIdentityResourceID,omitempty"`
}

// OIDCIssuerProfile is the OIDC issuer profile of the Managed Cluster.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/en-us/azure/aks/use-oidc-issuer
type OIDCIssuerProfile struct {
	// Enabled is whether the OIDC issuer is enabled.
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=azuremanagedcontrolplanes,scope=Namespaced,categories=cluster-api,shortName=amcp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// AzureManagedControlPlane is the Schema for the azuremanagedcontrolplanes API.
type AzureManagedControlPlane struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   AzureManagedControlPlaneSpec   `json:"spec,omitempty"`
	Status AzureManagedControlPlaneStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// AzureManagedControlPlaneList contains a list of AzureManagedControlPlane.
type AzureManagedControlPlaneList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AzureManagedControlPlane `json:"items"`
}

// GetConditions returns the list of conditions for an AzureManagedControlPlane API object.
func (m *AzureManagedControlPlane) GetConditions() clusterv1.Conditions {
	return m.Status.Conditions
}

// SetConditions will set the given conditions on an AzureManagedControlPlane object.
func (m *AzureManagedControlPlane) SetConditions(conditions clusterv1.Conditions) {
	m.Status.Conditions = conditions
}

// GetFutures returns the list of long running operation states for an AzureManagedControlPlane API object.
func (m *AzureManagedControlPlane) GetFutures() Futures {
	return m.Status.LongRunningOperationStates
}

// SetFutures will set the given long running operation states on an AzureManagedControlPlane object.
func (m *AzureManagedControlPlane) SetFutures(futures Futures) {
	m.Status.LongRunningOperationStates = futures
}

func init() {
	SchemeBuilder.Register(&AzureManagedControlPlane{}, &AzureManagedControlPlaneList{})
}
