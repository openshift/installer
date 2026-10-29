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
	"fmt"
	"net"
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ClusterFinalizer allows ReconcileVSphereCluster to clean up vSphere
	// resources associated with VSphereCluster before removing it from the
	// API server.
	ClusterFinalizer = "vspherecluster.vmware.infrastructure.cluster.x-k8s.io"

	// ProviderServiceAccountFinalizer allows ServiceAccountReconciler to clean up service accounts
	// resources associated with VSphereCluster from the SERVICE_ACCOUNTS_CM (service accounts ConfigMap).
	//
	// Deprecated: ProviderServiceAccountFinalizer will be removed in a future release.
	ProviderServiceAccountFinalizer = "providerserviceaccount.vmware.infrastructure.cluster.x-k8s.io"
)

// VSphereCluster's Ready condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterReadyCondition is true if the VSphereCluster's deletionTimestamp is not set, VSphereCluster's
	// ResourcePolicyReady, NetworkReady, LoadBalancerReady, FailureDomainsReady, ProviderServiceAccountsReady and ServiceDiscoveryReady conditions are true.
	VSphereClusterReadyCondition = clusterv1.ReadyCondition

	// VSphereClusterReadyReason surfaces when the VSphereCluster readiness criteria is met.
	VSphereClusterReadyReason = clusterv1.ReadyReason

	// VSphereClusterNotReadyReason surfaces when the VSphereCluster readiness criteria is not met.
	VSphereClusterNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterReadyUnknownReason surfaces when at least one VSphereCluster readiness criteria is unknown
	// and no VSphereCluster readiness criteria is not met.
	VSphereClusterReadyUnknownReason = clusterv1.ReadyUnknownReason
)

// VSphereCluster's ResourcePolicyReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterResourcePolicyReadyCondition documents the status of the ResourcePolicy for a VSphereCluster.
	VSphereClusterResourcePolicyReadyCondition = "ResourcePolicyReady"

	// VSphereClusterResourcePolicyReadyReason surfaces when the ResourcePolicy for a VSphereCluster is ready.
	VSphereClusterResourcePolicyReadyReason = clusterv1.ReadyReason

	// VSphereClusterResourcePolicyNotReadyReason surfaces when the ResourcePolicy for a VSphereCluster is not ready.
	VSphereClusterResourcePolicyNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterResourcePolicyReadyDeletingReason surfaces when the resource policy for a VSphereCluster is being deleted.
	VSphereClusterResourcePolicyReadyDeletingReason = clusterv1.DeletingReason
)

// VSphereCluster's NetworkReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterNetworkReadyCondition documents the status of the network for a VSphereCluster.
	VSphereClusterNetworkReadyCondition = "NetworkReady"

	// VSphereClusterNetworkReadyReason surfaces when the network for a VSphereCluster is ready.
	VSphereClusterNetworkReadyReason = clusterv1.ReadyReason

	// VSphereClusterNetworkNotReadyReason surfaces when the network for a VSphereCluster is not ready.
	VSphereClusterNetworkNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterNetworkReadyDeletingReason surfaces when the network for a VSphereCluster is being deleted.
	VSphereClusterNetworkReadyDeletingReason = clusterv1.DeletingReason
)

