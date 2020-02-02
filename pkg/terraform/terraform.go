package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/lineprinter"
	texec "github.com/openshift/installer/pkg/terraform/exec"
	"github.com/openshift/installer/pkg/terraform/exec/plugins"
)

const (
	// StateFileName is the default name for Terraform state files.
	StateFileName string = "terraform.tfstate"

	// VarFileName is the default name for Terraform var file.
	VarFileName string = "terraform.tfvars"
)

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.  It returns the absolute path of the tfstate file, rooted
// in the specified directory, along with any errors from Terraform.
func Apply(dir string, platform string, extraArgs ...string) (path string, err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return "", err
	}

	defaultArgs := []string{
		"-auto-approve",
		"-input=false",
		fmt.Sprintf("-state=%s", filepath.Join(dir, StateFileName)),
		fmt.Sprintf("-state-out=%s", filepath.Join(dir, StateFileName)),
	}
	args := append(defaultArgs, extraArgs...)
	args = append(args, dir)
	sf := filepath.Join(dir, StateFileName)

	tDebug := &lineprinter.Trimmer{WrappedPrint: logrus.Debug}
	tError := &lineprinter.Trimmer{WrappedPrint: logrus.Error}
	lpDebug := &lineprinter.LinePrinter{Print: tDebug.Print}
	lpError := &lineprinter.LinePrinter{Print: tError.Print}
	defer lpDebug.Close()
	defer lpError.Close()

	if exitCode := texec.Apply(dir, args, lpDebug, lpError); exitCode != 0 {
		return sf, errors.New("failed to apply using Terraform")
	}
	return sf, nil
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, extraArgs ...string) (err error) {
	err = unpackAndInit(dir, platform)
	if err != nil {
		return err
	}

	defaultArgs := []string{
		"-auto-approve",
		"-input=false",
		fmt.Sprintf("-state=%s", filepath.Join(dir, StateFileName)),
		fmt.Sprintf("-state-out=%s", filepath.Join(dir, StateFileName)),
	}
	args := append(defaultArgs, extraArgs...)
	args = append(args, dir)

	tDebug := &lineprinter.Trimmer{WrappedPrint: logrus.Debug}
	tError := &lineprinter.Trimmer{WrappedPrint: logrus.Error}
	lpDebug := &lineprinter.LinePrinter{Print: tDebug.Print}
	lpError := &lineprinter.LinePrinter{Print: tError.Print}
	defer lpDebug.Close()
	defer lpError.Close()

	if exitCode := texec.Destroy(dir, args, lpDebug, lpError); exitCode != 0 {
		return errors.New("failed to destroy using Terraform")
	}
	return nil
}

// unpack unpacks the platform-specific Terraform modules into the
// given directory.
func unpack(dir string, platform string) (err error) {
	err = data.Unpack(dir, platform)
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "config.tf"), "config.tf")
	if err != nil {
		return err
	}

	return nil
}

// unpackAndInit unpacks the platform-specific Terraform modules into
// the given directory and then runs 'terraform init'.
func unpackAndInit(dir string, platform string) (err error) {
	err = unpack(dir, platform)
	if err != nil {
		return errors.Wrap(err, "failed to unpack Terraform modules")
	}

	if err := setupEmbeddedPlugins(dir); err != nil {
		return errors.Wrap(err, "failed to setup embedded Terraform plugins")
	}

	tDebug := &lineprinter.Trimmer{WrappedPrint: logrus.Debug}
	tError := &lineprinter.Trimmer{WrappedPrint: logrus.Error}
	lpDebug := &lineprinter.LinePrinter{Print: tDebug.Print}
	lpError := &lineprinter.LinePrinter{Print: tError.Print}
	defer lpDebug.Close()
	defer lpError.Close()

	args := []string{
		"-get-plugins=true",  // TODO: template plugin is required. Handle this the same as the other plugins.
	}
	args = append(args, dir)
	if exitCode := texec.Init(dir, args, lpDebug, lpError); exitCode != 0 {
		return errors.New("failed to initialize Terraform")
	}
	return nil
}

func setupEmbeddedPlugins(dir string) error {
	execPath, err := os.Executable()
	if err != nil {
		return errors.Wrap(err, "failed to find path for the executable")
	}

	pdir := filepath.Join(dir, "plugins")
	if err := os.MkdirAll(pdir, 0777); err != nil {
		return err
	}
	for name := range plugins.KnownPlugins {
		dst := filepath.Join(pdir, name)
		if runtime.GOOS == "windows" {
			dst = fmt.Sprintf("%s.exe", dst)
		}
		if _, err := os.Stat(dst); err == nil {
			// stat succeeded, the plugin already exists.
			continue
		}
		logrus.Debugf("Symlinking plugin %s src: %q dst: %q", name, execPath, dst)
		if err := os.Symlink(execPath, dst); err != nil {
			return err
		}
	}
	return nil
}

// InternalPlugin handles internal Terraform plugins like the local-exec / remote-exec provisioners.
func InternalPlugin(dir string, extraArgs ...string) (err error) {
	defaultArgs := []string{}
	args := append(defaultArgs, extraArgs...)

	tDebug := &lineprinter.Trimmer{WrappedPrint: logrus.Debug}
	tError := &lineprinter.Trimmer{WrappedPrint: logrus.Error}
	lpDebug := &lineprinter.LinePrinter{Print: tDebug.Print}
	lpError := &lineprinter.LinePrinter{Print: tError.Print}
	defer lpDebug.Close()
	defer lpError.Close()

	if exitCode := texec.InternalPlugin(dir, args, lpDebug, lpError); exitCode != 0 {
		return errors.New("failed to call internal-plugin using Terraform")
	}
	return nil
}