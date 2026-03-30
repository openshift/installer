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
	"reflect"

	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
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
	// VSphereClusterReadyV1Beta2Condition is true if the VSphereCluster's deletionTimestamp is not set, VSphereCluster's
	// ResourcePolicyReady, NetworkReady, LoadBalancerReady, ProviderServiceAccountsReady and ServiceDiscoveryReady conditions are true.
	VSphereClusterReadyV1Beta2Condition = clusterv1beta1.ReadyV1Beta2Condition

	// VSphereClusterReadyV1Beta2Reason surfaces when the VSphereCluster readiness criteria is met.
	VSphereClusterReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterNotReadyV1Beta2Reason surfaces when the VSphereCluster readiness criteria is not met.
	VSphereClusterNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterReadyUnknownV1Beta2Reason surfaces when at least one VSphereCluster readiness criteria is unknown
	// and no VSphereCluster readiness criteria is not met.
	VSphereClusterReadyUnknownV1Beta2Reason = clusterv1beta1.ReadyUnknownV1Beta2Reason
)

// VSphereCluster's ResourcePolicyReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterResourcePolicyReadyV1Beta2Condition documents the status of the ResourcePolicy for a VSphereCluster.
	VSphereClusterResourcePolicyReadyV1Beta2Condition = "ResourcePolicyReady"

	// VSphereClusterResourcePolicyReadyV1Beta2Reason surfaces when the ResourcePolicy for a VSphereCluster is ready.
	VSphereClusterResourcePolicyReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterResourcePolicyNotReadyV1Beta2Reason surfaces when the ResourcePolicy for a VSphereCluster is not ready.
	VSphereClusterResourcePolicyNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterResourcePolicyReadyDeletingV1Beta2Reason surfaces when the resource policy for a VSphereCluster is being deleted.
	VSphereClusterResourcePolicyReadyDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereCluster's NetworkReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterNetworkReadyV1Beta2Condition documents the status of the network for a VSphereCluster.
	VSphereClusterNetworkReadyV1Beta2Condition = "NetworkReady"

	// VSphereClusterNetworkReadyV1Beta2Reason surfaces when the network for a VSphereCluster is ready.
	VSphereClusterNetworkReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterNetworkNotReadyV1Beta2Reason surfaces when the network for a VSphereCluster is not ready.
	VSphereClusterNetworkNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterNetworkReadyDeletingV1Beta2Reason surfaces when the network for a VSphereCluster is being deleted.
	VSphereClusterNetworkReadyDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereCluster's LoadBalancerReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterLoadBalancerReadyV1Beta2Condition documents the status of the LoadBalancer for a VSphereCluster.
	VSphereClusterLoadBalancerReadyV1Beta2Condition = "LoadBalancerReady"

	// VSphereClusterLoadBalancerReadyV1Beta2Reason surfaces when the LoadBalancer for a VSphereCluster is ready.
	VSphereClusterLoadBalancerReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterLoadBalancerNotReadyV1Beta2Reason surfaces when the LoadBalancer for a VSphereCluster is not ready.
	VSphereClusterLoadBalancerNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason

	// VSphereClusterLoadBalancerWaitingForIPV1Beta2Reason surfaces when the LoadBalancer for a VSphereCluster is waiting for an IP to be assigned.
	VSphereClusterLoadBalancerWaitingForIPV1Beta2Reason = "WaitingForIP"

	// VSphereClusterLoadBalancerDeletingV1Beta2Reason surfaces when the LoadBalancer for a VSphereCluster is being deleted.
	VSphereClusterLoadBalancerDeletingV1Beta2Reason = clusterv1beta1.DeletingV1Beta2Reason
)

// VSphereCluster's ProviderServiceAccountsReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterProviderServiceAccountsReadyV1Beta2Condition documents the status of the provider service accounts for a VSphereCluster.
	VSphereClusterProviderServiceAccountsReadyV1Beta2Condition = "ProviderServiceAccountsReady"

	// VSphereClusterProviderServiceAccountsReadyV1Beta2Reason surfaces when the provider service accounts for a VSphereCluster is ready.
	VSphereClusterProviderServiceAccountsReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterProviderServiceAccountsNotReadyV1Beta2Reason surfaces when the provider service accounts for a VSphereCluster is not ready.
	VSphereClusterProviderServiceAccountsNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason
)

