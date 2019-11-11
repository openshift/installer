package openstack

import (
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/utils/openstack/clientconfig"
)

// DeleteSwiftContainer deletes a container and all of its objects.
func DeleteSwiftContainer(name string, cloud string) error {
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	swiftClient, err := clientconfig.NewServiceClient("object-store", &opts)
	if err != nil {
		return err
	}

	listOpts := objects.ListOpts{
		Full: false,
	}

	allPages, err := objects.List(swiftClient, name, listOpts).AllPages()
	if err != nil {
		return err
	}

	allObjects, err := objects.ExtractNames(allPages)
	if err != nil {
		return err
	}

	for _, object := range allObjects {
		_, err := objects.Delete(swiftClient, name, object, objects.DeleteOpts{}).Extract()
		if err != nil {
			return err
		}
	}

	_, err = containers.Delete(swiftClient, name).Extract()
	return err
}
