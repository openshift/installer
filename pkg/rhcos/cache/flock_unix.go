//go:build !windows
// +build !windows

package cache

import (
	"os"

	"golang.org/x/sys/unix"
)

func flockFile(f *os.File, lock bool) error {
	if lock {
		return unix.Flock(int(f.Fd()), unix.LOCK_EX)
	}
	return unix.Flock(int(f.Fd()), unix.LOCK_UN)
}
