// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

// Hub marks VirtualMachineImage as a conversion hub.
func (*VirtualMachineImage) Hub() {}

// Hub marks VirtualMachineImageList as a conversion hub.
func (*VirtualMachineImageList) Hub() {}

// Hub marks ClusterVirtualMachineImage as a conversion hub.
func (*ClusterVirtualMachineImage) Hub() {}

// Hub marks ClusterVirtualMachineImageList as a conversion hub.
func (*ClusterVirtualMachineImageList) Hub() {}
