// Package rhcos contains assets for RHCOS.
package rhcos

import (
	"context"
	"fmt"
	"time"

	"github.com/coreos/stream-metadata-go/arch"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/types/azure"
)

// Release is a string which denotes the rhcos release, eg: 412.86.202208101040-0.
// Currently we need this only for Azure to set the image version in the gallery.
// In the future we could extend to other platforms as necessary.
type Release string

var _ asset.Asset = (*Release)(nil)

// Name returns the human-friendly name of the asset.
func (r *Release) Name() string {
	return "Release"
}

// Dependencies returns dependencies used by the RHCOS asset.
func (r *Release) Dependencies() []asset.Asset {
	return []asset.Asset{
		&installconfig.InstallConfig{},
	}
}

// Generate the Release string.
func (r *Release) Generate(ctx context.Context, p asset.Parents) error {
	ic := &installconfig.InstallConfig{}
	p.Get(ic)
	config := ic.Config
	release, err := release(ctx, config)
	if err != nil {
		return err
	}
	*r = Release(release)
	return nil
}

func release(ctx context.Context, config *types.InstallConfig) (string, error) {
	ctx, cancel := context.WithTimeout(ctx, 30*time.Second)
	defer cancel()

	archName := arch.RpmArch(string(config.ControlPlane.Architecture))

	st, err := rhcos.FetchCoreOSBuild(ctx)
	if err != nil {
		return "", err
	}
	streamArch, err := st.GetArchitecture(archName)
	if err != nil {
		return "", err
	}
	switch config.Platform.Name() {
	case azure.Name:
		ext := streamArch.RHELCoreOSExtensions
		if ext == nil {
			return "", fmt.Errorf("%s: No azure build found", st.FormatPrefix(archName))
		}
		azd := ext.AzureDisk
		if azd == nil {
			return "", fmt.Errorf("%s: No azure build found", st.FormatPrefix(archName))
		}
		return azd.Release, nil
	default:
		return "", nil
	}
}

// GetAzureReleaseVersion - generates a modified string for Azure image gallery images. Image gallery image versions cannot have
// a "-" in the name and must be between 0-2147483647, so we have to truncate the hour and minutes of the date.
func (r *Release) GetAzureReleaseVersion() string {
	imageVersion := string(*r)
	if imageVersion != "" {
		imageVersion = imageVersion[:len(imageVersion)-6]
	}
	return imageVersion
}
