package asset

import (
	"io/ioutil"
	"path/filepath"
	"sort"
)

// FileFetcher fetches the asset files from disk.
type FileFetcher interface {
	// FetchByName returns the file with the given name.
	FetchByName(string) (*File, error)
	// FetchByPattern returns the files whose name match the given glob.
	FetchByPattern(pattern string) ([]*File, error)
}

type fileFetcher struct {
	directory string
}

// FetchByName returns the file with the given name.
func (f *fileFetcher) FetchByName(name string) (*File, error) {
	data, err := ioutil.ReadFile(filepath.Join(f.directory, name))
	if err != nil {
		return nil, err
	}
	return &File{Filename: name, Data: data}, nil
}

// FetchByPattern returns the files whose name match the given regexp.
func (f *fileFetcher) FetchByPattern(pattern string) (files []*File, err error) {
	matches, err := filepath.Glob(filepath.Join(f.directory, pattern))
	if err != nil {
		return nil, err
	}

	files = make([]*File, 0, len(matches))
	for _, path := range matches {
		data, err := ioutil.ReadFile(path)
		if err != nil {
			return nil, err
		}

		filename, err := filepath.Rel(f.directory, path)
		if err != nil {
			return nil, err
		}

		files = append(files, &File{
			Filename: filename,
			Data:     data,
		})
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })
	return files, nil
}
