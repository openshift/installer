// +build !windows

package exec

import (
	"os"
	"syscall"
)

var ignoreSignals = []os.Signal{os.Interrupt}
var forwardSignals = []os.Signal{syscall.SIGTERM}
