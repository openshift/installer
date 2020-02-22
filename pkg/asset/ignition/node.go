package ignition

import (
	"path/filepath"

	igntypes "github.com/coreos/ignition/v2/config/v3_0/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset"
)

// FilesFromAsset creates ignition files for each of the files in the specified
// asset.
func FilesFromAsset(pathPrefix string, username string, mode int, asset asset.WritableAsset) []igntypes.File {
	var files []igntypes.File
	for _, f := range asset.Files() {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, f.Filename), username, mode, f.Data))
	}
	return files
}

// FileFromString creates an ignition-config file with the given contents.
func FileFromString(path string, username string, mode int, contents string) igntypes.File {
	return FileFromBytes(path, username, mode, []byte(contents))
}

// FileFromBytes creates an ignition-config file with the given contents.
func FileFromBytes(path string, username string, mode int, contents []byte) igntypes.File {
	contentsString := dataurl.EncodeBytes(contents)
	overwrite := true
	return igntypes.File{
		Node: igntypes.Node{
			Path:      path,
			Overwrite: &overwrite,
			User: igntypes.NodeUser{
				Name: &username,
			},
		},
		FileEmbedded1: igntypes.FileEmbedded1{
			Mode: &mode,
			Contents: igntypes.FileContents{
				Source: &contentsString,
			},
		},
	}
}
