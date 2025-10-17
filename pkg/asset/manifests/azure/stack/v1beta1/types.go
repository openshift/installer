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
	"github.com/pkg/errors"
	"k8s.io/apimachinery/pkg/api/resource"
	"k8s.io/utils/net"
)

const (
	// ControlPlane machine label.
	ControlPlane string = "control-plane"
	// Node machine label.
	Node string = "node"
	// Bastion subnet label.
	Bastion string = "bastion"
	// Cluster subnet label.
	Cluster string = "cluster"
)

// SecurityEncryptionType represents the Encryption Type when the virtual machine is a
// Confidential VM.
type SecurityEncryptionType string

const (
	// SecurityEncryptionTypeVMGuestStateOnly disables OS disk confidential encryption.
	SecurityEncryptionTypeVMGuestStateOnly SecurityEncryptionType = "VMGuestStateOnly"
	// SecurityEncryptionTypeDiskWithVMGuestState OS disk confidential encryption with a
	// platform-managed key (PMK) or a customer-managed key (CMK).
	SecurityEncryptionTypeDiskWithVMGuestState SecurityEncryptionType = "DiskWithVMGuestState"
)

// SecurityTypes represents the SecurityType of the virtual machine.
type SecurityTypes string

const (
	// SecurityTypesConfidentialVM defines the SecurityType of the virtual machine as a Confidential VM.
	SecurityTypesConfidentialVM SecurityTypes = "ConfidentialVM"
	// SecurityTypesTrustedLaunch defines the SecurityType of the virtual machine as a Trusted Launch VM.
	SecurityTypesTrustedLaunch SecurityTypes = "TrustedLaunch"
)

// Futures is a slice of Future.
type Futures []Future

const (
	// PatchFuture is a future that was derived from a PATCH request.
	PatchFuture string = "PATCH"
	// PutFuture is a future that was derived from a PUT request.
	PutFuture string = "PUT"
	// DeleteFuture is a future that was derived from a DELETE request.
	DeleteFuture string = "DELETE"
)

// Future contains the data needed for an Azure long-running operation to continue across reconcile loops.
type Future struct {
	// Type describes the type of future, such as update, create, delete, etc.
	Type string `json:"type"`

	// ResourceGroup is the Azure resource group for the resource.
	// +optional
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// ServiceName is the name of the Azure service.
	// Together with the name of the resource, this forms the unique identifier for the future.
	ServiceName string `json:"serviceName"`

	// Name is the name of the Azure resource.
	// Together with the service name, this forms the unique identifier for the future.
	Name string `json:"name"`

	// Data is the base64 url encoded json Azure AutoRest Future.
	Data string `json:"data"`
}

// NetworkSpec specifies what the Azure networking resources should look like.
type NetworkSpec struct {
	// Vnet is the configuration for the Azure virtual network.
	// +optional
	Vnet VnetSpec `json:"vnet,omitempty"`

	// Subnets is the configuration for the control-plane subnet and the node subnet.
	// +optional
	Subnets Subnets `json:"subnets,omitempty"`

	// APIServerLB is the configuration for the control-plane load balancer.
	// +optional
	APIServerLB *LoadBalancerSpec `json:"apiServerLB,omitempty"`

	// NodeOutboundLB is the configuration for the node outbound load balancer.
	// +optional
	NodeOutboundLB *LoadBalancerSpec `json:"nodeOutboundLB,omitempty"`

	// ControlPlaneOutboundLB is the configuration for the control-plane outbound load balancer.
	// This is different from APIServerLB, and is used only in private clusters (optionally) for enabling outbound traffic.
	// +optional
	ControlPlaneOutboundLB *LoadBalancerSpec `json:"controlPlaneOutboundLB,omitempty"`

	NetworkClassSpec `json:",inline"`
}

// VnetSpec configures an Azure virtual network.
type VnetSpec struct {
	// ResourceGroup is the name of the resource group of the existing virtual network
	// or the resource group where a managed virtual network should be created.
	// +optional
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// ID is the Azure resource ID of the virtual network.
	// READ-ONLY
	// +optional
	ID string `json:"id,omitempty"`

	// Name defines a name for the virtual network resource.
	Name string `json:"name"`

	// Peerings defines a list of peerings of the newly created virtual network with existing virtual networks.
	// +optional
	Peerings VnetPeerings `json:"peerings,omitempty"`

	VnetClassSpec `json:",inline"`
}

// VnetPeeringSpec specifies an existing remote virtual network to peer with the AzureCluster's virtual network.
type VnetPeeringSpec struct {
	VnetPeeringClassSpec `json:",inline"`
}

// VnetPeeringClassSpec specifies a virtual network peering class.
type VnetPeeringClassSpec struct {
	// ResourceGroup is the resource group name of the remote virtual network.
	// +optional
	ResourceGroup string `json:"resourceGroup,omitempty"`

	// RemoteVnetName defines name of the remote virtual network.
	RemoteVnetName string `json:"remoteVnetName"`

	// ForwardPeeringProperties specifies VnetPeeringProperties for peering from the cluster's virtual network to the
	// remote virtual network.
	// +optional
	ForwardPeeringProperties VnetPeeringProperties `json:"forwardPeeringProperties,omitempty"`

	// ReversePeeringProperties specifies VnetPeeringProperties for peering from the remote virtual network to the
	// cluster's virtual network.
	// +optional
	ReversePeeringProperties VnetPeeringProperties `json:"reversePeeringProperties,omitempty"`
}

// VnetPeeringProperties specifies virtual network peering properties.
type VnetPeeringProperties struct {
	// AllowForwardedTraffic specifies whether the forwarded traffic from the VMs in the local virtual network will be
	// allowed/disallowed in remote virtual network.
	// +optional
	AllowForwardedTraffic *bool `json:"allowForwardedTraffic,omitempty"`

	// AllowGatewayTransit specifies if gateway links can be used in remote virtual networking to link to this virtual
	// network.
	// +optional
	AllowGatewayTransit *bool `json:"allowGatewayTransit,omitempty"`

	// AllowVirtualNetworkAccess specifies whether the VMs in the local virtual network space would be able to access
	// the VMs in remote virtual network space.
	// +optional
	AllowVirtualNetworkAccess *bool `json:"allowVirtualNetworkAccess,omitempty"`

	// UseRemoteGateways specifies if remote gateways can be used on this virtual network.
	// If the flag is set to true, and allowGatewayTransit on remote peering is also set to true, the virtual network
	// will use the gateways of the remote virtual network for transit. Only one peering can have this flag set to true.
	// This flag cannot be set if virtual network already has a gateway.
	// +optional
	UseRemoteGateways *bool `json:"useRemoteGateways,omitempty"`
}

