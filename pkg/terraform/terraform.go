package terraform

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"
	"time"

	"github.com/gosuri/uiprogress"
	"github.com/gosuri/uiprogress/util/strutil"
	"github.com/openshift/installer/data"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"golang.org/x/crypto/ssh/terminal"

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
		fmt.Sprintf("-var-file=%s", filepath.Join(dir, VarFileName)),
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

	progressCh := make(chan texec.Progress)
	defer close(progressCh)
	if terminal.IsTerminal(int(os.Stdout.Fd())) {
		go reportProgress(progressCh)
	}

	if exitCode := texec.ApplyWithProgress(dir, args, lpDebug, lpError, progressCh); exitCode != 0 {
		return sf, errors.New("failed to apply using Terraform")
	}
	return sf, nil
}

func reportProgress(progressCh <-chan texec.Progress) {
	uiprogress.Start()
	var bar *uiprogress.Bar
	started := time.Now()
	elapsedPrepend := func(_ *uiprogress.Bar) string {
		return strutil.PadLeft(strutil.PrettyTime(time.Since(started)), 5, ' ')
	}
	for p := range progressCh {
		// add the bar when total becomes non-zero
		if bar == nil && p.Total > 0 {
			bar = uiprogress.AddBar(p.Total)
			bar.TimeStarted = time.Now()
			bar.AppendCompleted()
			bar.PrependFunc(elapsedPrepend)
		}
		// update bar when it has been initialized
		if bar != nil {
			done := p.Added + p.Changed + p.Removed
			bar.Total = p.Total
			bar.Set(done)
		}
	}
	uiprogress.Stop()
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
		fmt.Sprintf("-var-file=%s", filepath.Join(dir, VarFileName)),
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
		"-get-plugins=false",
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
