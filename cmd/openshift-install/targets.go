package main

import (
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
	"github.com/openshift/installer/pkg/asset/metadata"
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
	assets: []asset.WritableAsset{&cluster.TerraformVariables{}, &kubeconfig.Admin{}, &cluster.Cluster{}, &metadata.Metadata{}},
}}

func newTargetsCmd() []*cobra.Command {
	var cmds []*cobra.Command
	for _, t := range targets {
		t.command.RunE = runTargetCmd(t.assets...)
		cmds = append(cmds, t.command)
	}
	return cmds
}

func runTargetCmd(targets ...asset.WritableAsset) func(cmd *cobra.Command, args []string) error {
	return func(cmd *cobra.Command, args []string) error {
		assetStore := &asset.StoreImpl{}
		for _, a := range targets {
			err := assetStore.Fetch(a)
			if err != nil {
				if exitError, ok := errors.Cause(err).(*exec.ExitError); ok && len(exitError.Stderr) > 0 {
					logrus.Error(strings.Trim(string(exitError.Stderr), "\n"))
				}
				return errors.Wrapf(err, "failed to generate %s", a.Name())
			}

			if err := asset.PersistToFile(a, rootOpts.dir); err != nil {
				return errors.Wrapf(err, "failed to write asset (%s) to disk", a.Name())
			}
		}
		return nil
	}
}
