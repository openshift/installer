package azure

import (
	"fmt"
	"strings"

	capz "sigs.k8s.io/cluster-api-provider-azure/api/v1beta1"

	"github.com/openshift/installer/pkg/types/dns"
)

// aro is a setting to enable aro-only modifications
var aro bool

// OutboundType is a strategy for how egress from cluster is achieved.
// +kubebuilder:validation:Enum="";Loadbalancer;NATGatewaySingleZone;NATGatewayMultiZone;UserDefinedRouting
type OutboundType string

const (
	// LoadbalancerOutboundType uses Standard loadbalancer for egress from the cluster.
	// see https://docs.microsoft.com/en-us/azure/load-balancer/load-balancer-outbound-connections#lb
	LoadbalancerOutboundType OutboundType = "Loadbalancer"

	// NATGatewaySingleZoneOutboundType uses a single (non-zone-resilient) NAT Gateway for compute node outbound access.
	// see https://learn.microsoft.com/en-us/azure/virtual-network/nat-gateway/nat-gateway-resource
	NATGatewaySingleZoneOutboundType OutboundType = "NATGatewaySingleZone"

	// NATGatewayMultiZoneOutboundType uses NAT gateways in multiple zones in the compute node subnets for outbound access.
	NATGatewayMultiZoneOutboundType OutboundType = "NATGatewayMultiZone"

	// UserDefinedRoutingOutboundType uses user defined routing for egress from the cluster.
	// see https://docs.microsoft.com/en-us/azure/virtual-network/virtual-networks-udr-overview
	UserDefinedRoutingOutboundType OutboundType = "UserDefinedRouting"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {
	// Region specifies the Azure region where the cluster will be created.
	Region string `json:"region"`

	// ARMEndpoint is the endpoint for the Azure API when installing on Azure Stack.
	ARMEndpoint string `json:"armEndpoint,omitempty"`

	// ClusterOSImage is the url of a storage blob in the Azure Stack environment containing an RHCOS VHD. This field is required for Azure Stack and not applicable to Azure.
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// BaseDomainResourceGroupName specifies the resource group where the Azure DNS zone for the base domain is found. This field is optional when creating a private cluster, otherwise required.
	//
	// +optional
	BaseDomainResourceGroupName string `json:"baseDomainResourceGroupName,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Azure for machine pools which do not define their own
	// platform configuration.
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// NetworkResourceGroupName specifies the network resource group that contains an existing VNet
	//
	// +optional
	NetworkResourceGroupName string `json:"networkResourceGroupName,omitempty"`

	// VirtualNetwork specifies the name of an existing VNet for the installer to use
	//
	// +optional
	VirtualNetwork string `json:"virtualNetwork,omitempty"`

	// ControlPlaneSubnet specifies an existing subnet for use by the control plane nodes
	//
	// Deprecated: use platform.Azure.Subnets section
	// +optional
	DeprecatedControlPlaneSubnet string `json:"controlPlaneSubnet,omitempty"`

	// ComputeSubnet specifies an existing subnet for use by compute nodes
	//
	// Deprecated: use platform.Azure.Subnets section
	// +optional
	DeprecatedComputeSubnet string `json:"computeSubnet,omitempty"`

	// cloudName is the name of the Azure cloud environment which can be used to configure the Azure SDK
	// with the appropriate Azure API endpoints.
	// If empty, the value is equal to "AzurePublicCloud".
	// +optional
	CloudName CloudEnvironment `json:"cloudName,omitempty"`

	// OutboundType is a strategy for how egress from cluster is achieved. When not specified default is "Loadbalancer".
	//
	// +kubebuilder:default=Loadbalancer
	// +optional
	OutboundType OutboundType `json:"outboundType"`

	// Subnets is the list of subnets the user can bring into the cluster to be used.
	//
	// +optional
	Subnets []SubnetSpec `json:"subnets,omitempty"`

	// ResourceGroupName is the name of an already existing resource group where the cluster should be installed.
	// This resource group should only be used for this specific cluster and the cluster components will assume
	// ownership of all resources in the resource group. Destroying the cluster using installer will delete this
	// resource group.
	// This resource group must be empty with no other resources when trying to use it for creating a cluster.
	// If empty, a new resource group will created for the cluster.
	//
	// +optional
	ResourceGroupName string `json:"resourceGroupName,omitempty"`

	// UserTags has additional keys and values that the installer will add
	// as tags to all resources that it creates on AzurePublicCloud alone.
	// Resources created by the cluster itself may not include these tags.
	// +optional
	UserTags map[string]string `json:"userTags,omitempty"`

	// CustomerManagedKey has the keys needed to encrypt the storage account.
	CustomerManagedKey *CustomerManagedKey `json:"customerManagedKey,omitempty"`

	// UserProvisionedDNS indicates if the customer is providing their own DNS solution in place of the default
	// provisioned by the Installer.
	// +kubebuilder:default:="Disabled"
	// +default="Disabled"
	// +kubebuilder:validation:Enum="Enabled";"Disabled"
	UserProvisionedDNS dns.UserProvisionedDNS `json:"userProvisionedDNS,omitempty"`
}

// SubnetSpec specifies the properties the subnet needs to be used in the cluster.
type SubnetSpec struct {
	// Name of the subnet.
	Name string `json:"name"`
	// Role specifies the actual role which the subnet should be used in.
	// +kubebuilder:validation:Enum=node;control-plane
	Role capz.SubnetRole `json:"role"`
}

// KeyVault defines an Azure Key Vault.
type KeyVault struct {
	// ResourceGroup defines the Azure resource group used by the key
	// vault.
	ResourceGroup string `json:"resourceGroup"`
	// Name is the name of the key vault.
	Name string `json:"name"`
	// KeyName is the name of the key vault key.
	KeyName string `json:"keyName"`
}

// CustomerManagedKey defines the customer managed key settings for encryption of the Azure storage account.
type CustomerManagedKey struct {
	// KeyVault is the keyvault used for the customer created key required for encryption.
	KeyVault KeyVault `json:"keyVault,omitempty"`
	// UserAssignedIdentityKey is the name of the user identity that has access to the managed key.
	UserAssignedIdentityKey string `json:"userAssignedIdentityKey,omitempty"`
}

// CloudEnvironment is the name of the Azure cloud environment
// +kubebuilder:validation:Enum="";AzurePublicCloud;AzureUSGovernmentCloud;AzureChinaCloud;AzureGermanCloud;AzureStackCloud
type CloudEnvironment string

const (
	// PublicCloud is the general-purpose, public Azure cloud environment.
	PublicCloud CloudEnvironment = "AzurePublicCloud"

	// USGovernmentCloud is the Azure cloud environment for the US government.
	USGovernmentCloud CloudEnvironment = "AzureUSGovernmentCloud"

	// ChinaCloud is the Azure cloud environment used in China.
	ChinaCloud CloudEnvironment = "AzureChinaCloud"

	// GermanCloud is the Azure cloud environment used in Germany.
	GermanCloud CloudEnvironment = "AzureGermanCloud"

	// StackCloud is the Azure cloud environment used at the edge and on premises.
	StackCloud CloudEnvironment = "AzureStackCloud"
)

// Name returns name that Azure uses for the cloud environment.
// See https://github.com/Azure/go-autorest/blob/ec5f4903f77ed9927ac95b19ab8e44ada64c1356/autorest/azure/environments.go#L13
func (e CloudEnvironment) Name() string {
	return string(e)
}

// SetBaseDomain parses the baseDomainID and sets the related fields on azure.Platform
func (p *Platform) SetBaseDomain(baseDomainID string) error {
	parts := strings.Split(baseDomainID, "/")
	p.BaseDomainResourceGroupName = parts[4]
	return nil
}

// ClusterResourceGroupName returns the name of the resource group for the cluster.
func (p *Platform) ClusterResourceGroupName(infraID string) string {
	if len(p.ResourceGroupName) > 0 {
		return p.ResourceGroupName
	}
	return fmt.Sprintf("%s-rg", infraID)
}

// VirtualNetworkName returns the name of the virtual network for the cluster.
func (p *Platform) VirtualNetworkName(infraID string) string {
	if len(p.VirtualNetwork) > 0 {
		return p.VirtualNetwork
	}
	return fmt.Sprintf("%s-vnet", infraID)
}

// ControlPlaneSubnetName returns the name of the control plane subnet for the
// cluster.
func (p *Platform) ControlPlaneSubnetName(infraID string) string {
	return fmt.Sprintf("%s-master-subnet", infraID)
}

// ComputeSubnetName returns the name of the compute subnet for the cluster.
func (p *Platform) ComputeSubnetName(infraID string) string {
	return fmt.Sprintf("%s-worker-subnet", infraID)
}

// NetworkSecurityGroupName returns the name of the network security group.
func (p *Platform) NetworkSecurityGroupName(infraID string) string {
	return fmt.Sprintf("%s-nsg", infraID)
}

// GetStorageAccountName takes an infraID and generates a
// storage account name, which can't be more than 24 characters.
func GetStorageAccountName(infraID string) string {
	storageAccountNameMax := 24

	storageAccountName := strings.ReplaceAll(infraID, "-", "")
	if len(storageAccountName) > storageAccountNameMax-2 {
		storageAccountName = storageAccountName[:storageAccountNameMax-2]
	}
	storageAccountName = fmt.Sprintf("%ssa", storageAccountName)

	return storageAccountName
}
