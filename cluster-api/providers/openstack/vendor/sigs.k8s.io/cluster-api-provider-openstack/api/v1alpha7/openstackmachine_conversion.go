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
	apiconversion "k8s.io/apimachinery/pkg/conversion"
	ctrlconversion "sigs.k8s.io/controller-runtime/pkg/conversion"

	infrav1 "sigs.k8s.io/cluster-api-provider-openstack/api/v1beta1"
	"sigs.k8s.io/cluster-api-provider-openstack/pkg/utils/conversion"
)

var _ ctrlconversion.Convertible = &OpenStackMachine{}

func (r *OpenStackMachine) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachine)

	return conversion.ConvertAndRestore(
		r, dst,
		Convert_v1alpha7_OpenStackMachine_To_v1beta1_OpenStackMachine, Convert_v1beta1_OpenStackMachine_To_v1alpha7_OpenStackMachine,
		v1alpha7OpenStackMachineRestorer, v1beta1OpenStackMachineRestorer,
	)
}

func (r *OpenStackMachine) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachine)

	return conversion.ConvertAndRestore(
		src, r,
		Convert_v1beta1_OpenStackMachine_To_v1alpha7_OpenStackMachine, Convert_v1alpha7_OpenStackMachine_To_v1beta1_OpenStackMachine,
		v1beta1OpenStackMachineRestorer, v1alpha7OpenStackMachineRestorer,
	)
}

var _ ctrlconversion.Convertible = &OpenStackMachineList{}

func (r *OpenStackMachineList) ConvertTo(dstRaw ctrlconversion.Hub) error {
	dst := dstRaw.(*infrav1.OpenStackMachineList)
	return Convert_v1alpha7_OpenStackMachineList_To_v1beta1_OpenStackMachineList(r, dst, nil)
}

func (r *OpenStackMachineList) ConvertFrom(srcRaw ctrlconversion.Hub) error {
	src := srcRaw.(*infrav1.OpenStackMachineList)
	return Convert_v1beta1_OpenStackMachineList_To_v1alpha7_OpenStackMachineList(src, r, nil)
}

/* Restorers */

var v1alpha7OpenStackMachineRestorer = conversion.RestorerFor[*OpenStackMachine]{
	"spec": conversion.HashedFieldRestorer(
		func(c *OpenStackMachine) *OpenStackMachineSpec {
			return &c.Spec
		},
		restorev1alpha7MachineSpec,
		conversion.HashedFilterField[*OpenStackMachine, OpenStackMachineSpec](func(s *OpenStackMachineSpec) *OpenStackMachineSpec {
			// Despite being spec fields, ProviderID and InstanceID
			// are both set by the machine controller. If these are
			// the only changes to the spec, we still want to
			// restore the rest of the spec to its original state.
			if s.ProviderID != nil || s.InstanceID != nil {
				f := *s
				f.ProviderID = nil
				f.InstanceID = nil
				return &f
			}
			return s
		}),
	),
}

var v1beta1OpenStackMachineRestorer = conversion.RestorerFor[*infrav1.OpenStackMachine]{
	"spec": conversion.HashedFieldRestorer(
		func(c *infrav1.OpenStackMachine) *infrav1.OpenStackMachineSpec {
			return &c.Spec
		},
		restorev1beta1MachineSpec,
	),
	"depresources": conversion.UnconditionalFieldRestorer(
		func(c *infrav1.OpenStackMachine) **infrav1.MachineResources {
			return &c.Status.Resources
		},
	),
	// No equivalent in v1alpha7
	"refresources": conversion.UnconditionalFieldRestorer(
		func(c *infrav1.OpenStackMachine) **infrav1.ResolvedMachineSpec {
			return &c.Status.Resolved
		},
	),
}

/* OpenStackMachine */

func Convert_v1alpha7_OpenStackMachine_To_v1beta1_OpenStackMachine(in *OpenStackMachine, out *infrav1.OpenStackMachine, s apiconversion.Scope) error {
	err := autoConvert_v1alpha7_OpenStackMachine_To_v1beta1_OpenStackMachine(in, out, s)
	if err != nil {
		return err
	}

	out.Status.InstanceID = in.Spec.InstanceID
	return nil
}

func Convert_v1beta1_OpenStackMachine_To_v1alpha7_OpenStackMachine(in *infrav1.OpenStackMachine, out *OpenStackMachine, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_OpenStackMachine_To_v1alpha7_OpenStackMachine(in, out, s)
	if err != nil {
		return err
	}

	out.Spec.InstanceID = in.Status.InstanceID
	return nil
}

/* OpenStackMachineSpec */

func restorev1alpha7MachineSpec(previous *OpenStackMachineSpec, dst *OpenStackMachineSpec) {
	dst.FloatingIP = previous.FloatingIP

	// Conversion to v1beta1 truncates keys and values to 255 characters
	for k, v := range previous.ServerMetadata {
		kd := k
		if len(k) > 255 {
			kd = k[:255]
		}

		vd := v
		if len(v) > 255 {
			vd = v[:255]
		}

		if kd != k || vd != v {
			if dst.ServerMetadata == nil {
				dst.ServerMetadata = make(map[string]string)
			}
			delete(dst.ServerMetadata, kd)
			dst.ServerMetadata[k] = v
		}
	}

	// Conversion to v1beta1 removes the Kind field
	dst.IdentityRef = previous.IdentityRef

	if len(dst.Ports) == len(previous.Ports) {
		for i := range dst.Ports {
			restorev1alpha7Port(&previous.Ports[i], &dst.Ports[i])
		}
	}

	if len(dst.SecurityGroups) == len(previous.SecurityGroups) {
		for i := range dst.SecurityGroups {
			restorev1alpha7SecurityGroupFilter(&previous.SecurityGroups[i], &dst.SecurityGroups[i])
		}
	}

	// Conversion to v1beta1 removes Image when ImageUUID is set
	if dst.Image == "" && previous.Image != "" {
		dst.Image = previous.Image
	}
}

