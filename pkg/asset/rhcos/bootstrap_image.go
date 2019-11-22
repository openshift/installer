// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"
	"os"
	"time"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types/baremetal"
)

// BootstrapImage is location of the RHCOS image for the Bootstrap node
// This stores the location of the image based on the platform.
// eg. on AWS this contains ami-id, on Livirt this can be the URI for QEMU image etc.
// Note that for most platforms this is the same as rhcos.Image
type BootstrapImage string

var _ asset.Asset = (*BootstrapImage)(nil)

// Name returns the human-friendly name of the asset.
func (i *BootstrapImage) Name() string {
	return "BootstrapImage"
}

// Dependencies returns no dependencies.
func (i *BootstrapImage) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate the RHCOS Bootstrap image location.
func (i *BootstrapImage) Generate(p asset.Parents) error {
	if oi, ok := os.LookupEnv("OPENSHIFT_INSTALL_BOOTSTRAP_OS_IMAGE_OVERRIDE"); ok && oi != "" {
		logrus.Warn("Found override for Bootstrap OS Image. Please be warned, this is not advised")
		*i = BootstrapImage(oi)
		return nil
	}

	ic := &installconfig.InstallConfig{}
	p.Get(ic)
	config := ic.Config

	var osimage string
	var err error
	ctx, cancel := context.WithTimeout(context.TODO(), 30*time.Second)
	defer cancel()
	switch config.Platform.Name() {
	case baremetal.Name:
		// Baremetal IPI launches a local VM for the bootstrap node
		// Hence requires the QEMU image to use the libvirt backend
		osimage, err = rhcos.QEMU(ctx)
	default:
		// other platforms use the same image for all nodes
		osimage, err = osImage(config)
	}
	if err != nil {
		return err
	}
	*i = BootstrapImage(osimage)
	return nil
}
