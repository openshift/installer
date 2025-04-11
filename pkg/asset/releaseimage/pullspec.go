package releaseimage

import (
	"context"
	"os"

	dockerref "github.com/containers/image/v5/docker/reference"
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
func (a *Image) Generate(_ context.Context, dependencies asset.Parents) error {
	var pullSpec string
	if ri, ok := os.LookupEnv("OPENSHIFT_INSTALL_RELEASE_IMAGE_OVERRIDE"); ok && ri != "" {
		logrus.Warnf("Found override for release image (%s). Please be warned, this is not advised", ri)
		pullSpec = ri
	} else {
		var err error
		pullSpec, err = Default()
		if err != nil {
			return errors.Wrap(err, "failed to load default release image")
		}
		logrus.Debugf("Using internal constant for release image %s", pullSpec)
	}
	a.PullSpec = pullSpec

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
