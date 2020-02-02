package main

// This command is necessary for Terraform's internal plugins (remote-exec, local-exec, ...) to work.
// If they are used, Terraform will fork it's binary and call it with three parameters: 
//     'internal-plugin' <TYPE: e.g. provisioner> <NAME: e.g. remote-exec>
// As soon as the fork is running in the background, Terraform connects to it through RPC.

// This command is hidden in the help menu.

import (
	"os"
	"io/ioutil"
	"github.com/pkg/errors"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/terraform"	
)

func internalPluginCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "internal-plugin",
		Short: "Run internal Terraform commands (local-exec, remote-exec, ...).",
		Long:  "",
		Args:  cobra.ExactArgs(2),
		RunE:  runInternalPluginCmd,
		Hidden: true,
	}
}

func runInternalPluginCmd(cmd *cobra.Command, args []string) error {
	// Copy the terraform.tfvars to a temp directory where the terraform will be invoked within.
	tmpDir, err := ioutil.TempDir("", "openshift-install-")
	if err != nil {
		return errors.Wrap(err, "failed to create temp dir for terraform internal-plugin execution")
	}
	defer os.RemoveAll(tmpDir)

	err = terraform.InternalPlugin(tmpDir, args[0], args[1])
	if err != nil {
		return err
	}

	return nil
}