// VSphereCluster's LoadBalancerReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterLoadBalancerReadyCondition documents the status of the LoadBalancer for a VSphereCluster.
	VSphereClusterLoadBalancerReadyCondition = "LoadBalancerReady"

	// VSphereClusterLoadBalancerReadyReason surfaces when the LoadBalancer for a VSphereCluster is ready.
	VSphereClusterLoadBalancerReadyReason = clusterv1.ReadyReason

	// VSphereClusterLoadBalancerNotReadyReason surfaces when the LoadBalancer for a VSphereCluster is not ready.
	VSphereClusterLoadBalancerNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterLoadBalancerWaitingForIPReason surfaces when the LoadBalancer for a VSphereCluster is waiting for an IP to be assigned.
	VSphereClusterLoadBalancerWaitingForIPReason = "WaitingForIP"

	// VSphereClusterLoadBalancerWaitingForKubeadmControlPlaneSpecReadyReason surfaces when the LoadBalancer has IP(s) but CAPV is waiting for
	// KubeadmControlPlane to have correct certSANs and observedGeneration == generation (avoids extra rollout for dual stack).
	VSphereClusterLoadBalancerWaitingForKubeadmControlPlaneSpecReadyReason = "WaitingForKubeadmControlPlaneSpecReady"

	// VSphereClusterLoadBalancerDeletingReason surfaces when the LoadBalancer for a VSphereCluster is being deleted.
	VSphereClusterLoadBalancerDeletingReason = clusterv1.DeletingReason
)

// VSphereCluster's ProviderServiceAccountsReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterProviderServiceAccountsReadyCondition documents the status of the provider service accounts for a VSphereCluster.
	VSphereClusterProviderServiceAccountsReadyCondition = "ProviderServiceAccountsReady"

	// VSphereClusterProviderServiceAccountsReadyReason surfaces when the provider service accounts for a VSphereCluster is ready.
	VSphereClusterProviderServiceAccountsReadyReason = clusterv1.ReadyReason

	// VSphereClusterProviderServiceAccountsNotReadyReason surfaces when the provider service accounts for a VSphereCluster is not ready.
	VSphereClusterProviderServiceAccountsNotReadyReason = clusterv1.NotReadyReason
)

// VSphereCluster's ServiceDiscoveryReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterServiceDiscoveryReadyCondition documents the status of the service discovery for a VSphereCluster.
	VSphereClusterServiceDiscoveryReadyCondition = "ServiceDiscoveryReady"

	// VSphereClusterServiceDiscoveryReadyReason surfaces when the service discovery for a VSphereCluster is ready.
	VSphereClusterServiceDiscoveryReadyReason = clusterv1.ReadyReason

	// VSphereClusterServiceDiscoveryNotReadyReason surfaces when the service discovery for a VSphereCluster is not ready.
	VSphereClusterServiceDiscoveryNotReadyReason = clusterv1.NotReadyReason
)

// VSphereCluster's FailureDomainsReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterFailureDomainsReadyCondition documents the status of the FailureDomains discovery and filtering for a VSphereCluster.
	VSphereClusterFailureDomainsReadyCondition = "FailureDomainsReady"

	// VSphereClusterFailureDomainsReadyReason surfaces when the FailureDomains are successfully discovered and matched.
	VSphereClusterFailureDomainsReadyReason = clusterv1.ReadyReason

	// VSphereClusterFailureDomainsNotReadyReason surfaces when no failure domains match the specified control plane label selector.
	VSphereClusterFailureDomainsNotReadyReason = clusterv1.NotReadyReason

	// VSphereClusterFailureDomainsReadyInternalErrorReason surfaces unexpected failures during FailureDomains discovery.
	VSphereClusterFailureDomainsReadyInternalErrorReason = clusterv1.InternalErrorReason

	// VSphereClusterFailureDomainsReadyDeletingReason surfaces when the cluster is being deleted.
	VSphereClusterFailureDomainsReadyDeletingReason = clusterv1.DeletingReason
)

// NSXVPC defines the configuration when the network provider is NSX-VPC.
// +kubebuilder:validation:XValidation:rule="has(self.createSubnetSet) == has(oldSelf.createSubnetSet) && self.createSubnetSet == oldSelf.createSubnetSet",message="createSubnetSet value cannot be changed after creation"
// +kubebuilder:validation:MinProperties=1
type NSXVPC struct {
	// createSubnetSet is a flag to indicate whether to create a SubnetSet or not as the primary network. If not set, the default is true.
	// +optional
	CreateSubnetSet *bool `json:"createSubnetSet,omitempty"`
}

