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

import clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"

// AWSCluster v1beta2 condition types.
const (
	// AWSClusterReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSCluster.
	AWSClusterReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSClusterVpcReadyV1Beta2Condition reports on the successful reconciliation of a VPC.
	AWSClusterVpcReadyV1Beta2Condition = "VpcReady"

	// AWSClusterSubnetsReadyV1Beta2Condition reports on the successful reconciliation of subnets.
	AWSClusterSubnetsReadyV1Beta2Condition = "SubnetsReady"

	// AWSClusterInternetGatewayReadyV1Beta2Condition reports on the successful reconciliation of internet gateways.
	// Only applicable to managed clusters.
	AWSClusterInternetGatewayReadyV1Beta2Condition = "InternetGatewayReady"

	// AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition reports on the successful reconciliation of egress only internet gateways.
	// Only applicable to managed clusters.
	AWSClusterEgressOnlyInternetGatewayReadyV1Beta2Condition = "EgressOnlyInternetGatewayReady"

	// AWSClusterCarrierGatewayReadyV1Beta2Condition reports on the successful reconciliation of carrier gateways.
	// Only applicable to managed clusters.
	AWSClusterCarrierGatewayReadyV1Beta2Condition = "CarrierGatewayReady"

	// AWSClusterNatGatewaysReadyV1Beta2Condition reports on the successful reconciliation of NAT gateways.
	// Only applicable to managed clusters.
	AWSClusterNatGatewaysReadyV1Beta2Condition = "NatGatewaysReady"

	// AWSClusterRouteTablesReadyV1Beta2Condition reports on the successful reconciliation of route tables.
	// Only applicable to managed clusters.
	AWSClusterRouteTablesReadyV1Beta2Condition = "RouteTablesReady"

	// AWSClusterVpcEndpointsReadyV1Beta2Condition reports on the successful reconciliation of VPC endpoints.
	// Only applicable to managed clusters.
	AWSClusterVpcEndpointsReadyV1Beta2Condition = "VpcEndpointsReady"

	// AWSClusterSecondaryCidrsReadyV1Beta2Condition reports on the successful reconciliation of secondary CIDR blocks.
	// Only applicable to managed clusters.
	AWSClusterSecondaryCidrsReadyV1Beta2Condition = "SecondaryCidrsReady"

	// AWSClusterSecurityGroupsReadyV1Beta2Condition reports on the successful reconciliation of cluster security groups.
	AWSClusterSecurityGroupsReadyV1Beta2Condition = "ClusterSecurityGroupsReady"

	// AWSClusterBastionHostReadyV1Beta2Condition reports whether a bastion host is ready. Depending on the configuration,
	// a cluster may not require a bastion host and this condition will be skipped.
	AWSClusterBastionHostReadyV1Beta2Condition = "BastionHostReady"

	// AWSClusterLoadBalancerReadyV1Beta2Condition reports on whether a control plane load balancer was successfully reconciled.
	AWSClusterLoadBalancerReadyV1Beta2Condition = "LoadBalancerReady"

	// AWSClusterS3BucketReadyV1Beta2Condition reports on the successful reconciliation of an S3 bucket.
	AWSClusterS3BucketReadyV1Beta2Condition = "S3BucketReady"

	// AWSClusterPrincipalCredentialRetrievedV1Beta2Condition reports on whether principal credentials could be retrieved successfully.
	// A possible scenario, where retrieval is unsuccessful, is when SourcePrincipal is not authorized for assume role.
	AWSClusterPrincipalCredentialRetrievedV1Beta2Condition = "PrincipalCredentialRetrieved"

	// AWSClusterPrincipalUsageAllowedV1Beta2Condition reports on whether the principal and all nested source identities
	// are allowed to be used in the AWSCluster namespace.
	AWSClusterPrincipalUsageAllowedV1Beta2Condition = "PrincipalUsageAllowed"
)

