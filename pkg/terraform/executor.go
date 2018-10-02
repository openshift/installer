package terraform

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"

	"github.com/sirupsen/logrus"
)

// executor enables calling Terraform from Go, across platforms, with any
// additional providers/provisioners that the currently executing binary
// exposes.
//
// The Terraform binary is expected to be in the executing binary's folder, in
// the current working directory or in the PATH.
type executor struct {
	binaryPath string
}

// Set the binary names for different platforms
const (
	tfBinUnix    = "terraform"
	tfBinWindows = "terraform.exe"
)

// errBinaryNotFound denotes the fact that the Terraform binary could not be
// found on disk.
var errBinaryNotFound = errors.New(
	"terraform not in executable's folder, cwd nor PATH",
)

// newExecutor initializes a new Executor.
func newExecutor() (*executor, error) {
	ex := new(executor)

	// Find the Terraform binary.
	binPath, err := tfBinaryPath()
	if err != nil {
		return nil, err
	}

	ex.binaryPath = binPath
	return ex, nil
}

// Execute runs the given command and arguments against Terraform.
//
// An error is returned if the Terraform binary could not be found, or if the
// Terraform call itself failed, in which case, details can be found in the
// output.
func (ex *executor) execute(clusterDir string, args ...string) error {
	// Prepare Terraform command by setting up the command, configuration,
	// and the working directory
	if clusterDir == "" {
		return fmt.Errorf("clusterDir is unset. Quitting")
	}

	cmd := exec.Command(ex.binaryPath, args...)
	cmd.Dir = clusterDir
	if logrus.GetLevel() == logrus.DebugLevel {
		cmd.Stdin = os.Stdin
		cmd.Stdout = os.Stdout
		cmd.Stderr = os.Stderr
	}

	logrus.Debugf("Running %#v...", cmd)

	// Start Terraform.
	return cmd.Run()
}

// tfBinaryPath searches for a Terraform binary on disk:
// - in the executing binary's folder,
// - in the current working directory,
// - in the PATH.
// The first to be found is the one returned.
func tfBinaryPath() (string, error) {
	// Depending on the platform, the expected binary name is different.
	binaryFileName := tfBinUnix
	if runtime.GOOS == "windows" {
		binaryFileName = tfBinWindows
	}

	// Find the current executable's path, gets an absolute path or error
	execPath, err := os.Executable()
	if err == nil {
		// execPath could be a symlink
		if stat, err := os.Stat(execPath); err == nil && (stat.Mode()&os.ModeSymlink) == os.ModeSymlink {
			if evalExecPath, err := filepath.EvalSymlinks(execPath); err != nil {
				execPath = evalExecPath
			}
		}

		// Look into the executable's folder.
		path := filepath.Join(filepath.Dir(execPath), binaryFileName)
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			return path, nil
		}
	}

	// Look into cwd.
	if workingDirectory, err := os.Getwd(); err == nil {
		path := filepath.Join(workingDirectory, binaryFileName)
		if stat, err := os.Stat(path); err == nil && !stat.IsDir() {
			return path, nil
		}
	}

	// If we still haven't found the executable, look for it
	// in the PATH.
	if path, err := exec.LookPath(binaryFileName); err == nil {
		return filepath.Abs(path)
	}

	return "", errBinaryNotFound
}
