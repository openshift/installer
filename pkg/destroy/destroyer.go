package destroy

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/types"
)

// Destroyer allows multiple implementations of destroy
// for different platforms.
type Destroyer interface {
	Run() error
}

// NewFunc is an interface for creating platform-specific destroyers.
type NewFunc func(logger logrus.FieldLogger, metadata *types.ClusterMetadata) (Destroyer, error)

// Registry maps ClusterMetadata.Platform() to per-platform Destroyer creators.
var Registry = make(map[string]NewFunc)

// New returns a Destroyer based on `metadata.json` in `rootDir`.
func New(logger logrus.FieldLogger, rootDir string) (Destroyer, error) {
	metadata, err := cluster.LoadMetadata(rootDir)
	if err != nil {
		return nil, err
	}

	platform := metadata.Platform()
	if platform == "" {
		return nil, errors.New("no platform configured in metadata")
	}

	creator, ok := Registry[platform]
	if !ok {
		return nil, errors.Errorf("no destroyers registered for %q", platform)
	}
	return creator(logger, metadata)
}
