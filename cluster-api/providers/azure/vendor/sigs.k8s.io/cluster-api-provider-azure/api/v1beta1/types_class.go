/*
Copyright 2022 The Kubernetes Authors.

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
	"k8s.io/apimachinery/pkg/api/resource"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// AzureClusterClassSpec defines the AzureCluster properties that may be shared across several Azure clusters.
type AzureClusterClassSpec struct {
	// +optional
	SubscriptionID string `json:"subscriptionID,omitempty"`

	Location string `json:"location"`

	// ExtendedLocation is an optional set of ExtendedLocation properties for clusters on Azure public MEC.
	// +optional
	ExtendedLocation *ExtendedLocationSpec `json:"extendedLocation,omitempty"`

	// AdditionalTags is an optional set of tags to add to Azure resources managed by the Azure provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// IdentityRef is a reference to an AzureIdentity to be used when reconciling this cluster
	// +optional
	IdentityRef *corev1.ObjectReference `json:"identityRef,omitempty"`

	// AzureEnvironment is the name of the AzureCloud to be used.
	// The default value that would be used by most users is "AzurePublicCloud", other values are:
	// - ChinaCloud: "AzureChinaCloud"
	// - GermanCloud: "AzureGermanCloud"
	// - PublicCloud: "AzurePublicCloud"
	// - USGovernmentCloud: "AzureUSGovernmentCloud"
	//
	// Note that values other than the default must also be accompanied by corresponding changes to the
	// aso-controller-settings Secret to configure ASO to refer to the non-Public cloud. ASO currently does
	// not support referring to multiple different clouds in a single installation. The following fields must
	// be defined in the Secret:
	// - AZURE_AUTHORITY_HOST
	// - AZURE_RESOURCE_MANAGER_ENDPOINT
	// - AZURE_RESOURCE_MANAGER_AUDIENCE
	//
	// See the [ASO docs] for more details.
	//
	// [ASO docs]: https://azure.github.io/azure-service-operator/guide/aso-controller-settings-options/
	// +optional
	AzureEnvironment string `json:"azureEnvironment,omitempty"`

	// CloudProviderConfigOverrides is an optional set of configuration values that can be overridden in azure cloud provider config.
	// This is only a subset of options that are available in azure cloud provider config.
	// Some values for the cloud provider config are inferred from other parts of cluster api provider azure spec, and may not be available for overrides.
	// See: https://cloud-provider-azure.sigs.k8s.io/install/configs
	// Note: All cloud provider config values can be customized by creating the secret beforehand. CloudProviderConfigOverrides is only used when the secret is managed by the Azure Provider.
	// +optional
	CloudProviderConfigOverrides *CloudProviderConfigOverrides `json:"cloudProviderConfigOverrides,omitempty"`

	// FailureDomains is a list of failure domains in the cluster's region, used to restrict
	// eligibility to host the control plane. A FailureDomain maps to an availability zone,
	// which is a separated group of datacenters within a region.
	// See: https://learn.microsoft.com/azure/reliability/availability-zones-overview
	// +optional
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`
}

// AzureManagedControlPlaneClassSpec defines the AzureManagedControlPlane properties that may be shared across several azure managed control planes.
type AzureManagedControlPlaneClassSpec struct {
	// MachineTemplate contains information about how machines
	// should be shaped when creating or updating a control plane.
	// For the AzureManagedControlPlaneTemplate, this field is used
	// only to fulfill the CAPI contract.
	// +optional
	MachineTemplate *AzureManagedControlPlaneTemplateMachineTemplate `json:"machineTemplate,omitempty"`

	// ResourceGroupName is the name of the Azure resource group for this AKS Cluster.
	// Immutable.
	ResourceGroupName string `json:"resourceGroupName"`

	// Version defines the desired Kubernetes version.
	// +kubebuilder:validation:MinLength:=2
	Version string `json:"version"`

	// VirtualNetwork describes the virtual network for the AKS cluster. It will be created if it does not already exist.
	// +optional
	VirtualNetwork ManagedControlPlaneVirtualNetwork `json:"virtualNetwork,omitempty"`

	// SubscriptionID is the GUID of the Azure subscription that owns this cluster.
	// +optional
	SubscriptionID string `json:"subscriptionID,omitempty"`

	// Location is a string matching one of the canonical Azure region names. Examples: "westus2", "eastus".
	Location string `json:"location"`

	// AdditionalTags is an optional set of tags to add to Azure resources managed by the Azure provider, in addition to the
	// ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// NetworkPlugin used for building Kubernetes network.
	// +kubebuilder:validation:Enum=azure;kubenet;none
	// +kubebuilder:default:=azure
	// +optional
	NetworkPlugin *string `json:"networkPlugin,omitempty"`

	// NetworkPluginMode is the mode the network plugin should use.
	// Allowed value is "overlay".
	// +kubebuilder:validation:Enum=overlay
	// +optional
	NetworkPluginMode *NetworkPluginMode `json:"networkPluginMode,omitempty"`

	// NetworkPolicy used for building Kubernetes network.
	// +kubebuilder:validation:Enum=azure;calico;cilium
	// +optional
	NetworkPolicy *string `json:"networkPolicy,omitempty"`

	// NetworkDataplane is the dataplane used for building the Kubernetes network.
	// +kubebuilder:validation:Enum=azure;cilium
	// +optional
	NetworkDataplane *NetworkDataplaneType `json:"networkDataplane,omitempty"`

	// Outbound configuration used by Nodes.
	// +kubebuilder:validation:Enum=loadBalancer;managedNATGateway;userAssignedNATGateway;userDefinedRouting
	// +optional
	OutboundType *ManagedControlPlaneOutboundType `json:"outboundType,omitempty"`

	// DNSServiceIP is an IP address assigned to the Kubernetes DNS service.
	// It must be within the Kubernetes service address range specified in serviceCidr.
	// Immutable.
	// +optional
	DNSServiceIP *string `json:"dnsServiceIP,omitempty"`

	// LoadBalancerSKU is the SKU of the loadBalancer to be provisioned.
	// Immutable.
	// +kubebuilder:validation:Enum=Basic;Standard
	// +kubebuilder:default:=Standard
	// +optional
	LoadBalancerSKU *string `json:"loadBalancerSKU,omitempty"`

	// IdentityRef is a reference to a AzureClusterIdentity to be used when reconciling this cluster
	IdentityRef *corev1.ObjectReference `json:"identityRef"`

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
	//
	// Note that values other than the default must also be accompanied by corresponding changes to the
	// aso-controller-settings Secret to configure ASO to refer to the non-Public cloud. ASO currently does
	// not support referring to multiple different clouds in a single installation. The following fields must
	// be defined in the Secret:
	// - AZURE_AUTHORITY_HOST
	// - AZURE_RESOURCE_MANAGER_ENDPOINT
	// - AZURE_RESOURCE_MANAGER_AUDIENCE
	//
	// See the [ASO docs] for more details.
	//
	// [ASO docs]: https://azure.github.io/azure-service-operator/guide/aso-controller-settings-options/
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

	// DisableLocalAccounts disables getting static credentials for this cluster when set. Expected to only be used for AAD clusters.
	// +optional
	DisableLocalAccounts *bool `json:"disableLocalAccounts,omitempty"`

	// FleetsMember is the spec for the fleet this cluster is a member of.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/en-us/azure/templates/microsoft.containerservice/2023-03-15-preview/fleets/members
	// +optional
	FleetsMember *FleetsMemberClassSpec `json:"fleetsMember,omitempty"`

	// Extensions is a list of AKS extensions to be installed on the cluster.
	// +optional
	Extensions []AKSExtension `json:"extensions,omitempty"`

	// AutoUpgradeProfile defines the auto upgrade configuration.
	// +optional
	AutoUpgradeProfile *ManagedClusterAutoUpgradeProfile `json:"autoUpgradeProfile,omitempty"`

	// SecurityProfile defines the security profile for cluster.
	// +optional
	SecurityProfile *ManagedClusterSecurityProfile `json:"securityProfile,omitempty"`

	// ASOManagedClusterPatches defines JSON merge patches to be applied to the generated ASO ManagedCluster resource.
	// WARNING: This is meant to be used sparingly to enable features for development and testing that are not
	// otherwise represented in the CAPZ API. Misconfiguration that conflicts with CAPZ's normal mode of
	// operation is possible.
	// +optional
	ASOManagedClusterPatches []string `json:"asoManagedClusterPatches,omitempty"`

	// EnablePreviewFeatures enables preview features for the cluster.
	// +kubebuilder:default:=false
	// +optional
	EnablePreviewFeatures *bool `json:"enablePreviewFeatures,omitempty"`
}

// ManagedClusterAutoUpgradeProfile defines the auto upgrade profile for a managed cluster.
type ManagedClusterAutoUpgradeProfile struct {
	// UpgradeChannel determines the type of upgrade channel for automatically upgrading the cluster.
	// +kubebuilder:validation:Enum=node-image;none;patch;rapid;stable
	// +optional
	UpgradeChannel *UpgradeChannel `json:"upgradeChannel,omitempty"`
}

// AzureManagedMachinePoolClassSpec defines the AzureManagedMachinePool properties that may be shared across several Azure managed machinepools.
type AzureManagedMachinePoolClassSpec struct {
	// AdditionalTags is an optional set of tags to add to Azure resources managed by the
	// Azure provider, in addition to the ones added by default.
	// +optional
	AdditionalTags Tags `json:"additionalTags,omitempty"`

	// Name is the name of the agent pool. If not specified, CAPZ uses the name of the CR as the agent pool name.
	// Immutable.
	// +optional
	Name *string `json:"name,omitempty"`

	// Mode represents the mode of an agent pool. Possible values include: System, User.
	// +kubebuilder:validation:Enum=System;User
	Mode string `json:"mode"`

	// SKU is the size of the VMs in the node pool.
	// Immutable.
	SKU string `json:"sku"`

	// OSDiskSizeGB is the disk size for every machine in this agent pool.
	// If you specify 0, it will apply the default osDisk size according to the vmSize specified.
	// Immutable.
	// +optional
	OSDiskSizeGB *int `json:"osDiskSizeGB,omitempty"`

	// AvailabilityZones - Availability zones for nodes. Must use VirtualMachineScaleSets AgentPoolType.
	// Immutable.
	// +optional
	AvailabilityZones []string `json:"availabilityZones,omitempty"`

	// Node labels represent the labels for all of the nodes present in node pool.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/azure/aks/use-labels
	// +optional
	NodeLabels map[string]string `json:"nodeLabels,omitempty"`

	// Taints specifies the taints for nodes present in this agent pool.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/azure/aks/use-multiple-node-pools#setting-node-pool-taints
	// +optional
	Taints Taints `json:"taints,omitempty"`

	// Scaling specifies the autoscaling parameters for the node pool.
	// +optional
	Scaling *ManagedMachinePoolScaling `json:"scaling,omitempty"`

	// MaxPods specifies the kubelet `--max-pods` configuration for the node pool.
	// Immutable.
	// See also [AKS doc], [K8s doc].
	//
	// [AKS doc]: https://learn.microsoft.com/azure/aks/configure-azure-cni#configure-maximum---new-clusters
	// [K8s doc]: https://kubernetes.io/docs/reference/command-line-tools-reference/kubelet/
	// +optional
	MaxPods *int `json:"maxPods,omitempty"`

	// OsDiskType specifies the OS disk type for each node in the pool. Allowed values are 'Ephemeral' and 'Managed' (default).
	// Immutable.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/azure/aks/cluster-configuration#ephemeral-os
	// +kubebuilder:validation:Enum=Ephemeral;Managed
	// +kubebuilder:default=Managed
	// +optional
	OsDiskType *string `json:"osDiskType,omitempty"`

	// EnableUltraSSD enables the storage type UltraSSD_LRS for the agent pool.
	// Immutable.
	// +optional
	EnableUltraSSD *bool `json:"enableUltraSSD,omitempty"`

	// OSType specifies the virtual machine operating system. Default to Linux. Possible values include: 'Linux', 'Windows'.
	// 'Windows' requires the AzureManagedControlPlane's `spec.networkPlugin` to be `azure`.
	// Immutable.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/rest/api/aks/agent-pools/create-or-update?tabs=HTTP#ostype
	// +kubebuilder:validation:Enum=Linux;Windows
	// +kubebuilder:default:=Linux
	// +optional
	OSType *string `json:"osType,omitempty"`

	// EnableNodePublicIP controls whether or not nodes in the pool each have a public IP address.
	// Immutable.
	// +optional
	EnableNodePublicIP *bool `json:"enableNodePublicIP,omitempty"`

	// NodePublicIPPrefixID specifies the public IP prefix resource ID which VM nodes should use IPs from.
	// Immutable.
	// +optional
	NodePublicIPPrefixID *string `json:"nodePublicIPPrefixID,omitempty"`

	// ScaleSetPriority specifies the ScaleSetPriority value. Default to Regular. Possible values include: 'Regular', 'Spot'
	// Immutable.
	// +kubebuilder:validation:Enum=Regular;Spot
	// +optional
	ScaleSetPriority *string `json:"scaleSetPriority,omitempty"`

	// ScaleDownMode affects the cluster autoscaler behavior. Default to Delete. Possible values include: 'Deallocate', 'Delete'
	// +kubebuilder:validation:Enum=Deallocate;Delete
	// +kubebuilder:default=Delete
	// +optional
	ScaleDownMode *string `json:"scaleDownMode,omitempty"`

	// SpotMaxPrice defines max price to pay for spot instance. Possible values are any decimal value greater than zero or -1.
	// If you set the max price to be -1, the VM won't be evicted based on price. The price for the VM will be the current price
	// for spot or the price for a standard VM, which ever is less, as long as there's capacity and quota available.
	// +optional
	SpotMaxPrice *resource.Quantity `json:"spotMaxPrice,omitempty"`

	// KubeletConfig specifies the kubelet configurations for nodes.
	// Immutable.
	// +optional
	KubeletConfig *KubeletConfig `json:"kubeletConfig,omitempty"`

	// KubeletDiskType specifies the kubelet disk type. Default to OS. Possible values include: 'OS', 'Temporary'.
	// Requires Microsoft.ContainerService/KubeletDisk preview feature to be set.
	// Immutable.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/rest/api/aks/agent-pools/create-or-update?tabs=HTTP#kubeletdisktype
	// +kubebuilder:validation:Enum=OS;Temporary
	// +optional
	KubeletDiskType *KubeletDiskType `json:"kubeletDiskType,omitempty"`

	// LinuxOSConfig specifies the custom Linux OS settings and configurations.
	// Immutable.
	// +optional
	LinuxOSConfig *LinuxOSConfig `json:"linuxOSConfig,omitempty"`

	// SubnetName specifies the Subnet where the MachinePool will be placed
	// Immutable.
	// +optional
	SubnetName *string `json:"subnetName,omitempty"`

	// EnableFIPS indicates whether FIPS is enabled on the node pool.
	// Immutable.
	// +optional
	EnableFIPS *bool `json:"enableFIPS,omitempty"`

	// EnableEncryptionAtHost indicates whether host encryption is enabled on the node pool.
	// Immutable.
	// See also [AKS doc].
	//
	// [AKS doc]: https://learn.microsoft.com/en-us/azure/aks/enable-host-encryption
	// +optional
	EnableEncryptionAtHost *bool `json:"enableEncryptionAtHost,omitempty"`

	// ASOManagedClustersAgentPoolPatches defines JSON merge patches to be applied to the generated ASO ManagedClustersAgentPool resource.
	// WARNING: This is meant to be used sparingly to enable features for development and testing that are not
	// otherwise represented in the CAPZ API. Misconfiguration that conflicts with CAPZ's normal mode of
	// operation is possible.
	// +optional
	ASOManagedClustersAgentPoolPatches []string `json:"asoManagedClustersAgentPoolPatches,omitempty"`
}

// ManagedControlPlaneVirtualNetworkClassSpec defines the ManagedControlPlaneVirtualNetwork properties that may be shared across several managed control plane vnets.
type ManagedControlPlaneVirtualNetworkClassSpec struct {
	// +kubebuilder:default:="10.0.0.0/8"
	CIDRBlock string `json:"cidrBlock"`
	// +optional
	Subnet ManagedControlPlaneSubnet `json:"subnet,omitempty"`
}

// APIServerAccessProfileClassSpec defines the APIServerAccessProfile properties that may be shared across several API server access profiles.
type APIServerAccessProfileClassSpec struct {
	// EnablePrivateCluster indicates whether to create the cluster as a private cluster or not.
	// +optional
	EnablePrivateCluster *bool `json:"enablePrivateCluster,omitempty"`

	// PrivateDNSZone enables private dns zone mode for private cluster.
	// +optional
	PrivateDNSZone *string `json:"privateDNSZone,omitempty"`

	// EnablePrivateClusterPublicFQDN indicates whether to create additional public FQDN for private cluster or not.
	// +optional
	EnablePrivateClusterPublicFQDN *bool `json:"enablePrivateClusterPublicFQDN,omitempty"`
}

// ExtendedLocationSpec defines the ExtendedLocation properties to enable CAPZ for Azure public MEC.
type ExtendedLocationSpec struct {
	// Name defines the name for the extended location.
	Name string `json:"name"`

	// Type defines the type for the extended location.
	// +kubebuilder:validation:Enum=EdgeZone
	Type string `json:"type"`
}

// NetworkClassSpec defines the NetworkSpec properties that may be shared across several Azure clusters.
type NetworkClassSpec struct {
	// PrivateDNSZoneName defines the zone name for the Azure Private DNS.
	// +optional
	PrivateDNSZoneName string `json:"privateDNSZoneName,omitempty"`

	// PrivateDNSZoneResourceGroup defines the resource group to be used for Azure Private DNS Zone.
	// If not specified, the resource group of the cluster will be used to create the Azure Private DNS Zone.
	// +optional
	PrivateDNSZoneResourceGroup string `json:"privateDNSZoneResourceGroup,omitempty"`
}

// VnetClassSpec defines the VnetSpec properties that may be shared across several Azure clusters.
type VnetClassSpec struct {
	// CIDRBlocks defines the virtual network's address space, specified as one or more address prefixes in CIDR notation.
	// +optional
	CIDRBlocks []string `json:"cidrBlocks,omitempty"`

	// Tags is a collection of tags describing the resource.
	// +optional
	Tags Tags `json:"tags,omitempty"`
}

// SubnetClassSpec defines the SubnetSpec properties that may be shared across several Azure clusters.
type SubnetClassSpec struct {
	// Name defines a name for the subnet resource.
	Name string `json:"name"`

	// Role defines the subnet role (eg. Node, ControlPlane)
	// +kubebuilder:validation:Enum=node;control-plane;bastion;cluster
	Role SubnetRole `json:"role"`

	// CIDRBlocks defines the subnet's address space, specified as one or more address prefixes in CIDR notation.
	// +optional
	CIDRBlocks []string `json:"cidrBlocks,omitempty"`

	// ServiceEndpoints is a slice of Virtual Network service endpoints to enable for the subnets.
	// +optional
	ServiceEndpoints ServiceEndpoints `json:"serviceEndpoints,omitempty"`

	// PrivateEndpoints defines a list of private endpoints that should be attached to this subnet.
	// +optional
	PrivateEndpoints PrivateEndpoints `json:"privateEndpoints,omitempty"`
}

// LoadBalancerClassSpec defines the LoadBalancerSpec properties that may be shared across several Azure clusters.
type LoadBalancerClassSpec struct {
	// +optional
	SKU SKU `json:"sku,omitempty"`
	// +optional
	Type LBType `json:"type,omitempty"`
	// IdleTimeoutInMinutes specifies the timeout for the TCP idle connection.
	// +optional
	IdleTimeoutInMinutes *int32 `json:"idleTimeoutInMinutes,omitempty"`
}

// FleetsMemberClassSpec defines the FleetsMemberSpec properties that may be shared across several Azure clusters.
type FleetsMemberClassSpec struct {
	// Group is the group this member belongs to for multi-cluster update management.
	// +kubebuilder:default:=default
	// +optional
	Group string `json:"group,omitempty"`

	// ManagerName is the name of the fleet manager.
	ManagerName string `json:"managerName"`

	// ManagerResourceGroup is the resource group of the fleet manager.
	ManagerResourceGroup string `json:"managerResourceGroup"`
}

// SecurityGroupClass defines the SecurityGroup properties that may be shared across several Azure clusters.
type SecurityGroupClass struct {
	// +optional
	SecurityRules SecurityRules `json:"securityRules,omitempty"`
	// +optional
	Tags Tags `json:"tags,omitempty"`
}

// FrontendIPClass defines the FrontendIP properties that may be shared across several Azure clusters.
type FrontendIPClass struct {
	// +optional
	PrivateIPAddress string `json:"privateIP,omitempty"`
}

// setDefaults sets default values for AzureClusterClassSpec.
func (acc *AzureClusterClassSpec) setDefaults() {
	if acc.AzureEnvironment == "" {
		acc.AzureEnvironment = DefaultAzureCloud
	}
}

// setDefaults sets default values for VnetClassSpec.
func (vc *VnetClassSpec) setDefaults() {
	if len(vc.CIDRBlocks) == 0 {
		vc.CIDRBlocks = []string{DefaultVnetCIDR}
	}
}

// setDefaults sets default values for SubnetClassSpec.
func (sc *SubnetClassSpec) setDefaults(cidr string) {
	if len(sc.CIDRBlocks) == 0 {
		sc.CIDRBlocks = []string{cidr}
	}
}

// setDefaults sets default values for SecurityGroupClass.
func (sgc *SecurityGroupClass) setDefaults() {
	for i := range sgc.SecurityRules {
		if sgc.SecurityRules[i].Direction == "" {
			sgc.SecurityRules[i].Direction = SecurityRuleDirectionInbound
		}
	}
}
