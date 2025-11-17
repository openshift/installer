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

	capoerrors "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/errors"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const (
	// ClusterFinalizer allows ReconcileOpenStackCluster to clean up OpenStack resources associated with OpenStackCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "openstackcluster.infrastructure.cluster.x-k8s.io"
)

// OpenStackClusterSpec defines the desired state of OpenStackCluster.
// +kubebuilder:validation:XValidation:rule="has(self.disableExternalNetwork) && self.disableExternalNetwork ? !has(self.bastion) || !has(self.bastion.floatingIP) : true",message="bastion floating IP cannot be set when disableExternalNetwork is true"
// +kubebuilder:validation:XValidation:rule="has(self.disableExternalNetwork) && self.disableExternalNetwork ? has(self.disableAPIServerFloatingIP) && self.disableAPIServerFloatingIP : true",message="disableAPIServerFloatingIP cannot be false when disableExternalNetwork is true"
type OpenStackClusterSpec struct {
	// ManagedSubnets describe OpenStack Subnets to be created. Cluster actuator will create a network,
	// subnets with the defined CIDR, and a router connected to these subnets. Currently only one IPv4
	// subnet is supported. If you leave this empty, no network will be created.
	// +kubebuilder:validation:MaxItems=1
	// +listType=atomic
	// +optional
	ManagedSubnets []SubnetSpec `json:"managedSubnets,omitempty"`

	// Router specifies an existing router to be used if ManagedSubnets are
	// specified. If specified, no new router will be created.
	// +optional
	Router *RouterParam `json:"router,omitempty"`

	// Network specifies an existing network to use if no ManagedSubnets
	// are specified.
	// +optional
	Network *NetworkParam `json:"network,omitempty"`

	// Subnets specifies existing subnets to use if not ManagedSubnets are
	// specified. All subnets must be in the network specified by Network.
	// There can be zero, one, or two subnets. If no subnets are specified,
	// all subnets in Network will be used. If 2 subnets are specified, one
	// must be IPv4 and the other IPv6.
	// +kubebuilder:validation:MaxItems=2
	// +listType=atomic
	// +optional
	Subnets []SubnetParam `json:"subnets,omitempty"`

	// NetworkMTU sets the maximum transmission unit (MTU) value to address fragmentation for the private network ID.
	// This value will be used only if the Cluster actuator creates the network.
	// If left empty, the network will have the default MTU defined in Openstack network service.
	// To use this field, the Openstack installation requires the net-mtu neutron API extension.
	// +optional
	NetworkMTU optional.Int `json:"networkMTU,omitempty"`

	// ExternalRouterIPs is an array of externalIPs on the respective subnets.
	// This is necessary if the router needs a fixed ip in a specific subnet.
	// +listType=atomic
	// +optional
	ExternalRouterIPs []ExternalRouterIPParam `json:"externalRouterIPs,omitempty"`

	// ExternalNetwork is the OpenStack Network to be used to get public internet to the VMs.
	// This option is ignored if DisableExternalNetwork is set to true.
	//
	// If ExternalNetwork is defined it must refer to exactly one external network.
	//
	// If ExternalNetwork is not defined or is empty the controller will use any
	// existing external network as long as there is only one. It is an
	// error if ExternalNetwork is not defined and there are multiple
	// external networks unless DisableExternalNetwork is also set.
	//
	// If ExternalNetwork is not defined and there are no external networks
	// the controller will proceed as though DisableExternalNetwork was set.
	// +optional
	ExternalNetwork *NetworkParam `json:"externalNetwork,omitempty"`

	// DisableExternalNetwork specifies whether or not to attempt to connect the cluster
	// to an external network. This allows for the creation of clusters when connecting
	// to an external network is not possible or desirable, e.g. if using a provider network.
	// +optional
	DisableExternalNetwork optional.Bool `json:"disableExternalNetwork,omitempty"`

	// APIServerLoadBalancer configures the optional LoadBalancer for the APIServer.
	// If not specified, no load balancer will be created for the API server.
	// +optional
	APIServerLoadBalancer *APIServerLoadBalancer `json:"apiServerLoadBalancer,omitempty"`

	// DisableAPIServerFloatingIP determines whether or not to attempt to attach a floating
	// IP to the API server. This allows for the creation of clusters when attaching a floating
	// IP to the API server (and hence, in many cases, exposing the API server to the internet)
	// is not possible or desirable, e.g. if using a shared VLAN for communication between
	// management and workload clusters or when the management cluster is inside the
	// project network.
	// This option requires that the API server use a VIP on the cluster network so that the
	// underlying machines can change without changing ControlPlaneEndpoint.Host.
	// When using a managed load balancer, this VIP will be managed automatically.
	// If not using a managed load balancer, cluster configuration will fail without additional
	// configuration to manage the VIP on the control plane machines, which falls outside of
	// the scope of this controller.
	// +optional
	DisableAPIServerFloatingIP optional.Bool `json:"disableAPIServerFloatingIP,omitempty"`

	// APIServerFloatingIP is the floatingIP which will be associated with the API server.
	// The floatingIP will be created if it does not already exist.
	// If not specified, a new floatingIP is allocated.
	// This field is not used if DisableAPIServerFloatingIP is set to true.
	// +optional
	APIServerFloatingIP optional.String `json:"apiServerFloatingIP,omitempty"`

	// APIServerFixedIP is the fixed IP which will be associated with the API server.
	// In the case where the API server has a floating IP but not a managed load balancer,
	// this field is not used.
	// If a managed load balancer is used and this field is not specified, a fixed IP will
	// be dynamically allocated for the load balancer.
	// If a managed load balancer is not used AND the API server floating IP is disabled,
	// this field MUST be specified and should correspond to a pre-allocated port that
	// holds the fixed IP to be used as a VIP.
	// +optional
	APIServerFixedIP optional.String `json:"apiServerFixedIP,omitempty"`

	// APIServerPort is the port on which the listener on the APIServer
	// will be created. If specified, it must be an integer between 0 and 65535.
	// +optional
	APIServerPort optional.UInt16 `json:"apiServerPort,omitempty"`

	// ManagedSecurityGroups determines whether OpenStack security groups for the cluster
	// will be managed by the OpenStack provider or whether pre-existing security groups will
	// be specified as part of the configuration.
	// By default, the managed security groups have rules that allow the Kubelet, etcd, and the
	// Kubernetes API server to function correctly.
	// It's possible to add additional rules to the managed security groups.
	// When defined to an empty struct, the managed security groups will be created with the default rules.
	// +optional
	ManagedSecurityGroups *ManagedSecurityGroups `json:"managedSecurityGroups,omitempty"`

	// DisablePortSecurity disables the port security of the network created for the
	// Kubernetes cluster, which also disables SecurityGroups
	// +optional
	DisablePortSecurity optional.Bool `json:"disablePortSecurity,omitempty"`

	// Tags to set on all resources in cluster which support tags
	// +listType=set
	// +optional
	Tags []string `json:"tags,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// It is normally populated automatically by the OpenStackCluster
	// controller during cluster provisioning. If it is set on creation the
	// control plane endpoint will use the values set here in preference to
	// values set elsewhere.
	// ControlPlaneEndpoint cannot be modified after ControlPlaneEndpoint.Host has been set.
	// +optional
	ControlPlaneEndpoint *clusterv1.APIEndpoint `json:"controlPlaneEndpoint,omitempty"`

	// ControlPlaneAvailabilityZones is the set of availability zones which
	// control plane machines may be deployed to.
	// +listType=set
	// +optional
	ControlPlaneAvailabilityZones []string `json:"controlPlaneAvailabilityZones,omitempty"`

	// ControlPlaneOmitAvailabilityZone causes availability zone to be
	// omitted when creating control plane nodes, allowing the Nova
	// scheduler to make a decision on which availability zone to use based
	// on other scheduling constraints
	// +optional
	ControlPlaneOmitAvailabilityZone optional.Bool `json:"controlPlaneOmitAvailabilityZone,omitempty"`

	// Bastion is the OpenStack instance to login the nodes
	//
	// As a rolling update is not ideal during a bastion host session, we
	// prevent changes to a running bastion configuration. To make changes, it's required
	// to first set `enabled: false` which will remove the bastion and then changes can be made.
	//+optional
	Bastion *Bastion `json:"bastion,omitempty"`

	// IdentityRef is a reference to a secret holding OpenStack credentials
	// to be used when reconciling this cluster. It is also to reconcile
	// machines unless overridden in the machine spec.
	// +kubebuilder:validation:Required
	IdentityRef OpenStackIdentityReference `json:"identityRef"`
}

// OpenStackClusterStatus defines the observed state of OpenStackCluster.
type OpenStackClusterStatus struct {
	// Ready is true when the cluster infrastructure is ready.
	// +kubebuilder:default=false
	Ready bool `json:"ready"`

	// Network contains information about the created OpenStack Network.
	// +optional
	Network *NetworkStatusWithSubnets `json:"network,omitempty"`

	// ExternalNetwork contains information about the external network used for default ingress and egress traffic.
	// +optional
	ExternalNetwork *NetworkStatus `json:"externalNetwork,omitempty"`

	// Router describes the default cluster router
	// +optional
	Router *Router `json:"router,omitempty"`

	// APIServerLoadBalancer describes the api server load balancer if one exists
	// +optional
	APIServerLoadBalancer *LoadBalancer `json:"apiServerLoadBalancer,omitempty"`

	// FailureDomains represent OpenStack availability zones
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// ControlPlaneSecurityGroup contains the information about the
	// OpenStack Security Group that needs to be applied to control plane
	// nodes.
	// +optional
	ControlPlaneSecurityGroup *SecurityGroupStatus `json:"controlPlaneSecurityGroup,omitempty"`

	// WorkerSecurityGroup contains the information about the OpenStack
	// Security Group that needs to be applied to worker nodes.
	// +optional
	WorkerSecurityGroup *SecurityGroupStatus `json:"workerSecurityGroup,omitempty"`

	// BastionSecurityGroup contains the information about the OpenStack
	// Security Group that needs to be applied to worker nodes.
	// +optional
	BastionSecurityGroup *SecurityGroupStatus `json:"bastionSecurityGroup,omitempty"`

	// Bastion contains the information about the deployed bastion host
	// +optional
	Bastion *BastionStatus `json:"bastion,omitempty"`

	// FailureReason will be set in the event that there is a terminal problem
	// reconciling the OpenStackCluster and will contain a succinct value suitable
	// for machine interpretation.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the OpenStackCluster's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of
	// OpenStackClusters can be added as events to the OpenStackCluster object
	// and/or logged in the controller's output.
	// +optional
	FailureReason *capoerrors.DeprecatedCAPIClusterStatusError `json:"failureReason,omitempty"`

	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the OpenStackCluster and will contain a more verbose string suitable
	// for logging and human consumption.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the OpenStackCluster's spec or the configuration of
	// the controller, and that manual intervention is required. Examples
	// of terminal errors would be invalid combinations of settings in the
	// spec, values that are unsupported by the controller, or the
	// responsible controller itself being critically misconfigured.
	//
	// Any transient errors that occur during the reconciliation of
	// OpenStackClusters can be added as events to the OpenStackCluster object
	// and/or logged in the controller's output.
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`
}

