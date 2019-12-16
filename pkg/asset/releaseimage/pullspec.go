package releaseimage

import (
	"os"

	dockerref "github.com/containers/image/docker/reference"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset"
)

// Image asset generates the release-image pullspec for the cluster
type Image struct {
	PullSpec   string
	Repository string
}

var _ asset.Asset = (*Image)(nil)

// Dependencies is the list of assets required to generate ReleaseImage.
func (a *Image) Dependencies() []asset.Asset {
	return []asset.Asset{}
}

// Generate creates the asset using the dependencies.
func (a *Image) Generate(dependencies asset.Parents) error {
	var pullSpec string
	if ri, ok := os.LookupEnv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE"); ok && ri != "" {
		logrus.Warn("Found override for release image. Please be warned, this is not advised")
		pullSpec = ri
	} else {
		var err error
		pullSpec, err = Default()
		if err != nil {
			return errors.Wrap(err, "failed to load default release image")
		}
		logrus.Debugf("Using internal constant for release image %s", pullSpec)
	}
	a.PullSpec = "registry.svc.ci.openshift.org/rhcos/walters-release@sha256:0961fdf13635a0fbae48e93915982d64ba56c461c477a6fb333b4c97d1e5a082"

	ref, err := dockerref.ParseNamed(pullSpec)
	if err != nil {
		return errors.Wrap(err, "failed to parse release-image pull spec")
	}
	a.Repository = ref.Name()

	return nil
}

// Name is the human friendly name for the asset.
func (a *Image) Name() string {
	return "Release Image Pull Spec"
}
