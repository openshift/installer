package executer

import (
	"bytes"
	"context"
	"os"
	"os/exec"
)

//go:generate mockgen --build_flags=--mod=mod -package executer -destination mock_executer.go . Executer
type Executer interface {
	Execute(command string, args ...string) (stdout string, stderr string, exitCode int)
	ExecuteWithContext(ctx context.Context, command string, args ...string) (stdout string, stderr string, exitCode int)
	TempFile(dir, pattern string) (f *os.File, err error)
}

type CommonExecuter struct{}

func (e *CommonExecuter) TempFile(dir, pattern string) (f *os.File, err error) {
	return os.CreateTemp(dir, pattern)
}

func (e *CommonExecuter) Execute(command string, args ...string) (stdout string, stderr string, exitCode int) {
	cmd := exec.Command(command, args...)
	return e.execute(cmd)
}

func (e *CommonExecuter) ExecuteWithContext(ctx context.Context, command string, args ...string) (stdout string, stderr string, exitCode int) {
	cmd := exec.CommandContext(ctx, command, args...)
	return e.execute(cmd)
}

func (e *CommonExecuter) execute(cmd *exec.Cmd) (stdout string, stderr string, exitCode int) {
	var stdoutBytes, stderrBytes bytes.Buffer
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes
	err := cmd.Run()
	return stdoutBytes.String(), getErrorStr(err, &stderrBytes), getExitCode(err)
}

func getExitCode(err error) int {
	if err == nil {
		return 0
	}
	switch value := err.(type) {
	case *exec.ExitError:
		return value.ExitCode()
	default:
		return -1
	}
}

func getErrorStr(err error, stderr *bytes.Buffer) string {
	b := stderr.Bytes()
	if len(b) > 0 {
		return string(b)
	} else if err != nil {
		return err.Error()
	}
	return ""
}
