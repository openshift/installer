package ignition

import (
	"path/filepath"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset"
)

// FilesFromAsset creates ignition files for each of the files in the specified
// asset.
func FilesFromAsset(pathPrefix string, mode int, asset asset.WritableAsset) []ignition.File {
	var files []ignition.File
	for _, f := range asset.Files() {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, f.Filename), mode, f.Data))
	}
	return files
}

// FileFromString creates an ignition-config file with the given contents.
func FileFromString(path string, mode int, contents string) ignition.File {
	return FileFromBytes(path, mode, []byte(contents))
}

// FileFromBytes creates an ignition-config file with the given contents.
func FileFromBytes(path string, mode int, contents []byte) ignition.File {
	return ignition.File{
		Node: ignition.Node{
			Filesystem: "root",
			Path:       path,
		},
		FileEmbedded1: ignition.FileEmbedded1{
			Mode: &mode,
			Contents: ignition.FileContents{
				Source: dataurl.EncodeBytes(contents),
			},
		},
	}
}
