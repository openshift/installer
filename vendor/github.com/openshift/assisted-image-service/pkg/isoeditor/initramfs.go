package isoeditor

import (
	"fmt"
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
)

const initrdPathInISO = "images/pxeboot/initrd.img"

func NewInitRamFSStreamReader(irfsPath string, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	irfsReader, err := os.Open(irfsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open base initrd: %w", err)
	}
	return newInitRamFSStreamReaderFromStream(irfsReader, ignitionContent)
}

func NewInitRamFSStreamReaderFromISO(isoPath string, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	irfsReader, err := GetFileFromISO(isoPath, initrdPathInISO)
	if err != nil {
		return nil, fmt.Errorf("failed to open base initrd from ISO: %w", err)
	}
	return newInitRamFSStreamReaderFromStream(irfsReader, ignitionContent)
}

func newInitRamFSStreamReaderFromStream(irfsReader io.ReadSeekCloser, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	ignitionReader, err := ignitionContent.Archive()
	if err != nil {
		return nil, fmt.Errorf("failed to create ignition archive: %w", err)
	}

	r, err := overlay.NewAppendReader(irfsReader, ignitionReader)
	if err != nil {
		return nil, fmt.Errorf("failed to create append reader for ignition: %w", err)
	}
	return r, nil
}
