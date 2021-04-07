package terraform

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
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

	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	errBuf := &bytes.Buffer{}
	if exitCode := texec.Apply(dir, args, lpDebug, io.MultiWriter(errBuf, lpError)); exitCode != 0 {
		return sf, errors.Wrap(Diagnose(errBuf.String()), "failed to apply Terraform")
	}
	return sf, nil
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

	err = data.Unpack(filepath.Join(dir, "terraform.rc"), "terraform.rc")
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

	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	os.Setenv("TF_CLI_CONFIG_FILE", filepath.Join(dir, "terraform.rc"))
	os.Setenv("TERRAFORM_LOCK_FILE_PATH", dir)

	args := []string{
		fmt.Sprintf("-plugin-dir=%s", filepath.Join(dir, "plugins")),
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

	re := regexp.MustCompile(`^terraform-provider-(.+)$`)
	pdir := filepath.Join(dir, "plugins", "openshift", "local")

	for name := range plugins.KnownPlugins {
		matches := re.FindStringSubmatch(name)
		if matches == nil {
			logrus.Warnf("Failed to extract plugin name from %s", name)
			continue
		}
		pluginName := matches[1]

		// XXX: HACK: pretend all plugin versions are v1.0.0
		pluginVersion := "1.0.0"
		dstDir := filepath.Join(pdir, pluginName, pluginVersion, fmt.Sprintf("linux_%s", runtime.GOARCH))
		if err := os.MkdirAll(dstDir, 0777); err != nil {
			return err
		}

		dst := filepath.Join(dstDir, name)

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
