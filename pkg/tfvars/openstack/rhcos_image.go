package openstack

import (
	"fmt"
	"os"
	"time"

	"github.com/gophercloud/gophercloud"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imageimport"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// uploadBaseImage creates a new image in Glance and uploads the RHCOS image there
func uploadBaseImage(cloud string, localFilePath string, imageName string, clusterID string, useImageImport bool) error {
	logrus.Debugln("Creating a Glance image for RHCOS...")

	f, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		return err
	}

	imageCreateOpts := images.CreateOpts{
		Name:            imageName,
		ContainerFormat: "bare",
		DiskFormat:      "qcow2",
		Tags:            []string{fmt.Sprintf("openshiftClusterID=%s", clusterID)},
		// TODO(mfedosin): add Description when gophercloud supports it.
	}

	img, err := images.Create(conn, imageCreateOpts).Extract()
	if err != nil {
		return err
	}
	logrus.Debugf("Image %s was created.", img.Name)

	if useImageImport {
		err = importImage(img.Name, img.ID, conn, f)
		if err != nil {
			err = images.Delete(conn, img.ID).ExtractErr()
			if err != nil {
				return err
			}
			return err
		}
	} else {
		// Use classic legacy upload that doesn't support image conversion
		logrus.Debugf("Using legacy API to upload RHCOS to the image %q with ID %q", img.Name, img.ID)
		res := imagedata.Upload(conn, img.ID, f)
		if res.Err != nil {
			err = images.Delete(conn, img.ID).ExtractErr()
			if err != nil {
				return err
			}
			return err
		}
		logrus.Debugf("The data was uploaded.")
	}

	return nil
}

func importImage(imageName, imageID string, conn *gophercloud.ServiceClient, f *os.File) error {
	logrus.Debugf("Using Image Import API to upload RHCOS to the image %q with ID %q", imageName, imageID)
	stageRes := imagedata.Stage(conn, imageID, f)
	if stageRes.Err != nil {
		logrus.Debugf("Image uploading failed")
		return stageRes.Err
	}
	logrus.Debugf("The data was uploaded")

	logrus.Debugf("Begin image import for the image %q with ID %q", imageName, imageID)
	co := imageimport.CreateOpts{
		Name: imageimport.GlanceDirectMethod,
	}
	importRes := imageimport.Create(conn, imageID, co)
	if importRes.Err != nil {
		return importRes.Err
	}
	logrus.Debugf("Image import started")

	// Image import is an asynchronous operation, so we have to wait until the image becomes "active"
	const numRetries = 5000
	const timeSleepSeconds = 15
	for i := 0; i < numRetries; i++ {
		getRes, err := images.Get(conn, imageID).Extract()
		if err != nil {
			return err
		}

		// More information about Glance Image Status transitioning
		// https://docs.openstack.org/glance/latest/user/statuses.html
		if getRes.Status == images.ImageStatusActive {
			// Import succeed
			break
		} else if getRes.Status == images.ImageStatusQueued || getRes.Status == images.ImageStatusDeleted {
			// Import failed
			logrus.Debugf("RHCOS image import failed")
			return errors.New("RHCOS image import failed")
		}
		time.Sleep(timeSleepSeconds * time.Second)
	}

	logrus.Debugf("Image import finished.")

	return nil
}

// isImageImportSupported checks if we can use Image Import mechanism for image uploading
func isImageImportSupported(cloud string) (bool, error) {
	// More information about the Discovery API:
	// https://docs.openstack.org/api-ref/image/v2/?expanded=#image-service-info-discovery
	logrus.Debugln("Checking if the image import mechanism is supported")

	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		return false, err
	}

	s, err := imageimport.Get(conn).Extract()
	if err != nil {
		// ErrDefault404 means that image discovery API is not available for the cloud
		if _, ok := err.(gophercloud.ErrDefault404); ok {
			return false, nil
		}
		return false, err
	}

	// Next check is just to make sure the response data was not corrupted
	if s.ImportMethods.Type != "array" {
		return false, nil
	}

	for _, method := range s.ImportMethods.Value {
		if method == string(imageimport.GlanceDirectMethod) {
			logrus.Debugln("Glance Direct image import plugin was found")
			return true, nil
		}
	}

	logrus.Debugln("Glance Direct image import plugin was not found")
	return false, nil
}
