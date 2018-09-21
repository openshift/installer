package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
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
	for _, c := range s.Contents {
		if c.Name == "" {
			continue
		}
		path := filepath.Join(directory, c.Name)
		if err := os.MkdirAll(filepath.Dir(path), 0755); err != nil {
			return err
		}
		if err := ioutil.WriteFile(path, c.Data, 0644); err != nil {
			return err
		}
	}
	return nil
}
