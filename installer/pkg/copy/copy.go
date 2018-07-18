// Package copy supports copying a file to another path.
package copy

import (
	"io"
	"os"
)

// Copy creates a new file at toFilePath with with mode 0666 (before
// umask) and the same content as fromFilePath.
func Copy(fromFilePath, toFilePath string) error {
	from, err := os.Open(fromFilePath)
	if err != nil {
		return err
	}
	defer from.Close()

	to, err := os.OpenFile(toFilePath, os.O_WRONLY|os.O_CREATE, 0666)
	if err != nil {
		return err
	}
	defer to.Close()

	_, err = io.Copy(to, from)
	return err
}
