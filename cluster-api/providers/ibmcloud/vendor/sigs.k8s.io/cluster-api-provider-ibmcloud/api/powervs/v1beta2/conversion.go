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
	"reflect"

	corev1 "k8s.io/api/core/v1"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	apimachineryconversion "k8s.io/apimachinery/pkg/conversion"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/controller-runtime/pkg/conversion"

	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1" //nolint:staticcheck
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-ibmcloud/api/powervs/v1beta3"
)

//nolint:gocyclo // complexity is acceptable for conversion logic
func Convert_v1beta2_IBMPowerVSClusterStatus_To_v1beta3_IBMPowerVSClusterStatus(in *IBMPowerVSClusterStatus, out *infrav1.IBMPowerVSClusterStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSClusterStatus_To_v1beta3_IBMPowerVSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manual conversion for Status: ServiceInstance -> Workspace.
	// v1beta2 ResourceReference only has an ID, no Name.
	if in.ServiceInstance != nil && in.ServiceInstance.ID != nil {
		out.Workspace.ID = *in.ServiceInstance.ID
	}

	// Manual conversion for Status: ResourceGroup.
	// v1beta2 ResourceReference only has an ID, no Name.
	if in.ResourceGroup != nil && in.ResourceGroup.ID != nil {
		out.ResourceGroup.ID = *in.ResourceGroup.ID
	}

	// Convert Network status
	if in.Network != nil && in.Network.ID != nil {
		out.Network.ID = *in.Network.ID
	}

	// Convert DHCPServer status
	if in.DHCPServer != nil && in.DHCPServer.ID != nil {
		out.Network.DHCPServer.ID = *in.DHCPServer.ID
	}

	// Convert TransitGateway status
	if in.TransitGateway != nil {
		if err := Convert_v1beta2_TransitGatewayStatus_To_v1beta3_TransitGatewayStatus(in.TransitGateway, &out.TransitGateway, s); err != nil {
			return err
		}
	}

	if in.VPC != nil {
		if err := Convert_v1beta2_ResourceReference_To_v1beta3_VPCStatus(in.VPC, &out.VPC, s); err != nil {
			return err
		}
	}

	// Convert COSInstance status: v1beta2 *ResourceReference -> v1beta3 COSInstanceStatus
	if in.COSInstance != nil && in.COSInstance.ID != nil && *in.COSInstance.ID != "" {
		out.COSInstance.ID = *in.COSInstance.ID
	}
	if in.VPCSubnet != nil {
		out.VPCSubnets = make([]infrav1.VPCSubnetStatus, 0, len(in.VPCSubnet))
		for name, subnet := range in.VPCSubnet {
			if subnet.ID == nil || *subnet.ID == "" || name == "" {
				continue
			}
			status := infrav1.VPCSubnetStatus{
				ID:   *subnet.ID,
				Name: name,
			}
			out.VPCSubnets = append(out.VPCSubnets, status)
		}
		if len(out.VPCSubnets) == 0 {
			out.VPCSubnets = nil
		}
	}
	if in.LoadBalancers != nil {
		out.LoadBalancers = make([]infrav1.LoadBalancerStatus, 0, len(in.LoadBalancers))
		for name, lb := range in.LoadBalancers {
			status := infrav1.LoadBalancerStatus{Name: name}
			if err := Convert_v1beta2_VPCLoadBalancerStatus_To_v1beta3_LoadBalancerStatus(&lb, &status, s); err != nil {
				return err
			}
			if status.Name == "" {
				continue
			}
			out.LoadBalancers = append(out.LoadBalancers, status)
		}
		if len(out.LoadBalancers) == 0 {
			out.LoadBalancers = nil
		}
	}
	// Convert VPCSecurityGroups status: v1beta2 map[string]VPCSecurityGroupStatus -> v1beta3 []VPCSecurityGroupStatus.
	// The map key is the security group name; it is stored in v1beta3 VPCSecurityGroupStatus.Name.
	if in.VPCSecurityGroups != nil {
		out.VPCSecurityGroups = make([]infrav1.VPCSecurityGroupStatus, 0, len(in.VPCSecurityGroups))
		for name, sg := range in.VPCSecurityGroups {
			if name == "" {
				continue
			}
			status := infrav1.VPCSecurityGroupStatus{Name: name}
			if err := Convert_v1beta2_VPCSecurityGroupStatus_To_v1beta3_VPCSecurityGroupStatus(&sg, &status, s); err != nil {
				return err
			}
			out.VPCSecurityGroups = append(out.VPCSecurityGroups, status)
		}
		if len(out.VPCSecurityGroups) == 0 {
			out.VPCSecurityGroups = nil
		}
	}
	out.Conditions = nil
	if in.V1Beta2 != nil {
		out.Conditions = in.V1Beta2.Conditions
	}
	if in.Conditions == nil {
		return nil
	}

	if out.Deprecated == nil {
		out.Deprecated = &infrav1.IBMPowerVSClusterDeprecatedStatus{}
	}
	if out.Deprecated.V1Beta2 == nil {
		out.Deprecated.V1Beta2 = &infrav1.IBMPowerVSClusterV1Beta2DeprecatedStatus{}
	}
	clusterv1beta1.Convert_v1beta1_Conditions_To_v1beta2_Deprecated_V1Beta1_Conditions(&in.Conditions, &out.Deprecated.V1Beta2.Conditions)
	return nil
}

//nolint:gocyclo // complexity is acceptable for conversion logic
func Convert_v1beta3_IBMPowerVSClusterStatus_To_v1beta2_IBMPowerVSClusterStatus(in *infrav1.IBMPowerVSClusterStatus, out *IBMPowerVSClusterStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSClusterStatus_To_v1beta2_IBMPowerVSClusterStatus(in, out, s); err != nil {
		return err
	}

	// Manual conversion for Status: Workspace -> ServiceInstance.
	if in.Workspace.ID != "" {
		out.ServiceInstance = &ResourceReference{
			ID: ptr.To(in.Workspace.ID),
		}
	} else {
		out.ServiceInstance = nil
	}

	// Manual conversion for Status: ResourceGroup.
	if in.ResourceGroup.ID != "" {
		out.ResourceGroup = &ResourceReference{
			ID: ptr.To(in.ResourceGroup.ID),
		}
	} else {
		out.ResourceGroup = nil
	}

	// Convert Network status
	if in.Network.ID != "" || in.Network.Name != "" {
		out.Network = &ResourceReference{}
		if in.Network.ID != "" {
			out.Network.ID = ptr.To(in.Network.ID)
		}
	}

	// Convert DHCPServer status
	if in.Network.DHCPServer.ID != "" || in.Network.DHCPServer.Name != "" {
		out.DHCPServer = &ResourceReference{}
		if in.Network.DHCPServer.ID != "" {
			out.DHCPServer.ID = ptr.To(in.Network.DHCPServer.ID)
		}
	}

	// Convert TransitGateway status
	tgStatus := &TransitGatewayStatus{}
	if err := Convert_v1beta3_TransitGatewayStatus_To_v1beta2_TransitGatewayStatus(&in.TransitGateway, tgStatus, s); err != nil {
		return err
	}
	// If the converted TransitGatewayStatus is empty, set to nil
	if tgStatus.ID == nil && tgStatus.VPCConnection == nil && tgStatus.PowerVSConnection == nil {
		out.TransitGateway = nil
	} else {
		out.TransitGateway = tgStatus
	}

	out.Conditions = nil
	if in.Deprecated != nil && in.Deprecated.V1Beta2 != nil && in.Deprecated.V1Beta2.Conditions != nil {
		clusterv1beta1.Convert_v1beta2_Deprecated_V1Beta1_Conditions_To_v1beta1_Conditions(&in.Deprecated.V1Beta2.Conditions, &out.Conditions)
	}

	out.Ready = ptr.Deref(in.Initialization.Provisioned, false)

	if in.VPC.ID != "" {
		out.VPC = &ResourceReference{}
		if err := Convert_v1beta3_VPCStatus_To_v1beta2_ResourceReference(&in.VPC, out.VPC, s); err != nil {
			return err
		}
	}

	// Convert COSInstance status: v1beta3 COSInstanceStatus -> v1beta2 *ResourceReference
	if in.COSInstance.ID != "" {
		out.COSInstance = &ResourceReference{
			ID: ptr.To(in.COSInstance.ID),
		}
	} else {
		out.COSInstance = nil
	}

	if in.VPCSubnets != nil {
		out.VPCSubnet = make(map[string]ResourceReference, len(in.VPCSubnets))
		for i := range in.VPCSubnets {
			subnet := ResourceReference{}
			if in.VPCSubnets[i].ID != "" {
				subnet.ID = ptr.To(in.VPCSubnets[i].ID)
			}
			name := in.VPCSubnets[i].Name
			if name == "" {
				name = ptr.Deref(subnet.ID, "")
			}
			if name == "" {
				continue
			}
			out.VPCSubnet[name] = subnet
		}
		if len(out.VPCSubnet) == 0 {
			out.VPCSubnet = nil
		}
	}

	if in.LoadBalancers != nil {
		out.LoadBalancers = make(map[string]VPCLoadBalancerStatus, len(in.LoadBalancers))
		for i := range in.LoadBalancers {
			lb := VPCLoadBalancerStatus{}
			if err := Convert_v1beta3_LoadBalancerStatus_To_v1beta2_VPCLoadBalancerStatus(&in.LoadBalancers[i], &lb, s); err != nil {
				return err
			}
			name := in.LoadBalancers[i].Name
			if name == "" {
				name = ptr.Deref(lb.ID, "")
			}
			if name == "" {
				continue
			}
			out.LoadBalancers[name] = lb
		}
		if len(out.LoadBalancers) == 0 {
			out.LoadBalancers = nil
		}
	}
	// Convert VPCSecurityGroups status: v1beta3 []VPCSecurityGroupStatus -> v1beta2 map[string]VPCSecurityGroupStatus.
	// The v1beta3 Name field becomes the map key; entries without a Name or ID are skipped.
	if in.VPCSecurityGroups != nil {
		out.VPCSecurityGroups = make(map[string]VPCSecurityGroupStatus, len(in.VPCSecurityGroups))
		for i := range in.VPCSecurityGroups {
			sg := VPCSecurityGroupStatus{}
			if err := Convert_v1beta3_VPCSecurityGroupStatus_To_v1beta2_VPCSecurityGroupStatus(&in.VPCSecurityGroups[i], &sg, s); err != nil {
				return err
			}
			name := in.VPCSecurityGroups[i].Name
			if name == "" {
				name = ptr.Deref(sg.ID, "")
			}
			if name == "" {
				continue
			}
			out.VPCSecurityGroups[name] = sg
		}
		if len(out.VPCSecurityGroups) == 0 {
			out.VPCSecurityGroups = nil
		}
	}
	if in.Conditions == nil {
		return nil
	}
	out.V1Beta2 = &IBMPowerVSClusterV1Beta2Status{}
	out.V1Beta2.Conditions = in.Conditions
	return nil
}

