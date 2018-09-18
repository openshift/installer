package asset

import (
	"fmt"
	"path/filepath"
)

// Asset used to install OpenShift.
type Asset interface {
	// Dependencies returns the assets upon which this asset directly depends.
	Dependencies() []Asset

	// Generate generates this asset given the states of its dependent assets.
	Generate(map[Asset]*State) (*State, error)

	// Name returns the human-friendly name of the asset.
	Name() string
}

// GetDataByFilename searches the file in the asset.State.Contents, and returns its data.
// filename is the base name of the file.
func GetDataByFilename(a Asset, parents map[Asset]*State, filename string) ([]byte, error) {
	st, ok := parents[a]
	if !ok {
		return nil, fmt.Errorf("failed to find %T in parents", a)
	}

	for _, c := range st.Contents {
		if filepath.Base(c.Name) == filename {
			return c.Data, nil
		}
	}
	return nil, fmt.Errorf("failed to find data in %v with filename == %q", st, filename)
}
