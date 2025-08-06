/*
Copyright The Kubernetes Authors.

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

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
)

// ROSANetworkFinalizer allows the controller to clean up resources on delete.
const ROSANetworkFinalizer = "rosanetwork.infrastructure.cluster.x-k8s.io"

// ROSANetworkSpec defines the desired state of ROSANetwork
type ROSANetworkSpec struct {
	// The name of the cloudformation stack under which the network infrastructure would be created
	// +immutable
	StackName string `json:"stackName"`

	// The AWS region in which the components of ROSA network infrastruture are to be crated
	// +immutable
	Region string `json:"region"`

	// The number of availability zones to be used for creation of the network infrastructure.
	// You can specify anything between one and four, depending on the chosen AWS region.
	// +kubebuilder:default=1
	// +optional
	// +immutable
	AvailabilityZoneCount int `json:"availabilityZoneCount"`

	// The list of availability zones to be used for creation of the network infrastructure.
	// You can specify anything between one and four valid availability zones from a given region.
	// Should you specify both the availabilityZoneCount and availabilityZones, the list of availability zones takes preference.
	// +optional
	// +immutable
	AvailabilityZones []string `json:"availabilityZones"`

	// CIDR block to be used for the VPC
	// +kubebuilder:validation:Format=cidr
	// +immutable
	CIDRBlock string `json:"cidrBlock"`

	// IdentityRef is a reference to an identity to be used when reconciling rosa network.
	// If no identity is specified, the default identity for this controller will be used.
	//
	// +optional
	IdentityRef *infrav1.AWSIdentityReference `json:"identityRef,omitempty"`

	// StackTags is an optional set of tags to add to the created cloudformation stack.
	// The stack tags will then be automatically applied to the supported AWS resources (VPC, subnets, ...).
	//
	// +optional
	StackTags Tags `json:"stackTags,omitempty"`
}

// ROSANetworkSubnet groups public and private subnet and the availability zone in which the two subnets got created
type ROSANetworkSubnet struct {
	// Availability zone of the subnet pair, for example us-west-2a
	AvailabilityZone string `json:"availabilityZone"`

	// ID of the public subnet, for example subnet-0f7e49a3ce68ff338
	PublicSubnet string `json:"publicSubnet"`

	// ID of the private subnet, for example subnet-07a20d6c41af2b725
	PrivateSubnet string `json:"privateSubnet"`
}

// CFResource groups information pertaining to a resource created as a part of a cloudformation stack
type CFResource struct {
	// Type of the created resource: AWS::EC2::VPC, AWS::EC2::Subnet, ...
	ResourceType string `json:"resource"`

	// LogicalResourceID of the created resource.
	LogicalID string `json:"logicalId"`

	// PhysicalResourceID of the created resource.
	PhysicalID string `json:"physicalId"`

	// Status of the resource: CREATE_IN_PROGRESS, CREATE_COMPLETE, ...
	Status string `json:"status"`

	// Message pertaining to the status of the resource
	Reason string `json:"reason"`
}

// ROSANetworkStatus defines the observed state of ROSANetwork
type ROSANetworkStatus struct {
	// Array of created private, public subnets and availability zones, grouped by availability zones
	Subnets []ROSANetworkSubnet `json:"subnets,omitempty"`

	// Resources created in the cloudformation stack
	Resources []CFResource `json:"resources,omitempty"`

	// Conditions specifies the conditions for ROSANetwork
	Conditions clusterv1.Conditions `json:"conditions,omitempty"`
}

// +kubebuilder:object:root=true
// +kubebuilder:resource:path=rosanetworks,shortName=rosanet,scope=Namespaced,categories=cluster-api
// +kubebuilder:storageversion
// +kubebuilder:subresource:status

// ROSANetwork is the schema for the rosanetworks API
type ROSANetwork struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	Spec   ROSANetworkSpec   `json:"spec,omitempty"`
	Status ROSANetworkStatus `json:"status,omitempty"`
}

// +kubebuilder:object:root=true

// ROSANetworkList contains a list of ROSANetwork
type ROSANetworkList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []ROSANetwork `json:"items"`
}

// GetConditions returns the observations of the operational state of the ROSANetwork resource.
func (r *ROSANetwork) GetConditions() clusterv1.Conditions {
	return r.Status.Conditions
}

// SetConditions sets the underlying service state of the ROSANetwork to the predescribed clusterv1.Conditions.
func (r *ROSANetwork) SetConditions(conditions clusterv1.Conditions) {
	r.Status.Conditions = conditions
}

func init() {
	SchemeBuilder.Register(&ROSANetwork{}, &ROSANetworkList{})
}
