/*
Copyright 2018 The Kubernetes Authors.

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

package v1alpha1

import (
	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
)

// Annotation constants
const (
	// ClusterIDLabel is the label that a machineset must have to identify the
	// cluster to which it belongs.
	ClusterIDLabel   = "sigs.k8s.io/cluster-api-cluster"
	MachineRoleLabel = "sigs.k8s.io/cluster-api-machine-role"
	MachineTypeLabel = "sigs.k8s.io/cluster-api-machine-type"
)

// EDIT THIS FILE!  THIS IS SCAFFOLDING FOR YOU TO OWN!
// NOTE: json tags are required.  Any new fields you add must have json tags for the fields to be serialized.

// AWSMachineProviderStatus is the type that will be embedded in a Machine.Status.ProviderStatus field.
// It containsk AWS-specific status information.
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object
type AWSMachineProviderStatus struct {
	metav1.TypeMeta `json:",inline"`

	// InstanceID is the instance ID of the machine created in AWS
	InstanceID *string `json:"instanceId"`

	// InstanceState is the state of the AWS instance for this machine
	InstanceState *string `json:"instanceState"`

	// Conditions is a set of conditions associated with the Machine to indicate
	// errors or other status
	Conditions []AWSMachineProviderCondition `json:"conditions"`
}

// AWSMachineProviderConditionType is a valid value for AWSMachineProviderCondition.Type
type AWSMachineProviderConditionType string

// Valid conditions for an AWS machine instance
const (
	// MachineCreation indicates whether the machine has been created or not. If not,
	// it should include a reason and message for the failure.
	MachineCreation AWSMachineProviderConditionType = "MachineCreation"
)

// AWSMachineProviderCondition is a condition in a AWSMachineProviderStatus
type AWSMachineProviderCondition struct {
	// Type is the type of the condition.
	Type AWSMachineProviderConditionType `json:"type"`
	// Status is the status of the condition.
	Status corev1.ConditionStatus `json:"status"`
	// LastProbeTime is the last time we probed the condition.
	// +optional
	LastProbeTime metav1.Time `json:"lastProbeTime"`
	// LastTransitionTime is the last time the condition transitioned from one status to another.
	// +optional
	LastTransitionTime metav1.Time `json:"lastTransitionTime"`
	// Reason is a unique, one-word, CamelCase reason for the condition's last transition.
	// +optional
	Reason string `json:"reason"`
	// Message is a human-readable message indicating details about last transition.
	// +optional
	Message string `json:"message"`
}

// +genclient
// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineProviderConfig is the Schema for the awsmachineproviderconfigs API
// +k8s:openapi-gen=true
type AWSMachineProviderConfig struct {
	metav1.TypeMeta   `json:",inline"`
	metav1.ObjectMeta `json:"metadata,omitempty"`

	// AMI is the reference to the AMI from which to create the machine instance.
	AMI AWSResourceReference `json:"ami"`

	// InstanceType is the type of instance to create. Example: m4.xlarge
	InstanceType string `json:"instanceType"`

	// Tags is the set of tags to add to apply to an instance, in addition to the ones
	// added by default by the actuator. These tags are additive. The actuator will ensure
	// these tags are present, but will not remove any other tags that may exist on the
	// instance.
	Tags []TagSpecification `json:"tags"`

	// IAMInstanceProfile is a reference to an IAM role to assign to the instance
	IAMInstanceProfile *AWSResourceReference `json:"iamInstanceProfile"`

	// UserDataSecret contains a local reference to a secret that contains the
	// UserData to apply to the instance
	UserDataSecret *corev1.LocalObjectReference `json:"userDataSecret"`

	// CredentialsSecret is a reference to the secret with AWS credentials. Otherwise, defaults to permissions
	// provided by attached IAM role where the actuator is running.
	CredentialsSecret *corev1.LocalObjectReference `json:"credentialsSecret"`

	// KeyName is the name of the KeyPair to use for SSH
	KeyName *string `json:"keyName"`

	// DeviceIndex is the index of the device on the instance for the network interface attachment.
	// Defaults to 0.
	DeviceIndex int64 `json:"deviceIndex"`

	// PublicIP specifies whether the instance should get a public IP. If not present,
	// it should use the default of its subnet.
	PublicIP *bool `json:"publicIp"`

	// SecurityGroups is an array of references to security groups that should be applied to the
	// instance.
	SecurityGroups []AWSResourceReference `json:"securityGroups"`

	// Subnet is a reference to the subnet to use for this instance
	Subnet AWSResourceReference `json:"subnet"`

	// Placement specifies where to create the instance in AWS
	Placement Placement `json:"placement"`

	// LoadBalancers is the set of load balancers to which the new instance
	// should be added once it is created.
	LoadBalancers []LoadBalancerReference `json:"loadBalancers"`
}

// AWSResourceReference is a reference to a specific AWS resource by ID, ARN, or filters.
// Only one of ID, ARN or Filters may be specified. Specifying more than one will result in
// a validation error.
type AWSResourceReference struct {
	// ID of resource
	ID *string `json:"id"`

	// ARN of resource
	ARN *string `json:"arn"`

	// Filters is a set of filters used to identify a resource
	Filters []Filter `json:"filters"`
}

// Placement indicates where to create the instance in AWS
type Placement struct {
	// Region is the region to use to create the instance
	Region string `json:"region"`

	// AvailabilityZone is the availability zone of the instance
	AvailabilityZone string `json:"availabilityZone"`
}

// Filter is a filter used to identify an AWS resource
type Filter struct {
	// Name of the filter. Filter names are case-sensitive.
	Name string `json:"name"`

	// Values includes one or more filter values. Filter values are case-sensitive.
	Values []string `json:"values"`
}

// TagSpecification is the name/value pair for a tag
type TagSpecification struct {
	// Name of the tag
	Name string `json:"name"`

	// Value of the tag
	Value string `json:"value"`
}

// +k8s:deepcopy-gen:interfaces=k8s.io/apimachinery/pkg/runtime.Object

// AWSMachineProviderConfigList contains a list of AWSMachineProviderConfig
type AWSMachineProviderConfigList struct {
	metav1.TypeMeta `json:",inline"`
	metav1.ListMeta `json:"metadata,omitempty"`
	Items           []AWSMachineProviderConfig `json:"items"`
}

// LoadBalancerReference is a reference to a load balancer on AWS.
type LoadBalancerReference struct {
	Name string              `json:"name"`
	Type AWSLoadBalancerType `json:"type"`
}

// AWSLoadBalancerType is the type of LoadBalancer to use when registering
// an instance with load balancers specified in LoadBalancerNames
type AWSLoadBalancerType string

// Possible values for AWSLoadBalancerType. Add to this list as other types
// of load balancer are supported by the actuator.
const (
	ClassicLoadBalancerType AWSLoadBalancerType = "classic" // AWS classic ELB
	NetworkLoadBalancerType AWSLoadBalancerType = "network" // AWS Network Load Balancer (NLB)
)

func init() {
	SchemeBuilder.Register(&AWSMachineProviderConfig{}, &AWSMachineProviderConfigList{}, &AWSMachineProviderStatus{})
}
