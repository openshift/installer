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

	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// RosaMachinePoolSpec defines the desired state of RosaMachinePool.
type RosaMachinePoolSpec struct {
	// NodePoolName specifies the name of the nodepool in Rosa
	// must be a valid DNS-1035 label, so it must consist of lower case alphanumeric and have a max length of 15 characters.
	//
	// +immutable
	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="nodepoolName is immutable"
	// +kubebuilder:validation:MaxLength:=15
	// +kubebuilder:validation:Pattern:=`^[a-z]([-a-z0-9]*[a-z0-9])?$`
	NodePoolName string `json:"nodePoolName"`

	// Version specifies the penshift version of the nodes associated with this machinepool.
	// ROSAControlPlane version is used if not set.
	//
	// +optional
	// +kubebuilder:validation:XValidation:rule=`self.matches('^(0|[1-9]\\d*)\\.(0|[1-9]\\d*)\\.(0|[1-9]\\d*)$')`, message="version must be a valid semantic version"
	Version string `json:"version,omitempty"`

	// AvailabilityZone is an optinal field specifying the availability zone where instances of this machine pool should run
	// For Multi-AZ clusters, you can create a machine pool in a Single-AZ of your choice.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// +optional
	Subnet string `json:"subnet,omitempty"`

	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// AutoRepair specifies whether health checks should be enabled for machines
	// in the NodePool. The default is false.
	// +optional
	// +kubebuilder:default=false
	AutoRepair bool `json:"autoRepair,omitempty"`

	// InstanceType specifies the AWS instance type
	InstanceType string `json:"instanceType,omitempty"`

	// Autoscaling specifies auto scaling behaviour for this MachinePool.
	// required if Replicas is not configured
	// +optional
	Autoscaling *RosaMachinePoolAutoScaling `json:"autoscaling,omitempty"`

	// TODO(alberto): Enable and propagate this API input.
	// Taints           []*Taint                     `json:"taints,omitempty"`
	// TuningConfigs    []string                     `json:"tuningConfigs,omitempty"`
	// Version          *Version                     `json:"version,omitempty"`

	// ProviderIDList contain a ProviderID for each machine instance that's currently managed by this machine pool.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`
}

// RosaMachinePoolAutoScaling specifies scaling options.
type RosaMachinePoolAutoScaling struct {
	// +kubebuilder:validation:Minimum=1
	MinReplicas int `json:"minReplicas,omitempty"`
	// +kubebuilder:validation:Minimum=1
	MaxReplicas int `json:"maxReplicas,omitempty"`
}

// RosaMachinePoolStatus defines the observed state of RosaMachinePool.
type RosaMachinePoolStatus struct {
	// Ready denotes that the RosaMachinePool nodepool has joined
	// the cluster
	// +kubebuilder:default=false
	Ready bool `json:"ready"`
	// Replicas is the most recently observed number of replicas.
	// +optional
	Replicas int32 `json:"replicas"`
	// Conditions defines current service state of the managed machine pool
	// +optional
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
	// FailureMessage will be set in the event that there is a terminal problem
	// reconciling the state and will be set to a descriptive error message.
	//
	// This field should not be set for transitive errors that a controller
	// faces that are expected to be fixed automatically over
	// time (like service outages), but instead indicate that something is
	// fundamentally wrong with the spec or the configuration of
	// the controller, and that manual intervention is required.
	//
	// +optional
	FailureMessage *string `json:"failureMessage,omitempty"`

	// ID is the ID given by ROSA.
	ID string `json:"id,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosamachinepools,scope=Namespaced,categories=cluster-api,shortName=rosamp
// +kubebuilder:storageversion
// +kubebuilder:subresource:status
// +kubebuilder:printcolumn:name="Ready",type="string",JSONPath=".status.ready",description="MachinePool ready status"
// +kubebuilder:printcolumn:name="Replicas",type="integer",JSONPath=".status.replicas",description="Number of replicas"

// ROSAMachinePool is the Schema for the rosamachinepools API.
type ROSAMachinePool struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   RosaMachinePoolSpec   `json:"spec,omitempty"`
	Status RosaMachinePoolStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ROSAMachinePoolList contains a list of RosaMachinePools.
type ROSAMachinePoolList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSAMachinePool `json:"items"`
}

// GetConditions returns the observations of the operational state of the RosaMachinePool resource.
func (r *ROSAMachinePool) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the RosaMachinePool to the predescribed clusterv1.Conditions.
func (r *ROSAMachinePool) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSAMachinePool{}, &ROSAMachinePoolList{})
}
