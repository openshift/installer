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

package v1alpha5

import (
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	capierrors "sigs.k8s.io/cluster-api/errors"
)

const (
	// ClusterFinalizer allows ReconcileOpenStackCluster to clean up OpenStack resources associated with OpenStackCluster before
	// removing it from the apiserver.
	ClusterFinalizer = "openstackcluster.infrastructure.cluster.x-k8s.io"
)

// OpenStackClusterSpec defines the desired state of OpenStackCluster.
type OpenStackClusterSpec struct {
	// The name of the cloud to use from the clouds secret
	// +optional
	CloudName string `json:"cloudName"`

	// NodeCIDR is the OpenStack Subnet to be created. Cluster actuator will create a
	// network, a subnet with NodeCIDR, and a router connected to this subnet.
	// If you leave this empty, no network will be created.
	NodeCIDR string `json:"nodeCidr,omitempty"`

	// If NodeCIDR cannot be set this can be used to detect an existing network.
	Network NetworkFilter `json:"network,omitempty"`

	// If NodeCIDR cannot be set this can be used to detect an existing subnet.
	Subnet SubnetFilter `json:"subnet,omitempty"`

	// DNSNameservers is the list of nameservers for OpenStack Subnet being created.
	// Set this value when you need create a new network/subnet while the access
	// through DNS is required.
	DNSNameservers []string `json:"dnsNameservers,omitempty"`
	// ExternalRouterIPs is an array of externalIPs on the respective subnets.
	// This is necessary if the router needs a fixed ip in a specific subnet.
	ExternalRouterIPs []ExternalRouterIPParam `json:"externalRouterIPs,omitempty"`
	// ExternalNetworkID is the ID of an external OpenStack Network. This is necessary
	// to get public internet to the VMs.
	// +optional
	ExternalNetworkID string `json:"externalNetworkId,omitempty"`

	// APIServerLoadBalancer configures the optional LoadBalancer for the APIServer.
	// It must be activated by setting `enabled: true`.
	// +optional
	APIServerLoadBalancer APIServerLoadBalancer `json:"apiServerLoadBalancer,omitempty"`

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
	DisableAPIServerFloatingIP bool `json:"disableAPIServerFloatingIP"`

	// APIServerFloatingIP is the floatingIP which will be associated with the API server.
	// The floatingIP will be created if it does not already exist.
	// If not specified, a new floatingIP is allocated.
	// This field is not used if DisableAPIServerFloatingIP is set to true.
	APIServerFloatingIP string `json:"apiServerFloatingIP,omitempty"`

	// APIServerFixedIP is the fixed IP which will be associated with the API server.
	// In the case where the API server has a floating IP but not a managed load balancer,
	// this field is not used.
	// If a managed load balancer is used and this field is not specified, a fixed IP will
	// be dynamically allocated for the load balancer.
	// If a managed load balancer is not used AND the API server floating IP is disabled,
	// this field MUST be specified and should correspond to a pre-allocated port that
	// holds the fixed IP to be used as a VIP.
	APIServerFixedIP string `json:"apiServerFixedIP,omitempty"`

	// APIServerPort is the port on which the listener on the APIServer
	// will be created
	APIServerPort int `json:"apiServerPort,omitempty"`

	// ManagedSecurityGroups determines whether OpenStack security groups for the cluster
	// will be managed by the OpenStack provider or whether pre-existing security groups will
	// be specified as part of the configuration.
	// By default, the managed security groups have rules that allow the Kubelet, etcd, the
	// Kubernetes API server and the Calico CNI plugin to function correctly.
	// +optional
	ManagedSecurityGroups bool `json:"managedSecurityGroups"`

	// AllowAllInClusterTraffic is only used when managed security groups are in use.
	// If set to true, the rules for the managed security groups are configured so that all
	// ingress and egress between cluster nodes is permitted, allowing CNIs other than
	// Calico to be used.
	// +optional
	AllowAllInClusterTraffic bool `json:"allowAllInClusterTraffic"`

	// DisablePortSecurity disables the port security of the network created for the
	// Kubernetes cluster, which also disables SecurityGroups
	DisablePortSecurity bool `json:"disablePortSecurity,omitempty"`

	// Tags for all resources in cluster
	Tags []string `json:"tags,omitempty"`

	// ControlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint clusterv1.APIEndpoint `json:"controlPlaneEndpoint"`

	// ControlPlaneAvailabilityZones is the az to deploy control plane to
	ControlPlaneAvailabilityZones []string `json:"controlPlaneAvailabilityZones,omitempty"`

	// Bastion is the OpenStack instance to login the nodes
	//
	// As a rolling update is not ideal during a bastion host session, we
	// prevent changes to a running bastion configuration. Set `enabled: false` to
	// make changes.
	//+optional
	Bastion *Bastion `json:"bastion,omitempty"`

	// IdentityRef is a reference to a identity to be used when reconciling this cluster
	// +optional
	IdentityRef *OpenStackIdentityReference `json:"identityRef,omitempty"`
}