// VnetPeerings is a slice of VnetPeering.
type VnetPeerings []VnetPeeringSpec

// IsManaged returns true if the vnet is managed.
func (v *VnetSpec) IsManaged(clusterName string) bool {
	return v.ID == "" || v.Tags.HasOwned(clusterName)
}

// Subnets is a slice of Subnet.
// +listType=map
// +listMapKey=name
type Subnets []SubnetSpec

// ServiceEndpoints is a slice of string.
// +listType=map
// +listMapKey=service
type ServiceEndpoints []ServiceEndpointSpec

// PrivateEndpoints is a slice of PrivateEndpointSpec.
// +listType=map
// +listMapKey=name
type PrivateEndpoints []PrivateEndpointSpec

// SecurityGroup defines an Azure security group.
type SecurityGroup struct {
	// ID is the Azure resource ID of the security group.
	// READ-ONLY
	// +optional
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`

	SecurityGroupClass `json:",inline"`
}

// RouteTable defines an Azure route table.
type RouteTable struct {
	// ID is the Azure resource ID of the route table.
	// READ-ONLY
	// +optional
	ID   string `json:"id,omitempty"`
	Name string `json:"name"`
}

// NatGateway defines an Azure NAT gateway.
// NAT gateway resources are part of Vnet NAT and provide outbound Internet connectivity for subnets of a virtual network.
type NatGateway struct {
	// ID is the Azure resource ID of the NAT gateway.
	// READ-ONLY
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	NatGatewayIP PublicIPSpec `json:"ip,omitempty"`

	NatGatewayClassSpec `json:",inline"`
}

// NatGatewayClassSpec defines a NAT gateway class specification.
type NatGatewayClassSpec struct {
	Name string `json:"name"`
}

// SecurityGroupProtocol defines the protocol type for a security group rule.
type SecurityGroupProtocol string

const (
	// SecurityGroupProtocolAll is a wildcard for all IP protocols.
	SecurityGroupProtocolAll = SecurityGroupProtocol("*")
	// SecurityGroupProtocolTCP represents the TCP protocol.
	SecurityGroupProtocolTCP = SecurityGroupProtocol("Tcp")
	// SecurityGroupProtocolUDP represents the UDP protocol.
	SecurityGroupProtocolUDP = SecurityGroupProtocol("Udp")
	// SecurityGroupProtocolICMP represents the ICMP protocol.
	SecurityGroupProtocolICMP = SecurityGroupProtocol("Icmp")
)

// SecurityRuleDirection defines the direction type for a security group rule.
type SecurityRuleDirection string

const (
	// SecurityRuleDirectionInbound defines an ingress security rule.
	SecurityRuleDirectionInbound = SecurityRuleDirection("Inbound")

	// SecurityRuleDirectionOutbound defines an egress security rule.
	SecurityRuleDirectionOutbound = SecurityRuleDirection("Outbound")
)

// SecurityRuleAccess defines the action type for a security group rule.
type SecurityRuleAccess string

const (
	// SecurityRuleActionAllow allows traffic defined in the rule.
	SecurityRuleActionAllow SecurityRuleAccess = "Allow"

	// SecurityRuleActionDeny denies traffic defined in the rule.
	SecurityRuleActionDeny SecurityRuleAccess = "Deny"
)

// SecurityRule defines an Azure security rule for security groups.
type SecurityRule struct {
	// Name is a unique name within the network security group.
	Name string `json:"name"`
	// A description for this rule. Restricted to 140 chars.
	Description string `json:"description"`
	// Protocol specifies the protocol type. "Tcp", "Udp", "Icmp", or "*".
	// +kubebuilder:validation:Enum=Tcp;Udp;Icmp;*
	Protocol SecurityGroupProtocol `json:"protocol"`
	// Direction indicates whether the rule applies to inbound, or outbound traffic. "Inbound" or "Outbound".
	// +kubebuilder:validation:Enum=Inbound;Outbound
	Direction SecurityRuleDirection `json:"direction"`
	// Priority is a number between 100 and 4096. Each rule should have a unique value for priority. Rules are processed in priority order, with lower numbers processed before higher numbers. Once traffic matches a rule, processing stops.
	// +optional
	Priority int32 `json:"priority,omitempty"`
	// SourcePorts specifies source port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	// +optional
	SourcePorts *string `json:"sourcePorts,omitempty"`
	// DestinationPorts specifies the destination port or range. Integer or range between 0 and 65535. Asterix '*' can also be used to match all ports.
	// +optional
	DestinationPorts *string `json:"destinationPorts,omitempty"`
	// Source specifies the CIDR or source IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used. If this is an ingress rule, specifies where network traffic originates from.
	// +optional
	Source *string `json:"source,omitempty"`
	// Sources specifies The CIDR or source IP ranges.
	Sources []*string `json:"sources,omitempty"`
	// Destination is the destination address prefix. CIDR or destination IP range. Asterix '*' can also be used to match all source IPs. Default tags such as 'VirtualNetwork', 'AzureLoadBalancer' and 'Internet' can also be used.
	// +optional
	Destination *string `json:"destination,omitempty"`
	// Action specifies whether network traffic is allowed or denied. Can either be "Allow" or "Deny". Defaults to "Allow".
	// +kubebuilder:default=Allow
	// +kubebuilder:validation:Enum=Allow;Deny
	//+optional
	Action SecurityRuleAccess `json:"action"`
}

// SecurityRules is a slice of Azure security rules for security groups.
// +listType=map
// +listMapKey=name
type SecurityRules []SecurityRule

// LoadBalancerSpec defines an Azure load balancer.
type LoadBalancerSpec struct {
	// ID is the Azure resource ID of the load balancer.
	// READ-ONLY
	// +optional
	ID string `json:"id,omitempty"`
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	FrontendIPs []FrontendIP `json:"frontendIPs,omitempty"`
	// FrontendIPsCount specifies the number of frontend IP addresses for the load balancer.
	// +optional
	FrontendIPsCount *int32 `json:"frontendIPsCount,omitempty"`
	// BackendPool describes the backend pool of the load balancer.
	// +optional
	BackendPool BackendPool `json:"backendPool,omitempty"`

	LoadBalancerClassSpec `json:",inline"`
}

// SKU defines an Azure load balancer SKU.
type SKU string

const (
	// SKUStandard is the value for the Azure load balancer Standard SKU.
	SKUStandard = SKU("Standard")
)

// LBType defines an Azure load balancer Type.
type LBType string

const (
	// Internal is the value for the Azure load balancer internal type.
	Internal = LBType("Internal")
	// Public is the value for the Azure load balancer public type.
	Public = LBType("Public")
)

// FrontendIP defines a load balancer frontend IP configuration.
type FrontendIP struct {
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// +optional
	PublicIP *PublicIPSpec `json:"publicIP,omitempty"`

	FrontendIPClass `json:",inline"`
}

// PublicIPSpec defines the inputs to create an Azure public IP address.
type PublicIPSpec struct {
	Name string `json:"name"`
	// +optional
	DNSName string `json:"dnsName,omitempty"`
	// +optional
	IPTags []IPTag `json:"ipTags,omitempty"`
}

// IPTag contains the IpTag associated with the object.
type IPTag struct {
	// Type specifies the IP tag type. Example: FirstPartyUsage.
	Type string `json:"type"`
	// Tag specifies the value of the IP tag associated with the public IP. Example: SQL.
	Tag string `json:"tag"`
}

// VMState describes the state of an Azure virtual machine.
// Deprecated: use ProvisioningState.
type VMState string

// ProvisioningState describes the provisioning state of an Azure resource.
type ProvisioningState string

const (
	// Creating ...
	Creating ProvisioningState = "Creating"
	// Deleting ...
	Deleting ProvisioningState = "Deleting"
	// Failed ...
	Failed ProvisioningState = "Failed"
	// Migrating ...
	Migrating ProvisioningState = "Migrating"
	// Succeeded ...
	Succeeded ProvisioningState = "Succeeded"
	// Updating ...
	Updating ProvisioningState = "Updating"
	// Canceled represents an action which was initiated but terminated by the user before completion.
	Canceled ProvisioningState = "Canceled"
	// Deleted represents a deleted VM
	// NOTE: This state is specific to capz, and does not have corresponding mapping in Azure API (https://learn.microsoft.com/azure/virtual-machines/states-billing#provisioning-states)
	Deleted ProvisioningState = "Deleted"
)

// Image defines information about the image to use for VM creation.
// There are three ways to specify an image: by ID, Marketplace Image or SharedImageGallery
// One of ID, SharedImage or Marketplace should be set.
type Image struct {
	// ID specifies an image to use by ID
	// +optional
	ID *string `json:"id,omitempty"`

	// SharedGallery specifies an image to use from an Azure Shared Image Gallery
	// Deprecated: use ComputeGallery instead.
	// +optional
	SharedGallery *AzureSharedGalleryImage `json:"sharedGallery,omitempty"`

	// Marketplace specifies an image to use from the Azure Marketplace
	// +optional
	Marketplace *AzureMarketplaceImage `json:"marketplace,omitempty"`

	// ComputeGallery specifies an image to use from the Azure Compute Gallery
	// +optional
	ComputeGallery *AzureComputeGalleryImage `json:"computeGallery,omitempty"`
}

// AzureComputeGalleryImage defines an image in the Azure Compute Gallery to use for VM creation.
type AzureComputeGalleryImage struct {
	// Gallery specifies the name of the compute image gallery that contains the image
	// +kubebuilder:validation:MinLength=1
	Gallery string `json:"gallery"`
	// Name is the name of the image
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Version specifies the version of the marketplace image. The allowed formats
	// are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers.
	// Specify 'latest' to use the latest version of an image available at deploy time.
	// Even if you use 'latest', the VM image will not automatically update after deploy
	// time even if a new version becomes available.
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
	// SubscriptionID is the identifier of the subscription that contains the private compute gallery.
	// +optional
	SubscriptionID *string `json:"subscriptionID,omitempty"`
	// ResourceGroup specifies the resource group containing the private compute gallery.
	// +optional
	ResourceGroup *string `json:"resourceGroup,omitempty"`
	// Plan contains plan information.
	// +optional
	Plan *ImagePlan `json:"plan,omitempty"`
}

// ImagePlan contains plan information for marketplace images.
type ImagePlan struct {
	// Publisher is the name of the organization that created the image
	// +kubebuilder:validation:MinLength=1
	Publisher string `json:"publisher"`
	// Offer specifies the name of a group of related images created by the publisher.
	// For example, UbuntuServer, WindowsServer
	// +kubebuilder:validation:MinLength=1
	Offer string `json:"offer"`
	// SKU specifies an instance of an offer, such as a major release of a distribution.
	// For example, 18.04-LTS, 2019-Datacenter
	// +kubebuilder:validation:MinLength=1
	SKU string `json:"sku"`
}

// AzureMarketplaceImage defines an image in the Azure Marketplace to use for VM creation.
type AzureMarketplaceImage struct {
	ImagePlan `json:",inline"`

	// Version specifies the version of an image sku. The allowed formats
	// are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers.
	// Specify 'latest' to use the latest version of an image available at deploy time.
	// Even if you use 'latest', the VM image will not automatically update after deploy
	// time even if a new version becomes available.
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
	// ThirdPartyImage indicates the image is published by a third party publisher and a Plan
	// will be generated for it.
	// +kubebuilder:default=false
	// +optional
	ThirdPartyImage bool `json:"thirdPartyImage"`
}

// AzureSharedGalleryImage defines an image in a Shared Image Gallery to use for VM creation.
type AzureSharedGalleryImage struct {
	// SubscriptionID is the identifier of the subscription that contains the shared image gallery
	// +kubebuilder:validation:MinLength=1
	SubscriptionID string `json:"subscriptionID"`
	// ResourceGroup specifies the resource group containing the shared image gallery
	// +kubebuilder:validation:MinLength=1
	ResourceGroup string `json:"resourceGroup"`
	// Gallery specifies the name of the shared image gallery that contains the image
	// +kubebuilder:validation:MinLength=1
	Gallery string `json:"gallery"`
	// Name is the name of the image
	// +kubebuilder:validation:MinLength=1
	Name string `json:"name"`
	// Version specifies the version of the marketplace image. The allowed formats
	// are Major.Minor.Build or 'latest'. Major, Minor, and Build are decimal numbers.
	// Specify 'latest' to use the latest version of an image available at deploy time.
	// Even if you use 'latest', the VM image will not automatically update after deploy
	// time even if a new version becomes available.
	// +kubebuilder:validation:MinLength=1
	Version string `json:"version"`
	// Publisher is the name of the organization that created the image.
	// This value will be used to add a `Plan` in the API request when creating the VM/VMSS resource.
	// This is needed when the source image from which this SIG image was built requires the `Plan` to be used.
	// +optional
	Publisher *string `json:"publisher,omitempty"`
	// Offer specifies the name of a group of related images created by the publisher.
	// For example, UbuntuServer, WindowsServer
	// This value will be used to add a `Plan` in the API request when creating the VM/VMSS resource.
	// This is needed when the source image from which this SIG image was built requires the `Plan` to be used.
	// +optional
	Offer *string `json:"offer,omitempty"`
	// SKU specifies an instance of an offer, such as a major release of a distribution.
	// For example, 18.04-LTS, 2019-Datacenter
	// This value will be used to add a `Plan` in the API request when creating the VM/VMSS resource.
	// This is needed when the source image from which this SIG image was built requires the `Plan` to be used.
	// +optional
	SKU *string `json:"sku,omitempty"`
}

// VMIdentity defines the identity of the virtual machine, if configured.
// +kubebuilder:validation:Enum=None;SystemAssigned;UserAssigned
type VMIdentity string

const (
	// VMIdentityNone ...
	VMIdentityNone VMIdentity = "None"
	// VMIdentitySystemAssigned ...
	VMIdentitySystemAssigned VMIdentity = "SystemAssigned"
	// VMIdentityUserAssigned ...
	VMIdentityUserAssigned VMIdentity = "UserAssigned"
)

// SpotEvictionPolicy defines the eviction policy for spot VMs, if configured.
// +kubebuilder:validation:Enum=Deallocate;Delete
type SpotEvictionPolicy string

const (
	// SpotEvictionPolicyDeallocate is the default eviction policy and will deallocate the VM when the node is marked for eviction.
	SpotEvictionPolicyDeallocate SpotEvictionPolicy = "Deallocate"
	// SpotEvictionPolicyDelete will delete the VM when the node is marked for eviction.
	SpotEvictionPolicyDelete SpotEvictionPolicy = "Delete"
)

// UserAssignedIdentity defines the user-assigned identities provided
// by the user to be assigned to Azure resources.
type UserAssignedIdentity struct {
	// ProviderID is the identification ID of the user-assigned Identity, the format of an identity is:
	// 'azure:///subscriptions/{subscriptionId}/resourceGroups/{resourceGroupName}/providers/Microsoft.ManagedIdentity/userAssignedIdentities/{identityName}'
	ProviderID string `json:"providerID"`
}

// IdentityType represents different types of identities.
// +kubebuilder:validation:Enum=ServicePrincipal;UserAssignedMSI;ManualServicePrincipal;ServicePrincipalCertificate;WorkloadIdentity;UserAssignedIdentityCredential
type IdentityType string

const (
	// UserAssignedMSI represents a user-assigned managed identity.
	UserAssignedMSI IdentityType = "UserAssignedMSI"

	// ServicePrincipal represents a service principal using a client password as secret.
	ServicePrincipal IdentityType = "ServicePrincipal"

	// ManualServicePrincipal represents a manual service principal.
	ManualServicePrincipal IdentityType = "ManualServicePrincipal"

	// ServicePrincipalCertificate represents a service principal using a certificate as secret.
	ServicePrincipalCertificate IdentityType = "ServicePrincipalCertificate"

	// WorkloadIdentity represents a WorkloadIdentity.
	WorkloadIdentity IdentityType = "WorkloadIdentity"

	// UserAssignedIdentityCredential represents a UserAssignedIdentityCredential.
	UserAssignedIdentityCredential IdentityType = "UserAssignedIdentityCredential"
)

// OSDisk defines the operating system disk for a VM.
//
// WARNING: this requires any updates to ManagedDisk to be manually converted. This is due to the odd issue with
// conversion-gen where the warning message generated uses a relative directory import rather than the fully
// qualified import when generating outside of the GOPATH.
type OSDisk struct {
	// +kubebuilder:default:=Linux
	OSType string `json:"osType"`
	// DiskSizeGB is the size in GB to assign to the OS disk.
	// Will have a default of 30GB if not provided
	// +optional
	DiskSizeGB *int32 `json:"diskSizeGB,omitempty"`
	// ManagedDisk specifies the Managed Disk parameters for the OS disk.
	// +optional
	ManagedDisk *ManagedDiskParameters `json:"managedDisk,omitempty"`
	// +optional
	DiffDiskSettings *DiffDiskSettings `json:"diffDiskSettings,omitempty"`
	// CachingType specifies the caching requirements.
	// +optional
	// +kubebuilder:validation:Enum=None;ReadOnly;ReadWrite
	// +kubebuilder:default:=None
	CachingType string `json:"cachingType,omitempty"`
}

// DataDisk specifies the parameters that are used to add one or more data disks to the machine.
type DataDisk struct {
	// NameSuffix is the suffix to be appended to the machine name to generate the disk name.
	// Each disk name will be in format <machineName>_<nameSuffix>.
	NameSuffix string `json:"nameSuffix"`
	// DiskSizeGB is the size in GB to assign to the data disk.
	DiskSizeGB int32 `json:"diskSizeGB"`
	// ManagedDisk specifies the Managed Disk parameters for the data disk.
	// +optional
	ManagedDisk *ManagedDiskParameters `json:"managedDisk,omitempty"`
	// Lun Specifies the logical unit number of the data disk. This value is used to identify data disks within the VM and therefore must be unique for each data disk attached to a VM.
	// The value must be between 0 and 63.
	// +optional
	Lun *int32 `json:"lun,omitempty"`
	// CachingType specifies the caching requirements.
	// +optional
	// +kubebuilder:validation:Enum=None;ReadOnly;ReadWrite
	CachingType string `json:"cachingType,omitempty"`
}

// VMExtension specifies the parameters for a custom VM extension.
type VMExtension struct {
	// Name is the name of the extension.
	Name string `json:"name"`
	// Publisher is the name of the extension handler publisher.
	Publisher string `json:"publisher"`
	// Version specifies the version of the script handler.
	Version string `json:"version"`
	// Settings is a JSON formatted public settings for the extension.
	// +optional
	Settings Tags `json:"settings,omitempty"`
	// ProtectedSettings is a JSON formatted protected settings for the extension.
	// +optional
	ProtectedSettings Tags `json:"protectedSettings,omitempty"`
}

// ManagedDiskParameters defines the parameters of a managed disk.
type ManagedDiskParameters struct {
	// +optional
	StorageAccountType string `json:"storageAccountType,omitempty"`
	// DiskEncryptionSet specifies the customer-managed disk encryption set resource id for the managed disk.
	// +optional
	DiskEncryptionSet *DiskEncryptionSetParameters `json:"diskEncryptionSet,omitempty"`
	// SecurityProfile specifies the security profile for the managed disk.
	// +optional
	SecurityProfile *VMDiskSecurityProfile `json:"securityProfile,omitempty"`
}

// VMDiskSecurityProfile specifies the security profile settings for the managed disk.
// It can be set only for Confidential VMs.
type VMDiskSecurityProfile struct {
	// DiskEncryptionSet specifies the customer-managed disk encryption set resource id for the
	// managed disk that is used for Customer Managed Key encrypted ConfidentialVM OS Disk and
	// VMGuest blob.
	// +optional
	DiskEncryptionSet *DiskEncryptionSetParameters `json:"diskEncryptionSet,omitempty"`
	// SecurityEncryptionType specifies the encryption type of the managed disk.
	// It is set to DiskWithVMGuestState to encrypt the managed disk along with the VMGuestState
	// blob, and to VMGuestStateOnly to encrypt the VMGuestState blob only.
	// When set to VMGuestStateOnly, VirtualizedTrustedPlatformModule should be set to Enabled.
	// When set to DiskWithVMGuestState, EncryptionAtHost should be disabled, SecureBoot and
	// VirtualizedTrustedPlatformModule should be set to Enabled.
	// It can be set only for Confidential VMs.
	// +kubebuilder:validation:Enum=VMGuestStateOnly;DiskWithVMGuestState
	// +optional
	SecurityEncryptionType SecurityEncryptionType `json:"securityEncryptionType,omitempty"`
}

// DiskEncryptionSetParameters defines disk encryption options.
type DiskEncryptionSetParameters struct {
	// ID defines resourceID for diskEncryptionSet resource. It must be in the same subscription
	// +optional
	ID string `json:"id,omitempty"`
}

// DiffDiskPlacement - Specifies the ephemeral disk placement for operating system disk. This property can be used by user
// in the request to choose the location i.e, cache disk, resource disk or nvme disk space for
// Ephemeral OS disk provisioning. For more information on Ephemeral OS disk size requirements, please refer Ephemeral OS
// disk size requirements for Windows VM at
// https://docs.microsoft.com/azure/virtual-machines/windows/ephemeral-os-disks#size-requirements and Linux VM at
// https://docs.microsoft.com/azure/virtual-machines/linux/ephemeral-os-disks#size-requirements.
type DiffDiskPlacement string

const (
	// DiffDiskPlacementCacheDisk places the OsDisk on cache disk.
	DiffDiskPlacementCacheDisk DiffDiskPlacement = "CacheDisk"

	// DiffDiskPlacementNvmeDisk places the OsDisk on NVMe disk.
	DiffDiskPlacementNvmeDisk DiffDiskPlacement = "NvmeDisk"

	// DiffDiskPlacementResourceDisk places the OsDisk on temp disk.
	DiffDiskPlacementResourceDisk DiffDiskPlacement = "ResourceDisk"
)

// PossibleDiffDiskPlacementValues returns the possible values for the DiffDiskPlacement const type.
func PossibleDiffDiskPlacementValues() []DiffDiskPlacement {
	return []DiffDiskPlacement{
		DiffDiskPlacementCacheDisk,
		DiffDiskPlacementNvmeDisk,
		DiffDiskPlacementResourceDisk,
	}
}

// DiffDiskSettings describe ephemeral disk settings for the os disk.
type DiffDiskSettings struct {
	// Option enables ephemeral OS when set to "Local"
	// See https://learn.microsoft.com/azure/virtual-machines/ephemeral-os-disks for full details
	// +kubebuilder:validation:Enum=Local
	Option string `json:"option"`

	// Placement specifies the ephemeral disk placement for operating system disk. If placement is specified, Option must be set to "Local".
	// +kubebuilder:validation:Enum=CacheDisk;NvmeDisk;ResourceDisk
	// +optional
	Placement *DiffDiskPlacement `json:"placement,omitempty"`
}

// SubnetRole defines the unique role of a subnet.
type SubnetRole string

const (
	// SubnetNode defines a Kubernetes workload node role.
	SubnetNode = SubnetRole(Node)

	// SubnetControlPlane defines a Kubernetes control plane node role.
	SubnetControlPlane = SubnetRole(ControlPlane)

	// SubnetBastion defines a Bastion subnet role.
	SubnetBastion = SubnetRole(Bastion)

	// SubnetCluster defines a role that can be used for both Kubernetes control plane node and Kubernetes workload node.
	SubnetCluster = SubnetRole(Cluster)
)

// SubnetSpec configures an Azure subnet.
type SubnetSpec struct {
	// ID is the Azure resource ID of the subnet.
	// READ-ONLY
	// +optional
	ID string `json:"id,omitempty"`

	// SecurityGroup defines the NSG (network security group) that should be attached to this subnet.
	// +optional
	SecurityGroup SecurityGroup `json:"securityGroup,omitempty"`

	// RouteTable defines the route table that should be attached to this subnet.
	// +optional
	RouteTable RouteTable `json:"routeTable,omitempty"`

	// NatGateway associated with this subnet.
	// +optional
	NatGateway NatGateway `json:"natGateway,omitempty"`

	SubnetClassSpec `json:",inline"`
}

// ServiceEndpointSpec configures an Azure Service Endpoint.
type ServiceEndpointSpec struct {
	Service string `json:"service"`

	Locations []string `json:"locations"`
}

// PrivateLinkServiceConnection defines the specification for a private link service connection associated with a private endpoint.
type PrivateLinkServiceConnection struct {
	// Name specifies the name of the private link service.
	// +optional
	Name string `json:"name,omitempty"`
	// PrivateLinkServiceID specifies the resource ID of the private link service.
	PrivateLinkServiceID string `json:"privateLinkServiceID,omitempty"`
	// GroupIDs specifies the ID(s) of the group(s) obtained from the remote resource that this private endpoint should connect to.
	// +optional
	GroupIDs []string `json:"groupIDs,omitempty"`
	// RequestMessage specifies a message passed to the owner of the remote resource with the private endpoint connection request.
	// +kubebuilder:validation:MaxLength=140
	// +optional
	RequestMessage string `json:"requestMessage,omitempty"`
}

// PrivateEndpointSpec configures an Azure Private Endpoint.
type PrivateEndpointSpec struct {
	// Name specifies the name of the private endpoint.
	Name string `json:"name"`
	// Location specifies the region to create the private endpoint.
	// +optional
	Location string `json:"location,omitempty"`
	// PrivateLinkServiceConnections specifies Private Link Service Connections of the private endpoint.
	PrivateLinkServiceConnections []PrivateLinkServiceConnection `json:"privateLinkServiceConnections,omitempty"`
	// CustomNetworkInterfaceName specifies the network interface name associated with the private endpoint.
	// +optional
	CustomNetworkInterfaceName string `json:"customNetworkInterfaceName,omitempty"`
	// PrivateIPAddresses specifies the IP addresses for the network interface associated with the private endpoint.
	// They have to be part of the subnet where the private endpoint is linked.
	// +optional
	PrivateIPAddresses []string `json:"privateIPAddresses,omitempty"`
	// ApplicationSecurityGroups specifies the Application security group in which the private endpoint IP configuration is included.
	// +optional
	ApplicationSecurityGroups []string `json:"applicationSecurityGroups,omitempty"`
	// ManualApproval specifies if the connection approval needs to be done manually or not.
	// Set it true when the network admin does not have access to approve connections to the remote resource.
	// Defaults to false.
	// +optional
	ManualApproval bool `json:"manualApproval,omitempty"`
}

// NetworkInterface defines a network interface.
type NetworkInterface struct {
	// SubnetName specifies the subnet in which the new network interface will be placed.
	SubnetName string `json:"subnetName,omitempty"`

	// PrivateIPConfigs specifies the number of private IP addresses to attach to the interface.
	// Defaults to 1 if not specified.
	// +optional
	PrivateIPConfigs int `json:"privateIPConfigs,omitempty"`

	// AcceleratedNetworking enables or disables Azure accelerated networking. If omitted, it will be set based on
	// whether the requested VMSize supports accelerated networking.
	// If AcceleratedNetworking is set to true with a VMSize that does not support it, Azure will return an error.
	// +kubebuilder:validation:nullable
	// +optional
	AcceleratedNetworking *bool `json:"acceleratedNetworking,omitempty"`
}

// GetControlPlaneSubnet returns a subnet that has a role assigned to controlplane or all. Subnets with role controlplane are given higher priority.
func (n *NetworkSpec) GetControlPlaneSubnet() (SubnetSpec, error) {
	// Priority is given for subnet that have role assigned as controlplane
	if subnet, err := n.GetSubnet(SubnetControlPlane); err == nil {
		return subnet, nil
	}

	if subnet, err := n.GetSubnet(SubnetCluster); err == nil {
		return subnet, nil
	}

	return SubnetSpec{}, errors.Errorf("no subnet found with role %s", SubnetControlPlane)
}

// GetSubnet returns a subnet based on the subnet role.
func (n *NetworkSpec) GetSubnet(role SubnetRole) (SubnetSpec, error) {
	for _, sn := range n.Subnets {
		if sn.Role == role {
			return sn, nil
		}
	}
	return SubnetSpec{}, errors.Errorf("no subnet found with role %s", role)
}

// UpdateControlPlaneSubnet updates the cluster control plane subnets.
func (n *NetworkSpec) UpdateControlPlaneSubnet(subnet SubnetSpec) {
	n.UpdateSubnet(subnet, SubnetControlPlane)
	n.UpdateSubnet(subnet, SubnetCluster)
}

// UpdateSubnet updates the subnet based on the subnet role.
func (n *NetworkSpec) UpdateSubnet(subnet SubnetSpec, role SubnetRole) {
	for i, sn := range n.Subnets {
		if sn.Role == role {
			n.Subnets[i] = subnet
		}
	}
}

// IsNatGatewayEnabled returns whether or not a NAT gateway is enabled on the subnet.
func (s SubnetSpec) IsNatGatewayEnabled() bool {
	return s.NatGateway.Name != ""
}

// IsIPv6Enabled returns whether or not IPv6 is enabled on the subnet.
func (s SubnetSpec) IsIPv6Enabled() bool {
	for _, cidr := range s.CIDRBlocks {
		if net.IsIPv6CIDRString(cidr) {
			return true
		}
	}
	return false
}

// SecurityProfile specifies the Security profile settings for a
// virtual machine or virtual machine scale set.
type SecurityProfile struct {
	// This field indicates whether Host Encryption should be enabled
	// or disabled for a virtual machine or virtual machine scale set.
	// This should be disabled when SecurityEncryptionType is set to DiskWithVMGuestState.
	// Default is disabled.
	// +optional
	EncryptionAtHost *bool `json:"encryptionAtHost,omitempty"`
	// SecurityType specifies the SecurityType of the virtual machine. It has to be set to any specified value to
	// enable UefiSettings. The default behavior is: UefiSettings will not be enabled unless this property is set.
	// +kubebuilder:validation:Enum=ConfidentialVM;TrustedLaunch
	// +optional
	SecurityType SecurityTypes `json:"securityType,omitempty"`
	// UefiSettings specifies the security settings like secure boot and vTPM used while creating the virtual machine.
	// +optional
	UefiSettings *UefiSettings `json:"uefiSettings,omitempty"`
}

// UefiSettings specifies the security settings like secure boot and vTPM used while creating the virtual
// machine.
// +optional
type UefiSettings struct {
	// SecureBootEnabled specifies whether secure boot should be enabled on the virtual machine.
	// Secure Boot verifies the digital signature of all boot components and halts the boot process if signature verification fails.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is false.
	//+optional
	SecureBootEnabled *bool `json:"secureBootEnabled,omitempty"`
	// VTpmEnabled specifies whether vTPM should be enabled on the virtual machine.
	// When true it enables the virtualized trusted platform module measurements to create a known good boot integrity policy baseline.
	// The integrity policy baseline is used for comparison with measurements from subsequent VM boots to determine if anything has changed.
	// This is required to be set to Enabled if SecurityEncryptionType is defined.
	// If omitted, the platform chooses a default, which is subject to change over time, currently that default is false.
	// +optional
	VTpmEnabled *bool `json:"vTpmEnabled,omitempty"`
}

// AddressRecord specifies a DNS record mapping a hostname to an IPV4 or IPv6 address.
type AddressRecord struct {
	Hostname string
	IP       string
}

// CloudProviderConfigOverrides represents the fields that can be overridden in azure cloud provider config.
type CloudProviderConfigOverrides struct {
	// +optional
	RateLimits []RateLimitSpec `json:"rateLimits,omitempty"`
	// +optional
	BackOffs BackOffConfig `json:"backOffs,omitempty"`
}

// BackOffConfig indicates the back-off config options.
type BackOffConfig struct {
	// +optional
	CloudProviderBackoff bool `json:"cloudProviderBackoff,omitempty"`
	// +optional
	CloudProviderBackoffRetries int `json:"cloudProviderBackoffRetries,omitempty"`
	// +optional
	CloudProviderBackoffExponent *resource.Quantity `json:"cloudProviderBackoffExponent,omitempty"`
	// +optional
	CloudProviderBackoffDuration int `json:"cloudProviderBackoffDuration,omitempty"`
	// +optional
	CloudProviderBackoffJitter *resource.Quantity `json:"cloudProviderBackoffJitter,omitempty"`
}

// RateLimitSpec represents the rate limit configuration for a particular kind of resource.
// Eg. loadBalancerRateLimit is used to configure rate limits for load balancers.
// This eventually gets converted to CloudProviderRateLimitConfig that cloud-provider-azure expects.
// See: https://github.com/kubernetes-sigs/cloud-provider-azure/blob/d585c2031925b39c925624302f22f8856e29e352/pkg/provider/azure_ratelimit.go#L25
// We cannot use CloudProviderRateLimitConfig directly because floating point values are not supported in controller-tools.
// See: https://github.com/kubernetes-sigs/controller-tools/issues/245
type RateLimitSpec struct {
	// Name is the name of the rate limit spec.
	// +kubebuilder:validation:Enum=defaultRateLimit;routeRateLimit;subnetsRateLimit;interfaceRateLimit;routeTableRateLimit;loadBalancerRateLimit;publicIPAddressRateLimit;securityGroupRateLimit;virtualMachineRateLimit;storageAccountRateLimit;diskRateLimit;snapshotRateLimit;virtualMachineScaleSetRateLimit;virtualMachineSizesRateLimit;availabilitySetRateLimit
	Name string `json:"name"`
	// +optional
	Config RateLimitConfig `json:"config,omitempty"`
}

// RateLimitConfig indicates the rate limit config options.
type RateLimitConfig struct {
	// +optional
	CloudProviderRateLimit bool `json:"cloudProviderRateLimit,omitempty"`
	// +optional
	CloudProviderRateLimitQPS *resource.Quantity `json:"cloudProviderRateLimitQPS,omitempty"`
	// +optional
	CloudProviderRateLimitBucket int `json:"cloudProviderRateLimitBucket,omitempty"`
	// +optional
	CloudProviderRateLimitQPSWrite *resource.Quantity `json:"cloudProviderRateLimitQPSWrite,omitempty"`
	// +optional
	CloudProviderRateLimitBucketWrite int `json:"cloudProviderRateLimitBucketWrite,omitempty"`
}

const (
	// DefaultRateLimit ...
	DefaultRateLimit = "defaultRateLimit"
	// RouteRateLimit ...
	RouteRateLimit = "routeRateLimit"
	// SubnetsRateLimit ...
	SubnetsRateLimit = "subnetsRateLimit"
	// InterfaceRateLimit ...
	InterfaceRateLimit = "interfaceRateLimit"
	// RouteTableRateLimit ...
	RouteTableRateLimit = "routeTableRateLimit"
	// LoadBalancerRateLimit ...
	LoadBalancerRateLimit = "loadBalancerRateLimit"
	// PublicIPAddressRateLimit ...
	PublicIPAddressRateLimit = "publicIPAddressRateLimit"
	// SecurityGroupRateLimit ...
	SecurityGroupRateLimit = "securityGroupRateLimit"
	// VirtualMachineRateLimit ...
	VirtualMachineRateLimit = "virtualMachineRateLimit"
	// StorageAccountRateLimit ...
	StorageAccountRateLimit = "storageAccountRateLimit"
	// DiskRateLimit ...
	DiskRateLimit = "diskRateLimit"
	// SnapshotRateLimit ...
	SnapshotRateLimit = "snapshotRateLimit"
	// VirtualMachineScaleSetRateLimit ...
	VirtualMachineScaleSetRateLimit = "virtualMachineScaleSetRateLimit"
	// VirtualMachineSizesRateLimit ...
	VirtualMachineSizesRateLimit = "virtualMachineSizesRateLimit"
	// AvailabilitySetRateLimit ...
	AvailabilitySetRateLimit = "availabilitySetRateLimit"
)

// BastionHostSkuName is the name of the SKU used to specify the tier of Azure Bastion Host.
type BastionHostSkuName string

const (
	// BasicBastionHostSku SKU for the Azure Bastion Host.
	BasicBastionHostSku BastionHostSkuName = "Basic"
	// StandardBastionHostSku SKU for the Azure Bastion Host.
	StandardBastionHostSku BastionHostSkuName = "Standard"
)

// BastionSpec specifies how the Bastion feature should be set up for the cluster.
type BastionSpec struct {
	// +optional
	AzureBastion *AzureBastion `json:"azureBastion,omitempty"`
}

// AzureBastion specifies how the Azure Bastion cloud component should be configured.
type AzureBastion struct {
	// +optional
	Name string `json:"name,omitempty"`
	// +optional
	Subnet SubnetSpec `json:"subnet,omitempty"`
	// +optional
	PublicIP PublicIPSpec `json:"publicIP,omitempty"`
	// BastionHostSkuName configures the tier of the Azure Bastion Host. Can be either Basic or Standard. Defaults to Basic.
	// +kubebuilder:default=Basic
	// +kubebuilder:validation:Enum=Basic;Standard
	// +optional
	Sku BastionHostSkuName `json:"sku,omitempty"`
	// EnableTunneling enables the native client support feature for the Azure Bastion Host. Defaults to false.
	// +kubebuilder:default=false
	// +optional
	EnableTunneling bool `json:"enableTunneling,omitempty"`
}

// FleetsMember defines the fleets member configuration.
// See also [AKS doc].
//
// [AKS doc]: https://learn.microsoft.com/en-us/azure/templates/microsoft.containerservice/2023-03-15-preview/fleets/members
type FleetsMember struct {
	// Name is the name of the member.
	// +optional
	Name string `json:"name,omitempty"`

	FleetsMemberClassSpec `json:",inline"`
}

// BackendPool describes the backend pool of the load balancer.
type BackendPool struct {
	// Name specifies the name of backend pool for the load balancer. If not specified, the default name will
	// be set, depending on the load balancer role.
	// +optional
	Name string `json:"name,omitempty"`
}

// IsTerminalProvisioningState returns true if the ProvisioningState is a terminal state for an Azure resource.
func IsTerminalProvisioningState(state ProvisioningState) bool {
	return state == Failed || state == Succeeded
}

// Diagnostics is used to configure the diagnostic settings of the virtual machine.
type Diagnostics struct {
	// Boot configures the boot diagnostics settings for the virtual machine.
	// This allows to configure capturing serial output from the virtual machine on boot.
	// This is useful for debugging software based launch issues.
	// If not specified then Boot diagnostics (Managed) will be enabled.
	// +optional
	Boot *BootDiagnostics `json:"boot,omitempty"`
}

// BootDiagnostics configures the boot diagnostics settings for the virtual machine.
// This allows you to configure capturing serial output from the virtual machine on boot.
// This is useful for debugging software based launch issues.
// +union
type BootDiagnostics struct {
	// StorageAccountType determines if the storage account for storing the diagnostics data
	// should be disabled (Disabled), provisioned by Azure (Managed) or by the user (UserManaged).
	// +kubebuilder:validation:Required
	// +unionDiscriminator
	StorageAccountType BootDiagnosticsStorageAccountType `json:"storageAccountType"`

	// UserManaged provides a reference to the user-managed storage account.
	// +optional
	UserManaged *UserManagedBootDiagnostics `json:"userManaged,omitempty"`
}

// BootDiagnosticsStorageAccountType defines the list of valid storage account types
// for the boot diagnostics.
// +kubebuilder:validation:Enum:="Managed";"UserManaged";"Disabled"
type BootDiagnosticsStorageAccountType string

const (
	// DisabledDiagnosticsStorage is used to determine that the diagnostics storage account
	// should be disabled.
	DisabledDiagnosticsStorage BootDiagnosticsStorageAccountType = "Disabled"

	// ManagedDiagnosticsStorage is used to determine that the diagnostics storage account
	// should be provisioned by Azure.
	ManagedDiagnosticsStorage BootDiagnosticsStorageAccountType = "Managed"

	// UserManagedDiagnosticsStorage is used to determine that the diagnostics storage account
	// should be provisioned by the User.
	UserManagedDiagnosticsStorage BootDiagnosticsStorageAccountType = "UserManaged"
)

// UserManagedBootDiagnostics provides a reference to a user-managed
// storage account.
type UserManagedBootDiagnostics struct {
	// StorageAccountURI is the URI of the user-managed storage account.
	// The URI typically will be `https://<mystorageaccountname>.blob.core.windows.net/`
	// but may differ if you are using Azure DNS zone endpoints.
	// You can find the correct endpoint by looking for the Blob Primary Endpoint in the
	// endpoints tab in the Azure console or with the CLI by issuing
	// `az storage account list --query='[].{name: name, "resource group": resourceGroup, "blob endpoint": primaryEndpoints.blob}'`.
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Pattern=`^https://`
	// +kubebuilder:validation:MaxLength=1024
	StorageAccountURI string `json:"storageAccountURI"`
}

