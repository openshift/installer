//go:build !go1.18
// +build !go1.18

package testscript

import (
	"testing"
)

func mainStart() *testing.M {
	return testing.MainStart(nopTestDeps{}, nil, nil, nil)
}
