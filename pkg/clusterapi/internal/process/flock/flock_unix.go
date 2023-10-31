//go:build linux || darwin || freebsd || openbsd || netbsd || dragonfly
// +build linux darwin freebsd openbsd netbsd dragonfly

package flock

import (
	"errors"
	"fmt"
	"os"

	"golang.org/x/sys/unix"
)

// Acquire acquires a lock on a file for the duration of the process. This method
// is reentrant.
func Acquire(path string) error {
	fd, err := unix.Open(path, unix.O_CREAT|unix.O_RDWR|unix.O_CLOEXEC, 0600)
	if err != nil {
		if errors.Is(err, os.ErrExist) {
			return fmt.Errorf("cannot lock file %q: %w", path, ErrAlreadyLocked)
		}
		return err
	}

	// We don't need to close the fd since we should hold
	// it until the process exits.
	err = unix.Flock(fd, unix.LOCK_NB|unix.LOCK_EX)
	if errors.Is(err, unix.EWOULDBLOCK) { // This condition requires LOCK_NB.
		return fmt.Errorf("cannot lock file %q: %w", path, ErrAlreadyLocked)
	}
	return err
}