func Convert_v1beta2_IBMPowerVSMachineStatus_To_v1beta3_IBMPowerVSMachineStatus(in *IBMPowerVSMachineStatus, out *infrav1.IBMPowerVSMachineStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSMachineStatus_To_v1beta3_IBMPowerVSMachineStatus(in, out, s); err != nil {
		return err
	}
	out.Conditions = nil
	if in.V1Beta2 != nil {
		out.Conditions = in.V1Beta2.Conditions
	}
	if in.Conditions == nil {
		return nil
	}
	if out.Deprecated == nil {
		out.Deprecated = &infrav1.IBMPowerVSMachineDeprecatedStatus{}
	}
	if out.Deprecated.V1Beta2 == nil {
		out.Deprecated.V1Beta2 = &infrav1.IBMPowerVSMachineV1Beta2DeprecatedStatus{}
	}
	clusterv1beta1.Convert_v1beta1_Conditions_To_v1beta2_Deprecated_V1Beta1_Conditions(&in.Conditions, &out.Deprecated.V1Beta2.Conditions)
	return nil
}

func Convert_v1beta3_IBMPowerVSMachineStatus_To_v1beta2_IBMPowerVSMachineStatus(in *infrav1.IBMPowerVSMachineStatus, out *IBMPowerVSMachineStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSMachineStatus_To_v1beta2_IBMPowerVSMachineStatus(in, out, s); err != nil {
		return err
	}
	out.Conditions = nil
	if in.Deprecated != nil && in.Deprecated.V1Beta2 != nil && in.Deprecated.V1Beta2.Conditions != nil {
		clusterv1beta1.Convert_v1beta2_Deprecated_V1Beta1_Conditions_To_v1beta1_Conditions(&in.Deprecated.V1Beta2.Conditions, &out.Conditions)
	}
	out.Ready = ptr.Deref(in.Initialization.Provisioned, false)
	// Normalize empty-string pointers back to nil (empty string in v1beta3 == absent in v1beta2)
	if out.Region != nil && *out.Region == "" {
		out.Region = nil
	}
	if out.Zone != nil && *out.Zone == "" {
		out.Zone = nil
	}
	if in.Conditions == nil {
		return nil
	}
	out.V1Beta2 = &IBMPowerVSMachineV1Beta2Status{}
	out.V1Beta2.Conditions = in.Conditions
	return nil
}

func Convert_v1beta2_IBMPowerVSImageStatus_To_v1beta3_IBMPowerVSImageStatus(in *IBMPowerVSImageStatus, out *infrav1.IBMPowerVSImageStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSImageStatus_To_v1beta3_IBMPowerVSImageStatus(in, out, s); err != nil {
		return err
	}
	out.Conditions = nil
	if in.V1Beta2 != nil {
		out.Conditions = in.V1Beta2.Conditions
	}
	if in.Conditions == nil {
		return nil
	}
	if out.Deprecated == nil {
		out.Deprecated = &infrav1.IBMPowerVSImageDeprecatedStatus{}
	}
	if out.Deprecated.V1Beta2 == nil {
		out.Deprecated.V1Beta2 = &infrav1.IBMPowerVSImageV1Beta2DeprecatedStatus{}
	}
	clusterv1beta1.Convert_v1beta1_Conditions_To_v1beta2_Deprecated_V1Beta1_Conditions(&in.Conditions, &out.Deprecated.V1Beta2.Conditions)
	return nil
}

func Convert_v1beta3_IBMPowerVSImageStatus_To_v1beta2_IBMPowerVSImageStatus(in *infrav1.IBMPowerVSImageStatus, out *IBMPowerVSImageStatus, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSImageStatus_To_v1beta2_IBMPowerVSImageStatus(in, out, s); err != nil {
		return err
	}
	out.Conditions = nil
	if in.Deprecated != nil && in.Deprecated.V1Beta2 != nil && in.Deprecated.V1Beta2.Conditions != nil {
		clusterv1beta1.Convert_v1beta2_Deprecated_V1Beta1_Conditions_To_v1beta1_Conditions(&in.Deprecated.V1Beta2.Conditions, &out.Conditions)
	}
	if in.Conditions == nil {
		return nil
	}
	out.V1Beta2 = &IBMPowerVSImageV1Beta2Status{}
	out.V1Beta2.Conditions = in.Conditions
	return nil
}

func Convert_v1_Condition_To_v1beta1_Condition(in *metav1.Condition, out *clusterv1beta1.Condition, s apimachineryconversion.Scope) error {
	return clusterv1beta1.Convert_v1_Condition_To_v1beta1_Condition(in, out, s)
}

func Convert_v1beta1_Condition_To_v1_Condition(in *clusterv1beta1.Condition, out *metav1.Condition, s apimachineryconversion.Scope) error {
	return clusterv1beta1.Convert_v1beta1_Condition_To_v1_Condition(in, out, s)
}

func Convert_v1beta3_IBMPowerVSMachineTemplateResource_To_v1beta2_IBMPowerVSMachineTemplateResource(in *infrav1.IBMPowerVSMachineTemplateResource, out *IBMPowerVSMachineTemplateResource, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSMachineTemplateResource_To_v1beta2_IBMPowerVSMachineTemplateResource(in, out, s); err != nil {
		return err
	}
	// Network.RegEx is not preserved during round-trip conversion
	// because v1beta3 ResourceIdentifier doesn't support RegEx.
	out.Spec.Network.RegEx = nil

	// Ensure empty string pointers are converted to nil
	if out.Spec.Network.ID != nil && *out.Spec.Network.ID == "" {
		out.Spec.Network.ID = nil
	}
	if out.Spec.Network.Name != nil && *out.Spec.Network.Name == "" {
		out.Spec.Network.Name = nil
	}
	return nil
}

func Convert_v1beta1_APIEndpoint_To_v1beta3_APIEndpoint(in *clusterv1beta1.APIEndpoint, out *infrav1.APIEndpoint, _ apimachineryconversion.Scope) error {
	out.Host = in.Host
	out.Port = in.Port
	return nil
}

func Convert_v1beta3_APIEndpoint_To_v1beta1_APIEndpoint(in *infrav1.APIEndpoint, out *clusterv1beta1.APIEndpoint, _ apimachineryconversion.Scope) error {
	out.Host = in.Host
	out.Port = in.Port
	return nil
}

func Convert_v1beta1_ObjectMeta_To_v1beta2_ObjectMeta(in *clusterv1beta1.ObjectMeta, out *clusterv1.ObjectMeta, s apimachineryconversion.Scope) error {
	return clusterv1beta1.Convert_v1beta1_ObjectMeta_To_v1beta2_ObjectMeta(in, out, s)
}

func Convert_v1beta2_ObjectMeta_To_v1beta1_ObjectMeta(in *clusterv1.ObjectMeta, out *clusterv1beta1.ObjectMeta, s apimachineryconversion.Scope) error {
	return clusterv1beta1.Convert_v1beta2_ObjectMeta_To_v1beta1_ObjectMeta(in, out, s)
}

func (src *IBMPowerVSCluster) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.IBMPowerVSCluster)
	if err := Convert_v1beta2_IBMPowerVSCluster_To_v1beta3_IBMPowerVSCluster(src, dst, nil); err != nil {
		return err
	}
	restored := &infrav1.IBMPowerVSCluster{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	initialization := infrav1.IBMPowerVSClusterInitializationStatus{}
	clusterv1.Convert_bool_To_Pointer_bool(src.Status.Ready, ok, restored.Status.Initialization.Provisioned, &initialization.Provisioned)
	if !reflect.DeepEqual(initialization, infrav1.IBMPowerVSClusterInitializationStatus{}) {
		dst.Status.Initialization = initialization
	}

	// Restore fields with union types from annotation if available
	// This preserves the full v1beta3 structure including Type, Reference, Provision fields
	if ok {
		dst.Spec.Workspace = restored.Spec.Workspace
		dst.Spec.ResourceGroup = restored.Spec.ResourceGroup
		dst.Spec.Network = restored.Spec.Network
		dst.Spec.TransitGateway = restored.Spec.TransitGateway
		dst.Spec.VPC = restored.Spec.VPC
		dst.Spec.VPCSubnets = restored.Spec.VPCSubnets
		dst.Spec.LoadBalancers = restored.Spec.LoadBalancers
		dst.Spec.COSInstance = restored.Spec.COSInstance
		// If Type was lost (v1beta2 annotation has no Type field), infer it from provision/reference data
		if dst.Spec.COSInstance.Type == "" {
			if dst.Spec.COSInstance.Reference.ID != "" || dst.Spec.COSInstance.Reference.Name != "" {
				dst.Spec.COSInstance.Type = infrav1.SourceTypeReference
			} else if dst.Spec.COSInstance.BucketName != "" || dst.Spec.COSInstance.BucketRegion != "" || dst.Spec.COSInstance.Provision.Name != "" {
				dst.Spec.COSInstance.Type = infrav1.SourceTypeProvision
			}
		}
		dst.Spec.Ignition = restored.Spec.Ignition
		dst.Status.VPC = restored.Status.VPC
		dst.Status.VPCSubnets = restored.Status.VPCSubnets
		dst.Status.LoadBalancers = restored.Status.LoadBalancers
		dst.Status.COSInstance = restored.Status.COSInstance
		dst.Annotations = restored.Annotations
	}

	// Preserve empty/unknown topology when the legacy annotation is absent.
	if val, exists := src.Annotations["powervs.cluster.x-k8s.io/create-infra"]; exists {
		if val == "true" {
			dst.Spec.Topology = infrav1.PowerVSLoadBalancerTopology
		} else {
			dst.Spec.Topology = infrav1.PowerVSVirtualIPTopology
		}
	}

	// Clean up the annotation in v1beta3 so we don't have duplicated sources of truth
	if dst.Annotations != nil {
		delete(dst.Annotations, "powervs.cluster.x-k8s.io/create-infra")
		if len(dst.Annotations) == 0 {
			dst.Annotations = nil
		}
	}

	return nil
}

