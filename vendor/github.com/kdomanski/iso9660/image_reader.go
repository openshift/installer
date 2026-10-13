package iso9660

import (
	"fmt"
	"io"
	"os"
	"strings"
	"time"
)

// Image is a wrapper around an image file that allows reading its ISO9660 data
type Image struct {
	ra                io.ReaderAt
	volumeDescriptors []volumeDescriptor
}

// OpenImage returns an Image reader reating from a given file
func OpenImage(ra io.ReaderAt) (*Image, error) {
	i := &Image{ra: ra}

	if err := i.readVolumes(); err != nil {
		return nil, err
	}

	return i, nil
}

func (i *Image) readVolumes() error {
	buffer := make([]byte, sectorSize)
	// skip the 16 sectors of system area
	for sector := 16; ; sector++ {
		if _, err := i.ra.ReadAt(buffer, int64(sector)*int64(sectorSize)); err != nil {
			return err
		}

		var vd volumeDescriptor
		if err := vd.UnmarshalBinary(buffer); err != nil {
			return err
		}

		i.volumeDescriptors = append(i.volumeDescriptors, vd)
		if vd.Header.Type == volumeTypeTerminator {
			break
		}
	}

	return nil
}

// RootDir returns the File structure corresponding to the root directory
// of the first primary volume
func (i *Image) RootDir() (*File, error) {
	for _, vd := range i.volumeDescriptors {
		if vd.Type() == volumeTypePrimary {
			return &File{de: vd.Primary.RootDirectoryEntry, ra: i.ra, children: nil}, nil
		}
	}
	return nil, fmt.Errorf("no primary volumes found")
}

// File is a os.FileInfo-compatible wrapper around an ISO9660 directory entry
type File struct {
	ra       io.ReaderAt
	de       *DirectoryEntry
	children []*File
}

var _ os.FileInfo = &File{}

// IsDir returns true if the entry is a directory or false otherwise
func (f *File) IsDir() bool {
	return f.de.FileFlags&dirFlagDir != 0
}

// ModTime returns the entry's recording time
func (f *File) ModTime() time.Time {
	return time.Time(f.de.RecordingDateTime)
}

// Mode returns os.FileMode flag set with the os.ModeDir flag enabled in case of directories
func (f *File) Mode() os.FileMode {
	var mode os.FileMode
	if f.IsDir() {
		mode |= os.ModeDir
	}
	return mode
}

// Name returns the base name of the given entry
func (f *File) Name() string {
	if f.IsDir() {
		return f.de.Identifier
	}

	// drop the version part
	// assume only one ';'
	fileIdentifier := strings.Split(f.de.Identifier, ";")[0]

	// split into filename and extension
	// assume only only one '.'
	splitFileIdentifier := strings.Split(fileIdentifier, ".")

	// there's no dot in the name, thus no extension
	if len(splitFileIdentifier) == 1 {
		return splitFileIdentifier[0]
	}

	// extension is empty, return just the name without a dot
	if len(splitFileIdentifier[1]) == 0 {
		return splitFileIdentifier[0]
	}

	// return file with extension
	return fileIdentifier
}

// Size returns the size in bytes of the extent occupied by the file or directory
func (f *File) Size() int64 {
	return int64(f.de.ExtentLength)
}

// Sys returns nil
func (f *File) Sys() interface{} {
	return nil
}

// GetChildren returns the chilren entries in case of a directory
// or an error in case of a file
func (f *File) GetChildren() ([]*File, error) {
	if !f.IsDir() {
		return nil, fmt.Errorf("%s is not a directory", f.Name())
	}

	if f.children != nil {
		return f.children, nil
	}

	baseOffset := uint32(f.de.ExtentLocation) * sectorSize

	buffer := make([]byte, sectorSize)
	for bytesProcessed := uint32(0); bytesProcessed < uint32(f.de.ExtentLength); bytesProcessed += sectorSize {
		if _, err := f.ra.ReadAt(buffer, int64(baseOffset+bytesProcessed)); err != nil {
			return nil, nil
		}

		for i := uint32(0); i < sectorSize; {
			entryLength := uint32(buffer[i])
			if entryLength == 0 {
				break
			}

			if i+entryLength > sectorSize {
				return nil, fmt.Errorf("reading directory entries: DE outside of sector boundries")
			}

			newDE := &DirectoryEntry{}
			if err := newDE.UnmarshalBinary(buffer[i : i+entryLength]); err != nil {
				return nil, err
			}
			i += entryLength
			if newDE.Identifier == string([]byte{0}) || newDE.Identifier == string([]byte{1}) {
				continue
			}

			newFile := &File{ra: f.ra,
				de:       newDE,
				children: nil,
			}

			f.children = append(f.children, newFile)
		}
	}

	return f.children, nil
}

// Reader returns a reader that allows to read the file's data.
// If File is a directory, it returns nil.
func (f *File) Reader() io.Reader {
	if f.IsDir() {
		return nil
	}

	baseOffset := int64(f.de.ExtentLocation) * int64(sectorSize)
	return io.NewSectionReader(f.ra, baseOffset, int64(f.de.ExtentLength))
}
