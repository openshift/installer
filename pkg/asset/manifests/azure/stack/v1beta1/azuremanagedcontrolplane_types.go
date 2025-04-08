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

// UpgradeChannel determines the type of upgrade channel for automatically upgrading the cluster.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/en-us/azure/aks/auto-upgrade-cluster
type UpgradeChannel string

const (
	// UpgradeChannelNodeImage automatically upgrades the node image to the latest version available.
	// Consider using nodeOSUpgradeChannel instead as that allows you to configure node OS patching separate from Kubernetes version patching.
	UpgradeChannelNodeImage UpgradeChannel = "node-image"

	// UpgradeChannelNone disables auto-upgrades and keeps the cluster at its current version of Kubernetes.
	UpgradeChannelNone UpgradeChannel = "none"

	// UpgradeChannelPatch automatically upgrades the cluster to the latest supported patch version when it becomes available
	// while keeping the minor version the same. For example, if a cluster is running version 1.17.7 while versions 1.17.9, 1.18.4,
	// 1.18.6, and 1.19.1 are available, the cluster will be upgraded to 1.17.9.
	UpgradeChannelPatch UpgradeChannel = "patch"

	// UpgradeChannelRapid automatically upgrades the cluster to the latest supported patch release on the latest supported minor
	// version. In cases where the cluster is at a version of Kubernetes that is at an N-2 minor version where N is the latest
	// supported minor version, the cluster first upgrades to the latest supported patch version on N-1 minor version. For example,
	// if a cluster is running version 1.17.7 while versions 1.17.9, 1.18.4, 1.18.6, and 1.19.1 are available, the cluster
	// will first be upgraded to 1.18.6 and then to 1.19.1.
	UpgradeChannelRapid UpgradeChannel = "rapid"

	// UpgradeChannelStable automatically upgrade the cluster to the latest supported patch release on minor version N-1, where
	// N is the latest supported minor version. For example, if a cluster is running version 1.17.7 while versions 1.17.9, 1.18.4,
	// 1.18.6, and 1.19.1 are available, the cluster will be upgraded to 1.18.6.
	UpgradeChannelStable UpgradeChannel = "stable"
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

// NetworkDataplaneType is the type of network dataplane to use.
type NetworkDataplaneType string

const (
	// NetworkDataplaneTypeAzure is the Azure network dataplane type.
	NetworkDataplaneTypeAzure NetworkDataplaneType = "azure"
	// NetworkDataplaneTypeCilium is the Cilium network dataplane type.
	NetworkDataplaneTypeCilium NetworkDataplaneType = "cilium"
)

const (
	// LoadBalancerSKUStandard is the Standard load balancer SKU.
	LoadBalancerSKUStandard = "Standard"
	// LoadBalancerSKUBasic is the Basic load balancer SKU.
	LoadBalancerSKUBasic = "Basic"
)

// KeyVaultNetworkAccessTypes defines the types of network access of key vault.
// The possible values are Public and Private.
// The default value is Public.
type KeyVaultNetworkAccessTypes string

const (
	// KeyVaultNetworkAccessTypesPrivate means the key vault disables public access and enables private link.
	KeyVaultNetworkAccessTypesPrivate KeyVaultNetworkAccessTypes = "Private"

	// KeyVaultNetworkAccessTypesPublic means the key vault allows public access from all networks.
	KeyVaultNetworkAccessTypesPublic KeyVaultNetworkAccessTypes = "Public"
)

// AzureManagedControlPlaneSpec defines the desired state of AzureManagedControlPlane.
type AzureManagedControlPlaneSpec struct {
	AzureManagedControlPlaneClassSpec `json:",inline"`

	// NodeResourceGroupName is the name of the resource group
	// containing cluster IaaS resources. Will be populated to default
	// in webhook.
	// Immutable.
	// +optional
	NodeResourceGroupName string `json:"nodeResourceGroupName,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// Immutable, populated by the AKS API at create.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// SSHPublicKey is a string literal containing an ssh public key base64 encoded.
	// Use empty string to autogenerate new key. Use null value to not set key.
	// Immutable.
	// +optional
	SSHPublicKey *string `json:"sshPublicKey,omitempty"`

	// DNSPrefix allows the user to customize dns prefix.
	// Immutable.
	// +optional
	DNSPrefix *string `json:"dnsPrefix,omitempty"`

	// FleetsMember is the spec for the fleet this cluster is a member of.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/en-us/azure/templates/microsoft.containerservice/2023-03-15-preview/fleets/members
	// +optional
	FleetsMember *FleetsMember `json:"fleetsMember,omitempty"`
}

// ManagedClusterSecurityProfile defines the security profile for the cluster.
type ManagedClusterSecurityProfile struct {
	// AzureKeyVaultKms defines Azure Key Vault Management Services Profile for the security profile.
	// +optional
	AzureKeyVaultKms *AzureKeyVaultKms `json:"azureKeyVaultKms,omitempty"`

	// Defender settings for the security profile.
	// +optional
	Defender *ManagedClusterSecurityProfileDefender `json:"defender,omitempty"`

	// ImageCleaner settings for the security profile.
	// +optional
	ImageCleaner *ManagedClusterSecurityProfileImageCleaner `json:"imageCleaner,omitempty"`

	// Workloadidentity enables Kubernetes applications to access Azure cloud resources securely with Azure AD. Ensure to enable OIDC issuer while enabling Workload Identity
	// +optional
	WorkloadIdentity *ManagedClusterSecurityProfileWorkloadIdentity `json:"workloadIdentity,omitempty"`
}

// ManagedClusterSecurityProfileDefender defines Microsoft Defender settings for the security profile.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/defender-for-cloud/defender-for-containers-enable
type ManagedClusterSecurityProfileDefender struct {
	// LogAnalyticsWorkspaceResourceID is the ID of the Log Analytics workspace that has to be associated with Microsoft Defender.
	// When Microsoft Defender is enabled, this field is required and must be a valid workspace resource ID.
	// +kubebuilder:validation:Required
	LogAnalyticsWorkspaceResourceID string `json:"logAnalyticsWorkspaceResourceID"`

	// SecurityMonitoring profile defines the Microsoft Defender threat detection for Cloud settings for the security profile.
	// +kubebuilder:validation:Required
	SecurityMonitoring ManagedClusterSecurityProfileDefenderSecurityMonitoring `json:"securityMonitoring"`
}

// ManagedClusterSecurityProfileDefenderSecurityMonitoring settings for the security profile threat detection.
type ManagedClusterSecurityProfileDefenderSecurityMonitoring struct {
	// Enabled enables Defender threat detection
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`
}

// ManagedClusterSecurityProfileImageCleaner removes unused images from nodes, freeing up disk space and helping to reduce attack surface area.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/image-cleaner
type ManagedClusterSecurityProfileImageCleaner struct {
	// Enabled enables the Image Cleaner on AKS cluster.
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`

	// IntervalHours defines Image Cleaner scanning interval in hours. Default value is 24 hours.
	// +optional
	// +kubebuilder:validation:Minimum=24
	// +kubebuilder:validation:Maximum=2160
	IntervalHours *int `json:"intervalHours,omitempty"`
}

// ManagedClusterSecurityProfileWorkloadIdentity settings for the security profile.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/defender-for-cloud/defender-for-containers-enable
type ManagedClusterSecurityProfileWorkloadIdentity struct {
	// Enabled enables the workload identity.
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`
}

// AzureKeyVaultKms service settings for the security profile.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/use-kms-etcd-encryption#update-key-vault-mode
type AzureKeyVaultKms struct {
	// Enabled enables the Azure Key Vault key management service. The default is false.
	// +kubebuilder:validation:Required
	Enabled bool `json:"enabled"`

	// KeyID defines the Identifier of Azure Key Vault key.
	// When Azure Key Vault key management service is enabled, this field is required and must be a valid key identifier.
	// +kubebuilder:validation:Required
	KeyID string `json:"keyID"`

	// KeyVaultNetworkAccess defines the network access of key vault.
	// The possible values are Public and Private.
	// Public means the key vault allows public access from all networks.
	// Private means the key vault disables public access and enables private link. The default value is Public.
	// +optional
	// +kubebuilder:default:=Public
	KeyVaultNetworkAccess *KeyVaultNetworkAccessTypes `json:"keyVaultNetworkAccess,omitempty"`

	// KeyVaultResourceID is the Resource ID of key vault. When keyVaultNetworkAccess is Private, this field is required and must be a valid resource ID.
	// +optional
	KeyVaultResourceID *string `json:"keyVaultResourceID,omitempty"`
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
	// +kubebuilder:default:="Free"
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
	ManagedOutboundIPs *int `json:"managedOutboundIPs,omitempty"`

	// OutboundIPPrefixes - Desired outbound IP Prefix resources for the cluster load balancer.
	// +optional
	OutboundIPPrefixes []string `json:"outboundIPPrefixes,omitempty"`

	// OutboundIPs - Desired outbound IP resources for the cluster load balancer.
	// +optional
	OutboundIPs []string `json:"outboundIPs,omitempty"`

	// AllocatedOutboundPorts - Desired number of allocated SNAT ports per VM. Allowed values must be in the range of 0 to 64000 (inclusive). The default value is 0 which results in Azure dynamically allocating ports.
	// +optional
	AllocatedOutboundPorts *int `json:"allocatedOutboundPorts,omitempty"`

	// IdleTimeoutInMinutes - Desired outbound flow idle timeout in minutes. Allowed values must be in the range of 4 to 120 (inclusive). The default value is 30 minutes.
	// +optional
	IdleTimeoutInMinutes *int `json:"idleTimeoutInMinutes,omitempty"`
}

// APIServerAccessProfile tunes the accessibility of the cluster's control plane.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/azure/aks/api-server-authorized-ip-ranges
type APIServerAccessProfile struct {
	// AuthorizedIPRanges - Authorized IP Ranges to kubernetes API server.
	// +optional
	AuthorizedIPRanges []string `json:"authorizedIPRanges,omitempty"`

	APIServerAccessProfileClassSpec `json:",inline"`
}

// ManagedControlPlaneVirtualNetwork describes a virtual network required to provision AKS clusters.
type ManagedControlPlaneVirtualNetwork struct {
	// ResourceGroup is the name of the Azure resource group for the VNet and Subnet.
	// +optional
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// Name is the name of the virtual network.
	Name string `json:"name"`

	ManagedControlPlaneVirtualNetworkClassSpec `json:",inline"`
}

// ManagedControlPlaneSubnet describes a subnet for an AKS cluster.
type ManagedControlPlaneSubnet struct {
	Name string `json:"name"`

	// +kubebuilder:default:="10.240.0.0/16"
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
	// AutoUpgradeVersion is the Kubernetes version populated after auto-upgrade based on the upgrade channel.
	// +kubebuilder:validation:MinLength=2
	// +optional
	AutoUpgradeVersion string `json:"autoUpgradeVersion,omitempty"`

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

	// Version defines the Kubernetes version for the control plane instance.
	// +optional
	Version string `json:"version"`
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
// Default values are from https://learn.microsoft.com/azure/aks/cluster-autoscaler#using-the-autoscaler-profile
type AutoScalerProfile struct {
	// BalanceSimilarNodeGroups - Valid values are 'true' and 'false'. The default is false.
	// +kubebuilder:validation:Enum="true";"false"
	// +kubebuilder:default:="false"
	// +optional
	BalanceSimilarNodeGroups *BalanceSimilarNodeGroups `json:"balanceSimilarNodeGroups,omitempty"`
	// Expander - If not specified, the default is 'random'. See [expanders](https://github.com/kubernetes/autoscaler/blob/master/cluster-autoscaler/FAQ.md#what-are-expanders) for more information.
	// +kubebuilder:validation:Enum=least-waste;most-pods;priority;random
	// +kubebuilder:default:="random"
	// +optional
	Expander *Expander `json:"expander,omitempty"`
	// MaxEmptyBulkDelete - The default is 10.
	// +kubebuilder:default:="10"
	// +optional
	MaxEmptyBulkDelete *string `json:"maxEmptyBulkDelete,omitempty"`
	// MaxGracefulTerminationSec - The default is 600.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +kubebuilder:default:="600"
	// +optional
	MaxGracefulTerminationSec *string `json:"maxGracefulTerminationSec,omitempty"`
	// MaxNodeProvisionTime - The default is '15m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +kubebuilder:default:="15m"
	// +optional
	MaxNodeProvisionTime *string `json:"maxNodeProvisionTime,omitempty"`
	// MaxTotalUnreadyPercentage - The default is 45. The maximum is 100 and the minimum is 0.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +kubebuilder:validation:MaxLength=3
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:default:="45"
	// +optional
	MaxTotalUnreadyPercentage *string `json:"maxTotalUnreadyPercentage,omitempty"`
	// NewPodScaleUpDelay - For scenarios like burst/batch scale where you don't want CA to act before the kubernetes scheduler could schedule all the pods, you can tell CA to ignore unscheduled pods before they're a certain age. The default is '0s'. Values must be an integer followed by a unit ('s' for seconds, 'm' for minutes, 'h' for hours, etc).
	// +optional
	// +kubebuilder:default:="0s"
	NewPodScaleUpDelay *string `json:"newPodScaleUpDelay,omitempty"`
	// OkTotalUnreadyCount - This must be an integer. The default is 3.
	// +kubebuilder:validation:Pattern=`^(\d+)$`
	// +kubebuilder:default:="3"
	// +optional
	OkTotalUnreadyCount *string `json:"okTotalUnreadyCount,omitempty"`
	// ScanInterval - How often cluster is reevaluated for scale up or down. The default is '10s'.
	// +kubebuilder:validation:Pattern=`^(\d+)s$`
	// +kubebuilder:default:="10s"
	// +optional
	ScanInterval *string `json:"scanInterval,omitempty"`
	// ScaleDownDelayAfterAdd - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +kubebuilder:default:="10m"
	// +optional
	ScaleDownDelayAfterAdd *string `json:"scaleDownDelayAfterAdd,omitempty"`
	// ScaleDownDelayAfterDelete - The default is the scan-interval. Values must be an integer followed by an 's'. No unit of time other than seconds (s) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)s$`
	// +kubebuilder:default:="10s"
	// +optional
	ScaleDownDelayAfterDelete *string `json:"scaleDownDelayAfterDelete,omitempty"`
	// ScaleDownDelayAfterFailure - The default is '3m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +kubebuilder:default:="3m"
	// +optional
	ScaleDownDelayAfterFailure *string `json:"scaleDownDelayAfterFailure,omitempty"`
	// ScaleDownUnneededTime - The default is '10m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +kubebuilder:default:="10m"
	// +optional
	ScaleDownUnneededTime *string `json:"scaleDownUnneededTime,omitempty"`
	// ScaleDownUnreadyTime - The default is '20m'. Values must be an integer followed by an 'm'. No unit of time other than minutes (m) is supported.
	// +kubebuilder:validation:Pattern=`^(\d+)m$`
	// +kubebuilder:default:="20m"
	// +optional
	ScaleDownUnreadyTime *string `json:"scaleDownUnreadyTime,omitempty"`
	// ScaleDownUtilizationThreshold - The default is '0.5'.
	// +kubebuilder:default:="0.5"
	// +optional
	ScaleDownUtilizationThreshold *string `json:"scaleDownUtilizationThreshold,omitempty"`
	// SkipNodesWithLocalStorage - The default is false.
	// +kubebuilder:validation:Enum="true";"false"
	// +kubebuilder:default:="false"
	// +optional
	SkipNodesWithLocalStorage *SkipNodesWithLocalStorage `json:"skipNodesWithLocalStorage,omitempty"`
	// SkipNodesWithSystemPods - The default is true.
	// +kubebuilder:validation:Enum="true";"false"
	// +kubebuilder:default:="true"
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
	// +kubebuilder:default:=SystemAssigned
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
	// +kubebuilder:default:=false
	// +optional
	Enabled *bool `json:"enabled,omitempty"`
}

// AKSExtension represents the configuration for an AKS cluster extension.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/en-us/azure/aks/cluster-extensions
type AKSExtension struct {
	// Name is the name of the extension.
	Name string `json:"name"`

	// AKSAssignedIdentityType is the type of the AKS assigned identity.
	// +optional
	AKSAssignedIdentityType AKSAssignedIdentity `json:"aksAssignedIdentityType,omitempty"`

	// AutoUpgradeMinorVersion is a flag to note if this extension participates in auto upgrade of minor version, or not.
	// +kubebuilder:default=true
	// +optional
	AutoUpgradeMinorVersion *bool `json:"autoUpgradeMinorVersion,omitempty"`

	// ConfigurationSettings are the name-value pairs for configuring this extension.
	// +optional
	ConfigurationSettings map[string]string `json:"configurationSettings,omitempty"`

	// ExtensionType is the type of the Extension of which this resource is an instance.
	// It must be one of the Extension Types registered with Microsoft.KubernetesConfiguration by the Extension publisher.
	ExtensionType *string `json:"extensionType"`

	// Plan is the plan of the extension.
	// +optional
	Plan *ExtensionPlan `json:"plan,omitempty"`

	// ReleaseTrain is the release train this extension participates in for auto-upgrade (e.g. Stable, Preview, etc.)
	// This is only used if autoUpgradeMinorVersion is ‘true’.
	// +optional
	ReleaseTrain *string `json:"releaseTrain,omitempty"`

	// Scope is the scope at which this extension is enabled.
	// +optional
	Scope *ExtensionScope `json:"scope,omitempty"`

	// Version is the version of the extension.
	// +optional
	Version *string `json:"version,omitempty"`

	// Identity is the identity type of the Extension resource in an AKS cluster.
	// +optional
	Identity ExtensionIdentity `json:"identity,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this AzureManagedControlPlane belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].status"
// +kubebuilder:printcolumn:name="Severity",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].severity"
// +kubebuilder:printcolumn:name="Reason",type="string",JSONPath=".status.conditions[?(@.type=='Ready')].reason"
// +kubebuilder:printcolumn:name="Message",type="string",priority=1,JSONPath=".status.conditions[?(@.type=='Ready')].message"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of this AzureManagedControlPlane"
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
