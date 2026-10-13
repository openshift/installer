package isoeditor

import (
	"bytes"
	"compress/gzip"
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	"github.com/cavaliercoder/go-cpio"
	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/pkg/errors"
)

const (
	AMD64CPUArchitecture   = "amd64"
	X86CPUArchitecture     = "x86_64"
	ARM64CPUArchitecture   = "arm64"
	AARCH64CPUArchitecture = "aarch64"
)

// Extract unpacks the iso contents into the working directory
func Extract(isoPath string, workDir string) error {
	d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return err
	}

	fs, err := GetISO9660FileSystem(d)
	if err != nil {
		return err
	}

	files, err := fs.ReadDir("/")
	if err != nil {
		return err
	}
	err = copyAll(fs, "/", files, workDir)
	if err != nil {
		return err
	}

	return nil
}

// recursive function for unpacking all files and directores from the given iso filesystem starting at fsDir
func copyAll(fs filesystem.FileSystem, fsDir string, infos []os.FileInfo, targetDir string) error {
	for _, info := range infos {
		osName := filepath.Join(targetDir, info.Name())
		fsName := filepath.Join(fsDir, info.Name())

		if info.IsDir() {
			if err := os.Mkdir(osName, info.Mode().Perm()); err != nil {
				return err
			}

			files, err := fs.ReadDir(fsName)
			if err != nil {
				return err
			}
			if err := copyAll(fs, fsName, files[:], osName); err != nil {
				return err
			}
		} else {
			fsFile, err := fs.OpenFile(fsName, os.O_RDONLY)
			if err != nil {
				return err
			}
			osFile, err := os.Create(osName)
			if err != nil {
				return err
			}

			_, err = io.Copy(osFile, fsFile)
			if err != nil {
				osFile.Close()
				return err
			}

			if err := osFile.Sync(); err != nil {
				osFile.Close()
				return err
			}

			if err := osFile.Close(); err != nil {
				return err
			}
		}
	}
	return nil
}

// Create builds an iso file at outPath with the given volumeLabel using the contents of the working directory
func Create(outPath string, workDir string, volumeLabel string) error {
	// Use the minimum iso size that will satisfy diskfs validations here.
	// This value doesn't determine the final image size, but is used
	// to truncate the initial file. This value would be relevant if
	// we were writing to a particular partition on a device, but we are
	// not so the minimum iso size will work for us here
	minISOSize := 38 * 1024
	d, err := diskfs.Create(outPath, int64(minISOSize), diskfs.Raw, diskfs.SectorSizeDefault)
	if err != nil {
		return err
	}

	d.LogicalBlocksize = 2048
	fspec := disk.FilesystemSpec{
		Partition:   0,
		FSType:      filesystem.TypeISO9660,
		VolumeLabel: volumeLabel,
		WorkDir:     workDir,
	}
	fs, err := d.CreateFilesystem(fspec)
	if err != nil {
		return err
	}

	iso, ok := fs.(*iso9660.FileSystem)
	if !ok {
		return fmt.Errorf("not an iso9660 filesystem")
	}

	options := iso9660.FinalizeOptions{
		RockRidge:        true,
		VolumeIdentifier: volumeLabel,
	}

	if haveFiles, err := haveBootFiles(workDir); err != nil {
		return err
	} else if haveFiles {
		efiSectors, err := efiLoadSectors(workDir)
		if err != nil {
			return err
		}
		options.ElTorito = &iso9660.ElTorito{
			BootCatalog: "isolinux/boot.cat",
			Entries: []*iso9660.ElToritoEntry{
				{
					Platform:  iso9660.BIOS,
					Emulation: iso9660.NoEmulation,
					BootFile:  "isolinux/isolinux.bin",
					BootTable: true,
					LoadSize:  4,
				},
				{
					Platform:  iso9660.EFI,
					Emulation: iso9660.NoEmulation,
					BootFile:  "images/efiboot.img",
					LoadSize:  efiSectors,
				},
			},
		}
	} else if exists, _ := fileExists(filepath.Join(workDir, "images/efiboot.img")); exists {
		// Creating an ISO with EFI boot only
		efiSectors, err := efiLoadSectors(workDir)
		if err != nil {
			return err
		}
		if exists, _ := fileExists(filepath.Join(workDir, "boot.catalog")); !exists {
			return fmt.Errorf("missing boot.catalog file")
		}
		options.ElTorito = &iso9660.ElTorito{
			BootCatalog:     "boot.catalog",
			HideBootCatalog: true,
			Entries: []*iso9660.ElToritoEntry{
				{
					Platform:  iso9660.EFI,
					Emulation: iso9660.NoEmulation,
					BootFile:  "images/efiboot.img",
					LoadSize:  efiSectors,
				},
			},
		}
	} else if exists, _ := fileExists(filepath.Join(workDir, "images/cdboot.img")); exists {
		// Creating an ISO for S390 boot:
		cdbootSectors, err := cdbootLoadSectors(workDir)
		if err != nil {
			return err
		}
		if exists, _ := fileExists(filepath.Join(workDir, "boot.catalog")); !exists {
			return fmt.Errorf("missing boot.catalog file")
		}
		options.ElTorito = &iso9660.ElTorito{
			BootCatalog:     "boot.catalog",
			HideBootCatalog: true,
			Entries: []*iso9660.ElToritoEntry{
				{
					Platform:  iso9660.BIOS,
					Emulation: iso9660.NoEmulation,
					BootFile:  "images/cdboot.img",
					LoadSize:  cdbootSectors,
				},
			},
		}
	}

	return iso.Finalize(options)
}

