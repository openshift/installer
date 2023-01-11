//go:build go1.18
// +build go1.18

package testscript

import (
	"reflect"
	"testing"
	"time"
)

func mainStart() *testing.M {
	// Note: testing.MainStart acquired an extra argument in Go 1.18.
	return testing.MainStart(nopTestDeps{}, nil, nil, nil, nil)
}

// Note: corpusEntry is an anonymous struct type used by some method stubs.
type corpusEntry = struct {
	Parent     string
	Path       string
	Data       []byte
	Values     []interface{}
	Generation int
	IsSeed     bool
}

// Note: CoordinateFuzzing was added in Go 1.18.
func (nopTestDeps) CoordinateFuzzing(time.Duration, int64, time.Duration, int64, int, []corpusEntry, []reflect.Type, string, string) error {
	return nil
}

// Note: RunFuzzWorker was added in Go 1.18.
func (nopTestDeps) RunFuzzWorker(func(corpusEntry) error) error {
	return nil
}

// Note: ReadCorpus was added in Go 1.18.
func (nopTestDeps) ReadCorpus(string, []reflect.Type) ([]corpusEntry, error) {
	return nil, nil
}

// Note: CheckCorpus was added in Go 1.18.
func (nopTestDeps) CheckCorpus([]interface{}, []reflect.Type) error {
	return nil
}

// Note: ResetCoverage was added in Go 1.18.
func (nopTestDeps) ResetCoverage() {}

// Note: SnapshotCoverage was added in Go 1.18.
func (nopTestDeps) SnapshotCoverage() {}