// AWSCluster v1beta2 reason constants.
const (
	// AWSClusterReadyV1Beta2Reason indicates the AWSCluster infrastructure is ready.
	AWSClusterReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSClusterNotReadyV1Beta2Reason indicates the AWSCluster infrastructure is not ready.
	AWSClusterNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSClusterDeletingV1Beta2Reason indicates the AWSCluster is being deleted.
	AWSClusterDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSClusterVpcReconciliationFailedV1Beta2Reason used when errors occur during VPC reconciliation.
	AWSClusterVpcReconciliationFailedV1Beta2Reason = "VpcReconciliationFailed"

	// AWSClusterVpcCreationStartedV1Beta2Reason used when attempting to create a VPC for a managed cluster.
	// Will not be applied to unmanaged clusters.
	AWSClusterVpcCreationStartedV1Beta2Reason = "VpcCreationStarted"

	// AWSClusterSubnetsReconciliationFailedV1Beta2Reason used to report failures while reconciling subnets.
	AWSClusterSubnetsReconciliationFailedV1Beta2Reason = "SubnetsReconciliationFailed"

	// AWSClusterInternetGatewayFailedV1Beta2Reason used when errors occur during internet gateway reconciliation.
	AWSClusterInternetGatewayFailedV1Beta2Reason = "InternetGatewayFailed"

	// AWSClusterEgressOnlyInternetGatewayFailedV1Beta2Reason used when errors occur during egress only internet gateway reconciliation.
	AWSClusterEgressOnlyInternetGatewayFailedV1Beta2Reason = "EgressOnlyInternetGatewayFailed"

	// AWSClusterCarrierGatewayFailedV1Beta2Reason used when errors occur during carrier gateway reconciliation.
	AWSClusterCarrierGatewayFailedV1Beta2Reason = "CarrierGatewayFailed"

	// AWSClusterNatGatewaysReconciliationFailedV1Beta2Reason used when any errors occur during reconciliation of NAT gateways.
	AWSClusterNatGatewaysReconciliationFailedV1Beta2Reason = "NatGatewaysReconciliationFailed"

	// AWSClusterNatGatewaysCreationStartedV1Beta2Reason set once when creating new NAT gateways.
	AWSClusterNatGatewaysCreationStartedV1Beta2Reason = "NatGatewaysCreationStarted"

	// AWSClusterRouteTableReconciliationFailedV1Beta2Reason used when any errors occur during reconciliation of route tables.
	AWSClusterRouteTableReconciliationFailedV1Beta2Reason = "RouteTableReconciliationFailed"

	// AWSClusterVpcEndpointsReconciliationFailedV1Beta2Reason used when any errors occur during reconciliation of VPC endpoints.
	AWSClusterVpcEndpointsReconciliationFailedV1Beta2Reason = "VpcEndpointsReconciliationFailed"

	// AWSClusterSecondaryCidrReconciliationFailedV1Beta2Reason used when any errors occur during reconciliation of secondary CIDR blocks.
	AWSClusterSecondaryCidrReconciliationFailedV1Beta2Reason = "SecondaryCidrReconciliationFailed"

	// AWSClusterSecurityGroupReconciliationFailedV1Beta2Reason used when any errors occur during reconciliation of security groups.
	AWSClusterSecurityGroupReconciliationFailedV1Beta2Reason = "SecurityGroupReconciliationFailed"

	// AWSClusterBastionHostFailedV1Beta2Reason used when an error occurs during the creation of a bastion host.
	AWSClusterBastionHostFailedV1Beta2Reason = "BastionHostFailed"

	// AWSClusterBastionCreationStartedV1Beta2Reason used when creating a new bastion host.
	AWSClusterBastionCreationStartedV1Beta2Reason = "BastionCreationStarted"

	// AWSClusterLoadBalancerFailedV1Beta2Reason used when an error occurs during load balancer reconciliation.
	AWSClusterLoadBalancerFailedV1Beta2Reason = "LoadBalancerFailed"

	// AWSClusterWaitForDNSNameV1Beta2Reason used while waiting for a DNS name for the API server to be populated.
	AWSClusterWaitForDNSNameV1Beta2Reason = "WaitForDNSName"

	// AWSClusterWaitForExternalControlPlaneEndpointV1Beta2Reason is set when the AWSCluster is waiting for an externally managed
	// load balancer, such as an external control plane provider.
	AWSClusterWaitForExternalControlPlaneEndpointV1Beta2Reason = "WaitForExternalControlPlaneEndpoint"

	// AWSClusterWaitForDNSNameResolveV1Beta2Reason used while waiting for DNS name to resolve.
	AWSClusterWaitForDNSNameResolveV1Beta2Reason = "WaitForDNSNameResolve"

	// AWSClusterS3BucketFailedV1Beta2Reason used when any errors occur during reconciliation of an S3 bucket.
	AWSClusterS3BucketFailedV1Beta2Reason = "S3BucketCreationFailed"

	// AWSClusterPrincipalCredentialRetrievalFailedV1Beta2Reason used when errors occur during identity credential retrieval.
	AWSClusterPrincipalCredentialRetrievalFailedV1Beta2Reason = "PrincipalCredentialRetrievalFailed"

	// AWSClusterCredentialProviderBuildFailedV1Beta2Reason used when errors occur during building providers before trying credential retrieval.
	//nolint:gosec
	AWSClusterCredentialProviderBuildFailedV1Beta2Reason = "CredentialProviderBuildFailed"

	// AWSClusterPrincipalUsageUnauthorizedV1Beta2Reason used when AWSCluster namespace is not in the identity's allowed namespaces list.
	AWSClusterPrincipalUsageUnauthorizedV1Beta2Reason = "PrincipalUsageUnauthorized"

	// AWSClusterSourcePrincipalUsageUnauthorizedV1Beta2Reason used when AWSCluster namespace is not in the intersection
	// of source identity allowed namespaces and allowed namespaces of the identities that source identity depends on.
	AWSClusterSourcePrincipalUsageUnauthorizedV1Beta2Reason = "SourcePrincipalUsageUnauthorized"
)

