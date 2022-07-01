package ovirtclientlog

import (
	"context"
	"testing"
)

// NewTestLogger returns a logger that logs via the Go test facility
func NewTestLogger(t *testing.T) Logger {
	return &testLogger{
		t: t,
	}
}

type testLogger struct {
	t *testing.T
}

func (t *testLogger) WithContext(_ context.Context) Logger {
	return t
}

func (t *testLogger) Debugf(format string, args ...interface{}) {
	t.t.Helper()
	t.t.Logf(format, args...)
}

func (t *testLogger) Infof(format string, args ...interface{}) {
	t.t.Helper()
	t.t.Logf(format, args...)
}

func (t *testLogger) Warningf(format string, args ...interface{}) {
	t.t.Helper()
	t.t.Logf(format, args...)
}

func (t *testLogger) Errorf(format string, args ...interface{}) {
	t.t.Helper()
	t.t.Logf(format, args...)
}
