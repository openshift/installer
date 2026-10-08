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
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
)

const (
	// ResourcePolicyReadyV1Beta1Condition reports the successful creation of a
	// Resource Policy.
	ResourcePolicyReadyV1Beta1Condition clusterv1.ConditionType = "ResourcePolicyReady"

	// ResourcePolicyCreationFailedV1Beta1Reason used when any errors occur during
	// ResourcePolicy creation.
	ResourcePolicyCreationFailedV1Beta1Reason = "ResourcePolicyCreationFailed"
)

const (
	// ClusterNetworkReadyV1Beta1Condition reports the successful provision of a
	// Cluster Network.
	ClusterNetworkReadyV1Beta1Condition clusterv1.ConditionType = "ClusterNetworkReady"

	// ClusterNetworkProvisionStartedV1Beta1Reason is used when waiting for Cluster
	// Network to be Ready.
	ClusterNetworkProvisionStartedV1Beta1Reason = "ClusterNetworkProvisionStarted"

	// ClusterNetworkProvisionFailedV1Beta1Reason is used when any errors occur
	// during network provision.
	ClusterNetworkProvisionFailedV1Beta1Reason = "ClusterNetworkProvisionFailed"
)

const (
	// LoadBalancerReadyV1Beta1Condition reports the successful reconciliation of
	// a static control plane endpoint.
	LoadBalancerReadyV1Beta1Condition clusterv1.ConditionType = "LoadBalancerReady"

	// LoadBalancerCreationFailedV1Beta1Reason is used when load balancer related
	// resources creation fails.
	LoadBalancerCreationFailedV1Beta1Reason = "LoadBalancerCreationFailed"

	// WaitingForLoadBalancerIPV1Beta1Reason is used when waiting for load
	// balancer IP to exist.
	WaitingForLoadBalancerIPV1Beta1Reason = "WaitingForLoadBalancerIP"

	// WaitingForKubeadmControlPlaneSpecReadyV1Beta1Reason is used when waiting for KubeadmControlPlane to have
	// correct certSANs and observedGeneration equal to generation.
	WaitingForKubeadmControlPlaneSpecReadyV1Beta1Reason = "WaitingForKubeadmControlPlaneSpecReady"
)

// Conditions and condition Reasons for VSphereMachine.
// ConditionType VMProvisionedCondition is shared with infrav1.VSphereMachine.
const (
	// VMCreationFailedV1Beta1Reason reports that creating VM CRD or corresponding bootstrap ConfigMap failed.
	VMCreationFailedV1Beta1Reason = "VMCreationFailed"
	// VMProvisionStartedV1Beta1Reason documents (Severity=Info) a Virtual Machine currently is in creation process.
	VMProvisionStartedV1Beta1Reason = "VMProvisionStarted"
	// PoweringOnV1Beta1Reason documents (Severity=Info) a Virtual Machine currently executing the power on sequence.
	PoweringOnV1Beta1Reason = "PoweringOn"
	// WaitingForNetworkAddressV1Beta1Reason (Severity=Info) documents a VSphereMachine waiting for the machine network
	// settings to be reported after machine being powered on.
	WaitingForNetworkAddressV1Beta1Reason = "WaitingForNetworkAddress"
	// WaitingForBIOSUUIDV1Beta1Reason (Severity=Info) documents a VSphereMachine waiting for the machine to have a BIOS UUID.
	WaitingForBIOSUUIDV1Beta1Reason = "WaitingForBIOSUUID"
)

const (
	// ProviderServiceAccountsReadyV1Beta1Condition documents the status of provider service accounts
	// and related Roles, RoleBindings and Secrets are created.
	ProviderServiceAccountsReadyV1Beta1Condition clusterv1.ConditionType = "ProviderServiceAccountsReady"

	// ProviderServiceAccountsReconciliationFailedV1Beta1Reason reports that provider service accounts related resources reconciliation failed.
	ProviderServiceAccountsReconciliationFailedV1Beta1Reason = "ProviderServiceAccountsReconciliationFailed"
)

const (
	// SupervisorLoadBalancerSvcNamespace is the namespace for the Supervisor load balancer service.
	SupervisorLoadBalancerSvcNamespace = "kube-system"
	// SupervisorLoadBalancerSvcName is the name for the Supervisor load balancer service.
	SupervisorLoadBalancerSvcName = "kube-apiserver-lb-svc"
	// SupervisorAPIServerPort is the port for the Supervisor apiserver when using the load balancer service.
	SupervisorAPIServerPort = 6443
	// SupervisorHeadlessSvcNamespace is the namespace for the Supervisor headless service.
	SupervisorHeadlessSvcNamespace = "default"
	// SupervisorHeadlessSvcName is the name for the Supervisor headless service.
	SupervisorHeadlessSvcName = "supervisor"
	// SupervisorHeadlessSvcPort is the port for the Supervisor apiserver when using the headless service.
	SupervisorHeadlessSvcPort = 6443

	// ServiceDiscoveryReadyV1Beta1Condition documents the status of service discoveries.
	ServiceDiscoveryReadyV1Beta1Condition clusterv1.ConditionType = "ServiceDiscoveryReady"

	// SupervisorHeadlessServiceSetupFailedV1Beta1Reason documents the headless service setup for svc api server failed.
	SupervisorHeadlessServiceSetupFailedV1Beta1Reason = "SupervisorHeadlessServiceSetupFailed"
)
