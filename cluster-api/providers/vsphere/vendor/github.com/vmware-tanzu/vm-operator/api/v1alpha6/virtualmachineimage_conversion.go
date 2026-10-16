// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

// Hub marks VirtualMachineImage as a conversion hub.
func (*VirtualMachineImage) Hub() {}

// Hub marks VirtualMachineImageList as a conversion hub.
func (*VirtualMachineImageList) Hub() {}

// Hub marks ClusterVirtualMachineImage as a conversion hub.
func (*ClusterVirtualMachineImage) Hub() {}

// Hub marks ClusterVirtualMachineImageList as a conversion hub.
func (*ClusterVirtualMachineImageList) Hub() {}