func (dst *IBMPowerVSCluster) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.IBMPowerVSCluster)
	if err := Convert_v1beta3_IBMPowerVSCluster_To_v1beta2_IBMPowerVSCluster(src, dst, nil); err != nil {
		return err
	}

	// Map the v1beta3 Topology explicit enum back to the v1beta2 annotation.
	// Preserve empty/unknown hub topology by not forcing the legacy annotation state.
	switch src.Spec.Topology {
	case infrav1.PowerVSLoadBalancerTopology:
		if dst.Annotations == nil {
			dst.Annotations = make(map[string]string)
		}
		dst.Annotations["powervs.cluster.x-k8s.io/create-infra"] = "true"
	case infrav1.PowerVSVirtualIPTopology:
		if dst.Annotations != nil {
			delete(dst.Annotations, "powervs.cluster.x-k8s.io/create-infra")
		}
	}

	restored := &IBMPowerVSCluster{
		Spec: IBMPowerVSClusterSpec{
			VPC:           dst.Spec.VPC,
			VPCSubnets:    dst.Spec.VPCSubnets,
			LoadBalancers: dst.Spec.LoadBalancers,
			CosInstance:   dst.Spec.CosInstance,
			Ignition:      dst.Spec.Ignition,
		},
		Status: IBMPowerVSClusterStatus{
			VPC:           dst.Status.VPC,
			VPCSubnet:     dst.Status.VPCSubnet,
			LoadBalancers: dst.Status.LoadBalancers,
			COSInstance:   dst.Status.COSInstance,
		},
	}
	if err := utilconversion.MarshalData(restored, dst); err != nil {
		return err
	}

	if ok, err := utilconversion.UnmarshalData(dst, restored); err != nil {
		return err
	} else if ok {
		dst.Spec.VPC = restored.Spec.VPC
		dst.Spec.VPCSubnets = restored.Spec.VPCSubnets
		dst.Spec.LoadBalancers = restored.Spec.LoadBalancers
		dst.Spec.CosInstance = restored.Spec.CosInstance
		dst.Spec.Ignition = restored.Spec.Ignition
		dst.Status.VPC = restored.Status.VPC
		dst.Status.VPCSubnet = restored.Status.VPCSubnet
		dst.Status.LoadBalancers = restored.Status.LoadBalancers
		dst.Status.COSInstance = restored.Status.COSInstance
	}

	// Fix annotation discrepancy during round-trip
	if len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	return nil
}

func (src *IBMPowerVSClusterTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.IBMPowerVSClusterTemplate)
	if err := Convert_v1beta2_IBMPowerVSClusterTemplate_To_v1beta3_IBMPowerVSClusterTemplate(src, dst, nil); err != nil {
		return err
	}
	restored := &infrav1.IBMPowerVSClusterTemplate{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	if ok {
		// Restore any fields that were lost in conversion
		dst.Spec = restored.Spec
		// Re-infer COSInstance.Type if it was lost in the v1beta2 annotation (which has no Type field)
		spec := &dst.Spec.Template.Spec
		if spec.COSInstance.Type == "" {
			if spec.COSInstance.Reference.ID != "" || spec.COSInstance.Reference.Name != "" {
				spec.COSInstance.Type = infrav1.SourceTypeReference
			} else if spec.COSInstance.BucketName != "" || spec.COSInstance.BucketRegion != "" || spec.COSInstance.Provision.Name != "" {
				spec.COSInstance.Type = infrav1.SourceTypeProvision
			}
		}
	}
	if dst.Annotations != nil && len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	return nil
}

func (dst *IBMPowerVSClusterTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.IBMPowerVSClusterTemplate)
	if err := Convert_v1beta3_IBMPowerVSClusterTemplate_To_v1beta2_IBMPowerVSClusterTemplate(src, dst, nil); err != nil {
		return err
	}
	restored := &IBMPowerVSClusterTemplate{
		Spec: IBMPowerVSClusterTemplateSpec{
			Template: IBMPowerVSClusterTemplateResource{
				ObjectMeta: dst.Spec.Template.ObjectMeta,
				Spec: IBMPowerVSClusterSpec{
					VPC:           dst.Spec.Template.Spec.VPC,
					VPCSubnets:    dst.Spec.Template.Spec.VPCSubnets,
					LoadBalancers: dst.Spec.Template.Spec.LoadBalancers,
					CosInstance:   dst.Spec.Template.Spec.CosInstance,
					Ignition:      dst.Spec.Template.Spec.Ignition,
				},
			},
		},
	}
	if err := utilconversion.MarshalData(restored, dst); err != nil {
		return err
	}

	if ok, err := utilconversion.UnmarshalData(dst, restored); err != nil {
		return err
	} else if ok {
		dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta
		dst.Spec.Template.Spec.VPC = restored.Spec.Template.Spec.VPC
		dst.Spec.Template.Spec.VPCSubnets = restored.Spec.Template.Spec.VPCSubnets
		dst.Spec.Template.Spec.LoadBalancers = restored.Spec.Template.Spec.LoadBalancers
		dst.Spec.Template.Spec.CosInstance = restored.Spec.Template.Spec.CosInstance
		dst.Spec.Template.Spec.Ignition = restored.Spec.Template.Spec.Ignition
	}

	if len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	if dst.Spec.Template.Spec.ResourceGroup != nil &&
		dst.Spec.Template.Spec.ResourceGroup.ID == nil &&
		dst.Spec.Template.Spec.ResourceGroup.Name == nil &&
		dst.Spec.Template.Spec.ResourceGroup.RegEx == nil {
		dst.Spec.Template.Spec.ResourceGroup = nil
	}
	return nil
}

func (src *IBMPowerVSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.IBMPowerVSMachine)
	if err := Convert_v1beta2_IBMPowerVSMachine_To_v1beta3_IBMPowerVSMachine(src, dst, nil); err != nil {
		return err
	}
	restored := &infrav1.IBMPowerVSMachine{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	initialization := infrav1.IBMPowerVSMachineInitializationStatus{}
	clusterv1.Convert_bool_To_Pointer_bool(src.Status.Ready, ok, restored.Status.Initialization.Provisioned, &initialization.Provisioned)
	if !reflect.DeepEqual(initialization, infrav1.IBMPowerVSMachineInitializationStatus{}) {
		dst.Status.Initialization = initialization
	}
	return nil
}

func (dst *IBMPowerVSMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.IBMPowerVSMachine)
	if err := Convert_v1beta3_IBMPowerVSMachine_To_v1beta2_IBMPowerVSMachine(src, dst, nil); err != nil {
		return err
	}
	if dst.Spec.ProviderID != nil && *dst.Spec.ProviderID == "" {
		dst.Spec.ProviderID = nil
	}
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}
	if len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	return nil
}

func (src *IBMPowerVSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.IBMPowerVSMachineTemplate)
	if err := Convert_v1beta2_IBMPowerVSMachineTemplate_To_v1beta3_IBMPowerVSMachineTemplate(src, dst, nil); err != nil {
		return err
	}
	restored := &infrav1.IBMPowerVSMachineTemplate{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	if ok {
		dst.Status = restored.Status
	}
	return nil
}

func (dst *IBMPowerVSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.IBMPowerVSMachineTemplate)
	if err := Convert_v1beta3_IBMPowerVSMachineTemplate_To_v1beta2_IBMPowerVSMachineTemplate(src, dst, nil); err != nil {
		return err
	}
	if dst.Spec.Template.Spec.ProviderID != nil && *dst.Spec.Template.Spec.ProviderID == "" {
		dst.Spec.Template.Spec.ProviderID = nil
	}
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}
	if len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	// Network.RegEx is not preserved during round-trip conversion
	// because v1beta3 ResourceIdentifier doesn't support RegEx.
	// Clear it after MarshalData to ensure it's not restored.
	dst.Spec.Template.Spec.Network.RegEx = nil
	return nil
}

func (src *IBMPowerVSImage) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.IBMPowerVSImage)
	if err := Convert_v1beta2_IBMPowerVSImage_To_v1beta3_IBMPowerVSImage(src, dst, nil); err != nil {
		return err
	}
	restored := &infrav1.IBMPowerVSImage{}
	ok, err := utilconversion.UnmarshalData(src, restored)
	if err != nil {
		return err
	}
	if ok {
		// Restore any fields that were lost in conversion
		dst.Spec = restored.Spec
	}
	return nil
}

func (dst *IBMPowerVSImage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.IBMPowerVSImage)
	if err := Convert_v1beta3_IBMPowerVSImage_To_v1beta2_IBMPowerVSImage(src, dst, nil); err != nil {
		return err
	}
	// Ready is derived from ImageState: active → true.
	dst.Status.Ready = src.Status.ImageState == infrav1.PowerVSImageStateACTIVE
	if err := utilconversion.MarshalData(src, dst); err != nil {
		return err
	}
	if len(dst.Annotations) == 0 {
		dst.Annotations = nil
	}
	return nil
}

