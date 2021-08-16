package main

import (
	"context"
	"path/filepath"

	"github.com/openshift/installer/pkg/asset/installconfig"
	assetstore "github.com/openshift/installer/pkg/asset/store"
	"github.com/openshift/installer/pkg/metrics/gatherer"
	timer "github.com/openshift/installer/pkg/metrics/timer"
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
			timer.StartTimer(timer.TotalTimeElapsed)
			ctx := context.Background()
			gatherer.InitializeInvocationMetrics(gatherer.WaitforMetricName)

			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}
			timer.StartTimer("Bootstrap Complete")
			if err := waitForBootstrapComplete(ctx, config); err != nil {
				if err2 := logClusterOperatorConditions(ctx, config); err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}

				logrus.Info("Use the following commands to gather logs from the cluster")
				logrus.Info("openshift-install gather bootstrap --help")
				logrus.Error("Bootstrap failed to complete: ", err.Unwrap())
				logrus.Error(err.Error())
				logrus.Fatal("Bootstrap failed to complete")
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
			gatherer.InitializeInvocationMetrics(gatherer.WaitforMetricName)

			cleanup := setupFileHook(rootOpts.dir)
			defer cleanup()

			config, err := clientcmd.BuildConfigFromFlags("", filepath.Join(rootOpts.dir, "auth", "kubeconfig"))
			if err != nil {
				logrus.Fatal(errors.Wrap(err, "loading kubeconfig"))
			}

			err = waitForInstallComplete(ctx, config, rootOpts.dir)
			if err != nil {
				if err2 := logClusterOperatorConditions(ctx, config); err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				logTroubleshootingLink()
				logrus.Fatal(err)
			}
			if assetStore, err := assetstore.NewStore(rootOpts.dir); err == nil {
				if asset, err := assetStore.Load(&installconfig.InstallConfig{}); err == nil && asset != nil {
					gatherer.AddLabelValue(gatherer.CurrentInvocationContext, "platform", asset.(*installconfig.InstallConfig).Config.Platform.Name())
				}
			}
			timer.StopTimer(timer.TotalTimeElapsed)
			timer.LogSummary()
		},
	}
}
