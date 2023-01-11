// Package lineprinter wraps a Print implementation to provide an io.WriteCloser.
package lineprinter

import (
	"bytes"
	"io"
	"sync"
)

// Print is a type that can hold fmt.Print and other implementations
// which match that signature.  For example, you can use:
//
//	trimmer := &lineprinter.Trimmer{WrappedPrint: logrus.StandardLogger().Debug}
//	linePrinter := &linePrinter{Print: trimmer.Print}
//
// to connect the line printer to logrus at the debug level.
type Print func(args ...interface{})

// LinePrinter is an io.WriteCloser that buffers written bytes.
// During each Write, newline-terminated lines are removed from the
// buffer and passed to Print.  On Close, any content remaining in the
// buffer is also passed to Print.
//
// One use-case is connecting a subprocess's standard streams to a
// logger:
//
//	linePrinter := &linePrinter{
//	  Print: &Trimmer{WrappedPrint: logrus.StandardLogger().Debug}.Print,
//	}
//	defer linePrinter.Close()
//	cmd := exec.Command(...)
//	cmd.Stdout = linePrinter
//
// LinePrinter buffers the subcommand's byte stream and splits it into
// lines for the logger.  Sometimes we might have a partial line
// written to the buffer.  We don't want to push that partial line into
// the logs if the next read attempt will pull in the remainder of the
// line.  But we do want to push that partial line into the logs if there
// will never be a next read.  So LinePrinter.Write pushes any
// complete lines to the wrapped printer, and LinePrinter.Close pushes
// any remaining partial line.
type LinePrinter struct {
	buf   bytes.Buffer
	Print Print

	sync.Mutex
}

// Write writes len(p) bytes from p to an internal buffer.  Then it
// retrieves any newline-terminated lines from the internal buffer and
// prints them with lp.Print.  Partial lines are left in the internal
// buffer.
func (lp *LinePrinter) Write(p []byte) (int, error) {
	lp.Lock()
	defer lp.Unlock()
	n, err := lp.buf.Write(p)
	if err != nil {
		return n, err
	}

	for {
		line, err := lp.buf.ReadString(byte('\n'))
		if err == io.EOF {
			_, err = lp.buf.Write([]byte(line))
			return n, err
		} else if err != nil {
			return n, err
		}

		lp.Print(line)
	}
}

// Close prints anything that remains in the buffer.
func (lp *LinePrinter) Close() error {
	lp.Lock()
	defer lp.Unlock()
	line := lp.buf.String()
	if len(line) > 0 {
		lp.Print(line)
	}
	return nil
}
