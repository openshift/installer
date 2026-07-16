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

import (
	"sort"
	"strings"
	"unsafe"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/ptr"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta2"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/optional"
)

const securityGroupRuleDirectionIngress = "ingress"

// ConvertTo converts this OpenStackCluster to the Hub version (v1beta2).
func (src *OpenStackCluster) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackCluster)
	if err := Convert_v1beta1_OpenStackCluster_To_v1beta2_OpenStackCluster(src, dst, nil); err != nil {
		return err
	}
	for i := range dst.Status.Conditions {
		dst.Status.Conditions[i].ObservedGeneration = src.Generation
	}
	return utilconversion.MarshalData(src, dst)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackCluster) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackCluster)
	if err := Convert_v1beta2_OpenStackCluster_To_v1beta1_OpenStackCluster(src, dst, nil); err != nil {
		return err
	}
	_, err := utilconversion.UnmarshalData(src, dst)
	return err
}

// ConvertTo converts this OpenStackMachine to the Hub version (v1beta2).
func (src *OpenStackMachine) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachine)
	if err := Convert_v1beta1_OpenStackMachine_To_v1beta2_OpenStackMachine(src, dst, nil); err != nil {
		return err
	}
	for i := range dst.Status.Conditions {
		dst.Status.Conditions[i].ObservedGeneration = src.Generation
	}
	return utilconversion.MarshalData(src, dst)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackMachine) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachine)
	if err := Convert_v1beta2_OpenStackMachine_To_v1beta1_OpenStackMachine(src, dst, nil); err != nil {
		return err
	}
	_, err := utilconversion.UnmarshalData(src, dst)
	return err
}

// ConvertTo converts this OpenStackClusterTemplate to the Hub version (v1beta2).
func (src *OpenStackClusterTemplate) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterTemplate)
	if err := Convert_v1beta1_OpenStackClusterTemplate_To_v1beta2_OpenStackClusterTemplate(src, dst, nil); err != nil {
		return err
	}
	return utilconversion.MarshalData(src, dst)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackClusterTemplate) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterTemplate)
	if err := Convert_v1beta2_OpenStackClusterTemplate_To_v1beta1_OpenStackClusterTemplate(src, dst, nil); err != nil {
		return err
	}
	_, err := utilconversion.UnmarshalData(src, dst)
	return err
}

// ConvertTo converts this OpenStackMachineTemplate to the Hub version (v1beta2).
func (src *OpenStackMachineTemplate) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineTemplate)
	if err := Convert_v1beta1_OpenStackMachineTemplate_To_v1beta2_OpenStackMachineTemplate(src, dst, nil); err != nil {
		return err
	}
	return utilconversion.MarshalData(src, dst)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackMachineTemplate) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineTemplate)
	if err := Convert_v1beta2_OpenStackMachineTemplate_To_v1beta1_OpenStackMachineTemplate(src, dst, nil); err != nil {
		return err
	}
	_, err := utilconversion.UnmarshalData(src, dst)
	return err
}

// ConvertTo converts this OpenStackClusterList to the Hub version (v1beta2).
func (src *OpenStackClusterList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterList)
	return Convert_v1beta1_OpenStackClusterList_To_v1beta2_OpenStackClusterList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackClusterList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterList)
	return Convert_v1beta2_OpenStackClusterList_To_v1beta1_OpenStackClusterList(src, dst, nil)
}

// ConvertTo converts this OpenStackMachineList to the Hub version (v1beta2).
func (src *OpenStackMachineList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineList)
	return Convert_v1beta1_OpenStackMachineList_To_v1beta2_OpenStackMachineList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackMachineList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineList)
	return Convert_v1beta2_OpenStackMachineList_To_v1beta1_OpenStackMachineList(src, dst, nil)
}