func Convert_v1beta2_IBMPowerVSClusterSpec_To_v1beta3_IBMPowerVSClusterSpec(in *IBMPowerVSClusterSpec, out *infrav1.IBMPowerVSClusterSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSClusterSpec_To_v1beta3_IBMPowerVSClusterSpec(in, out, s); err != nil {
		return err
	}

	convertV1beta2WorkspaceToV1beta3(in, out)
	convertV1beta2ResourceGroupToV1beta3(in, out)

	if err := convertV1beta2NetworkToV1beta3(in, out); err != nil {
		return err
	}

	// Convert TransitGateway
	if in.TransitGateway != nil {
		if err := Convert_v1beta2_TransitGateway_To_v1beta3_TransitGatewaySource(in.TransitGateway, &out.TransitGateway, s); err != nil {
			return err
		}
	}

	if in.VPC != nil {
		if err := Convert_v1beta2_VPCResourceReference_To_v1beta3_VPCSource(in.VPC, &out.VPC, s); err != nil {
			return err
		}
	} else {
		out.VPC.Type = infrav1.SourceTypeProvision
	}

	// Convert CosInstance: v1beta2 *CosInstance -> v1beta3 COSInstanceSource (value)
	if in.CosInstance != nil {
		convertV1beta2CosInstanceToV1beta3(in.CosInstance, &out.COSInstance)
	}

	// Convert Ignition: v1beta2 *Ignition -> v1beta3 Ignition (value)
	if in.Ignition != nil {
		out.Ignition = infrav1.Ignition{Version: in.Ignition.Version}
	}

	return nil
}

func convertV1beta2WorkspaceToV1beta3(in *IBMPowerVSClusterSpec, out *infrav1.IBMPowerVSClusterSpec) {
	if in.ServiceInstance == nil && in.ServiceInstanceID == "" {
		out.Workspace.Type = infrav1.SourceTypeProvision
		return
	}

	out.Workspace.Type = infrav1.SourceTypeReference
	if in.ServiceInstance != nil {
		if in.ServiceInstance.ID != nil {
			out.Workspace.Reference.ID = *in.ServiceInstance.ID
		}
		if in.ServiceInstance.Name != nil {
			out.Workspace.Reference.Name = *in.ServiceInstance.Name
		}
	}
	if in.ServiceInstanceID != "" {
		out.Workspace.Reference.ID = in.ServiceInstanceID
	}
}

func convertV1beta2ResourceGroupToV1beta3(in *IBMPowerVSClusterSpec, out *infrav1.IBMPowerVSClusterSpec) {
	if in.ResourceGroup == nil {
		return
	}

	out.ResourceGroup.Type = infrav1.SourceTypeReference
	if in.ResourceGroup.ID != nil {
		out.ResourceGroup.Reference.ID = *in.ResourceGroup.ID
	}
	if in.ResourceGroup.Name != nil {
		out.ResourceGroup.Reference.Name = *in.ResourceGroup.Name
	}
}

func convertV1beta2NetworkToV1beta3(in *IBMPowerVSClusterSpec, out *infrav1.IBMPowerVSClusterSpec) error {
	if (in.Network.ID != nil && *in.Network.ID != "") || (in.Network.Name != nil && *in.Network.Name != "") {
		out.Network.Type = infrav1.SourceTypeReference
		if in.Network.ID != nil && *in.Network.ID != "" {
			out.Network.Reference.ID = *in.Network.ID
		}
		if in.Network.Name != nil && *in.Network.Name != "" {
			out.Network.Reference.Name = *in.Network.Name
		}
		return nil
	}

	if in.DHCPServer == nil {
		return nil
	}

	out.Network.Type = infrav1.SourceTypeProvision
	return Convert_v1beta2_DHCPServer_To_v1beta3_DHCPServer(in.DHCPServer, &out.Network.Provision.DHCPServer, nil)
}

func convertV1beta2CosInstanceToV1beta3(in *CosInstance, out *infrav1.COSInstanceSource) {
	// v1beta2 CosInstance has no Type field; it is always provision-style (Name/BucketName/BucketRegion).
	// Any non-empty CosInstance is treated as SourceTypeProvision.
	out.BucketName = in.BucketName
	out.BucketRegion = in.BucketRegion
	out.Provision.Name = in.Name
	if in.Name != "" || in.BucketName != "" || in.BucketRegion != "" {
		out.Type = infrav1.SourceTypeProvision
	}
}

func convertV1beta3CosInstanceToV1beta2(in *infrav1.COSInstanceSource, out *CosInstance) {
	out.BucketName = in.BucketName
	out.BucketRegion = in.BucketRegion
	switch in.Type {
	case infrav1.SourceTypeProvision:
		out.Name = in.Provision.Name
	case infrav1.SourceTypeReference:
		// Reference mode: no Name in v1beta2 CosInstance; keep BucketName/BucketRegion only
	default:
		out.Name = in.Provision.Name
	}
}

//nolint:gocyclo
func Convert_v1beta3_IBMPowerVSClusterSpec_To_v1beta2_IBMPowerVSClusterSpec(in *infrav1.IBMPowerVSClusterSpec, out *IBMPowerVSClusterSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSClusterSpec_To_v1beta2_IBMPowerVSClusterSpec(in, out, s); err != nil {
		return err
	}

	if in.Zone == "" {
		out.Zone = nil
	}

	if in.Workspace.Type == infrav1.SourceTypeReference || in.Workspace.Reference.ID != "" || in.Workspace.Reference.Name != "" {
		out.ServiceInstance = &IBMPowerVSResourceReference{}
		if in.Workspace.Reference.ID != "" {
			out.ServiceInstance.ID = ptr.To(in.Workspace.Reference.ID)
			out.ServiceInstanceID = in.Workspace.Reference.ID
		}
		if in.Workspace.Reference.Name != "" {
			out.ServiceInstance.Name = ptr.To(in.Workspace.Reference.Name)
		}
	} else {
		out.ServiceInstance = nil
		out.ServiceInstanceID = ""
	}

	if in.ResourceGroup.Type == infrav1.SourceTypeReference || in.ResourceGroup.Reference.ID != "" || in.ResourceGroup.Reference.Name != "" {
		out.ResourceGroup = &IBMPowerVSResourceReference{}
		if in.ResourceGroup.Reference.ID != "" {
			out.ResourceGroup.ID = ptr.To(in.ResourceGroup.Reference.ID)
		}
		if in.ResourceGroup.Reference.Name != "" {
			out.ResourceGroup.Name = ptr.To(in.ResourceGroup.Reference.Name)
		}
	} else {
		out.ResourceGroup = nil
	}

	// Convert Network field
	switch in.Network.Type {
	case infrav1.SourceTypeReference:
		// Convert reference
		if in.Network.Reference.ID != "" {
			out.Network.ID = ptr.To(in.Network.Reference.ID)
		}
		if in.Network.Reference.Name != "" {
			out.Network.Name = ptr.To(in.Network.Reference.Name)
		}
		out.DHCPServer = nil
	case infrav1.SourceTypeProvision:
		// Convert provision to DHCPServer
		dhcp := &DHCPServer{}
		if err := Convert_v1beta3_DHCPServer_To_v1beta2_DHCPServer(&in.Network.Provision.DHCPServer, dhcp, nil); err != nil {
			return err
		}
		out.DHCPServer = dhcp
		// Clear network reference when provisioning
		out.Network.ID = nil
		out.Network.Name = nil
	default:
		// Empty NetworkSource - clear both fields
		out.Network.ID = nil
		out.Network.Name = nil
		out.DHCPServer = nil
	}
	// Note: v1beta2 Network.RegEx field is dropped in v1beta3, cannot be restored
	out.Network.RegEx = nil

	// Convert TransitGateway
	tg := &TransitGateway{}
	if err := Convert_v1beta3_TransitGatewaySource_To_v1beta2_TransitGateway(&in.TransitGateway, tg, s); err != nil {
		return err
	}
	// If the converted TransitGateway is empty, set to nil
	if tg.ID == nil && tg.Name == nil && tg.GlobalRouting == nil {
		out.TransitGateway = nil
	} else {
		out.TransitGateway = tg
	}

	if in.VPC.Type != "" || in.VPC.Reference.ID != "" || in.VPC.Reference.Name != "" || in.VPC.Provision.Name != "" || in.VPC.Region != "" {
		out.VPC = &VPCResourceReference{}
		if err := Convert_v1beta3_VPCSource_To_v1beta2_VPCResourceReference(&in.VPC, out.VPC, s); err != nil {
			return err
		}
		if out.VPC.ID == nil && out.VPC.Name == nil && out.VPC.Region == nil {
			if in.VPC.Reference.ID != "" {
				out.VPC.ID = ptr.To(in.VPC.Reference.ID)
			}
			if in.VPC.Reference.Name != "" {
				out.VPC.Name = ptr.To(in.VPC.Reference.Name)
			}
			if in.VPC.Provision.Name != "" && out.VPC.Name == nil {
				out.VPC.Name = ptr.To(in.VPC.Provision.Name)
			}
			if in.VPC.Region != "" {
				out.VPC.Region = ptr.To(in.VPC.Region)
			}
		}
		if out.VPC.ID == nil && out.VPC.Name == nil && out.VPC.Region == nil {
			out.VPC = nil
		}
	}

	// Convert COSInstance: v1beta3 COSInstanceSource (value) -> v1beta2 *CosInstance
	if in.COSInstance.BucketName != "" || in.COSInstance.BucketRegion != "" || in.COSInstance.Type != "" {
		cos := &CosInstance{}
		convertV1beta3CosInstanceToV1beta2(&in.COSInstance, cos)
		out.CosInstance = cos
	} else {
		out.CosInstance = nil
	}

	// Convert Ignition: v1beta3 Ignition (value) -> v1beta2 *Ignition
	if in.Ignition.Version != "" {
		out.Ignition = &Ignition{Version: in.Ignition.Version}
	} else {
		out.Ignition = nil
	}

	return nil
}

