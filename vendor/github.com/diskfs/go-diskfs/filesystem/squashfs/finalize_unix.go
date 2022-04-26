// +build aix darwin dragonfly freebsd js,wasm linux nacl netbsd openbsd solaris

package squashfs

import (
	"os"

	"golang.org/x/sys/unix"
)

func getDeviceNumbers(path string) (uint32, uint32, error) {
	stat := unix.Stat_t{}
	err := unix.Stat(path, &stat)
	if err != nil {
		return 0, 0, err
	}
	return uint32(stat.Rdev / 256), uint32(stat.Rdev % 256), nil
}

func getFileProperties(fi os.FileInfo) (uint32, uint32, uint32) {
	var nlink, uid, gid uint32
	if sys := fi.Sys(); sys != nil {
		if stat, ok := sys.(*unix.Stat_t); ok {
			nlink = uint32(stat.Nlink)
			uid = stat.Uid
			gid = stat.Gid
		}
	}
	return nlink, uid, gid
}
