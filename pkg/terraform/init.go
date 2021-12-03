package terraform

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/hashicorp/terraform-exec/tfexec"
	"github.com/pkg/errors"

	"github.com/openshift/installer/data"
	prov "github.com/openshift/installer/pkg/terraform/providers"
)

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
func unpackAndInit(dir string, platform string, target string, providers []prov.Provider) (err error) {
	err = unpack(dir, platform, target)
	if err != nil {
		return errors.Wrap(err, "failed to unpack Terraform modules")
	}

	if err := setupEmbeddedPlugins(dir, providers); err != nil {
		return errors.Wrap(err, "failed to setup embedded Terraform plugins")
	}

	tf, err := newTFExec(dir)
	if err != nil {
		return errors.Wrap(err, "failed to create a new tfexec")
	}

	// Explicitly specify the CLI config file to use so that we control the providers that are used.
	tf.SetEnv(map[string]string{
		"TF_CLI_CONFIG_FILE": filepath.Join(dir, "terraform.rc"),
	})

	return errors.Wrap(
		tf.Init(context.Background(), tfexec.PluginDir(filepath.Join(dir, "plugins"))),
		"failed doing terraform init",
	)
}

func setupEmbeddedPlugins(dir string, providers []prov.Provider) error {
	// Unpack the terraform binary.
	if err := prov.UnpackTerraformBinary(filepath.Join(dir, "bin")); err != nil {
		return err
	}

	// Unpack the providers.
	for _, provider := range providers {
		if err := provider.Extract(filepath.Join(dir, "plugins")); err != nil {
			return err
		}
	}

	return nil
}