func Convert_v1beta2_IBMPowerVSMachineSpec_To_v1beta3_IBMPowerVSMachineSpec(in *IBMPowerVSMachineSpec, out *infrav1.IBMPowerVSMachineSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSMachineSpec_To_v1beta3_IBMPowerVSMachineSpec(in, out, s); err != nil {
		return err
	}

	// v1beta2 ImageRef (LocalObjectReference → IBMPowerVSImage CRD) maps to v1beta3 Image.Type=Import
	if in.ImageRef != nil && in.ImageRef.Name != "" {
		out.Image.Type = infrav1.ImageSourceTypeImport
		out.Image.Import = infrav1.ImageReference{Name: in.ImageRef.Name}
	} else if in.Image != nil {
		// v1beta2 Image (IBMPowerVSResourceReference → existing image) maps to v1beta3 Image.Type=Reference
		out.Image.Type = infrav1.ImageSourceTypeReference
		if in.Image.ID != nil {
			out.Image.Reference.ID = *in.Image.ID
		}
		if in.Image.Name != nil {
			out.Image.Reference.Name = *in.Image.Name
		}
		// RegEx is intentionally dropped — not preserved in v1beta3
	}

	// Machine uses ResourceIdentifier, NO Type/Provision fields!
	if in.ServiceInstance != nil {
		if in.ServiceInstance.ID != nil {
			out.Workspace.ID = *in.ServiceInstance.ID
		}
		if in.ServiceInstance.Name != nil {
			out.Workspace.Name = *in.ServiceInstance.Name
		}
	} else if in.ServiceInstanceID != "" {
		out.Workspace.ID = in.ServiceInstanceID
	}

	return nil
}

func Convert_v1beta3_IBMPowerVSMachineSpec_To_v1beta2_IBMPowerVSMachineSpec(in *infrav1.IBMPowerVSMachineSpec, out *IBMPowerVSMachineSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSMachineSpec_To_v1beta2_IBMPowerVSMachineSpec(in, out, s); err != nil {
		return err
	}

	switch in.Image.Type {
	case infrav1.ImageSourceTypeImport:
		// v1beta3 Import → v1beta2 ImageRef (IBMPowerVSImage CRD reference)
		if in.Image.Import.Name != "" {
			out.ImageRef = &corev1.LocalObjectReference{Name: in.Image.Import.Name}
		} else {
			out.ImageRef = nil
		}
		out.Image = nil
	case infrav1.ImageSourceTypeReference:
		// v1beta3 Reference → v1beta2 Image (existing PowerVS image)
		out.Image = &IBMPowerVSResourceReference{}
		if in.Image.Reference.ID != "" {
			out.Image.ID = ptr.To(in.Image.Reference.ID)
		}
		if in.Image.Reference.Name != "" {
			out.Image.Name = ptr.To(in.Image.Reference.Name)
		}
		out.ImageRef = nil
	default:
		out.Image = nil
		out.ImageRef = nil
	}

	if in.Workspace.ID != "" || in.Workspace.Name != "" {
		out.ServiceInstance = &IBMPowerVSResourceReference{}
		if in.Workspace.ID != "" {
			out.ServiceInstance.ID = ptr.To(in.Workspace.ID)
			out.ServiceInstanceID = in.Workspace.ID
		}
		if in.Workspace.Name != "" {
			out.ServiceInstance.Name = ptr.To(in.Workspace.Name)
		}
	} else {
		out.ServiceInstance = nil
		out.ServiceInstanceID = ""
	}

	// Network.RegEx is not preserved during round-trip conversion
	// because v1beta3 ResourceIdentifier doesn't support RegEx.
	// Set RegEx to nil to ensure clean conversion.
	out.Network.RegEx = nil

	return nil
}

func Convert_v1beta2_IBMPowerVSImageSpec_To_v1beta3_IBMPowerVSImageSpec(in *IBMPowerVSImageSpec, out *infrav1.IBMPowerVSImageSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta2_IBMPowerVSImageSpec_To_v1beta3_IBMPowerVSImageSpec(in, out, s); err != nil {
		return err
	}

	// Image uses ResourceIdentifier, NO Type/Provision fields!
	if in.ServiceInstance != nil {
		if in.ServiceInstance.ID != nil {
			out.Workspace.ID = *in.ServiceInstance.ID
		}
		if in.ServiceInstance.Name != nil {
			out.Workspace.Name = *in.ServiceInstance.Name
		}
	} else if in.ServiceInstanceID != "" {
		out.Workspace.ID = in.ServiceInstanceID
	}

	return nil
}

func Convert_v1beta3_IBMPowerVSImageSpec_To_v1beta2_IBMPowerVSImageSpec(in *infrav1.IBMPowerVSImageSpec, out *IBMPowerVSImageSpec, s apimachineryconversion.Scope) error {
	if err := autoConvert_v1beta3_IBMPowerVSImageSpec_To_v1beta2_IBMPowerVSImageSpec(in, out, s); err != nil {
		return err
	}

	// autoConvert uses Convert_string_To_Pointer_string which turns "" into &"".
	// Preserve the v1beta2 convention that nil means "not set".
	if out.Bucket != nil && *out.Bucket == "" {
		out.Bucket = nil
	}
	if out.Object != nil && *out.Object == "" {
		out.Object = nil
	}
	if out.Region != nil && *out.Region == "" {
		out.Region = nil
	}

	if in.Workspace.ID != "" || in.Workspace.Name != "" {
		out.ServiceInstance = &IBMPowerVSResourceReference{}
		if in.Workspace.ID != "" {
			out.ServiceInstance.ID = ptr.To(in.Workspace.ID)
			out.ServiceInstanceID = in.Workspace.ID
		}
		if in.Workspace.Name != "" {
			out.ServiceInstance.Name = ptr.To(in.Workspace.Name)
		}
	} else {
		out.ServiceInstance = nil
		out.ServiceInstanceID = ""
	}

	return nil
}

// Convert_v1beta2_DHCPServer_To_v1beta3_DHCPServer handles the conversion from v1beta2 to v1beta3 DHCPServer.
func Convert_v1beta2_DHCPServer_To_v1beta3_DHCPServer(in *DHCPServer, out *infrav1.DHCPServer, _ apimachineryconversion.Scope) error {
	// Convert Cidr (pointer) to CIDR (string)
	if in.Cidr != nil {
		out.CIDR = *in.Cidr
	}

	// Convert DNSServer (pointer) to DNSServer (string)
	if in.DNSServer != nil {
		out.DNSServer = *in.DNSServer
	}

	// Convert Name (pointer) to Name (string)
	if in.Name != nil {
		out.Name = *in.Name
	}

	// Convert Snat (bool pointer) to Snat (DHCPSnatPolicy enum)
	if in.Snat != nil {
		if *in.Snat {
			out.Snat = infrav1.DHCPSnatPolicyEnabled
		} else {
			out.Snat = infrav1.DHCPSnatPolicyDisabled
		}
	}

	// Note: v1beta2.ID field is dropped in v1beta3 as it's not part of the spec
	return nil
}

// Convert_v1beta3_DHCPServer_To_v1beta2_DHCPServer handles the conversion from v1beta3 to v1beta2 DHCPServer.
func Convert_v1beta3_DHCPServer_To_v1beta2_DHCPServer(in *infrav1.DHCPServer, out *DHCPServer, _ apimachineryconversion.Scope) error {
	// Convert CIDR (string) to Cidr (pointer)
	if in.CIDR != "" {
		out.Cidr = ptr.To(in.CIDR)
	}

	// Convert DNSServer (string) to DNSServer (pointer)
	if in.DNSServer != "" {
		out.DNSServer = ptr.To(in.DNSServer)
	}

	// Convert Name (string) to Name (pointer)
	if in.Name != "" {
		out.Name = ptr.To(in.Name)
	}

	// Convert Snat (DHCPSnatPolicy enum) to Snat (bool pointer)
	if in.Snat != "" {
		out.Snat = ptr.To(in.Snat == infrav1.DHCPSnatPolicyEnabled)
	}

	// Note: v1beta2.ID field cannot be populated from v1beta3 as it doesn't exist there
	return nil
}

// Convert_v1beta2_IBMPowerVSResourceReference_To_v1beta3_NetworkSource is a stub conversion function.
// The actual conversion is handled in Convert_v1beta2_IBMPowerVSClusterSpec_To_v1beta3_IBMPowerVSClusterSpec
// because Network and DHCPServer fields need to be converted together.
func Convert_v1beta2_IBMPowerVSResourceReference_To_v1beta3_NetworkSource(_ *IBMPowerVSResourceReference, _ *infrav1.NetworkSource, _ apimachineryconversion.Scope) error {
	// Conversion handled in Spec conversion
	return nil
}

// Convert_v1beta3_NetworkSource_To_v1beta2_IBMPowerVSResourceReference is a stub conversion function.
// The actual conversion is handled in Convert_v1beta3_IBMPowerVSClusterSpec_To_v1beta2_IBMPowerVSClusterSpec
// because Network and DHCPServer fields need to be converted together.
func Convert_v1beta3_NetworkSource_To_v1beta2_IBMPowerVSResourceReference(_ *infrav1.NetworkSource, _ *IBMPowerVSResourceReference, _ apimachineryconversion.Scope) error {
	// Conversion handled in Spec conversion
	return nil
}

// Convert_v1beta2_IBMPowerVSResourceReference_To_v1beta3_ResourceIdentifier converts v1beta2 IBMPowerVSResourceReference to v1beta3 ResourceIdentifier.
func Convert_v1beta2_IBMPowerVSResourceReference_To_v1beta3_ResourceIdentifier(in *IBMPowerVSResourceReference, out *infrav1.ResourceIdentifier, _ apimachineryconversion.Scope) error {
	if in.ID != nil && *in.ID != "" {
		out.ID = *in.ID
	}
	if in.Name != nil && *in.Name != "" {
		out.Name = *in.Name
	}
	// Note: RegEx field is dropped in v1beta3
	return nil
}

