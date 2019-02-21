// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"
	"time"

	"github.com/pkg/errors"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	releaseasset "github.com/openshift/installer/pkg/asset/release"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/rhcos/release"
	"github.com/openshift/installer/pkg/types/aws"
	"github.com/openshift/installer/pkg/types/libvirt"
	"github.com/openshift/installer/pkg/types/none"
	"github.com/openshift/installer/pkg/types/openstack"
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

// Dependencies returns the assets on which the Image asset depends.
func (i *Image) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
		new(releaseasset.Image),
	}
}

// Generate the RHCOS image location.
func (i *Image) Generate(p asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	releaseImage := new(releaseasset.Image)
	p.Get(ic, releaseImage)
	config := ic.Config

	ctx := context.TODO()
	build, err := release.RHCOSBuild(ctx, string(*releaseImage), []byte(ic.Config.PullSecret))
	if err != nil {
		return err
	}

	var osimage string
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()
	switch config.Platform.Name() {
	case aws.Name:
		osimage, err = rhcos.AMI(ctx, rhcos.DefaultChannel, build, config.Platform.AWS.Region)
	case libvirt.Name:
		osimage, err = rhcos.QEMU(ctx, rhcos.DefaultChannel, build)
	case openstack.Name:
		osimage = "rhcos"
	case none.Name:
	default:
		return errors.New("invalid Platform")
	}
	if err != nil {
		return err
	}
	*i = Image(osimage)
	return nil
}
