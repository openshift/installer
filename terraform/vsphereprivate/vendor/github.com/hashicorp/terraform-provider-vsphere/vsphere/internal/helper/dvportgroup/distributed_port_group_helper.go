package dvportgroup

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// MissingPortGroupReferenceError is an error that gets returned when a port
// group lookup succeeds but does not return a MO to a
// DistributedVirtualPortgroup.
type MissingPortGroupReferenceError struct {
	message string
}

// NewMissingPortGroupReferenceError returns a MissingPortGroupReferenceError
// with the supplied message.
func NewMissingPortGroupReferenceError(message string) error {
	return &MissingPortGroupReferenceError{message}
}

func (e *MissingPortGroupReferenceError) Error() string {
	return e.message
}

// FromKey gets a portgroup object from its key.
func FromKey(client *govmomi.Client, dvsUUID, pgKey string) (*object.DistributedVirtualPortgroup, error) {
	dvsm := types.ManagedObjectReference{Type: "DistributedVirtualSwitchManager", Value: "DVSManager"}
	req := &types.DVSManagerLookupDvPortGroup{
		This:         dvsm,
		SwitchUuid:   dvsUUID,
		PortgroupKey: pgKey,
	}
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	resp, err := methods.DVSManagerLookupDvPortGroup(ctx, client, req)
	if err != nil {
		return nil, err
	}

	if resp.Returnval == nil {
		return nil, NewMissingPortGroupReferenceError(
			fmt.Sprintf(
				"portgroup lookup by key returned nil result for DVS UUID %q and portgroup key %q",
				dvsUUID,
				pgKey,
			),
		)
	}

	return FromMOID(client, resp.Returnval.Reference().Value)
}

// FromMOID locates a portgroup by its managed object reference ID.
func FromMOID(client *govmomi.Client, id string) (*object.DistributedVirtualPortgroup, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "DistributedVirtualPortgroup",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	// Should be safe to return here. If our reference returned here and is not a
	// DistributedVirtualPortgroup, then we have bigger problems and to be
	// honest we should be panicking anyway.
	return ds.(*object.DistributedVirtualPortgroup), nil
}

// FromPath gets a portgroup object from its path.
func FromPath(client *govmomi.Client, name string, dc *object.Datacenter) (*object.DistributedVirtualPortgroup, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		finder.SetDatacenter(dc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	net, err := finder.Network(ctx, name)
	if err != nil {
		return nil, err
	}
	if net.Reference().Type != "DistributedVirtualPortgroup" {
		return nil, fmt.Errorf("network at path %q is not a portgroup (type %s)", name, net.Reference().Type)
	}
	return FromMOID(client, net.Reference().Value)
}

// Properties is a convenience method that wraps fetching the
// portgroup MO from its higher-level object.
func Properties(pg *object.DistributedVirtualPortgroup) (*mo.DistributedVirtualPortgroup, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.DistributedVirtualPortgroup
	if err := pg.Properties(ctx, pg.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// Create exposes the CreateDVPortgroup_Task method of the
// DistributedVirtualSwitch MO.  This local implementation may go away if this
// is exposed in the higher-level object upstream.
func Create(client *govmomi.Client, dvs *object.VmwareDistributedVirtualSwitch, spec types.DVPortgroupConfigSpec) (*object.Task, error) {
	req := &types.CreateDVPortgroup_Task{
		This: dvs.Reference(),
		Spec: spec,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	resp, err := methods.CreateDVPortgroup_Task(ctx, client, req)
	if err != nil {
		return nil, err
	}

	return object.NewTask(client.Client, resp.Returnval.Reference()), nil
}
