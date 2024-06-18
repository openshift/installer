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

package v1alpha6

import (
	"reflect"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversion"
	optional "sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

var _ ctrlconversion.Convertible = &OpenStackCluster{}

func (r *OpenStackCluster) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackCluster)

	return conversion.ConvertAndRestore(
		r, dst,
		Convert_v1alpha6_OpenStackCluster_To_v1beta1_OpenStackCluster, Convert_v1beta1_OpenStackCluster_To_v1alpha6_OpenStackCluster,
		v1alpha6OpenStackClusterRestorer, v1beta1OpenStackClusterRestorer,
	)
}

func (r *OpenStackCluster) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackCluster)

	return conversion.ConvertAndRestore(
		src, r,
		Convert_v1beta1_OpenStackCluster_To_v1alpha6_OpenStackCluster, Convert_v1alpha6_OpenStackCluster_To_v1beta1_OpenStackCluster,
		v1beta1OpenStackClusterRestorer, v1alpha6OpenStackClusterRestorer,
	)
}

var _ ctrlconversion.Convertible = &OpenStackClusterList{}

func (r *OpenStackClusterList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterList)

	return Convert_v1alpha6_OpenStackClusterList_To_v1beta1_OpenStackClusterList(r, dst, nil)
}

func (r *OpenStackClusterList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterList)

	return Convert_v1beta1_OpenStackClusterList_To_v1alpha6_OpenStackClusterList(src, r, nil)
}

/* Restorers */

var v1alpha6OpenStackClusterRestorer = conversion.RestorerFor[*OpenStackCluster]{
	"spec": conversion.HashedFieldRestorer(
		func(c *OpenStackCluster) *OpenStackClusterSpec {
			return &c.Spec
		},
		restorev1alpha6ClusterSpec,
		// Filter out ControlPlaneEndpoint, which is written by the
		// cluster controller
		conversion.HashedFilterField[*OpenStackCluster](
			func(s *OpenStackClusterSpec) *OpenStackClusterSpec {
				if s.ControlPlaneEndpoint != (clusterv1.APIEndpoint{}) {
					f := *s
					f.ControlPlaneEndpoint = clusterv1.APIEndpoint{}
					return &f
				}
				return s
			},
		),
	),
	"status": conversion.HashedFieldRestorer(
		func(c *OpenStackCluster) *OpenStackClusterStatus {
			return &c.Status
		},
		restorev1alpha6ClusterStatus,
	),
}

var v1beta1OpenStackClusterRestorer = conversion.RestorerFor[*infrav1.OpenStackCluster]{
	"bastion": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackCluster) **infrav1.Bastion {
			return &c.Spec.Bastion
		},
		restorev1beta1Bastion,
	),
	"spec": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackCluster) *infrav1.OpenStackClusterSpec {
			return &c.Spec
		},
		restorev1beta1ClusterSpec,
		// Filter out Bastion, which is restored separately
		conversion.HashedFilterField[*infrav1.OpenStackCluster, infrav1.OpenStackClusterSpec](
			func(s *infrav1.OpenStackClusterSpec) *infrav1.OpenStackClusterSpec {
				if s.Bastion != nil {
					f := *s
					f.Bastion = nil
					return &f
				}
				return s
			},
		),
	),
	"status": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackCluster) *infrav1.OpenStackClusterStatus {
			return &c.Status
		},
		restorev1beta1ClusterStatus,
	),
}

/* OpenStackClusterSpec */