// Returns the number of sectors to load for efi boot
// Load Sectors * 2048 should be the size of efiboot.img rounded up to a multiple of 2048
// For UEFI boot, the sector size is 512
// To support iso9660 (2048) and UEFI (512), sectors must be in blocks of 512, but must also be a multiple of 2048
func efiLoadSectors(workDir string) (uint16, error) {
	efiStat, err := os.Stat(filepath.Join(workDir, "images/efiboot.img"))
	if err != nil {
		return 0, err
	}
	return uint16(math.Ceil(float64(efiStat.Size())/2048) * 4), nil
}

func cdbootLoadSectors(workDir string) (result uint16, err error) {
	// Calculate the number of 512 sectors that would be needed for the boot image:
	info, err := os.Stat(filepath.Join(workDir, "images/cdboot.img"))
	if err != nil {
		return
	}
	size := info.Size()
	sectors := size / 512
	if size%512 != 0 {
		sectors++
	}

	// Some BIOS may have problems if this number isn't a multiple of four, because then it
	// won't fit into complete the 2048 byte blocks used by CD devices. So we need to round
	// up to the closest multiple of four.
	if sectors%4 != 0 {
		sectors += 4 - sectors%4
	}

	// The resulting number will not fit in a 16 bit unsigned number if the file is 32 MiB
	// or larger. For example, in OpenShift 4.14.0 the s390x boot image is 64 MiB:
	//
	//	# mount -o loop,ro rhcos-4.14.0-s390x-live.s390x.iso /mnt
	//	# ls -lh /mnt/images/cdboot.img
	//	-r--r--r--. 1 root root 64M Sep 20 20:16 /mnt/images/cdboot.img
	//
	// In that case the 'xorrisofs' tool that is used to generate those ISO files sets the
	// field to 65535, the maximum value for a 16 bits unsigned integer. That is incorrect,
	// but seems to work, so we do the same.
	if sectors > math.MaxUint16 {
		sectors = math.MaxUint16
	}
	// nolint: gosec
	result = uint16(sectors)
	return
}

func haveBootFiles(workDir string) (bool, error) {
	files := []string{"isolinux/boot.cat", "isolinux/isolinux.bin", "images/efiboot.img"}
	for _, f := range files {
		name := filepath.Join(workDir, f)
		if exists, err := fileExists(name); err != nil {
			return false, err
		} else if !exists {
			return false, nil
		}
	}

	return true, nil
}

