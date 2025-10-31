//go:build windows
// +build windows

package cache

import (
	"os"

	"golang.org/x/sys/windows"
)

func flockFile(f *os.File, lock bool) error {
	if lock {
		var overlapped windows.Overlapped
		return windows.LockFileEx(windows.Handle(f.Fd()), windows.LOCKFILE_EXCLUSIVE_LOCK, 0, 1, 0, &overlapped)
	}
	var overlapped windows.Overlapped
	return windows.UnlockFileEx(windows.Handle(f.Fd()), 0, 1, 0, &overlapped)
}
