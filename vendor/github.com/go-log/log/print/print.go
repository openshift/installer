// Package print allows users to create a Logger interface from any
// object that supports Print and Printf.
package print

// Printer is an interface for Print and Printf.
type Printer interface {
	Print(v ...interface{})
	Printf(format string, v ...interface{})
}

type logger struct{
	printer Printer
}

func (logger *logger) Log(v ...interface{}) {
	logger.printer.Print(v...)
}

func (logger *logger) Logf(format string, v ...interface{}) {
	logger.printer.Printf(format, v...)
}

// New creates a new logger wrapping printer.
func New(printer Printer) *logger {
	return &logger{
		printer: printer,
	}
}
