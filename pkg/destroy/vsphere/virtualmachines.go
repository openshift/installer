package vsphere

import (
	"context"
	"time"

	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/property"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	utilerrors "k8s.io/apimachinery/pkg/util/errors"
)

func (o *ClusterUninstaller) getVirtualMachineManagedObjects(ctx context.Context, moRef []types.ManagedObjectReference) ([]mo.VirtualMachine, error) {
	var virtualMachineMoList []mo.VirtualMachine
	if len(moRef) > 0 {
		pc := property.DefaultCollector(o.Client)
		err := pc.Retrieve(ctx, moRef, nil, &virtualMachineMoList)
		if err != nil {
			return nil, err
		}
	}
	return virtualMachineMoList, nil
}

func (o *ClusterUninstaller) listVirtualMachines(ctx context.Context) ([]mo.VirtualMachine, error) {
	virtualMachineList, err := o.getAttachedObjectsOnTag("VirtualMachine")
	if err != nil {
		return nil, err
	}

	return o.getVirtualMachineManagedObjects(ctx, virtualMachineList)
}

func isPoweredOff(vmMO mo.VirtualMachine) bool {
	return vmMO.Summary.Runtime.PowerState == "poweredOff"
}

func (o *ClusterUninstaller) stopVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name)
	if !isPoweredOff(vmMO) {
		vm := object.NewVirtualMachine(o.Client, vmMO.Reference())
		task, err := vm.PowerOff(ctx)
		if err == nil {
			err = task.Wait(ctx)
		}
		if err != nil {
			virtualMachineLogger.Debug(err)
			return err
		}
	}
	virtualMachineLogger.Debug("Powered off")
	return nil
}

func (o *ClusterUninstaller) stopVirtualMachines() error {
	ctx, cancel := context.WithTimeout(o.context, time.Minute*30)
	defer cancel()

	o.Logger.Debug("Power Off Virtual Machines")
	found, err := o.listVirtualMachines(ctx)
	if err != nil {
		o.Logger.Debug(err)
		return err
	}

	var errs []error
	for _, vmMO := range found {
		if !isPoweredOff(vmMO) {
			if err := o.stopVirtualMachine(ctx, vmMO); err != nil {
				errs = append(errs, err)
			}
		}
	}

	return utilerrors.NewAggregate(errs)
}

func (o *ClusterUninstaller) deleteVirtualMachine(ctx context.Context, vmMO mo.VirtualMachine) error {
	virtualMachineLogger := o.Logger.WithField("VirtualMachine", vmMO.Name)
	vm := object.NewVirtualMachine(o.Client, vmMO.Reference())
	task, err := vm.Destroy(ctx)
	if err == nil {
		err = task.Wait(ctx)
	}
	if err != nil {
		virtualMachineLogger.Debug(err)
		return err
	}
	virtualMachineLogger.Info("Destroyed")
	return nil
}

func (o *ClusterUninstaller) deleteVirtualMachines() error {
	ctx, cancel := context.WithTimeout(o.context, time.Minute*30)
	defer cancel()

	o.Logger.Debug("Delete Virtual Machines")
	found, err := o.listVirtualMachines(ctx)
	if err != nil {
		o.Logger.Debug(err)
		return err
	}

	var errs []error
	for _, vmMO := range found {
		if err := o.deleteVirtualMachine(ctx, vmMO); err != nil {
			errs = append(errs, err)
		}
	}

	return utilerrors.NewAggregate(errs)
}
