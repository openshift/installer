//go:build windows
// +build windows

package exec

import (
	"os"
)

var ignoreSignals = []os.Signal{os.Interrupt}
var forwardSignals []os.Signal
