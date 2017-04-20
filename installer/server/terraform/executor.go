package terraform

import (
	"bufio"
	"errors"
	"fmt"
	"io"
	"io/ioutil"
	"os"
	"os/exec"
	"path/filepath"
	"runtime"
	"strings"

	"github.com/hashicorp/terraform/terraform"
	"github.com/kardianos/osext"
	"github.com/shirou/gopsutil/process"
)

const (
	configFileName = ".terraformrc"
	stateFileName  = "terraform.tfstate"
	tfVarsFileName = "terraform.tfvars"
	logsFolderName = "logs"

	logsFileSuffix = ".log"
	failFileSuffix = ".fail"
)

// ErrBinaryNotFound denotes the fact that the TerraForm binary could not be
// found on disk.
var ErrBinaryNotFound = errors.New(
	"TerraForm not in executable's folder, cwd nor PATH",
)

// ExecutionStatus describes whether an execution succeeded, failed or is still
// in progress.
type ExecutionStatus string

const (
	// ExecutionStatusUnknown indicates that the status of execution is unknown.
	ExecutionStatusUnknown ExecutionStatus = "Unknown"
	// ExecutionStatusRunning indicates that the the execution is still in
	// process.
	ExecutionStatusRunning ExecutionStatus = "Running"
	// ExecutionStatusSuccess indicates that the execution succeeded.
	ExecutionStatusSuccess ExecutionStatus = "Success"
	// ExecutionStatusFailure indicates that the execution failed.
	ExecutionStatusFailure ExecutionStatus = "Failure"
)

// Executor enables calling TerraForm from Go, across platforms, with any
// additional providers/provisioners that the currently executing binary
// exposes.
//
// The TerraForm binary is expected to be in the executing binary's folder, in
// the current working directory or in the PATH.
// Each Executor runs in a temporary folder, so each Executor should only be
// used for one TF project.
//
// TODO: Ideally, we would use TerraForm as a Go library, so we can monitor a
// hook and report the current state in real-time when
// Apply/Refresh/Destroy are used. While technically possible today, because
// TerraForm currently hides the providers/provisioners list construction in
// their main package, it would require to reproduce a bunch of their logic,
// which is out of the scope of the first-version of the Executor. With a bit of
// efforts, we could actually even stop requiring having a TerraForm binary
// altogether, by linking the builtin providers/provisioners to this particular
// binary and re-implemeting the routing here. Alternatively, we could
// contribute upstream to add a 'debug' flag that would enable a hook that would
// expose the live state to a file (or else).
type Executor struct {
	executionPath string
	configPath    string
	binaryPath    string
	envVariables  map[string]string
}

// NewExecutor initializes a new Executor.
func NewExecutor() (*Executor, error) {
	// Create a temporary folder in which the new Executor can run.
	executionPath, err := ioutil.TempDir(os.TempDir(), "tectonic")
	if err != nil {
		return nil, err
	}

	// Create an executor in that path.
	ex, err := NewExecutorFromPath(executionPath)
	if err != nil {
		os.RemoveAll(executionPath)
		return nil, err
	}

	return ex, err
}

// NewExecutorFromPath creates an Executor from an existing path.
func NewExecutorFromPath(executionPath string) (*Executor, error) {
	var err error

	ex := new(Executor)
	ex.executionPath = executionPath

	// Create the folder in which the logs will be stored, if not existing.
	os.Mkdir(filepath.Join(ex.executionPath, logsFolderName), 0770)

	// Create a Executor CLI configuration file, that contains the list of
	// vendored providers/provisioners.
	config, err := BuildPluginsConfig()
	if err != nil {
		return nil, err
	}

	ex.configPath = filepath.Join(ex.WorkingDirectory(), configFileName)
	if err := ioutil.WriteFile(ex.configPath, []byte(config), 0660); err != nil {
		return nil, err
	}

	// Find the TerraForm binary.
	ex.binaryPath, err = tfBinaryPath()
	if err != nil {
		return nil, err
	}

	return ex, nil
}

// AddFile is a convenience function that writes a single file in the Executor's
// working directory using the given content. It may replace an existing file.
func (ex *Executor) AddFile(name string, content []byte) error {
	filePath := filepath.Join(ex.WorkingDirectory(), name)
	return ioutil.WriteFile(filePath, content, 0660)
}

// AddVariables writes the `terraform.tfvars` file in the Executor's working
// directory using the given content. It may replace an existing file.
func (ex *Executor) AddVariables(content []byte) error {
	return ex.AddFile(tfVarsFileName, content)
}

// AddEnvironmentVariables adds extra environment variables that will be set
// during the execution.
// Existing variables are replaced. This function is not thread-safe.
func (ex *Executor) AddEnvironmentVariables(envVars map[string]string) {
	if ex.envVariables == nil {
		ex.envVariables = make(map[string]string)
	}
	for k, v := range envVars {
		ex.envVariables[k] = v
	}
}

