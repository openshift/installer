package store

import (
	"errors"
	"io/fs"
	"os"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset"
)

type fileFetcher struct {
	directory string
	modCheck  modificationChecker
}

// FetchByName returns the file with the given name.
func (f *fileFetcher) FetchByName(name string) (*asset.File, error) {
	path := filepath.Join(f.directory, name)
	data, err := os.ReadFile(path)
	if err != nil {
		return nil, err
	}
	file := &asset.File{
		Filename: name,
		Data:     data,
	}
	if err := f.modCheck.CheckModified(path, file); err != nil {
		return nil, err
	}
	return file, nil
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

		file := &asset.File{
			Filename: filename,
			Data:     data,
		}

		if err := f.modCheck.CheckModified(path, file); err != nil {
			if errors.Is(err, fs.ErrNotExist) {
				continue
			}
			return nil, err
		}

		files = append(files, file)
	}

	return files, nil
}
