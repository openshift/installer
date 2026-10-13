package logging

import (
	"fmt"
	"path/filepath"
	"strings"

	"k8s.io/apimachinery/pkg/util/sets"

	"github.com/openshift/installer/pkg/asset"
)

// LogCreatedFiles checks all the asset files created and logs it for the user to see.
func LogCreatedFiles(cmdName string, directory string, targets []asset.WritableAsset) string {
	directory = filepath.Dir(fmt.Sprintf("%s//", directory))

	assetDirs := sets.NewString()
	for _, a := range targets {
		for _, f := range a.Files() {
			index := strings.Index(f.Filename, "/")
			path := directory

			if index != -1 {
				path = filepath.Join(path, f.Filename[:index])
			}
			assetDirs.Insert(path)
		}
	}

	if len(assetDirs) == 0 {
		return ""
	}

	var directories string
	keys := assetDirs.List()

	if len(keys) == 1 {
		directories = keys[0]
	} else {
		maxIndex := 3
		if maxIndex >= len(keys) {
			maxIndex = len(keys)
			directories = strings.Join(keys[:maxIndex-1], ", ")
			directories = fmt.Sprintf("%s and %s", strings.TrimRight(directories, ", "), keys[maxIndex-1])
		} else {
			directories = directory
		}

	}

	return fmt.Sprintf("%s created in: %s", strings.Title(strings.ToLower(cmdName)), directories)
}
