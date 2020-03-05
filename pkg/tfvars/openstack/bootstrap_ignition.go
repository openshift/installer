package openstack

import (
	"fmt"
	"strings"

	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/imagedata"
	"github.com/gophercloud/gophercloud/openstack/imageservice/v2/images"
	"github.com/gophercloud/utils/openstack/clientconfig"
	"github.com/sirupsen/logrus"
)

// Starting from OpenShift 4.4 we store bootstrap Ignition configs in Glance.

// uploadBootstrapConfig uploads the bootstrap Ignition config in Glance and returns its location
func uploadBootstrapConfig(cloud string, bootstrapIgn string, clusterID string) (string, error) {
	logrus.Debugln("Creating a Glance image for your bootstrap ignition config...")
	opts := clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("image", &opts)
	if err != nil {
		return "", err
	}

	imageCreateOpts := images.CreateOpts{
		Name:            fmt.Sprintf("%s-ignition", clusterID),
		ContainerFormat: "bare",
		DiskFormat:      "raw",
		Tags:            []string{fmt.Sprintf("openshiftClusterID=%s", clusterID)},
		// TODO(mfedosin): add Description when gophercloud supports it.
	}

	img, err := images.Create(conn, imageCreateOpts).Extract()
	if err != nil {
		return "", err
	}
	logrus.Debugf("Image %s was created.", img.Name)

	logrus.Debugf("Uploading bootstrap config to the image %v with ID %v", img.Name, img.ID)
	res := imagedata.Upload(conn, img.ID, strings.NewReader(bootstrapIgn))
	if res.Err != nil {
		return "", res.Err
	}
	logrus.Debugf("The config was uploaded.")

	// img.File contains location of the uploaded data
	return img.File, nil
}

// getAuthToken fetches valid OpenStack authentication token ID
func getAuthToken(cloud string) (string, error) {
	opts := &clientconfig.ClientOpts{
		Cloud: cloud,
	}

	conn, err := clientconfig.NewServiceClient("identity", opts)
	if err != nil {
		return "", err
	}

	token, err := conn.GetAuthResult().ExtractTokenID()
	if err != nil {
		return "", err
	}

	return token, nil
}
