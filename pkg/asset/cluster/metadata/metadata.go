package metadata

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/types"
)

const (
	// FileName is the filename for the cluster metadata.json file.
	FileName     = "metadata.json"
	AuthFileName = "auth/kubeconfig"
)

// Load loads the cluster metadata from an asset directory.
func Load(dir string) (*types.ClusterMetadata, error) {
	path := filepath.Join(dir, FileName)
	raw, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}

	var metadata *types.ClusterMetadata
	if err = json.Unmarshal(raw, &metadata); err != nil {
		return nil, fmt.Errorf("failed to Unmarshal data from %q to types.ClusterMetadata: %w", path, err)
	}
	authFilePath := filepath.Join(dir, AuthFileName)
	if _, err := os.Stat(authFilePath); err != nil {
		if !os.IsNotExist(err) {
			return nil, err
		}
	} else {
		auth, err := os.ReadFile(authFilePath)
		if err != nil {
			return nil, err
		}
		metadata.Auth = &auth
		return metadata, nil

	}
	return metadata, err
}
