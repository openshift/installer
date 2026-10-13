package store

import (
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
)

type fileFetcher struct {
	directory string
}

// FetchByName returns the file with the given name.
func (f *fileFetcher) FetchByName(name string) (*asset.File, error) {
	data, err := os.ReadFile(filepath.Join(f.directory, name))
	if err != nil {
		return nil, err
	}
	return &asset.File{Filename: name, Data: data}, nil
}

// FetchByPattern returns the files whose name match the given regexp.
func (f *fileFetcher) FetchByPattern(pattern string) (files []*asset.File, err error) {
	matches, err := filepath.Glob(filepath.Join(f.directory, pattern))
	if err != nil {
		return nil, err
	}

	files = make([]*asset.File, 0, len(matches))
	for _, path := range matches {
		data, err := os.ReadFile(path)
		if err != nil {
			return nil, err
		}

		filename, err := filepath.Rel(f.directory, path)
		if err != nil {
			return nil, err
		}

		files = append(files, &asset.File{
			Filename: filename,
			Data:     data,
		})
	}

	return files, nil
}
