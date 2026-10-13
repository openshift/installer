package gci

import (
	"github.com/daixiang0/gci/pkg/gci"
)

// writeCmd represents the write command
func (e *Executor) initWrite() {
	e.newGciCommand(
		"write path...",
		"Formats the specified files in-place",
		"Write modifies the specified files in-place",
		[]string{"overwrite"},
		false,
		gci.WriteFormattedFiles)
}
