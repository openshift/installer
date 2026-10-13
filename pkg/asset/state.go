package asset

import (
	"os"
	"path/filepath"

	"github.com/pkg/errors"
)

// State is the state of an Asset.
type State struct {
	Contents []Content
}

// Content is a generated portion of an Asset.
type Content struct {
	Name string // the path on disk for this content.
	Data []byte
}

// PersistToFile persists the data in the State to files. Each Content entry that
// has a non-empty Name will be persisted to a file with that name.
func (s *State) PersistToFile(directory string) error {
	if s == nil {
		return nil
	}

	for _, c := range s.Contents {
		if c.Name == "" {
			continue
		}
		path := filepath.Join(directory, c.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return errors.Wrap(err, "failed to create dir")
		}
		if err := os.WriteFile(path, c.Data, 0o644); err != nil { //nolint:gosec // no sensitive info
			return errors.Wrap(err, "failed to write file")
		}
	}
	return nil
}
