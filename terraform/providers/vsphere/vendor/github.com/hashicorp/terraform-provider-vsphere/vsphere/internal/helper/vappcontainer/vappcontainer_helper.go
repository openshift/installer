package vappcontainer

import (
	"context"
	"log"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// FromPath returns a VirtualApp via its supplied path.
func FromPath(client *govmomi.Client, name string, dc *object.Datacenter) (*object.VirtualApp, error) {
	finder := find.NewFinder(client.Client, false)

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	if dc != nil {
		finder.SetDatacenter(dc)
	}
	return finder.VirtualApp(ctx, name)
}

// FromID locates a VirtualApp by its managed object reference ID.
func FromID(client *govmomi.Client, id string) (*object.VirtualApp, error) {
	log.Printf("[DEBUG] Locating vApp container with ID %s", id)
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "VirtualApp",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	obj, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	log.Printf("[DEBUG] vApp container found: %s", obj.Reference().Value)
	return obj.(*object.VirtualApp), nil
}

// IsVApp checks if a given managed object ID is a vApp. This is useful
// deciding if a given resource pool is a vApp or a standard resource pool.
func IsVApp(client *govmomi.Client, rp string) bool {
	_, err := FromID(client, rp)
	if err != nil {
		return false
	}
	return true
}

// Properties returns the VirtualApp managed object from its higher-level
// object.
func Properties(obj *object.VirtualApp) (*mo.VirtualApp, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.VirtualApp
	if err := obj.Properties(ctx, obj.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// Create creates a VirtualApp.
func Create(rp *object.ResourcePool, name string, resSpec *types.ResourceConfigSpec, vSpec *types.VAppConfigSpec, folder *object.Folder) (*object.VirtualApp, error) {
	log.Printf("[DEBUG] Creating vApp container %s/%s", rp.InventoryPath, name)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	nva, err := rp.CreateVApp(ctx, name, *resSpec, *vSpec, folder)
	if err != nil {
		return nil, err
	}
	return nva, nil
}

// Update updates a VirtualApp.
func Update(vc *object.VirtualApp, spec types.VAppConfigSpec) error {
	log.Printf("[DEBUG] Updating vApp container %q", vc.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return vc.UpdateConfig(ctx, spec)
}

// Delete destroys a VirtualApp.
func Delete(vc *object.VirtualApp) error {
	log.Printf("[DEBUG] Deleting vApp container %q", vc.InventoryPath)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := vc.Destroy(ctx)
	if err != nil {
		return err
	}
	return task.Wait(ctx)
}

// HasChildren checks to see if a vApp container has any child items (virtual
// machines, vApps, or resource pools) and returns true if that is the case.
// This is useful when checking to see if a vApp container is safe to delete.
// Destroying a vApp container in vSphere destroys *all* children if at all
// possible, so extra verification is necessary to prevent accidental removal.
func HasChildren(vc *object.VirtualApp) (bool, error) {
	props, err := Properties(vc)
	if err != nil {
		return false, err
	}
	if len(props.Vm) > 0 || len(props.ResourcePool.ResourcePool) > 0 {
		return true, nil
	}
	return false, nil
}
