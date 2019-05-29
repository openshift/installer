package main

import (
	"context"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"
)

func newWaitForCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait-for",
		Short: "Wait for install-time events",
		Long: `Wait for install-time events.

'create cluster' has a few stages that wait for cluster events.  But
these waits can also be useful on their own.  This subcommand exposes
them directly.`,
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}
	cmd.AddCommand(newWaitForBootstrapCompleteCmd())
	cmd.AddCommand(newWaitForInstallCompleteCmd())
	return cmd
}

func newWaitForBootstrapCompleteCmd() *cobra.Command {
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
				_, _, err2 := getGatherBootstrapIPs(rootOpts.dir)
				if err2 == nil {
					logrus.Infof("Gather logs from the cluster with: %q --dir %q gather bootstrap", os.Args[0], rootOpts.dir)
				} else {
					logrus.Infof("Gather logs from the cluster with: %q gather bootstrap --bootstrap PUBLIC_BOOTSTRAP_IP [--master MASTER_IP ...]", os.Args[0])
					logrus.Infof("See also: %q gather bootstrap --help", os.Args[0])
				}
				logrus.Fatal(err)
			}

			logrus.Info("It is now safe to remove the bootstrap resources")
		},
	}
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster is ready",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			ctx := context.Background()

			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}

			err = waitForInstallComplete(ctx, config, rootOpts.dir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}
