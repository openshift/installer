package isoeditor

import (
	"bytes"
	"fmt"
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
	"github.com/pkg/errors"
)

const ignitionImagePath = "/images/ignition.img"

type ImageReader = overlay.OverlayReader

type BoundariesFinder func(filePath, isoPath string) (int64, int64, error)

type StreamGeneratorFunc func(isoPath string, ignitionContent *IgnitionContent, ramdiskContent, kargs []byte) (ImageReader, error)

func NewRHCOSStreamReader(isoPath string, ignitionContent *IgnitionContent, ramdiskContent []byte, kargs []byte) (ImageReader, error) {
	isoReader, err := os.Open(isoPath)
	if err != nil {
		return nil, err
	}

	ignitionReader, err := ignitionContent.Archive()
	if err != nil {
		return nil, err
	}

	r, err := readerForFileContent(isoPath, ignitionImagePath, isoReader, ignitionReader)
	if err != nil {
		return nil, errors.Wrap(err, "failed to create overwrite reader for ignition")
	}

	if ramdiskContent != nil {
		r, err = readerForFileContent(isoPath, ramDiskImagePath, r, bytes.NewReader(ramdiskContent))
		if err != nil {
			return nil, errors.Wrap(err, "failed to create overwrite reader for ramdisk")
		}
	}

	if kargs != nil {
		files, err := KargsFiles(isoPath)
		if err != nil {
			return nil, errors.Wrap(err, "failed to read files to patch for kernel arguments")
		}
		for _, file := range files {
			r, err = readerForKargsContent(isoPath, file, r, bytes.NewReader(kargs))
			if err != nil {
				return nil, errors.Wrapf(err, "failed to create overwrite reader for kernel arguments in file \"%s\"", file)
			}
		}
	}

	return r, nil
}

func readerForContent(isoPath, filePath string, base io.ReadSeeker, contentReader *bytes.Reader, boundariesFinder BoundariesFinder) (overlay.OverlayReader, error) {
	start, length, err := boundariesFinder(filePath, isoPath)
	if err != nil {
		return nil, err
	}

	if length < contentReader.Size() {
		return nil, errors.New(fmt.Sprintf("content length (%d) exceeds embed area size (%d)", contentReader.Size(), length))
	}

	rdOverlay := overlay.Overlay{
		Reader: contentReader,
		Offset: start,
		Length: contentReader.Size(),
	}
	r, err := overlay.NewOverlayReader(base, rdOverlay)
	if err != nil {
		return nil, err
	}

	return r, nil
}

func readerForFileContent(isoPath string, filePath string, base io.ReadSeeker, contentReader *bytes.Reader) (overlay.OverlayReader, error) {
	return readerForContent(isoPath, filePath, base, contentReader, GetISOFileInfo)
}
