package releaseimage

import (
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
       logrus.Debugf("Found override for release image. Using the OKD registry")
       pullSpec := "quay.io/openshift-release-dev/ocp-release@sha256:3a516480dfd68e0f87f702b4d7bdd6f6a0acfdac5cd2e9767b838ceede34d70d"
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