// Convert_v1beta3_ResourceIdentifier_To_v1beta2_IBMPowerVSResourceReference converts v1beta3 ResourceIdentifier to v1beta2 IBMPowerVSResourceReference.
func Convert_v1beta3_ResourceIdentifier_To_v1beta2_IBMPowerVSResourceReference(in *infrav1.ResourceIdentifier, out *IBMPowerVSResourceReference, _ apimachineryconversion.Scope) error {
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	} else {
		out.ID = nil
	}
	if in.Name != "" {
		out.Name = ptr.To(in.Name)
	} else {
		out.Name = nil
	}
	// RegEx field cannot be restored from v1beta3
	out.RegEx = nil
	return nil
}

// Convert_v1beta2_TransitGateway_To_v1beta3_TransitGatewaySource converts v1beta2 TransitGateway to v1beta3 TransitGatewaySource.
func Convert_v1beta2_TransitGateway_To_v1beta3_TransitGatewaySource(in *TransitGateway, out *infrav1.TransitGatewaySource, _ apimachineryconversion.Scope) error {
	// Determine if this is Reference or Provision based on the fields present
	// Priority: ID presence indicates Reference, GlobalRouting presence indicates Provision

	hasID := in.ID != nil && *in.ID != ""
	hasName := in.Name != nil && *in.Name != ""
	hasGlobalRouting := in.GlobalRouting != nil

	// If ID is set, it's definitely a Reference
	if hasID {
		out.Type = infrav1.SourceTypeReference
		out.Reference.ID = *in.ID
		if hasName {
			out.Reference.Name = *in.Name
		}
		return nil
	}

	// If GlobalRouting is set, it's Provision (GlobalRouting only exists for Provision)
	if hasGlobalRouting {
		out.Type = infrav1.SourceTypeProvision
		if hasName {
			out.Provision.Name = *in.Name
		}
		if *in.GlobalRouting {
			out.Provision.GlobalRouting = infrav1.TransitGatewayRoutingGlobal
		} else {
			out.Provision.GlobalRouting = infrav1.TransitGatewayRoutingLocal
		}
		return nil
	}

	// If only Name is set (no ID, no GlobalRouting), we can't determine the type reliably
	// This is a known limitation: v1beta2 Name field is ambiguous
	// Default to Reference for backward compatibility
	if hasName {
		out.Type = infrav1.SourceTypeReference
		out.Reference.Name = *in.Name
		return nil
	}

	// No fields set, leave empty
	return nil
}

// Convert_v1beta3_TransitGatewaySource_To_v1beta2_TransitGateway converts v1beta3 TransitGatewaySource to v1beta2 TransitGateway.
func Convert_v1beta3_TransitGatewaySource_To_v1beta2_TransitGateway(in *infrav1.TransitGatewaySource, out *TransitGateway, _ apimachineryconversion.Scope) error {
	// If Type is empty, check if there's any data in Reference or Provision
	if in.Type == "" {
		// If there's reference data, treat as reference
		if in.Reference.ID != "" || in.Reference.Name != "" {
			if in.Reference.ID != "" {
				out.ID = ptr.To(in.Reference.ID)
			}
			if in.Reference.Name != "" {
				out.Name = ptr.To(in.Reference.Name)
			}
			return nil
		}
		// If there's provision data, treat as provision
		if in.Provision.Name != "" || in.Provision.GlobalRouting != "" {
			if in.Provision.Name != "" {
				out.Name = ptr.To(in.Provision.Name)
			}
			if in.Provision.GlobalRouting != "" {
				out.GlobalRouting = ptr.To(in.Provision.GlobalRouting == infrav1.TransitGatewayRoutingGlobal)
			}
			return nil
		}
		// No data, return empty
		return nil
	}

	switch in.Type {
	case infrav1.SourceTypeReference:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		}
	case infrav1.SourceTypeProvision:
		// For provision, set name if specified
		if in.Provision.Name != "" {
			out.Name = ptr.To(in.Provision.Name)
		}
		// Convert GlobalRouting enum to bool
		if in.Provision.GlobalRouting != "" {
			out.GlobalRouting = ptr.To(in.Provision.GlobalRouting == infrav1.TransitGatewayRoutingGlobal)
		}
	}

	return nil
}

// Convert_v1beta2_TransitGatewayStatus_To_v1beta3_TransitGatewayStatus converts v1beta2 TransitGatewayStatus to v1beta3 TransitGatewayStatus.
func Convert_v1beta2_TransitGatewayStatus_To_v1beta3_TransitGatewayStatus(in *TransitGatewayStatus, out *infrav1.TransitGatewayStatus, _ apimachineryconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
	}

	// Convert VPC connection status
	if in.VPCConnection != nil && in.VPCConnection.ID != nil {
		out.VPCConnection.ID = *in.VPCConnection.ID
	}

	// Convert PowerVS connection status
	if in.PowerVSConnection != nil && in.PowerVSConnection.ID != nil {
		out.PowerVSConnection.ID = *in.PowerVSConnection.ID
	}

	return nil
}

// Convert_v1beta3_TransitGatewayStatus_To_v1beta2_TransitGatewayStatus converts v1beta3 TransitGatewayStatus to v1beta2 TransitGatewayStatus.
func Convert_v1beta3_TransitGatewayStatus_To_v1beta2_TransitGatewayStatus(in *infrav1.TransitGatewayStatus, out *TransitGatewayStatus, _ apimachineryconversion.Scope) error {
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	}

	// Convert VPC connection status
	if in.VPCConnection.ID != "" {
		out.VPCConnection = &ResourceReference{
			ID: ptr.To(in.VPCConnection.ID),
		}
	}

	// Convert PowerVS connection status
	if in.PowerVSConnection.ID != "" {
		out.PowerVSConnection = &ResourceReference{
			ID: ptr.To(in.PowerVSConnection.ID),
		}
	}

	return nil
}

func Convert_v1beta2_VPCResourceReference_To_v1beta3_VPCSource(in *VPCResourceReference, out *infrav1.VPCSource, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	if in.ID != nil && *in.ID != "" {
		out.Type = infrav1.SourceTypeReference
		out.Reference.ID = *in.ID
		if in.Name != nil && *in.Name != "" {
			out.Reference.Name = *in.Name
		}
	} else {
		out.Type = infrav1.SourceTypeProvision
		if in.Name != nil && *in.Name != "" {
			out.Provision.Name = *in.Name
		}
	}

	if in.Region != nil && *in.Region != "" {
		out.Region = *in.Region
	}

	return nil
}

func Convert_v1beta3_VPCSource_To_v1beta2_VPCResourceReference(in *infrav1.VPCSource, out *VPCResourceReference, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	switch in.Type {
	case infrav1.SourceTypeReference:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		}
	default:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Provision.Name != "" {
			out.Name = ptr.To(in.Provision.Name)
		}
		if out.Name == nil && in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		}
	}

	if in.Region != "" {
		out.Region = ptr.To(in.Region)
	}

	return nil
}

// Convert_v1beta2_ResourceReference_To_v1beta3_ResourceReference converts a v1beta2 ResourceReference to a v1beta3 ResourceReference.
// ControllerCreated is dropped as it does not exist in v1beta3.
func Convert_v1beta2_ResourceReference_To_v1beta3_ResourceReference(in *ResourceReference, out *infrav1.ResourceReference, _ apimachineryconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
	}
	return nil
}

// Convert_v1beta3_ResourceReference_To_v1beta2_ResourceReference converts a v1beta3 ResourceReference to a v1beta2 ResourceReference.
// Name is dropped as it does not exist in v1beta2.
func Convert_v1beta3_ResourceReference_To_v1beta2_ResourceReference(in *infrav1.ResourceReference, out *ResourceReference, _ apimachineryconversion.Scope) error {
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	}
	return nil
}

func Convert_v1beta2_ResourceReference_To_v1beta3_VPCStatus(in *ResourceReference, out *infrav1.VPCStatus, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}
	if in.ID != nil && *in.ID != "" {
		out.ID = *in.ID
	}
	return nil
}

func Convert_v1beta3_VPCStatus_To_v1beta2_ResourceReference(in *infrav1.VPCStatus, out *ResourceReference, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	}
	return nil
}

func Convert_v1beta2_VPCLoadBalancerStatus_To_v1beta3_LoadBalancerStatus(in *VPCLoadBalancerStatus, out *infrav1.LoadBalancerStatus, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}
	if in.ID != nil && *in.ID != "" {
		out.ID = *in.ID
	}
	out.State = infrav1.LoadBalancerState(in.State)
	if in.Hostname != nil && *in.Hostname != "" {
		out.Hostname = *in.Hostname
	}
	return nil
}

func Convert_v1beta3_LoadBalancerStatus_To_v1beta2_VPCLoadBalancerStatus(in *infrav1.LoadBalancerStatus, out *VPCLoadBalancerStatus, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	}
	out.State = VPCLoadBalancerState(in.State)
	if in.Hostname != "" {
		out.Hostname = ptr.To(in.Hostname)
	}
	return nil
}

// Convert_v1beta2_Subnet_To_v1beta3_VPCSubnetSource converts v1beta2 Subnet to v1beta3 VPCSubnetSource.
func Convert_v1beta2_Subnet_To_v1beta3_VPCSubnetSource(in *Subnet, out *infrav1.VPCSubnetSource, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	if in.ID != nil && *in.ID != "" {
		out.Type = infrav1.SourceTypeReference
		out.Reference.ID = *in.ID
		if in.Name != nil && *in.Name != "" {
			out.Reference.Name = *in.Name
		}
	} else if in.Name != nil && *in.Name != "" {
		out.Type = infrav1.SourceTypeProvision
		out.Provision.Name = *in.Name
	} else {
		out.Type = infrav1.SourceTypeProvision
	}

	if in.Zone != nil && *in.Zone != "" {
		out.Zone = *in.Zone
	}

	return nil
}

