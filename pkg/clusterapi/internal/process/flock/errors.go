package flock

import "errors"

var (
	// ErrAlreadyLocked is returned when the file is already locked.
	ErrAlreadyLocked = errors.New("the file is already locked")
)
