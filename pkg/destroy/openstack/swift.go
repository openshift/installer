package openstack

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/containers"
	"github.com/gophercloud/gophercloud/openstack/objectstorage/v1/objects"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
)

// DeleteSwiftContainer deletes a container and all of its objects.
func DeleteSwiftContainer(name string, cloud string) error {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 20,
		Steps:    30,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return deleteSwiftContainer(name, cloud)
	})
	if err != nil {
		return errors.Errorf("Unrecoverable error/timed out: %v", err)
	}

	return nil
}

func deleteSwiftContainer(name string, cloud string) (bool, error) {
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	swiftClient, err := clientconfig.NewServiceClient("object-store", &opts)
	if err != nil {
		logrus.Debugf("There was an error during the container removal: %v", err)
		return false, nil
	}

	listOpts := objects.ListOpts{
		Full: false,
	}

	allPages, err := objects.List(swiftClient, name, listOpts).AllPages()
	if err != nil {
		logrus.Debugf("There was an error during the container removal: %v", err)
		return false, nil
	}

	allObjects, err := objects.ExtractNames(allPages)
	if err != nil {
		logrus.Debugf("There was an error during the container removal: %v", err)
		return false, nil
	}

	for _, object := range allObjects {
		_, err := objects.Delete(swiftClient, name, object, objects.DeleteOpts{}).Extract()
		if err != nil {
			logrus.Debugf("There was an error during the container removal: %v", err)
			return false, nil
		}
	}

	_, err = containers.Delete(swiftClient, name).Extract()
	if err != nil {
		logrus.Debugf("There was an error during the container removal: %v", err)
		return false, nil
	}

	return true, nil
}
