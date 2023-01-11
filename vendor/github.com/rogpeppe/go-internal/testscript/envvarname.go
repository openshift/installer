//go:build !windows
// +build !windows

package testscript

func envvarname(k string) string {
	return k
}
