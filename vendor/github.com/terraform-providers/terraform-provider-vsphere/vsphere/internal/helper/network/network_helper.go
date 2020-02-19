package network

import (
	"context"
	"fmt"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
)

// FromPath loads a network via its path.
//
// A network is a usually one of three kinds of networks: a DVS port group, a
// host port group, or a "opaque" network, provided externally from something
// like NSX. All three of these can be used as a backing for a virtual ethernet
// card, which is usually what these helpers are used with.
//
// Datacenter is optional here - if not provided, it's expected that the path
// is sufficient enough for finder to determine the datacenter required.
func FromPath(client *govmomi.Client, name string, dc *object.Datacenter) (object.NetworkReference, error) {
	finder := find.NewFinder(client.Client, false)
	if dc != nil {
		finder.SetDatacenter(dc)
	}

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return finder.Network(ctx, name)
}

// FromID loads a network via its managed object reference ID.
func FromID(client *govmomi.Client, id string) (object.NetworkReference, error) {
	// I'm not too sure if a more efficient method to do this exists, but if this
	// becomes a pain point we might want to change this logic a bit.
	//
	// This is pretty much the direct example from
	// github.com/vmware/govmomi/examples/networks/main.go.
	m := view.NewManager(client.Client)

	vctx, vcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer vcancel()
	v, err := m.CreateContainerView(vctx, client.ServiceContent.RootFolder, []string{"Network"}, true)
	if err != nil {
		return nil, err
	}

	defer func() {
		dctx, dcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
		defer dcancel()
		v.Destroy(dctx)
	}()

	var networks []mo.Network
	err = v.Retrieve(vctx, []string{"Network"}, []string{"name"}, &networks)
	if err != nil {
		return nil, err
	}

	for _, net := range networks {
		ref := net.Reference()
		if ref.Value == id {
			finder := find.NewFinder(client.Client, false)
			fctx, fcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
			defer fcancel()
			nref, err := finder.ObjectReference(fctx, ref)
			if err != nil {
				return nil, err
			}
			// Should be safe to return here, as we have already asserted that this type
			// should be a NetworkReference by using ContainerView.
			return nref.(object.NetworkReference), nil
		}
	}
	return nil, fmt.Errorf("could not find network with ID %q", id)
}

// ReferenceProperties is a convenience method that wraps fetching the Network
// MO from a NetworkReference.
//
// Note that regardless of the network type, this only fetches the Network MO
// and not any of the extended properties of that network.
func ReferenceProperties(client *govmomi.Client, net object.NetworkReference) (*mo.Network, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.Network
	nc := object.NewCommon(client.Client, net.Reference())
	if err := nc.Properties(ctx, nc.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}

// Properties gets the properties for a specific Network.
//
// By itself, the Network type usually represents a standard port group in
// vCenter - it has been set up on a host or a set of hosts, and is usually
// configured via through an appropriate HostNetworkSystem. vCenter, however,
// groups up these networks and displays them as a single network that VM can
// use across hosts, facilitating HA and vMotion for VMs that use standard port
// groups versus DVS port groups. Hence the "Network" object is mainly a
// read-only MO and is only useful for checking some very base level
// attributes.
//
// While other network MOs extend the base network object (such as DV port
// groups and opaque networks), this only works with the base object only.
// Refer to functions more specific to the MO to get a fully extended property
// set for the extended objects if you are dealing with those object types.
func Properties(net *object.Network) (*mo.Network, error) {
	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	var props mo.Network
	if err := net.Properties(ctx, net.Reference(), nil, &props); err != nil {
		return nil, err
	}
	return &props, nil
}
