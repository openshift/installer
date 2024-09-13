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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	capiv1beta1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

const (
	// IBMPowerVSClusterFinalizer allows IBMPowerVSClusterReconciler to clean up resources associated with IBMPowerVSCluster before
	// removing it from the apiserver.
	IBMPowerVSClusterFinalizer = "ibmpowervscluster.infrastructure.cluster.x-k8s.io"
)

// IBMPowerVSClusterSpec defines the desired state of IBMPowerVSCluster.
type IBMPowerVSClusterSpec struct {
	// ServiceInstanceID is the id of the power cloud instance where the vsi instance will get deployed.
	// Deprecated: use ServiceInstance instead
	ServiceInstanceID string `json:"serviceInstanceID"`

	// Network is the reference to the Network to use for this cluster.
	// when the field is omitted, A DHCP service will be created in the Power VS workspace and its private network will be used.
	// the DHCP service created network will have the following name format
	// 1. in the case of DHCPServer.Name is not set the name will be DHCPSERVER<CLUSTER_NAME>_Private.
	// 2. if DHCPServer.Name is set the name will be DHCPSERVER<DHCPServer.Name>_Private.
	// when Network.ID is set, its expected that there exist a network in PowerVS workspace with id or else system will give error.
	// when Network.Name is set, system will first check for network with Name in PowerVS workspace, if not exist network will be created by DHCP service.
	// Network.RegEx is not yet supported and system will ignore the value.
	Network IBMPowerVSResourceReference `json:"network"`

	// dhcpServer is contains the configuration to be used while creating a new DHCP server in PowerVS workspace.
	// when the field is omitted, CLUSTER_NAME will be used as DHCPServer.Name and DHCP server will be created.
	// it will automatically create network with name DHCPSERVER<DHCPServer.Name>_Private in PowerVS workspace.
	// +optional
	DHCPServer *DHCPServer `json:"dhcpServer,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint capiv1beta1.APIEndpoint `json:"controlPlaneEndpoint"`

	// serviceInstance is the reference to the Power VS server workspace on which the server instance(VM) will be created.
	// Power VS server workspace is a container for all Power VS instances at a specific geographic region.
	// serviceInstance can be created via IBM Cloud catalog or CLI.
	// supported serviceInstance identifier in PowerVSResource are Name and ID and that can be obtained from IBM Cloud UI or IBM Cloud cli.
	// More detail about Power VS service instance.
	// https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server
	// when omitted system will dynamically create the service instance with name CLUSTER_NAME-serviceInstance.
	// when ServiceInstance.ID is set, its expected that there exist a service instance in PowerVS workspace with id or else system will give error.
	// when ServiceInstance.Name is set, system will first check for service instance with Name in PowerVS workspace, if not exist system will create new instance.
	// if there are more than one service instance exist with the ServiceInstance.Name in given Zone, installation fails with an error. Use ServiceInstance.ID in those situations to use the specific service instance.
	// ServiceInstance.Regex is not yet supported not yet supported and system will ignore the value.
	// +optional
	ServiceInstance *IBMPowerVSResourceReference `json:"serviceInstance,omitempty"`

	// zone is the name of Power VS zone where the cluster will be created
	// possible values can be found here https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server.
	// when powervs.cluster.x-k8s.io/create-infra=true annotation is set on IBMPowerVSCluster resource,
	// 1. it is expected to set the zone, not setting will result in webhook error.
	// 2. the zone should have PER capabilities, or else system will give error.
	// +optional
	Zone *string `json:"zone,omitempty"`

	// resourceGroup name under which the resources will be created.
	// when powervs.cluster.x-k8s.io/create-infra=true annotation is set on IBMPowerVSCluster resource,
	// 1. it is expected to set the ResourceGroup.Name, not setting will result in webhook error.
	// ResourceGroup.ID and ResourceGroup.Regex is not yet supported and system will ignore the value.
	// +optional
	ResourceGroup *IBMPowerVSResourceReference `json:"resourceGroup,omitempty"`

	// vpc contains information about IBM Cloud VPC resources.
	// when omitted system will dynamically create the VPC with name CLUSTER_NAME-vpc.
	// when VPC.ID is set, its expected that there exist a VPC with ID or else system will give error.
	// when VPC.Name is set, system will first check for VPC with Name, if not exist system will create new VPC.
	// when powervs.cluster.x-k8s.io/create-infra=true annotation is set on IBMPowerVSCluster resource,
	// 1. it is expected to set the VPC.Region, not setting will result in webhook error.
	// +optional
	VPC *VPCResourceReference `json:"vpc,omitempty"`

	// vpcSubnets contains information about IBM Cloud VPC Subnet resources.
	// when omitted system will create the subnets in all the zone corresponding to VPC.Region, with name CLUSTER_NAME-vpcsubnet-ZONE_NAME.
	// possible values can be found here https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server.
	// when VPCSubnets[].ID is set, its expected that there exist a subnet with ID or else system will give error.
	// when VPCSubnets[].Zone is not set, a random zone is picked from available zones of VPC.Region.
	// when VPCSubnets[].Name is not set, system will set name as CLUSTER_NAME-vpcsubnet-INDEX.
	// if subnet with name VPCSubnets[].Name not found, system will create new subnet in VPCSubnets[].Zone.
	// +optional
	VPCSubnets []Subnet `json:"vpcSubnets,omitempty"`

	// VPCSecurityGroups to attach it to the VPC resource
	// +optional
	VPCSecurityGroups []VPCSecurityGroup `json:"vpcSecurityGroups,omitempty"`

	// transitGateway contains information about IBM Cloud TransitGateway
	// IBM Cloud TransitGateway helps in establishing network connectivity between IBM Cloud Power VS and VPC infrastructure
	// more information about TransitGateway can be found here https://www.ibm.com/products/transit-gateway.
	// when TransitGateway.ID is set, its expected that there exist a TransitGateway with ID or else system will give error.
	// when TransitGateway.Name is set, system will first check for TransitGateway with Name, if not exist system will create new TransitGateway.
	// +optional
	TransitGateway *TransitGateway `json:"transitGateway,omitempty"`

	// loadBalancers is optional configuration for configuring loadbalancers to control plane or data plane nodes.
	// when omitted system will create a default public loadbalancer with name CLUSTER_NAME-loadbalancer.
	// when specified a vpc loadbalancer will be created and controlPlaneEndpoint will be set with associated hostname of loadbalancer.
	// ControlPlaneEndpoint will be set with associated hostname of public loadbalancer.
	// when LoadBalancers[].ID is set, its expected that there exist a loadbalancer with ID or else system will give error.
	// when LoadBalancers[].Name is set, system will first check for loadbalancer with Name, if not exist system will create new loadbalancer.
	// For each loadbalancer a default backed pool and front listener will be configured with port 6443.
	// +optional
	LoadBalancers []VPCLoadBalancerSpec `json:"loadBalancers,omitempty"`

	// cosInstance contains options to configure a supporting IBM Cloud COS bucket for this
	// cluster - currently used for nodes requiring Ignition
	// (https://coreos.github.io/ignition/) for bootstrapping (requires
	// BootstrapFormatIgnition feature flag to be enabled).
	// when powervs.cluster.x-k8s.io/create-infra=true annotation is set on IBMPowerVSCluster resource and Ignition is set, then
	// 1. CosInstance.Name should be set not setting will result in webhook error.
	// 2. CosInstance.BucketName should be set not setting will result in webhook error.
	// 3. CosInstance.BucketRegion should be set not setting will result in webhook error.
	// +optional
	CosInstance *CosInstance `json:"cosInstance,omitempty"`

	// Ignition defined options related to the bootstrapping systems where Ignition is used.
	// +optional
	Ignition *Ignition `json:"ignition,omitempty"`
}

// Ignition defines options related to the bootstrapping systems where Ignition is used.
type Ignition struct {
	// Version defines which version of Ignition will be used to generate bootstrap data.
	//
	// +optional
	// +kubebuilder:default="2.3"
	// +kubebuilder:validation:Enum="2.3";"2.4";"3.0";"3.1";"3.2";"3.3";"3.4"
	Version string `json:"version,omitempty"`
}

// DHCPServer contains the DHCP server configurations.
type DHCPServer struct {
	// Optional cidr for DHCP private network
	Cidr *string `json:"cidr,omitempty"`

	// Optional DNS Server for DHCP service
	// +kubebuilder:default="1.1.1.1"
	DNSServer *string `json:"dnsServer,omitempty"`

	// Optional name of DHCP Service. Only alphanumeric characters and dashes are allowed.
	Name *string `json:"name,omitempty"`

	// Optional id of the existing DHCPServer
	ID *string `json:"id,omitempty"`

	// Optional indicates if SNAT will be enabled for DHCP service
	// +kubebuilder:default=true
	Snat *bool `json:"snat,omitempty"`
}

// ResourceReference identifies a resource with id.
type ResourceReference struct {
	// id represents the id of the resource.
	ID *string `json:"id,omitempty"`
	// +kubebuilder:default=false
	// controllerCreated indicates whether the resource is created by the controller.
	ControllerCreated *bool `json:"controllerCreated,omitempty"`
}

// TransitGatewayStatus defines the status of transit gateway as well as it's connection's status.
type TransitGatewayStatus struct {
	// id represents the id of the resource.
	ID *string `json:"id,omitempty"`
	// +kubebuilder:default=false
	// controllerCreated indicates whether the resource is created by the controller.
	ControllerCreated *bool `json:"controllerCreated,omitempty"`
	// vpcConnection defines the vpc connection status in transit gateway.
	VPCConnection *ResourceReference `json:"vpcConnection,omitempty"`
	// powerVSConnection defines the powervs connection status in transit gateway.
	PowerVSConnection *ResourceReference `json:"powerVSConnection,omitempty"`
}

// IBMPowerVSClusterStatus defines the observed state of IBMPowerVSCluster.
type IBMPowerVSClusterStatus struct {
	// ready is true when the provider resource is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// ResourceGroup is the reference to the Power VS resource group under which the resources will be created.
	ResourceGroup *ResourceReference `json:"resourceGroupID,omitempty"`

	// serviceInstance is the reference to the Power VS service on which the server instance(VM) will be created.
	ServiceInstance *ResourceReference `json:"serviceInstance,omitempty"`

	// networkID is the reference to the Power VS network to use for this cluster.
	Network *ResourceReference `json:"network,omitempty"`

	// dhcpServer is the reference to the Power VS DHCP server.
	DHCPServer *ResourceReference `json:"dhcpServer,omitempty"`

	// vpc is reference to IBM Cloud VPC resources.
	VPC *ResourceReference `json:"vpc,omitempty"`

	// vpcSubnet is reference to IBM Cloud VPC subnet.
	VPCSubnet map[string]ResourceReference `json:"vpcSubnet,omitempty"`

	// vpcSecurityGroups is reference to IBM Cloud VPC security group.
	VPCSecurityGroups map[string]VPCSecurityGroupStatus `json:"vpcSecurityGroups,omitempty"`

	// transitGateway is reference to IBM Cloud TransitGateway.
	TransitGateway *TransitGatewayStatus `json:"transitGateway,omitempty"`

	// cosInstance is reference to IBM Cloud COS Instance resource.
	COSInstance *ResourceReference `json:"cosInstance,omitempty"`

	// loadBalancers reference to IBM Cloud VPC Loadbalancer.
	LoadBalancers map[string]VPCLoadBalancerStatus `json:"loadBalancers,omitempty"`

	// Conditions defines current service state of the IBMPowerVSCluster.
	Conditions capiv1beta1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this IBMPowerVSCluster belongs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSCluster"
// +kubebuilder:printcolumn:name="PowerVS Cloud Instance ID",type="string",priority=1,JSONPath=".spec.serviceInstanceID"
// +kubebuilder:printcolumn:name="Endpoint",type="string",priority=1,JSONPath=".spec.controlPlaneEndpoint.host",description="Control Plane Endpoint"
// +kubebuilder:printcolumn:name="Port",type="string",priority=1,JSONPath=".spec.controlPlaneEndpoint.port",description="Control Plane Port"

// IBMPowerVSCluster is the Schema for the ibmpowervsclusters API.
type IBMPowerVSCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   IBMPowerVSClusterSpec   `json:"spec,omitempty"`
	Status IBMPowerVSClusterStatus `json:"status,omitempty"`
}

//+kubebuilder:object:root=true

// IBMPowerVSClusterList contains a list of IBMPowerVSCluster.
type IBMPowerVSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []IBMPowerVSCluster `json:"items"`
}

// TransitGateway holds the TransitGateway information.
type TransitGateway struct {
	// name of resource.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`
	// id of resource.
	// +optional
	ID *string `json:"id,omitempty"`
	// globalRouting indicates whether to set global routing true or not while creating the transit gateway.
	// set this field to true only when PowerVS and VPC are from different regions, if they are same it's suggested to use local routing by setting the field to false.
	// when the field is omitted,  based on PowerVS region (region associated with IBMPowerVSCluster.Spec.Zone) and VPC region(IBMPowerVSCluster.Spec.VPC.Region) system will decide whether to enable globalRouting or not.
	// +optional
	GlobalRouting *bool `json:"globalRouting,omitempty"`
}

