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
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/apimachinery/pkg/util/intstr"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	rosacontrolplanev1 "sigs.k8s.io/cluster-api-provider-aws/v2/controlplane/rosa/api/v1beta2"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
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

	// Version specifies the OpenShift version of the nodes associated with this machinepool.
	// ROSAControlPlane version is used if not set.
	//
	// +optional
	Version string `json:"version,omitempty"`

	// AvailabilityZone is an optinal field specifying the availability zone where instances of this machine pool should run
	// For Multi-AZ clusters, you can create a machine pool in a Single-AZ of your choice.
	// +optional
	AvailabilityZone string `json:"availabilityZone,omitempty"`

	// +kubebuilder:validation:XValidation:rule="self == oldSelf", message="subnet is immutable"
	// +immutable
	// +optional
	Subnet string `json:"subnet,omitempty"`

	// Labels specifies labels for the Kubernetes node objects
	// +optional
	Labels map[string]string `json:"labels,omitempty"`

	// Taints specifies the taints to apply to the nodes of the machine pool
	// +optional
	Taints []RosaTaint `json:"taints,omitempty"`

	// AdditionalTags are user-defined tags to be added on the underlying EC2 instances associated with this machine pool.
	// +immutable
	// +optional
	AdditionalTags infrav1.Tags `json:"additionalTags,omitempty"`

	// AutoRepair specifies whether health checks should be enabled for machines
	// in the NodePool. The default is true.
	// +kubebuilder:default=true
	// +optional
	AutoRepair bool `json:"autoRepair,omitempty"`

	// InstanceType specifies the AWS instance type, for example `r5.xlarge`. Instance type ref; https://aws.amazon.com/ec2/instance-types/
	//
	// +kubebuilder:validation:Required
	InstanceType string `json:"instanceType"`

	// ImageType is the AMI (Amazon Machine Image) to use for running the associated NodePool (i.e. Windows or Default/Linux).
	//
	// +kubebuilder:validation:Enum=Windows;Default
	// +optional
	ImageType string `json:"imageType,omitempty"`

	// Autoscaling specifies auto scaling behaviour for this MachinePool.
	// required if Replicas is not configured
	// +optional
	Autoscaling *rosacontrolplanev1.AutoScaling `json:"autoscaling,omitempty"`

	// TuningConfigs specifies the names of the tuning configs to be applied to this MachinePool.
	// Tuning configs must already exist.
	// +optional
	TuningConfigs []string `json:"tuningConfigs,omitempty"`

	// AdditionalSecurityGroups is an optional set of security groups to associate
	// with all node instances of the machine pool.
	//
	// +immutable
	// +optional
	AdditionalSecurityGroups []string `json:"additionalSecurityGroups,omitempty"`

	// VolumeSize set the disk volume size for the machine pool, in Gib. The default is 300 GiB.
	// +kubebuilder:validation:Minimum=75
	// +kubebuilder:validation:Maximum=16384
	// +immutable
	// +optional
	VolumeSize int `json:"volumeSize,omitempty"`

	// ProviderIDList contain a ProviderID for each machine instance that's currently managed by this machine pool.
	// +optional
	ProviderIDList []string `json:"providerIDList,omitempty"`

	// NodeDrainGracePeriod is grace period for how long Pod Disruption Budget-protected workloads will be
	// respected during upgrades. After this grace period, any workloads protected by Pod Disruption
	// Budgets that have not been successfully drained from a node will be forcibly evicted.
	//
	// Valid values are from 0 to 1 week(10080m|168h) .
	// 0 or empty value means that the MachinePool can be drained without any time limitation.
	//
	// +optional
	NodeDrainGracePeriod *metav1.Duration `json:"nodeDrainGracePeriod,omitempty"`

	// UpdateConfig specifies update configurations.
	//
	// +optional
	UpdateConfig *RosaUpdateConfig `json:"updateConfig,omitempty"`

	// CapacityReservationID specifies the ID of an AWS On-Demand Capacity Reservation and Capacity Blocks for ML.
	// The CapacityReservationID must be pre-created in advance, before creating a NodePool.
	//
	// +optional
	CapacityReservationID string `json:"capacityReservationID,omitempty"`
}

