package folder

import (
	"context"
	"errors"
	"fmt"
	"log"
	"path"
	"strings"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/viapi"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// VSphereFolderType is an enumeration type for vSphere folder types.
type VSphereFolderType string

// The following are constants for the 5 vSphere folder types - these are used
// to help determine base paths and also to validate folder types in the
// vsphere_folder resource.
const (
	VSphereFolderTypeVM        = VSphereFolderType("vm")
	VSphereFolderTypeNetwork   = VSphereFolderType("network")
	VSphereFolderTypeHost      = VSphereFolderType("host")
	VSphereFolderTypeDatastore = VSphereFolderType("datastore")

	// VSphereFolderTypeDatacenter is a special folder type - it does not get a
	// root path particle generated for it as it is an integral part of the path
	// generation process, but is defined so that it can be properly referenced
	// and used in validation.
	VSphereFolderTypeDatacenter = VSphereFolderType("datacenter")
)

// RootPathParticle is the section of a vSphere inventory path that denotes a
// specific kind of inventory item.
type RootPathParticle VSphereFolderType

// String implements Stringer for RootPathParticle.
func (p RootPathParticle) String() string {
	return string(p)
}

// Delimiter returns the path delimiter for the particle, which is basically
// just a particle with a leading slash.
func (p RootPathParticle) Delimiter() string {
	return string("/" + p)
}

// RootFromDatacenter returns the root path for the particle from the given
// datacenter's inventory path.
func (p RootPathParticle) RootFromDatacenter(dc *object.Datacenter) string {
	return dc.InventoryPath + "/" + string(p)
}

// PathFromDatacenter returns the combined result of RootFromDatacenter plus a
// relative path for a given particle and datacenter object.
func (p RootPathParticle) PathFromDatacenter(dc *object.Datacenter, relative string) string {
	return p.RootFromDatacenter(dc) + "/" + relative
}

// SplitDatacenter is a convenience method that splits out the datacenter path
// from the supplied path for the particle.
func (p RootPathParticle) SplitDatacenter(inventoryPath string) (string, error) {
	s := strings.SplitN(inventoryPath, p.Delimiter(), 2)
	if len(s) != 2 {
		return inventoryPath, fmt.Errorf("could not split path %q on %q", inventoryPath, p.Delimiter())
	}
	return s[0], nil
}

// SplitRelative is a convenience method that splits out the relative path from
// the supplied path for the particle.
func (p RootPathParticle) SplitRelative(inventoryPath string) (string, error) {
	s := strings.SplitN(inventoryPath, p.Delimiter(), 2)
	if len(s) != 2 {
		return inventoryPath, fmt.Errorf("could not split path %q on %q", inventoryPath, p.Delimiter())
	}
	return s[1], nil
}

// SplitRelativeFolder is a convenience method that returns the parent folder
// for the result of SplitRelative on the supplied path.
//
// This is generally useful to get the folder for a managed entity, versus getting a full relative path. If you want that, use SplitRelative instead.
func (p RootPathParticle) SplitRelativeFolder(inventoryPath string) (string, error) {
	relative, err := p.SplitRelative(inventoryPath)
	if err != nil {
		return inventoryPath, err
	}
	return path.Dir(relative), nil
}

// NewRootFromPath takes the datacenter path for a specific entity, and then
// appends the new particle supplied.
func (p RootPathParticle) NewRootFromPath(inventoryPath string, newParticle RootPathParticle) (string, error) {
	dcPath, err := p.SplitDatacenter(inventoryPath)
	if err != nil {
		return inventoryPath, err
	}
	return fmt.Sprintf("%s/%s", dcPath, newParticle), nil
}

// PathFromNewRoot takes the datacenter path for a specific entity, and then
// appends the new particle supplied with the new relative path.
//
// As an example, consider a supplied host path "/dc1/host/cluster1/esxi1", and
// a supplied datastore folder relative path of "/foo/bar".  This function will
// split off the datacenter section of the path (/dc1) and combine it with the
// datastore folder with the proper delimiter. The resulting path will be
// "/dc1/datastore/foo/bar".
func (p RootPathParticle) PathFromNewRoot(inventoryPath string, newParticle RootPathParticle, relative string) (string, error) {
	rootPath, err := p.NewRootFromPath(inventoryPath, newParticle)
	if err != nil {
		return inventoryPath, err
	}
	return path.Clean(fmt.Sprintf("%s/%s", rootPath, relative)), nil
}

const (
	// RootPathParticleVM provides root path parsing functionality for VM paths.
	RootPathParticleVM = RootPathParticle(VSphereFolderTypeVM)

	// RootPathParticleNetwork provides root path parsing functionality for
	// network paths.
	RootPathParticleNetwork = RootPathParticle(VSphereFolderTypeNetwork)

	// RootPathParticleHost provides root path parsing functionality for host and
	// cluster paths.
	RootPathParticleHost = RootPathParticle(VSphereFolderTypeHost)

	// RootPathParticleDatastore provides root path parsing functionality for
	// datastore paths.
	RootPathParticleDatastore = RootPathParticle(VSphereFolderTypeDatastore)
)

// datacenterPathFromHostSystemID returns the datacenter section of a
// HostSystem's inventory path.
func datacenterPathFromHostSystemID(client *govmomi.Client, hsID string) (string, error) {
	hs, err := hostsystem.FromID(client, hsID)
	if err != nil {
		return "", err
	}
	return RootPathParticleHost.SplitDatacenter(hs.InventoryPath)
}

// datastoreRootPathFromHostSystemID returns the root datastore folder path
// for a specific host system ID.
func datastoreRootPathFromHostSystemID(client *govmomi.Client, hsID string) (string, error) {
	hs, err := hostsystem.FromID(client, hsID)
	if err != nil {
		return "", err
	}
	return RootPathParticleHost.NewRootFromPath(hs.InventoryPath, RootPathParticleDatastore)
}

// FromAbsolutePath returns an *object.Folder from a given absolute path.
// If no such folder is found, an appropriate error will be returned.
func FromAbsolutePath(client *govmomi.Client, path string) (*object.Folder, error) {
	finder := find.NewFinder(client.Client, false)
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	folder, err := finder.Folder(ctx, path)
	if err != nil {
		return nil, err
	}
	return folder, nil
}

// folderFromObject returns an *object.Folder from a given object of specific
// types, and relative path of a type defined in folderType. If no such folder
// is found, an appropriate error will be returned.
//
// The list of supported object types will grow as the provider supports more
// resources.
func folderFromObject(client *govmomi.Client, obj interface{}, folderType RootPathParticle, relative string) (*object.Folder, error) {
	// If we are using this for anything else other than the root folder on ESXi,
	// return an error.
	if err := viapi.ValidateVirtualCenter(client); err != nil && relative != "" {
		return nil, errors.New("folders are only supported vCenter only")
	}
	var p string
	var err error
	switch o := obj.(type) {
	case *object.VmwareDistributedVirtualSwitch:
		p, err = RootPathParticleNetwork.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.Datastore:
		p, err = RootPathParticleDatastore.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.StoragePod:
		p, err = RootPathParticleDatastore.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.HostSystem:
		p, err = RootPathParticleHost.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.ResourcePool:
		p, err = RootPathParticleHost.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.ComputeResource:
		p, err = RootPathParticleHost.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.ClusterComputeResource:
		p, err = RootPathParticleHost.PathFromNewRoot(o.InventoryPath, folderType, relative)
	case *object.VirtualMachine:
		p, err = RootPathParticleVM.PathFromNewRoot(o.InventoryPath, folderType, relative)
	default:
		return nil, fmt.Errorf("unsupported object type %T", o)
	}
	if err != nil {
		return nil, err
	}
	return FromAbsolutePath(client, p)
}

// DatastoreFolderFromObject returns an *object.Folder from a given object,
// and relative datastore folder path. If no such folder is found, of if it is
// not a datastore folder, an appropriate error will be returned.
func DatastoreFolderFromObject(client *govmomi.Client, obj interface{}, relative string) (*object.Folder, error) {
	folder, err := folderFromObject(client, obj, RootPathParticleDatastore, relative)
	if err != nil {
		return nil, err
	}

	return validateDatastoreFolder(folder)
}

// HostFolderFromObject returns an *object.Folder from a given object, and
// relative host folder path. If no such folder is found, or if it is not a
// host folder, an appropriate error will be returned.
func HostFolderFromObject(client *govmomi.Client, obj interface{}, relative string) (*object.Folder, error) {
	folder, err := folderFromObject(client, obj, RootPathParticleHost, relative)
	if err != nil {
		return nil, err
	}

	return validateHostFolder(folder)
}

// VirtualMachineFolderFromObject returns an *object.Folder from a given
// object, and relative datastore folder path. If no such folder is found, or
// if it is not a VM folder, an appropriate error will be returned.
func VirtualMachineFolderFromObject(client *govmomi.Client, obj interface{}, relative string) (*object.Folder, error) {
	log.Printf("[DEBUG] Locating folder at path %q relative to virtual machine root", relative)
	folder, err := folderFromObject(client, obj, RootPathParticleVM, relative)
	if err != nil {
		return nil, err
	}

	return validateVirtualMachineFolder(folder)
}

// networkFolderFromObject returns an *object.Folder from a given object,
// and relative network folder path. If no such folder is found, of if it is
// not a network folder, an appropriate error will be returned.
func networkFolderFromObject(client *govmomi.Client, obj interface{}, relative string) (*object.Folder, error) {
	folder, err := folderFromObject(client, obj, RootPathParticleNetwork, relative)
	if err != nil {
		return nil, err
	}

	return validateNetworkFolder(folder)
}

// validateDatastoreFolder checks to make sure the folder is a datastore
// folder, and returns it if it is, or an error if it isn't.
func validateDatastoreFolder(folder *object.Folder) (*object.Folder, error) {
	ft, err := FindType(folder)
	if err != nil {
		return nil, err
	}
	if ft != VSphereFolderTypeDatastore {
		return nil, fmt.Errorf("%q is not a datastore folder", folder.InventoryPath)
	}
	return folder, nil
}

// validateHostFolder checks to make sure the folder is a host
// folder, and returns it if it is, or an error if it isn't.
func validateHostFolder(folder *object.Folder) (*object.Folder, error) {
	ft, err := FindType(folder)
	if err != nil {
		return nil, err
	}
	if ft != VSphereFolderTypeHost {
		return nil, fmt.Errorf("%q is not a host folder", folder.InventoryPath)
	}
	return folder, nil
}

// validateVirtualMachineFolder checks to make sure the folder is a VM folder,
// and returns it if it is, or an error if it isn't.
func validateVirtualMachineFolder(folder *object.Folder) (*object.Folder, error) {
	ft, err := FindType(folder)
	if err != nil {
		return nil, err
	}
	if ft != VSphereFolderTypeVM {
		return nil, fmt.Errorf("%q is not a VM folder", folder.InventoryPath)
	}
	log.Printf("[DEBUG] Folder located: %q", folder.InventoryPath)
	return folder, nil
}

// validateNetworkFolder checks to make sure the folder is a network folder,
// and returns it if it is, or an error if it isn't.
func validateNetworkFolder(folder *object.Folder) (*object.Folder, error) {
	ft, err := FindType(folder)
	if err != nil {
		return nil, err
	}
	if ft != VSphereFolderTypeNetwork {
		return nil, fmt.Errorf("%q is not a network folder", folder.InventoryPath)
	}
	return folder, nil
}

// PathIsEmpty checks a folder path to see if it's "empty" (ie: would resolve
// to the root inventory path for a given type in a datacenter - "" or "/").
func PathIsEmpty(path string) bool {
	return path == "" || path == "/"
}

// NormalizePath is a SchemaStateFunc that normalizes a folder path.
func NormalizePath(v interface{}) string {
	p := v.(string)
	if PathIsEmpty(p) {
		return ""
	}
	return strings.TrimPrefix(path.Clean(p), "/")
}

// MoveObjectTo moves a object by reference into a folder.
func MoveObjectTo(ref types.ManagedObjectReference, folder *object.Folder) error {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := folder.MoveInto(ctx, []types.ManagedObjectReference{ref})
	if err != nil {
		return err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	return task.Wait(tctx)
}

// FromPath takes a relative folder path, an object type, and an optional
// supplied datacenter, and returns the respective *object.Folder if it exists.
//
// The datacenter supplied in dc cannot be nil if the folder type supplied by
// ft is something else other than VSphereFolderTypeDatacenter.
func FromPath(c *govmomi.Client, p string, ft VSphereFolderType, dc *object.Datacenter) (*object.Folder, error) {
	var fp string
	if ft == VSphereFolderTypeDatacenter {
		fp = "/" + p
	} else {
		pt := RootPathParticle(ft)
		fp = pt.PathFromDatacenter(dc, p)
	}
	return FromAbsolutePath(c, fp)
}

// ParentFromPath takes a relative object path (usually a folder), an
// object type, and an optional supplied datacenter, and returns the parent
// *object.Folder if it exists.
//
// The datacenter supplied in dc cannot be nil if the folder type supplied by
// ft is something else other than VSphereFolderTypeDatacenter.
func ParentFromPath(c *govmomi.Client, p string, ft VSphereFolderType, dc *object.Datacenter) (*object.Folder, error) {
	return FromPath(c, path.Dir(p), ft, dc)
}

// FromID locates a Folder by its managed object reference ID.
func FromID(client *govmomi.Client, id string) (*object.Folder, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "Folder",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	folder, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	return folder.(*object.Folder), nil
}

func List(client *govmomi.Client) ([]*object.Folder, error) {
	return getFolders(client, "/*")
}

func getFolders(client *govmomi.Client, path string) ([]*object.Folder, error) {
	ctx := context.TODO()
	var folders []*object.Folder
	finder := find.NewFinder(client.Client, false)
	es, err := finder.ManagedObjectListChildren(ctx, path+"/*", "folder")
	if err != nil {
		return nil, err
	}
	for _, id := range es {
		switch {
		case id.Object.Reference().Type == "Folder":
			newFolders, err := getFolders(client, id.Path)
			if err != nil {
				return nil, err
			}
			folders = append(folders, newFolders...)
		default:
			continue
		}
	}
	return folders, nil
}

// Properties is a convenience method that wraps fetching the
// Folder MO from its higher-level object.
func Properties(folder *object.Folder) (*mo.Folder, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.Folder
	if err := folder.Properties(ctx, folder.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// FindType returns a proper VSphereFolderType for a folder object by checking its child type.
func FindType(folder *object.Folder) (VSphereFolderType, error) {
	var ft VSphereFolderType

	props, err := Properties(folder)
	if err != nil {
		return ft, err
	}

	// Depending on the container type, the actual folder type may be contained
	// in either the first or second element, the former for clusters, datastore
	// clusters, or standalone ESXi for VMs, and the latter in the case of actual
	// folders that can contain subfolders.
	var ct string
	if props.ChildType[0] != "Folder" {
		ct = props.ChildType[0]
	} else {
		ct = props.ChildType[1]
	}

	switch ct {
	case "Datacenter":
		ft = VSphereFolderTypeDatacenter
	case "ComputeResource":
		ft = VSphereFolderTypeHost
	case "VirtualMachine":
		ft = VSphereFolderTypeVM
	case "Datastore":
		ft = VSphereFolderTypeDatastore
	case "Network":
		ft = VSphereFolderTypeNetwork
	default:
		return ft, fmt.Errorf("unknown folder type: %#v", ct)
	}

	return ft, nil
}

// HasChildren checks to see if a folder has any child items and returns
// true if that is the case. This is useful when checking to see if a folder is
// safe to delete - destroying a folder in vSphere destroys *all* children if
// at all possible (including removing virtual machines), so extra verification
// is necessary to prevent accidental removal.
func HasChildren(f *object.Folder) (bool, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	children, err := f.Children(ctx)
	if err != nil {
		return false, err
	}
	return len(children) > 0, nil
}
