package datastore

import (
	"context"
	"fmt"
	"log"
	"path"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/folder"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/hostsystem"
	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// FromID locates a Datastore by its managed object reference ID.
func FromID(client *govmomi.Client, id string) (*object.Datastore, error) {
	log.Printf("[DEBUG] Locating datastore with ID %q", id)
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "Datastore",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, err
	}
	// Should be safe to return here. If our reference returned here and is not a
	// datastore, then we have bigger problems and to be honest we should be
	// panicking anyway.
	log.Printf("[DEBUG] Datastore with ID %q found", ds.Reference().Value)
	return ds.(*object.Datastore), nil
}

// FromPath loads a datastore from its path. The datacenter is optional if the
// path is specific enough to not require it.
func FromPath(client *govmomi.Client, name string, dc *object.Datacenter) (*object.Datastore, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		finder.SetDatacenter(dc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return finder.Datastore(ctx, name)
}

func List(client *govmomi.Client) ([]*object.Datastore, error) {
	return getDatastores(client, "/*")
}

func getDatastores(client *govmomi.Client, path string) ([]*object.Datastore, error) {
	ctx := context.TODO()
	var dss []*object.Datastore
	finder := find.NewFinder(client.Client, false)
	es, err := finder.ManagedObjectListChildren(ctx, path+"/*", "datastore", "folder", "storagepod")
	if err != nil {
		return nil, err
	}
	for _, id := range es {
		switch {
		case id.Object.Reference().Type == "Datastore":
			ds, err := FromID(client, id.Object.Reference().Value)
			if err != nil {
				return nil, err
			}
			dss = append(dss, ds)
		case id.Object.Reference().Type == "Folder" || id.Object.Reference().Type == "Storagepod":
			newDSs, err := getDatastores(client, id.Path)
			if err != nil {
				return nil, err
			}
			dss = append(dss, newDSs...)
		default:
			continue
		}
	}
	return dss, nil
}

// Properties is a convenience method that wraps fetching the
// Datastore MO from its higher-level object.
func Properties(ds *object.Datastore) (*mo.Datastore, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.Datastore
	if err := ds.Properties(ctx, ds.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

func Unmount(client *govmomi.Client, ds *object.Datastore) error {
	dsprops, err := Properties(ds)
	if err != nil {
		return err
	}
	for _, h := range dsprops.Host {
		host, err := hostsystem.FromID(client, h.Key.Value)
		if err != nil {
			return err
		}
		hds, err := host.ConfigManager().DatastoreSystem(context.TODO())
		if err != nil {
			return err
		}
		err = hds.Remove(context.TODO(), ds)
		if err != nil {
			return err
		}
	}
	return nil
}

// MoveToFolder is a complex method that moves a datastore to a given
// relative datastore folder path. "Relative" here means relative to a
// datacenter, which is discovered from the current datastore path.
func MoveToFolder(client *govmomi.Client, ds *object.Datastore, relative string) error {
	f, err := folder.DatastoreFolderFromObject(client, ds, relative)
	if err != nil {
		return err
	}
	return folder.MoveObjectTo(ds.Reference(), f)
}

// MoveToFolderRelativeHostSystemID is a complex method that moves a
// datastore to a given datastore path, similar to MoveToFolder,
// except the path is relative to a HostSystem supplied by ID instead of the
// datastore.
func MoveToFolderRelativeHostSystemID(client *govmomi.Client, ds *object.Datastore, hsID, relative string) error {
	hs, err := hostsystem.FromID(client, hsID)
	if err != nil {
		return err
	}
	f, err := folder.DatastoreFolderFromObject(client, hs, relative)
	if err != nil {
		return err
	}
	return folder.MoveObjectTo(ds.Reference(), f)
}

// Browser returns the HostDatastoreBrowser for a certain datastore. This is a
// convenience method that exists to abstract the context.
func Browser(ds *object.Datastore) (*object.HostDatastoreBrowser, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return ds.Browser(ctx)
}

// SearchDatastore searches a datastore using the supplied HostDatastoreBrowser
// and a supplied path. The current implementation only returns the basic
// information, so all FileQueryFlags set, but not any flags for specific types
// of files.
func SearchDatastore(ds *object.Datastore, name string) ([]*types.FileInfo, error) {
	result, err := searchDatastore(ds, name)
	if err != nil {
		return nil, err
	}
	var files []*types.FileInfo
	for _, bfi := range result.File {
		files = append(files, bfi.GetFileInfo())
	}
	return files, nil
}

func searchDatastore(ds *object.Datastore, name string) (*types.HostDatastoreBrowserSearchResults, error) {
	browser, err := Browser(ds)
	if err != nil {
		return nil, err
	}
	var p, m string

	switch {
	case path.Dir(name) == ".":
		fallthrough
	case path.Base(name) == "":
		p = name
		m = "*"
	default:
		p = path.Dir(name)
		m = path.Base(name)
	}
	dp := &object.DatastorePath{
		Datastore: ds.Name(),
		Path:      p,
	}
	spec := &types.HostDatastoreBrowserSearchSpec{
		MatchPattern: []string{m},
		Details: &types.FileQueryFlags{
			FileType:     true,
			FileSize:     true,
			FileOwner:    types.NewBool(true),
			Modification: true,
		},
	}
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	task, err := browser.SearchDatastore(ctx, dp.String(), spec)
	if err != nil {
		return nil, err
	}
	tctx, tcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer tcancel()
	info, err := task.WaitForResult(tctx, nil)
	if err != nil {
		return nil, err
	}
	r := info.Result.(types.HostDatastoreBrowserSearchResults)
	return &r, nil
}

// FileExists takes a path in the datastore and checks to see if it exists.
//
// The path should be a bare path, not a datastore path. Globs are not allowed.
func FileExists(ds *object.Datastore, name string) (bool, error) {
	files, err := SearchDatastore(ds, name)
	if err != nil {
		return false, err
	}
	if len(files) > 1 {
		return false, fmt.Errorf("multiple results returned for %q in datastore %q, use a more specific search", name, ds)
	}
	if len(files) < 1 {
		return false, nil
	}
	return path.Base(name) == files[0].Path, nil
}