func fileExists(name string) (bool, error) {
	if _, err := os.Stat(name); os.IsNotExist(err) {
		return false, nil
	} else if err != nil {
		return false, err
	}
	return true, nil
}

func VolumeIdentifier(isoPath string) (string, error) {
	// Need to get the volume id from the ISO provided
	iso, err := os.Open(isoPath)
	if err != nil {
		return "", err
	}
	defer iso.Close()

	// Need a method to identify the ISO provided
	// The first 32768 bytes are unused by the ISO 9660 standard, typically for bootable media
	// This is where the data area begins and the 32 byte string representing the volume identifier
	// is offset 40 bytes into the primary volume descriptor
	volumeId := make([]byte, 32)
	_, err = iso.ReadAt(volumeId, 32808)
	if err != nil {
		return "", err
	}

	return strings.TrimSpace(string(volumeId)), nil
}

func GetISOFileInfo(filePath, isoPath string) (int64, int64, error) {
	d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return 0, 0, err
	}

	fs, err := GetISO9660FileSystem(d)
	if err != nil {
		return 0, 0, err
	}

	fsFile, err := fs.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		return 0, 0, errors.Wrapf(err, "Failed to open file %s", filePath)
	}

	defer fsFile.Close()
	isoFile := fsFile.(*iso9660.File)
	defaultSectorSize := uint32(2 * 1024)
	return int64(isoFile.Location() * defaultSectorSize), isoFile.Size(), nil
}

// Gets a readWrite seeker of a specific file from the ISO image
func GetFileFromISO(isoPath, filePath string) (filesystem.File, error) {
	d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return nil, err
	}

	fs, err := GetISO9660FileSystem(d)
	if err != nil {
		return nil, err
	}

	file, err := fs.OpenFile(filePath, os.O_RDONLY)
	if err != nil {
		return nil, err
	}
	return file, nil
}

// Reads a whole specific file from the ISO image
func ReadFileFromISO(isoPath, filePath string) ([]byte, error) {
	f, err := GetFileFromISO(isoPath, filePath)
	if err != nil {
		return nil, err
	}
	defer f.Close()
	ret, err := io.ReadAll(f)
	if err != nil {
		return nil, err
	}
	return ret, nil
}

// Gets directly the ISO 9660 filesystem (equivalent to GetFileSystem(0)).
func GetISO9660FileSystem(d *disk.Disk) (filesystem.FileSystem, error) {
	return iso9660.Read(d.File, d.Size, 0, 0)
}

// fileEntry represents a single file to be added to a CPIO archive
type fileEntry struct {
	Content []byte
	Path    string
	Mode    cpio.FileMode
}

func generateCompressedCPIO(files []fileEntry) ([]byte, error) {
	// Run gzip compression
	compressedBuffer := new(bytes.Buffer)
	gzipWriter := gzip.NewWriter(compressedBuffer)
	// Create CPIO archive
	cpioWriter := cpio.NewWriter(gzipWriter)

	// Add each file to the archive
	for _, file := range files {
		if err := cpioWriter.WriteHeader(&cpio.Header{
			Name: file.Path,
			Mode: file.Mode,
			Size: int64(len(file.Content)),
		}); err != nil {
			return nil, errors.Wrap(err, "Failed to write CPIO header")
		}
		if _, err := cpioWriter.Write(file.Content); err != nil {
			return nil, errors.Wrap(err, "Failed to write CPIO archive")
		}
	}

	if err := cpioWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to close CPIO archive")
	}
	if err := gzipWriter.Close(); err != nil {
		return nil, errors.Wrap(err, "Failed to gzip ignition config")
	}

	padSize := (4 - (compressedBuffer.Len() % 4)) % 4
	for i := 0; i < padSize; i++ {
		if err := compressedBuffer.WriteByte(0); err != nil {
			return nil, err
		}
	}

	return compressedBuffer.Bytes(), nil
}
