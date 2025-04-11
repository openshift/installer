package preprovision

import (
	"context"
	"fmt"
	"net/url"
	"os"
	"path/filepath"

	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/imagedata"
	"github.com/gophercloud/gophercloud/v2/openstack/image/v2/images"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/rhcos/cache"
	openstackdefaults "github.com/openshift/installer/pkg/types/openstack/defaults"
)

// UploadBaseImage creates a new image in Glance and uploads the RHCOS image there.
func UploadBaseImage(ctx context.Context, cloud string, rhcosImage string, imageName string, infraID string, imageProperties map[string]string) error {
	url, err := url.Parse(rhcosImage)
	if err != nil {
		return err
	}

	// We support 'http(s)' and 'file' schemes. If the scheme is http(s), then we will upload a file from that
	// location. Otherwise will take local file path from the URL.
	var localFilePath string
	switch url.Scheme {
	case "http", "https":
		localFilePath, err = cache.DownloadImageFile(rhcosImage, cache.InstallerApplicationName)
		if err != nil {
			return err
		}
	case "file":
		localFilePath = filepath.FromSlash(url.Path)
	default:
		return fmt.Errorf("unsupported URL scheme: %q", url.Scheme)
	}

	logrus.Debugln("Creating a Glance image for RHCOS...")

	f, err := os.Open(localFilePath)
	if err != nil {
		return err
	}
	defer f.Close()

	conn, err := openstackdefaults.NewServiceClient(ctx, "image", openstackdefaults.DefaultClientOpts(cloud))
	if err != nil {
		return err
	}

	// By default we use "qcow2" disk format, but if the file extension is "raw",
	// then we set the disk format as "raw".
	diskFormat := "qcow2"
	if extension := filepath.Ext(localFilePath); extension == "raw" {
		diskFormat = "raw"
	}

	img, err := images.Create(ctx, conn, images.CreateOpts{
		Name:            imageName,
		ContainerFormat: "bare",
		DiskFormat:      diskFormat,
		Tags:            []string{"openshiftClusterID=" + infraID},
		Properties:      imageProperties,
	}).Extract()
	if err != nil {
		return err
	}

	// Use direct upload (see
	// https://github.com/openshift/installer/issues/3403 for a discussion
	// on web-download)
	logrus.Debugf("Upload RHCOS to the image %q (%s)", img.Name, img.ID)
	res := imagedata.Upload(ctx, conn, img.ID, f)
	if res.Err != nil {
		return err
	}
	logrus.Debugf("RHCOS image upload completed.")

	return nil
}
