// © Broadcom. All Rights Reserved.
// The term “Broadcom” refers to Broadcom Inc. and/or its subsidiaries.
// SPDX-License-Identifier: Apache-2.0

package v1alpha6

// Hub marks VirtualMachine as a conversion hub.
func (*VirtualMachine) Hub() {}

// Hub marks VirtualMachineList as a conversion hub.
func (*VirtualMachineList) Hub() {}