// ConvertTo converts this OpenStackClusterTemplateList to the Hub version (v1beta2).
func (src *OpenStackClusterTemplateList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackClusterTemplateList)
	return Convert_v1beta1_OpenStackClusterTemplateList_To_v1beta2_OpenStackClusterTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackClusterTemplateList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackClusterTemplateList)
	return Convert_v1beta2_OpenStackClusterTemplateList_To_v1beta1_OpenStackClusterTemplateList(src, dst, nil)
}

// ConvertTo converts this OpenStackMachineTemplateList to the Hub version (v1beta2).
func (src *OpenStackMachineTemplateList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineTemplateList)
	return Convert_v1beta1_OpenStackMachineTemplateList_To_v1beta2_OpenStackMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts from the Hub version (v1beta2) to this version.
//
//nolint:revive // dst is the receiver here (converting FROM hub TO spoke)
func (dst *OpenStackMachineTemplateList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineTemplateList)
	return Convert_v1beta2_OpenStackMachineTemplateList_To_v1beta1_OpenStackMachineTemplateList(src, dst, nil)
}

// Manual conversion functions for Status types that conversion-gen cannot
// auto-generate due to FailureDomains (map↔slice), Conditions (CAPI↔metav1),
// and deprecated fields (Ready, FailureReason, FailureMessage).

func Convert_v1beta1_OpenStackClusterStatus_To_v1beta2_OpenStackClusterStatus(in *OpenStackClusterStatus, out *infrav1.OpenStackClusterStatus, _ apiconversion.Scope) error {
	out.Initialization = (*infrav1.ClusterInitialization)(unsafe.Pointer(in.Initialization))
	out.Network = (*infrav1.NetworkStatusWithSubnets)(unsafe.Pointer(in.Network))
	out.ExternalNetwork = (*infrav1.NetworkStatus)(unsafe.Pointer(in.ExternalNetwork))
	out.Router = (*infrav1.Router)(unsafe.Pointer(in.Router))
	out.APIServerManagedLoadBalancer = (*infrav1.LoadBalancer)(unsafe.Pointer(in.APIServerLoadBalancer))
	out.ControlPlaneSecurityGroup = (*infrav1.SecurityGroupStatus)(unsafe.Pointer(in.ControlPlaneSecurityGroup))
	out.WorkerSecurityGroup = (*infrav1.SecurityGroupStatus)(unsafe.Pointer(in.WorkerSecurityGroup))
	out.BastionSecurityGroup = (*infrav1.SecurityGroupStatus)(unsafe.Pointer(in.BastionSecurityGroup))
	out.Bastion = (*infrav1.BastionStatus)(unsafe.Pointer(in.Bastion))

	if len(in.FailureDomains) > 0 {
		out.FailureDomains = make([]clusterv1.FailureDomain, 0, len(in.FailureDomains))
		for name, fd := range in.FailureDomains {
			out.FailureDomains = append(out.FailureDomains, clusterv1.FailureDomain{
				Name:         name,
				ControlPlane: ptr.To(fd.ControlPlane),
				Attributes:   fd.Attributes,
			})
		}
		sort.Slice(out.FailureDomains, func(i, j int) bool {
			return out.FailureDomains[i].Name < out.FailureDomains[j].Name
		})
	}

	out.Conditions = infrav1.ConvertConditionsToV1Beta2(in.Conditions, 0)

	return nil
}

