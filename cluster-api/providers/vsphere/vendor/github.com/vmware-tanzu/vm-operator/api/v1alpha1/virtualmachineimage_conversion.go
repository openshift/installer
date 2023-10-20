// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha1

import (
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	"github.com/vmware-tanzu/vm-operator/api/v1alpha2/common"
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	"sigs.k8s.io/controller-runtime/pkg/conversion"
)

func Convert_v1alpha1_ContentProviderReference_To_common_LocalObjectRef(
	in *ContentProviderReference, out *common.LocalObjectRef, s apiconversion.Scope) error {

	out.APIVersion = in.APIVersion
	out.Kind = in.Kind
	out.Name = in.Name

	return nil
}

func Convert_common_LocalObjectRef_To_v1alpha1_ContentProviderReference(
	in *common.LocalObjectRef, out *ContentProviderReference, s apiconversion.Scope) error {

	out.APIVersion = in.APIVersion
	out.Kind = in.Kind
	out.Name = in.Name
	// out.Namespace =

	return nil
}

func Convert_v1alpha2_VirtualMachineImageOSInfo_To_v1alpha1_VirtualMachineImageOSInfo(
	in *v1alpha2.VirtualMachineImageOSInfo, out *VirtualMachineImageOSInfo, s apiconversion.Scope) error {

	// = in.ID
	out.Type = in.Type
	out.Version = in.Version

	return nil
}

func convert_v1alpha1_VirtualMachineImageOSInfo_To_v1alpha2_VirtualMachineImageOSInfo(
	in *VirtualMachineImageOSInfo, out *v1alpha2.VirtualMachineImageOSInfo, s apiconversion.Scope) error {

	out.Type = in.Type
	out.Version = in.Version

	return nil
}

func convert_v1alpah2_VirtualMachineImage_OVFProperties_To_v1alpha1_VirtualMachineImage_OVFEnv(
	in []v1alpha2.OVFProperty, out *map[string]OvfProperty, s apiconversion.Scope) error {

	if in != nil {
		*out = map[string]OvfProperty{}
		for _, p := range in {
			(*out)[p.Key] = OvfProperty{
				Key:         p.Key,
				Type:        p.Type,
				Default:     p.Default,
				Description: "",
				Label:       "",
			}
		}
	}

	return nil
}

func convert_v1alpha1_VirtualMachineImage_OVFEnv_To_v1alpha2_VirtualMachineImage_OVFProperties(
	in map[string]OvfProperty, out *[]v1alpha2.OVFProperty, s apiconversion.Scope) error {

	if in != nil {
		*out = make([]v1alpha2.OVFProperty, 0, len(in))
		for _, v := range in {
			*out = append(*out, v1alpha2.OVFProperty{
				Key:     v.Key,
				Type:    v.Type,
				Default: v.Default,
			})
		}
	}

	return nil
}

func Convert_v1alpha1_VirtualMachineImageSpec_To_v1alpha2_VirtualMachineImageSpec(
	in *VirtualMachineImageSpec, out *v1alpha2.VirtualMachineImageSpec, s apiconversion.Scope) error {

	// in.Type
	// in.ImageSourceType
	// in.ImageID
	// in.ProductInfo
	// in.OSInfo
	// in.OVFEnv
	// in.HardwareVersion

	return autoConvert_v1alpha1_VirtualMachineImageSpec_To_v1alpha2_VirtualMachineImageSpec(in, out, s)
}

func Convert_v1alpha1_VirtualMachineImageStatus_To_v1alpha2_VirtualMachineImageStatus(
	in *VirtualMachineImageStatus, out *v1alpha2.VirtualMachineImageStatus, s apiconversion.Scope) error {

	out.Name = in.ImageName
	out.ProviderContentVersion = in.ContentVersion
	// in.ImageSupported
	// in.ContentLibraryRef

	// Deprecated:
	// in.Uuid
	// in.PowerState
	// in.InternalId

	return autoConvert_v1alpha1_VirtualMachineImageStatus_To_v1alpha2_VirtualMachineImageStatus(in, out, s)
}

func Convert_v1alpha2_VirtualMachineImageStatus_To_v1alpha1_VirtualMachineImageStatus(
	in *v1alpha2.VirtualMachineImageStatus, out *VirtualMachineImageStatus, s apiconversion.Scope) error {

	out.ImageName = in.Name
	out.ContentVersion = in.ProviderContentVersion
	// out.ContentLibraryRef =

	// in.Capabilities

	// Deprecated:
	out.Uuid = ""
	out.InternalId = ""
	out.PowerState = ""

	return autoConvert_v1alpha2_VirtualMachineImageStatus_To_v1alpha1_VirtualMachineImageStatus(in, out, s)
}

