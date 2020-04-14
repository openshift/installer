package openstack

import (
	"time"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"k8s.io/apimachinery/pkg/util/wait"
)

// DeleteGlanceImage deletes the image with the specified name
func DeleteGlanceImage(name string, cloud string) error {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 20,
		Steps:    30,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return deleteGlanceImage(name, cloud)
	})
	if err != nil {
		return errors.Errorf("Unrecoverable error/timed out: %v", err)
	}

	return nil
}

func deleteGlanceImage(name string, cloud string) (bool, error) {
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		logrus.Warningf("There was an error during the image removal: %v", err)
		return false, nil
	}

	listOpts := images.ListOpts{
		Name: name,
	}

	allPages, err := images.List(conn, listOpts).AllPages()
	if err != nil {
		logrus.Warningf("There was an error during the image removal: %v", err)
		return false, nil
	}

	allImages, err := images.ExtractImages(allPages)
	if err != nil {
		logrus.Warningf("There was an error during the image removal: %v", err)
		return false, nil
	}

	for _, image := range allImages {
		err := images.Delete(conn, image.ID).ExtractErr()
		if err != nil {
			logrus.Warningf("There was an error during the image removal: %v", err)
			return false, nil
		}
	}
	return true, nil
}