// Convert_v1beta3_VPCSubnetSource_To_v1beta2_Subnet converts v1beta3 VPCSubnetSource to v1beta2 Subnet.
func Convert_v1beta3_VPCSubnetSource_To_v1beta2_Subnet(in *infrav1.VPCSubnetSource, out *Subnet, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	if in.Zone != "" {
		out.Zone = ptr.To(in.Zone)
	}

	switch in.Type {
	case infrav1.SourceTypeReference:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		}
	case infrav1.SourceTypeProvision:
		if in.Provision.Name != "" {
			out.Name = ptr.To(in.Provision.Name)
		}
	default:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
			if in.Reference.Name != "" {
				out.Name = ptr.To(in.Reference.Name)
			}
		} else if in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		} else if in.Provision.Name != "" {
			out.Name = ptr.To(in.Provision.Name)
		}
	}

	return nil
}

// Convert_v1beta2_VPCLoadBalancerSpec_To_v1beta3_LoadBalancerSource converts v1beta2 VPCLoadBalancerSpec to v1beta3 LoadBalancerSource.
//
//nolint:gocyclo // complexity is acceptable for conversion logic
func Convert_v1beta2_VPCLoadBalancerSpec_To_v1beta3_LoadBalancerSource(in *VPCLoadBalancerSpec, out *infrav1.LoadBalancerSource, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	hasProvisionData := in.Public != nil || len(in.AdditionalListeners) > 0 || len(in.BackendPools) > 0 || len(in.SecurityGroups) > 0 || len(in.Subnets) > 0
	if hasProvisionData || (in.Name != "" && (in.ID == nil || *in.ID == "")) {
		out.Type = infrav1.SourceTypeProvision
		out.Provision.Name = in.Name

		if in.Public != nil {
			if *in.Public {
				out.Provision.Type = infrav1.LoadBalancerTypePublic
			} else {
				out.Provision.Type = infrav1.LoadBalancerTypePrivate
			}
		}

		if len(in.AdditionalListeners) > 0 {
			out.Provision.AdditionalListeners = make([]infrav1.AdditionalListener, len(in.AdditionalListeners))
			for i := range in.AdditionalListeners {
				listener := in.AdditionalListeners[i]
				out.Provision.AdditionalListeners[i] = infrav1.AdditionalListener{
					Port:     listener.Port,
					Selector: listener.Selector,
				}
				if listener.DefaultPoolName != nil {
					out.Provision.AdditionalListeners[i].DefaultPoolName = *listener.DefaultPoolName
				}
				if listener.Protocol != nil {
					out.Provision.AdditionalListeners[i].Protocol = infrav1.LoadBalancerListenerProtocol(*listener.Protocol)
				}
			}
		}

		if len(in.BackendPools) > 0 {
			out.Provision.BackendPools = make([]infrav1.LoadBalancerBackendPool, len(in.BackendPools))
			for i := range in.BackendPools {
				pool := in.BackendPools[i]
				out.Provision.BackendPools[i] = infrav1.LoadBalancerBackendPool{
					Algorithm: infrav1.LoadBalancerBackendPoolAlgorithm(pool.Algorithm),
					Protocol:  infrav1.LoadBalancerBackendPoolProtocol(pool.Protocol),
					HealthMonitor: infrav1.LoadBalancerHealthMonitor{
						Delay:   pool.HealthMonitor.Delay,
						Retries: pool.HealthMonitor.Retries,
						Timeout: pool.HealthMonitor.Timeout,
						Type:    infrav1.LoadBalancerBackendPoolHealthMonitorType(pool.HealthMonitor.Type),
					},
				}
				if pool.Name != nil {
					out.Provision.BackendPools[i].Name = *pool.Name
				}
				if pool.HealthMonitor.Port != nil {
					out.Provision.BackendPools[i].HealthMonitor.Port = *pool.HealthMonitor.Port
				}
				if pool.HealthMonitor.URLPath != nil {
					out.Provision.BackendPools[i].HealthMonitor.URLPath = *pool.HealthMonitor.URLPath
				}
			}
		}

		if len(in.SecurityGroups) > 0 {
			out.Provision.SecurityGroups = make([]infrav1.ResourceIdentifier, len(in.SecurityGroups))
			for i := range in.SecurityGroups {
				out.Provision.SecurityGroups[i] = infrav1.ResourceIdentifier{}
				if in.SecurityGroups[i].ID != nil && *in.SecurityGroups[i].ID != "" {
					out.Provision.SecurityGroups[i].ID = *in.SecurityGroups[i].ID
				}
				if in.SecurityGroups[i].Name != nil && *in.SecurityGroups[i].Name != "" {
					out.Provision.SecurityGroups[i].Name = *in.SecurityGroups[i].Name
				}
			}
		}

		if len(in.Subnets) > 0 {
			out.Provision.Subnets = make([]infrav1.ResourceIdentifier, len(in.Subnets))
			for i := range in.Subnets {
				out.Provision.Subnets[i] = infrav1.ResourceIdentifier{}
				if in.Subnets[i].ID != nil && *in.Subnets[i].ID != "" {
					out.Provision.Subnets[i].ID = *in.Subnets[i].ID
				}
				if in.Subnets[i].Name != nil && *in.Subnets[i].Name != "" {
					out.Provision.Subnets[i].Name = *in.Subnets[i].Name
				}
			}
		}

		return nil
	}

	if in.ID != nil && *in.ID != "" {
		out.Type = infrav1.SourceTypeReference
		out.Reference.ID = *in.ID
		if in.Name != "" {
			out.Reference.Name = in.Name
		}
		return nil
	}

	if in.Name != "" {
		out.Type = infrav1.SourceTypeProvision
		out.Provision.Name = in.Name
		return nil
	}

	return nil
}

// Convert_v1beta3_LoadBalancerSource_To_v1beta2_VPCLoadBalancerSpec converts v1beta3 LoadBalancerSource to v1beta2 VPCLoadBalancerSpec.
//
//nolint:gocyclo // complexity is acceptable for conversion logic
func Convert_v1beta3_LoadBalancerSource_To_v1beta2_VPCLoadBalancerSpec(in *infrav1.LoadBalancerSource, out *VPCLoadBalancerSpec, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	switch in.Type {
	case infrav1.SourceTypeReference:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Reference.Name != "" {
			out.Name = in.Reference.Name
		}
	case infrav1.SourceTypeProvision:
		if in.Provision.Name != "" {
			out.Name = in.Provision.Name
		}
		if in.Provision.Type != "" {
			out.Public = ptr.To(in.Provision.Type == infrav1.LoadBalancerTypePublic)
		}
		if len(in.Provision.AdditionalListeners) > 0 {
			out.AdditionalListeners = make([]AdditionalListenerSpec, len(in.Provision.AdditionalListeners))
			for i := range in.Provision.AdditionalListeners {
				listener := in.Provision.AdditionalListeners[i]
				out.AdditionalListeners[i] = AdditionalListenerSpec{
					Port:     listener.Port,
					Selector: listener.Selector,
				}
				if listener.DefaultPoolName != "" {
					out.AdditionalListeners[i].DefaultPoolName = ptr.To(listener.DefaultPoolName)
				}
				if listener.Protocol != "" {
					protocol := VPCLoadBalancerListenerProtocol(listener.Protocol)
					out.AdditionalListeners[i].Protocol = &protocol
				}
			}
		}
		if len(in.Provision.BackendPools) > 0 {
			out.BackendPools = make([]VPCLoadBalancerBackendPoolSpec, len(in.Provision.BackendPools))
			for i := range in.Provision.BackendPools {
				pool := in.Provision.BackendPools[i]
				out.BackendPools[i] = VPCLoadBalancerBackendPoolSpec{
					Algorithm: VPCLoadBalancerBackendPoolAlgorithm(pool.Algorithm),
					Protocol:  VPCLoadBalancerBackendPoolProtocol(pool.Protocol),
					HealthMonitor: VPCLoadBalancerHealthMonitorSpec{
						Delay:   pool.HealthMonitor.Delay,
						Retries: pool.HealthMonitor.Retries,
						Timeout: pool.HealthMonitor.Timeout,
						Type:    VPCLoadBalancerBackendPoolHealthMonitorType(pool.HealthMonitor.Type),
					},
				}
				if pool.Name != "" {
					out.BackendPools[i].Name = ptr.To(pool.Name)
				}
				if pool.HealthMonitor.Port != 0 {
					out.BackendPools[i].HealthMonitor.Port = ptr.To(pool.HealthMonitor.Port)
				}
				if pool.HealthMonitor.URLPath != "" {
					out.BackendPools[i].HealthMonitor.URLPath = ptr.To(pool.HealthMonitor.URLPath)
				}
			}
		}
		if len(in.Provision.SecurityGroups) > 0 {
			out.SecurityGroups = make([]VPCResource, len(in.Provision.SecurityGroups))
			for i := range in.Provision.SecurityGroups {
				if in.Provision.SecurityGroups[i].ID != "" {
					out.SecurityGroups[i].ID = ptr.To(in.Provision.SecurityGroups[i].ID)
				}
				if in.Provision.SecurityGroups[i].Name != "" {
					out.SecurityGroups[i].Name = ptr.To(in.Provision.SecurityGroups[i].Name)
				}
			}
		}
		if len(in.Provision.Subnets) > 0 {
			out.Subnets = make([]VPCResource, len(in.Provision.Subnets))
			for i := range in.Provision.Subnets {
				if in.Provision.Subnets[i].ID != "" {
					out.Subnets[i].ID = ptr.To(in.Provision.Subnets[i].ID)
				}
				if in.Provision.Subnets[i].Name != "" {
					out.Subnets[i].Name = ptr.To(in.Provision.Subnets[i].Name)
				}
			}
		}
	}

	return nil
}

