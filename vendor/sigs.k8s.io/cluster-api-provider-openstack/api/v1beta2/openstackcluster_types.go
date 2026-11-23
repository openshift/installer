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

package v1beta2

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const (
	// ClusterFinalizer allows ReconcileOpenStackCluster to clean up OpenStack resources associated with OpenStackCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "openstackcluster.infrastructure.cluster.x-k8s.io"
)

// OpenStackClusterSpec defines the desired state of OpenStackCluster.
// +kubebuilder:validation:XValidation:rule="has(self.enableExternalNetwork) && !self.enableExternalNetwork ? !has(self.bastion) || !has(self.bastion.floatingIP) : true",message="bastion floating IP cannot be set when enableExternalNetwork is false"
// +kubebuilder:validation:XValidation:rule="has(self.enableExternalNetwork) && !self.enableExternalNetwork ? has(self.apiServer) && has(self.apiServer.enableFloatingIP) && !self.apiServer.enableFloatingIP : true",message="apiServer.enableFloatingIP cannot be true when enableExternalNetwork is false"
type OpenStackClusterSpec struct {
	// managedSubnets describe OpenStack Subnets to be created. Cluster actuator will create a network,
	// subnets with the defined CIDR, and a router connected to these subnets. Currently only one IPv4
	// subnet is supported. If you leave this empty, no network will be created.
	// +kubebuilder:validation:MaxItems=1
	// +listType=atomic
	// +optional
	ManagedSubnets []SubnetSpec `json:"managedSubnets,omitempty"`

	// subnets specifies existing subnets to use if not ManagedSubnets are
	// specified. All subnets must be in the network specified by Network.
	// If no subnets are specified, all subnets in Network will be used.
	// Multiple subnets of the same IP version are supported when primarySubnet
	// is also set to identify which subnet should be used for services like
	// load balancer VIP allocation.
	// +listType=atomic
	// +optional
	Subnets []SubnetParam `json:"subnets,omitempty"`

	// primarySubnet identifies the primary subnet for the cluster when multiple
	// subnets are specified in Subnets. It is used to determine the subnet for
	// load balancer VIP allocation and node member registration.
	// If not specified and multiple subnets exist, the first subnet in the
	// resolved Subnets list is used.
	// +optional
	PrimarySubnet *SubnetParam `json:"primarySubnet,omitempty"`

	// managedRouter specifies attributes of the router. The values are used only
	// if the Cluster actuator creates the router.
	// +kubebuilder:validation:XValidation:rule="has(self.externalIPs)",message="managedRouter must not be empty if set"
	// +optional
	ManagedRouter *ManagedRouter `json:"managedRouter,omitempty"`

	// router specifies an existing router to be used if ManagedSubnets are
	// specified. If specified, no new router will be created.
	// +optional
	Router *RouterParam `json:"router,omitempty"`

	// managedNetwork specifies attributes of the network. The values are used only
	// if the Cluster actuator creates the network.
	// +kubebuilder:validation:XValidation:rule="self == null || has(self.mtu) || has(self.enablePortSecurity)",message="managedNetwork must not be empty if set"
	// +optional
	ManagedNetwork *ManagedNetwork `json:"managedNetwork,omitempty"`

	// network specifies an existing network to use if no ManagedSubnets
	// are specified.
	// +optional
	Network *NetworkParam `json:"network,omitempty"`

	// externalNetwork is the OpenStack Network to be used to get public internet to the VMs.
	// This option is ignored if EnableExternalNetwork is set to false.
	//
	// If ExternalNetwork is defined it must refer to exactly one external network.
	//
	// If ExternalNetwork is not defined or is empty the controller will use any
	// existing external network as long as there is only one. It is an
	// error if ExternalNetwork is not defined and there are multiple
	// external networks unless EnableExternalNetwork is also set to false.
	//
	// If ExternalNetwork is not defined and there are no external networks
	// the controller will proceed as though EnableExternalNetwork was set to false.
	// +optional
	ExternalNetwork *NetworkParam `json:"externalNetwork,omitempty"`

	// enableExternalNetwork specifies whether to connect the cluster to an external network.
	// Set this to false when connecting to an external network is not possible or desirable,
	// e.g. if using a provider network.
	// +optional
	EnableExternalNetwork optional.Bool `json:"enableExternalNetwork,omitempty"`

	// apiServer configures the API server endpoint and its associated
	// load balancer and floating IP.
	// +optional
	APIServer *APIServer `json:"apiServer,omitempty"`

	// managedSecurityGroups determines whether OpenStack security groups for the cluster
	// will be managed by the OpenStack provider or whether pre-existing security groups will
	// be specified as part of the configuration.
	// By default, the managed security groups have rules that allow the Kubelet, etcd, and the
	// Kubernetes API server to function correctly.
	// It's possible to add additional rules to the managed security groups.
	// When defined to an empty struct, the managed security groups will be created with the default rules.
	// +optional
	ManagedSecurityGroups *ManagedSecurityGroups `json:"managedSecurityGroups,omitempty"`

	// tags to set on all resources in cluster which support tags
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// controlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// It is normally populated automatically by the OpenStackCluster
	// controller during cluster provisioning. If it is set on creation the
	// control plane endpoint will use the values set here in preference to
	// values set elsewhere.
	// ControlPlaneEndpoint cannot be modified after ControlPlaneEndpoint.Host has been set.
	// +optional
	ControlPlaneEndpoint *clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// controlPlaneAvailabilityZones is the set of availability zones which
	// control plane machines may be deployed to.
	// +listType=set
	// +optional
	ControlPlaneAvailabilityZones []string `json:"controlPlaneAvailabilityZones,omitempty"`

	// controlPlaneOmitAvailabilityZone causes availability zone to be
	// omitted when creating control plane nodes, allowing the Nova
	// scheduler to make a decision on which availability zone to use based
	// on other scheduling constraints
	// +optional
	ControlPlaneOmitAvailabilityZone optional.Bool `json:"controlPlaneOmitAvailabilityZone,omitempty"`

	// bastion is the OpenStack instance to login the nodes
	//
	// As a rolling update is not ideal during a bastion host session, we
	// prevent changes to a running bastion configuration. To make changes, it's required
	// to first set `enabled: false` which will remove the bastion and then changes can be made.
	// +optional
	Bastion *Bastion `json:"bastion,omitempty"`

	// identityRef is a reference to a secret holding OpenStack credentials
	// to be used when reconciling this cluster. It is also to reconcile
	// machines unless overridden in the machine spec.
	// +required
	IdentityRef OpenStackIdentityReference `json:"identityRef,omitzero"`
}

