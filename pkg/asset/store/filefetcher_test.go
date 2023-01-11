package store

import (
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"

	"github.com/openshift/installer/pkg/asset"
)

func TestFetchByName(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string][]byte
		input      string
		expectFile *asset.File
	}{
		{
			name:       "input doesn't match",
			files:      map[string][]byte{"foo.bar": []byte("some data")},
			input:      "bar.foo",
			expectFile: nil,
		},
		{
			name:  "with contents",
			files: map[string][]byte{"foo.bar": []byte("some data")},
			input: "foo.bar",
			expectFile: &asset.File{
				Filename: "foo.bar",
				Data:     []byte("some data"),
			},
		},
		{
			name:  "match one file",
			files: map[string][]byte{"foo.bar": []byte("some data")},
			input: "foo.bar",
			expectFile: &asset.File{
				Filename: "foo.bar",
				Data:     []byte("some data"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir := t.TempDir()

			for filename, data := range tt.files {
				err := os.WriteFile(filepath.Join(tempDir, filename), data, 0o666) //nolint:gosec // no sensitive data
				if err != nil {
					t.Fatal(err)
				}
			}

			f := &fileFetcher{directory: tempDir}
			file, err := f.FetchByName(tt.input)
			if err != nil {
				if os.IsNotExist(err) && tt.expectFile == nil {
					return
				}
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectFile, file)
		})
	}
}

func TestFetchByPattern(t *testing.T) {
	tempDir := t.TempDir()

	files := map[string][]byte{
		"master-0.ign":   []byte("some data 0"),
		"master-1.ign":   []byte("some data 1"),
		"master-2.ign":   []byte("some data 2"),
		"master-10.ign":  []byte("some data 3"),
		"master-20.ign":  []byte("some data 4"),
		"master-00.ign":  []byte("some data 5"),
		"master-01.ign":  []byte("some data 6"),
		"amaster-0.ign":  []byte("some data 7"),
		"master-.ign":    []byte("some data 8"),
		"master-.igni":   []byte("some data 9"),
		"master-.ignign": []byte("some data 10"),
		"manifests/0":    []byte("some data 11"),
		"manifests/some": []byte("some data 12"),
		"amanifests/a":   []byte("some data 13"),
	}

	for path, data := range files {
		dir := filepath.Dir(path)
		if dir != "." {
			err := os.MkdirAll(filepath.Join(tempDir, dir), 0777)
			if err != nil {
				t.Fatal(err)
			}
		}
		err := os.WriteFile(filepath.Join(tempDir, path), data, 0o666) //nolint:gosec // no sensitive data
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		input       string
		expectFiles []*asset.File
	}{
		{
			input: "master-[0-9]*.ign",
			expectFiles: []*asset.File{
				{
					Filename: "master-0.ign",
					Data:     []byte("some data 0"),
				},
				{
					Filename: "master-00.ign",
					Data:     []byte("some data 5"),
				},
				{
					Filename: "master-01.ign",
					Data:     []byte("some data 6"),
				},
				{
					Filename: "master-1.ign",
					Data:     []byte("some data 1"),
				},
				{
					Filename: "master-10.ign",
					Data:     []byte("some data 3"),
				},
				{
					Filename: "master-2.ign",
					Data:     []byte("some data 2"),
				},
				{
					Filename: "master-20.ign",
					Data:     []byte("some data 4"),
				},
			},
		},
		{
			input: filepath.Join("manifests", "*"),
			expectFiles: []*asset.File{
				{
					Filename: "manifests/0",
					Data:     []byte("some data 11"),
				},
				{
					Filename: "manifests/some",
					Data:     []byte("some data 12"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.input, func(t *testing.T) {
			f := &fileFetcher{directory: tempDir}
			files, err := f.FetchByPattern(tt.input)
			if err != nil {
				t.Fatal(err)
			}

			assert.Equal(t, tt.expectFiles, files)
		})
	}
}
