package images

import (
	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
)

// IDFromName is a convienience function that returns an image's ID given its name.
func IDFromName(client *gophercloud.ServiceClient, name string) (string, error) {
	allPages, err := images.List(client, images.ListOpts{
		Name: name,
	}).AllPages()
	if err != nil {
		return "", err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return "", err
	}

	switch count := len(allImages); count {
	case 0:
		return "", gophercloud.ErrResourceNotFound{
			Name:         name,
			ResourceType: "image",
		}
	case 1:
		return allImages[0].ID, nil
	default:
		return "", gophercloud.ErrMultipleResourcesFound{
			Name:         name,
			Count:        count,
			ResourceType: "image",
		}
	}
}
