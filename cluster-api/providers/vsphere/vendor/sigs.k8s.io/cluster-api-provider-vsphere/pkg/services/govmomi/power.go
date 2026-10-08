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
	"fmt"
	"time"

	pkgerrors "github.com/pkg/errors"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/core/v1beta2"
	"sigs.k8s.io/cluster-api/util/conditions"
	deprecatedv1beta1conditions "sigs.k8s.io/cluster-api/util/conditions/deprecated/v1beta1"
	"sigs.k8s.io/controller-runtime/pkg/client"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/api/govmomi/v1beta2"
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
		return "", pkgerrors.Errorf("unexpected power state %q for vm %s", powerState, virtualMachineCtx)
	}
}

func (vms *VMService) isSoftPowerOffTimeoutExceeded(vm *infrav1.VSphereVM) bool {
	if !conditions.Has(vm, infrav1.VSphereVMGuestSoftPowerOffSucceededCondition) {
		// The SoftPowerOff never got triggered, so it can't be timed out yet.
		return false
	}
	if vm.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeSoft {
		// Timeout only applies to trySoft mode.
		// For soft mode it will wait infinitely.
		return false
	}
	now := time.Now()
	timeSoftPowerOff := conditions.GetLastTransitionTime(vm, infrav1.VSphereVMGuestSoftPowerOffSucceededCondition)
	diff := now.Sub(timeSoftPowerOff.Time)
	var timeout time.Duration
	if vm.Spec.GuestSoftPowerOffTimeoutSeconds != 0 {
		timeout = time.Duration(vm.Spec.GuestSoftPowerOffTimeoutSeconds) * time.Second
	} else {
		timeout = infrav1.GuestSoftPowerOffDefaultTimeoutSeconds * time.Second
	}
	return timeout.Seconds() > 0 && diff.Seconds() >= timeout.Seconds()
}

// triggerSoftPowerOff tries to trigger a soft power off for a VM to shut down the guest.
// It returns true if the soft power off operation is pending.
func (vms *VMService) triggerSoftPowerOff(ctx context.Context, virtualMachineCtx *virtualMachineContext) (bool, error) {
	if virtualMachineCtx.VSphereVM.Spec.PowerOffMode == "" || // hard is default
		virtualMachineCtx.VSphereVM.Spec.PowerOffMode == infrav1.VirtualMachinePowerOpModeHard {
		// hard power off is expected.
		return false, nil
	}

	if conditions.Has(virtualMachineCtx.VSphereVM, infrav1.VSphereVMGuestSoftPowerOffSucceededCondition) {
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

		deprecatedv1beta1conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededV1Beta1Condition, infrav1.GuestSoftPowerOffFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning,
			"VMware Tools not installed on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
		conditions.Set(virtualMachineCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMGuestSoftPowerOffSucceededCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMGuestSoftPowerOffFailedReason,
			Message: fmt.Sprintf("VMware Tools not installed on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM)),
		})
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

		deprecatedv1beta1conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededV1Beta1Condition, infrav1.GuestSoftPowerOffFailedV1Beta1Reason, clusterv1.ConditionSeverityWarning,
			"unable to trigger soft power off because guest state change is not supported on VM %s.", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
		conditions.Set(virtualMachineCtx.VSphereVM, metav1.Condition{
			Type:    infrav1.VSphereVMGuestSoftPowerOffSucceededCondition,
			Status:  metav1.ConditionFalse,
			Reason:  infrav1.VSphereVMGuestSoftPowerOffFailedReason,
			Message: fmt.Sprintf("Unable to trigger soft power off because guest state change is not supported on VM %s.", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM)),
		})
		// we are not able to trigger the soft power off so returning true to wait infinitely
		return true, nil
	}

	err = virtualMachineCtx.Obj.ShutdownGuest(ctx)
	if err != nil {
		return false, err
	}

	deprecatedv1beta1conditions.MarkFalse(virtualMachineCtx.VSphereVM, infrav1.GuestSoftPowerOffSucceededV1Beta1Condition, infrav1.GuestSoftPowerOffInProgressV1Beta1Reason, clusterv1.ConditionSeverityInfo,
		"guest soft power off initiated on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM))
	conditions.Set(virtualMachineCtx.VSphereVM, metav1.Condition{
		Type:    infrav1.VSphereVMGuestSoftPowerOffSucceededCondition,
		Status:  metav1.ConditionFalse,
		Reason:  infrav1.VSphereVMGuestSoftPowerOffInProgressReason,
		Message: fmt.Sprintf("guest soft power off initiated on VM %s", client.ObjectKeyFromObject(virtualMachineCtx.VSphereVM)),
	})
	return true, nil
}