type APIServer struct {
	// port is the port on which the API server listener will be created.
	// If specified, it must be an integer between 0 and 65535.
	// +optional
	Port optional.UInt16 `json:"port,omitempty"`

	// fixedIP is the fixed IP which will be associated with the API server.
	// In the case where the API server has a floating IP but not a managed
	// load balancer, this field is not used.
	// If a managed load balancer is used and this field is not specified, a
	// fixed IP will be dynamically allocated for the load balancer.
	// If a managed load balancer is not used AND the floating IP is disabled,
	// this field MUST be specified and should correspond to a pre-allocated
	// port that holds the fixed IP to be used as a VIP.
	// +optional
	FixedIP optional.String `json:"fixedIP,omitempty"`

	// floatingIP is the floating IP which will be associated with the API server.
	// The floating IP will be created if it does not already exist.
	// If not specified, a new floating IP is allocated.
	// This field is not used if EnableFloatingIP is set to false.
	// +optional
	FloatingIP optional.String `json:"floatingIP,omitempty"`

	// enableFloatingIP determines whether to attach a floating IP to the API server.
	// +optional
	EnableFloatingIP optional.Bool `json:"enableFloatingIP,omitempty"`

	// managedLoadBalancer configures the optional LoadBalancer for the API server.
	// If not specified, no load balancer will be created.
	// +optional
	ManagedLoadBalancer *APIServerLoadBalancer `json:"managedLoadBalancer,omitempty"`
}

// ClusterInitialization represents the initialization status of the cluster.
type ClusterInitialization struct {
	// provisioned is set to true when the initial provisioning of the cluster infrastructure is completed.
	// The value of this field is never updated after provisioning is completed.
	// +optional
	Provisioned bool `json:"provisioned,omitempty"`
}

