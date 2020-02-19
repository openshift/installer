package vsphere

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/network"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

var dvsVersions = []string{
	"5.0.0",
	"5.1.0",
	"5.5.0",
	"6.0.0",
	"6.5.0",
}

// dvsFromUUID gets a DVS object from its UUID.
func dvsFromUUID(client *govmomi.Client, uuid string) (*object.VmwareDistributedVirtualSwitch, error) {
	dvsm := types.ManagedObjectReference{Type: "DistributedVirtualSwitchManager", Value: "DVSManager"}
	req := &types.QueryDvsByUuid{
		This: dvsm,
		Uuid: uuid,
	}
	resp, err := methods.QueryDvsByUuid(context.TODO(), client, req)
	if err != nil {
		return nil, err
	}

	return dvsFromMOID(client, resp.Returnval.Reference().Value)
}

// dvsFromMOID locates a DVS by its managed object reference ID.
func dvsFromMOID(client *govmomi.Client, id string) (*object.VmwareDistributedVirtualSwitch, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "VmwareDistributedVirtualSwitch",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	// Should be safe to return here. If our reference returned here and is not a
	// VmwareDistributedVirtualSwitch, then we have bigger problems and to be
	// honest we should be panicking anyway.
	return ds.(*object.VmwareDistributedVirtualSwitch), nil
}

// dvsFromPath gets a DVS object from its path.
func dvsFromPath(client *govmomi.Client, name string, dc *object.Datacenter) (*object.VmwareDistributedVirtualSwitch, error) {
	net, err := network.FromPath(client, name, dc)
	if err != nil {
		return nil, err
	}
	if net.Reference().Type != "VmwareDistributedVirtualSwitch" {
		return nil, fmt.Errorf("network at path %q is not a VMware distributed virtual switch (type %s)", name, net.Reference().Type)
	}
	return dvsFromMOID(client, net.Reference().Value)
}

// dvsProperties is a convenience method that wraps fetching the DVS MO from
// its higher-level object.
func dvsProperties(dvs *object.VmwareDistributedVirtualSwitch) (*mo.VmwareDistributedVirtualSwitch, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	var props mo.VmwareDistributedVirtualSwitch
	if err := dvs.Properties(ctx, dvs.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// upgradeDVS upgrades a DVS to a specific version. Downgrades are not
// supported and will result in an error. This should be checked before running
// this function.
func upgradeDVS(client *govmomi.Client, dvs *object.VmwareDistributedVirtualSwitch, version string) error {
	req := &types.PerformDvsProductSpecOperation_Task{
		This:      dvs.Reference(),
		Operation: "upgrade",
		ProductSpec: &types.DistributedVirtualSwitchProductSpec{
			Version: version,
		},
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	resp, err := methods.PerformDvsProductSpecOperation_Task(ctx, client, req)
	if err != nil {
		return err
	}
	task := object.NewTask(client.Client, resp.Returnval)
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return err
	}

	return nil
}

// updateDVSConfiguration contains the atomic update/wait operation for a DVS.
func updateDVSConfiguration(client *govmomi.Client, dvs *object.VmwareDistributedVirtualSwitch, spec *types.VMwareDVSConfigSpec) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	task, err := dvs.Reconfigure(ctx, spec)
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer tcancel()
	if err := task.Wait(tctx); err != nil {
		return err
	}
	return nil
}

// enableDVSNetworkResourceManagement exposes the
// EnableNetworkResourceManagement method of the DistributedVirtualSwitch MO.
// This local implementation may go away if this is exposed in the higher-level
// object upstream.
func enableDVSNetworkResourceManagement(client *govmomi.Client, dvs *object.VmwareDistributedVirtualSwitch, enabled bool) error {
	req := &types.EnableNetworkResourceManagement{
		This:   dvs.Reference(),
		Enable: enabled,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	_, err := methods.EnableNetworkResourceManagement(ctx, client, req)
	if err != nil {
		return err
	}

	return nil
}
