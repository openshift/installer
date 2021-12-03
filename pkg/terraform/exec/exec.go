// Package exec is glue between the vendored terraform codebase and installer.
package exec

import (
	"context"
	"fmt"
	"github.com/hashicorp/terraform-exec/tfexec"
	"io"
	"os"
	"path/filepath"
)

func getPluginPath() string {
	userCacheDir, _ := os.UserCacheDir()
	return filepath.Join(userCacheDir, "openshift-installer", "terraform")
}

func getPluginBinPath() string {
	return filepath.Join(getPluginPath(), "bin")
}

func getTerraformPath() string {
	terraformPath := filepath.Join(getPluginBinPath(), "terraform")
	_, err := os.Stat(terraformPath)
	if err != nil {
		panic(fmt.Sprintf("Failed to find terraform: %s", err))
	}

	return terraformPath
}

// Apply is wrapper around `terraform apply` subcommand.
func Apply(datadir string, stdout, stderr io.Writer, opts ...tfexec.ApplyOption) int {
	tfPath := getTerraformPath()
	tf, err := tfexec.NewTerraform(datadir, tfPath)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: new terraform: %s\n", err)
	}

	tf.SetStdout(stdout)
	tf.SetStderr(stderr)

	err = tf.Apply(context.Background(), opts...)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: terraform apply: %s\n", err)
		return 1
	}

	return 0
}

// Destroy is wrapper around `terraform destroy` subcommand.
func Destroy(datadir string, stdout, stderr io.Writer, opts ...tfexec.DestroyOption) int {
	tfPath := getTerraformPath()
	tf, err := tfexec.NewTerraform(datadir, tfPath)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: new terraform: %s\n", err)
	}

	tf.SetStdout(stdout)
	tf.SetStderr(stderr)

	err = tf.Destroy(context.Background(), opts...)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: terraform destroy: %s\n", err)
		return 1
	}

	return 0
}

// Init is wrapper around `terraform init` subcommand.
func Init(datadir string, stdout, stderr io.Writer, opts ...tfexec.InitOption) int {
	tfPath := getTerraformPath()
	tf, err := tfexec.NewTerraform(datadir, tfPath)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: new terraform: %s\n", err)
	}

	tf.SetStdout(stdout)
	tf.SetStderr(stderr)

	err = tf.Init(context.Background(), opts...)
	if err != nil {
		fmt.Fprintf(stderr, "Failed: terraform init: %s\n", err)
		return 1
	}

	return 0
}