func restorev1alpha6ClusterSpec(previous *OpenStackClusterSpec, dst *OpenStackClusterSpec) {
	if previous == nil || dst == nil {
		return
	}

	if len(previous.ExternalRouterIPs) == len(dst.ExternalRouterIPs) {
		for i := range previous.ExternalRouterIPs {
			dstIP := &dst.ExternalRouterIPs[i]
			previousIP := &previous.ExternalRouterIPs[i]

			// Subnet.Filter.ID was overwritten in up-conversion by Subnet.UUID
			dstIP.Subnet.Filter.ID = previousIP.Subnet.Filter.ID

			// If Subnet.UUID was previously unset, we overwrote it with the value of Subnet.Filter.ID
			// Don't unset it again if it doesn't have the previous value of Subnet.Filter.ID, because that means it was genuinely changed
			if previousIP.Subnet.UUID == "" && dstIP.Subnet.UUID == previousIP.Subnet.Filter.ID {
				dstIP.Subnet.UUID = ""
			}
		}
	}

	// We only restore DNSNameservers when these were lossly converted when NodeCIDR is empty.
	if len(previous.DNSNameservers) > 0 && dst.NodeCIDR == "" {
		dst.DNSNameservers = previous.DNSNameservers
	}

	prevBastion := previous.Bastion
	dstBastion := dst.Bastion
	if prevBastion != nil && dstBastion != nil {
		restorev1alpha6MachineSpec(&prevBastion.Instance, &dstBastion.Instance)
		dstBastion.Instance.InstanceID = prevBastion.Instance.InstanceID
	}

	// To avoid lossy conversion, we need to restore AllowAllInClusterTraffic
	// even if ManagedSecurityGroups is set to false
	if previous.AllowAllInClusterTraffic && !previous.ManagedSecurityGroups {
		dst.AllowAllInClusterTraffic = true
	}

	// Conversion to v1beta1 removes the Kind field
	dst.IdentityRef = previous.IdentityRef

	if len(dst.ExternalRouterIPs) == len(previous.ExternalRouterIPs) {
		for i := range dst.ExternalRouterIPs {
			restorev1alpha6SubnetParam(&previous.ExternalRouterIPs[i].Subnet, &dst.ExternalRouterIPs[i].Subnet)
		}
	}

	restorev1alpha6SubnetFilter(&previous.Subnet, &dst.Subnet)

	restorev1alpha6NetworkFilter(&previous.Network, &dst.Network)
}

func restorev1beta1ClusterSpec(previous *infrav1.OpenStackClusterSpec, dst *infrav1.OpenStackClusterSpec) {
	// Bastion is restored separately

	restorev1beta1NetworkParam(previous.Network, dst.Network)

	// ExternalNetwork by filter will be been lost in down-conversion
	if previous.ExternalNetwork != nil {
		if dst.ExternalNetwork == nil {
			dst.ExternalNetwork = &infrav1.NetworkParam{}
		}
		dst.ExternalNetwork.Filter = previous.ExternalNetwork.Filter
	}

	// Restore fields not present in v1alpha6
	dst.Router = previous.Router
	dst.NetworkMTU = previous.NetworkMTU
	dst.DisableExternalNetwork = previous.DisableExternalNetwork

	if len(previous.Subnets) > 0 && len(dst.Subnets) > 0 {
		restorev1beta1SubnetParam(&previous.Subnets[0], &dst.Subnets[0])
	}
	if len(previous.Subnets) > 1 {
		dst.Subnets = append(dst.Subnets, previous.Subnets[1:]...)
	}

	if len(previous.ExternalRouterIPs) == len(dst.ExternalRouterIPs) {
		for i := range dst.ExternalRouterIPs {
			restorev1beta1SubnetParam(&previous.ExternalRouterIPs[i].Subnet, &dst.ExternalRouterIPs[i].Subnet)
		}
	}

	dst.ManagedSubnets = previous.ManagedSubnets

	if previous.ManagedSecurityGroups != nil && dst.ManagedSecurityGroups != nil {
		dst.ManagedSecurityGroups.AllNodesSecurityGroupRules = previous.ManagedSecurityGroups.AllNodesSecurityGroupRules
	}

	if dst.APIServerLoadBalancer != nil && previous.APIServerLoadBalancer != nil {
		if dst.APIServerLoadBalancer.Enabled == nil || !*dst.APIServerLoadBalancer.Enabled {
			dst.APIServerLoadBalancer.Enabled = previous.APIServerLoadBalancer.Enabled
		}
		optional.RestoreString(&previous.APIServerLoadBalancer.Provider, &dst.APIServerLoadBalancer.Provider)

		if previous.APIServerLoadBalancer.Network != nil {
			dst.APIServerLoadBalancer.Network = previous.APIServerLoadBalancer.Network
		}
		if previous.APIServerLoadBalancer.Subnets != nil {
			dst.APIServerLoadBalancer.Subnets = previous.APIServerLoadBalancer.Subnets
		}
	}
	if dst.APIServerLoadBalancer.IsZero() {
		dst.APIServerLoadBalancer = previous.APIServerLoadBalancer
	}

	if dst.ControlPlaneEndpoint == nil || *dst.ControlPlaneEndpoint == (clusterv1.APIEndpoint{}) {
		dst.ControlPlaneEndpoint = previous.ControlPlaneEndpoint
	}

	optional.RestoreString(&previous.APIServerFloatingIP, &dst.APIServerFloatingIP)
	optional.RestoreString(&previous.APIServerFixedIP, &dst.APIServerFixedIP)
	optional.RestoreInt(&previous.APIServerPort, &dst.APIServerPort)
	optional.RestoreBool(&previous.DisableAPIServerFloatingIP, &dst.DisableAPIServerFloatingIP)
	optional.RestoreBool(&previous.ControlPlaneOmitAvailabilityZone, &dst.ControlPlaneOmitAvailabilityZone)
	optional.RestoreBool(&previous.DisablePortSecurity, &dst.DisablePortSecurity)

	restorev1beta1APIServerLoadBalancer(previous.APIServerLoadBalancer, dst.APIServerLoadBalancer)
}

