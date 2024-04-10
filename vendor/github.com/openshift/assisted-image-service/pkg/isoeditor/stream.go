package isoeditor

import (
	"bytes"
	"encoding/json"
	"fmt"
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
	"github.com/pkg/errors"
)

const ignitionImagePath = "/images/ignition.img"
const ignitionInfoPath = "/coreos/igninfo.json"

type ImageReader = overlay.OverlayReader

type BoundariesFinder func(filePath, isoPath string) (int64, int64, error)

type StreamGeneratorFunc func(isoPath string, ignitionContent *IgnitionContent, ramdiskContent, kargs []byte) (ImageReader, error)

type ignitionInfo struct {
	File   string `json:"file,omitempty"`
	Length int64  `json:"length,omitempty"`
	Offset int64  `json:"offset,omitempty"`
}

func NewRHCOSStreamReader(isoPath string, ignitionContent *IgnitionContent, ramdiskContent []byte, kargs []byte) (ImageReader, error) {
	_, r, err := ignitionOverlay(isoPath, ignitionContent, false)
	if err != nil {
		return nil, err
	}

	if ramdiskContent != nil {
		r, err = readerForContent(isoPath, ramDiskImagePath, r, bytes.NewReader(ramdiskContent), GetISOFileInfo)
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

func ignitionOverlay(isoPath string, ignitionContent *IgnitionContent, allowOverflow bool) (*ignitionInfo, overlay.OverlayReader, error) {
	isoReader, err := os.Open(isoPath)
	if err != nil {
		return nil, nil, err
	}

	ignitionReader, err := ignitionContent.Archive()
	if err != nil {
		return nil, nil, err
	}

	ibf := &ignitionBoundaryFinder{
		allowOverflow: allowOverflow,
		dataSize:      ignitionReader.Size(),
	}

	r, err := readerForContent(isoPath, ignitionImagePath, isoReader, ignitionReader, ibf.findBoundaries)
	if err != nil {
		return nil, nil, errors.Wrap(err, "failed to create overwrite reader for ignition")
	}

	if ibf.info.Length > ibf.dataSize {
		offset, _, err := GetISOFileInfo(ibf.info.File, isoPath)
		if err != nil {
			r.Close()
			return nil, nil, err
		}
		paddingLen := ibf.info.Length - ibf.dataSize
		paddingOverlay := overlay.Overlay{
			Reader: bytes.NewReader(bytes.Repeat([]byte{0}, int(paddingLen))),
			Offset: offset + ibf.info.Offset + ibf.dataSize,
			Length: paddingLen,
		}
		if r2, err := overlay.NewOverlayReader(r, paddingOverlay); err == nil {
			r = r2
		} else {
			r.Close()
			return nil, nil, errors.Wrap(err, "failed to create overwrite reader for padding")
		}
	}
	return &ibf.info, r, nil
}

type ignitionBoundaryFinder struct {
	info          ignitionInfo
	allowOverflow bool
	dataSize      int64
}

func (ibf *ignitionBoundaryFinder) findBoundaries(filePath, isoPath string) (int64, int64, error) {
	info := &ibf.info

	ignitionInfoData, err := ReadFileFromISO(isoPath, ignitionInfoPath)
	// If the igninfo.json file doesn't exist or we fail to access it, fall back to using the given ignition file
	// This will be the case for earlier versions of RHCOS
	if err != nil {
		info.File = filePath
	} else {
		err = json.Unmarshal(ignitionInfoData, info)
		if err != nil {
			return 0, 0, err
		}
	}

	isoFileOffset, isoFileLength, err := GetISOFileInfo(info.File, isoPath)
	if err != nil {
		return 0, 0, err
	}

	// use the entire file offset and length if they are not specified in the info struct
	if info.Length == 0 && info.Offset == 0 {
		info.Length = isoFileLength
	}
	// allow overflow if requested and if the embed area extends all the way
	// to the end of the file
	if ibf.allowOverflow && ((info.Offset + info.Length) >= isoFileLength) {
		chunkSize := (info.Length + 3) / 4
		if (ibf.dataSize + chunkSize) > info.Length {
			// increase size in chunks equal to a quarter of the original size,
			// ensuring that there is always at least one full chunk free
			info.Length = (1 + ((ibf.dataSize + chunkSize - 1) / chunkSize)) * chunkSize
		}
	}

	// the final offset is the file offset within the ISO plus the offset within the file
	return isoFileOffset + info.Offset, info.Length, nil
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
