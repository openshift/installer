package image

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"os"
	"path/filepath"

	"github.com/cavaliercoder/go-cpio"
)

// CpioArchive simplifies the creation of a compressed cpio archive.
type CpioArchive struct {
	buffer     *bytes.Buffer
	gzipWriter *gzip.Writer
	cpioWriter *cpio.Writer
}

// NewCpioArchive creates a new CpioArchive instance.
func NewCpioArchive() *CpioArchive {
	buf := new(bytes.Buffer)
	gw := gzip.NewWriter(buf)
	cw := cpio.NewWriter(gw)

	return &CpioArchive{
		buffer:     buf,
		gzipWriter: gw,
		cpioWriter: cw,
	}
}

// StoreBytes appends to the current archive the given content using
// the specified filename.
func (ca *CpioArchive) StoreBytes(filename string, content []byte, mode int) error {
	header := cpio.Header{
		Name: filename,
		Mode: cpio.FileMode(mode),
		Size: int64(len(content)),
	}

	err := ca.cpioWriter.WriteHeader(&header)
	if err != nil {
		return err
	}

	_, err = ca.cpioWriter.Write(content)
	if err != nil {
		return err
	}

	return nil
}

// StorePath adds a new path in the archive.
func (ca *CpioArchive) StorePath(path string) error {
	header := cpio.Header{
		Name: path,
		Mode: cpio.ModeDir | 0o755,
		Size: 0,
	}

	err := ca.cpioWriter.WriteHeader(&header)
	if err != nil {
		return err
	}

	return nil
}

// StoreFile appends to the current archive the specified file.
func (ca *CpioArchive) StoreFile(filename string, dstPath string) error {
	fileInfo, err := os.Lstat(filename)
	if err != nil {
		return err
	}

	header, err := cpio.FileInfoHeader(fileInfo, "")
	if err != nil {
		return err
	}
	header.Name = filepath.Join(dstPath, header.Name)

	if err := ca.cpioWriter.WriteHeader(header); err != nil {
		return err
	}

	fm := fileInfo.Mode()
	switch {
	case fm.IsRegular():
		// Copy the file content
		f, err := os.Open(filename)
		if err != nil {
			return err
		}
		defer f.Close()

		_, err = io.Copy(ca.cpioWriter, f)
		if err != nil {
			return err
		}
	case fm&os.ModeSymlink != 0:
		// In case of a symbolic link, copy the link text
		s, err := os.Readlink(filename)
		if err != nil {
			return err
		}
		_, err = ca.cpioWriter.Write([]byte(s))
		if err != nil {
			return err
		}
	default:
		return fmt.Errorf("unsupported file mode %v for file %s", fm, filename)
	}

	return nil
}

// SaveBuffer saves the content of the current archive and returns
// the buffer content.
func (ca *CpioArchive) SaveBuffer() ([]byte, error) {
	err := ca.cpioWriter.Close()
	if err != nil {
		return nil, err
	}

	err = ca.gzipWriter.Close()
	if err != nil {
		return nil, err
	}

	return ca.buffer.Bytes(), nil
}

// Save the content of the current archive on the disk.
func (ca *CpioArchive) Save(archivePath string) error {
	out, err := os.Create(archivePath)
	if err != nil {
		return err
	}
	defer out.Close()

	bs, err := ca.SaveBuffer()
	if err != nil {
		return err
	}

	_, err = out.Write(bs)
	if err != nil {
		return err
	}

	return nil
}
