package gci

import (
	"github.com/daixiang0/gci/pkg/gci"
)

// listCmd represents the list command
func (e *Executor) initList() {
	e.newGciCommand(
		"list path...",
		"Prints filenames that need to be formatted to STDOUT",
		"Prints the filenames that need to be formatted. If you want to show the diff use diff instead, and if you want to apply the changes use write instead",
		[]string{},
		false,
		gci.ListUnFormattedFiles)
}
