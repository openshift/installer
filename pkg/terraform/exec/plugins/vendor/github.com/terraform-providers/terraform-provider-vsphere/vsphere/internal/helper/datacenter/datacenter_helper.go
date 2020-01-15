package datacenter

import (
	"context"

	"github.com/terraform-providers/terraform-provider-vsphere/vsphere/internal/helper/provider"
	"github.com/vmware/govmomi"
	"github.com/vmware/govmomi/find"
	"github.com/vmware/govmomi/object"
)

// FromPath returns a Datacenter via its supplied path.
func FromPath(client *govmomi.Client, path string) (*object.Datacenter, error) {
	finder := find.NewFinder(client.Client, false)

	ctx, cancel := context.WithTimeout(context.Background(), provider.DefaultAPITimeout)
	defer cancel()
	return finder.Datacenter(ctx, path)
}
