// +build darwin linux solaris aix freebsd illumos netbsd openbsd plan9

package diskfs

import (
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// getSectorSizes get the logical and physical sector sizes for a block device
func getSectorSizes(f *os.File) (int64, int64, error) {
	/*
		ioctl(fd, BLKPBSZGET, &physicalsectsize);

	*/
	fd := f.Fd()
	logicalSectorSize, err := unix.IoctlGetInt(int(fd), blksszGet)
	if err != nil {
		return 0, 0, fmt.Errorf("Unable to get device logical sector size: %v", err)
	}
	physicalSectorSize, err := unix.IoctlGetInt(int(fd), blkpbszGet)
	if err != nil {
		return 0, 0, fmt.Errorf("Unable to get device physical sector size: %v", err)
	}
	return int64(logicalSectorSize), int64(physicalSectorSize), nil
}
