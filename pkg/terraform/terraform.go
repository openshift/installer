package terraform

import (
	"context"
	"os"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"

	"github.com/openshift/installer/pkg/lineprinter"
)

// newTFExec creates a tfexec.Terraform for executing Terraform CLI commands.
// The `datadir` is the location to which the terraform plan (tf files, etc) has been unpacked.
// The `terraformDir` is the location to which the terraform and provider binaries have been unpacked.
// The stdout and stderr will be sent to the logger at the debug and error levels,
// respectively.
func newTFExec(datadir string, terraformDir string) (*tfexec.Terraform, error) {
	tfPath := filepath.Join(terraformDir, "bin", "terraform")
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

	// Opt in to the beta of v3.0 of the azurerm terraform provider.
	// Note that it would have been nice to use tf.SetEnv. Unfortunately, that prevents the environment variables from
	// the OS from being used.
	os.Setenv("ARM_THREEPOINTZERO_BETA", "true")

	return tf, nil
}

// Apply unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// apply'.
func Apply(dir string, platform string, stage Stage, terraformDir string, extraOpts ...tfexec.ApplyOption) error {
	if err := unpackAndInit(dir, platform, stage.Name(), terraformDir, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}
	err = tf.Apply(context.Background(), extraOpts...)
	return errors.Wrap(diagnoseApplyError(err), "failed to apply Terraform")
}

// Destroy unpacks the platform-specific Terraform modules into the
// given directory and then runs 'terraform init' and 'terraform
// destroy'.
func Destroy(dir string, platform string, stage Stage, terraformDir string, extraOpts ...tfexec.DestroyOption) error {
	if err := unpackAndInit(dir, platform, stage.Name(), terraformDir, stage.Providers()); err != nil {
		return err
	}

	tf, err := newTFExec(dir, terraformDir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}
	return errors.Wrap(
		tf.Destroy(context.Background(), extraOpts...),
		"failed doing terraform destroy",
	)
}
