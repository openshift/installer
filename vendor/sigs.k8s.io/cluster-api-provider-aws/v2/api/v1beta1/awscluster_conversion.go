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

package v1beta1

import (
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	infrav2 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1beta1 AWSCluster receiver to a v1beta2 AWSCluster.
func (src *AWSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav2.AWSCluster)

	if err := Convert_v1beta1_AWSCluster_To_v1beta2_AWSCluster(src, dst, nil); err != nil {
		return err
	}
	// Manually restore data.
	restored := &infrav2.AWSCluster{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	if restored.Spec.ControlPlaneLoadBalancer != nil {
		if dst.Spec.ControlPlaneLoadBalancer == nil {
			dst.Spec.ControlPlaneLoadBalancer = &infrav2.AWSLoadBalancerSpec{}
		}
		restoreControlPlaneLoadBalancer(restored.Spec.ControlPlaneLoadBalancer, dst.Spec.ControlPlaneLoadBalancer)
	}
	restoreControlPlaneLoadBalancerStatus(&restored.Status.Network.APIServerELB, &dst.Status.Network.APIServerELB)

	if restored.Spec.SecondaryControlPlaneLoadBalancer != nil {
		if dst.Spec.SecondaryControlPlaneLoadBalancer == nil {
			dst.Spec.SecondaryControlPlaneLoadBalancer = &infrav2.AWSLoadBalancerSpec{}
		}
		restoreControlPlaneLoadBalancer(restored.Spec.SecondaryControlPlaneLoadBalancer, dst.Spec.SecondaryControlPlaneLoadBalancer)
	}
	restoreControlPlaneLoadBalancerStatus(&restored.Status.Network.SecondaryAPIServerELB, &dst.Status.Network.SecondaryAPIServerELB)

	dst.Spec.S3Bucket = restored.Spec.S3Bucket
	if restored.Status.Bastion != nil {
		dst.Status.Bastion.InstanceMetadataOptions = restored.Status.Bastion.InstanceMetadataOptions
		dst.Status.Bastion.PlacementGroupName = restored.Status.Bastion.PlacementGroupName
		dst.Status.Bastion.PlacementGroupPartition = restored.Status.Bastion.PlacementGroupPartition
		dst.Status.Bastion.PrivateDNSName = restored.Status.Bastion.PrivateDNSName
		dst.Status.Bastion.PublicIPOnLaunch = restored.Status.Bastion.PublicIPOnLaunch
		dst.Status.Bastion.NetworkInterfaceType = restored.Status.Bastion.NetworkInterfaceType
		dst.Status.Bastion.CapacityReservationID = restored.Status.Bastion.CapacityReservationID
		dst.Status.Bastion.MarketType = restored.Status.Bastion.MarketType
	}
	dst.Spec.Partition = restored.Spec.Partition

	for role, sg := range restored.Status.Network.SecurityGroups {
		dst.Status.Network.SecurityGroups[role] = sg
	}
	dst.Status.Network.NatGatewaysIPs = restored.Status.Network.NatGatewaysIPs

	if restored.Spec.NetworkSpec.VPC.IPAMPool != nil {
		if dst.Spec.NetworkSpec.VPC.IPAMPool == nil {
			dst.Spec.NetworkSpec.VPC.IPAMPool = &infrav2.IPAMPool{}
		}

		restoreIPAMPool(restored.Spec.NetworkSpec.VPC.IPAMPool, dst.Spec.NetworkSpec.VPC.IPAMPool)
	}

	if restored.Spec.NetworkSpec.VPC.IsIPv6Enabled() && restored.Spec.NetworkSpec.VPC.IPv6.IPAMPool != nil {
		if dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool == nil {
			dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool = &infrav2.IPAMPool{}
		}

		restoreIPAMPool(restored.Spec.NetworkSpec.VPC.IPv6.IPAMPool, dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool)
	}

	dst.Spec.NetworkSpec.AdditionalControlPlaneIngressRules = restored.Spec.NetworkSpec.AdditionalControlPlaneIngressRules
	dst.Spec.NetworkSpec.NodePortIngressRuleCidrBlocks = restored.Spec.NetworkSpec.NodePortIngressRuleCidrBlocks

	if restored.Spec.NetworkSpec.VPC.IPAMPool != nil {
		if dst.Spec.NetworkSpec.VPC.IPAMPool == nil {
			dst.Spec.NetworkSpec.VPC.IPAMPool = &infrav2.IPAMPool{}
		}

		restoreIPAMPool(restored.Spec.NetworkSpec.VPC.IPAMPool, dst.Spec.NetworkSpec.VPC.IPAMPool)
	}

	if restored.Spec.NetworkSpec.VPC.IsIPv6Enabled() && restored.Spec.NetworkSpec.VPC.IPv6.IPAMPool != nil {
		if dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool == nil {
			dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool = &infrav2.IPAMPool{}
		}

		restoreIPAMPool(restored.Spec.NetworkSpec.VPC.IPv6.IPAMPool, dst.Spec.NetworkSpec.VPC.IPv6.IPAMPool)
	}

	dst.Spec.NetworkSpec.VPC.EmptyRoutesDefaultVPCSecurityGroup = restored.Spec.NetworkSpec.VPC.EmptyRoutesDefaultVPCSecurityGroup
	dst.Spec.NetworkSpec.VPC.PrivateDNSHostnameTypeOnLaunch = restored.Spec.NetworkSpec.VPC.PrivateDNSHostnameTypeOnLaunch
	dst.Spec.NetworkSpec.VPC.CarrierGatewayID = restored.Spec.NetworkSpec.VPC.CarrierGatewayID
	dst.Spec.NetworkSpec.VPC.SubnetSchema = restored.Spec.NetworkSpec.VPC.SubnetSchema
	dst.Spec.NetworkSpec.VPC.SecondaryCidrBlocks = restored.Spec.NetworkSpec.VPC.SecondaryCidrBlocks

	if restored.Spec.NetworkSpec.VPC.ElasticIPPool != nil {
		if dst.Spec.NetworkSpec.VPC.ElasticIPPool == nil {
			dst.Spec.NetworkSpec.VPC.ElasticIPPool = &infrav2.ElasticIPPool{}
		}
		if restored.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4Pool != nil {
			dst.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4Pool = restored.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4Pool
		}
		if restored.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4PoolFallBackOrder != nil {
			dst.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4PoolFallBackOrder = restored.Spec.NetworkSpec.VPC.ElasticIPPool.PublicIpv4PoolFallBackOrder
		}
	}

	// Restore SubnetSpec.ResourceID, SubnetSpec.ParentZoneName, and SubnetSpec.ZoneType fields, if any.
	for _, subnet := range restored.Spec.NetworkSpec.Subnets {
		for i, dstSubnet := range dst.Spec.NetworkSpec.Subnets {
			if dstSubnet.ID == subnet.ID {
				if len(subnet.ResourceID) > 0 {
					dstSubnet.ResourceID = subnet.ResourceID
				}
				if subnet.ParentZoneName != nil {
					dstSubnet.ParentZoneName = subnet.ParentZoneName
				}
				if subnet.ZoneType != nil {
					dstSubnet.ZoneType = subnet.ZoneType
				}
				dstSubnet.DeepCopyInto(&dst.Spec.NetworkSpec.Subnets[i])
			}
		}
	}

	return nil
}

// restoreControlPlaneLoadBalancerStatus manually restores the control plane loadbalancer status data.
// Assumes restored and dst are non-nil.
func restoreControlPlaneLoadBalancerStatus(restored, dst *infrav2.LoadBalancer) {
	dst.ARN = restored.ARN
	dst.LoadBalancerType = restored.LoadBalancerType
	dst.ELBAttributes = restored.ELBAttributes
	dst.ELBListeners = restored.ELBListeners
	dst.Name = restored.Name
	dst.DNSName = restored.DNSName
	dst.Scheme = restored.Scheme
	dst.SubnetIDs = restored.SubnetIDs
	dst.SecurityGroupIDs = restored.SecurityGroupIDs
	dst.HealthCheck = restored.HealthCheck
	dst.ClassicElbAttributes = restored.ClassicElbAttributes
	dst.Tags = restored.Tags
	dst.ClassicELBListeners = restored.ClassicELBListeners
	dst.AvailabilityZones = restored.AvailabilityZones
}

// restoreIPAMPool manually restores the ipam pool data.
// Assumes restored and dst are non-nil.
func restoreIPAMPool(restored, dst *infrav2.IPAMPool) {
	dst.ID = restored.ID
	dst.Name = restored.Name
	dst.NetmaskLength = restored.NetmaskLength
}

// restoreControlPlaneLoadBalancer manually restores the control plane loadbalancer data.
// Assumes restored and dst are non-nil.
func restoreControlPlaneLoadBalancer(restored, dst *infrav2.AWSLoadBalancerSpec) {
	dst.Name = restored.Name
	dst.HealthCheckProtocol = restored.HealthCheckProtocol
	dst.HealthCheck = restored.HealthCheck
	dst.LoadBalancerType = restored.LoadBalancerType
	dst.DisableHostsRewrite = restored.DisableHostsRewrite
	dst.PreserveClientIP = restored.PreserveClientIP
	dst.IngressRules = restored.IngressRules
	dst.AdditionalListeners = restored.AdditionalListeners
	dst.AdditionalSecurityGroups = restored.AdditionalSecurityGroups
	dst.Scheme = restored.Scheme
	dst.CrossZoneLoadBalancing = restored.CrossZoneLoadBalancing
	dst.Subnets = restored.Subnets
}

// ConvertFrom converts the v1beta1 AWSCluster receiver to a v1beta1 AWSCluster.
func (r *AWSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav2.AWSCluster)

	if err := Convert_v1beta2_AWSCluster_To_v1beta1_AWSCluster(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 AWSClusterList receiver to a v1beta2 AWSClusterList.
func (src *AWSClusterList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav2.AWSClusterList)

	return Convert_v1beta1_AWSClusterList_To_v1beta2_AWSClusterList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSClusterList receiver to a v1beta1 AWSClusterList.
func (r *AWSClusterList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav2.AWSClusterList)

	return Convert_v1beta2_AWSClusterList_To_v1beta1_AWSClusterList(src, r, nil)
}

func Convert_v1beta2_SubnetSpec_To_v1beta1_SubnetSpec(in *infrav2.SubnetSpec, out *SubnetSpec, s apiconversion.Scope) error {
	return autoConvert_v1beta2_SubnetSpec_To_v1beta1_SubnetSpec(in, out, s)
}
