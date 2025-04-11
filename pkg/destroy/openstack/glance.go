package openstack

import (
	"context"
	"fmt"
	"time"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/sirupsen/logrus"
	"k8s.io/apimachinery/pkg/util/wait"

	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// DeleteGlanceImage deletes the image with the specified name
func DeleteGlanceImage(ctx context.Context, name string, cloud string) error {
	backoffSettings := wait.Backoff{
		Duration: time.Second * 20,
		Steps:    30,
	}

	err := wait.ExponentialBackoff(backoffSettings, func() (bool, error) {
		return deleteGlanceImage(ctx, name, cloud)
	})
	if err != nil {
		return fmt.Errorf("unrecoverable error/timed out: %w", err)
	}

	return nil
}

func deleteGlanceImage(ctx context.Context, name string, cloud string) (bool, error) {
	conn, err := openstackdefaults.NewServiceClient(ctx, "image", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		logrus.Warningf("There was an error during the image removal: %v", err)
		return false, nil
	}

	listOpts := images.ListOpts{
		Name: name,
	}

	allPages, err := images.List(conn, listOpts).AllPages(ctx)
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
		err := images.Delete(ctx, conn, image.ID).ExtractErr()
		if err != nil {
			logrus.Warningf("There was an error during the image removal: %v", err)
			return false, nil
		}
	}
	return true, nil
}
