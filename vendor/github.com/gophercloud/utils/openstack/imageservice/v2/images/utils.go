package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// IDFromName is a convienience function that returns an image's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	count := 0
	id := ""

	listOpts := images.ListOpts{
		Name: name,
	}

	pages, err := images.List(client, listOpts).AllPages()
	if err != nil {
		return "", err
	}

	all, err := images.ExtractImages(pages)
	if err != nil {
		return "", err
	}

	for _, s := range all {
		if s.Name == name {
			count++
			id = s.ID
		}
	}

	switch count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{Name: name, ResourceType: "image"}
	case 1:
		return id, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{Name: name, Count: count, ResourceType: "image"}
	}
}
