package powervs

import (
	configv1 "github.com/openshift/api/config/v1"
)

// NetworkType defines which PowerVS network provisioning strategy to use.
// +kubebuilder:validation:Enum=DHCPSubnet;DHCPServer
type NetworkType string

const (
	// NetworkTypeDHCPSubnet creates a plain PowerVS network with DHCP natively enabled
	// (NetworkCreate.EnableDHCP). No separate DHCP server VM is created.
	// This is the default and recommended mode for new clusters.
	NetworkTypeDHCPSubnet NetworkType = "DHCPSubnet"

	// NetworkTypeDHCPServer creates a dedicated PowerVS DHCP Service resource (a managed VM)
	// and its associated private network. This is the legacy mode used by clusters created
	// before OPENSHIFTP-855.
	// When this mode is selected the route-mode NLB, VPC Routing Table, and Public Gateway
	// are NOT provisioned, because the DHCP Server handles its own SNAT/egress.
	NetworkTypeDHCPServer NetworkType = "DHCPServer"
)

// Platform stores all the global configuration that all machinesets
// use.
type Platform struct {

	// PowerVSResourceGroup is the resource group in which Power VS resources will be created.
	PowerVSResourceGroup string `json:"powervsResourceGroup"`

	// Region specifies the IBM Cloud colo region where the cluster will be created.
	Region string `json:"region,omitempty"`

	// Zone specifies the IBM Cloud colo region where the cluster will be created.
	// At this time, only single-zone clusters are supported.
	Zone string `json:"zone"`

	// VPCRegion specifies the IBM Cloud region in which to create VPC resources.
	// Leave unset to allow installer to select the closest VPC region.
	//
	// +optional
	VPCRegion string `json:"vpcRegion,omitempty"`

	// UserID is the login for the user's IBM Cloud account.
	UserID string `json:"userID"`

	// vpcName is the name or id of a pre-created VPC inside IBM Cloud.
	//
	// +optional
	VPC string `json:"vpcName,omitempty"`

	// VPCSubnets specifies existing subnets (by ID) where cluster
	// resources will be created.  Leave unset to have the installer
	// create subnets in a new VPC on your behalf.
	//
	// +optional
	VPCSubnets []string `json:"vpcSubnets,omitempty"`

	// ClusterOSImage is a pre-created Power VS boot image that overrides the
	// default image for cluster nodes.
	//
	// +optional
	ClusterOSImage string `json:"clusterOSImage,omitempty"`

	// DefaultMachinePlatform is the default configuration used when
	// installing on Power VS for machine pools which do not define their own
	// platform configuration.
	//
	// +optional
	DefaultMachinePlatform *MachinePool `json:"defaultMachinePlatform,omitempty"`

	// ServiceInstanceGUID is the GUID of the Power IAAS instance created from the IBM Cloud Catalog
	// before the cluster is completed.  Leave unset to allow the installer to create a service
	// instance during cluster creation.
	//
	// +optional
	ServiceInstanceGUID string `json:"serviceInstanceGUID,omitempty"`

	// COSInstanceCRN is the CRN of an existing IBM Cloud Object Storage instance
	// to use for storing the RHCOS image and bootstrap ignition data.
	// If not specified, a new COS instance will be created for the cluster.
	// The CRN format is: crn:v1:bluemix:public:cloud-object-storage:global:a/<account_id>:<instance_id>::
	//
	// +optional
	COSInstanceCRN string `json:"cosInstanceCRN,omitempty"`

	// ServiceEndpoints is a list which contains custom endpoints to override default
	// service endpoints of IBM Cloud Services.
	// There must only be one ServiceEndpoint for a service (no duplicates).
	//
	// +optional
	ServiceEndpoints []configv1.PowerVSServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// ProvisionTimeoutMinutes overrides the default machine provisioning timeout.
	// Use this when the target environment provisions machines more slowly than
	// the default 15-minute window (e.g. IBM Cloud staging/test environments).
	// When unset or zero, the default of 15 minutes is used.
	//
	// +optional
	ProvisionTimeoutMinutes int `json:"provisionTimeoutMinutes,omitempty"`

	// tgName is the name or id of a pre-created TransitGateway inside IBM Cloud.
	//
	// +optional
	TransitGateway string `json:"tgName,omitempty"`

	// NetworkType selects the PowerVS network provisioning strategy.
	// Allowed values are "DHCPSubnet" (default) and "DHCPServer" (legacy).
	//
	// DHCPSubnet creates a plain PowerVS network with DHCP natively enabled and is the
	// recommended mode. It enables the full new networking model: route-mode Network Load
	// Balancer, VPC Routing Table advertising 0.0.0.0/0, and a VPC Public Gateway.
	//
	// DHCPServer creates a dedicated DHCP Service VM (legacy). When this mode is selected
	// the route-mode NLB, Routing Table, and Public Gateway are suppressed because the
	// DHCP Server manages its own SNAT and egress.
	//
	// Leave unset to use DHCPSubnet.
	//
	// +optional
	NetworkType NetworkType `json:"networkType,omitempty"`
}
