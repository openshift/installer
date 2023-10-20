// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

// Hub marks VirtualMachineClass as a conversion hub.
func (*VirtualMachineClass) Hub() {}

// Hub marks VirtualMachineClassList as a conversion hub.
func (*VirtualMachineClassList) Hub() {}
