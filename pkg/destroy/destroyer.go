package destroy

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path/filepath"

	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/asset/metadata"
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
	path := filepath.Join(rootDir, metadata.MetadataFilename)
	raw, err := ioutil.ReadFile(filepath.Join(rootDir, metadata.MetadataFilename))
	if err != nil {
		return nil, err
	}

	var cmetadata *types.ClusterMetadata
	if err := json.Unmarshal(raw, &cmetadata); err != nil {
		return nil, err
	}

	platform := cmetadata.Platform()
	if platform == "" {
		return nil, fmt.Errorf("no platform configured in %q", path)
	}

	creator, ok := Registry[platform]
	if !ok {
		return nil, fmt.Errorf("no destroyers registered for %q", platform)
	}
	return creator(logger, cmetadata)
}
