// Copyright (c) 2023 VMware, Inc. All Rights Reserved.
// SPDX-License-Identifier: Apache-2.0

package v1alpha2

// Hub marks VirtualMachineService as a conversion hub.
func (*VirtualMachineService) Hub() {}

// Hub marks VirtualMachineServiceList as a conversion hub.
func (*VirtualMachineServiceList) Hub() {}