// VPCResourceReference is a reference to a specific VPC resource by ID or Name
// Only one of ID or Name may be specified. Specifying more than one will result in
// a validation error.
type VPCResourceReference struct {
	// id of resource.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=64
	// +kubebuilder:validation:Pattern=`^[-0-9a-z_]+$`
	// +optional
	ID *string `json:"id,omitempty"`

	// name of resource.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name *string `json:"name,omitempty"`

	// region of IBM Cloud VPC.
	// when powervs.cluster.x-k8s.io/create-infra=true annotation is set on IBMPowerVSCluster resource,
	// it is expected to set the region, not setting will result in webhook error.
	Region *string `json:"region,omitempty"`
}

// CosInstance represents IBM Cloud COS instance.
type CosInstance struct {
	// name defines name of IBM cloud COS instance to be created.
	// when IBMPowerVSCluster.Ignition is set
	// +kubebuilder:validation:MinLength:=3
	// +kubebuilder:validation:MaxLength:=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`
	Name string `json:"name,omitempty"`

	// bucketName is IBM cloud COS bucket name
	BucketName string `json:"bucketName,omitempty"`

	// bucketRegion is IBM cloud COS bucket region
	BucketRegion string `json:"bucketRegion,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSCluster resource.
func (r *IBMPowerVSCluster) GetConditions() capiv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the IBMPowerVSCluster to the predescribed clusterv1.Conditions.
func (r *IBMPowerVSCluster) SetConditions(conditions capiv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// Set sets the details of the resource.
func (rf *ResourceReference) Set(resource ResourceReference) {
	rf.ID = resource.ID
	if !*rf.ControllerCreated {
		rf.ControllerCreated = resource.ControllerCreated
	}
}

func init() {
	SchemeBuilder.Register(&IBMPowerVSCluster{}, &IBMPowerVSClusterList{})
}
