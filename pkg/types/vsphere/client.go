package vsphere

import (
	"context"
	"fmt"
	"net/url"
	"time"

	"github.com/pkg/errors"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"

	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/vapi/rest"
	"github.com/vmware/govmomi/vim25"
	"github.com/vmware/govmomi/vim25/soap"
)

//go:generate mockgen -source=./client.go -destination=mock/vsphereclient_generated.go -package=mock

// CreateVSphereClients creates the SOAP and REST client to access
// different portions of the vSphere API
// e.g. tags are only available in REST
func CreateVSphereClients(ctx context.Context, vcenter, username, password string) (*vim25.Client, *rest.Client, error) {
	ctx, cancel := context.WithTimeout(ctx, 60*time.Second)
	defer cancel()

	u, err := soap.ParseURL(vcenter)
	if err != nil {
		return nil, nil, err
	}
	u.User = url.UserPassword(username, password)
	c, err := govmomi.NewClient(ctx, u, false)

	if err != nil {
		return nil, nil, err
	}

	restClient := rest.NewClient(c.Client)
	err = restClient.Login(ctx, u.User)
	if err != nil {
		return nil, nil, err
	}

	return c.Client, restClient, nil
}

// NetworkIdentifier interface used to help identify a Network's Managed
// Object ID. Note that this interface is mocked as well.
// See also go:generate at top of file.
type NetworkIdentifier interface {
	GetNetworkName(ref types.ManagedObjectReference) (string, error)
	GetNetworks(ccr *object.ClusterComputeResource) ([]types.ManagedObjectReference, error)
}

// NetworkUtil is the runtime implementation of NetworkIdentifier.
type NetworkUtil struct {
	client *vim25.Client
}

func (n *NetworkUtil) GetNetworkName(ref types.ManagedObjectReference) (string, error) {
	netObj := object.NewNetwork(n.client, ref)
	name, err := netObj.ObjectName(context.TODO())
	if err != nil {
		return "", errors.Wrap(err, "could not get network name")
	}
	return name, nil
}

func (n *NetworkUtil) GetNetworks(ccr *object.ClusterComputeResource) ([]types.ManagedObjectReference, error) {
	var ccrMo mo.ClusterComputeResource
	err := ccr.Properties(context.TODO(), ccr.Reference(), []string{"network"}, &ccrMo)
	if err != nil {
		return nil, errors.Wrap(err, "could not get properties of cluster")
	}
	return ccrMo.Network, nil
}

func NewNetworkUtil(client *vim25.Client) NetworkIdentifier {
	return &NetworkUtil{client: client}
}

func GetClusterNetworks(networkIdentifier NetworkIdentifier, finder Finder, datacenter, cluster string) ([]types.ManagedObjectReference, error) {
	// Get vSphere Cluster resource in the given Datacenter.
	path := fmt.Sprintf("/%s/host/%s", datacenter, cluster)
	ccr, err := finder.ClusterComputeResource(context.TODO(), path)
	if err != nil {
		return nil, errors.Wrapf(err, "could not find vSphere cluster")
	}

	// Get list of Networks inside vSphere Cluster
	networks, err := networkIdentifier.GetNetworks(ccr)
	if err != nil {
		return nil, err
	}

	return networks, nil
}

func GetNetworkMoID(networkIdentifier NetworkIdentifier, finder Finder, datacenter, cluster, network string) (string, error) {
	networks, err := GetClusterNetworks(networkIdentifier, finder, datacenter, cluster)
	if err != nil {
		return "", err
	}

	for _, net := range networks {
		name, err := networkIdentifier.GetNetworkName(net)
		if err != nil {
			return "", errors.Wrap(err, "could not get network name")
		}
		if name == network {
			return net.Value, nil
		}
	}

	return "", errors.Errorf("unable to find network provided")
}

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