func Convert_v1beta2_OpenStackClusterStatus_To_v1beta1_OpenStackClusterStatus(in *infrav1.OpenStackClusterStatus, out *OpenStackClusterStatus, _ apiconversion.Scope) error {
	out.Initialization = (*ClusterInitialization)(unsafe.Pointer(in.Initialization))
	out.Network = (*NetworkStatusWithSubnets)(unsafe.Pointer(in.Network))
	out.ExternalNetwork = (*NetworkStatus)(unsafe.Pointer(in.ExternalNetwork))
	out.Router = (*Router)(unsafe.Pointer(in.Router))
	out.APIServerLoadBalancer = (*LoadBalancer)(unsafe.Pointer(in.APIServerManagedLoadBalancer))
	out.ControlPlaneSecurityGroup = (*SecurityGroupStatus)(unsafe.Pointer(in.ControlPlaneSecurityGroup))
	out.WorkerSecurityGroup = (*SecurityGroupStatus)(unsafe.Pointer(in.WorkerSecurityGroup))
	out.BastionSecurityGroup = (*SecurityGroupStatus)(unsafe.Pointer(in.BastionSecurityGroup))
	out.Bastion = (*BastionStatus)(unsafe.Pointer(in.Bastion))

	if len(in.FailureDomains) > 0 {
		out.FailureDomains = make(clusterv1beta1.FailureDomains, len(in.FailureDomains))
		for _, fd := range in.FailureDomains {
			out.FailureDomains[fd.Name] = clusterv1beta1.FailureDomainSpec{
				ControlPlane: ptr.Deref(fd.ControlPlane, false),
				Attributes:   fd.Attributes,
			}
		}
	}

	out.Conditions = infrav1.ConvertConditionsFromV1Beta2(in.Conditions)
	out.Ready = infrav1.IsReady(in.Conditions)

	return nil
}

func Convert_v1beta1_OpenStackMachineStatus_To_v1beta2_OpenStackMachineStatus(in *OpenStackMachineStatus, out *infrav1.OpenStackMachineStatus, _ apiconversion.Scope) error {
	out.Initialization = (*infrav1.MachineInitialization)(unsafe.Pointer(in.Initialization))
	out.InstanceID = in.InstanceID
	out.Addresses = in.Addresses
	out.InstanceState = (*infrav1.InstanceState)(unsafe.Pointer(in.InstanceState))
	out.Resolved = (*infrav1.ResolvedMachineSpec)(unsafe.Pointer(in.Resolved))
	out.Resources = (*infrav1.MachineResources)(unsafe.Pointer(in.Resources))

	out.Conditions = infrav1.ConvertConditionsToV1Beta2(in.Conditions, 0)

	return nil
}

func Convert_v1beta2_OpenStackMachineStatus_To_v1beta1_OpenStackMachineStatus(in *infrav1.OpenStackMachineStatus, out *OpenStackMachineStatus, _ apiconversion.Scope) error {
	out.Initialization = (*MachineInitialization)(unsafe.Pointer(in.Initialization))
	out.InstanceID = in.InstanceID
	out.Addresses = in.Addresses
	out.InstanceState = (*InstanceState)(unsafe.Pointer(in.InstanceState))
	out.Resolved = (*ResolvedMachineSpec)(unsafe.Pointer(in.Resolved))
	out.Resources = (*MachineResources)(unsafe.Pointer(in.Resources))

	out.Conditions = infrav1.ConvertConditionsFromV1Beta2(in.Conditions)
	out.Ready = infrav1.IsReady(in.Conditions)

	return nil
}

// Element-level Condition conversion functions required by conversion-gen's
// autoConvert functions for Status types. The actual condition conversion is
// handled at the Status level by the manual Convert_*_Status_* functions above.

