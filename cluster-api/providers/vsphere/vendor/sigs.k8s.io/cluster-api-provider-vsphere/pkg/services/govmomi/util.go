/*
Copyright 2019 The Kubernetes Authors.

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
	gonet "net"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1beta1 "sigs.k8s.io/cluster-api/api/core/v1beta1"
	v1beta1conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions"
	v1beta2conditions "sigs.k8s.io/cluster-api/util/deprecated/v1beta1/conditions/v1beta2"
	ctrl "sigs.k8s.io/controller-runtime"
	"sigs.k8s.io/controller-runtime/pkg/event"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	capvcontext "sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/net"
)

func sanitizeIPAddrs(ctx context.Context, ipAddrs []string) []string {
	log := ctrl.LoggerFrom(ctx)

	if len(ipAddrs) == 0 {
		return nil
	}
	newIPAddrs := []string{}
	for _, addr := range ipAddrs {
		if err := net.ErrOnLocalOnlyIPAddr(addr); err != nil {
			log.V(4).Info("Ignoring IP address", "reason", err.Error())
		} else {
			newIPAddrs = append(newIPAddrs, addr)
		}
	}
	return newIPAddrs
}

// findVM searches for a VM in one of two ways:
//  1. If the BIOS UUID is available, then it is used to find the VM.
//  2. Lacking the BIOS UUID, the VM is queried by its instance UUID,
//     which was assigned the value of the VSphereVM resource's UID string.
//  3. If it is not found by instance UUID, fallback to an inventory path search
//     using the vm folder path and the VSphereVM name
func findVM(ctx context.Context, vmCtx *capvcontext.VMContext) (types.ManagedObjectReference, error) {
	log := ctrl.LoggerFrom(ctx)

	if biosUUID := vmCtx.VSphereVM.Spec.BiosUUID; biosUUID != "" {
		objRef, err := vmCtx.Session.FindByBIOSUUID(ctx, biosUUID)
		if err != nil {
			return types.ManagedObjectReference{}, err
		}
		if objRef == nil {
			log.Info("VM not found by bios uuid", "biosUUID", biosUUID)
			return types.ManagedObjectReference{}, errNotFound{uuid: biosUUID}
		}
		log.Info("VM found by bios uuid", "vmRef", objRef.Reference())
		return objRef.Reference(), nil
	}

	instanceUUID := string(vmCtx.VSphereVM.UID)
	objRef, err := vmCtx.Session.FindByInstanceUUID(ctx, instanceUUID)
	if err != nil {
		return types.ManagedObjectReference{}, err
	}
	if objRef == nil {
		// fallback to use inventory paths
		folder, err := vmCtx.Session.Finder.FolderOrDefault(ctx, vmCtx.VSphereVM.Spec.Folder)
		if err != nil {
			return types.ManagedObjectReference{}, err
		}
		inventoryPath := path.Join(folder.InventoryPath, vmCtx.VSphereVM.Name)
		log.Info("Using inventory path to find VM", "inventoryPath", inventoryPath)
		vm, err := vmCtx.Session.Finder.VirtualMachine(ctx, inventoryPath)
		if err != nil {
			if isVirtualMachineNotFound(err) {
				log.Info("VM not found by instance uuid or inventory path")
				return types.ManagedObjectReference{}, errNotFound{byInventoryPath: inventoryPath}
			}
			return types.ManagedObjectReference{}, err
		}
		log.Info("VM found by inventory path", "vmRef", vm.Reference())
		return vm.Reference(), nil
	}
	log.Info("VM found by instance uuid", "vmRef", objRef.Reference())
	return objRef.Reference(), nil
}

func getTask(ctx context.Context, vmCtx *capvcontext.VMContext) *mo.Task {
	if vmCtx.VSphereVM.Status.TaskRef == "" {
		return nil
	}
	var obj mo.Task
	moRef := types.ManagedObjectReference{
		Type:  morefTypeTask,
		Value: vmCtx.VSphereVM.Status.TaskRef,
	}
	if err := vmCtx.Session.RetrieveOne(ctx, moRef, []string{"info"}, &obj); err != nil {
		return nil
	}
	return &obj
}

// reconcileInFlightTask determines if a task associated to the VSphereVM object
// is in flight or not.
func reconcileInFlightTask(ctx context.Context, vmCtx *capvcontext.VMContext) (bool, error) {
	// Check to see if there is an in-flight task.
	task := getTask(ctx, vmCtx)
	return checkAndRetryTask(ctx, vmCtx, task)
}

// checkAndRetryTask verifies whether the task exists and if the
// task should be reconciled which is determined by the task state retryAfter value set.
func checkAndRetryTask(ctx context.Context, vmCtx *capvcontext.VMContext, task *mo.Task) (bool, error) {
	log := ctrl.LoggerFrom(ctx)

	// If no task was found then make sure to clear the VSphereVM
	// resource's Status.TaskRef field.
	if task == nil {
		vmCtx.VSphereVM.Status.TaskRef = ""
		return false, nil
	}

	// Since RetryAfter is set, the last task failed. Wait for the RetryAfter time duration to expire
	// before checking/resetting the task.
	if !vmCtx.VSphereVM.Status.RetryAfter.IsZero() && time.Now().Before(vmCtx.VSphereVM.Status.RetryAfter.Time) {
		return false, errors.Errorf("last task failed retry after %v", vmCtx.VSphereVM.Status.RetryAfter)
	}

	// Otherwise the course of action is determined by the state of the task.
	log = log.WithValues("taskRef", task.Reference().Value, "taskState", task.Info.State, "taskDescriptionID", task.Info.DescriptionId)
	ctx = ctrl.LoggerInto(ctx, log) //nolint:ineffassign,staticcheck // ensure the logger is up-to-date in ctx, even if we currently don't use ctx below.
	switch task.Info.State {
	case types.TaskInfoStateQueued:
		log.Info("Task found: Task is still pending")
		return true, nil
	case types.TaskInfoStateRunning:
		log.Info("Task found: Task is still running")
		return true, nil
	case types.TaskInfoStateSuccess:
		log.Info("Task found: Task is a success")
		vmCtx.VSphereVM.Status.TaskRef = ""
		return false, nil
	case types.TaskInfoStateError:
		// NOTE: When a task fails there is no simple way to understand which operation is failing (e.g. cloning or powering on)
		// so we are reporting failures using a dedicated reason until we find a better solution.
		var errorMessage string

		if task.Info.Error != nil {
			// If the result is InvalidPowerState and the VM's current and expected state are the same, it means the VM already
			// was in the state we wanted it to be.
			if isInvalidPowerStateAndExpectedPowerState(task.Info.Error.Fault) {
				log.Info("Task found: Task failed to power vm, but VM is already in expected state")
				vmCtx.VSphereVM.Status.TaskRef = ""
				return false, nil
			}
			errorMessage = task.Info.Error.LocalizedMessage
		}

		log.Info("Task found: Task failed")
		v1beta1conditions.MarkFalse(vmCtx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.TaskFailure, clusterv1beta1.ConditionSeverityInfo, "%s", errorMessage)
		v1beta2conditions.Set(vmCtx.VSphereVM, metav1.Condition{
			Type:   infrav1.VSphereVMVirtualMachineProvisionedV1Beta2Condition,
			Status: metav1.ConditionFalse,
			Reason: infrav1.VSphereVMVirtualMachineTaskFailedV1Beta2Reason,
		})

		// Instead of directly requeuing the failed task, wait for the RetryAfter duration to pass
		// before resetting the taskRef from the VSphereVM status.
		if vmCtx.VSphereVM.Status.RetryAfter.IsZero() {
			vmCtx.VSphereVM.Status.RetryAfter = metav1.Time{Time: time.Now().Add(1 * time.Minute)}
		} else {
			vmCtx.VSphereVM.Status.TaskRef = ""
			vmCtx.VSphereVM.Status.RetryAfter = metav1.Time{}
		}
		return true, nil
	default:
		return false, errors.Errorf("unknown task state %q for %q", task.Info.State, vmCtx)
	}
}

func isInvalidPowerStateAndExpectedPowerState(f types.BaseMethodFault) bool {
	invalidPowerState, ok := f.(*types.InvalidPowerState)
	if !ok {
		return false
	}

	return invalidPowerState.ExistingState == invalidPowerState.RequestedState
}

func reconcileVSphereVMWhenNetworkIsReady(ctx context.Context, virtualMachineCtx *virtualMachineContext, powerOnTask *object.Task) {
	reconcileVSphereVMOnChannel(
		ctx,
		&virtualMachineCtx.VMContext,
		func() (<-chan []interface{}, <-chan error, error) {
			// Wait for the VM to be powered on.
			powerOnTaskInfo, err := powerOnTask.WaitForResult(ctx)
			if err != nil && powerOnTaskInfo == nil {
				return nil, nil, errors.Wrapf(err, "failed to wait for power on op for vm %s", virtualMachineCtx)
			}
			powerState, err := virtualMachineCtx.Obj.PowerState(ctx)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to get power state for vm %s", virtualMachineCtx)
			}
			if powerState != types.VirtualMachinePowerStatePoweredOn {
				return nil, nil, errors.Errorf(
					"unexpected power state %v for vm %s",
					powerState, ctx)
			}

			// Wait for all NICs to have valid MAC addresses.
			if err := waitForMacAddresses(ctx, virtualMachineCtx); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to wait for mac addresses for vm %s", virtualMachineCtx)
			}

			// Get all the MAC addresses. This is done separately from waiting
			// for all NICs to have MAC addresses in order to ensure the order
			// of the retrieved MAC addresses matches the order of the device
			// specs, and not the propery change order.
			_, macToDeviceIndex, deviceToMacIndex, err := getMacAddresses(ctx, virtualMachineCtx)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to get mac addresses for vm %s", virtualMachineCtx)
			}

			// Wait for the IP addresses to show up for the VM.
			chanIPAddresses, chanErrs := waitForIPAddresses(ctx, virtualMachineCtx, macToDeviceIndex, deviceToMacIndex)

			// Trigger a reconcile every time a new IP is discovered.
			chanOfLoggerKeysAndValues := make(chan []interface{})
			go func() {
				for ip := range chanIPAddresses {
					chanOfLoggerKeysAndValues <- []interface{}{
						"reason", "network",
						"ipAddress", ip,
					}
				}
			}()
			return chanOfLoggerKeysAndValues, chanErrs, nil
		})
}

func reconcileVSphereVMOnTaskCompletion(ctx context.Context, vmCtx *capvcontext.VMContext) {
	log := ctrl.LoggerFrom(ctx)

	task := getTask(ctx, vmCtx)
	if task == nil {
		log.V(4).Info("Skipping reconcile VSphereVM on task completion, because there is no task")
		return
	}
	taskRef := task.Reference()
	taskHelper := object.NewTask(vmCtx.Session.Client.Client, taskRef)

	log.Info("Enqueuing reconcile request on task completion",
		"taskRef", taskRef,
		"taskName", task.Info.Name,
		"taskEntityName", task.Info.EntityName,
		"taskDescriptionID", task.Info.DescriptionId)

	reconcileVSphereVMOnFuncCompletion(ctx, vmCtx, func() ([]interface{}, error) {
		taskInfo, err := taskHelper.WaitForResult(ctx)

		// An error is only returned if the process of waiting for the result
		// failed, *not* if the task itself failed.
		if err != nil && taskInfo == nil {
			return nil, err
		}
		// do not queue in the event channel when task fails as we don't
		// want to retry right away
		if taskInfo.State == types.TaskInfoStateError {
			return nil, errors.Errorf("task failed: task is in state error")
		}

		return []interface{}{
			"reason", "task",
			"taskRef", taskRef,
			"taskName", taskInfo.Name,
			"taskEntityName", taskInfo.EntityName,
			"taskState", taskInfo.State,
			"taskDescriptionID", taskInfo.DescriptionId,
		}, nil
	})
}

func reconcileVSphereVMOnFuncCompletion(ctx context.Context, vmCtx *capvcontext.VMContext, waitFn func() (loggerKeysAndValues []interface{}, _ error)) {
	log := ctrl.LoggerFrom(ctx)

	obj := vmCtx.VSphereVM.DeepCopy()
	gvk := infrav1.GroupVersion.WithKind("VSphereVM")

	// Wait on the function to complete in a background goroutine.
	go func() {
		loggerKeysAndValues, err := waitFn()
		if err != nil {
			log.Error(err, "failed to wait on func")
			return
		}

		// Once the task has completed (successfully or otherwise), trigger
		// a reconcile event for the associated resource by sending a
		// GenericEvent into the event channel for the resource type.
		log.Info("Triggering GenericEvent", loggerKeysAndValues...)
		eventChannel := vmCtx.GetGenericEventChannelFor(gvk)
		eventChannel <- event.GenericEvent{
			Object: obj,
		}
	}()
}

func reconcileVSphereVMOnChannel(ctx context.Context, vmCtx *capvcontext.VMContext, waitFn func() (<-chan []interface{}, <-chan error, error)) {
	log := ctrl.LoggerFrom(ctx)

	obj := vmCtx.VSphereVM.DeepCopy()
	gvk := obj.GetObjectKind().GroupVersionKind()

	// Send a generic event for every set of logger keys/values received
	// on the channel.
	go func() {
		chanOfLoggerKeysAndValues, chanErrs, err := waitFn()
		if err != nil {
			log.Error(err, "failed to wait on func")
			return
		}
		for {
			select {
			case loggerKeysAndValues := <-chanOfLoggerKeysAndValues:
				if loggerKeysAndValues == nil {
					return
				}
				go func() {
					// Trigger a reconcile event for the associated resource by
					// sending a GenericEvent into the event channel for the resource
					// type.
					log.Info("Triggering GenericEvent", loggerKeysAndValues...)
					eventChannel := vmCtx.GetGenericEventChannelFor(gvk)
					eventChannel <- event.GenericEvent{
						Object: obj,
					}
				}()
			case err := <-chanErrs:
				if err != nil {
					log.Error(err, "error occurred while waiting to trigger a generic event")
				}
				return
			}
		}
	}()
}

// waitForMacAddresses waits for all configured network devices to have
// valid MAC addresses.
func waitForMacAddresses(ctx context.Context, virtualMachineCtx *virtualMachineContext) error {
	return property.Wait(
		ctx, property.DefaultCollector(virtualMachineCtx.Session.Client.Client),
		virtualMachineCtx.Obj.Reference(), []string{"config.hardware.device"},
		func(propertyChanges []types.PropertyChange) bool {
			for _, propChange := range propertyChanges {
				if propChange.Op != types.PropertyChangeOpAssign {
					continue
				}
				deviceList := object.VirtualDeviceList(propChange.Val.(types.ArrayOfVirtualDevice).VirtualDevice)
				for _, dev := range deviceList {
					if nic, ok := dev.(types.BaseVirtualEthernetCard); ok {
						mac := nic.GetVirtualEthernetCard().MacAddress
						if mac == "" {
							return false
						}
					}
				}
			}
			return true
		})
}

// getMacAddresses gets the MAC addresses for all network devices.
// This happens separately from waitForMacAddresses to ensure returned order of
// devices matches the spec and not order in which the property changes were
// noticed.
func getMacAddresses(ctx context.Context, virtualMachineCtx *virtualMachineContext) ([]string, map[string]int, map[int]string, error) {
	var (
		vm                   mo.VirtualMachine
		macAddresses         = make([]string, 0)
		macToDeviceSpecIndex = map[string]int{}
		deviceSpecIndexToMac = map[int]string{}
	)
	if err := virtualMachineCtx.Obj.Properties(ctx, virtualMachineCtx.Obj.Reference(), []string{"config.hardware.device"}, &vm); err != nil {
		return nil, nil, nil, err
	}
	i := 0
	for _, device := range vm.Config.Hardware.Device {
		nic, ok := device.(types.BaseVirtualEthernetCard)
		if !ok {
			continue
		}
		mac := nic.GetVirtualEthernetCard().MacAddress
		macAddresses = append(macAddresses, mac)
		macToDeviceSpecIndex[mac] = i
		deviceSpecIndexToMac[i] = mac
		i++
	}
	return macAddresses, macToDeviceSpecIndex, deviceSpecIndexToMac, nil
}

// waitForIPAddresses waits for all network devices that should be getting an
// IP address to have an IP address. This is any network device that specifies a
// network name and DHCP for v4 or v6 or one or more static IP addresses.
func waitForIPAddresses(
	ctx context.Context,
	virtualMachineCtx *virtualMachineContext,
	macToDeviceIndex map[string]int,
	deviceToMacIndex map[int]string) (<-chan string, <-chan error) {
	log := ctrl.LoggerFrom(ctx)

	var (
		chanErrs          = make(chan error)
		chanIPAddresses   = make(chan string)
		macToHasIPv4Lease = map[string]struct{}{}
		macToHasIPv6Lease = map[string]struct{}{}
		macToSkipped      = map[string]map[string]struct{}{}
		macToHasStaticIP  = map[string]map[string]struct{}{}
		propCollector     = property.DefaultCollector(virtualMachineCtx.Session.Client.Client)
	)

	// Initialize the nested maps early.
	for mac := range macToDeviceIndex {
		macToSkipped[mac] = map[string]struct{}{}
		macToHasStaticIP[mac] = map[string]struct{}{}
	}

	onPropertyChange := func(propertyChanges []types.PropertyChange) bool {
		for _, propChange := range propertyChanges {
			if propChange.Op != types.PropertyChangeOpAssign {
				continue
			}
			nics := propChange.Val.(types.ArrayOfGuestNicInfo).GuestNicInfo
			for _, nic := range nics {
				mac := nic.MacAddress
				if mac == "" || nic.IpConfig == nil {
					continue
				}
				// Ignore any that don't correspond to a network
				// device spec.
				deviceSpecIndex, ok := macToDeviceIndex[mac]
				if !ok {
					chanErrs <- errors.Errorf("unknown device spec index for mac %s while waiting for ip addresses for vm %s", mac, virtualMachineCtx)
					// Return true to stop the property collector from waiting
					// on any more changes.
					return true
				}
				if deviceSpecIndex < 0 || deviceSpecIndex >= len(virtualMachineCtx.VSphereVM.Spec.Network.Devices) {
					chanErrs <- errors.Errorf("invalid device spec index %d for mac %s while waiting for ip addresses for vm %s", deviceSpecIndex, mac, virtualMachineCtx)
					// Return true to stop the property collector from waiting
					// on any more changes.
					return true
				}

				// Get the network device spec that corresponds to the MAC.
				deviceSpec := virtualMachineCtx.VSphereVM.Spec.Network.Devices[deviceSpecIndex]

				// Look at each IP and determine whether a reconcile has
				// been triggered for the IP.
				for _, discoveredIPInfo := range nic.IpConfig.IpAddress {
					discoveredIP := discoveredIPInfo.IpAddress

					// Ignore link-local addresses.
					if err := net.ErrOnLocalOnlyIPAddr(discoveredIP); err != nil {
						if _, ok := macToSkipped[mac][discoveredIP]; !ok {
							log.Info("Ignoring IP address", "reason", err.Error())
							macToSkipped[mac][discoveredIP] = struct{}{}
						}
						continue
					}

					// Check to see if the IP is in the list of the device
					// spec's static IP addresses.
					isStatic := false
					for _, specIP := range deviceSpec.IPAddrs {
						// The static IP assigned to the VM is required in the CIDR format
						ip, _, _ := gonet.ParseCIDR(specIP)
						if discoveredIP == ip.String() {
							isStatic = true
							break
						}
					}

					// If it's a static IP then check to see if the IP has
					// triggered a reconcile yet.
					switch {
					case isStatic:
						if _, ok := macToHasStaticIP[mac][discoveredIP]; !ok {
							// No reconcile yet. Record the IP send it to the
							// channel.
							log.Info("Discovered IP address",
								"addressType", "static",
								"addressValue", discoveredIP)
							macToHasStaticIP[mac][discoveredIP] = struct{}{}
							chanIPAddresses <- discoveredIP
						}
					case gonet.ParseIP(discoveredIP).To4() != nil:
						// An IPv4 address...
						if deviceSpec.DHCP4 {
							// Has an IPv4 lease been discovered yet?
							if _, ok := macToHasIPv4Lease[mac]; !ok {
								log.Info("Discovered IP address",
									"addressType", "dhcp4",
									"addressValue", discoveredIP)
								macToHasIPv4Lease[mac] = struct{}{}
								chanIPAddresses <- discoveredIP
							}
						}
					default:
						// An IPv6 address..
						if deviceSpec.DHCP6 {
							// Has an IPv6 lease been discovered yet?
							if _, ok := macToHasIPv6Lease[mac]; !ok {
								log.Info("Discovered IP address",
									"addressType", "dhcp6",
									"addressValue", discoveredIP)
								macToHasIPv6Lease[mac] = struct{}{}
								chanIPAddresses <- discoveredIP
							}
						}
					}
				}
			}
		}

		// Determine whether or not the wait operation is over by whether
		// or not the VM has all of the requested IP addresses.
		for i, deviceSpec := range virtualMachineCtx.VSphereVM.Spec.Network.Devices {
			// If the device spec has SkipIPAllocation set true then
			// the wait is not required
			if deviceSpec.SkipIPAllocation {
				continue
			}

			mac, ok := deviceToMacIndex[i]
			if !ok {
				chanErrs <- errors.Errorf("invalid mac index %d waiting for ip addresses for vm %s", i, virtualMachineCtx)

				// Return true to stop the property collector from waiting
				// on any more changes.
				return true
			}
			// If the device spec requires DHCP4 then the Wait is not
			// over if there is no IPv4 lease.
			if deviceSpec.DHCP4 {
				if _, ok := macToHasIPv4Lease[mac]; !ok {
					log.Info("The VM is missing the requested IP address",
						"addressType", "dhcp4")
					return false
				}
			}
			// If the device spec requires DHCP6 then the Wait is not
			// over if there is no IPv6 lease.
			if deviceSpec.DHCP6 {
				if _, ok := macToHasIPv6Lease[mac]; !ok {
					log.Info("The VM is missing the requested IP address",
						"addressType", "dhcp6")
					return false
				}
			}
			// If the device spec requires static IP addresses, the wait
			// is not over if the device lacks one of those addresses.
			for _, specIP := range deviceSpec.IPAddrs {
				ip, _, _ := gonet.ParseCIDR(specIP)
				if _, ok := macToHasStaticIP[mac][ip.String()]; !ok {
					log.Info("The VM is missing the requested IP address",
						"addressType", "static",
						"addressValue", specIP)
					return false
				}
			}
		}

		log.Info("The VM has all of the requested IP addresses")
		return true
	}

	// The wait function will not return true until all the VM's
	// network devices have IP assignments that match the requested
	// network device specs. However, every time a new IP is discovered,
	// a reconcile request will be triggered for the VSphereVM.
	go func() {
		// Note: We intentionally don't use the context from the Reconcile
		// so this go routine continues independent of the current Reconcile.
		ctx := context.Background()
		if err := property.Wait(
			ctx, propCollector, virtualMachineCtx.Obj.Reference(),
			[]string{"guest.net"}, onPropertyChange); err != nil {
			chanErrs <- errors.Wrapf(err, "failed to wait for ip addresses for vm %s", virtualMachineCtx)
		}
		close(chanIPAddresses)
		close(chanErrs)
	}()

	return chanIPAddresses, chanErrs
}