// IsDefined returns true if the NSXVPC is defined.
func (r *NSXVPC) IsDefined() bool {
	return !reflect.DeepEqual(r, &NSXVPC{})
}

// Network defines the network configuration for the cluster with different network providers.
// +kubebuilder:validation:XValidation:rule="has(self.nsxVPC) == has(oldSelf.nsxVPC)",message="field 'nsxVPC' cannot be added or removed after creation"
// +kubebuilder:validation:MinProperties=1
type Network struct {
	// nsxVPC defines the configuration when the network provider is NSX-VPC.
	// +optional
	NSXVPC NSXVPC `json:"nsxVPC,omitempty,omitzero"`
}

// IsDefined returns true if the Network is defined.
func (r *Network) IsDefined() bool {
	return !reflect.DeepEqual(r, &Network{})
}

// VSphereClusterSpec defines the desired state of VSphereCluster.
// +kubebuilder:validation:XValidation:rule="has(self.network) == has(oldSelf.network)",message="field 'network' cannot be added or removed after creation"
// +kubebuilder:validation:MinProperties=1
type VSphereClusterSpec struct {
	// controlPlaneEndpoint represents the endpoint used to communicate with the control plane.
	// +optional
	ControlPlaneEndpoint APIEndpoint `json:"controlPlaneEndpoint,omitempty,omitzero"`

	// network defines the network configuration for the cluster with different network providers.
	// +optional
	Network Network `json:"network,omitempty,omitzero"`

	// failureDomains defines the failure domains.
	// +optional
	FailureDomains FailureDomainsSpec `json:"failureDomains,omitempty,omitzero"`
}

// FailureDomainsSpec defines the desired state of FailureDomains.
// +kubebuilder:validation:MinProperties=1
type FailureDomainsSpec struct {
	// controlPlane defines the failure domains for controlPlane.
	// +optional
	ControlPlane FailureDomainsControlPlaneSpec `json:"controlPlane,omitempty,omitzero"`
}