func convert_v1alpha1_VirtualMachineImageSpec_To_v1alpha2_VirtualMachineImageStatus(
	in *VirtualMachineImageSpec, out *v1alpha2.VirtualMachineImageStatus, s apiconversion.Scope) error {

	// Some fields of the v1a1 ImageSpec moved into the v1a2 ImageStatus.
	// conversion-gen doesn't handle that so do those here.

	if in.HardwareVersion != 0 {
		out.HardwareVersion = &in.HardwareVersion
	}

	if err := convert_v1alpha1_VirtualMachineImageOSInfo_To_v1alpha2_VirtualMachineImageOSInfo(&in.OSInfo, &out.OSInfo, s); err != nil {
		return err
	}

	if err := convert_v1alpha1_VirtualMachineImage_OVFEnv_To_v1alpha2_VirtualMachineImage_OVFProperties(in.OVFEnv, &out.OVFProperties, s); err != nil {
		return err
	}

	if err := Convert_v1alpha1_VirtualMachineImageProductInfo_To_v1alpha2_VirtualMachineImageProductInfo(&in.ProductInfo, &out.ProductInfo, s); err != nil {
		return err
	}

	return nil
}

func convert_v1alpha2_VirtualMachineImageStatus_To_v1alpha1_VirtualMachineImageSpec(
	in *v1alpha2.VirtualMachineImageStatus, out *VirtualMachineImageSpec, s apiconversion.Scope) error {

	// Some fields of the v1a1 ImageSpec moved into the v1a2 ImageStatus.
	// conversion-gen doesn't handle that so do those here.

	if in.HardwareVersion != nil {
		out.HardwareVersion = *in.HardwareVersion
	}

	if err := Convert_v1alpha2_VirtualMachineImageOSInfo_To_v1alpha1_VirtualMachineImageOSInfo(&in.OSInfo, &out.OSInfo, s); err != nil {
		return err
	}

	if err := convert_v1alpah2_VirtualMachineImage_OVFProperties_To_v1alpha1_VirtualMachineImage_OVFEnv(in.OVFProperties, &out.OVFEnv, s); err != nil {
		return err
	}

	if err := Convert_v1alpha2_VirtualMachineImageProductInfo_To_v1alpha1_VirtualMachineImageProductInfo(&in.ProductInfo, &out.ProductInfo, s); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this VirtualMachineImage to the Hub version.
func (src *VirtualMachineImage) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachineImage)
	if err := Convert_v1alpha1_VirtualMachineImage_To_v1alpha2_VirtualMachineImage(src, dst, nil); err != nil {
		return err
	}

	if err := convert_v1alpha1_VirtualMachineImageSpec_To_v1alpha2_VirtualMachineImageStatus(&src.Spec, &dst.Status, nil); err != nil {
		return err
	}

	return nil
}

// ConvertFrom converts the hub version to this VirtualMachineImage.
func (dst *VirtualMachineImage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachineImage)
	if err := Convert_v1alpha2_VirtualMachineImage_To_v1alpha1_VirtualMachineImage(src, dst, nil); err != nil {
		return err
	}

	if err := convert_v1alpha2_VirtualMachineImageStatus_To_v1alpha1_VirtualMachineImageSpec(&src.Status, &dst.Spec, nil); err != nil {
		return err
	}

	return nil
}

// ConvertTo converts this VirtualMachineImageList to the Hub version.
func (src *VirtualMachineImageList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.VirtualMachineImageList)
	return Convert_v1alpha1_VirtualMachineImageList_To_v1alpha2_VirtualMachineImageList(src, dst, nil)
}

// ConvertFrom converts the hub version to this VirtualMachineImageList.
func (dst *VirtualMachineImageList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.VirtualMachineImageList)
	return Convert_v1alpha2_VirtualMachineImageList_To_v1alpha1_VirtualMachineImageList(src, dst, nil)
}

// ConvertTo converts this ClusterVirtualMachineImage to the Hub version.
func (src *ClusterVirtualMachineImage) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.ClusterVirtualMachineImage)
	return Convert_v1alpha1_ClusterVirtualMachineImage_To_v1alpha2_ClusterVirtualMachineImage(src, dst, nil)
}

// ConvertFrom converts the hub version to this ClusterVirtualMachineImage.
func (dst *ClusterVirtualMachineImage) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.ClusterVirtualMachineImage)
	return Convert_v1alpha2_ClusterVirtualMachineImage_To_v1alpha1_ClusterVirtualMachineImage(src, dst, nil)
}

// ConvertTo converts this ClusterVirtualMachineImageList to the Hub version.
func (src *ClusterVirtualMachineImageList) ConvertTo(dstRaw conversion.Hub) error {
	dst := dstRaw.(*v1alpha2.ClusterVirtualMachineImageList)
	return Convert_v1alpha1_ClusterVirtualMachineImageList_To_v1alpha2_ClusterVirtualMachineImageList(src, dst, nil)
}

// ConvertFrom converts the hub version to this ClusterVirtualMachineImageList.
func (dst *ClusterVirtualMachineImageList) ConvertFrom(srcRaw conversion.Hub) error {
	src := srcRaw.(*v1alpha2.ClusterVirtualMachineImageList)
	return Convert_v1alpha2_ClusterVirtualMachineImageList_To_v1alpha1_ClusterVirtualMachineImageList(src, dst, nil)
}
