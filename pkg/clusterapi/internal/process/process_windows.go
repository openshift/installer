//go:build windows
// +build windows

package process

import (
	"os/exec"
)

func setSysProcAttr(cmd *exec.Cmd) {
}

func stopProcess(ps *State) error {
	return ps.Cmd.Process.Kill()
}

func killProcess(ps *State) error {
	return ps.Cmd.Process.Kill()
}
