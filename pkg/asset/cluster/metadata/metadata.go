package metadata

import (
	"encoding/json"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/types"
	"github.com/pkg/errors"
)

const (
	FileName = "metadata.json"
)

// LoadMetadata loads the cluster metadata from an asset directory.
func Load(dir string) (*types.ClusterMetadata, error) {
	path := filepath.Join(dir, FileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata *types.ClusterMetadata
	if err = json.Unmarshal(raw, &metadata); err != nil {
		return nil, errors.Wrapf(err, "failed to Unmarshal data from %q to types.ClusterMetadata", path)
	}

	return metadata, err
}
