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

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	// InstanceReadyCondition reports on current status of the OpenStack instance. Ready indicates the instance is in a Running state.
	InstanceReadyCondition clusterv1.ConditionType = "InstanceReady"

	// WaitingForClusterInfrastructureReason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	WaitingForClusterInfrastructureReason = "WaitingForClusterInfrastructure"
	// WaitingForBootstrapDataReason used when machine is waiting for bootstrap data to be ready before proceeding.
	WaitingForBootstrapDataReason = "WaitingForBootstrapData"
	// InvalidMachineSpecReason used when the machine spec is invalid.
	InvalidMachineSpecReason = "InvalidMachineSpec"
	// InstanceCreateFailedReason used when creating the instance failed.
	InstanceCreateFailedReason = "InstanceCreateFailed"
	// InstanceNotFoundReason used when the instance couldn't be retrieved.
	InstanceNotFoundReason = "InstanceNotFound"
	// InstanceStateErrorReason used when the instance is in error state.
	InstanceStateErrorReason = "InstanceStateError"
	// InstanceDeletedReason used when the instance is in a deleted state.
	InstanceDeletedReason = "InstanceDeleted"
	// InstanceNotReadyReason used when the instance is in a pending state.
	InstanceNotReadyReason = "InstanceNotReady"
	// InstanceDeleteFailedReason used when deleting the instance failed.
	InstanceDeleteFailedReason = "InstanceDeleteFailed"
	// OpenstackErrorReason used when there is an error communicating with OpenStack.
	OpenStackErrorReason = "OpenStackError"
	// DependencyFailedReason indicates that a dependent object failed.
	DependencyFailedReason = "DependencyFailed"

	// ServerUnexpectedDeletedMessage is the message used when the server is unexpectedly deleted via an external agent.
	ServerUnexpectedDeletedMessage = "The server was unexpectedly deleted"
)

const (
	// APIServerIngressReadyCondition reports on the current status of the network ingress (Loadbalancer, Floating IP) for Control Plane machines. Ready indicates that the instance can receive requests.
	APIServerIngressReadyCondition clusterv1.ConditionType = "APIServerIngressReadyCondition"

	// LoadBalancerMemberErrorReason used when the instance could not be added as a loadbalancer member.
	LoadBalancerMemberErrorReason = "LoadBalancerMemberError"
	// FloatingIPErrorReason used when the floating ip could not be created or attached.
	FloatingIPErrorReason = "FloatingIPError"
)

const (
	// FloatingAddressFromPoolReadyCondition reports on the current status of the Floating IPs from ipam pool.
	FloatingAddressFromPoolReadyCondition clusterv1.ConditionType = "FloatingAddressFromPoolReady"
	// WaitingForIpamProviderReason used when machine is waiting for ipam provider to be ready before proceeding.
	FloatingAddressFromPoolWaitingForIpamProviderReason = "WaitingForIPAMProvider"
	// FloatingAddressFromPoolErrorReason is used when there is an error attaching an IP from the pool to an machine.
	FloatingAddressFromPoolErrorReason = "FloatingIPError"
	// UnableToFindFloatingIPNetworkReason is used when the floating ip network is not found.
	UnableToFindFloatingIPNetworkReason = "UnableToFindFloatingIPNetwork"
)