// OpenStackClusterStatus defines the observed state of OpenStackCluster.
type OpenStackClusterStatus struct {
	// conditions defines current service state of the OpenStackCluster.
	// This field surfaces into Cluster's status.conditions[InfrastructureReady] condition.
	// The Ready condition must surface issues during the entire lifecycle of the OpenStackCluster
	// (both during initial provisioning and after the initial provisioning is completed).
	// +optional
	// +listType=map
	// +listMapKey=type
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization contains information about the initialization status of the cluster.
	// +optional
	Initialization *ClusterInitialization `json:"initialization,omitempty"`

	// network contains information about the created OpenStack Network.
	// +optional
	Network *NetworkStatusWithSubnets `json:"network,omitempty"`

	// externalNetwork contains information about the external network used for default ingress and egress traffic.
	// +optional
	ExternalNetwork *NetworkStatus `json:"externalNetwork,omitempty"`

	// router describes the default cluster router
	// +optional
	Router *Router `json:"router,omitempty"`

	// apiServerManagedLoadBalancer describes the api server load balancer if one exists
	// +optional
	APIServerManagedLoadBalancer *LoadBalancer `json:"apiServerManagedLoadBalancer,omitempty"`

	// failureDomains represent OpenStack availability zones
	// +listType=map
	// +listMapKey=name
	// +optional
	FailureDomains []clusterv1.FailureDomain `json:"failureDomains,omitempty"`

	// controlPlaneSecurityGroup contains the information about the
	// OpenStack Security Group that needs to be applied to control plane
	// nodes.
	// +optional
	ControlPlaneSecurityGroup *SecurityGroupStatus `json:"controlPlaneSecurityGroup,omitempty"`

	// workerSecurityGroup contains the information about the OpenStack
	// Security Group that needs to be applied to worker nodes.
	// +optional
	WorkerSecurityGroup *SecurityGroupStatus `json:"workerSecurityGroup,omitempty"`

	// bastionSecurityGroup contains the information about the OpenStack
	// Security Group that needs to be applied to worker nodes.
	// +optional
	BastionSecurityGroup *SecurityGroupStatus `json:"bastionSecurityGroup,omitempty"`

	// bastion contains the information about the deployed bastion host
	// +optional
	Bastion *BastionStatus `json:"bastion,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:storageversion
// +kubebuilder:resource:path=openstackclusters,scope=Namespaced,categories=cluster-api,shortName=osc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackCluster belongs"
// +kubebuilder:printcolumn:name="Network",type="string",JSONPath=".status.network.id",description="Network the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.floatingIP",description="Bastion address for breakglass access"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStackCluster"

// OpenStackCluster is the Schema for the openstackclusters API.
type OpenStackCluster struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object metadata.
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of the OpenStackCluster.
	// +optional
	Spec OpenStackClusterSpec `json:"spec,omitempty"`
	// status is the observed state of the OpenStackCluster.
	// +optional
	Status OpenStackClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackClusterList contains a list of OpenStackCluster.
type OpenStackClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	// +required
	Items []OpenStackCluster `json:"items"`
}

// ManagedRouter specifies attributes of the router.
type ManagedRouter struct {
	// externalIPs is a list of external IPs to assign to the router.
	// This is necessary if the router needs a fixed ip in a specific subnet.
	// Each entry specifies a fixed IP and the subnet it should be allocated from.
	// +kubebuilder:validation:MinItems=1
	// +optional
	// +listType=atomic
	ExternalIPs []ExternalRouterIPParam `json:"externalIPs,omitempty"`
}

// ManagedNetwork specifies attributes of the network.
type ManagedNetwork struct {
	// mtu sets the maximum transmission unit (MTU) value to address fragmentation for the private network ID.
	// This value will be used only if the Cluster actuator creates the network.
	// If left empty, the network will have the default MTU defined in Openstack network service.
	// To use this field, the Openstack installation requires the net-mtu neutron API extension.
	// +optional
	MTU *int32 `json:"mtu,omitempty"`

	// enablePortSecurity enables port security for the network created for the
	// Kubernetes cluster, which also enables SecurityGroups.
	// If left empty, the network will have port security setting enabled.
	// +optional
	EnablePortSecurity optional.Bool `json:"enablePortSecurity,omitempty"`
}

// ManagedSecurityGroups defines the desired state of security groups and rules for the cluster.
type ManagedSecurityGroups struct {
	// clusterNodesSecurityGroupRules defines the rules that should be applied to all cluster nodes, excluding the bastion host.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	ClusterNodesSecurityGroupRules []SecurityGroupRuleSpec `json:"clusterNodesSecurityGroupRules,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// controlPlaneNodesSecurityGroupRules defines the rules that should be applied to control plane nodes.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	ControlPlaneNodesSecurityGroupRules []SecurityGroupRuleSpec `json:"controlPlaneNodesSecurityGroupRules,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// workerNodesSecurityGroupRules defines the rules that should be applied to worker nodes.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	WorkerNodesSecurityGroupRules []SecurityGroupRuleSpec `json:"workerNodesSecurityGroupRules,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

	// allowAllInClusterTraffic allows all ingress and egress traffic between cluster nodes when set to true.
	// +kubebuilder:default=false
	// +optional
	AllowAllInClusterTraffic bool `json:"allowAllInClusterTraffic,omitempty"`
}

var _ IdentityRefProvider = &OpenStackCluster{}

// GetConditions returns the observations of the operational state of the OpenStackCluster resource.
func (c *OpenStackCluster) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets the underlying service state of the OpenStackCluster to the predescribed clusterv1.Conditions.
func (c *OpenStackCluster) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

// GetIdentifyRef returns the cluster's namespace and IdentityRef.
func (c *OpenStackCluster) GetIdentityRef() (*string, *OpenStackIdentityReference) {
	return &c.Namespace, &c.Spec.IdentityRef
}

func init() {
	objectTypes = append(objectTypes, &OpenStackCluster{}, &OpenStackClusterList{})
}

func (a *APIServer) GetManagedLoadBalancer() *APIServerLoadBalancer {
	if a == nil {
		return nil
	}
	return a.ManagedLoadBalancer
}

func (a *APIServer) GetEnableFloatingIP() *bool {
	if a == nil {
		return nil
	}
	return a.EnableFloatingIP
}

func (a *APIServer) GetFloatingIP() *string {
	if a == nil {
		return nil
	}
	return a.FloatingIP
}

func (a *APIServer) GetFixedIP() *string {
	if a == nil {
		return nil
	}
	return a.FixedIP
}

func (a *APIServer) GetPort() *uint16 {
	if a == nil {
		return nil
	}
	return a.Port
}