// OrchestrationModeType represents the orchestration mode for a Virtual Machine Scale Set backing an AzureMachinePool.
// +kubebuilder:validation:Enum=Flexible;Uniform
type OrchestrationModeType string

const (
	// FlexibleOrchestrationMode treats VMs as individual resources accessible by standard VM APIs.
	FlexibleOrchestrationMode OrchestrationModeType = "Flexible"
	// UniformOrchestrationMode treats VMs as identical instances accessible by the VMSS VM API.
	UniformOrchestrationMode OrchestrationModeType = "Uniform"
)

// ExtensionPlan represents the plan for an AKS marketplace extension.
type ExtensionPlan struct {
	// Name is the user-defined name of the 3rd Party Artifact that is being procured.
	// +optional
	Name string `json:"name,omitempty"`

	// Product is the name of the 3rd Party artifact that is being procured.
	// +optional
	Product string `json:"product,omitempty"`

	// PromotionCode is a publisher-provided promotion code as provisioned in Data Market for the said product/artifact.
	// +optional
	PromotionCode string `json:"promotionCode,omitempty"`

	// Publisher is the name of the publisher of the 3rd Party Artifact that is being bought.
	// +optional
	Publisher string `json:"publisher,omitempty"`

	// Version is the version of the plan.
	// +optional
	Version string `json:"version,omitempty"`
}