// RosaTaint represents a taint to be applied to a node.
type RosaTaint struct {
	// The taint key to be applied to a node.
	//
	// +kubebuilder:validation:Required
	Key string `json:"key"`
	// The taint value corresponding to the taint key.
	//
	// +kubebuilder:validation:Pattern:=`^(([A-Za-z0-9][-A-Za-z0-9_.]*)?[A-Za-z0-9])?$`
	// +optional
	Value string `json:"value,omitempty"`
	// The effect of the taint on pods that do not tolerate the taint.
	// Valid effects are NoSchedule, PreferNoSchedule and NoExecute.
	//
	// +kubebuilder:validation:Required
	// +kubebuilder:validation:Enum=NoSchedule;PreferNoSchedule;NoExecute
	Effect corev1.TaintEffect `json:"effect"`
}

// RosaUpdateConfig specifies update configuration
type RosaUpdateConfig struct {
	// RollingUpdate specifies MaxUnavailable & MaxSurge number of nodes during update.
	//
	// +optional
	RollingUpdate *RollingUpdate `json:"rollingUpdate,omitempty"`
}

// RollingUpdate specifies MaxUnavailable & MaxSurge number of nodes during update.
type RollingUpdate struct {
	// MaxUnavailable is the maximum number of nodes that can be unavailable during the update.
	// Value can be an absolute number (ex: 5) or a percentage of desired nodes (ex: 10%).
	// Absolute number is calculated from percentage by rounding down.
	//
	// MaxUnavailable can not be 0 if MaxSurge is 0, default is 0.
	// Both MaxUnavailable & MaxSurge must use the same units (absolute value or percentage).
	//
	// Example: when MaxUnavailable is set to 30%, old nodes can be deleted down to 70% of
	// desired nodes immediately when the rolling update starts. Once new nodes
	// are ready, more old nodes be deleted, followed by provisioning new nodes,
	// ensuring that the total number of nodes available at all times during the
	// update is at least 70% of desired nodes.
	//
	// +kubebuilder:validation:Pattern="^((100|[0-9]{1,2})%|[0-9]+)$"
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:default=0
	// +optional
	MaxUnavailable *intstr.IntOrString `json:"maxUnavailable,omitempty"`

	// MaxSurge is the maximum number of nodes that can be provisioned above the desired number of nodes.
	// Value can be an absolute number (ex: 5) or a percentage of desired nodes (ex: 10%).
	// Absolute number is calculated from percentage by rounding up.
	//
	// MaxSurge can not be 0 if MaxUnavailable is 0, default is 1.
	// Both MaxSurge & MaxUnavailable must use the same units (absolute value or percentage).
	//
	// Example: when MaxSurge is set to 30%, new nodes can be provisioned immediately
	// when the rolling update starts, such that the total number of old and new
	// nodes do not exceed 130% of desired nodes. Once old nodes have been
	// deleted, new nodes can be provisioned, ensuring that total number of nodes
	// running at any time during the update is at most 130% of desired nodes.
	//
	// +kubebuilder:validation:Pattern="^((100|[0-9]{1,2})%|[0-9]+)$"
	// +kubebuilder:validation:XIntOrString
	// +kubebuilder:default=1
	// +optional
	MaxSurge *intstr.IntOrString `json:"maxSurge,omitempty"`
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
	Conditions clusterv1beta1.Conditions `json:"conditions,omitempty"`
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

	// Available upgrades for the ROSA MachinePool.
	AvailableUpgrades []string `json:"availableUpgrades,omitempty"`
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
func (r *ROSAMachinePool) GetConditions() clusterv1beta1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the RosaMachinePool to the predescribed clusterv1beta1.Conditions.
func (r *ROSAMachinePool) SetConditions(conditions clusterv1beta1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSAMachinePool{}, &ROSAMachinePoolList{})
}
