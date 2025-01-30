package destroy

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/cluster/metadata"
	"github.com/openshift/installer/pkg/destroy/providers"
)

// New returns a Destroyer based on `metadata.json` in `rootDir`.
func New(logger logrus.FieldLogger, rootDir string, destroyVolumes bool, kubeConfig string) (providers.Destroyer, error) {
	clusterMetadata, err := metadata.Load(rootDir)
	if err != nil {
		return nil, err
	}

	if destroyVolumes {
		ctx := context.Background()

		v, err := NewVolume(ctx, logger, kubeConfig)
		if err != nil {
			return nil, err
		}

		// if there are PVs continue with destroy processes
		if v.persistentVolumeList != nil && len(v.persistentVolumeClaimList.Items) != 0 {
			if err := v.drainNodes(ctx); err != nil {
				return nil, err
			}

			if err := v.deletePersistentVolumes(ctx); err != nil {
				return nil, err
			}
		}
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