func Convert_v1alpha6_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(in *OpenStackClusterSpec, out *infrav1.OpenStackClusterSpec, s apiconversion.Scope) error {
	err := autoConvert_v1alpha6_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.Network != (NetworkFilter{}) {
		out.Network = &infrav1.NetworkParam{}
		if err := Convert_v1alpha6_NetworkFilter_To_v1beta1_NetworkParam(&in.Network, out.Network, s); err != nil {
			return err
		}
	}

	if in.ExternalNetworkID != "" {
		out.ExternalNetwork = &infrav1.NetworkParam{
			ID: &in.ExternalNetworkID,
		}
	}

	emptySubnet := SubnetFilter{}
	if in.Subnet != emptySubnet {
		subnet := infrav1.SubnetParam{}
		if err := Convert_v1alpha6_SubnetFilter_To_v1beta1_SubnetParam(&in.Subnet, &subnet, s); err != nil {
			return err
		}
		out.Subnets = []infrav1.SubnetParam{subnet}
	}

	// DNSNameservers without NodeCIDR doesn't make sense, so we drop that.
	if len(in.NodeCIDR) > 0 {
		out.ManagedSubnets = []infrav1.SubnetSpec{
			{
				CIDR:           in.NodeCIDR,
				DNSNameservers: in.DNSNameservers,
			},
		}
	}

	if in.ManagedSecurityGroups {
		out.ManagedSecurityGroups = &infrav1.ManagedSecurityGroups{}
		if !in.AllowAllInClusterTraffic {
			out.ManagedSecurityGroups.AllNodesSecurityGroupRules = infrav1.LegacyCalicoSecurityGroupRules()
		} else {
			out.ManagedSecurityGroups.AllowAllInClusterTraffic = true
		}
	}

	if in.ControlPlaneEndpoint != (clusterv1.APIEndpoint{}) {
		out.ControlPlaneEndpoint = &in.ControlPlaneEndpoint
	}

	out.IdentityRef.CloudName = in.CloudName
	if in.IdentityRef != nil {
		out.IdentityRef.Name = in.IdentityRef.Name
	}

	apiServerLoadBalancer := &infrav1.APIServerLoadBalancer{}
	if err := Convert_v1alpha6_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(&in.APIServerLoadBalancer, apiServerLoadBalancer, s); err != nil {
		return err
	}
	if !apiServerLoadBalancer.IsZero() {
		out.APIServerLoadBalancer = apiServerLoadBalancer
	}

	return nil
}

func Convert_v1beta1_APIServerLoadBalancer_To_v1alpha6_APIServerLoadBalancer(in *infrav1.APIServerLoadBalancer, out *APIServerLoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1beta1_APIServerLoadBalancer_To_v1alpha6_APIServerLoadBalancer(in, out, s)
}

func Convert_v1beta1_LoadBalancer_To_v1alpha6_LoadBalancer(in *infrav1.LoadBalancer, out *LoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1beta1_LoadBalancer_To_v1alpha6_LoadBalancer(in, out, s)
}

func Convert_v1alpha6_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(in *APIServerLoadBalancer, out *infrav1.APIServerLoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1alpha6_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(in, out, s)
}

func Convert_v1beta1_OpenStackClusterSpec_To_v1alpha6_OpenStackClusterSpec(in *infrav1.OpenStackClusterSpec, out *OpenStackClusterSpec, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_OpenStackClusterSpec_To_v1alpha6_OpenStackClusterSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.Network != nil {
		if err := Convert_v1beta1_NetworkParam_To_v1alpha6_NetworkFilter(in.Network, &out.Network, s); err != nil {
			return err
		}
	}

	if in.ExternalNetwork != nil && in.ExternalNetwork.ID != nil {
		out.ExternalNetworkID = *in.ExternalNetwork.ID
	}

	if len(in.Subnets) >= 1 {
		if err := Convert_v1beta1_SubnetParam_To_v1alpha6_SubnetFilter(&in.Subnets[0], &out.Subnet, s); err != nil {
			return err
		}
	}

	if len(in.ManagedSubnets) > 0 {
		out.NodeCIDR = in.ManagedSubnets[0].CIDR
		out.DNSNameservers = in.ManagedSubnets[0].DNSNameservers
	}

	if in.ManagedSecurityGroups != nil {
		out.ManagedSecurityGroups = true
		out.AllowAllInClusterTraffic = in.ManagedSecurityGroups.AllowAllInClusterTraffic
	}

	if in.ControlPlaneEndpoint != nil {
		out.ControlPlaneEndpoint = *in.ControlPlaneEndpoint
	}

	out.CloudName = in.IdentityRef.CloudName
	out.IdentityRef = &OpenStackIdentityReference{Name: in.IdentityRef.Name}

	if in.APIServerLoadBalancer != nil {
		if err := Convert_v1beta1_APIServerLoadBalancer_To_v1alpha6_APIServerLoadBalancer(in.APIServerLoadBalancer, &out.APIServerLoadBalancer, s); err != nil {
			return err
		}
	}

	return nil
}

