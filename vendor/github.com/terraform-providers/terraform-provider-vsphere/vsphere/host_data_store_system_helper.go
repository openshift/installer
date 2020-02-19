package vsphere

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/datastore"
	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/methods"
	"github.com/vmware/govmomi/vim25/types"
)

// hostDatastoreSystemFromHostSystemID locates a HostDatastoreSystem from a
// specified HostSystem managed object ID.
func hostDatastoreSystemFromHostSystemID(client *govmomi.Client, hsID string) (*object.HostDatastoreSystem, error) {
	hs, err := hostsystem.FromID(client, hsID)
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return hs.ConfigManager().DatastoreSystem(ctx)
}

// availableScsiDisk checks to make sure that a disk is available for use in a
// VMFS datastore, and returns the ScsiDisk.
func availableScsiDisk(dss *object.HostDatastoreSystem, name string) (*types.HostScsiDisk, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	disks, err := dss.QueryAvailableDisksForVmfs(ctx)
	if err != nil {
		return nil, fmt.Errorf("cannot query available disks: %s", err)
	}

	var disk *types.HostScsiDisk
	for _, d := range disks {
		if d.CanonicalName == name {
			disk = &d
			break
		}
	}
	if disk == nil {
		return nil, fmt.Errorf("%s does not seem to be a disk available for VMFS", name)
	}
	return disk, nil
}

// diskSpecForCreate checks to make sure that a disk is available to be used to
// create a VMFS datastore, specifically in its entirety, and returns a
// respective VmfsDatastoreCreateSpec.
func diskSpecForCreate(dss *object.HostDatastoreSystem, name string) (*types.VmfsDatastoreCreateSpec, error) {
	disk, err := availableScsiDisk(dss, name)
	if err != nil {
		return nil, err
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	options, err := dss.QueryVmfsDatastoreCreateOptions(ctx, disk.DevicePath)
	if err != nil {
		return nil, fmt.Errorf("could not get disk creation options for %q: %s", name, err)
	}
	var option *types.VmfsDatastoreOption
	for _, o := range options {
		if _, ok := o.Info.(*types.VmfsDatastoreAllExtentOption); ok {
			option = &o
			break
		}
	}
	if option == nil {
		return nil, fmt.Errorf("device %q is not available as a new whole-disk device for datastore", name)
	}
	return option.Spec.(*types.VmfsDatastoreCreateSpec), nil
}

// diskSpecForExtend checks to make sure that a disk is available to be
// used to extend a VMFS datastore, specifically in its entirety, and returns a
// respective VmfsDatastoreExtendSpec if it is. An error is returned if it's
// not.
func diskSpecForExtend(dss *object.HostDatastoreSystem, ds *object.Datastore, name string) (*types.VmfsDatastoreExtendSpec, error) {
	disk, err := availableScsiDisk(dss, name)
	if err != nil {
		return nil, err
	}

	props, err := datastore.Properties(ds)
	if err != nil {
		return nil, fmt.Errorf("error getting properties for datastore ID %q: %s", ds.Reference().Value, err)
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	options, err := queryVmfsDatastoreExtendOptions(ctx, dss, ds, disk.DevicePath, true)
	if err != nil {
		return nil, fmt.Errorf("could not get disk extension options for %q: %s", name, err)
	}
	var option *types.VmfsDatastoreOption
	for _, o := range options {
		if _, ok := o.Info.(*types.VmfsDatastoreAllExtentOption); ok {
			option = &o
			break
		}
	}
	if option == nil {
		return nil, fmt.Errorf("device %q cannot be used as a new whole-disk device for datastore %q", name, props.Summary.Name)
	}
	return option.Spec.(*types.VmfsDatastoreExtendSpec), nil
}

// removeDatastore is a convenience method for removing a referenced datastore.
func removeDatastore(s *object.HostDatastoreSystem, ds *object.Datastore) error {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	return s.Remove(ctx, ds)
}

// queryVmfsDatastoreExtendOptions is a stop-gap method that implements
// QueryVmfsDatastoreExtendOptions. It will be removed once the higher level
// HostDatastoreSystem object supports this method.
func queryVmfsDatastoreExtendOptions(ctx context.Context, s *object.HostDatastoreSystem, ds *object.Datastore, devicePath string, suppressExpandCandidates bool) ([]types.VmfsDatastoreOption, error) {
	req := types.QueryVmfsDatastoreExtendOptions{
		This:                     s.Reference(),
		Datastore:                ds.Reference(),
		DevicePath:               devicePath,
		SuppressExpandCandidates: &suppressExpandCandidates,
	}

	res, err := methods.QueryVmfsDatastoreExtendOptions(ctx, s.Client(), &req)
	if err != nil {
		return nil, err
	}

	return res.Returnval, nil
}

// extendVmfsDatastore is a stop-gap method that implements
// ExtendVmfsDatastore. It will be removed once the higher level
// HostDatastoreSystem object supports this method.
func extendVmfsDatastore(ctx context.Context, s *object.HostDatastoreSystem, ds *object.Datastore, spec types.VmfsDatastoreExtendSpec) (*object.Datastore, error) {
	req := types.ExtendVmfsDatastore{
		This:      s.Reference(),
		Datastore: ds.Reference(),
		Spec:      spec,
	}

	res, err := methods.ExtendVmfsDatastore(ctx, s.Client(), &req)
	if err != nil {
		return nil, err
	}

	return object.NewDatastore(s.Client(), res.Returnval), nil
}
