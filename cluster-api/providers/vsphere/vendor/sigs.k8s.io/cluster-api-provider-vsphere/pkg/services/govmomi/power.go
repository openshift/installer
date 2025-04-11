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

package govmomi

import (
	"context"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
)

func (vms *VMService) getPowerState(ctx context.Context, virtualMachineCtx *virtualMachineContext) (infrav1.VirtualMachinePowerState, error) {
	powerState, err := virtualMachineCtx.Obj.PowerState(ctx)
	if err != nil {
		return "", err
	}

	switch powerState {
	case types.VirtualMachinePowerStatePoweredOn:
		return infrav1.VirtualMachinePowerStatePoweredOn, nil
	case types.VirtualMachinePowerStatePoweredOff:
		return infrav1.VirtualMachinePowerStatePoweredOff, nil
	case types.VirtualMachinePowerStateSuspended:
		return infrav1.VirtualMachinePowerStateSuspended, nil
	default:
		return "", errors.Errorf("unexpected power state %q for vm %s", powerState, virtualMachineCtx)
	}
}

func (vms *VMService) isSoftPowerOffTimeoutExceeded(vm *infrav1.VSphereVM) bool {
	if !conditions.Has(vm, infrav1.GuestSoftPowerOffSucceededCondition) {
		// The SoftPowerOff never got triggered, so it can't be timed out yet.
		return false
	}
	if vm.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeSoft {
		// Timeout only applies to trySoft mode.
		// For soft mode it will wait infinitely.
		return false
	}
	now := time.Now()
	timeSoftPowerOff := conditions.GetLastTransitionTime(vm, infrav1.GuestSoftPowerOffSucceededCondition)
	diff := now.Sub(timeSoftPowerOff.Time)
	var timeout time.Duration
	if vm.Spec.GuestSoftPowerOffTimeout != nil {
		timeout = vm.Spec.GuestSoftPowerOffTimeout.Duration
	} else {
		timeout = infrav1.GuestSoftPowerOffDefaultTimeout
	}
	return timeout.Seconds() > 0 && diff.Seconds() >= timeout.Seconds()
}

// triggerSoftPowerOff tries to trigger a soft power off for a VM to shut down the guest.
// It returns true if the soft power off operation is pending.
func (vms *VMService) triggerSoftPowerOff(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	if virtualMachineCtx.VSphereVM.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeHard {
		// hard power off is expected.
		return false, nil
	}

	if conditions.Has(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition) {
		// soft power off operation has been triggered.
		if virtualMachineCtx.VSphereVM.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeSoft {
			return true, nil
		}

		return !vms.isSoftPowerOffTimeoutExceeded(virtualMachineCtx.VSphereVM), nil
	}

	vmwareToolsRunning, err := virtualMachineCtx.Obj.IsToolsRunning(ctx)
	if err != nil {
		return false, err
	}

	if !vmwareToolsRunning {
		// VMware tools is not installed.
		if virtualMachineCtx.VSphereVM.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeTrySoft {
			// Returning false to force a power off.
			return false, nil
		}

		conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition, infrav1.GuestSoftPowerOffFailedReason, clusterv1.ConditionSeverityWarning,
			"VMware Tools not installed on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
		// we are not able to trigger the soft power off so returning true to wait infinitely
		return true, nil
	}

	var o mo.VirtualMachine
	if err := virtualMachineCtx.Obj.Properties(ctx, virtualMachineCtx.Obj.Reference(), []string{"guest.guestStateChangeSupported"}, &o); err != nil {
		return false, err
	}

	if o.Guest.GuestStateChangeSupported == nil || !*o.Guest.GuestStateChangeSupported {
		if virtualMachineCtx.VSphereVM.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeTrySoft {
			// Returning false to force a power off.
			return false, nil
		}

		conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition, infrav1.GuestSoftPowerOffFailedReason, clusterv1.ConditionSeverityWarning,
			"unable to trigger soft power off because guest state change is not supported on VM %s.", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
		// we are not able to trigger the soft power off so returning true to wait infinitely
		return true, nil
	}

	err = virtualMachineCtx.Obj.ShutdownGuest(ctx)
	if err != nil {
		return false, err
	}

	conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededCondition, infrav1.GuestSoftPowerOffInProgressReason, clusterv1.ConditionSeverityInfo,
		"guest soft power off initiated on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
	return true, nil
}
