package flavors

import (
	"context"

	"github.com/gophercloud/gophercloud/v2"
	"github.com/gophercloud/gophercloud/v2/openstack/compute/v2/flavors"
)

// IDFromName is a convenience function that returns a flavor's ID given its
// name. Errors when the number of items found is not one.
func IDFromName(ctx context.Context, client *gophercloud.ServiceClient, name string) (string, error) {
	IDs, err := IDsFromName(ctx, client, name)
	if err != nil {
		return "", err
	}

	switch count := len(IDs); count {
	case 0:
		return "", &gophercloud.ErrResourceNotFound{Name: name, ResourceType: "flavor"}
	case 1:
		return IDs[0], nil
	default:
		return "", &gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "flavor"}
	}
}

// IDsFromName returns zero or more IDs corresponding to a name. The returned
// error is only non-nil in case of failure.
func IDsFromName(ctx context.Context, client *gophercloud.ServiceClient, name string) ([]string, error) {
	pages, err := flavors.ListDetail(client, nil).AllPages(ctx)
	if err != nil {
		return nil, err
	}

	all, err := flavors.ExtractFlavors(pages)
	if err != nil {
		return nil, err
	}

	IDs := make([]string, 0, len(all))
	for _, s := range all {
		if s.Name == name {
			IDs = append(IDs, s.ID)
		}
	}

	return IDs, nil
}
