package main

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/tools/clientcmd"

	"github.com/openshift/installer/cmd/openshift-install/command"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	timer "github.com/openshift/installer/pkg/metrics/timer"
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
			timer.StartTimer(timer.TotalTimeElapsed)
			ctx := context.Background()

			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(command.RootOpts.Dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}
			timer.StartTimer("Bootstrap Complete")
			if err := waitForBootstrapComplete(ctx, config); err != nil {
				if err2 := command.LogClusterOperatorConditions(ctx, config); err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}

				logrus.Info("Use the following commands to gather logs from the cluster")
				logrus.Info("openshift-install gather bootstrap --help")
				logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
				logrus.Error(err.Error())
				logrus.Exit(command.ExitCodeBootstrapFailed)
			}

			logrus.Info("It is now safe to remove the bootstrap resources")
			timer.StopTimer("Bootstrap Complete")
			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster is ready",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			timer.StartTimer(timer.TotalTimeElapsed)
			ctx := context.Background()

			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(command.RootOpts.Dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}

			assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
			if err != nil {
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallFailed)
			}

			err = command.WaitForInstallComplete(ctx, config, assetStore)
			if err != nil {
				if err2 := command.LogClusterOperatorConditions(ctx, config); err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				command.LogTroubleshootingLink()
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallFailed)
			}
			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
}
