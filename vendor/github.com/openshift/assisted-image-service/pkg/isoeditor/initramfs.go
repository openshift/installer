package isoeditor

import (
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
	"github.com/pkg/errors"
)

func NewInitRamFSStreamReader(irfsPath string, ignitionContent *IgnitionContent) (overlay.OverlayReader, error) {
	irfsReader, err := os.Open(irfsPath)
	if err != nil {
		return nil, err
	}

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
