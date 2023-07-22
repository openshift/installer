package isoeditor

import (
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
	"github.com/pkg/errors"
)

const initrdPathInISO = "images/pxeboot/initrd.img"

func NewInitRamFSStreamReader(irfsPath string, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	irfsReader, err := os.Open(irfsPath)
	if err != nil {
		return nil, err
	}
	return newInitRamFSStreamReaderFromStream(irfsReader, ignitionContent)
}

func NewInitRamFSStreamReaderFromISO(isoPath string, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	irfsReader, err := GetFileFromISO(isoPath, initrdPathInISO)
	if err != nil {
		return nil, err
	}
	return newInitRamFSStreamReaderFromStream(irfsReader, ignitionContent)
}

func newInitRamFSStreamReaderFromStream(irfsReader io.ReadSeekCloser, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	ignitionReader, err := ignitionContent.Archive()
	if err != nil {
		return nil, err
	}

	r, err := overlay.NewAppendReader(irfsReader, ignitionReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create append reader for ignition")
	}
	return r, nil
}
