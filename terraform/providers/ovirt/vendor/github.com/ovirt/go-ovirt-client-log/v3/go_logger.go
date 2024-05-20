package ovirtclientlog

import (
	"context"
	"fmt"
	"log"
)

// NewGoLogger creates a logger that writes to the Go log facility. The optional logger parameter can be
// used to pass one scoped logger, otherwise the global logger is used. If multiple loggers are passed the
// function will panic.
func NewGoLogger(logger ...*log.Logger) Logger {
	var l *log.Logger = nil
	if len(logger) == 1 {
		l = logger[0]
	} else {
		panic(fmt.Sprintf("Only one logger may be passed to NewGoLogger, %d were passed.", len(logger)))
	}
	return &goLogger{
		logger: l,
	}
}

type goLogger struct {
	logger *log.Logger
}

func (g *goLogger) WithContext(_ context.Context) Logger {
	return g
}

func (g *goLogger) write(format string, args ...interface{}) {
	if g.logger == nil {
		log.Printf(fmt.Sprintf("%s\n", format), args...)
	} else {
		g.logger.Printf(fmt.Sprintf("%s\n", format), args...)
	}
}

func (g *goLogger) Debugf(format string, args ...interface{}) {
	g.write(format, args...)
}

func (g *goLogger) Infof(format string, args ...interface{}) {
	g.write(format, args...)
}

func (g *goLogger) Warningf(format string, args ...interface{}) {
	g.write(format, args...)
}

func (g *goLogger) Errorf(format string, args ...interface{}) {
	g.write(format, args...)
}
