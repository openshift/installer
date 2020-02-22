package openstack

import (
	"fmt"
	"os"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
)

// uploadBaseImage creates a new image in Glance and uploads the RHCOS image there
func uploadBaseImage(cloud string, localFilePath string, imageName string, clusterID string) error {
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

	logrus.Debugf("Uploading RHCOS to the image %v with ID %v", img.Name, img.ID)
	res := imagedata.Upload(conn, img.ID, f)
	if res.Err != nil {
		return err
	}
	logrus.Debugf("The data was uploaded.")

	return nil
}
