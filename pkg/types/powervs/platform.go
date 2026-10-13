package powervs

import (
	configv1 "github.com/openshift/api/config/v1"
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

	// ServiceEndpoints is a list which contains custom endpoints to override default
	// service endpoints of IBM Cloud Services.
	// There must only be one ServiceEndpoint for a service (no duplicates).
	//
	// +optional
	ServiceEndpoints []configv1.PowerVSServiceEndpoint `json:"serviceEndpoints,omitempty"`

	// tgName is the name or id of a pre-created TransitGateway inside IBM Cloud.
	//
	// +optional
	TransitGateway string `json:"tgName,omitempty"`
}