// +genclient
// +kubebuilder:object:root=true
// +kubebuilder:resource:path=openstackclusters,scope=Namespaced,categories=cluster-api,shortName=osc
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for OpenStack instances"
// +kubebuilder:printcolumn:name="Network",type="string",JSONPath=".status.network.id",description="Network the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.floatingIP",description="Bastion address for breakglass access"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStackCluster"

// OpenStackCluster is the Schema for the openstackclusters API.
type OpenStackCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackClusterSpec   `json:"spec,omitempty"`
	Status OpenStackClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackClusterList contains a list of OpenStackCluster.
type OpenStackClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackCluster `json:"items"`
}

// ManagedSecurityGroups defines the desired state of security groups and rules for the cluster.
type ManagedSecurityGroups struct {
	// allNodesSecurityGroupRules defines the rules that should be applied to all nodes.
	// +patchMergeKey=name
	// +patchStrategy=merge
	// +listType=map
	// +listMapKey=name
	// +optional
	AllNodesSecurityGroupRules []SecurityGroupRuleSpec `json:"allNodesSecurityGroupRules,omitempty" patchStrategy:"merge" patchMergeKey:"name"`

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

	// AllowAllInClusterTraffic allows all ingress and egress traffic between cluster nodes when set to true.
	// +kubebuilder:default=false
	// +kubebuilder:validation:Required
	AllowAllInClusterTraffic bool `json:"allowAllInClusterTraffic"`
}

var _ IdentityRefProvider = &OpenStackCluster{}

// GetIdentifyRef returns the cluster's namespace and IdentityRef.
func (c *OpenStackCluster) GetIdentityRef() (*string, *OpenStackIdentityReference) {
	return &c.Namespace, &c.Spec.IdentityRef
}

func init() {
	objectTypes = append(objectTypes, &OpenStackCluster{}, &OpenStackClusterList{})
}
