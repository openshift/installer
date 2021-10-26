package vsphere

import (
	"context"

	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25"
)

// Finder interface represents the client that is used to connect to VSphere to get specific
// information from the resources in the VCenter. This interface just describes all the useful
// functions used by the installer from the finder function in vmware govmomi package and is
// mostly used to create a mock client that can be used for testing.
type Finder interface {
	Datacenter(ctx context.Context, path string) (*object.Datacenter, error)
	DatacenterList(ctx context.Context, path string) ([]*object.Datacenter, error)
	DatastoreList(ctx context.Context, path string) ([]*object.Datastore, error)
	ClusterComputeResource(ctx context.Context, path string) (*object.ClusterComputeResource, error)
	ClusterComputeResourceList(ctx context.Context, path string) ([]*object.ClusterComputeResource, error)
	Folder(ctx context.Context, path string) (*object.Folder, error)
	NetworkList(ctx context.Context, path string) ([]object.NetworkReference, error)
	Network(ctx context.Context, path string) (object.NetworkReference, error)
	ResourcePool(ctx context.Context, path string) (*object.ResourcePool, error)
}

// NewFinder creates a new client that conforms with the Finder interface and returns a
// vmware govmomi finder object that can be used to search for resources in vsphere.
func NewFinder(client *vim25.Client, all ...bool) Finder {
	return find.NewFinder(client, all...)
}
