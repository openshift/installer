package gci

import (
	"github.com/daixiang0/gci/pkg/gci"
)

// printCmd represents the print command
func (e *Executor) initPrint() {
	e.newGciCommand(
		"print path...",
		"Outputs the formatted file to STDOUT",
		"Print outputs the formatted file. If you want to apply the changes to a file use write instead!",
		[]string{"output"},
		true,
		gci.PrintFormattedFiles)
}