// AddCredentials is a convenience function that converts the given Credentials
// into environment variables and add them to the Executor.
//
// If the credentials parameter is nil, nothing is done.
// An error is returned if the credentials are invalid.
func (ex *Executor) AddCredentials(credentials *Credentials) error {
	if credentials == nil {
		return nil
	}

	env, err := credentials.ToEnvironment()
	if err != nil {
		return err
	}
	ex.AddEnvironmentVariables(env)

	return nil
}

// Execute runs the given command and arguments against TerraForm, and returns
// an identifier that can be used to read the output of the process as it is
// executed and after.
//
// Execute is non-blocking, and takes a lock in the execution path.
// Locking is handled by TerraForm itself.
//
// An error is returned if the TerraForm binary could not be found, or if the
// TerraForm call itself failed, in which case, details can be found in the
// output.
func (ex *Executor) Execute(args ...string) (int, chan struct{}, error) {
	// Prepare TerraForm command by setting up the command, configuration,
	// working directory (so the files such as terraform.tfstate are stored at
	// the right place), extra environment variables and outputs.
	cmd := exec.Command(ex.binaryPath, args...)
	cmd.Env = append(cmd.Env, fmt.Sprintf("TERRAFORM_CONFIG=%s", ex.configPath))
	for k, v := range ex.envVariables {
		cmd.Env = append(cmd.Env, fmt.Sprintf("%s=%s", strings.ToUpper(k), v))
	}
	cmd.Dir = ex.executionPath

	rPipe, wPipe := io.Pipe()
	cmd.Stdout = wPipe
	cmd.Stderr = wPipe

	// Start TerraForm.
	err := cmd.Start()
	if err != nil {
		// The process failed to start, we can't even save that it started since we
		// don't have a PID yet.
		return -1, nil, err
	}

	// Create a log file and pipe stdout/stderr to it.
	logFile, err := os.Create(ex.logPath(cmd.Process.Pid))
	if err != nil {
		return -1, nil, err
	}
	go io.Copy(logFile, rPipe)

	done := make(chan struct{})
	go func() {
		// Wait for the process to finish.
		if err := cmd.Wait(); err != nil {
			// The process did not end cleanly. Write the failure file.
			ioutil.WriteFile(ex.failPath(cmd.Process.Pid), []byte(err.Error()), 0660)
		}

		// Close descriptors.
		wPipe.Close()
		logFile.Close()
		close(done)
	}()

	return cmd.Process.Pid, done, nil
}

// WorkingDirectory returns the directory in which TerraForm runs, which can be
// useful for inspection or to retrieve any generated files.
func (ex *Executor) WorkingDirectory() string {
	return ex.executionPath
}

// Output returns a ReadCloser on the output file of an execution, or an error
// if no output for that execution identifier can be found.
func (ex *Executor) Output(id int) (io.ReadCloser, error) {
	return os.Open(ex.logPath(id))
}

// Status returns the status of a given execution process.
//
// An error can be returned if the running processes could not be listed, or if
// the process failed, in which case the exit message is returned in an error
// of type ExecutionError.
//
// Note that if the identifier is invalid, the current implementation will
// return ExecutionStatusSuccess rather than ExecutionStatusUnknown.
func (ex *Executor) Status(id int) (ExecutionStatus, error) {
	isRunning, err := process.PidExists(int32(id))
	if err != nil {
		return ExecutionStatusUnknown, err
	}
	if isRunning {
		return ExecutionStatusRunning, nil
	}

	if failErr, err := ioutil.ReadFile(ex.failPath(id)); err == nil {
		return ExecutionStatusFailure, errors.New(string(failErr))
	}
	return ExecutionStatusSuccess, nil
}

// State returns the current TerraForm State.
//
// The returned value can be nil if there is currently no state held.
func (ex *Executor) State() *terraform.State {
	f, err := os.Open(filepath.Join(ex.executionPath, stateFileName))
	if err != nil {
		return nil
	}
	defer f.Close()

	s, err := terraform.ReadState(bufio.NewReader(f))
	if err != nil {
		return nil
	}

	return s
}

// Cleanup removes resources that were allocated by the Executor.
func (ex *Executor) Cleanup() {
	if ex.executionPath != "" {
		os.RemoveAll(ex.executionPath)
	}
}

// tfBinatyPath searches for a TerraForm binary on disk:
// - in the executing binary's folder,
// - in the current working directory,
// - in the PATH.
// The first to be found is the one returned.
func tfBinaryPath() (string, error) {
	// Depending on the platform, the expected binary name is different.
	binaryFileName := "terraform"
	if runtime.GOOS == "windows" {
		binaryFileName = "terraform.exe"
	}

	// Look into the executable's folder.
	if execFolderPath, err := osext.ExecutableFolder(); err == nil {
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

	return "", ErrBinaryNotFound
}

// failPath returns the path to the failure file of a given execution process.
func (ex *Executor) failPath(id int) string {
	failFileName := fmt.Sprintf("%d%s", id, failFileSuffix)
	return filepath.Join(ex.executionPath, logsFolderName, failFileName)
}

// logPath returns the path to the log file of a given execution process.
func (ex *Executor) logPath(id int) string {
	logFileName := fmt.Sprintf("%d%s", id, logsFileSuffix)
	return filepath.Join(ex.executionPath, logsFolderName, logFileName)
}
