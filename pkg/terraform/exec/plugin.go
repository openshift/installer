//go:build !release
// +build !release

package exec

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func globalPluginDirs(datadir string) ([]string, error) {
	var ret []string
	// Look in ~/.terraform.d/plugins/ , or its equivalent on non-UNIX
	cdir, err := configDir()
	if err != nil {
		return ret, fmt.Errorf("error finding global config directory: %s", err)
	}

	for _, d := range []string{cdir, datadir} {
		machineDir := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
		ret = append(ret, filepath.Join(d, "plugins"))
		ret = append(ret, filepath.Join(d, "plugins", machineDir))
	}

	return ret, nil
}
