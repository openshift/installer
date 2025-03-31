// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
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
		new(Image),
	}
}

// Generate the RHCOS Bootstrap image location.
func (i *BootstrapImage) Generate(ctx context.Context, p asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	rhcosImage := new(Image)
	p.Get(ic, rhcosImage)
	config := ic.Config

	switch config.Platform.Name() {
	case baremetal.Name:
		// Check for CoreOS image URL override
		if boi := config.Platform.BareMetal.BootstrapOSImage; boi != "" {
			*i = BootstrapImage(boi)
		}
		return nil
	default:
		// other platforms use the same image for all nodes
		*i = BootstrapImage(rhcosImage.ControlPlane)
		return nil
	}
}
