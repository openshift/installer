// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"
	"os"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/azure"
	"github.com/openshift/installer/pkg/types/baremetal"
	"github.com/openshift/installer/pkg/types/gcp"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
	"github.com/openshift/installer/pkg/types/vsphere"
)

// Image is location of RHCOS image.
// This stores the location of the image based on the platform.
// eg. on AWS this contains ami-id, on Livirt this can be the URI for QEMU image etc.
type Image string

var _ asset.Asset = (*Image)(nil)

// Name returns the human-friendly name of the asset.
func (i *Image) Name() string {
	return "Image"
}

// Dependencies returns no dependencies.
func (i *Image) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate the RHCOS image location.
func (i *Image) Generate(p asset.Parents) error {
	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_OS_IMAGE_OVERRIDE"); ok && oi != "" {
		logrus.Warn("Found override for OS Image. Please be warned, this is not advised")
		*i = Image(oi)
		return nil
	}

	ic := &installconfig.InstallConfig{}
	p.Get(ic)
	config := ic.Config
	osimage, err := osImage(config)
	if err != nil {
		return err
	}
	*i = Image(osimage)
	return nil
}

func osImage(config *types.InstallConfig) (string, error) {
	var osimage string
	var err error
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	switch config.Platform.Name() {
	case aws.Name:
		if len(config.Platform.AWS.AMIID) > 0 {
			osimage = config.Platform.AWS.AMIID
			break
		}
		osimage, err = rhcos.AMI(ctx, config.Platform.AWS.Region)
	case gcp.Name:
		osimage, err = rhcos.GCP(ctx)
	case libvirt.Name:
		osimage, err = rhcos.QEMU(ctx)
	case openstack.Name:
		if oi := config.Platform.OpenStack.ClusterOSImage; oi != "" {
			osimage = oi
			break
		}
		osimage, err = rhcos.OpenStack(ctx)
	case azure.Name:
		osimage, err = rhcos.VHD(ctx)
	case baremetal.Name:
		// Check for RHCOS image URL override
		if oi := config.Platform.BareMetal.ClusterOSImage; oi != "" {
			osimage = oi
			break
		}

		// Note that baremetal IPI currently uses the OpenStack image
		// because this contains the necessary ironic config drive
		// ignition support, which isn't enabled in the UPI BM images
		osimage, err = rhcos.OpenStack(ctx)
	case none.Name, vsphere.Name:
	default:
		return "", errors.New("invalid Platform")
	}
	if err != nil {
		return "", err
	}
	return osimage, nil
}