/* OpenStackClusterStatus */

func restorev1alpha6ClusterStatus(previous *OpenStackClusterStatus, dst *OpenStackClusterStatus) {
	// PortOpts.SecurityGroups have been removed in v1beta1
	// We restore the whole PortOpts/Networks since they are anyway immutable.
	if previous.ExternalNetwork != nil && dst.ExternalNetwork != nil {
		dst.ExternalNetwork.PortOpts = previous.ExternalNetwork.PortOpts
	}
	if previous.Network != nil {
		dst.Network = previous.Network
	}
	if previous.Bastion != nil && previous.Bastion.Networks != nil && dst.Bastion != nil {
		dst.Bastion.Networks = previous.Bastion.Networks
	}

	restorev1alpha6SecurityGroup(previous.ControlPlaneSecurityGroup, dst.ControlPlaneSecurityGroup)
	restorev1alpha6SecurityGroup(previous.WorkerSecurityGroup, dst.WorkerSecurityGroup)
	restorev1alpha6SecurityGroup(previous.BastionSecurityGroup, dst.BastionSecurityGroup)
}

func restorev1beta1ClusterStatus(previous *infrav1.OpenStackClusterStatus, dst *infrav1.OpenStackClusterStatus) {
	// It's (theoretically) possible in v1beta1 to have Network nil but
	// Router or APIServerLoadBalancer not nil. In hub-spoke-hub conversion this will
	// result in Network being a pointer to an empty object.
	if previous.Network == nil && dst.Network != nil && reflect.ValueOf(*dst.Network).IsZero() {
		dst.Network = nil
	}

	dst.ControlPlaneSecurityGroup = previous.ControlPlaneSecurityGroup
	dst.WorkerSecurityGroup = previous.WorkerSecurityGroup
	dst.BastionSecurityGroup = previous.BastionSecurityGroup

	restorev1beta1BastionStatus(previous.Bastion, dst.Bastion)

	if previous.APIServerLoadBalancer != nil {
		dst.APIServerLoadBalancer = previous.APIServerLoadBalancer
	}
}

func Convert_v1beta1_OpenStackClusterStatus_To_v1alpha6_OpenStackClusterStatus(in *infrav1.OpenStackClusterStatus, out *OpenStackClusterStatus, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_OpenStackClusterStatus_To_v1alpha6_OpenStackClusterStatus(in, out, s)
	if err != nil {
		return err
	}

	// Router and APIServerLoadBalancer have been moved out of Network in v1beta1
	if in.Router != nil || in.APIServerLoadBalancer != nil {
		if out.Network == nil {
			out.Network = &Network{}
		}

		out.Network.Router = (*Router)(in.Router)
		if in.APIServerLoadBalancer != nil {
			out.Network.APIServerLoadBalancer = &LoadBalancer{
				Name:         in.APIServerLoadBalancer.Name,
				ID:           in.APIServerLoadBalancer.ID,
				IP:           in.APIServerLoadBalancer.IP,
				InternalIP:   in.APIServerLoadBalancer.InternalIP,
				AllowedCIDRs: in.APIServerLoadBalancer.AllowedCIDRs,
				Tags:         in.APIServerLoadBalancer.Tags,
			}
		}
	}

	return nil
}

