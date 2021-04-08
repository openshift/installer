package terraform

import (
	"bytes"
	"fmt"
	"io"
	"os"
	"path/filepath"
	"regexp"
	"runtime"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/data"
	"github.com/openshift/installer/pkg/lineprinter"
	texec "github.com/openshift/installer/pkg/terraform/exec"
	"github.com/openshift/installer/pkg/terraform/exec/plugins"
)

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.  It returns the absolute path of the tfstate file, rooted
// in the specified directory, along with any errors from Terraform.
func Apply(dir string, platform string, stage Stage, extraOpts ...tfexec.ApplyOption) (path string, err error) {
	err = unpackAndInit(dir, platform, stage.Name())
	if err != nil {
		return "", err
	}

	sf := filepath.Join(dir, stage.StateFilename())
	defaultOpts := []tfexec.ApplyOption{
		tfexec.State(sf),
		tfexec.StateOut(sf),
	}
	opts := append(defaultOpts, extraOpts...)

	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	errBuf := &bytes.Buffer{}
	if exitCode := texec.Apply(dir, lpDebug, io.MultiWriter(errBuf, lpError), opts...); exitCode != 0 {
		return sf, errors.Wrap(Diagnose(errBuf.String()), "failed to apply Terraform")
	}
	return sf, nil
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, stage Stage, extraOpts ...tfexec.DestroyOption) (err error) {
	err = unpackAndInit(dir, platform, stage.Name())
	if err != nil {
		return err
	}

	sf := filepath.Join(dir, stage.StateFilename())
	defaultOpts := []tfexec.DestroyOption{
		tfexec.State(sf),
		tfexec.StateOut(sf),
	}
	opts := append(defaultOpts, extraOpts...)

	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	if exitCode := texec.Destroy(dir, lpDebug, lpError, opts...); exitCode != 0 {
		return errors.New("failed to destroy using Terraform")
	}
	return nil
}

// unpack unpacks the platform-specific Terraform modules into the
// given directory.
func unpack(dir string, platform string, target string) (err error) {
	err = data.Unpack(dir, filepath.Join(platform, target))
	if err != nil {
		return err
	}

	err = data.Unpack(filepath.Join(dir, "config.tf"), "config.tf")
	if err != nil {
		return err
	}

	platformVarFile := fmt.Sprintf("variables-%s.tf", platform)

	err = data.Unpack(filepath.Join(dir, platformVarFile), filepath.Join(platform, platformVarFile))
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
func unpackAndInit(dir string, platform string, target string) (err error) {
	err = unpack(dir, platform, target)
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

	if exitCode := texec.Init(dir, lpDebug, lpError, tfexec.PluginDir(filepath.Join(dir, "plugins"))); exitCode != 0 {
		return errors.New("failed to initialize Terraform")
	}
	return nil
}

func setupEmbeddedPlugins(dir string) error {
	re := regexp.MustCompile(`^terraform-provider-(.+)$`)
	pdir := filepath.Join(dir, "plugins", "openshift", "local")

	for name, pluginPath := range plugins.KnownPlugins {
		dst := filepath.Join(pdir, name)
		matches := re.FindStringSubmatch(name)
		if matches == nil {
			logrus.Warnf("Failed to extract plugin name from %s", name)
			continue
		}
		pluginName := matches[1]

		// XXX: HACK: pretend all plugin versions are v1.0.0
		pluginVersion := "1.0.0"
		dstDir := filepath.Join(pdir, pluginName, pluginVersion, fmt.Sprintf("%s_%s", runtime.GOOS, runtime.GOARCH))
		if err := os.MkdirAll(dstDir, 0777); err != nil {
			return err
		}

		dst = filepath.Join(dstDir, name)

		if runtime.GOOS == "windows" {
			dst = fmt.Sprintf("%s.exe", dst)
		}
		if _, err := os.Stat(dst); err == nil {
			// stat succeeded, the plugin already exists.
			continue
		}
		logrus.Debugf("Symlinking plugin %s src: %q dst: %q", name, pluginPath, dst)
		if err := os.Symlink(pluginPath, dst); err != nil {
			return err
		}
	}

	return nil
}
