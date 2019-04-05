package main

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

var (
	upiLong = `Entry-points for user-provided infrastructure.

Most users will want to use 'create cluster' to have the installer
create the required infrastructure for their cluster.  But in some
installations the infrastructure needs to be adapted in ways that
installer-created infrastructure does not support.  This command
provides entry points to support the following workflow:

1. Call 'create ignition-configs' to create the bootstrap Ignition
   config and admin kubeconfig.
2. Creates all required cluster resources, after which the cluster
   will being bootstrapping.
3. Call 'user-provided-infrastructure bootstrap-complete' to wait
   until the bootstrap phase has completed.
4. Destroy the bootstrap resources.
5. Call 'user-provided-infrastructure finish' to wait until the
   cluster finishes deploying its initial version.`
)

func newUPICmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:     "user-provided-infrastructure",
		Aliases: []string{"upi"},
		Short:   "Entry-points for user-provided infrastructure",
		Long:    upiLong,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newUPIBootstrapCompleteCmd())
	cmd.AddCommand(newUPIFinishCmd())
	return cmd
}

func newUPIBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until cluster bootstrapping has completed",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			ctx := context.Background()

			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}

			err = waitForBootstrapComplete(ctx, config, rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}

			logrus.Info("It is now safe to remove the bootstrap resources")
		},
	}
}

func newUPIFinishCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "finish",
		Short: "Wait for the cluster to finish updating",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}

			err = finish(ctx, config, rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}
