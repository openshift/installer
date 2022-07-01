package ovirtclientlog

import "context"

// NewNOOPLogger returns a logger that does nothing.
func NewNOOPLogger() Logger {
	return &noopLogger{}
}

type noopLogger struct {
}

func (n noopLogger) WithContext(_ context.Context) Logger {
	return n
}

func (n noopLogger) Debugf(_ string, _ ...interface{}) {

}

func (n noopLogger) Infof(_ string, _ ...interface{}) {

}

func (n noopLogger) Warningf(_ string, _ ...interface{}) {

}

func (n noopLogger) Errorf(_ string, _ ...interface{}) {

}
