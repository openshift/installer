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

package v1alpha7

import (
	"math"
	"reflect"

	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/ptr"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversion"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

var _ ctrlconversion.Convertible = &OpenStackCluster{}

func (r *OpenStackCluster) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackCluster)

	return conversion.ConvertAndRestore(
		r, dst,
		Convert_v1alpha7_OpenStackCluster_To_v1beta1_OpenStackCluster, Convert_v1beta1_OpenStackCluster_To_v1alpha7_OpenStackCluster,
		v1alpha7OpenStackClusterRestorer, v1beta1OpenStackClusterRestorer,
	)
}

func (r *OpenStackCluster) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackCluster)

	return conversion.ConvertAndRestore(
		src, r,
		Convert_v1beta1_OpenStackCluster_To_v1alpha7_OpenStackCluster, Convert_v1alpha7_OpenStackCluster_To_v1beta1_OpenStackCluster,
		v1beta1OpenStackClusterRestorer, v1alpha7OpenStackClusterRestorer,
	)
}

var _ ctrlconversion.Convertible = &OpenStackClusterList{}

func (r *OpenStackClusterList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterList)

	return Convert_v1alpha7_OpenStackClusterList_To_v1beta1_OpenStackClusterList(r, dst, nil)
}

func (r *OpenStackClusterList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterList)

	return Convert_v1beta1_OpenStackClusterList_To_v1alpha7_OpenStackClusterList(src, r, nil)
}

/* Restorers */

