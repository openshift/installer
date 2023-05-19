package isoeditor

import (
	"fmt"
	"io"
	"math"
	"os"
	"path/filepath"
	"strings"

	diskfs "github.com/diskfs/go-diskfs"
	"github.com/diskfs/go-diskfs/disk"
	"github.com/diskfs/go-diskfs/filesystem"
	"github.com/diskfs/go-diskfs/filesystem/iso9660"
	"github.com/pkg/errors"
)

// Extract unpacks the iso contents into the working directory
func Extract(isoPath string, workDir string) error {
	d, err := diskfs.Open(isoPath, diskfs.WithOpenMode(diskfs.ReadOnly))
	if err != nil {
		return err
	}

	fs, err := d.GetFilesystem(0)
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
	}

	return iso.Finalize(options)
}

// Returns the number of sectors to load for efi boot
// Load Sectors * 2048 should be the size of efiboot.img rounded up to a multiple of 2048
func efiLoadSectors(workDir string) (uint16, error) {
	efiStat, err := os.Stat(filepath.Join(workDir, "images/efiboot.img"))
	if err != nil {
		return 0, err
	}
	return uint16(math.Ceil(float64(efiStat.Size()) / 2048)), nil
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

	fs, err := d.GetFilesystem(0)
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

	fs, err := d.GetFilesystem(0)
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
