package destroy

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"os"

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
		clusterMetadata.DestroyVolumes = destroyVolumes

		// todo: jcallen - do we need this?
		clusterMetadata.KubeConfig = kubeConfig
		auth, err := getKubeConfigAuth(kubeConfig)
		if err != nil {
			return nil, errors.Wrapf(err, "could not get auth kubeconfig")

		}
		clusterMetadata.Auth = &auth
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

func getKubeConfigAuth(kubeConfig string) ([]byte, error) {
	_, err := os.Stat(kubeConfig)
	if err != nil {
		return nil, err
	}

	auth, err := os.ReadFile(kubeConfig)
	if err != nil {
		return nil, err
	}
	return auth, nil
}