// VSphereCluster's ServiceDiscoveryReady condition and corresponding reasons that will be used in v1Beta2 API version.
const (
	// VSphereClusterServiceDiscoveryReadyV1Beta2Condition documents the status of the service discovery for a VSphereCluster.
	VSphereClusterServiceDiscoveryReadyV1Beta2Condition = "ServiceDiscoveryReady"

	// VSphereClusterServiceDiscoveryReadyV1Beta2Reason surfaces when the service discovery for a VSphereCluster is ready.
	VSphereClusterServiceDiscoveryReadyV1Beta2Reason = clusterv1beta1.ReadyV1Beta2Reason

	// VSphereClusterServiceDiscoveryNotReadyV1Beta2Reason surfaces when the service discovery for a VSphereCluster is not ready.
	VSphereClusterServiceDiscoveryNotReadyV1Beta2Reason = clusterv1beta1.NotReadyV1Beta2Reason
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
type VSphereClusterSpec struct {
	// +optional
	ControlPlaneEndpoint clusterv1beta1.APIEndpoint `json:"controlPlaneEndpoint"`
	// network defines the network configuration for the cluster with different network providers.
	// +optional
	Network Network `json:"network,omitempty,omitzero"`
}

// VSphereClusterStatus defines the observed state of VSphereClusterSpec.
type VSphereClusterStatus struct {
	// Ready indicates the infrastructure required to deploy this cluster is
	// ready.
	// +optional
	Ready bool `json:"ready"`

	// ResourcePolicyName is the name of the VirtualMachineSetResourcePolicy for
	// the cluster, if one exists
	// +optional
	ResourcePolicyName string `json:"resourcePolicyName,omitempty"`

	// Conditions defines current service state of the VSphereCluster.
	// +optional
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`

	// FailureDomains is a list of failure domain objects synced from the
	// infrastructure provider.
	FailureDomains clusterv1beta1.FailureDomains `json:"failureDomains,omitempty"`

	// v1beta2 groups all the fields that will be added or modified in VSphereCluster's status with the V1Beta2 version.
	// +optional
	V1Beta2 *VSphereClusterV1Beta2Status `json:"v1beta2,omitempty"`
}

// VSphereClusterV1Beta2Status groups all the fields that will be added or modified in VSphereClusterStatus with the V1Beta2 version.
// See https://github.com/kubernetes-sigs/cluster-api/blob/main/docs/proposals/20240916-improve-status-in-CAPI-resources.md for more context.
type VSphereClusterV1Beta2Status struct {
	// conditions represents the observations of a VSphereCluster's current state.
	// Known condition types are Ready, ResourcePolicyReady, NetworkReady, LoadBalancerReady,
	// ProviderServiceAccountsReady, ServiceDiscoveryReady and Paused.
	// +optional
	// +listType=map
	// +listMapKey=type
	// +kubebuilder:validation:MaxItems=32
	Conditions []metav1.Condition `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=vsphereclusters,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// VSphereCluster is the Schema for the VSphereClusters API.
type VSphereCluster struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   VSphereClusterSpec   `json:"spec,omitempty"`
	Status VSphereClusterStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// VSphereClusterList contains a list of VSphereCluster.
type VSphereClusterList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []VSphereCluster `json:"items"`
}

// GetConditions returns conditions for VSphereCluster.
func (r *VSphereCluster) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets conditions on the VSphereCluster.
func (r *VSphereCluster) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

// GetV1Beta2Conditions returns the set of conditions for this object.
func (r *VSphereCluster) GetV1Beta2Conditions() []metav1.Condition {
	if r.Status.V1Beta2 == nil {
		return nil
	}
	return r.Status.V1Beta2.Conditions
}

// SetV1Beta2Conditions sets conditions for an API object.
func (r *VSphereCluster) SetV1Beta2Conditions(conditions []metav1.Condition) {
	if r.Status.V1Beta2 == nil {
		r.Status.V1Beta2 = &VSphereClusterV1Beta2Status{}
	}
	r.Status.V1Beta2.Conditions = conditions
}

func init() {
	objectTypes = append(objectTypes, &VSphereCluster{}, &VSphereClusterList{})
}
