/*
Copyright 2026 The Kubernetes Authors.

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

package v1beta3

import (
	"github.com/IBM/vpc-go-sdk/vpcv1"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"

	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

// ClusterTopology defines the external access architecture of the cluster.
type ClusterTopology string

const (
	// PowerVSVirtualIPTopology uses a pure PowerVS network and Virtual IP for access.
	PowerVSVirtualIPTopology ClusterTopology = "VirtualIP"

	// PowerVSLoadBalancerTopology integrates the PowerVS workspace with an IBM Cloud VPC and LoadBalancer.
	PowerVSLoadBalancerTopology ClusterTopology = "LoadBalancer"
)

const (
	// IBMPowerVSClusterFinalizer allows IBMPowerVSClusterReconciler to clean up resources associated with IBMPowerVSCluster before
	// removing it from the apiserver.
	IBMPowerVSClusterFinalizer = "ibmpowervscluster.infrastructure.cluster.x-k8s.io"
)

// SourceType defines the provisioning strategy for a resource.
type SourceType string

const (
	// SourceTypeReference indicates the controller should use an existing resource.
	SourceTypeReference SourceType = "Reference"

	// SourceTypeProvision indicates the controller should create a new resource.
	SourceTypeProvision SourceType = "Provision"
)

// DHCPSnatPolicy defines the SNAT policy for the DHCP service.
type DHCPSnatPolicy string

const (
	// DHCPSnatPolicyEnabled indicates that SNAT is enabled for the DHCP service.
	DHCPSnatPolicyEnabled DHCPSnatPolicy = "Enabled"

	// DHCPSnatPolicyDisabled indicates that SNAT is disabled for the DHCP service.
	DHCPSnatPolicyDisabled DHCPSnatPolicy = "Disabled"
)

// TransitGatewayRouting defines the routing behavior for the Transit Gateway.
type TransitGatewayRouting string

const (
	// TransitGatewayRoutingLocal forces local routing.
	TransitGatewayRoutingLocal TransitGatewayRouting = "Local"

	// TransitGatewayRoutingGlobal forces global routing.
	TransitGatewayRoutingGlobal TransitGatewayRouting = "Global"
)

// LoadBalancerType defines the network visibility of the VPC Load Balancer.
// +kubebuilder:validation:Enum=Public;Private
type LoadBalancerType string

const (
	// LoadBalancerTypePublic indicates the load balancer is accessible from the internet.
	LoadBalancerTypePublic LoadBalancerType = "Public"

	// LoadBalancerTypePrivate indicates the load balancer is only accessible internally within the VPC.
	LoadBalancerTypePrivate LoadBalancerType = "Private"
)

func init() {
	objectTypes = append(objectTypes, &IBMPowerVSCluster{}, &IBMPowerVSClusterList{})
}

// IBMPowerVSClusterSpec defines the desired state of IBMPowerVSCluster.
//
// Zone Validation:
// +kubebuilder:validation:XValidation:rule="!has(self.topology) || self.topology != 'LoadBalancer' || (has(self.zone) && size(self.zone) > 0)",message="zone is required when topology is set to LoadBalancer"
//
// ResourceGroup Validation (LoadBalancer):
// +kubebuilder:validation:XValidation:rule="!has(self.topology) || self.topology != 'LoadBalancer' || (has(self.resourceGroup) && self.resourceGroup.type == 'Reference' && has(self.resourceGroup.reference) && ((has(self.resourceGroup.reference.id) && size(self.resourceGroup.reference.id) > 0) || (has(self.resourceGroup.reference.name) && size(self.resourceGroup.reference.name) > 0)))",message="resourceGroup is required and must include either an id or name when topology is set to LoadBalancer"
//
// Workspace Validation (VirtualIP):
// +kubebuilder:validation:XValidation:rule="!has(self.topology) || self.topology != 'VirtualIP' || (has(self.workspace) && self.workspace.type == 'Reference' && has(self.workspace.reference) && ((has(self.workspace.reference.id) && size(self.workspace.reference.id) > 0) || (has(self.workspace.reference.name) && size(self.workspace.reference.name) > 0)))",message="When topology is VirtualIP, workspace type must be 'Reference' and include either an id or name"
//
// Network Validation (VirtualIP):
// +kubebuilder:validation:XValidation:rule="!has(self.topology) || self.topology != 'VirtualIP' || (has(self.network) && self.network.type == 'Reference' && has(self.network.reference) && ((has(self.network.reference.id) && size(self.network.reference.id) > 0) || (has(self.network.reference.name) && size(self.network.reference.name) > 0)))",message="When topology is VirtualIP, network type must be 'Reference' and include either an id or name"
//
// TransitGateway Validation (VirtualIP):
// +kubebuilder:validation:XValidation:rule="!has(self.topology) || self.topology != 'VirtualIP' || !has(self.transitGateway)",message="TransitGateway must not be configured when topology is set to VirtualIP"
//
// Ignition Validation:
// +kubebuilder:validation:XValidation:rule="has(self.ignition) ? has(self.cosInstance) : true",message="cosInstance configuration is required when ignition is specified"
type IBMPowerVSClusterSpec struct {
	// controlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint,omitempty,omitzero"`

	// topology defines the architectural mode for external cluster access.
	// +required
	// +kubebuilder:validation:Enum=VirtualIP;LoadBalancer
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="topology is immutable once set"
	Topology ClusterTopology `json:"topology,omitempty"`

	// workspace specifies how the PowerVS workspace is sourced.
	// A PowerVS workspace is a container for PowerVS resources in a specific zone.
	// More details: https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server
	// +optional
	Workspace WorkspaceSource `json:"workspace,omitempty,omitzero"`

	// network specifies how the PowerVS network should be sourced.
	// +optional
	Network NetworkSource `json:"network,omitempty,omitzero"`

	// zone is the name of PowerVS zone where the cluster will be created
	// possible values can be found here https://cloud.ibm.com/docs/power-iaas?topic=power-iaas-creating-power-virtual-server.
	// +optional
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="zone is immutable"
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	// +kubebuilder:validation:Pattern=^[a-zA-Z0-9\-_]+$
	Zone string `json:"zone,omitempty"`

	// resourceGroup defines the IBM Cloud Resource Group for the cluster.
	// +optional
	ResourceGroup ResourceGroupSource `json:"resourceGroup,omitempty,omitzero"`

	// transitGateway contains information about the IBM Cloud TransitGateway.
	// IBM Cloud TransitGateway helps in establishing network connectivity between IBM Cloud PowerVS and VPC infrastructure.
	// This field is rejected by the API if the Topology is set to VirtualIP.
	// +optional
	TransitGateway TransitGatewaySource `json:"transitGateway,omitempty,omitzero"`

	// vpc specifies how the IBM Cloud VPC should be sourced.
	// +optional
	VPC VPCSource `json:"vpc,omitempty,omitzero"`

	// subnets configures the VPC Subnets bound to this cluster environment.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=15
	VPCSubnets []VPCSubnetSource `json:"subnets,omitempty"`

	// loadBalancers contains information about IBM Cloud VPC Load Balancer resources.
	// This field is rejected by the API if the Topology is set to VirtualIP.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=50
	LoadBalancers []LoadBalancerSource `json:"loadBalancers,omitempty"`

	// vpcSecurityGroups defines the VPC Security Groups that should exist or be created.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=25
	VPCSecurityGroups []VPCSecurityGroupSource `json:"vpcSecurityGroups,omitempty"`

	// cosInstance contains options to configure a supporting IBM Cloud COS instance and bucket
	// for this cluster. It is currently used for nodes requiring Ignition for bootstrapping.
	// +optional
	COSInstance COSInstanceSource `json:"cosInstance,omitempty,omitzero"`

	// ignition defines options related to the bootstrapping systems where Ignition is used.
	// +optional
	Ignition Ignition `json:"ignition,omitempty,omitzero"`
}

// IBMPowerVSClusterStatus defines the observed state of IBMPowerVSCluster.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSClusterStatus struct {
	// conditions represents the observations of a IBMPowerVSCluster's current state.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the IBMPowerVSCluster initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Cluster provisioning.
	// +optional
	Initialization IBMPowerVSClusterInitializationStatus `json:"initialization,omitempty,omitzero"`

	// workspace is the reference to the PowerVS workspace.
	// +optional
	Workspace ResourceReference `json:"workspace,omitempty,omitzero"`

	// network tracks the status of the PowerVS network and its associated resources.
	// +optional
	Network NetworkStatus `json:"network,omitempty,omitzero"`

	// resourceGroup is the reference to the IBM Cloud Resource Group where the cluster resources are provisioned.
	// +optional
	ResourceGroup ResourceReference `json:"resourceGroup,omitempty,omitzero"`

	// transitGateway is reference to IBM Cloud TransitGateway.
	// +optional
	TransitGateway TransitGatewayStatus `json:"transitGateway,omitempty,omitzero"`

	// vpc tracks the observed state of the provisioned or referenced IBM Cloud VPC.
	// +optional
	VPC VPCStatus `json:"vpc,omitempty,omitzero"`

	// vpcSubnets tracks the current status of the VPC subnets.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=15
	VPCSubnets []VPCSubnetStatus `json:"vpcSubnets,omitempty"`

	// loadBalancers tracks the status of the IBM Cloud VPC Load Balancers.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=50
	LoadBalancers []LoadBalancerStatus `json:"loadBalancers,omitempty"`

	// vpcSecurityGroups tracks the live observed states of all managed or referenced VPC Security Groups.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MaxItems=25
	VPCSecurityGroups []VPCSecurityGroupStatus `json:"vpcSecurityGroups,omitempty"`

	// cosInstance tracks the observed state of the provisioned or referenced IBM Cloud COS instance.
	// +optional
	COSInstance COSInstanceStatus `json:"cosInstance,omitempty,omitzero"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *IBMPowerVSClusterDeprecatedStatus `json:"deprecated,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:subresource:status
// +kubebuilder:storageversion
// +kubebuilder:resource:path=ibmpowervsclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this IBMPowerVSCluster belongs"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of IBMPowerVSCluster"
// +kubebuilder:printcolumn:name="Endpoint",type="string",priority=1,JSONPath=".spec.controlPlaneEndpoint.host",description="Control Plane Endpoint"
// +kubebuilder:printcolumn:name="Port",type="string",priority=1,JSONPath=".spec.controlPlaneEndpoint.port",description="Control Plane Port"

// IBMPowerVSCluster is the Schema for the ibmpowervsclusters API.
type IBMPowerVSCluster struct {
	metav1.TypeMeta `json:",inline"`

	// metadata is a standard object metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitzero"`

	// spec defines the desired state of IBMPowerVSCluster
	// +required
	Spec IBMPowerVSClusterSpec `json:"spec,omitempty,omitzero"`

	// status defines the observed state of IBMPowerVSCluster
	// +optional
	Status IBMPowerVSClusterStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// IBMPowerVSClusterList contains a list of IBMPowerVSCluster.
type IBMPowerVSClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitzero"`
	Items           []IBMPowerVSCluster `json:"items"`
}

// APIEndpoint represents a reachable Kubernetes API endpoint.
// +kubebuilder:validation:MinProperties=1
type APIEndpoint struct {
	// host is the hostname on which the API server is serving.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	Host string `json:"host,omitempty"`

	// port is the port on which the API server is serving.
	// +optional
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	Port int32 `json:"port,omitempty"`
}

// IBMPowerVSClusterInitializationStatus provides observations of the IBMPowerVSCluster initialization process.
// +kubebuilder:validation:MinProperties=1
type IBMPowerVSClusterInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Cluster's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Cluster provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// TransitGatewaySource holds the TransitGateway information and determines how it is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision' ? !has(self.provision) : true",message="provision configuration is forbidden when type is not Provision"
type TransitGatewaySource struct {
	// type defines whether to use an existing Transit Gateway or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	Type SourceType `json:"type,omitempty"`

	// reference contains the information to identify an existing Transit Gateway.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new Transit Gateway.
	// +optional
	Provision TransitGatewayProvision `json:"provision,omitempty,omitzero"`

	// vpcConnection defines how the VPC connection to the Transit Gateway is sourced.
	// +optional
	VPCConnection TransitGatewayConnectionSource `json:"vpcConnection,omitempty,omitzero"`

	// powerVSConnection defines how the PowerVS connection to the Transit Gateway is sourced.
	// +optional
	PowerVSConnection TransitGatewayConnectionSource `json:"powerVSConnection,omitempty,omitzero"`
}

// TransitGatewayProvision holds the configuration for a new Transit Gateway.
// +kubebuilder:validation:MinProperties=1
type TransitGatewayProvision struct {
	// name of the transit gateway to be created.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-zA-Z]|[a-zA-Z][-_a-zA-Z0-9]*[a-zA-Z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`

	// globalRouting indicates whether to use Local or Global routing.
	// If omitted, the system will automatically decide based on the PowerVS and VPC regions.
	// +kubebuilder:validation:Enum=Local;Global
	// +optional
	GlobalRouting TransitGatewayRouting `json:"globalRouting,omitempty"`
}

// TransitGatewayConnectionSource defines how a Transit Gateway connection is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision' ? !has(self.provision) : true",message="provision configuration is forbidden when type is not Provision"
type TransitGatewayConnectionSource struct {
	// type defines whether to use an existing connection or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	Type SourceType `json:"type,omitempty"`

	// reference contains the information to identify an existing connection.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new connection.
	// +optional
	Provision TransitGatewayConnectionProvision `json:"provision,omitempty,omitzero"`
}

// TransitGatewayConnectionProvision holds the configuration for a new Transit Gateway connection.
// +kubebuilder:validation:MinProperties=1
type TransitGatewayConnectionProvision struct {
	// name of the connection to be created.
	// If omitted, the system will dynamically create the connection with a default name.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`
}

// TransitGatewayStatus defines the status of the transit gateway as well as its connections.
// +kubebuilder:validation:MinProperties=1
type TransitGatewayStatus struct {
	// id represents the id of the resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name represents the name of the resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`

	// vpcConnection defines the vpc connection status in the transit gateway.
	// +optional
	VPCConnection ResourceConnectionStatus `json:"vpcConnection,omitempty,omitzero"`

	// powerVSConnection defines the powervs connection status in the transit gateway.
	// +optional
	PowerVSConnection ResourceConnectionStatus `json:"powerVSConnection,omitempty,omitzero"`
}

// ResourceConnectionStatus identifies a connection resource.
// +kubebuilder:validation:MinProperties=1
type ResourceConnectionStatus struct {
	// id represents the id of the connection resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name represents the name of the connection resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`

	// state indicates the current state of the connection (e.g., pending, attached).
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	State string `json:"state,omitempty"`
}

// Ignition defines options related to the bootstrapping systems where Ignition is used.
// +kubebuilder:validation:MinProperties=1
type Ignition struct {
	// version defines which version of Ignition will be used to generate bootstrap data.
	//
	// +optional
	// +kubebuilder:validation:Enum="2.3";"2.4";"3.0";"3.1";"3.2";"3.3";"3.4"
	Version string `json:"version,omitempty"`
}

// IBMPowerVSClusterDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSClusterDeprecatedStatus struct {
	// v1beta2 groups all the status fields that are deprecated and will be removed when support for v1beta2 will be dropped.
	// +optional
	V1Beta2 *IBMPowerVSClusterV1Beta2DeprecatedStatus `json:"v1beta2,omitempty"`
}

// IBMPowerVSClusterV1Beta2DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type IBMPowerVSClusterV1Beta2DeprecatedStatus struct {
	// conditions defines current service state of the VSphereCluster.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// ResourceReference identifies a resource with id and name.
// +kubebuilder:validation:MinProperties=1
type ResourceReference struct {
	// id represents the id of the resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name is the name of the resource.
	// When used in a list, this field acts as the unique correlation key (listMapKey)
	// to map the Status object back to its corresponding Spec definition.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`
}

// WorkspaceSource defines how the PowerVS workspace is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Provision' ? has(self.provision) : !has(self.provision)",message="provision configuration is required when type is Provision, and forbidden otherwise"
type WorkspaceSource struct {
	// type defines how the workspace is sourced.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="workspace type is immutable once set"
	Type SourceType `json:"type,omitempty"`

	// reference tells the controller to use an existing PowerVS workspace.
	// Supported identifiers are name and id.
	// If more than one workspace has the same name, use id.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision defines the configuration for creating a new PowerVS workspace.
	// +optional
	Provision WorkspaceProvisionConfig `json:"provision,omitempty,omitzero"`
}

// WorkspaceProvisionConfig defines the parameters for creating a new workspace.
// +kubebuilder:validation:MinProperties=1
type WorkspaceProvisionConfig struct {
	// name is the explicit name of the workspace to be created.
	// If omitted, the system will dynamically create the workspace with the name <CLUSTER_NAME>-workspace.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`
}

// ResourceIdentifier defines the identification of a specific PowerVS resource by ID or Name.
// +kubebuilder:validation:MinProperties=1
// +kubebuilder:validation:XValidation:rule="(has(self.id) ? 1 : 0) + (has(self.name) ? 1 : 0) == 1",message="exactly one of id or name must be specified"
type ResourceIdentifier struct {
	// id of the resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name of the resource.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`
}

// NetworkSource defines how to source the PowerVS network.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Provision' ? has(self.provision) : !has(self.provision)",message="provision configuration is required when type is Provision, and forbidden otherwise"
type NetworkSource struct {
	// type defines how the Network is sourced.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="Network type is immutable once set"
	Type SourceType `json:"type,omitempty"`

	// reference tells the controller to look up an EXISTING PowerVS network.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision provides the configuration for the controller to CREATE a new Network and DHCP Server.
	// +optional
	Provision NetworkProvisionConfig `json:"provision,omitempty,omitzero"`
}

// NetworkProvisionConfig defines the parameters for creating a new PowerVS Network.
// +kubebuilder:validation:MinProperties=1
type NetworkProvisionConfig struct {
	// dhcpServer contains the configuration for the DHCP server that will be created.
	// +optional
	DHCPServer DHCPServer `json:"dhcpServer,omitempty,omitzero"`
}

// DHCPServer contains the configuration for a NEW DHCP server.
// +kubebuilder:validation:MinProperties=1
type DHCPServer struct {
	// name is the name of the DHCP Service to be created. Only alphanumeric characters and dashes are allowed.
	// If omitted, the name will default to DHCPSERVER<CLUSTER_NAME>_Private.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`

	// cidr is the CIDR for the DHCP private network.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=49
	// +kubebuilder:validation:Pattern=`^([0-9]{1,3}\.){3}[0-9]{1,3}($|/[0-9]{1,2})$`
	CIDR string `json:"cidr,omitempty"`

	// dnsServer is the DNS Server for the DHCP service.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=45
	DNSServer string `json:"dnsServer,omitempty"`

	// snat indicates the SNAT policy for the DHCP service.
	// Allowed values are "Enabled" and "Disabled".
	// If omitted, the system will choose a Enabled policy by default.
	// +optional
	// +kubebuilder:validation:Enum=Enabled;Disabled
	Snat DHCPSnatPolicy `json:"snat,omitempty"`
}

// NetworkStatus defines the observed state of the PowerVS network and its associated components.
// +kubebuilder:validation:MinProperties=1
type NetworkStatus struct {
	// id is the unique identifier of the network.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name is the name of the network.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`

	// dhcpServer tracks the provisioned DHCP server identity, if one was created.
	// +optional
	DHCPServer ResourceReference `json:"dhcpServer,omitempty,omitzero"`
}

// ResourceGroupSource represents the source of an IBM Cloud Resource Group.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : true",message="reference configuration is required when type is Reference"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision'",message="Provisioning a Resource Group is not yet supported in this API version"
type ResourceGroupSource struct {
	// type defines the intended action for the Resource Group.
	// Currently, only "Reference" is supported.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	Type SourceType `json:"type,omitempty"`

	// reference specifies the existing Resource Group to use by Name or ID.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`
}

// VPCSource defines how the IBM Cloud VPC is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Provision' ? has(self.provision) : !has(self.provision)",message="provision configuration is required when type is Provision, and forbidden otherwise"
type VPCSource struct {
	// type defines whether to use an existing VPC or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="VPC type is immutable once set"
	Type SourceType `json:"type,omitempty"`

	// region is the IBM Cloud region where the VPC is or will be located.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	Region string `json:"region,omitempty"`

	// reference contains the information to identify an existing VPC.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new VPC.
	// +optional
	Provision VPCProvision `json:"provision,omitempty,omitzero"`
}

// VPCProvision holds the configuration for creating a new VPC.
// +kubebuilder:validation:MinProperties=1
type VPCProvision struct {
	// name of the VPC to be created.
	// If omitted, the system will dynamically create the VPC with the name <CLUSTER_NAME>-vpc.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`
}

// VPCStatus tracks the live observed state of the IBM Cloud VPC.
type VPCStatus struct {
	// id is the validated string identifier returned by the IBM Cloud API.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name is the unique name identifying the VPC in the cloud.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`

	// region is the IBM Cloud region where the VPC resides.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	Region string `json:"region,omitempty"`
}

// VPCSubnetSource defines how the IBM Cloud VPC Subnet is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision' ? !has(self.provision) : true",message="provision configuration is forbidden when type is not Provision"
type VPCSubnetSource struct {
	// type defines whether to use an existing VPC Subnet or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	Type SourceType `json:"type,omitempty"`

	// zone of the IBM Cloud VPC Subnet.
	// When provisioning, if omitted, a random zone is picked from available zones of the VPC.Region.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Zone string `json:"zone,omitempty"`

	// reference contains the information to identify an existing VPC Subnet.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new VPC Subnet.
	// +optional
	Provision VPCSubnetProvision `json:"provision,omitempty,omitzero"`
}

// VPCSubnetProvision holds the configuration for a new VPC Subnet.
// +kubebuilder:validation:MinProperties=1
type VPCSubnetProvision struct {
	// name of the VPC Subnet to be created.
	// If omitted, the system will dynamically create the VPC subnet with name <CLUSTER_NAME>-vpcsubnet-<INDEX>.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`
}

// VPCSubnetStatus defines the observed state of an IBM Cloud VPC Subnet.
type VPCSubnetStatus struct {
	// id is the validated string identifier returned by the IBM Cloud API.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// name is the unique name identifying the subnet in the cloud.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`

	// zone is the actual IBM Cloud zone where the subnet resides.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Zone string `json:"zone,omitempty"`
}

// LoadBalancerSource defines how the IBM Cloud VPC Load Balancer is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision' ? !has(self.provision) : true",message="provision configuration is forbidden when type is not Provision"
type LoadBalancerSource struct {
	// type defines whether to use an existing Load Balancer or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	Type SourceType `json:"type,omitempty"`

	// reference contains the information to identify an existing Load Balancer.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new Load Balancer.
	// +optional
	Provision LoadBalancerProvision `json:"provision,omitempty,omitzero"`
}

// LoadBalancerProvision holds the configuration for a new VPC Load Balancer.
// +kubebuilder:validation:MinProperties=1
type LoadBalancerProvision struct {
	// name sets the name of the VPC load balancer.
	// If omitted, the system will dynamically create it.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`

	// type indicates whether the load balancer is public or private.
	// When omitted, defaults to Public.
	// +optional
	Type LoadBalancerType `json:"type,omitempty"`

	// additionalListeners sets the additional listeners for the load balancer.
	// +listType=map
	// +listMapKey=port
	// +optional
	// +kubebuilder:validation:MaxItems=10
	AdditionalListeners []AdditionalListener `json:"additionalListeners,omitempty"`

	// backendPools defines the load balancer's backend pools.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=10
	BackendPools []LoadBalancerBackendPool `json:"backendPools,omitempty"`

	// securityGroups defines the Security Groups to attach to the load balancer.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MaxItems=5
	SecurityGroups []ResourceIdentifier `json:"securityGroups,omitempty"`

	// subnets defines the VPC Subnets to attach to the load balancer.
	// +optional
	// +listType=atomic
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=15
	Subnets []ResourceIdentifier `json:"subnets,omitempty"`
}

// AdditionalListener defines the desired state of an
// additional listener on a VPC load balancer.
type AdditionalListener struct {
	// defaultPoolName defines the name of a VPC Load Balancer Backend Pool to use for the VPC Load Balancer Listener.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	DefaultPoolName string `json:"defaultPoolName,omitempty"`

	// port sets the port for the additional listener.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +required
	Port int64 `json:"port,omitempty"`

	// protocol defines the protocol to use for the VPC Load Balancer Listener.
	// Will default to TCP protocol if not specified.
	// +optional
	Protocol LoadBalancerListenerProtocol `json:"protocol,omitempty"`

	// selector is used to find IBMPowerVSMachines with matching labels.
	// If the label matches, the machine is then added to the load balancer listener configuration.
	// +optional
	Selector metav1.LabelSelector `json:"selector,omitempty"`
}

// LoadBalancerBackendPool defines the desired configuration of a VPC Load Balancer Backend Pool.
type LoadBalancerBackendPool struct {
	// name defines the name of the Backend Pool.
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$`
	// +optional
	Name string `json:"name,omitempty"`

	// algorithm defines the load balancing algorithm to use.
	// +required
	Algorithm LoadBalancerBackendPoolAlgorithm `json:"algorithm,omitempty"`

	// healthMonitor defines the backend pool's health monitor.
	// +required
	HealthMonitor LoadBalancerHealthMonitor `json:"healthMonitor,omitempty,omitzero"`

	// protocol defines the protocol to use for the Backend Pool.
	// +required
	Protocol LoadBalancerBackendPoolProtocol `json:"protocol,omitempty"`
}

// LoadBalancerHealthMonitor defines the desired state of a Health Monitor resource for a VPC Load Balancer Backend Pool.
// +kubebuilder:validation:XValidation:rule="self.delay > self.timeout",message="health monitor's delay must be greater than the timeout"
type LoadBalancerHealthMonitor struct {
	// delay defines the seconds to wait between health checks.
	// +kubebuilder:validation:Minimum=2
	// +kubebuilder:validation:Maximum=60
	// +required
	Delay int64 `json:"delay,omitempty"`

	// retries defines the max retries for health check.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=10
	// +required
	Retries int64 `json:"retries,omitempty"`

	// port defines the port to perform health monitoring on.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	// +optional
	Port int64 `json:"port,omitempty"`

	// timeout defines the seconds to wait for a health check response.
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=59
	// +required
	Timeout int64 `json:"timeout,omitempty"`

	// type defines the protocol used for health checks.
	// +required
	Type LoadBalancerBackendPoolHealthMonitorType `json:"type,omitempty"`

	// urlPath defines the URL to use for health monitoring.
	// +kubebuilder:validation:Pattern=`^\/(([a-zA-Z0-9-._~!$&'()*+,;=:@]|%[a-fA-F0-9]{2})+(\/([a-zA-Z0-9-._~!$&'()*+,;=:@]|%[a-fA-F0-9]{2})*)*)?(\\?([a-zA-Z0-9-._~!$&'()*+,;=:@\/?]|%[a-fA-F0-9]{2})*)?$`
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=8192
	URLPath string `json:"urlPath,omitempty"`
}

// LoadBalancerStatus defines the status of a VPC load balancer.
type LoadBalancerStatus struct {
	// name is the unique identifier for the load balancer configuration.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	Name string `json:"name,omitempty"`

	// id of the VPC load balancer.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	ID string `json:"id,omitempty"`

	// state is the status of the load balancer.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	State LoadBalancerState `json:"state,omitempty"`

	// hostname is the hostname of load balancer.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	Hostname string `json:"hostname,omitempty"`
}

// COSInstanceSource defines how the IBM Cloud COS instance is sourced.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type != 'Provision' ? !has(self.provision) : true",message="provision configuration is forbidden when type is not Provision; it is optional when type is Provision"
type COSInstanceSource struct {
	// type defines whether to use an existing COS instance or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	// +kubebuilder:validation:XValidation:rule="self == oldSelf",message="COS instance type is immutable once set"
	Type SourceType `json:"type,omitempty"`

	// bucketName is the name of the IBM Cloud COS bucket used for Ignition bootstrapping.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	BucketName string `json:"bucketName,omitempty"`

	// bucketRegion is the IBM Cloud region where the COS bucket resides or will be created.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	BucketRegion string `json:"bucketRegion,omitempty"`

	// reference contains the information to identify an existing COS instance.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new COS instance and bucket.
	// +optional
	Provision COSInstanceProvision `json:"provision,omitempty,omitzero"`
}

// COSInstanceProvision holds the configuration for creating a new COS instance.
// +kubebuilder:validation:MinProperties=1
type COSInstanceProvision struct {
	// name defines the explicit name of the IBM Cloud COS instance to be created.
	// If omitted, the system will dynamically create it using the cluster name.
	// +kubebuilder:validation:MinLength=3
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=`^[a-z0-9][a-z0-9.-]{1,61}[a-z0-9]$`
	// +optional
	Name string `json:"name,omitempty"`
}

// COSInstanceStatus tracks the live observed state of the IBM Cloud COS instance.
type COSInstanceStatus struct {
	// id is the validated string identifier (CRN or GUID) returned by the IBM Cloud API.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=512
	ID string `json:"id,omitempty"`

	// name is the unique name identifying the COS instance in the cloud.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Name string `json:"name,omitempty"`

	// bucketName tracks the confirmed bucket used for bootstrapping.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	BucketName string `json:"bucketName,omitempty"`

	// bucketRegion tracks the confirmed region where the bucket resides.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=32
	BucketRegion string `json:"bucketRegion,omitempty"`

	// hmacSecretName is the name of the Kubernetes Secret (in the same namespace as the
	// IBMPowerVSCluster) that holds the HMAC credentials for the COS instance.
	// The Secret contains keys "access_key_id" and "secret_access_key".
	// These are used by the machine scope to generate short-lived pre-signed URLs for
	// ignition bootstrap data instead of embedding a long-lived IAM Bearer token.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=253
	HMACSecretName string `json:"hmacSecretName,omitempty"`
}

const (
	// VPCSecurityGroupRuleProtocolAnyType is a string representation of the 'SecurityGroupRuleProtocolAny' type.
	VPCSecurityGroupRuleProtocolAnyType = "*vpcv1.SecurityGroupRuleProtocolAny"

	// VPCSecurityGroupRuleProtocolIcmptcpudpType is a string representation of the 'SecurityGroupRuleProtocolIcmptcpudp' type.
	VPCSecurityGroupRuleProtocolIcmptcpudpType = "*vpcv1.SecurityGroupRuleProtocolIcmptcpudp"

	// VPCSecurityGroupRuleProtocolIcmpType is a string representation of the 'SecurityGroupRuleSecurityGroupRuleProtocolIcmp' type.
	VPCSecurityGroupRuleProtocolIcmpType = "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolIcmp"

	// VPCSecurityGroupRuleProtocolTcpudpType is a string representation of the 'SecurityGroupRuleSecurityGroupRuleProtocolTcpudp' type.
	VPCSecurityGroupRuleProtocolTcpudpType = "*vpcv1.SecurityGroupRuleSecurityGroupRuleProtocolTcpudp"

	// VPCSecurityGroupRuleProtocolIndividualType is a string representation of the 'SecurityGroupRuleProtocolIndividual' type.
	VPCSecurityGroupRuleProtocolIndividualType = "*vpcv1.SecurityGroupRuleProtocolIndividual"
)

// VPCSecurityGroupRuleDirection represents the directions for a Security Group Rule.
// +kubebuilder:validation:Enum=inbound;outbound
type VPCSecurityGroupRuleDirection string

const (
	// VPCSecurityGroupRuleDirectionInbound defines the Rule is for inbound traffic.
	VPCSecurityGroupRuleDirectionInbound VPCSecurityGroupRuleDirection = vpcv1.NetworkACLRuleDirectionInboundConst
	// VPCSecurityGroupRuleDirectionOutbound defines the Rule is for outbound traffic.
	VPCSecurityGroupRuleDirectionOutbound VPCSecurityGroupRuleDirection = vpcv1.NetworkACLRuleDirectionOutboundConst
)

// VPCSecurityGroupRuleProtocol represents the protocols for a Security Group Rule.
// +kubebuilder:validation:Pattern=`^(any|all|icmp_tcp_udp|icmp|tcp|udp|ah|esp|gre|ip_in_ip|l2tp|rsvp|sctp|vrrp|number_(?:0|2|3|5|[7-9]|1[0-6]|1[8-9]|[2-3][0-9]|4[0-5]|4[89]|5[2-9]|[6-9][0-9]|10[0-9]|11[0-1]|11[3-4]|11[6-9]|12[0-9]|13[0-1]|13[3-9]|1[4-9][0-9]|2[0-4][0-9]|25[0-5]))$`
type VPCSecurityGroupRuleProtocol string

const (
	// VPCSecurityGroupRuleProtocolAny defines the Rule is for any network protocols.
	VPCSecurityGroupRuleProtocolAny VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolAnyConst
	// VPCSecurityGroupRuleProtocolAll is DEPRECATED: Use VPCSecurityGroupRuleProtocolIcmpTCPUDP instead.
	// This constant is maintained for backward compatibility.
	// It will be removed in a future version.
	VPCSecurityGroupRuleProtocolAll VPCSecurityGroupRuleProtocol = "all"
	// VPCSecurityGroupRuleProtocolIcmpTCPUDP defines the Rule is for ICMP, TCP and UDP protocols.
	VPCSecurityGroupRuleProtocolIcmpTCPUDP VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolIcmpTCPUDPConst
	// VPCSecurityGroupRuleProtocolIcmp defiens the Rule is for ICMP network protocol.
	VPCSecurityGroupRuleProtocolIcmp VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolIcmpConst
	// VPCSecurityGroupRuleProtocolTCP defines the Rule is for TCP network protocol.
	VPCSecurityGroupRuleProtocolTCP VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolTCPConst
	// VPCSecurityGroupRuleProtocolUDP defines the Rule is for UDP network protocol.
	VPCSecurityGroupRuleProtocolUDP VPCSecurityGroupRuleProtocol = vpcv1.NetworkACLRuleProtocolUDPConst
)

// VPCSecurityGroupRuleRemoteType represents the type of Security Group Rule's destination or source is
// intended. This is intended to define the VPCSecurityGroupRulePrototype subtype.
// For example:
// - any - Any source or destination (0.0.0.0/0)
// - cidr - A CIDR representing a set of IP's (10.0.0.0/28)
// - address - A specific address (192.168.0.1)
// - sg - A Security Group.
// +kubebuilder:validation:Enum=any;cidr;address;sg
type VPCSecurityGroupRuleRemoteType string

const (
	// VPCSecurityGroupRuleRemoteTypeAny defines the destination or source for the Rule is anything/anywhere.
	VPCSecurityGroupRuleRemoteTypeAny VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("any")
	// VPCSecurityGroupRuleRemoteTypeCIDR defines the destination or source for the Rule is a CIDR block.
	VPCSecurityGroupRuleRemoteTypeCIDR VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("cidr")
	// VPCSecurityGroupRuleRemoteTypeAddress defines the destination or source for the Rule is an address.
	VPCSecurityGroupRuleRemoteTypeAddress VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("address")
	// VPCSecurityGroupRuleRemoteTypeSG defines the destination or source for the Rule is a VPC Security Group.
	VPCSecurityGroupRuleRemoteTypeSG VPCSecurityGroupRuleRemoteType = VPCSecurityGroupRuleRemoteType("sg")
)

// VPCSecurityGroupSource defines a VPC Security Group that should exist or be created.
// +kubebuilder:validation:XValidation:rule="self.type == 'Reference' ? has(self.reference) : !has(self.reference)",message="reference configuration is required when type is Reference, and forbidden otherwise"
// +kubebuilder:validation:XValidation:rule="self.type == 'Provision' ? has(self.provision) : !has(self.provision)",message="provision configuration is required when type is Provision, and forbidden otherwise"
type VPCSecurityGroupSource struct {
	// type defines whether to use an existing Security Group or provision a new one.
	// +required
	// +kubebuilder:validation:Enum=Reference;Provision
	Type SourceType `json:"type,omitempty"`

	// reference contains the information to identify an existing Security Group.
	// CAPI will not manage rules for referenced Security Groups.
	// +optional
	Reference ResourceIdentifier `json:"reference,omitempty,omitzero"`

	// provision contains the configuration for provisioning a new Security Group.
	// +optional
	Provision VPCSecurityGroupProvision `json:"provision,omitempty,omitzero"`
}

// VPCSecurityGroupProvision holds the configuration for creating a new Security Group.
// +kubebuilder:validation:MinProperties=1
type VPCSecurityGroupProvision struct {
	// name of the Security Group.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$
	Name string `json:"name,omitempty"`

	// rules are the Security Group Rules for the Security Group.
	// +optional
	// +kubebuilder:validation:MaxItems=250
	// +listType=atomic
	Rules []VPCSecurityGroupRule `json:"rules,omitempty"`

	// tags are tag names to create for the Security Group.
	// +optional
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	// +kubebuilder:validation:items:MaxLength=128
	// +kubebuilder:validation:items:Pattern=`^[ ]*[A-Za-z0-9:_.\-][A-Za-z0-9 :_.\-]*$`
	// +listType=set
	Tags []string `json:"tags,omitempty"`
}

// VPCSecurityGroupRule defines a VPC Security Group Rule for a specified Security Group.
// +kubebuilder:validation:XValidation:rule="(has(self.destination) && !has(self.source)) || (!has(self.destination) && has(self.source))",message="both destination and source cannot be provided"
// +kubebuilder:validation:XValidation:rule="self.direction == 'inbound' ? has(self.source) : true",message="source must be set for VPCSecurityGroupRuleDirectionInbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'inbound' ? !has(self.destination) : true",message="destination is not valid for VPCSecurityGroupRuleDirectionInbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'outbound' ? has(self.destination) : true",message="destination must be set for VPCSecurityGroupRuleDirectionOutbound direction"
// +kubebuilder:validation:XValidation:rule="self.direction == 'outbound' ? !has(self.source) : true",message="source is not valid for VPCSecurityGroupRuleDirectionOutbound direction"
// +kubebuilder:validation:XValidation:rule="has(self.source) ? has(self.source.protocol) : true",message="protocol is required in source"
// +kubebuilder:validation:XValidation:rule="has(self.destination) ? has(self.destination.protocol) : true",message="protocol is required in destination"
type VPCSecurityGroupRule struct {
	// destination defines the destination of outbound traffic for the Security Group Rule.
	// Only used when direction is VPCSecurityGroupRuleDirectionOutbound.
	// +optional
	Destination VPCSecurityGroupRulePrototype `json:"destination,omitempty,omitzero"`

	// direction is the direction of traffic to allow.
	// Allowable values: inbound, outbound.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	// +kubebuilder:validation:Pattern=^[a-z][a-z0-9]*(_[a-z0-9]+)*$
	Direction VPCSecurityGroupRuleDirection `json:"direction,omitempty"`

	// securityGroupID is the ID of the Security Group for the Security Group Rule.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=^[-0-9a-z_]+$
	SecurityGroupID string `json:"securityGroupID,omitempty"`

	// source defines the source of inbound traffic for the Security Group Rule.
	// Only used when direction is VPCSecurityGroupRuleDirectionInbound.
	// +optional
	Source VPCSecurityGroupRulePrototype `json:"source,omitempty,omitzero"`
}

// VPCSecurityGroupRuleRemote defines a VPC Security Group Rule's remote details.
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'any' ? (!has(self.cidrSubnetName) && !has(self.address) && !has(self.securityGroupName)) : true",message="cidrSubnetName, address, and securityGroupName are not valid for VPCSecurityGroupRuleRemoteTypeAny remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'cidr' ? (has(self.cidrSubnetName) && !has(self.address) && !has(self.securityGroupName)) : true",message="only cidrSubnetName is valid for VPCSecurityGroupRuleRemoteTypeCIDR remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'address' ? (has(self.address) && !has(self.cidrSubnetName) && !has(self.securityGroupName)) : true",message="only address is valid for VPCSecurityGroupRuleRemoteTypeAddress remoteType"
// +kubebuilder:validation:XValidation:rule="self.remoteType == 'sg' ? (has(self.securityGroupName) && !has(self.cidrSubnetName) && !has(self.address)) : true",message="only securityGroupName is valid for VPCSecurityGroupRuleRemoteTypeSG remoteType"
type VPCSecurityGroupRuleRemote struct {
	// cidrSubnetName is the name of the VPC Subnet to retrieve the CIDR from.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=^([a-z]|[a-z][-a-z0-9]*[a-z0-9])$
	CIDRSubnetName string `json:"cidrSubnetName,omitempty"`

	// address is the address to use for the remote's destination/source.
	// +optional
	// +kubebuilder:validation:MinLength=7
	// +kubebuilder:validation:MaxLength=15
	// +kubebuilder:validation:Pattern=`^(([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])\.){3}([0-9]|[1-9][0-9]|1[0-9]{2}|2[0-4][0-9]|25[0-5])$`
	Address string `json:"address,omitempty"`

	// remoteType defines the type of filter to define for the remote's destination/source.
	// +required
	RemoteType VPCSecurityGroupRuleRemoteType `json:"remoteType,omitempty"`

	// securityGroupName is the name of the VPC Security Group to use for the remote.
	// +optional
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$
	SecurityGroupName string `json:"securityGroupName,omitempty"`
}

// VPCSecurityGroupPortRange represents a range of ports, minimum to maximum.
// +kubebuilder:validation:XValidation:rule="self.maximumPort >= self.minimumPort",message="maximum port must be greater than or equal to minimum port"
type VPCSecurityGroupPortRange struct {
	// maximumPort is the inclusive upper range of ports.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	MaximumPort int64 `json:"maximumPort,omitempty"`

	// minimumPort is the inclusive lower range of ports.
	// +required
	// +kubebuilder:validation:Minimum=1
	// +kubebuilder:validation:Maximum=65535
	MinimumPort int64 `json:"minimumPort,omitempty"`
}

// VPCSecurityGroupRulePrototype defines a VPC Security Group Rule's traffic specifics.
// +kubebuilder:validation:XValidation:rule="self.protocol != 'icmp' ? (!has(self.icmpCode) && !has(self.icmpType)) : true",message="icmpCode and icmpType are only supported for VPCSecurityGroupRuleProtocolIcmp protocol"
// +kubebuilder:validation:XValidation:rule="self.protocol == 'icmp' ? !has(self.portRange) : true",message="portRange is not valid for VPCSecurityGroupRuleProtocolIcmp protocol"
// +kubebuilder:validation:XValidation:rule="(self.protocol != 'tcp' && self.protocol != 'udp') ? !has(self.portRange) : true",message="portRange is not valid for protocol"
// +kubebuilder:validation:XValidation:rule="self.protocol == 'icmp_tcp_udp' ? !has(self.portRange) : true",message="portRange is not valid for VPCSecurityGroupRuleProtocolIcmpTCPUDP protocol"
type VPCSecurityGroupRulePrototype struct {
	// icmpCode is the ICMP code for the Rule.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=255
	ICMPCode *int64 `json:"icmpCode,omitempty"`

	// icmpType is the ICMP type for the Rule.
	// +optional
	// +kubebuilder:validation:Minimum=0
	// +kubebuilder:validation:Maximum=254
	ICMPType *int64 `json:"icmpType,omitempty"`

	// portRange is a range of ports allowed for the Rule's remote.
	// If specified, both minimumPort and maximumPort must be specified.
	// If unspecified, traffic on all destination ports is allowed.
	// +optional
	PortRange VPCSecurityGroupPortRange `json:"portRange,omitempty,omitzero"`

	// protocol defines the traffic protocol used for the Security Group Rule.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=128
	Protocol VPCSecurityGroupRuleProtocol `json:"protocol,omitempty"`

	// remotes is a set of VPCSecurityGroupRuleRemote's that define the traffic allowed.
	// +optional
	// +kubebuilder:validation:MaxItems=25
	// +listType=atomic
	Remotes []VPCSecurityGroupRuleRemote `json:"remotes,omitempty"`
}

// VPCSecurityGroupStatus tracks the observed state of an individual VPC security group.
type VPCSecurityGroupStatus struct {
	// id is the unique cloud identifier generated by IBM Cloud for this Security Group.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=^[-0-9a-z_]+$
	ID string `json:"id,omitempty"`

	// name is the name for this security group.
	// The name is unique across all security groups for the VPC.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=63
	// +kubebuilder:validation:Pattern=^-?([a-z]|[a-z][-a-z0-9]*[a-z0-9]|[0-9][-a-z0-9]*([a-z]|[-a-z][-a-z0-9]*[a-z0-9]))$
	Name string `json:"name,omitempty"`

	// rules tracks the synchronized IDs of the rules belonging to this security group.
	// Tracking rule IDs ensures we can cleanly reconcile, update, or remove rules later.
	// +optional
	// +kubebuilder:validation:MaxItems=250
	// +listType=map
	// +listMapKey=id
	Rules []VPCSecurityGroupRuleStatus `json:"rules,omitempty"`
}

// VPCSecurityGroupRuleStatus tracks individual security group rule identifiers returned by the API.
type VPCSecurityGroupRuleStatus struct {
	// id is the unique string identifier generated by IBM Cloud for this specific rule.
	// +required
	// +kubebuilder:validation:MinLength=1
	// +kubebuilder:validation:MaxLength=64
	// +kubebuilder:validation:Pattern=^[-0-9a-z_]+$
	ID string `json:"id,omitempty"`
}

// GetConditions returns the observations of the operational state of the IBMPowerVSCluster resource.
func (r *IBMPowerVSCluster) GetConditions() []metav1.Condition {
	return r.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (r *IBMPowerVSCluster) SetConditions(conditions []metav1.Condition) {
	r.Status.Conditions = conditions
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (r *IBMPowerVSCluster) GetV1Beta1Conditions() clusterv1.Conditions {
	if r.Status.Deprecated == nil || r.Status.Deprecated.V1Beta2 == nil {
		return nil
	}
	return r.Status.Deprecated.V1Beta2.Conditions
}

// SetV1Beta1Conditions sets conditions for an API object.
func (r *IBMPowerVSCluster) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if r.Status.Deprecated == nil {
		r.Status.Deprecated = &IBMPowerVSClusterDeprecatedStatus{}
	}
	if r.Status.Deprecated.V1Beta2 == nil {
		r.Status.Deprecated.V1Beta2 = &IBMPowerVSClusterV1Beta2DeprecatedStatus{}
	}
	r.Status.Deprecated.V1Beta2.Conditions = conditions
}
