package destroy

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	"github.com/openshift/installer/pkg/destroy/providers"
)

// New returns a Destroyer based on `metadata.json` in `rootDir`.
func New(logger logrus.FieldLogger, rootDir string, deleteVolumes bool) (providers.Destroyer, error) {
	clusterMetadata, err := metadata.Load(rootDir)
	if err != nil {
		return nil, err
	}

	// todo: jcallen: need to think if this makes sense because we could
	// todo: still remove cns volumes for the vSphere case
	if len(*clusterMetadata.Auth) > 0 {
		clusterMetadata.DeleteVolumes = deleteVolumes
	} else {
		clusterMetadata.DeleteVolumes = false
	}

	platform := clusterMetadata.Platform()
	if platform == "" {
		return nil, errors.New("no platform configured in metadata")
	}

	creator, ok := providers.Registry[platform]
	if !ok {
		return nil, errors.Errorf("no destroyers registered for %q", platform)
	}
	return creator(logger, clusterMetadata)
}
