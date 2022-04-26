package isoeditor

import (
	"bytes"
	"compress/gzip"

	"github.com/cavaliercoder/go-cpio"
	"github.com/pkg/errors"
)

type IgnitionContent struct {
	Config []byte
}

func (ic *IgnitionContent) Archive() (*bytes.Reader, error) {
	// Run gzip compression
	compressedBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(compressedBuffer)
	// Create CPIO archive
	cpioWriter := cpio.NewWriter(gzipWriter)

	if err := cpioWriter.WriteHeader(&cpio.Header{
		Name: "config.ign",
		Mode: 0o100_644,
		Size: int64(len(ic.Config)),
	}); err != nil {
		return nil, errors.Wrap(err, "Failed to write CPIO header")
	}
	if _, err := cpioWriter.Write(ic.Config); err != nil {
		return nil, errors.Wrap(err, "Failed to write CPIO archive")
	}

	if err := cpioWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to close CPIO archive")
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to gzip ignition config")
	}

	return bytes.NewReader(compressedBuffer.Bytes()), nil
}
