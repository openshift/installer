package openstack

import (
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

// DeleteGlanceImage deletes the image with the specified name
func DeleteGlanceImage(name string, cloud string) error {
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		return err
	}

	listOpts := images.ListOpts{
		Name: name,
	}

	allPages, err := images.List(conn, listOpts).AllPages()
	if err != nil {
		return err
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		return err
	}

	for _, image := range allImages {
		err := images.Delete(conn, image.ID).ExtractErr()
		if err != nil {
			return err
		}
	}
	return nil
}
