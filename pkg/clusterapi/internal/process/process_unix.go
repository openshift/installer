//go:build !windows
// +build !windows

package process

import (
	"os/exec"
	"syscall"
)

func setSysProcAttr(cmd *exec.Cmd) {
	cmd.SysProcAttr = &syscall.SysProcAttr{
		Setpgid: true,
	}
}

func stopProcess(ps *State) error {
	return ps.Cmd.Process.Signal(syscall.SIGTERM)
}

func killProcess(ps *State) error {
	return ps.Cmd.Process.Signal(syscall.SIGKILL)
}
