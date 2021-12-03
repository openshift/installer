package terraform

import (
	"context"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/lineprinter"
)

// newTFExec creates a tfexec.Terraform for executing Terraform CLI commands.
// The `datadir` is the location where the terraform parts (binaries, tf files, etc) have been unpacked to.
// The stdout and stderr will be sent to the logger at the debug and error levels,
// respectively.
func newTFExec(datadir string) (*tfexec.Terraform, error) {
	tfPath := filepath.Join(datadir, "bin", "terraform")
	tf, err := tfexec.NewTerraform(datadir, tfPath)
	if err != nil {
		return nil, err
	}

	lpDebug := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Debug}).Print}
	lpError := &lineprinter.LinePrinter{Print: (&lineprinter.Trimmer{WrappedPrint: logrus.Error}).Print}
	defer lpDebug.Close()
	defer lpError.Close()

	tf.SetStdout(lpDebug)
	tf.SetStderr(lpError)
	tf.SetLogger(newPrintfer())

	return tf, nil
}

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.  It returns the absolute path of the tfstate file, rooted
// in the specified directory, along with any errors from Terraform.
func Apply(dir string, platform string, stage Stage, extraOpts ...tfexec.ApplyOption) (path string, err error) {
	err = unpackAndInit(dir, platform, stage.Name(), stage.Providers())
	if err != nil {
		return "", err
	}

	sf := filepath.Join(dir, stage.StateFilename())
	opts := append(
		extraOpts,
		tfexec.State(sf),
		tfexec.StateOut(sf),
	)

	tf, err := newTFExec(dir)
	if err != nil {
		return "", errors.Wrap(err, "failed to create a new tfexec")
	}
	err = tf.Apply(context.Background(), opts...)
	err = diagnoseApplyError(err)
	return sf, errors.Wrap(err, "failed to apply Terraform")
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, stage Stage, extraOpts ...tfexec.DestroyOption) (err error) {
	err = unpackAndInit(dir, platform, stage.Name(), stage.Providers())
	if err != nil {
		return err
	}

	sf := filepath.Join(dir, stage.StateFilename())
	opts := append(
		extraOpts,
		tfexec.State(sf),
		tfexec.StateOut(sf),
	)

	tf, err := newTFExec(dir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}
	return errors.Wrap(
		tf.Destroy(context.Background(), opts...),
		"failed doing terraform destroy",
	)
}