// FailureDomainsControlPlaneSpec defines the control plane failure domains.
// +kubebuilder:validation:MinProperties=1
type FailureDomainsControlPlaneSpec struct {
	// selector is a label selector to dynamically match supervisor's Zone CR.
	// Note: This feature requires the NamespaceScopedZones feature gate to be enabled.
	// If a selector is provided while the feature gate is disabled, cluster
	// reconciliation will fail and requeue.
	// +optional
	Selector *metav1.LabelSelector `json:"selector,omitempty"`
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

// IsZero returns true if both host and port are zero values.
func (v APIEndpoint) IsZero() bool {
	return v.Host == "" && v.Port == 0
}

// String returns a formatted version HOST:PORT of this APIEndpoint.
func (v APIEndpoint) String() string {
	return net.JoinHostPort(v.Host, fmt.Sprintf("%d", v.Port))
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterStatus struct {
	// conditions represents the observations of a VSphereCluster's current state.
	// Known condition types are Ready, ResourcePolicyReady, NetworkReady, LoadBalancerReady,
	// FailureDomainsReady, ProviderServiceAccountsReady, ServiceDiscoveryReady and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`

	// initialization provides observations of the VSphereCluster initialization process.
	// NOTE: Fields in this struct are part of the Cluster API contract and are used to orchestrate initial Cluster provisioning.
	// +optional
	Initialization VSphereClusterInitializationStatus `json:"initialization,omitempty,omitzero"`

	// failureDomains is a list of failure domain objects synced from the infrastructure provider.
	// +optional
	// +listType=map
	// +listMapKey=name
	// +kubebuilder:validation:MinItems=1
	// +kubebuilder:validation:MaxItems=100
	FailureDomains []clusterv1.FailureDomain `json:"failureDomains,omitempty"`

	// deprecated groups all the status fields that are deprecated and will be removed when all the nested field are removed.
	// +optional
	Deprecated *VSphereClusterDeprecatedStatus `json:"deprecated,omitempty"`
}

// VSphereClusterInitializationStatus provides observations of the VSphereCluster initialization process.
// +kubebuilder:validation:MinProperties=1
type VSphereClusterInitializationStatus struct {
	// provisioned is true when the infrastructure provider reports that the Cluster's infrastructure is fully provisioned.
	// NOTE: this field is part of the Cluster API contract, and it is used to orchestrate initial Cluster provisioning.
	// +optional
	Provisioned *bool `json:"provisioned,omitempty"`
}

// VSphereClusterDeprecatedStatus groups all the status fields that are deprecated and will be removed in a future version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterDeprecatedStatus struct {
	// v1beta1 groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
	// +optional
	V1Beta1 *VSphereClusterV1Beta1DeprecatedStatus `json:"v1beta1,omitempty"`
}

// VSphereClusterV1Beta1DeprecatedStatus groups all the status fields that are deprecated and will be removed when support for v1beta1 will be dropped.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterV1Beta1DeprecatedStatus struct {
	// conditions defines current service state of the VSphereCluster.
	//
	// Deprecated: This field is deprecated and is going to be removed when support for v1beta1 will be dropped. Please see https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more details.
	//
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Cluster",type="string",JSONPath=".metadata.labels['cluster\\.x-k8s\\.io/cluster-name']",description="Cluster"
// +kubebuilder:printcolumn:name="ControlPlane Endpoint",type="string",JSONPath=".spec.controlPlaneEndpoint.host",description="API Endpoint"
// +kubebuilder:printcolumn:name="Paused",type="string",JSONPath=`.status.conditions[?(@.type=="Paused")].status`,description="Reconciliation paused",priority=10
// +kubebuilder:printcolumn:name="Provisioned",type="string",JSONPath=".status.initialization.provisioned",description="VSphereCluster is provisioned"
// +kubebuilder:printcolumn:name="Age",type="date",JSONPath=".metadata.creationTimestamp",description="Time duration since creation of VSphereCluster"

// VSphereCluster is the Schema for the VSphereClusters API.
type VSphereCluster struct {
	metav1.TypeMeta `json:",inline"`
	// metadata is the standard object's metadata.
	// More info: https://git.k8s.io/community/contributors/devel/sig-architecture/api-conventions.md#metadata
	// +optional
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// spec is the desired state of VSphereCluster.
	// +optional
	Spec VSphereClusterSpec `json:"spec,omitempty,omitzero"`

	// status is the observed state of VSphereCluster.
	// +optional
	Status VSphereClusterStatus `json:"status,omitempty,omitzero"`
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster.
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

// GetV1Beta1Conditions returns the set of conditions for this object.
func (c *VSphereCluster) GetV1Beta1Conditions() clusterv1.Conditions {
	if c.Status.Deprecated == nil || c.Status.Deprecated.V1Beta1 == nil {
		return nil
	}
	return c.Status.Deprecated.V1Beta1.Conditions
}

// SetV1Beta1Conditions sets the conditions on this object.
func (c *VSphereCluster) SetV1Beta1Conditions(conditions clusterv1.Conditions) {
	if c.Status.Deprecated == nil {
		c.Status.Deprecated = &VSphereClusterDeprecatedStatus{}
	}
	if c.Status.Deprecated.V1Beta1 == nil {
		c.Status.Deprecated.V1Beta1 = &VSphereClusterV1Beta1DeprecatedStatus{}
	}
	c.Status.Deprecated.V1Beta1.Conditions = conditions
}

// GetConditions returns the set of conditions for this object.
func (c *VSphereCluster) GetConditions() []metav1.Condition {
	return c.Status.Conditions
}

// SetConditions sets conditions for an API object.
func (c *VSphereCluster) SetConditions(conditions []metav1.Condition) {
	c.Status.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &VSphereCluster{}, &VSphereClusterList{})
}