// Convert_v1beta2_VPCSecurityGroup_To_v1beta3_VPCSecurityGroupSource converts v1beta2 VPCSecurityGroup to v1beta3 VPCSecurityGroupSource.
// v1beta2 VPCSecurityGroup has ID/Name pointers and a Rules slice; v1beta3 uses a Type/Reference/Provision union.
// If an ID is set the entry is treated as a Reference, otherwise as a Provision entry.
func Convert_v1beta2_VPCSecurityGroup_To_v1beta3_VPCSecurityGroupSource(in *VPCSecurityGroup, out *infrav1.VPCSecurityGroupSource, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	if in.ID != nil && *in.ID != "" {
		out.Type = infrav1.SourceTypeReference
		out.Reference.ID = *in.ID
		if in.Name != nil {
			out.Reference.Name = *in.Name
		}
		return nil
	}

	out.Type = infrav1.SourceTypeProvision
	if in.Name != nil {
		out.Provision.Name = *in.Name
	}
	for _, tag := range in.Tags {
		if tag != nil {
			out.Provision.Tags = append(out.Provision.Tags, *tag)
		}
	}
	for _, rule := range in.Rules {
		if rule == nil {
			continue
		}
		v3Rule := infrav1.VPCSecurityGroupRule{}
		if err := Convert_v1beta2_VPCSecurityGroupRule_To_v1beta3_VPCSecurityGroupRule(rule, &v3Rule, nil); err != nil {
			return err
		}
		out.Provision.Rules = append(out.Provision.Rules, v3Rule)
	}
	return nil
}

// Convert_v1beta3_VPCSecurityGroupSource_To_v1beta2_VPCSecurityGroup converts v1beta3 VPCSecurityGroupSource to v1beta2 VPCSecurityGroup.
func Convert_v1beta3_VPCSecurityGroupSource_To_v1beta2_VPCSecurityGroup(in *infrav1.VPCSecurityGroupSource, out *VPCSecurityGroup, _ apimachineryconversion.Scope) error {
	if in == nil {
		return nil
	}

	switch in.Type {
	case infrav1.SourceTypeReference:
		if in.Reference.ID != "" {
			out.ID = ptr.To(in.Reference.ID)
		}
		if in.Reference.Name != "" {
			out.Name = ptr.To(in.Reference.Name)
		}
	case infrav1.SourceTypeProvision:
		if in.Provision.Name != "" {
			out.Name = ptr.To(in.Provision.Name)
		}
		for _, tag := range in.Provision.Tags {
			t := tag
			out.Tags = append(out.Tags, &t)
		}
		for i := range in.Provision.Rules {
			v2Rule := &VPCSecurityGroupRule{}
			if err := Convert_v1beta3_VPCSecurityGroupRule_To_v1beta2_VPCSecurityGroupRule(&in.Provision.Rules[i], v2Rule, nil); err != nil {
				return err
			}
			out.Rules = append(out.Rules, v2Rule)
		}
	}
	return nil
}

// Convert_v1beta2_VPCSecurityGroupRule_To_v1beta3_VPCSecurityGroupRule converts v1beta2 VPCSecurityGroupRule to v1beta3.
// The key difference is that v1beta2 Destination/Source are pointers while v1beta3 uses value types.
func Convert_v1beta2_VPCSecurityGroupRule_To_v1beta3_VPCSecurityGroupRule(in *VPCSecurityGroupRule, out *infrav1.VPCSecurityGroupRule, _ apimachineryconversion.Scope) error {
	out.Direction = infrav1.VPCSecurityGroupRuleDirection(in.Direction)
	if in.SecurityGroupID != nil {
		out.SecurityGroupID = *in.SecurityGroupID
	}
	if in.Destination != nil {
		if err := Convert_v1beta2_VPCSecurityGroupRulePrototype_To_v1beta3_VPCSecurityGroupRulePrototype(in.Destination, &out.Destination, nil); err != nil {
			return err
		}
	}
	if in.Source != nil {
		if err := Convert_v1beta2_VPCSecurityGroupRulePrototype_To_v1beta3_VPCSecurityGroupRulePrototype(in.Source, &out.Source, nil); err != nil {
			return err
		}
	}
	return nil
}

// Convert_v1beta3_VPCSecurityGroupRule_To_v1beta2_VPCSecurityGroupRule converts v1beta3 VPCSecurityGroupRule to v1beta2.
func Convert_v1beta3_VPCSecurityGroupRule_To_v1beta2_VPCSecurityGroupRule(in *infrav1.VPCSecurityGroupRule, out *VPCSecurityGroupRule, _ apimachineryconversion.Scope) error {
	out.Direction = VPCSecurityGroupRuleDirection(in.Direction)
	if in.SecurityGroupID != "" {
		out.SecurityGroupID = ptr.To(in.SecurityGroupID)
	}
	if in.Direction == infrav1.VPCSecurityGroupRuleDirectionOutbound {
		out.Destination = &VPCSecurityGroupRulePrototype{}
		if err := Convert_v1beta3_VPCSecurityGroupRulePrototype_To_v1beta2_VPCSecurityGroupRulePrototype(&in.Destination, out.Destination, nil); err != nil {
			return err
		}
	}
	if in.Direction == infrav1.VPCSecurityGroupRuleDirectionInbound {
		out.Source = &VPCSecurityGroupRulePrototype{}
		if err := Convert_v1beta3_VPCSecurityGroupRulePrototype_To_v1beta2_VPCSecurityGroupRulePrototype(&in.Source, out.Source, nil); err != nil {
			return err
		}
	}
	return nil
}

// Convert_v1beta2_VPCSecurityGroupRulePrototype_To_v1beta3_VPCSecurityGroupRulePrototype converts v1beta2 VPCSecurityGroupRulePrototype to v1beta3.
// The key difference is that v1beta2 PortRange is a pointer while v1beta3 uses a value type with omitzero.
func Convert_v1beta2_VPCSecurityGroupRulePrototype_To_v1beta3_VPCSecurityGroupRulePrototype(in *VPCSecurityGroupRulePrototype, out *infrav1.VPCSecurityGroupRulePrototype, _ apimachineryconversion.Scope) error {
	out.ICMPCode = in.ICMPCode
	out.ICMPType = in.ICMPType
	out.Protocol = infrav1.VPCSecurityGroupRuleProtocol(in.Protocol)
	if in.PortRange != nil {
		out.PortRange = infrav1.VPCSecurityGroupPortRange{
			MaximumPort: in.PortRange.MaximumPort,
			MinimumPort: in.PortRange.MinimumPort,
		}
	}
	for i := range in.Remotes {
		v3Remote := infrav1.VPCSecurityGroupRuleRemote{
			RemoteType: infrav1.VPCSecurityGroupRuleRemoteType(in.Remotes[i].RemoteType),
		}
		if in.Remotes[i].CIDRSubnetName != nil {
			v3Remote.CIDRSubnetName = *in.Remotes[i].CIDRSubnetName
		}
		if in.Remotes[i].Address != nil {
			v3Remote.Address = *in.Remotes[i].Address
		}
		if in.Remotes[i].SecurityGroupName != nil {
			v3Remote.SecurityGroupName = *in.Remotes[i].SecurityGroupName
		}
		out.Remotes = append(out.Remotes, v3Remote)
	}
	return nil
}

// Convert_v1beta3_VPCSecurityGroupRulePrototype_To_v1beta2_VPCSecurityGroupRulePrototype converts v1beta3 VPCSecurityGroupRulePrototype to v1beta2.
func Convert_v1beta3_VPCSecurityGroupRulePrototype_To_v1beta2_VPCSecurityGroupRulePrototype(in *infrav1.VPCSecurityGroupRulePrototype, out *VPCSecurityGroupRulePrototype, _ apimachineryconversion.Scope) error {
	out.ICMPCode = in.ICMPCode
	out.ICMPType = in.ICMPType
	out.Protocol = VPCSecurityGroupRuleProtocol(in.Protocol)
	if in.PortRange.MaximumPort != 0 || in.PortRange.MinimumPort != 0 {
		out.PortRange = &VPCSecurityGroupPortRange{
			MaximumPort: in.PortRange.MaximumPort,
			MinimumPort: in.PortRange.MinimumPort,
		}
	}
	for i := range in.Remotes {
		v2Remote := VPCSecurityGroupRuleRemote{
			RemoteType: VPCSecurityGroupRuleRemoteType(in.Remotes[i].RemoteType),
		}
		if in.Remotes[i].CIDRSubnetName != "" {
			v2Remote.CIDRSubnetName = ptr.To(in.Remotes[i].CIDRSubnetName)
		}
		if in.Remotes[i].Address != "" {
			v2Remote.Address = ptr.To(in.Remotes[i].Address)
		}
		if in.Remotes[i].SecurityGroupName != "" {
			v2Remote.SecurityGroupName = ptr.To(in.Remotes[i].SecurityGroupName)
		}
		out.Remotes = append(out.Remotes, v2Remote)
	}
	return nil
}

// Convert_v1beta2_VPCSecurityGroupStatus_To_v1beta3_VPCSecurityGroupStatus converts a single v1beta2 VPCSecurityGroupStatus to v1beta3.
// v1beta2 tracks RuleIDs as []*string; v1beta3 tracks them as []VPCSecurityGroupRuleStatus.
// ControllerCreated is a v1beta2-only field and is dropped.
func Convert_v1beta2_VPCSecurityGroupStatus_To_v1beta3_VPCSecurityGroupStatus(in *VPCSecurityGroupStatus, out *infrav1.VPCSecurityGroupStatus, _ apimachineryconversion.Scope) error {
	if in.ID != nil {
		out.ID = *in.ID
	}
	for _, rid := range in.RuleIDs {
		if rid != nil {
			out.Rules = append(out.Rules, infrav1.VPCSecurityGroupRuleStatus{ID: *rid})
		}
	}
	return nil
}

// Convert_v1beta3_VPCSecurityGroupStatus_To_v1beta2_VPCSecurityGroupStatus converts a single v1beta3 VPCSecurityGroupStatus to v1beta2.
// Name and full rule details are v1beta3-only fields and are dropped.
func Convert_v1beta3_VPCSecurityGroupStatus_To_v1beta2_VPCSecurityGroupStatus(in *infrav1.VPCSecurityGroupStatus, out *VPCSecurityGroupStatus, _ apimachineryconversion.Scope) error {
	if in.ID != "" {
		out.ID = ptr.To(in.ID)
	}
	for i := range in.Rules {
		if in.Rules[i].ID != "" {
			out.RuleIDs = append(out.RuleIDs, ptr.To(in.Rules[i].ID))
		}
	}
	return nil
}