// ExtensionScope defines the scope of the AKS marketplace extension, if configured.
type ExtensionScope struct {
	// ScopeType is the scope of the extension. It can be either Cluster or Namespace, but not both.
	ScopeType ExtensionScopeType `json:"scopeType"`

	// ReleaseNamespace is the namespace where the extension Release must be placed, for a Cluster-scoped extension.
	// Required for Cluster-scoped extensions.
	// +optional
	ReleaseNamespace string `json:"releaseNamespace,omitempty"`

	// TargetNamespace is the namespace where the extension will be created for a Namespace-scoped extension.
	// Required for Namespace-scoped extensions.
	// +optional
	TargetNamespace string `json:"targetNamespace,omitempty"`
}

// ExtensionScopeType defines the scope type of the AKS marketplace extension, if configured.
// +kubebuilder:validation:Enum=Cluster;Namespace
type ExtensionScopeType string

const (
	// ExtensionScopeCluster ...
	ExtensionScopeCluster ExtensionScopeType = "Cluster"
	// ExtensionScopeNamespace ...
	ExtensionScopeNamespace ExtensionScopeType = "Namespace"
)

// ExtensionIdentity defines the identity of the AKS marketplace extension, if configured.
// +kubebuilder:validation:Enum=SystemAssigned
type ExtensionIdentity string

const (
	// ExtensionIdentitySystemAssigned ...
	ExtensionIdentitySystemAssigned ExtensionIdentity = "SystemAssigned"
)

// AKSAssignedIdentity defines the AKS assigned-identity of the aks marketplace extension, if configured.
// +kubebuilder:validation:Enum=SystemAssigned;UserAssigned
type AKSAssignedIdentity string

const (
	// AKSAssignedIdentitySystemAssigned ...
	AKSAssignedIdentitySystemAssigned AKSAssignedIdentity = "SystemAssigned"

	// AKSAssignedIdentityUserAssigned ...
	AKSAssignedIdentityUserAssigned AKSAssignedIdentity = "UserAssigned"
)