func Convert_v1alpha6_OpenStackClusterStatus_To_v1beta1_OpenStackClusterStatus(in *OpenStackClusterStatus, out *infrav1.OpenStackClusterStatus, s apiconversion.Scope) error {
	err := autoConvert_v1alpha6_OpenStackClusterStatus_To_v1beta1_OpenStackClusterStatus(in, out, s)
	if err != nil {
		return err
	}

	// Router and APIServerLoadBalancer have been moved out of Network in v1beta1
	if in.Network != nil {
		out.Router = (*infrav1.Router)(in.Network.Router)
		if in.Network.APIServerLoadBalancer != nil {
			out.APIServerLoadBalancer = &infrav1.LoadBalancer{
				Name:         in.Network.APIServerLoadBalancer.Name,
				ID:           in.Network.APIServerLoadBalancer.ID,
				IP:           in.Network.APIServerLoadBalancer.IP,
				InternalIP:   in.Network.APIServerLoadBalancer.InternalIP,
				AllowedCIDRs: in.Network.APIServerLoadBalancer.AllowedCIDRs,
				Tags:         in.Network.APIServerLoadBalancer.Tags,
			}
		}
	}

	return nil
}

/* Bastion */

func restorev1beta1Bastion(previous **infrav1.Bastion, dst **infrav1.Bastion) {
	if previous == nil || dst == nil || *previous == nil || *dst == nil {
		return
	}
	if *dst != nil && (*previous).Spec != nil && (*dst).Spec != nil {
		restorev1beta1MachineSpec((*previous).Spec, (*dst).Spec)
	}

	optional.RestoreString(&(*previous).FloatingIP, &(*dst).FloatingIP)
	optional.RestoreString(&(*previous).AvailabilityZone, &(*dst).AvailabilityZone)

	if (*dst).Enabled != nil && !*(*dst).Enabled {
		(*dst).Enabled = (*previous).Enabled
	}
}

func Convert_v1alpha6_Bastion_To_v1beta1_Bastion(in *Bastion, out *infrav1.Bastion, s apiconversion.Scope) error {
	err := autoConvert_v1alpha6_Bastion_To_v1beta1_Bastion(in, out, s)
	if err != nil {
		return err
	}

	if !reflect.ValueOf(in.Instance).IsZero() {
		out.Spec = &infrav1.OpenStackMachineSpec{}

		err = Convert_v1alpha6_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(&in.Instance, out.Spec, s)
		if err != nil {
			return err
		}

		if in.Instance.ServerGroupID != "" {
			out.Spec.ServerGroup = &infrav1.ServerGroupParam{ID: &in.Instance.ServerGroupID}
		} else {
			out.Spec.ServerGroup = nil
		}

		err = optional.Convert_string_To_optional_String(&in.Instance.FloatingIP, &out.FloatingIP, s)
		if err != nil {
			return err
		}
	}

	// nil the Spec if it's basically an empty object.
	if out.Spec != nil && reflect.ValueOf(*out.Spec).IsZero() {
		out.Spec = nil
	}
	return nil
}

func Convert_v1beta1_Bastion_To_v1alpha6_Bastion(in *infrav1.Bastion, out *Bastion, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_Bastion_To_v1alpha6_Bastion(in, out, s)
	if err != nil {
		return err
	}

	if in.Spec != nil {
		err = Convert_v1beta1_OpenStackMachineSpec_To_v1alpha6_OpenStackMachineSpec(in.Spec, &out.Instance, s)
		if err != nil {
			return err
		}

		if in.Spec.ServerGroup != nil && in.Spec.ServerGroup.ID != nil {
			out.Instance.ServerGroupID = *in.Spec.ServerGroup.ID
		}
	}

	return optional.Convert_optional_String_To_string(&in.FloatingIP, &out.Instance.FloatingIP, s)
}

/* Bastion status */

func restorev1beta1BastionStatus(previous *infrav1.BastionStatus, dst *infrav1.BastionStatus) {
	if previous == nil || dst == nil {
		return
	}

	// Resolved and resources have no equivalents
	dst.Resolved = previous.Resolved
	dst.Resources = previous.Resources
}

func Convert_v1alpha6_Instance_To_v1beta1_BastionStatus(in *Instance, out *infrav1.BastionStatus, _ apiconversion.Scope) error {
	// BastionStatus is the same as Spec with unused fields removed
	out.ID = in.ID
	out.Name = in.Name
	out.SSHKeyName = in.SSHKeyName
	out.State = infrav1.InstanceState(in.State)
	out.IP = in.IP
	out.FloatingIP = in.FloatingIP
	return nil
}

func Convert_v1beta1_BastionStatus_To_v1alpha6_Instance(in *infrav1.BastionStatus, out *Instance, _ apiconversion.Scope) error {
	// BastionStatus is the same as Spec with unused fields removed
	out.ID = in.ID
	out.Name = in.Name
	out.SSHKeyName = in.SSHKeyName
	out.State = InstanceState(in.State)
	out.IP = in.IP
	out.FloatingIP = in.FloatingIP
	return nil
}
