package isoeditor

import (
	"bytes"
	"encoding/binary"
	"fmt"
	"io"
	"os"

	"github.com/openshift/assisted-image-service/pkg/overlay"
)

const initrdAddrsizePathInISO = "images/initrd.addrsize"

func NewInitrdAddrsizeReader(iafsPath string, initrdFile overlay.OverlayReader) (*bytes.Reader, error) {
	iafsReader, err := os.Open(iafsPath)
	if err != nil {
		return nil, fmt.Errorf("failed to open base initrd.addrsize: %w", err)
	}
	return NewInitrdAddrsizeReaderFromStream(iafsReader, initrdFile)
}

func NewInitrdAddrsizeReaderFromISO(isoPath string, initrdFile overlay.OverlayReader) (*bytes.Reader, error) {
	iafsReader, err := GetFileFromISO(isoPath, initrdAddrsizePathInISO)
	if err != nil {
		return nil, fmt.Errorf("failed to open base initrd.addrsize from ISO: %w", err)
	}
	return NewInitrdAddrsizeReaderFromStream(iafsReader, initrdFile)
}

func NewInitrdAddrsizeReaderFromStream(irfsReader io.ReadSeekCloser, initrdFile overlay.OverlayReader) (*bytes.Reader, error) {
	// get the size of the initrd including the embedded ignition
	sizeOfInitrd, err := initrdFile.Seek(0, io.SeekEnd)
	if err != nil {
		return nil, fmt.Errorf("failed to determine size of initrd: %v", err)
	}

	addrsizeBytes := new(bytes.Buffer)
	err = binary.Write(addrsizeBytes, binary.BigEndian, sizeOfInitrd)
	if err != nil {
		return nil, fmt.Errorf("error during write buffer: %v", err)
	}
	initrdPSW := make([]byte, 8)
	m, err := irfsReader.Read(initrdPSW)
	if err != nil || m != 8 {
		return nil, fmt.Errorf("failed to read initrd.addrsize: %v", err)
	}

	return bytes.NewReader(append(initrdPSW, addrsizeBytes.Bytes()...)), nil
}
