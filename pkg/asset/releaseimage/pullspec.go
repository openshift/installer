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
	pullSpec := "quay.io/openshift-release-dev/ocp-release@sha256:64320fbf95d968fc6b9863581a92d373bc75f563a13ae1c727af37450579f61a"
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
