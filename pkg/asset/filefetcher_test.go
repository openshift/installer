package asset

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"testing"

	"github.com/stretchr/testify/assert"
)

func TestFetchByName(t *testing.T) {
	tests := []struct {
		name       string
		files      map[string][]byte
		input      string
		expectFile *File
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
			expectFile: &File{
				Filename: "foo.bar",
				Data:     []byte("some data"),
			},
		},
		{
			name:  "match one file",
			files: map[string][]byte{"foo.bar": []byte("some data")},
			input: "foo.bar",
			expectFile: &File{
				Filename: "foo.bar",
				Data:     []byte("some data"),
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			tempDir, err := ioutil.TempDir("", "openshift-install-")
			if err != nil {
				t.Fatal(err)
			}
			defer os.RemoveAll(tempDir)

			for filename, data := range tt.files {
				err = ioutil.WriteFile(filepath.Join(tempDir, filename), data, 0666)
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
	tempDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		t.Fatal(err)
	}
	defer os.RemoveAll(tempDir)

	files := map[string][]byte{
		"controlplane-0.ign":   []byte("some data 0"),
		"controlplane-1.ign":   []byte("some data 1"),
		"controlplane-2.ign":   []byte("some data 2"),
		"controlplane-10.ign":  []byte("some data 3"),
		"controlplane-20.ign":  []byte("some data 4"),
		"controlplane-00.ign":  []byte("some data 5"),
		"controlplane-01.ign":  []byte("some data 6"),
		"acontrolplane-0.ign":  []byte("some data 7"),
		"controlplane-.ign":    []byte("some data 8"),
		"controlplane-.igni":   []byte("some data 9"),
		"controlplane-.ignign": []byte("some data 10"),
		"manifests/0":          []byte("some data 11"),
		"manifests/some":       []byte("some data 12"),
		"amanifests/a":         []byte("some data 13"),
	}

	for path, data := range files {
		dir := filepath.Dir(path)
		if dir != "." {
			err := os.MkdirAll(filepath.Join(tempDir, dir), 0777)
			if err != nil {
				t.Fatal(err)
			}
		}
		err = ioutil.WriteFile(filepath.Join(tempDir, path), data, 0666)
		if err != nil {
			t.Fatal(err)
		}
	}
	tests := []struct {
		input       string
		expectFiles []*File
	}{
		{
			input: "controlplane-[0-9]*.ign",
			expectFiles: []*File{
				{
					Filename: "controlplane-0.ign",
					Data:     []byte("some data 0"),
				},
				{
					Filename: "controlplane-00.ign",
					Data:     []byte("some data 5"),
				},
				{
					Filename: "controlplane-01.ign",
					Data:     []byte("some data 6"),
				},
				{
					Filename: "controlplane-1.ign",
					Data:     []byte("some data 1"),
				},
				{
					Filename: "controlplane-10.ign",
					Data:     []byte("some data 3"),
				},
				{
					Filename: "controlplane-2.ign",
					Data:     []byte("some data 2"),
				},
				{
					Filename: "controlplane-20.ign",
					Data:     []byte("some data 4"),
				},
			},
		},
		{
			input: filepath.Join("manifests", "*"),
			expectFiles: []*File{
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
