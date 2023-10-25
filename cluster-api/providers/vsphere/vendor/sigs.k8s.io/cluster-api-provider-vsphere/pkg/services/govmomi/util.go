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
	gonet "net"
	"path"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
	metav1 "k8s.io/apimachinery/pkg/apis/meta/v1"
	clusterv1 "sigs.k8s.io/cluster-api/api/v1beta1"
	"sigs.k8s.io/cluster-api/util/conditions"
	"sigs.k8s.io/controller-runtime/pkg/event"

	infrav1 "sigs.k8s.io/cluster-api-provider-vsphere/apis/v1beta1"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/context"
	"sigs.k8s.io/cluster-api-provider-vsphere/pkg/services/govmomi/net"
)

func sanitizeIPAddrs(ctx *context.VMContext, ipAddrs []string) []string {
	if len(ipAddrs) == 0 {
		return nil
	}
	newIPAddrs := []string{}
	for _, addr := range ipAddrs {
		if err := net.ErrOnLocalOnlyIPAddr(addr); err != nil {
			ctx.Logger.V(4).Info("ignoring IP address", "reason", err.Error())
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
func findVM(ctx *context.VMContext) (types.ManagedObjectReference, error) {
	if biosUUID := ctx.VSphereVM.Spec.BiosUUID; biosUUID != "" {
		objRef, err := ctx.Session.FindByBIOSUUID(ctx, biosUUID)
		if err != nil {
			return types.ManagedObjectReference{}, err
		}
		if objRef == nil {
			ctx.Logger.Info("vm not found by bios uuid", "biosuuid", biosUUID)
			return types.ManagedObjectReference{}, errNotFound{uuid: biosUUID}
		}
		ctx.Logger.Info("vm found by bios uuid", "vmref", objRef.Reference())
		return objRef.Reference(), nil
	}

	instanceUUID := string(ctx.VSphereVM.UID)
	objRef, err := ctx.Session.FindByInstanceUUID(ctx, instanceUUID)
	if err != nil {
		return types.ManagedObjectReference{}, err
	}
	if objRef == nil {
		// fallback to use inventory paths
		folder, err := ctx.Session.Finder.FolderOrDefault(ctx, ctx.VSphereVM.Spec.Folder)
		if err != nil {
			return types.ManagedObjectReference{}, err
		}
		inventoryPath := path.Join(folder.InventoryPath, ctx.VSphereVM.Name)
		ctx.Logger.Info("using inventory path to find vm", "path", inventoryPath)
		vm, err := ctx.Session.Finder.VirtualMachine(ctx, inventoryPath)
		if err != nil {
			if isVirtualMachineNotFound(err) {
				return types.ManagedObjectReference{}, errNotFound{byInventoryPath: inventoryPath}
			}
			return types.ManagedObjectReference{}, err
		}
		ctx.Logger.Info("vm found by name", "vmref", vm.Reference())
		return vm.Reference(), nil
	}
	ctx.Logger.Info("vm found by instance uuid", "vmref", objRef.Reference())
	return objRef.Reference(), nil
}

func getTask(ctx *context.VMContext) *mo.Task {
	if ctx.VSphereVM.Status.TaskRef == "" {
		return nil
	}
	var obj mo.Task
	moRef := types.ManagedObjectReference{
		Type:  morefTypeTask,
		Value: ctx.VSphereVM.Status.TaskRef,
	}
	if err := ctx.Session.RetrieveOne(ctx, moRef, []string{"info"}, &obj); err != nil {
		return nil
	}
	return &obj
}

// reconcileInFlightTask determines if a task associated to the VSphereVM object
// is in flight or not.
func reconcileInFlightTask(ctx *context.VMContext) (bool, error) {
	// Check to see if there is an in-flight task.
	task := getTask(ctx)
	return checkAndRetryTask(ctx, task)
}

// checkAndRetryTask verifies whether the task exists and if the
// task should be reconciled which is determined by the task state retryAfter value set.
func checkAndRetryTask(ctx *context.VMContext, task *mo.Task) (bool, error) {
	// If no task was found then make sure to clear the VSphereVM
	// resource's Status.TaskRef field.
	if task == nil {
		ctx.VSphereVM.Status.TaskRef = ""
		return false, nil
	}

	// Since RetryAfter is set, the last task failed. Wait for the RetryAfter time duration to expire
	// before checking/resetting the task.
	if !ctx.VSphereVM.Status.RetryAfter.IsZero() && time.Now().Before(ctx.VSphereVM.Status.RetryAfter.Time) {
		return false, errors.Errorf("last task failed retry after %v", ctx.VSphereVM.Status.RetryAfter)
	}

	// Otherwise the course of action is determined by the state of the task.
	logger := ctx.Logger.WithName(task.Reference().Value)
	logger.Info("task found", "state", task.Info.State, "description-id", task.Info.DescriptionId)
	switch task.Info.State {
	case types.TaskInfoStateQueued:
		logger.Info("task is still pending", "description-id", task.Info.DescriptionId)
		return true, nil
	case types.TaskInfoStateRunning:
		logger.Info("task is still running", "description-id", task.Info.DescriptionId)
		return true, nil
	case types.TaskInfoStateSuccess:
		logger.Info("task is a success", "description-id", task.Info.DescriptionId)
		ctx.VSphereVM.Status.TaskRef = ""
		return false, nil
	case types.TaskInfoStateError:
		logger.Info("task failed", "description-id", task.Info.DescriptionId)

		// NOTE: When a task fails there is not simple way to understand which operation is failing (e.g. cloning or powering on)
		// so we are reporting failures using a dedicated reason until we find a better solution.
		var errorMessage string

		if task.Info.Error != nil {
			errorMessage = task.Info.Error.LocalizedMessage
		}
		conditions.MarkFalse(ctx.VSphereVM, infrav1.VMProvisionedCondition, infrav1.TaskFailure, clusterv1.ConditionSeverityInfo, errorMessage)

		// Instead of directly requeuing the failed task, wait for the RetryAfter duration to pass
		// before resetting the taskRef from the VSphereVM status.
		if ctx.VSphereVM.Status.RetryAfter.IsZero() {
			ctx.VSphereVM.Status.RetryAfter = metav1.Time{Time: time.Now().Add(1 * time.Minute)}
		} else {
			ctx.VSphereVM.Status.TaskRef = ""
			ctx.VSphereVM.Status.RetryAfter = metav1.Time{}
		}
		return true, nil
	default:
		return false, errors.Errorf("unknown task state %q for %q", task.Info.State, ctx)
	}
}

func reconcileVSphereVMWhenNetworkIsReady(ctx *virtualMachineContext, powerOnTask *object.Task) {
	reconcileVSphereVMOnChannel(
		&ctx.VMContext,
		func() (<-chan []interface{}, <-chan error, error) {
			// Wait for the VM to be powered on.
			powerOnTaskInfo, err := powerOnTask.WaitForResult(ctx)
			if err != nil && powerOnTaskInfo == nil {
				return nil, nil, errors.Wrapf(err, "failed to wait for power on op for vm %s", ctx)
			}
			powerState, err := ctx.Obj.PowerState(ctx)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to get power state for vm %s", ctx)
			}
			if powerState != types.VirtualMachinePowerStatePoweredOn {
				return nil, nil, errors.Errorf(
					"unexpected power state %v for vm %s",
					powerState, ctx)
			}

			// Wait for all NICs to have valid MAC addresses.
			if err := waitForMacAddresses(ctx); err != nil {
				return nil, nil, errors.Wrapf(err, "failed to wait for mac addresses for vm %s", ctx)
			}

			// Get all the MAC addresses. This is done separately from waiting
			// for all NICs to have MAC addresses in order to ensure the order
			// of the retrieved MAC addresses matches the order of the device
			// specs, and not the propery change order.
			_, macToDeviceIndex, deviceToMacIndex, err := getMacAddresses(ctx)
			if err != nil {
				return nil, nil, errors.Wrapf(err, "failed to get mac addresses for vm %s", ctx)
			}

			// Wait for the IP addresses to show up for the VM.
			chanIPAddresses, chanErrs := waitForIPAddresses(ctx, macToDeviceIndex, deviceToMacIndex)

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

func reconcileVSphereVMOnTaskCompletion(ctx *context.VMContext) {
	task := getTask(ctx)
	if task == nil {
		ctx.Logger.V(4).Info(
			"skipping reconcile VSphereVM on task completion",
			"reason", "no-task")
		return
	}
	taskRef := task.Reference()
	taskHelper := object.NewTask(ctx.Session.Client.Client, taskRef)

	ctx.Logger.Info(
		"enqueuing reconcile request on task completion",
		"task-ref", taskRef,
		"task-name", task.Info.Name,
		"task-entity-name", task.Info.EntityName,
		"task-description-id", task.Info.DescriptionId)

	reconcileVSphereVMOnFuncCompletion(ctx, func() ([]interface{}, error) {
		taskInfo, err := taskHelper.WaitForResult(ctx)

		// An error is only returned if the process of waiting for the result
		// failed, *not* if the task itself failed.
		if err != nil && taskInfo == nil {
			return nil, err
		}
		// do not queue in the event channel when task fails as we don't
		// want to retry right away
		if taskInfo.State == types.TaskInfoStateError {
			ctx.Logger.Info("async task wait failed")
			return nil, errors.Errorf("task failed")
		}

		return []interface{}{
			"reason", "task",
			"task-ref", taskRef,
			"task-name", taskInfo.Name,
			"task-entity-name", taskInfo.EntityName,
			"task-state", taskInfo.State,
			"task-description-id", taskInfo.DescriptionId,
		}, nil
	})
}

func reconcileVSphereVMOnFuncCompletion(ctx *context.VMContext, waitFn func() (loggerKeysAndValues []interface{}, _ error)) {
	obj := ctx.VSphereVM.DeepCopy()
	gvk := obj.GetObjectKind().GroupVersionKind()

	// Wait on the function to complete in a background goroutine.
	go func() {
		loggerKeysAndValues, err := waitFn()
		if err != nil {
			ctx.Logger.Error(err, "failed to wait on func")
			return
		}

		// Once the task has completed (successfully or otherwise), trigger
		// a reconcile event for the associated resource by sending a
		// GenericEvent into the event channel for the resource type.
		ctx.Logger.Info("triggering GenericEvent", loggerKeysAndValues...)
		eventChannel := ctx.GetGenericEventChannelFor(gvk)
		eventChannel <- event.GenericEvent{
			Object: obj,
		}
	}()
}

func reconcileVSphereVMOnChannel(ctx *context.VMContext, waitFn func() (<-chan []interface{}, <-chan error, error)) {
	obj := ctx.VSphereVM.DeepCopy()
	gvk := obj.GetObjectKind().GroupVersionKind()

	// Send a generic event for every set of logger keys/values received
	// on the channel.
	go func() {
		chanOfLoggerKeysAndValues, chanErrs, err := waitFn()
		if err != nil {
			ctx.Logger.Error(err, "failed to wait on func")
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
					ctx.Logger.Info("triggering GenericEvent", loggerKeysAndValues...)
					eventChannel := ctx.GetGenericEventChannelFor(gvk)
					eventChannel <- event.GenericEvent{
						Object: obj,
					}
				}()
			case err := <-chanErrs:
				if err != nil {
					ctx.Logger.Error(err, "error occurred while waiting to trigger a generic event")
				}
				return
			case <-ctx.Done():
				return
			}
		}
	}()
}

// waitForMacAddresses waits for all configured network devices to have
// valid MAC addresses.
func waitForMacAddresses(ctx *virtualMachineContext) error {
	return property.Wait(
		ctx, property.DefaultCollector(ctx.Session.Client.Client),
		ctx.Obj.Reference(), []string{"config.hardware.device"},
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
func getMacAddresses(ctx *virtualMachineContext) ([]string, map[string]int, map[int]string, error) {
	var (
		vm                   mo.VirtualMachine
		macAddresses         []string
		macToDeviceSpecIndex = map[string]int{}
		deviceSpecIndexToMac = map[int]string{}
	)
	if err := ctx.Obj.Properties(ctx, ctx.Obj.Reference(), []string{"config.hardware.device"}, &vm); err != nil {
		return nil, nil, nil, err
	}
	i := 0
	for _, device := range vm.Config.Hardware.Device {
		if nic, ok := device.(types.BaseVirtualEthernetCard); ok {
			mac := nic.GetVirtualEthernetCard().MacAddress
			macAddresses = append(macAddresses, mac)
			macToDeviceSpecIndex[mac] = i
			deviceSpecIndexToMac[i] = mac
			i++
		}
	}
	return macAddresses, macToDeviceSpecIndex, deviceSpecIndexToMac, nil
}

// waitForIPAddresses waits for all network devices that should be getting an
// IP address to have an IP address. This is any network device that specifies a
// network name and DHCP for v4 or v6 or one or more static IP addresses.
// The gocyclo detector is disabled for this function as it is difficult to
// rewrite much simpler due to the maps used to track state and the lambdas
// that use the maps.
//
//nolint:gocyclo,gocognit
func waitForIPAddresses(
	ctx *virtualMachineContext,
	macToDeviceIndex map[string]int,
	deviceToMacIndex map[int]string) (<-chan string, <-chan error) {
	var (
		chanErrs          = make(chan error)
		chanIPAddresses   = make(chan string)
		macToHasIPv4Lease = map[string]struct{}{}
		macToHasIPv6Lease = map[string]struct{}{}
		macToSkipped      = map[string]map[string]struct{}{}
		macToHasStaticIP  = map[string]map[string]struct{}{}
		propCollector     = property.DefaultCollector(ctx.Session.Client.Client)
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
					chanErrs <- errors.Errorf("unknown device spec index for mac %s while waiting for ip addresses for vm %s", mac, ctx)
					// Return true to stop the property collector from waiting
					// on any more changes.
					return true
				}
				if deviceSpecIndex < 0 || deviceSpecIndex >= len(ctx.VSphereVM.Spec.Network.Devices) {
					chanErrs <- errors.Errorf("invalid device spec index %d for mac %s while waiting for ip addresses for vm %s", deviceSpecIndex, mac, ctx)
					// Return true to stop the property collector from waiting
					// on any more changes.
					return true
				}

				// Get the network device spec that corresponds to the MAC.
				deviceSpec := ctx.VSphereVM.Spec.Network.Devices[deviceSpecIndex]

				// Look at each IP and determine whether a reconcile has
				// been triggered for the IP.
				for _, discoveredIPInfo := range nic.IpConfig.IpAddress {
					discoveredIP := discoveredIPInfo.IpAddress

					// Ignore link-local addresses.
					if err := net.ErrOnLocalOnlyIPAddr(discoveredIP); err != nil {
						if _, ok := macToSkipped[mac][discoveredIP]; !ok {
							ctx.Logger.Info("ignoring IP address", "reason", err.Error())
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
							ctx.Logger.Info(
								"discovered IP address",
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
								ctx.Logger.Info(
									"discovered IP address",
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
								ctx.Logger.Info(
									"discovered IP address",
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
		for i, deviceSpec := range ctx.VSphereVM.Spec.Network.Devices {
			mac, ok := deviceToMacIndex[i]
			if !ok {
				chanErrs <- errors.Errorf("invalid mac index %d waiting for ip addresses for vm %s", i, ctx)

				// Return true to stop the property collector from waiting
				// on any more changes.
				return true
			}
			// If the device spec requires DHCP4 then the Wait is not
			// over if there is no IPv4 lease.
			if deviceSpec.DHCP4 {
				if _, ok := macToHasIPv4Lease[mac]; !ok {
					ctx.Logger.Info(
						"the VM is missing the requested IP address",
						"addressType", "dhcp4")
					return false
				}
			}
			// If the device spec requires DHCP6 then the Wait is not
			// over if there is no IPv4 lease.
			if deviceSpec.DHCP6 {
				if _, ok := macToHasIPv6Lease[mac]; !ok {
					ctx.Logger.Info(
						"the VM is missing the requested IP address",
						"addressType", "dhcp6")
					return false
				}
			}
			// If the device spec requires static IP addresses, the wait
			// is not over if the device lacks one of those addresses.
			for _, specIP := range deviceSpec.IPAddrs {
				ip, _, _ := gonet.ParseCIDR(specIP)
				if _, ok := macToHasStaticIP[mac][ip.String()]; !ok {
					ctx.Logger.Info(
						"the VM is missing the requested IP address",
						"addressType", "static",
						"addressValue", specIP)
					return false
				}
			}
		}

		ctx.Logger.Info("the VM has all of the requested IP addresses")
		return true
	}

	// The wait function will not return true until all the VM's
	// network devices have IP assignments that match the requested
	// network device specs. However, every time a new IP is discovered,
	// a reconcile request will be triggered for the VSphereVM.
	go func() {
		if err := property.Wait(
			ctx, propCollector, ctx.Obj.Reference(),
			[]string{"guest.net"}, onPropertyChange); err != nil {
			chanErrs <- errors.Wrapf(err, "failed to wait for ip addresses for vm %s", ctx)
		}
		close(chanIPAddresses)
		close(chanErrs)
	}()

	return chanIPAddresses, chanErrs
}
