package vsphere

import (
	"context"
	"fmt"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// getDatacenter gets the higher-level datacenter object for the datacenter
// name supplied by dc.
//
// The default datacenter is denoted by using an empty string. When working
// with ESXi directly, the default datacenter is always selected.
func getDatacenter(c *govmomi.Client, dc string) (*object.Datacenter, error) {
	finder := find.NewFinder(c.Client, true)
	t := c.ServiceContent.About.ApiType
	switch t {
	case "HostAgent":
		return finder.DefaultDatacenter(context.TODO())
	case "VirtualCenter":
		if dc != "" {
			return finder.Datacenter(context.TODO(), dc)
		}
		return finder.DefaultDatacenter(context.TODO())
	}
	return nil, fmt.Errorf("unsupported ApiType: %s", t)
}

// datacenterFromID locates a Datacenter by its managed object reference ID.
func datacenterFromID(client *govmomi.Client, id string) (*object.Datacenter, error) {
	finder := find.NewFinder(client.Client, false)

	ref := types.ManagedObjectReference{
		Type:  "Datacenter",
		Value: id,
	}

	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	ds, err := finder.ObjectReference(ctx, ref)
	if err != nil {
		return nil, fmt.Errorf("could not find datacenter with id: %s: %s", id, err)
	}
	return ds.(*object.Datacenter), nil
}

func datacenterCustomAttributes(dc *object.Datacenter) (*mo.Datacenter, error) {
	ctx, cancel := context.WithTimeout(context.Background(), defaultAPITimeout)
	defer cancel()
	var props mo.Datacenter
	if err := dc.Properties(ctx, dc.Reference(), []string{"customValue"}, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

func listDatacenters(client *govmomi.Client) ([]*object.Datacenter, error) {
	return dcsByPath(client, "/")
}

func dcsByPath(client *govmomi.Client, path string) ([]*object.Datacenter, error) {
	ctx := context.TODO()
	var dcs []*object.Datacenter
	finder := find.NewFinder(client.Client, false)
	if path != "/" {
		path = path + "/*"
	}
	es, err := finder.ManagedObjectListChildren(ctx, path, "datacenter", "folder")
	if err != nil {
		return nil, err
	}
	for _, id := range es {
		switch {
		case id.Object.Reference().Type == "Datacenter":
			dc, err := datacenterFromID(client, id.Object.Reference().Value)
			if err != nil {
				return nil, err
			}
			dcs = append(dcs, dc)
		case id.Object.Reference().Type == "Folder":
			newDCs, err := dcsByPath(client, id.Path)
			if err != nil {
				return nil, err
			}
			dcs = append(dcs, newDCs...)
		default:
			continue
		}
	}
	return dcs, nil
}
