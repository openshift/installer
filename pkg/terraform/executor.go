package terraform

import (
	"bytes"
	"os/exec"
	"path/filepath"

	"github.com/openshift/installer/pkg/lineprinter"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
)

// Execute runs the given command and arguments against Terraform.
//
// An error is returned if the Terraform binary could not be found, or if the
// Terraform call itself failed, in which case, details can be found in the
// output.
func Execute(clusterDir string, args ...string) error {
	// Prepare Terraform command by setting up the command, configuration,
	// and the working directory
	if clusterDir == "" {
		return errors.Errorf("clusterDir is unset. Quitting")
	}

	trimmer := &lineprinter.Trimmer{WrappedPrint: logrus.StandardLogger().Debug}
	linePrinter := &lineprinter.LinePrinter{Print: trimmer.Print}
	defer linePrinter.Close()

	stderr := &bytes.Buffer{}

	cmd := exec.Command(filepath.Join(clusterDir, executablePath), args...)
	cmd.Dir = clusterDir
	cmd.Stdout = linePrinter
	cmd.Stderr = stderr

	logrus.Debugf("Running %#v...", cmd)
	err := cmd.Run()
	if err != nil {
		exitError := err.(*exec.ExitError)
		exitError.Stderr = stderr.Bytes()
	}
	return err
}
