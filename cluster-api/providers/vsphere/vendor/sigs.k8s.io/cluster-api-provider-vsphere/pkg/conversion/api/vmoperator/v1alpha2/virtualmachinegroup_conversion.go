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

package v1alpha2

import (
	"context"

	vmoprv1alpha2 "github.com/vmware-tanzu/vm-operator/api/v1alpha2"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	"k8s.io/utils/ptr"

	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion"
	vmoprvhub "sigs.k8s.io/cluster-api-provider-vsphere/pkg/conversion/api/vmoperator/hub"
)

func convert_v1alpha2_VirtualMachineGroup_To_hub_VirtualMachineGroup(_ context.Context, src *vmoprv1alpha2.VirtualMachineGroup, dst *vmoprvhub.VirtualMachineGroup) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.BootOrder != nil {
		dst.Spec.BootOrder = []vmoprvhub.VirtualMachineGroupBootOrderGroup{}
		for _, bootOrderGroup := range src.Spec.BootOrder {
			bg := vmoprvhub.VirtualMachineGroupBootOrderGroup{}
			if bootOrderGroup.Members != nil {
				bg.Members = []vmoprvhub.GroupMember{}
				for _, member := range bootOrderGroup.Members {
					bg.Members = append(bg.Members, vmoprvhub.GroupMember{
						Name: member.Name,
						Kind: member.Kind,
					})
				}
			}
			bg.PowerOnDelay = bootOrderGroup.PowerOnDelay
			dst.Spec.BootOrder = append(dst.Spec.BootOrder, bg)
		}
	}
	if src.Status.Members != nil {
		dst.Status.Members = []vmoprvhub.VirtualMachineGroupMemberStatus{}
		for _, member := range src.Status.Members {
			m := vmoprvhub.VirtualMachineGroupMemberStatus{
				Name: member.Name,
				Kind: member.Kind,
				UID:  member.UID,
			}
			if member.Placement != nil {
				m.Placement = &vmoprvhub.VirtualMachinePlacementStatus{
					Zone: member.Placement.Zone,
					Node: member.Placement.Node,
					Pool: member.Placement.Pool,
				}
				if member.Placement.Datastores != nil {
					m.Placement.Datastores = []vmoprvhub.VirtualMachineGroupPlacementDatastoreStatus{}
					for _, datastore := range member.Placement.Datastores {
						d := vmoprvhub.VirtualMachineGroupPlacementDatastoreStatus{
							Name:                 datastore.Name,
							ID:                   datastore.ID,
							URL:                  datastore.URL,
							SupportedDiskFormats: datastore.SupportedDiskFormats,
						}
						if datastore.DiskKey != nil {
							d.DiskKey = ptr.To(*datastore.DiskKey)
						}
						m.Placement.Datastores = append(m.Placement.Datastores, d)
					}
				}
			}
			if member.PowerState != nil {
				m.PowerState = ptr.To(vmoprvhub.VirtualMachinePowerState(*member.PowerState))
			}
			if member.Conditions != nil {
				m.Conditions = []metav1.Condition{}
				for _, condition := range member.Conditions {
					m.Conditions = append(m.Conditions, condition)
				}
			}
			dst.Status.Members = append(dst.Status.Members, m)
		}
	}

	return nil
}

func convert_hub_VirtualMachineGroup_To_v1alpha2_VirtualMachineGroup(_ context.Context, src *vmoprvhub.VirtualMachineGroup, dst *vmoprv1alpha2.VirtualMachineGroup) error {
	dst.ObjectMeta = src.ObjectMeta

	if src.Spec.BootOrder != nil {
		dst.Spec.BootOrder = []vmoprv1alpha2.VirtualMachineGroupBootOrderGroup{}
		for _, bootOrderGroup := range src.Spec.BootOrder {
			bg := vmoprv1alpha2.VirtualMachineGroupBootOrderGroup{}
			if bootOrderGroup.Members != nil {
				bg.Members = []vmoprv1alpha2.GroupMember{}
				for _, member := range bootOrderGroup.Members {
					bg.Members = append(bg.Members, vmoprv1alpha2.GroupMember{
						Name: member.Name,
						Kind: member.Kind,
					})
				}
			}
			bg.PowerOnDelay = bootOrderGroup.PowerOnDelay
			dst.Spec.BootOrder = append(dst.Spec.BootOrder, bg)
		}
	}
	if src.Status.Members != nil {
		dst.Status.Members = []vmoprv1alpha2.VirtualMachineGroupMemberStatus{}
		for _, member := range src.Status.Members {
			m := vmoprv1alpha2.VirtualMachineGroupMemberStatus{
				Name: member.Name,
				Kind: member.Kind,
				UID:  member.UID,
			}
			if member.Placement != nil {
				m.Placement = &vmoprv1alpha2.VirtualMachinePlacementStatus{
					Zone: member.Placement.Zone,
					Node: member.Placement.Node,
					Pool: member.Placement.Pool,
				}
				if member.Placement.Datastores != nil {
					m.Placement.Datastores = []vmoprv1alpha2.VirtualMachineGroupPlacementDatastoreStatus{}
					for _, datastore := range member.Placement.Datastores {
						d := vmoprv1alpha2.VirtualMachineGroupPlacementDatastoreStatus{
							Name:                 datastore.Name,
							ID:                   datastore.ID,
							URL:                  datastore.URL,
							SupportedDiskFormats: datastore.SupportedDiskFormats,
						}
						if datastore.DiskKey != nil {
							d.DiskKey = ptr.To(*datastore.DiskKey)
						}
						m.Placement.Datastores = append(m.Placement.Datastores, d)
					}
				}
			}
			if member.PowerState != nil {
				m.PowerState = ptr.To(vmoprv1alpha2.VirtualMachinePowerState(*member.PowerState))
			}
			if member.Conditions != nil {
				m.Conditions = []metav1.Condition{}
				for _, condition := range member.Conditions {
					m.Conditions = append(m.Conditions, condition)
				}
			}
			dst.Status.Members = append(dst.Status.Members, m)
		}
	}

	return nil
}

func init() {
	converterBuilder.AddConversion(
		conversion.NewAddConversionBuilder(convert_hub_VirtualMachineGroup_To_v1alpha2_VirtualMachineGroup, convert_v1alpha2_VirtualMachineGroup_To_hub_VirtualMachineGroup),
	)
}
