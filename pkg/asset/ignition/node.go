package ignition

import (
	"path/filepath"

	ignition "github.com/coreos/ignition/config/v2_2/types"
	"github.com/vincent-petithory/dataurl"

	"github.com/openshift/installer/pkg/asset"
)

const (
	// keyCertAssetKeyIndex is the index of the private key in a key-pair asset.
	keyCertAssetKeyIndex = 0
	// keyCertAssetCrtIndex is the index of the public key in a key-pair asset.
	keyCertAssetCrtIndex = 1
)

// FilesFromContents creates an ignition-config file with the contents from the
// specified index in the specified asset state.
func FilesFromContents(pathPrefix string, mode int, contents []asset.Content) []ignition.File {
	var files []ignition.File
	for _, c := range contents {
		files = append(files, FileFromBytes(filepath.Join(pathPrefix, c.Name), mode, c.Data))
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
