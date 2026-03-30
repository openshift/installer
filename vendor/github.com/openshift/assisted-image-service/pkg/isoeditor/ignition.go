package isoeditor

import (
	"bytes"
	"fmt"
	"path"
	"strings"
)

type IgnitionContent struct {
	Config        []byte
	SystemConfigs map[string][]byte
}

func (ic *IgnitionContent) Archive() (*bytes.Reader, error) {
	var files []fileEntry

	if len(ic.Config) > 0 {
		files = append(files, fileEntry{
			Content: ic.Config,
			Path:    "config.ign",
			Mode:    0o100_644,
		})
	}

	for filename, content := range ic.SystemConfigs {
		if strings.Contains(filename, "/") {
			return nil, fmt.Errorf("system config filename %q contains path separators", filename)
		}
		files = append(files, fileEntry{
			Content: content,
			Path:    path.Join("usr", "lib", "ignition", "base.d", filename),
			Mode:    0o100_644,
		})
	}

	compressedCpio, err := generateCompressedCPIO(files)
	if err != nil {
		return nil, err
	}
	return bytes.NewReader(compressedCpio), nil
}
