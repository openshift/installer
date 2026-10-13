package agent

import (
	"bytes"
	"errors"
	"fmt"
	"os"
	"os/exec"
	"strings"
)

// ExecuteOC will execute an oc command.
func ExecuteOC(pullSecret string, command []string) (string, error) {
	// create registry config
	ps, err := os.CreateTemp("", "registry-config")
	if err != nil {
		return "", err
	}
	defer func() {
		ps.Close()
		os.Remove(ps.Name())
	}()
	_, err = ps.Write([]byte(pullSecret))
	if err != nil {
		return "", err
	}
	// flush the buffer to ensure the file can be read
	ps.Close()
	registryConfig := "--registry-config=" + ps.Name()
	command = append(command, registryConfig)
	var stdoutBytes, stderrBytes bytes.Buffer
	cmd := exec.Command(command[0], command[1:]...) // #nosec G204
	cmd.Stdout = &stdoutBytes
	cmd.Stderr = &stderrBytes

	err = cmd.Run()
	if err == nil {
		return strings.TrimSpace(stdoutBytes.String()), nil
	}

	var exitErr *exec.ExitError
	if errors.As(err, &exitErr) {
		err = fmt.Errorf("command '%s' exited with non-zero exit code %d: %s\n%s", command, exitErr.ExitCode(), stdoutBytes.String(), stderrBytes.String())
	} else {
		err = fmt.Errorf("command '%s' failed: %w", command, err)
	}
	return "", err
}
