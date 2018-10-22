package main

import (
	"fmt"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/asset"
	"github.com/openshift/installer/pkg/asset/cluster"
	"github.com/openshift/installer/pkg/asset/ignition/bootstrap"
	"github.com/openshift/installer/pkg/asset/ignition/machine"
	"github.com/openshift/installer/pkg/asset/installconfig"
	"github.com/openshift/installer/pkg/asset/kubeconfig"
	"github.com/openshift/installer/pkg/asset/manifests"
)

type target struct {
	name    string
	command *cobra.Command
	assets  []asset.WritableAsset
}

var targets = []target{{
	name: "Install Config",
	command: &cobra.Command{
		Use:   "install-config",
		Short: "Generates the Install Config asset",
		Long:  "",
	},
	assets: []asset.WritableAsset{&installconfig.InstallConfig{}},
}, {
	name: "Manifests",
	command: &cobra.Command{
		Use:   "manifests",
		Short: "Generates the Kubernetes manifests",
		Long:  "",
	},
	assets: []asset.WritableAsset{&manifests.Manifests{}, &manifests.Tectonic{}},
}, {
	name: "Ignition Configs",
	command: &cobra.Command{
		Use:   "ignition-configs",
		Short: "Generates the Ignition Config asset",
		Long:  "",
	},
	assets: []asset.WritableAsset{&bootstrap.Bootstrap{}, &machine.Master{}, &machine.Worker{}},
}, {
	name: "Cluster",
	command: &cobra.Command{
		Use:   "cluster",
		Short: "Create an OpenShift cluster",
		Long:  "",
	},
	assets: []asset.WritableAsset{&cluster.TerraformVariables{}, &kubeconfig.Admin{}, &cluster.Cluster{}},
}}

// Deprecated: Use 'create' subcommands instead.
func newTargetsCmd() []*cobra.Command {
	var cmds []*cobra.Command
	for _, t := range targets {
		cmd := *t.command
		cmd.Short = fmt.Sprintf("DEPRECATED: USE 'create %s' instead.", cmd.Use)
		cmd.RunE = runTargetCmd(t.assets...)
		cmds = append(cmds, &cmd)
	}
	return cmds
}

func newCreateCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "create",
		Short: "Create part of an OpenShift cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	for _, t := range targets {
		t.command.RunE = runTargetCmd(t.assets...)
		cmd.AddCommand(t.command)
	}

	return cmd
}

func runTargetCmd(targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		assetStore, err := asset.NewStore(rootOpts.dir)
		if err != nil {
			return errors.Wrapf(err, "failed to create asset store")
		}

		for _, a := range targets {
			err := assetStore.Fetch(a)
			if err != nil {
				if exitError, ok := errors.Cause(err).(*exec.ExitError); ok && len(exitError.Stderr) > 0 {
					logrus.Error(strings.Trim(string(exitError.Stderr), "\n"))
				}
				err = errors.Wrapf(err, "failed to fetch %s", a.Name())
			}

			if err2 := asset.PersistToFile(a, rootOpts.dir); err2 != nil {
				err2 = errors.Wrapf(err2, "failed to write asset (%s) to disk", a.Name())
				if err != nil {
					logrus.Error(err2)
					return err
				}
				return err2
			}

			if err != nil {
				return err
			}
		}

		if err := assetStore.Save(rootOpts.dir); err != nil {
			return errors.Wrapf(err, "failed to write to state file")
		}

		if err := assetStore.Purge(targets); err != nil {
			return errors.Wrapf(err, "failed to delete existing on-disk files")
		}

		return nil
	}
}
