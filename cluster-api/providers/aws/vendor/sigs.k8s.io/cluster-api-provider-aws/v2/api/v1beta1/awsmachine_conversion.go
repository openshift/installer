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
	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
	utilconversion "sigs.k8s.io/cluster-api/util/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

// ConvertTo converts the v1beta1 AWSMachine receiver to a v1beta2 AWSMachine.
func (src *AWSMachine) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachine)
	if err := Convert_v1beta1_AWSMachine_To_v1beta2_AWSMachine(src, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.AWSMachine{}
	if ok, err := utilconversion.UnmarshalData(src, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Ignition = restored.Spec.Ignition
	dst.Spec.InstanceMetadataOptions = restored.Spec.InstanceMetadataOptions
	dst.Spec.PlacementGroupName = restored.Spec.PlacementGroupName
	dst.Spec.PlacementGroupPartition = restored.Spec.PlacementGroupPartition
	dst.Spec.PrivateDNSName = restored.Spec.PrivateDNSName
	dst.Spec.SecurityGroupOverrides = restored.Spec.SecurityGroupOverrides
	dst.Spec.CapacityReservationID = restored.Spec.CapacityReservationID
	dst.Spec.MarketType = restored.Spec.MarketType
	dst.Spec.NetworkInterfaceType = restored.Spec.NetworkInterfaceType
	if restored.Spec.ElasticIPPool != nil {
		if dst.Spec.ElasticIPPool == nil {
			dst.Spec.ElasticIPPool = &infrav1.ElasticIPPool{}
		}
		if restored.Spec.ElasticIPPool.PublicIpv4Pool != nil {
			dst.Spec.ElasticIPPool.PublicIpv4Pool = restored.Spec.ElasticIPPool.PublicIpv4Pool
		}
		if restored.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder != nil {
			dst.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder = restored.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder
		}
	}

	return nil
}

// ConvertFrom converts the v1beta2 AWSMachine to a v1beta1 AWSMachine.
func (dst *AWSMachine) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachine)

	if err := Convert_v1beta2_AWSMachine_To_v1beta1_AWSMachine(src, dst, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion except for metadata.
	return utilconversion.MarshalData(src, dst)
}

// ConvertTo converts the v1beta1 AWSMachineList receiver to a v1beta2 AWSMachineList.
func (src *AWSMachineList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineList)
	return Convert_v1beta1_AWSMachineList_To_v1beta2_AWSMachineList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSMachineList to a v1beta1 AWSMachineList.
func (dst *AWSMachineList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineList)

	return Convert_v1beta2_AWSMachineList_To_v1beta1_AWSMachineList(src, dst, nil)
}

// ConvertTo converts the v1beta1 AWSMachineTemplate receiver to a v1beta2 AWSMachineTemplate.
func (r *AWSMachineTemplate) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineTemplate)

	if err := Convert_v1beta1_AWSMachineTemplate_To_v1beta2_AWSMachineTemplate(r, dst, nil); err != nil {
		return err
	}

	// Manually restore data.
	restored := &infrav1.AWSMachineTemplate{}
	if ok, err := utilconversion.UnmarshalData(r, restored); err != nil || !ok {
		return err
	}

	dst.Spec.Template.ObjectMeta = restored.Spec.Template.ObjectMeta
	dst.Spec.Template.Spec.Ignition = restored.Spec.Template.Spec.Ignition
	dst.Spec.Template.Spec.InstanceMetadataOptions = restored.Spec.Template.Spec.InstanceMetadataOptions
	dst.Spec.Template.Spec.PlacementGroupName = restored.Spec.Template.Spec.PlacementGroupName
	dst.Spec.Template.Spec.PlacementGroupPartition = restored.Spec.Template.Spec.PlacementGroupPartition
	dst.Spec.Template.Spec.PrivateDNSName = restored.Spec.Template.Spec.PrivateDNSName
	dst.Spec.Template.Spec.SecurityGroupOverrides = restored.Spec.Template.Spec.SecurityGroupOverrides
	dst.Spec.Template.Spec.CapacityReservationID = restored.Spec.Template.Spec.CapacityReservationID
	dst.Spec.Template.Spec.MarketType = restored.Spec.Template.Spec.MarketType
	dst.Spec.Template.Spec.NetworkInterfaceType = restored.Spec.Template.Spec.NetworkInterfaceType
	if restored.Spec.Template.Spec.ElasticIPPool != nil {
		if dst.Spec.Template.Spec.ElasticIPPool == nil {
			dst.Spec.Template.Spec.ElasticIPPool = &infrav1.ElasticIPPool{}
		}
		if restored.Spec.Template.Spec.ElasticIPPool.PublicIpv4Pool != nil {
			dst.Spec.Template.Spec.ElasticIPPool.PublicIpv4Pool = restored.Spec.Template.Spec.ElasticIPPool.PublicIpv4Pool
		}
		if restored.Spec.Template.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder != nil {
			dst.Spec.Template.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder = restored.Spec.Template.Spec.ElasticIPPool.PublicIpv4PoolFallBackOrder
		}
	}

	return nil
}

// ConvertFrom converts the v1beta2 AWSCluster receiver to a v1beta1 AWSCluster.
func (r *AWSMachineTemplate) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineTemplate)

	if err := Convert_v1beta2_AWSMachineTemplate_To_v1beta1_AWSMachineTemplate(src, r, nil); err != nil {
		return err
	}

	// Preserve Hub data on down-conversion.
	if err := utilconversion.MarshalData(src, r); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts the v1beta1 AWSMachineTemplateList receiver to a v1beta2 AWSMachineTemplateList.
func (src *AWSMachineTemplateList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*infrav1.AWSMachineTemplateList)
	return Convert_v1beta1_AWSMachineTemplateList_To_v1beta2_AWSMachineTemplateList(src, dst, nil)
}

// ConvertFrom converts the v1beta2 AWSMachineTemplateList to a v1beta1 AWSMachineTemplateList.
func (dst *AWSMachineTemplateList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*infrav1.AWSMachineTemplateList)

	return Convert_v1beta2_AWSMachineTemplateList_To_v1beta1_AWSMachineTemplateList(src, dst, nil)
}
