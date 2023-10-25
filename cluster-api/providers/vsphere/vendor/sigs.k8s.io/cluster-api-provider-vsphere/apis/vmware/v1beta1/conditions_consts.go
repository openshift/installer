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

//nolint:godot
package v1beta1

import clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"

const (
	// ResourcePolicyReadyCondition reports the successful creation of a
	// Resource Policy.
	ResourcePolicyReadyCondition clusterv1.ConditionType = "ResourcePolicyReady"

	// ResourcePolicyCreationFailedReason used when any errors occur during
	// ResourcePolicy creation.
	ResourcePolicyCreationFailedReason = "ResourcePolicyCreationFailed"
)

const (
	// ClusterNetworkReadyCondition reports the successful provision of a
	// Cluster Network.
	ClusterNetworkReadyCondition clusterv1.ConditionType = "ClusterNetworkReady"

	// ClusterNetworkProvisionStarted is used when waiting for Cluster
	// Network to be Ready.
	ClusterNetworkProvisionStartedReason = "ClusterNetworkProvisionStarted"
	// ClusterNetworkProvisionFailedReason is used when any errors occur
	// during network provision.
	ClusterNetworkProvisionFailedReason = "ClusterNetworkProvisionFailed"
)

const (
	// LoadBalancerReadyCondition reports the successful reconciliation of
	// a static control plane endpoint.
	LoadBalancerReadyCondition clusterv1.ConditionType = "LoadBalancerReady"

	// LoadBalancerCreationFailedReason is used when load balancer related
	// resources creation fails.
	LoadBalancerCreationFailedReason = "LoadBalancerCreationFailed"
	// WaitingForLoadBalancerIPReason is used when waiting for load
	// balancer IP to exist.
	WaitingForLoadBalancerIPReason = "WaitingForLoadBalancerIP"
)

// Conditions and condition Reasons for VSphereMachine.
const (
	// ConditionType VMProvisionedCondition is shared with infrav1.VSPhereMachine
	// VMCreationFailedReason reports that creating VM CRD or corresponding bootstrap ConfigMap failed.
	VMCreationFailedReason = "VMCreationFailed"
	// VMProvisionStartedReason documents (Severity=Info) a Virtual Machine currently is in creation process.
	VMProvisionStartedReason = "VMProvisionStarted"
	// PoweringOnReason documents (Severity=Info) a Virtual Machine currently executing the power on sequence.
	PoweringOnReason = "PoweringOn"
	// WaitingForNetworkAddressReason (Severity=Info) documents a VSphereMachine waiting for the the machine network
	// settings to be reported after machine being powered on.
	WaitingForNetworkAddressReason = "WaitingForNetworkAddress"
	// WaitingForBIOSUUIDReason (Severity=Info) documents a VSphereMachine waiting for the the machine to have a BIOS UUID.
	WaitingForBIOSUUIDReason = "WaitingForBIOSUUID"
)

const (
	// ProviderServiceAccountsReadyCondition documents the status of provider service accounts
	// and related Roles, RoleBindings and Secrets are created
	ProviderServiceAccountsReadyCondition clusterv1.ConditionType = "ProviderServiceAccountsReady"

	// ProviderServiceAccountsReconciliationFailedReason reports that provider service accounts related resources reconciliation failed
	ProviderServiceAccountsReconciliationFailedReason = "ProviderServiceAccountsReconciliationFailed"
)

const (
	SupervisorLoadBalancerSvcNamespace = "kube-system"
	SupervisorLoadBalancerSvcName      = "kube-apiserver-lb-svc"
	SupervisorAPIServerPort            = 6443

	SupervisorHeadlessSvcNamespace = "default"
	SupervisorHeadlessSvcName      = "supervisor"
	SupervisorHeadlessSvcPort      = 6443

	// ServiceDiscoveryReadyCondition documents the status of service discoveries
	ServiceDiscoveryReadyCondition clusterv1.ConditionType = "ServiceDiscoveryReady"

	// SupervisorHeadlessServiceSetupFailedReason documents the headless service setup for svc api server failed
	SupervisorHeadlessServiceSetupFailedReason = "SupervisorHeadlessServiceSetupFailed"
)
