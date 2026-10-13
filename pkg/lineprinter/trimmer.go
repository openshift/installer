package lineprinter

import (
	"strings"
)

// Trimmer is a Print wrapper that removes trailing newlines from the
// final argument (if it is a string argument).  This is useful for
// connecting a LinePrinter to a logger whose Print-analog does not
// expect trailing newlines.
type Trimmer struct {
	WrappedPrint Print
}

// Print removes trailing newlines from the final argument (if it is a
// string argument) and then passes the arguments through to
// WrappedPrint.
func (t *Trimmer) Print(args ...interface{}) {
	if len(args) > 0 {
		i := len(args) - 1
		arg, ok := args[i].(string)
		if ok {
			args[i] = strings.TrimRight(arg, "\n")
		}
	}
	t.WrappedPrint(args...)
}