var v1alpha7OpenStackClusterRestorer = conversion.RestorerFor[*OpenStackCluster]{
	"bastion": conversion.HashedFieldRestorer(
		func(c *OpenStackCluster) **Bastion {
			return &c.Spec.Bastion
		},
		restorev1alpha7Bastion,
	),
	"spec": conversion.HashedFieldRestorer(
		func(c *OpenStackCluster) *OpenStackClusterSpec {
			return &c.Spec
		},
		restorev1alpha7ClusterSpec,
		// Filter out Bastion, which is restored separately, and
		// ControlPlaneEndpoint, which is written by the cluster controller
		conversion.HashedFilterField[*OpenStackCluster](
			func(s *OpenStackClusterSpec) *OpenStackClusterSpec {
				if s.Bastion != nil || s.ControlPlaneEndpoint != (clusterv1.APIEndpoint{}) {
					f := *s
					f.Bastion = nil
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
		restorev1alpha7ClusterStatus,
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

func restorev1alpha7ClusterSpec(previous *OpenStackClusterSpec, dst *OpenStackClusterSpec) {
	prevBastion := previous.Bastion
	dstBastion := dst.Bastion
	if prevBastion != nil && dstBastion != nil {
		restorev1alpha7MachineSpec(&prevBastion.Instance, &dstBastion.Instance)
	}

	// We only restore DNSNameservers when these were lossly converted when NodeCIDR is empty.
	if len(previous.DNSNameservers) > 0 && dst.NodeCIDR == "" {
		dst.DNSNameservers = previous.DNSNameservers
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
			restorev1alpha7SubnetFilter(&previous.ExternalRouterIPs[i].Subnet, &dst.ExternalRouterIPs[i].Subnet)
		}
	}

	restorev1alpha7SubnetFilter(&previous.Subnet, &dst.Subnet)

	if dst.Router != nil && previous.Router != nil {
		restorev1alpha7RouterFilter(previous.Router, dst.Router)
	}

	restorev1alpha7NetworkFilter(&previous.Network, &dst.Network)

	// APIServerPort is not uint16
	if previous.APIServerPort > math.MaxUint16 {
		dst.APIServerPort = previous.APIServerPort
	}
}

func restorev1beta1ClusterSpec(previous *infrav1.OpenStackClusterSpec, dst *infrav1.OpenStackClusterSpec) {
	if previous == nil || dst == nil {
		return
	}

	// Bastion is restored separately

	if dst.Network == nil {
		dst.Network = previous.Network
	}

	// ExternalNetwork by filter will be been lost in down-conversion
	if previous.ExternalNetwork != nil {
		if dst.ExternalNetwork == nil {
			dst.ExternalNetwork = &infrav1.NetworkParam{}
		}
		dst.ExternalNetwork.Filter = previous.ExternalNetwork.Filter
	}

	dst.DisableExternalNetwork = previous.DisableExternalNetwork

	restorev1beta1RouterParam(previous.Router, dst.Router)
	restorev1beta1NetworkParam(previous.Network, dst.Network)

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
		optional.RestoreString(&previous.APIServerLoadBalancer.Flavor, &dst.APIServerLoadBalancer.Flavor)

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
	optional.RestoreUInt16(&previous.APIServerPort, &dst.APIServerPort)
	optional.RestoreBool(&previous.DisableAPIServerFloatingIP, &dst.DisableAPIServerFloatingIP)
	optional.RestoreBool(&previous.ControlPlaneOmitAvailabilityZone, &dst.ControlPlaneOmitAvailabilityZone)
	optional.RestoreBool(&previous.DisablePortSecurity, &dst.DisablePortSecurity)

	restorev1beta1APIServerLoadBalancer(previous.APIServerLoadBalancer, dst.APIServerLoadBalancer)
}

func Convert_v1alpha7_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(in *OpenStackClusterSpec, out *infrav1.OpenStackClusterSpec, s apiconversion.Scope) error {
	err := autoConvert_v1alpha7_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.Network != (NetworkFilter{}) {
		out.Network = &infrav1.NetworkParam{}
		if err := Convert_v1alpha7_NetworkFilter_To_v1beta1_NetworkParam(&in.Network, out.Network, s); err != nil {
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
		if err := Convert_v1alpha7_SubnetFilter_To_v1beta1_SubnetParam(&in.Subnet, &subnet, s); err != nil {
			return err
		}
		out.Subnets = []infrav1.SubnetParam{subnet}
	}

	// DNSNameservers without NodeCIDR doesn't make sense, so we drop that.
	if in.NodeCIDR != "" {
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
	if err := Convert_v1alpha7_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(&in.APIServerLoadBalancer, apiServerLoadBalancer, s); err != nil {
		return err
	}
	if !apiServerLoadBalancer.IsZero() {
		out.APIServerLoadBalancer = apiServerLoadBalancer
	}

	if in.APIServerPort > 0 && in.APIServerPort < math.MaxUint16 {
		out.APIServerPort = ptr.To(uint16(in.APIServerPort)) //nolint:gosec
	}

	return nil
}

func Convert_v1beta1_LoadBalancer_To_v1alpha7_LoadBalancer(in *infrav1.LoadBalancer, out *LoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1beta1_LoadBalancer_To_v1alpha7_LoadBalancer(in, out, s)
}

func Convert_v1beta1_APIServerLoadBalancer_To_v1alpha7_APIServerLoadBalancer(in *infrav1.APIServerLoadBalancer, out *APIServerLoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1beta1_APIServerLoadBalancer_To_v1alpha7_APIServerLoadBalancer(in, out, s)
}

func Convert_v1alpha7_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(in *APIServerLoadBalancer, out *infrav1.APIServerLoadBalancer, s apiconversion.Scope) error {
	return autoConvert_v1alpha7_APIServerLoadBalancer_To_v1beta1_APIServerLoadBalancer(in, out, s)
}

func Convert_v1beta1_OpenStackClusterSpec_To_v1alpha7_OpenStackClusterSpec(in *infrav1.OpenStackClusterSpec, out *OpenStackClusterSpec, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_OpenStackClusterSpec_To_v1alpha7_OpenStackClusterSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.Network != nil {
		if err := Convert_v1beta1_NetworkParam_To_v1alpha7_NetworkFilter(in.Network, &out.Network, s); err != nil {
			return err
		}
	}

	if in.ExternalNetwork != nil && in.ExternalNetwork.ID != nil {
		out.ExternalNetworkID = *in.ExternalNetwork.ID
	}

	if len(in.Subnets) >= 1 {
		if err := Convert_v1beta1_SubnetParam_To_v1alpha7_SubnetFilter(&in.Subnets[0], &out.Subnet, s); err != nil {
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
	out.IdentityRef = &OpenStackIdentityReference{}
	if err := Convert_v1beta1_OpenStackIdentityReference_To_v1alpha7_OpenStackIdentityReference(&in.IdentityRef, out.IdentityRef, s); err != nil {
		return err
	}

	if in.APIServerLoadBalancer != nil {
		if err := Convert_v1beta1_APIServerLoadBalancer_To_v1alpha7_APIServerLoadBalancer(in.APIServerLoadBalancer, &out.APIServerLoadBalancer, s); err != nil {
			return err
		}
	}

	out.APIServerPort = int(ptr.Deref(in.APIServerPort, 0))

	return nil
}

/* OpenStackClusterStatus */

func restorev1alpha7ClusterStatus(previous *OpenStackClusterStatus, dst *OpenStackClusterStatus) {
	restorev1alpha7SecurityGroup(previous.ControlPlaneSecurityGroup, dst.ControlPlaneSecurityGroup)
	restorev1alpha7SecurityGroup(previous.WorkerSecurityGroup, dst.WorkerSecurityGroup)
	restorev1alpha7SecurityGroup(previous.BastionSecurityGroup, dst.BastionSecurityGroup)
}

func restorev1beta1ClusterStatus(previous *infrav1.OpenStackClusterStatus, dst *infrav1.OpenStackClusterStatus) {
	if previous == nil || dst == nil {
		return
	}

	restorev1beta1BastionStatus(previous.Bastion, dst.Bastion)

	if previous.APIServerLoadBalancer != nil {
		dst.APIServerLoadBalancer = previous.APIServerLoadBalancer
	}
}

func Convert_v1beta1_OpenStackClusterStatus_To_v1alpha7_OpenStackClusterStatus(in *infrav1.OpenStackClusterStatus, out *OpenStackClusterStatus, s apiconversion.Scope) error {
	return autoConvert_v1beta1_OpenStackClusterStatus_To_v1alpha7_OpenStackClusterStatus(in, out, s)
}

/* Bastion */

func restorev1alpha7Bastion(previous **Bastion, dst **Bastion) {
	if previous == nil || dst == nil || *previous == nil || *dst == nil {
		return
	}
	prevMachineSpec := &(*previous).Instance
	dstMachineSpec := &(*dst).Instance
	restorev1alpha7MachineSpec(prevMachineSpec, dstMachineSpec)
	dstMachineSpec.InstanceID = prevMachineSpec.InstanceID
}

func restorev1beta1Bastion(previous **infrav1.Bastion, dst **infrav1.Bastion) {
	if previous == nil || dst == nil || *previous == nil || *dst == nil {
		return
	}

	restorev1beta1MachineSpec((*previous).Spec, (*dst).Spec)
	optional.RestoreString(&(*previous).FloatingIP, &(*dst).FloatingIP)
	optional.RestoreString(&(*previous).AvailabilityZone, &(*dst).AvailabilityZone)

	if (*dst).Enabled != nil && !*(*dst).Enabled {
		(*dst).Enabled = (*previous).Enabled
	}
}

func Convert_v1alpha7_Bastion_To_v1beta1_Bastion(in *Bastion, out *infrav1.Bastion, s apiconversion.Scope) error {
	err := autoConvert_v1alpha7_Bastion_To_v1beta1_Bastion(in, out, s)
	if err != nil {
		return err
	}

	if !reflect.ValueOf(in.Instance).IsZero() {
		out.Spec = &infrav1.OpenStackMachineSpec{}

		err = Convert_v1alpha7_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(&in.Instance, out.Spec, s)
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

func Convert_v1beta1_Bastion_To_v1alpha7_Bastion(in *infrav1.Bastion, out *Bastion, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_Bastion_To_v1alpha7_Bastion(in, out, s)
	if err != nil {
		return err
	}

	if in.Spec != nil {
		err = Convert_v1beta1_OpenStackMachineSpec_To_v1alpha7_OpenStackMachineSpec(in.Spec, &out.Instance, s)
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

func Convert_v1beta1_BastionStatus_To_v1alpha7_BastionStatus(in *infrav1.BastionStatus, out *BastionStatus, s apiconversion.Scope) error {
	// ReferencedResources have no equivalent in v1alpha7
	return autoConvert_v1beta1_BastionStatus_To_v1alpha7_BastionStatus(in, out, s)
}
