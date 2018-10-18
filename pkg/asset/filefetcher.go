package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"regexp"
	"sort"
)

// FileFetcher fetches the asset files from disk.
type FileFetcher interface {
	// FetchByName returns the file with the given name.
	FetchByName(string) *File
	// FetchByPattern returns the files whose name match the given regexp.
	FetchByPattern(*regexp.Regexp) []*File
}

type fileFetcher struct {
	onDiskAssets map[string][]byte
}

func newFileFetcher(clusterDir string) (*fileFetcher, error) {
	fileMap := make(map[string][]byte)

	// Don't bother if the clusterDir is not created yet because that
	// means there's no assets generated yet.
	_, err := os.Stat(clusterDir)
	if err != nil && os.IsNotExist(err) {
		return &fileFetcher{}, nil
	}

	if err := filepath.Walk(clusterDir, func(path string, info os.FileInfo, err error) error {
		if err != nil {
			return err
		}

		if info.IsDir() {
			return nil
		}

		filename, err := filepath.Rel(clusterDir, path)
		if err != nil {
			return err
		}

		data, err := ioutil.ReadFile(path)
		if err != nil {
			return err
		}

		fileMap[filename] = data
		return nil
	}); err != nil {
		return nil, err
	}
	return &fileFetcher{onDiskAssets: fileMap}, nil
}

// FetchByName returns the file with the given name.
func (f *fileFetcher) FetchByName(name string) *File {
	data, ok := f.onDiskAssets[name]
	if !ok {
		return nil
	}
	return &File{Filename: name, Data: data}
}

// FetchByPattern returns the files whose name match the given regexp.
func (f *fileFetcher) FetchByPattern(re *regexp.Regexp) []*File {
	var files []*File

	for filename, data := range f.onDiskAssets {
		if re.MatchString(filename) {
			files = append(files, &File{
				Filename: filename,
				Data:     data,
			})
		}
	}

	sort.Slice(files, func(i, j int) bool { return files[i].Filename < files[j].Filename })
	return files
}
