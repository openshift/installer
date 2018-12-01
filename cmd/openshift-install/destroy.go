package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/destroy"
	"github.com/openshift/installer/pkg/destroy/bootstrap"
	_ "github.com/openshift/installer/pkg/destroy/libvirt"
	_ "github.com/openshift/installer/pkg/destroy/openstack"
)

func newDestroyCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "destroy",
		Short: "Destroy part of an OpenShift cluster",
		Long:  "",
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newDestroyBootstrapCmd())
	cmd.AddCommand(newDestroyClusterCmd())
	return cmd
}

func newDestroyClusterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster",
		Short: "Destroy an OpenShift cluster",
		RunE:  runDestroyCmd,
	}
}

func runDestroyCmd(cmd *cobra.Command, args []string) error {
	cleanup, err := setupFileHook(rootOpts.dir)
	if err != nil {
		return errors.Wrap(err, "failed to setup logging hook")
	}
	defer cleanup()

	destroyer, err := destroy.New(logrus.StandardLogger(), rootOpts.dir)
	if err != nil {
		return errors.Wrap(err, "Failed while preparing to destroy cluster")
	}
	if err := destroyer.Run(); err != nil {
		return errors.Wrap(err, "Failed to destroy cluster")
	}
	return nil
}

func newDestroyBootstrapCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap",
		Short: "Destroy the bootstrap resources",
		RunE: func(cmd *cobra.Command, args []string) error {
			return bootstrap.Destroy(rootOpts.dir)
		},
	}
}
