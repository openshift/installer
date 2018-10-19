package asset

import (
	"fmt"
	"regexp"
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
			name:       "only dirs",
			files:      nil,
			input:      "",
			expectFile: nil,
		},
		{
			name:       "input empty",
			files:      map[string][]byte{"foo.bar": []byte("some data")},
			input:      "",
			expectFile: nil,
		},
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
			f := &fileFetcher{onDiskAssets: tt.files}
			file := f.FetchByName(tt.input)
			assert.Equal(t, tt.expectFile, file)
		})
	}
}

func TestFetchByPattern(t *testing.T) {
	tests := []struct {
		name        string
		files       map[string][]byte
		input       *regexp.Regexp
		expectFiles []*File
	}{
		{
			name: "match master configs",
			files: map[string][]byte{
				"master-0.ign":   []byte("some data 0"),
				"master-1.ign":   []byte("some data 1"),
				"master-2.ign":   []byte("some data 2"),
				"master-10.ign":  []byte("some data 3"),
				"master-20.ign":  []byte("some data 4"),
				"master-00.ign":  []byte("some data 5"),
				"master-01.ign":  []byte("some data 6"),
				"master-0x.ign":  []byte("some data 7"),
				"master-1x.ign":  []byte("some data 8"),
				"amaster-0.ign":  []byte("some data 9"),
				"master-.ign":    []byte("some data 10"),
				"master-.igni":   []byte("some data 11"),
				"master-.ignign": []byte("some data 12"),
			},
			input: regexp.MustCompile(`^(master-(0|([1-9]\d*))\.ign)$`),
			expectFiles: []*File{
				{
					Filename: "master-0.ign",
					Data:     []byte("some data 0"),
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
			name: "match directory",
			files: map[string][]byte{
				"manifests/":     []byte("some data 0"),
				"manifests/0":    []byte("some data 1"),
				"manifests/some": []byte("some data 2"),
				"manifest/":      []byte("some data 3"),
				"manifests":      []byte("some data 4"),
				"amanifests/a":   []byte("some data 5"),
			},
			input: regexp.MustCompile(fmt.Sprintf(`^%s\%c.*`, "manifests", '/')),
			expectFiles: []*File{
				{
					Filename: "manifests/",
					Data:     []byte("some data 0"),
				},
				{
					Filename: "manifests/0",
					Data:     []byte("some data 1"),
				},
				{
					Filename: "manifests/some",
					Data:     []byte("some data 2"),
				},
			},
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			f := &fileFetcher{onDiskAssets: tt.files}
			assert.Equal(t, tt.expectFiles, f.FetchByPattern(tt.input))
		})
	}
}
