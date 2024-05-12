package ops

import (
	"bytes"
	"fmt"
	"io"
	"os/exec"
	"strings"
	"syscall"

	"github.com/sirupsen/logrus"
)

// Execute is an interface for executing external commands and capturing their output.
//
//go:generate mockgen -source=execute.go -package=ops -destination=mock_execute.go
type Execute interface {
	Execute(command string, args ...string) (string, error)
	ExecuteWithLiveLogger(command string, args ...string) (string, error)
}

type executor struct {
	log     *logrus.Logger
	verbose bool
}

func (e *executor) execute(liveLogger io.Writer, root, command string, args ...string) (string, error) {
	e.log.Infof("Executing %s with args %s", command, args)
	cmd := exec.Command(command, args...)
	var stdoutBytes bytes.Buffer
	if liveLogger != nil {
		cmd.Stdout = io.MultiWriter(liveLogger, &stdoutBytes)
		cmd.Stderr = io.MultiWriter(liveLogger, &stdoutBytes)
	} else {
		cmd.Stdout = &stdoutBytes
		cmd.Stderr = &stdoutBytes
	}
	if root != "" {
		cmd.SysProcAttr = &syscall.SysProcAttr{Chroot: root}
		cmd.Dir = "/"
	}

	err := cmd.Run()
	stdoutBytesTrimmed := strings.TrimSpace(stdoutBytes.String())
	if err != nil {
		return stdoutBytesTrimmed, fmt.Errorf("%s: %w", stdoutBytes.String(), err)
	}
	return stdoutBytesTrimmed, nil
}

type regularExecutor struct {
	executor
}

func NewRegularExecutor(logger *logrus.Logger, verbose bool) Execute {
	return &regularExecutor{executor: executor{logger, verbose}}
}

func (e *regularExecutor) Execute(command string, args ...string) (string, error) {
	return e.executor.execute(nil, "", command, args...)
}

func (e *regularExecutor) ExecuteWithLiveLogger(command string, args ...string) (string, error) {
	return e.executor.execute(e.executor.log.Writer(), "", command, args...)
}

type nsenterExecutor struct {
	executor
}

func NewNsenterExecutor(logger *logrus.Logger, verbose bool) Execute {
	return &nsenterExecutor{executor: executor{logger, verbose}}
}

func (e *nsenterExecutor) ExecuteWithLiveLogger(command string, args ...string) (string, error) {
	return e.baseExecute(e.executor.log.Writer(), command, args...)
}

func (e *nsenterExecutor) baseExecute(writer io.Writer, command string, args ...string) (string, error) {
	// nsenter is used here to launch processes inside the container in a way that makes said processes feel
	// and behave as if they're running on the host directly rather than inside the container
	commandBase := "nsenter"

	arguments := []string{
		"--target", "1",
		// Entering the cgroup namespace is not required for podman on CoreOS (where the
		// agent typically runs), but it's needed on some Fedora versions and
		// some other systemd based systems. Those systems are used to run dry-mode
		// agents for load testing. If this flag is not used, Podman will sometimes
		// have trouble creating a systemd cgroup slice for new containers.
		"--cgroup",
		// The mount namespace is required for podman to access the host's container
		// storage
		"--mount",
		// TODO: Document why we need the IPC namespace
		"--ipc",
		"--pid",
		"--",
		command,
	}

	arguments = append(arguments, args...)
	return e.executor.execute(writer, "", commandBase, arguments...)
}

func (e *nsenterExecutor) Execute(command string, args ...string) (string, error) {
	return e.baseExecute(nil, command, args...)
}

type chrootExecutor struct {
	executor
	root string
}

func NewChrootExecutor(logger *logrus.Logger, verbose bool, root string) Execute {
	return &chrootExecutor{executor: executor{logger, verbose}, root: root}
}

// Running a command with chroot using exec.Command runs into issues with exec.LookPath,
// if an absolute path is not used for the "command", as it does not account for the chroot dir.
// To workaround this issue, prefix the command with /usr/bin/env.
func (e *chrootExecutor) baseExecute(writer io.Writer, command string, args ...string) (string, error) {
	commandBase := "/usr/bin/env"
	args = append([]string{"--", command}, args...)
	return e.executor.execute(writer, e.root, commandBase, args...)
}

func (e *chrootExecutor) Execute(command string, args ...string) (string, error) {
	return e.baseExecute(nil, command, args...)
}

func (e *chrootExecutor) ExecuteWithLiveLogger(command string, args ...string) (string, error) {
	return e.baseExecute(e.executor.log.Writer(), command, args...)
}
