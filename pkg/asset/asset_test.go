package asset

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

type persistAsset struct{}

func (a *persistAsset) Name() string {
	return "persist-asset"
}

func (a *persistAsset) Dependencies() []Asset {
	return []Asset{}
}

func (a *persistAsset) Generate(context.Context, Parents) error {
	return nil
}

type writablePersistAsset struct {
	persistAsset
	FileList []*File
}

func (a *writablePersistAsset) Files() []*File {
	return a.FileList
}

func (a *writablePersistAsset) Load(FileFetcher) (bool, error) {
	return false, nil
}

func TestPersistToFile(t *testing.T) {
	cases := []struct {
		name      string
		filenames []string
	}{
		{
			name:      "no files",
			filenames: []string{},
		},
		{
			name:      "single file",
			filenames: []string{"file1"},
		},
		{
			name:      "multiple files",
			filenames: []string{"file1", "file2"},
		},
		{
			name:      "new directory",
			filenames: []string{"dir1/file1"},
		},
	}
	for _, tc := range cases {
		t.Run(tc.name, func(t *testing.T) {
			dir := t.TempDir()

			asset := &writablePersistAsset{
				FileList: make([]*File, len(tc.filenames)),
			}
			expectedFiles := map[string][]byte{}
			for i, filename := range tc.filenames {
				data := []byte(fmt.Sprintf("data%d", i))
				asset.FileList[i] = &File{
					Filename: filename,
					Data:     data,
				}
				expectedFiles[filepath.Join(dir, filename)] = data
			}
			err := PersistToFile(asset, dir)
			assert.NoError(t, err, "unexpected error persisting state to file")
			verifyFilesCreated(t, dir, expectedFiles)
		})
	}
}

func verifyFilesCreated(t *testing.T, dir string, expectedFiles map[string][]byte) {
	dirContents, err := os.ReadDir(dir)
	assert.NoError(t, err, "could not read contents of directory %q", dir)
	for _, fileinfo := range dirContents {
		fullPath := filepath.Join(dir, fileinfo.Name())
		if fileinfo.IsDir() {
			verifyFilesCreated(t, fullPath, expectedFiles)
		} else {
			expectedData, fileExpected := expectedFiles[fullPath]
			if !fileExpected {
				t.Errorf("Unexpected file created: %v", fullPath)
				continue
			}
			actualData, err := os.ReadFile(fullPath)
			assert.NoError(t, err, "unexpected error reading created file %q", fullPath)
			assert.Equal(t, expectedData, actualData, "unexpected data in created file %q", fullPath)
			delete(expectedFiles, fullPath)
		}
	}
	for f := range expectedFiles {
		t.Errorf("Expected file %q not created", f)
	}
}
