package workflow

import (
	"errors"
	"fmt"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
)

// executor enables calling TerraForm from Go, across platforms, with any
// additional providers/provisioners that the currently executing binary
// exposes.
//
// The TerraForm binary is expected to be in the executing binary's folder, in
// the current working directory or in the PATH.
type executor struct {
	binaryPath string
}

// Set the binary names for different platforms
const (
	tfBinUnix    = "terraform"
	tfBinWindows = "terraform.exe"
)

// errBinaryNotFound denotes the fact that the TerraForm binary could not be
// found on disk.
var errBinaryNotFound = errors.New(
	"TerraForm not in executable's folder, cwd nor PATH",
)

// newExecutor initializes a new Executor.
func newExecutor() (*executor, error) {
	ex := new(executor)

	// Find the TerraForm binary.
	binPath, err := tfBinaryPath()
	if err != nil {
		return nil, err
	}

	ex.binaryPath = binPath
	return ex, nil
}

// Execute runs the given command and arguments against TerraForm.
//
// An error is returned if the TerraForm binary could not be found, or if the
// TerraForm call itself failed, in which case, details can be found in the
// output.
func (ex *executor) execute(clusterDir string, args ...string) error {
	// Prepare TerraForm command by setting up the command, configuration,
	// and the working directory
	if clusterDir == "" {
		return fmt.Errorf("clusterDir is unset. Quitting")
	}

	cmd := exec.Command(ex.binaryPath, args...)
	cmd.Stdin = os.Stdin
	cmd.Stdout = os.Stdout
	cmd.Stderr = os.Stderr
	cmd.Dir = clusterDir

	// Start TerraForm.
	return cmd.Run()
}

// tfBinatyPath searches for a TerraForm binary on disk:
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

	// Look into the executable's folder.
	if execFolderPath, err := filepath.Abs(filepath.Dir(os.Args[0])); err == nil {
		path := filepath.Join(execFolderPath, binaryFileName)
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
