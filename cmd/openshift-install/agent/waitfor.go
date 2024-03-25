package agent

import (
	"context"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/cmd/openshift-install/command"
	agentpkg "github.com/openshift/installer/pkg/agent"
)

// NewWaitForCmd create the commands for waiting the completion of the agent based cluster installation.
func NewWaitForCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "wait-for",
		Short: "Wait for install-time events",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newWaitForBootstrapCompleteCmd())
	cmd.AddCommand(newWaitForInstallCompleteCmd())
	return cmd
}

func handleBootstrapError(cluster *agentpkg.Cluster, err error) {
	logrus.Debug("Printing the event list gathered from the Agent Rest API")
	cluster.PrintInfraEnvRestAPIEventList()
	err2 := cluster.API.OpenShift.LogClusterOperatorConditions()
	if err2 != nil {
		logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
	}
	logrus.Info("Use the following commands to gather logs from the cluster")
	logrus.Info("openshift-install gather bootstrap --help")
	logrus.Error(errors.Wrap(err, "Bootstrap failed to complete: "))
	logrus.Exit(command.ExitCodeBootstrapFailed)
}

func newWaitForBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until the cluster bootstrap is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}

			ctx := context.Background()
			cluster, err := agentpkg.NewCluster(ctx, assetDir)
			if err != nil {
				logrus.Exit(command.ExitCodeBootstrapFailed)
			}

			if err := agentpkg.WaitForBootstrapComplete(cluster); err != nil {
				handleBootstrapError(cluster, err)
			}
		},
	}
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster installation is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			cleanup := command.SetupFileHook(command.RootOpts.Dir)
			defer cleanup()

			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}

			ctx := context.Background()
			cluster, err := agentpkg.NewCluster(ctx, assetDir)
			if err != nil {
				logrus.Exit(command.ExitCodeBootstrapFailed)
			}

			if err := agentpkg.WaitForBootstrapComplete(cluster); err != nil {
				handleBootstrapError(cluster, err)
			}

			config := cluster.API.Kube.Config

			err = command.WaitForInstallComplete(ctx, config, command.RootOpts.Dir)
			if err != nil {
				if err2 := command.LogClusterOperatorConditions(ctx, config); err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				command.LogTroubleshootingLink()
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallFailed)
			}
		},
	}
}
