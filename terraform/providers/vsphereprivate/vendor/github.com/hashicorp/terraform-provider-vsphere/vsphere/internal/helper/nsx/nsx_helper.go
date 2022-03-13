package nsx

import (
	"context"
	"fmt"

	"github.com/hashicorp/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
	"github.com/vmware/govmomi/view"
	"github.com/vmware/govmomi/vim25/mo"
	"github.com/vmware/govmomi/vim25/types"
)

// OpaqueNetworkFromNetworkID looks for an opaque network via its opaque network ID.
//
// As NSX support in the Terraform provider is not 100% as of the time of this
// writing (October 2017), this function may require some extra love in order
// to be 100% functional. This is in place to support using NSX with the
// vsphere_virtual_machine resource, as there is no direct path from an opaque
// network backing to the managed object reference that represents the opaque
// network in vCenter.
func OpaqueNetworkFromNetworkID(client *govmomi.Client, id string) (*object.OpaqueNetwork, error) {
	// We use the same ContainerView logic that we use with networkFromID, but we
	// go a step further and limit it to opaque networks only.
	m := view.NewManager(client.Client)

	vctx, vcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer vcancel()
	v, err := m.CreateContainerView(vctx, client.ServiceContent.RootFolder, []string{"OpaqueNetwork"}, true)
	if err != nil {
		return nil, err
	}

	defer func() {
		dctx, dcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
		defer dcancel()
		v.Destroy(dctx)
	}()

	var networks []mo.OpaqueNetwork
	err = v.Retrieve(vctx, []string{"OpaqueNetwork"}, nil, &networks)
	if err != nil {
		return nil, err
	}

	for _, net := range networks {
		if net.Summary.(*types.OpaqueNetworkSummary).OpaqueNetworkId == id {
			ref := net.Reference()
			finder := find.NewFinder(client.Client, false)
			fctx, fcancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
			defer fcancel()
			nref, err := finder.ObjectReference(fctx, ref)
			if err != nil {
				return nil, err
			}
			// Should be safe to return here, as we have already asserted that this type
			// should be a OpaqueNetwork by using ContainerView, along with relying
			// on several fields that only an opaque network would have.
			return nref.(*object.OpaqueNetwork), nil
		}
	}
	return nil, fmt.Errorf("could not find opaque network with ID %q", id)
}