// AWSMachine v1beta2 condition types.
const (
	// AWSMachineReadyV1Beta2Condition defines the Ready condition type that summarizes the operational state of an AWSMachine.
	AWSMachineReadyV1Beta2Condition = clusterv1.ReadyCondition

	// AWSMachineInstanceReadyV1Beta2Condition reports on current status of the EC2 instance. Ready indicates the instance is in a Running state.
	AWSMachineInstanceReadyV1Beta2Condition = "InstanceReady"

	// AWSMachineSecurityGroupsReadyV1Beta2Condition indicates the security groups are up to date on the AWSMachine.
	AWSMachineSecurityGroupsReadyV1Beta2Condition = "SecurityGroupsReady"

	// AWSMachineELBAttachedV1Beta2Condition will report true when a control plane is successfully registered with an ELB.
	// When set to false, the subnet may not be found or unavailable in the instance's AZ.
	// Only applicable to control plane machines.
	AWSMachineELBAttachedV1Beta2Condition = "ELBAttached"

	// AWSMachineDedicatedHostReleaseV1Beta2Condition reports on the status of dedicated host release operations.
	// This condition tracks whether the dedicated host has been successfully released or if there are failures.
	AWSMachineDedicatedHostReleaseV1Beta2Condition = "DedicatedHostRelease"
)

// AWSMachine v1beta2 reason constants.
const (
	// AWSMachineReadyV1Beta2Reason indicates the AWSMachine is ready.
	AWSMachineReadyV1Beta2Reason = clusterv1.ReadyReason

	// AWSMachineNotReadyV1Beta2Reason indicates the AWSMachine is not ready.
	AWSMachineNotReadyV1Beta2Reason = clusterv1.NotReadyReason

	// AWSMachineDeletingV1Beta2Reason indicates the AWSMachine is being deleted.
	AWSMachineDeletingV1Beta2Reason = clusterv1.DeletingReason

	// AWSMachineInstanceNotFoundV1Beta2Reason used when the instance couldn't be retrieved.
	AWSMachineInstanceNotFoundV1Beta2Reason = "InstanceNotFound"

	// AWSMachineInstanceTerminatedV1Beta2Reason used when the instance is in a terminated state.
	AWSMachineInstanceTerminatedV1Beta2Reason = "InstanceTerminated"

	// AWSMachineInstanceStoppedV1Beta2Reason used when the instance is in a stopped state.
	AWSMachineInstanceStoppedV1Beta2Reason = "InstanceStopped"

	// AWSMachineInstanceNotReadyV1Beta2Reason used when the instance is in a pending state.
	AWSMachineInstanceNotReadyV1Beta2Reason = "InstanceNotReady"

	// AWSMachineInstanceProvisionStartedV1Beta2Reason set when the provisioning of an instance started.
	AWSMachineInstanceProvisionStartedV1Beta2Reason = "InstanceProvisionStarted"

	// AWSMachineInstanceProvisionFailedV1Beta2Reason used for failures during instance provisioning.
	AWSMachineInstanceProvisionFailedV1Beta2Reason = "InstanceProvisionFailed"

	// AWSMachineWaitingForClusterInfrastructureV1Beta2Reason used when machine is waiting for cluster infrastructure to be ready before proceeding.
	AWSMachineWaitingForClusterInfrastructureV1Beta2Reason = clusterv1.WaitingForClusterInfrastructureReadyReason

	// AWSMachineWaitingForBootstrapDataV1Beta2Reason used when machine is waiting for bootstrap data to be ready before proceeding.
	AWSMachineWaitingForBootstrapDataV1Beta2Reason = clusterv1.WaitingForBootstrapDataReason

	// AWSMachineSecurityGroupsFailedV1Beta2Reason used when the security groups could not be synced.
	AWSMachineSecurityGroupsFailedV1Beta2Reason = "SecurityGroupsSyncFailed"

	// AWSMachineELBAttachFailedV1Beta2Reason used when a control plane node fails to attach to the ELB.
	AWSMachineELBAttachFailedV1Beta2Reason = "ELBAttachFailed"

	// AWSMachineELBDetachFailedV1Beta2Reason used when a control plane node fails to detach from an ELB.
	AWSMachineELBDetachFailedV1Beta2Reason = "ELBDetachFailed"

	// AWSMachineDedicatedHostReleaseFailedV1Beta2Reason used when the dedicated host release fails.
	AWSMachineDedicatedHostReleaseFailedV1Beta2Reason = "DedicatedHostReleaseFailed"
)
