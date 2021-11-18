//go:build release
// +build release

package exec

import (
	"fmt"
	"path/filepath"
	"runtime"
)

func globalPluginDirs(datadir string) ([]string, error) {
	var ret []string
	for _, d := range []string{datadir} {
		machineDir := fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH)
		ret = append(ret, filepath.Join(d, "plugins"))
		ret = append(ret, filepath.Join(d, "plugins", machineDir))
	}
	return ret, nil
}
