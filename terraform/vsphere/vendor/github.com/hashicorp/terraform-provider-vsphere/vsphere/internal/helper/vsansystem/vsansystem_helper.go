package vsansystem

import (
	"context"
	"time"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// Properties Returns the HostVsanSystem ManagedObject for the HostVsanSystem object.
func Properties(client *govmomi.Client, hss *object.HostVsanSystem, apiTimeout time.Duration) (*mo.HostVsanSystem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()
	var hvsProps mo.HostVsanSystem
	err := hss.Properties(ctx, hss.Reference(), nil, &hvsProps)
	if err != nil {
		return nil, err
	}
	return &hvsProps, err
}

// FromHost returns a host's HostVsanSystem object.
func FromHost(client *govmomi.Client, host *object.HostSystem, apiTimeout time.Duration) (*object.HostVsanSystem, error) {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()
	return host.ConfigManager().VsanSystem(ctx)
}

// RemoveDiskMapping removes the disks specified in diskMap from the disk group
// on host.
func RemoveDiskMapping(client *govmomi.Client, host *object.HostSystem, hvs *object.HostVsanSystem, diskMap *types.VsanHostDiskMapping, apiTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	// If the SSD name is set, then the whole disk group needs to be removed using
	// RemoveDiskMapping, otherwise use RemoveDisk to just remove storage disks
	// from the group.
	if diskMap.Ssd.CanonicalName != "" {
		ntask := types.RemoveDiskMapping_Task{
			This:    hvs.Reference(),
			Mapping: []types.VsanHostDiskMapping{*diskMap},
		}

		resp, err := methods.RemoveDiskMapping_Task(ctx, host.Client().RoundTripper, &ntask)
		if err != nil {
			return err
		}
		task := object.NewTask(client.Client, resp.Returnval)
		if err := task.Wait(ctx); err != nil {
			return err
		}
	} else {
		ntask := types.RemoveDisk_Task{
			This: hvs.Reference(),
			Disk: diskMap.NonSsd,
		}

		resp, err := methods.RemoveDisk_Task(ctx, host.Client().RoundTripper, &ntask)
		if err != nil {
			return err
		}
		task := object.NewTask(client.Client, resp.Returnval)
		if err := task.Wait(ctx); err != nil {
			return err
		}

	}
	return nil
}

// InitializeDisks initializes and adds disks to the specified host disk group.
func InitializeDisks(client *govmomi.Client, host *object.HostSystem, hvs *object.HostVsanSystem, diskMap *types.VsanHostDiskMapping, apiTimeout time.Duration) error {
	ctx, cancel := context.WithTimeout(context.Background(), apiTimeout)
	defer cancel()

	ntask := types.InitializeDisks_Task{
		This:    hvs.Reference(),
		Mapping: []types.VsanHostDiskMapping{*diskMap},
	}

	resp, err := methods.InitializeDisks_Task(ctx, host.Client().RoundTripper, &ntask)
	if err != nil {
		return err
	}
	task := object.NewTask(client.Client, resp.Returnval)
	if err := task.Wait(ctx); err != nil {
		return err
	}
	return nil
}