func Convert_v1beta1_Condition_To_v1_Condition(in *clusterv1beta1.Condition, out *metav1.Condition, _ apiconversion.Scope) error {
	out.Type = string(in.Type)
	out.Status = metav1.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

func Convert_v1_Condition_To_v1beta1_Condition(in *metav1.Condition, out *clusterv1beta1.Condition, _ apiconversion.Scope) error {
	out.Type = clusterv1beta1.ConditionType(in.Type)
	out.Status = corev1.ConditionStatus(in.Status)
	out.LastTransitionTime = in.LastTransitionTime
	out.Reason = in.Reason
	out.Message = in.Message
	return nil
}

// Convert_v1beta1_OpenStackMachineSpec_To_v1beta2_OpenStackMachineSpec
// handles manual conversion for Flavor/FlavorID -> FlavorParam.
func Convert_v1beta1_OpenStackMachineSpec_To_v1beta2_OpenStackMachineSpec(
	in *OpenStackMachineSpec,
	out *infrav1.OpenStackMachineSpec,
	s apiconversion.Scope,
) error {
	// First copy all identical fields
	if err := autoConvert_v1beta1_OpenStackMachineSpec_To_v1beta2_OpenStackMachineSpec(in, out, s); err != nil {
		return err
	}

	switch {
	case in.FlavorID != nil:
		var id optional.String
		_ = optional.Convert_string_To_optional_String(in.FlavorID, &id, s)

		out.Flavor = infrav1.FlavorParam{
			ID: id,
		}

	case in.Flavor != nil:
		var name optional.String
		_ = optional.Convert_string_To_optional_String(in.Flavor, &name, s)

		out.Flavor = infrav1.FlavorParam{
			Filter: &infrav1.FlavorFilter{
				Name: name,
			},
		}

	default:
		// This would fail the validation webhook, but we need to handle
		// conversion anyway. The apiserver may send objects without
		// a spec e.g. during managedField conversion.
		// Therefore we return no error.
	}

	return nil
}

// Convert_v1beta2_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec
// handles manual conversion for FlavorParam -> Flavor/FlavorID.
func Convert_v1beta2_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(
	in *infrav1.OpenStackMachineSpec,
	out *OpenStackMachineSpec,
	s apiconversion.Scope,
) error {
	if err := autoConvert_v1beta2_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(in, out, s); err != nil {
		return err
	}

	switch {
	case in.Flavor.ID != nil && *in.Flavor.ID != "":
		id := *in.Flavor.ID
		out.FlavorID = &id
		out.Flavor = nil

	case in.Flavor.Filter != nil &&
		in.Flavor.Filter.Name != nil &&
		*in.Flavor.Filter.Name != "":
		name := *in.Flavor.Filter.Name
		out.Flavor = &name
		out.FlavorID = nil

	default:
		// This would fail the validation webhook, but we need to handle
		// conversion anyway. The apiserver may send objects without
		// a spec e.g. during managedField conversion.
		// Therefore we return no error.
	}

	return nil
}

func Convert_v1beta1_OpenStackClusterSpec_To_v1beta2_OpenStackClusterSpec(
	in *OpenStackClusterSpec,
	out *infrav1.OpenStackClusterSpec,
	s apiconversion.Scope,
) error {
	if err := autoConvert_v1beta1_OpenStackClusterSpec_To_v1beta2_OpenStackClusterSpec(in, out, s); err != nil {
		return err
	}

	if in.NetworkMTU != nil || in.DisablePortSecurity != nil {
		managed := &infrav1.ManagedNetwork{}
		if in.NetworkMTU != nil {
			mtu := int32(*in.NetworkMTU) //nolint:gosec // MTU values are always within int32 range
			managed.MTU = &mtu
		}
		if in.DisablePortSecurity != nil {
			managed.EnablePortSecurity = ptr.To(!*in.DisablePortSecurity)
		}
		out.ManagedNetwork = managed
	}

	if in.DisableExternalNetwork != nil {
		out.EnableExternalNetwork = ptr.To(!*in.DisableExternalNetwork)
	}

	// ExternalRouterIPs and ExternalRouterIPParam are structurally identical between versions,
	// but the nested SubnetParam is a different Go type per package, preventing a direct cast.
	// Loop and convert element-wise.
	if len(in.ExternalRouterIPs) > 0 {
		externalIPs := make([]infrav1.ExternalRouterIPParam, len(in.ExternalRouterIPs))
		for i := range in.ExternalRouterIPs {
			externalIPs[i].FixedIP = in.ExternalRouterIPs[i].FixedIP
			if err := Convert_v1beta1_SubnetParam_To_v1beta2_SubnetParam(&in.ExternalRouterIPs[i].Subnet, &externalIPs[i].Subnet, s); err != nil {
				return err
			}
		}
		out.ManagedRouter = &infrav1.ManagedRouter{
			ExternalIPs: externalIPs,
		}
	}

	// Consolidate the flat v1beta1 APIServer fields into the new APIServer struct.
	if in.APIServerLoadBalancer != nil ||
		in.DisableAPIServerFloatingIP != nil ||
		in.APIServerFloatingIP != nil ||
		in.APIServerFixedIP != nil ||
		in.APIServerPort != nil {
		out.APIServer = &infrav1.APIServer{
			FloatingIP: in.APIServerFloatingIP,
			FixedIP:    in.APIServerFixedIP,
			Port:       in.APIServerPort,
		}
		if in.DisableAPIServerFloatingIP != nil {
			out.APIServer.EnableFloatingIP = ptr.To(!*in.DisableAPIServerFloatingIP)
		}
		// APIServerLoadBalancer fields are converted field-by-field to handle
		// the int → int32 type changes in AdditionalPorts and Monitor fields.
		if in.APIServerLoadBalancer != nil {
			lb := in.APIServerLoadBalancer
			out.APIServer.ManagedLoadBalancer = &infrav1.APIServerLoadBalancer{
				Enabled:          lb.Enabled,
				AllowedCIDRs:     lb.AllowedCIDRs,
				Provider:         lb.Provider,
				Network:          (*infrav1.NetworkParam)(unsafe.Pointer(lb.Network)),
				Subnets:          *(*[]infrav1.SubnetParam)(unsafe.Pointer(&lb.Subnets)),
				AvailabilityZone: lb.AvailabilityZone,
				Flavor:           lb.Flavor,
			}
			for _, p := range lb.AdditionalPorts {
				out.APIServer.ManagedLoadBalancer.AdditionalPorts = append(out.APIServer.ManagedLoadBalancer.AdditionalPorts, int32(p)) //nolint:gosec // Port values are always within int32 range
			}
			if lb.Monitor != nil {
				out.APIServer.ManagedLoadBalancer.Monitor = &infrav1.APIServerLoadBalancerMonitor{
					Delay:          int32(lb.Monitor.Delay),          //nolint:gosec // Monitor values are always within int32 range
					Timeout:        int32(lb.Monitor.Timeout),        //nolint:gosec // Monitor values are always within int32 range
					MaxRetries:     int32(lb.Monitor.MaxRetries),     //nolint:gosec // Monitor values are always within int32 range
					MaxRetriesDown: int32(lb.Monitor.MaxRetriesDown), //nolint:gosec // Monitor values are always within int32 range
				}
			}
		}
	}

	return nil
}

func Convert_v1beta2_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(
	in *infrav1.OpenStackClusterSpec,
	out *OpenStackClusterSpec,
	s apiconversion.Scope,
) error {
	if err := autoConvert_v1beta2_OpenStackClusterSpec_To_v1beta1_OpenStackClusterSpec(in, out, s); err != nil {
		return err
	}

	if in.ManagedNetwork != nil {
		if in.ManagedNetwork.MTU != nil {
			mtu := int(*in.ManagedNetwork.MTU)
			out.NetworkMTU = &mtu
		}
		if in.ManagedNetwork.EnablePortSecurity != nil {
			out.DisablePortSecurity = ptr.To(!*in.ManagedNetwork.EnablePortSecurity)
		}
	}

	if in.EnableExternalNetwork != nil {
		out.DisableExternalNetwork = ptr.To(!*in.EnableExternalNetwork)
	}

	// ExternalRouterIPs and ExternalRouterIPParam are structurally identical between versions,
	// but the nested SubnetParam is a different Go type per package, preventing a direct cast.
	// Loop and convert element-wise.
	if in.ManagedRouter != nil {
		out.ExternalRouterIPs = make([]ExternalRouterIPParam, len(in.ManagedRouter.ExternalIPs))
		for i := range in.ManagedRouter.ExternalIPs {
			out.ExternalRouterIPs[i].FixedIP = in.ManagedRouter.ExternalIPs[i].FixedIP
			if err := Convert_v1beta2_SubnetParam_To_v1beta1_SubnetParam(&in.ManagedRouter.ExternalIPs[i].Subnet, &out.ExternalRouterIPs[i].Subnet, s); err != nil {
				return err
			}
		}
	}

	// Expand the v1beta2 APIServer struct back into the flat v1beta1 fields.
	if in.APIServer != nil {
		out.APIServerFloatingIP = in.APIServer.FloatingIP
		out.APIServerFixedIP = in.APIServer.FixedIP
		out.APIServerPort = in.APIServer.Port
		if in.APIServer.EnableFloatingIP != nil {
			out.DisableAPIServerFloatingIP = ptr.To(!*in.APIServer.EnableFloatingIP)
		}

		if in.APIServer.ManagedLoadBalancer != nil {
			lb := in.APIServer.ManagedLoadBalancer
			out.APIServerLoadBalancer = &APIServerLoadBalancer{
				Enabled:          lb.Enabled,
				AllowedCIDRs:     lb.AllowedCIDRs,
				Provider:         lb.Provider,
				Network:          (*NetworkParam)(unsafe.Pointer(lb.Network)),
				Subnets:          *(*[]SubnetParam)(unsafe.Pointer(&lb.Subnets)),
				AvailabilityZone: lb.AvailabilityZone,
				Flavor:           lb.Flavor,
			}
			for _, p := range lb.AdditionalPorts {
				out.APIServerLoadBalancer.AdditionalPorts = append(out.APIServerLoadBalancer.AdditionalPorts, int(p))
			}
			if lb.Monitor != nil {
				out.APIServerLoadBalancer.Monitor = &APIServerLoadBalancerMonitor{
					Delay:          int(lb.Monitor.Delay),
					Timeout:        int(lb.Monitor.Timeout),
					MaxRetries:     int(lb.Monitor.MaxRetries),
					MaxRetriesDown: int(lb.Monitor.MaxRetriesDown),
				}
			}
		}
	}

	return nil
}

func Convert_v1beta1_ManagedSecurityGroups_To_v1beta2_ManagedSecurityGroups(
	in *ManagedSecurityGroups,
	out *infrav1.ManagedSecurityGroups,
	s apiconversion.Scope,
) error {
	if err := autoConvert_v1beta1_ManagedSecurityGroups_To_v1beta2_ManagedSecurityGroups(in, out, s); err != nil {
		return err
	}

	if len(in.AllNodesSecurityGroupRules) > 0 {
		out.ClusterNodesSecurityGroupRules = make([]infrav1.SecurityGroupRuleSpec, len(in.AllNodesSecurityGroupRules))
		for i := range in.AllNodesSecurityGroupRules {
			if err := Convert_v1beta1_SecurityGroupRuleSpec_To_v1beta2_SecurityGroupRuleSpec(&in.AllNodesSecurityGroupRules[i], &out.ClusterNodesSecurityGroupRules[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

func Convert_v1beta2_ManagedSecurityGroups_To_v1beta1_ManagedSecurityGroups(
	in *infrav1.ManagedSecurityGroups,
	out *ManagedSecurityGroups,
	s apiconversion.Scope,
) error {
	if err := autoConvert_v1beta2_ManagedSecurityGroups_To_v1beta1_ManagedSecurityGroups(in, out, s); err != nil {
		return err
	}

	if len(in.ClusterNodesSecurityGroupRules) > 0 {
		out.AllNodesSecurityGroupRules = make([]SecurityGroupRuleSpec, len(in.ClusterNodesSecurityGroupRules))
		for i := range in.ClusterNodesSecurityGroupRules {
			if err := Convert_v1beta2_SecurityGroupRuleSpec_To_v1beta1_SecurityGroupRuleSpec(&in.ClusterNodesSecurityGroupRules[i], &out.AllNodesSecurityGroupRules[i], s); err != nil {
				return err
			}
		}
	}
	return nil
}

// LegacyCalicoSecurityGroupRules returns a list of security group rules for calico
// that need to be applied to the control plane and worker security groups when
// managed security groups are enabled and upgrading to v1beta1.
func LegacyCalicoSecurityGroupRules() []SecurityGroupRuleSpec {
	return []SecurityGroupRuleSpec{
		{
			Name:                "BGP (calico)",
			Description:         ptr.To("Created by cluster-api-provider-openstack API conversion - BGP (calico)"),
			Direction:           securityGroupRuleDirectionIngress,
			EtherType:           ptr.To("IPv4"),
			PortRangeMin:        ptr.To(179),
			PortRangeMax:        ptr.To(179),
			Protocol:            ptr.To("tcp"),
			RemoteManagedGroups: []ManagedSecurityGroupName{"controlplane", "worker"},
		},
		{
			Name:                "IP-in-IP (calico)",
			Description:         ptr.To("Created by cluster-api-provider-openstack API conversion - IP-in-IP (calico)"),
			Direction:           securityGroupRuleDirectionIngress,
			EtherType:           ptr.To("IPv4"),
			Protocol:            ptr.To("4"),
			RemoteManagedGroups: []ManagedSecurityGroupName{"controlplane", "worker"},
		},
	}
}

// splitTags splits a comma separated list of tags into a slice of tags.
// If the input is an empty string, it returns nil representing no list rather
// than an empty list.
func splitTags(tags string) []NeutronTag {
	if tags == "" {
		return nil
	}

	var ret []NeutronTag
	for _, tag := range strings.Split(tags, ",") {
		if tag != "" {
			ret = append(ret, NeutronTag(tag))
		}
	}

	return ret
}

// JoinTags joins a slice of tags into a comma separated list of tags.
func JoinTags(tags []NeutronTag) string {
	var b strings.Builder
	for i := range tags {
		if i > 0 {
			b.WriteString(",")
		}
		b.WriteString(string(tags[i]))
	}
	return b.String()
}

func ConvertAllTagsTo(tags, tagsAny, notTags, notTagsAny string, neutronTags *FilterByNeutronTags) {
	neutronTags.Tags = splitTags(tags)
	neutronTags.TagsAny = splitTags(tagsAny)
	neutronTags.NotTags = splitTags(notTags)
	neutronTags.NotTagsAny = splitTags(notTagsAny)
}

func ConvertAllTagsFrom(neutronTags *FilterByNeutronTags, tags, tagsAny, notTags, notTagsAny *string) {
	*tags = JoinTags(neutronTags.Tags)
	*tagsAny = JoinTags(neutronTags.TagsAny)
	*notTags = JoinTags(neutronTags.NotTags)
	*notTagsAny = JoinTags(neutronTags.NotTagsAny)
}

func Convert_v1beta1_ResolvedPortSpecFields_To_v1beta2_ResolvedPortSpecFields(in *ResolvedPortSpecFields, out *infrav1.ResolvedPortSpecFields, s apiconversion.Scope) error {
	if err := autoConvert_v1beta1_ResolvedPortSpecFields_To_v1beta2_ResolvedPortSpecFields(in, out, s); err != nil {
		return err
	}

	if in.DisablePortSecurity != nil {
		out.EnablePortSecurity = ptr.To(!*in.DisablePortSecurity)
	}

	return nil
}

func Convert_v1beta2_ResolvedPortSpecFields_To_v1beta1_ResolvedPortSpecFields(in *infrav1.ResolvedPortSpecFields, out *ResolvedPortSpecFields, s apiconversion.Scope) error {
	if err := autoConvert_v1beta2_ResolvedPortSpecFields_To_v1beta1_ResolvedPortSpecFields(in, out, s); err != nil {
		return err
	}

	if in.EnablePortSecurity != nil {
		out.DisablePortSecurity = ptr.To(!*in.EnablePortSecurity)
	}

	return nil
}