// OpenStackClusterStatus defines the observed state of OpenStackCluster.
type OpenStackClusterStatus struct {
	Ready bool `json:"ready"`

	// Network contains all information about the created OpenStack Network.
	// It includes Subnets and Router.
	Network *Network `json:"network,omitempty"`

	// External Network contains information about the created OpenStack external network.
	ExternalNetwork *Network `json:"externalNetwork,omitempty"`

	// FailureDomains represent OpenStack availability zones
	FailureDomains clusterv1.FailureDomains `json:"failureDomains,omitempty"`

	// ControlPlaneSecurityGroups contains all the information about the OpenStack
	// Security Group that needs to be applied to control plane nodes.
	// TODO: Maybe instead of two properties, we add a property to the group?
	ControlPlaneSecurityGroup *SecurityGroup `json:"controlPlaneSecurityGroup,omitempty"`

	// WorkerSecurityGroup contains all the information about the OpenStack Security
	// Group that needs to be applied to worker nodes.
	WorkerSecurityGroup *SecurityGroup `json:"workerSecurityGroup,omitempty"`

	BastionSecurityGroup *SecurityGroup `json:"bastionSecurityGroup,omitempty"`

	Bastion *Instance `json:"bastion,omitempty"`

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
	FailureReason *capierrors.ClusterStatusError `json:"failureReason,omitempty"`

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

// +kubebuilder:object:root=true
// +kubebuilder:unservedversion
// +kubebuilder:deprecatedversion:warning="The v1alpha5 version of OpenStackCluster has been deprecated and will be removed in a future release of the API. Please upgrade."
// +kubebuilder:resource:path=openstackclusters,scope=Namespaced,categories=cluster-api,shortName=osc
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels.cluster\\.x-k8s\\.io/cluster-name",description="Cluster to which this OpenStackCluster belongs"
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="Cluster infrastructure is ready for OpenStack instances"
// +kubebuilder:printcolumn:name="Network",type="string",JSONPath=".status.network.id",description="Network the cluster is using"
// +kubebuilder:printcolumn:name="Subnet",type="string",JSONPath=".status.network.subnet.id",description="Subnet the cluster is using"
// +kubebuilder:printcolumn:name="Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint",priority=1
// +kubebuilder:printcolumn:name="Bastion IP",type="string",JSONPath=".status.bastion.floatingIP",description="Bastion address for breakglass access"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of OpenStackCluster"

// OpenStackCluster is the Schema for the openstackclusters API.
//
// Deprecated: This type will be removed in one of the next releases.
type OpenStackCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   OpenStackClusterSpec   `json:"spec,omitempty"`
	Status OpenStackClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// OpenStackClusterList contains a list of OpenStackCluster.
//
// Deprecated: This type will be removed in one of the next releases.
type OpenStackClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []OpenStackCluster `json:"items"`
}

func init() {
	SchemeBuilder.Register(&OpenStackCluster{}, &OpenStackClusterList{})
}