func restorev1beta1MachineSpec(previous *infrav1.OpenStackMachineSpec, dst *infrav1.OpenStackMachineSpec) {
	if previous == nil || dst == nil {
		return
	}

	dst.ServerGroup = previous.ServerGroup
	dst.Image = previous.Image

	if len(dst.SecurityGroups) == len(previous.SecurityGroups) {
		for i := range dst.SecurityGroups {
			restorev1beta1SecurityGroupParam(&previous.SecurityGroups[i], &dst.SecurityGroups[i])
		}
	}

	if len(dst.Ports) == len(previous.Ports) {
		for i := range dst.Ports {
			restorev1beta1Port(&previous.Ports[i], &dst.Ports[i])
		}
	}
	dst.FloatingIPPoolRef = previous.FloatingIPPoolRef
	dst.SchedulerHintAdditionalProperties = previous.SchedulerHintAdditionalProperties

	if dst.RootVolume != nil && previous.RootVolume != nil {
		restorev1beta1BlockDeviceVolume(
			&previous.RootVolume.BlockDeviceVolume,
			&dst.RootVolume.BlockDeviceVolume,
		)
	}

	if len(dst.AdditionalBlockDevices) == len(previous.AdditionalBlockDevices) {
		for i := range dst.AdditionalBlockDevices {
			restorev1beta1BlockDeviceVolume(
				previous.AdditionalBlockDevices[i].Storage.Volume,
				dst.AdditionalBlockDevices[i].Storage.Volume,
			)
		}
	}
}

func Convert_v1alpha7_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(in *OpenStackMachineSpec, out *infrav1.OpenStackMachineSpec, s apiconversion.Scope) error {
	err := autoConvert_v1alpha7_OpenStackMachineSpec_To_v1beta1_OpenStackMachineSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.ServerGroupID != "" {
		out.ServerGroup = &infrav1.ServerGroupParam{ID: &in.ServerGroupID}
	} else {
		out.ServerGroup = nil
	}

	imageParam := infrav1.ImageParam{}
	if in.ImageUUID != "" {
		imageParam.ID = &in.ImageUUID
	} else if in.Image != "" { // Only add name when ID is not set, in v1beta1 it's not possible to set both.
		imageParam.Filter = &infrav1.ImageFilter{Name: &in.Image}
	}
	out.Image = imageParam

	if len(in.ServerMetadata) > 0 {
		serverMetadata := make([]infrav1.ServerMetadata, 0, len(in.ServerMetadata))
		for k, v := range in.ServerMetadata {
			// Truncate key and value to 255 characters if required, as this
			// was not validated prior to v1beta1
			if len(k) > 255 {
				k = k[:255]
			}
			if len(v) > 255 {
				v = v[:255]
			}

			serverMetadata = append(serverMetadata, infrav1.ServerMetadata{Key: k, Value: v})
		}
		out.ServerMetadata = serverMetadata
	}

	if in.CloudName != "" {
		if out.IdentityRef == nil {
			out.IdentityRef = &infrav1.OpenStackIdentityReference{}
		}
		out.IdentityRef.CloudName = in.CloudName
	}

	return nil
}

func Convert_v1beta1_OpenStackMachineSpec_To_v1alpha7_OpenStackMachineSpec(in *infrav1.OpenStackMachineSpec, out *OpenStackMachineSpec, s apiconversion.Scope) error {
	err := autoConvert_v1beta1_OpenStackMachineSpec_To_v1alpha7_OpenStackMachineSpec(in, out, s)
	if err != nil {
		return err
	}

	if in.ServerGroup != nil && in.ServerGroup.ID != nil {
		out.ServerGroupID = *in.ServerGroup.ID
	}

	if in.Image.ID != nil {
		out.ImageUUID = *in.Image.ID
	} else if in.Image.Filter != nil && in.Image.Filter.Name != nil {
		out.Image = *in.Image.Filter.Name
	}

	if len(in.ServerMetadata) > 0 {
		serverMetadata := make(map[string]string, len(in.ServerMetadata))
		for i := range in.ServerMetadata {
			key := in.ServerMetadata[i].Key
			value := in.ServerMetadata[i].Value
			serverMetadata[key] = value
		}
		out.ServerMetadata = serverMetadata
	}

	if in.IdentityRef != nil {
		out.CloudName = in.IdentityRef.CloudName
	}

	return nil
}

/* OpenStackMachineStatus */

func Convert_v1beta1_OpenStackMachineStatus_To_v1alpha7_OpenStackMachineStatus(in *infrav1.OpenStackMachineStatus, out *OpenStackMachineStatus, s apiconversion.Scope) error {
	// ReferencedResources have no equivalent in v1alpha7
	return autoConvert_v1beta1_OpenStackMachineStatus_To_v1alpha7_OpenStackMachineStatus(in, out, s)
}
