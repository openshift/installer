// Package release contains assets for the release image (also known
// as the update payload).
package release

import (
	"os"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
)

var (
	defaultImage = "registry.svc.ci.openshift.org/openshift/origin-release:v4.0"
)

// Image is the pull-spec for the release image.
type Image string

var _ asset.Asset = (*Image)(nil)

// Name returns the human-friendly name of the asset.
func (i *Image) Name() string {
	return "Release Image"
}

// Dependencies returns no dependencies.
func (i *Image) Dependencies() []asset.Asset {
	return nil
}

// Generate the release image.
func (i *Image) Generate(p asset.Parents) error {
	releaseImage := defaultImage
	if ri := os.Getenv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE"); ri != "" {
		logrus.Warn("Found override for Image. Please be warned, this is not advised")
		releaseImage = ri
	}
	*i = Image(releaseImage)
	return nil
}
