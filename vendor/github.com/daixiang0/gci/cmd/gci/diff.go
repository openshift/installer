package gci

import (
	"github.com/daixiang0/gci/pkg/gci"
)

// diffCmd represents the diff command
func (e *Executor) initDiff() {
	e.newGciCommand(
		"diff path...",
		"Prints a git style diff to STDOUT",
		"Diff prints a patch in the style of the diff tool that contains the required changes to the file to make it adhere to the specified formatting.",
		[]string{},
		true,
		gci.DiffFormattedFiles)
}
