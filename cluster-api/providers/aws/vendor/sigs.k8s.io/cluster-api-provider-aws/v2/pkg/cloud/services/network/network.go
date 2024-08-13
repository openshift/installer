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

package network

import (
	"k8s.io/klog/v2"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-aws/v2/pkg/cloud/awserrors"
	infrautilconditions "sigs.k8s.io/cluster-api-provider-aws/v2/util/conditions"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
)

// ReconcileNetwork reconciles the network of the given cluster.
func (s *Service) ReconcileNetwork() (err error) {
	s.scope.Debug("Reconciling network for cluster", "cluster", klog.KRef(s.scope.Namespace(), s.scope.Name()))

	// VPC.
	if err := s.reconcileVPC(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, infrav1.VpcReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}
	conditions.MarkTrue(s.scope.InfraCluster(), infrav1.VpcReadyCondition)

	// Secondary CIDRs
	if err := s.associateSecondaryCidrs(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, infrav1.SecondaryCidrReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// Subnets.
	if err := s.reconcileSubnets(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, infrav1.SubnetsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// Internet Gateways.
	if err := s.reconcileInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, infrav1.InternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// Carrier Gateway.
	if err := s.reconcileCarrierGateway(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, infrav1.CarrierGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// Egress Only Internet Gateways.
	if err := s.reconcileEgressOnlyInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, infrav1.EgressOnlyInternetGatewayFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// NAT Gateways.
	if err := s.reconcileNatGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, infrav1.NatGatewaysReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// Routing tables.
	if err := s.reconcileRouteTables(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, infrav1.RouteTableReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	// VPC Endpoints.
	if err := s.reconcileVPCEndpoints(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, infrav1.VpcEndpointsReconciliationFailedReason, infrautilconditions.ErrorConditionAfterInit(s.scope.ClusterObj()), err.Error())
		return err
	}

	s.scope.Debug("Reconcile network completed successfully")
	return nil
}

// DeleteNetwork deletes the network of the given cluster.
func (s *Service) DeleteNetwork() (err error) {
	s.scope.Debug("Deleting network")

	vpc := &infrav1.VPCSpec{}
	// Get VPC used for the cluster
	if s.scope.VPC().ID != "" {
		var err error
		vpc, err = s.describeVPCByID()
		if err != nil {
			if awserrors.IsNotFound(err) {
				// If the VPC does not exist, nothing to do
				return nil
			}
			return err
		}
	} else {
		s.scope.Error(err, "non-fatal: VPC ID is missing, ")
	}

	vpc.DeepCopyInto(s.scope.VPC())

	// VPC Endpoints.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPCEndpoints(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcEndpointsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// Routing tables.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteRouteTables(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.RouteTablesReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// NAT Gateways.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteNatGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.NatGatewaysReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// EIPs.
	if err := s.releaseAddresses(); err != nil {
		return err
	}

	// Internet Gateways.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.InternetGatewayReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// Carrier Gateway.
	if s.scope.VPC().CarrierGatewayID != nil {
		if err := s.deleteCarrierGateway(); err != nil {
			conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
			return err
		}
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.CarrierGatewayReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")
	}

	// Egress Only Internet Gateways.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteEgressOnlyInternetGateways(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.EgressOnlyInternetGatewayReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// Subnets.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteSubnets(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SubnetsReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	// Secondary CIDR.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.disassociateSecondaryCidrs(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.SecondaryCidrsReadyCondition, "DisassociateFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}

	// VPC.
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1.DeletingReason, clusterv1.ConditionSeverityInfo, "")
	if err := s.scope.PatchObject(); err != nil {
		return err
	}

	if err := s.deleteVPC(); err != nil {
		conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, "DeletingFailed", clusterv1.ConditionSeverityWarning, err.Error())
		return err
	}
	conditions.MarkFalse(s.scope.InfraCluster(), infrav1.VpcReadyCondition, clusterv1.DeletedReason, clusterv1.ConditionSeverityInfo, "")

	s.scope.Debug("Delete network completed successfully")
	return nil
}
