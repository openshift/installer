package vsphere

import (
	"context"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/soap"
	"github.com/vmware/govmomi/vim25/types"
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
	VirtualMachine(ctx context.Context, path string) (*object.VirtualMachine, error)
	VirtualMachineList(ctx context.Context, path string) ([]*object.VirtualMachine, error)
	HostSystemList(ctx context.Context, path string) ([]*object.HostSystem, error)
}

// NewFinder creates a new client that conforms with the Finder interface and returns a
// vmware govmomi finder object that can be used to search for resources in vsphere.
func NewFinder(client *vim25.Client, all ...bool) Finder {
	return find.NewFinder(client, all...)
}

// ClientLogout is empty function that logs out of vSphere clients
type ClientLogout func()

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string) (*vim25.Client, *rest.Client, ClientLogout, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, nil, err
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(ctx, u, false)

	if err != nil {
		return nil, nil, nil, err
	}

	restClient := rest.NewClient(c.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		logoutErr := c.Logout(context.TODO())
		if logoutErr != nil {
			err = logoutErr
		}
		return nil, nil, nil, err
	}

	return c.Client, restClient, func() {
		c.Logout(context.TODO())
		restClient.Logout(context.TODO())
	}, nil
}

// getNetworks returns a slice of Managed Object references for networks in the given vSphere Cluster.
func getNetworks(ctx context.Context, ccr *object.ClusterComputeResource) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()
	var ccrMo mo.ClusterComputeResource

	err := ccr.Properties(ctx, ccr.Reference(), []string{"network"}, &ccrMo)
	if err != nil {
		return nil, errors.Wrap(err, "could not get properties of cluster")
	}
	return ccrMo.Network, nil
}

// GetClusterNetworks returns a slice of Managed Object references for vSphere networks in the given Datacenter
// and Cluster.
func GetClusterNetworks(ctx context.Context, finder Finder, datacenter, cluster string) ([]types.ManagedObjectReference, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	ccr, err := finder.ClusterComputeResource(context.TODO(), cluster)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find vSphere cluster at %s", cluster)
	}

	// Get list of Networks inside vSphere Cluster
	networks, err := getNetworks(ctx, ccr)
	if err != nil {
		return nil, err
	}

	return networks, nil
}

// GetNetworkName returns the name of a vSphere network given its Managed Object reference.
func GetNetworkName(ctx context.Context, client *vim25.Client, ref types.ManagedObjectReference) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	netObj := object.NewNetwork(client, ref)
	name, err := netObj.ObjectName(ctx)
	if err != nil {
		return "", errors.Wrapf(err, "could not get network name for %s", ref.String())
	}
	return name, nil
}

// GetNetworkMo returns the unique Managed Object for given network name inside of the given Datacenter
// and Cluster.
func GetNetworkMo(ctx context.Context, client *vim25.Client, finder Finder, datacenter, cluster, network string) (*types.ManagedObjectReference, error) {
	networks, err := GetClusterNetworks(ctx, finder, datacenter, cluster)
	if err != nil {
		return nil, err
	}
	for _, net := range networks {
		name, err := GetNetworkName(ctx, client, net)
		if err != nil {
			return nil, err
		}
		if name == network {
			return &net, nil
		}
	}

	return nil, errors.Errorf("unable to find network provided")
}

// GetNetworkMoID returns the unique Managed Object ID for given network name inside of the given Datacenter
// and Cluster.
func GetNetworkMoID(ctx context.Context, client *vim25.Client, finder Finder, datacenter, cluster, network string) (string, error) {
	mo, err := GetNetworkMo(ctx, client, finder, datacenter, cluster, network)
	if err != nil {
		return "", err
	}
	return mo.Value, nil
}
