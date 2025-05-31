package agent

import (
	"context"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"k8s.io/client-go/rest"

	"github.com/openshift/installer/cmd/openshift-install/command"
	agentpkg "github.com/openshift/installer/pkg/agent"
	"github.com/openshift/installer/pkg/asset/agent/workflow"
	assetstore "github.com/openshift/installer/pkg/asset/store"
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

func handleBootstrapError(ctx context.Context, config *rest.Config, cluster *agentpkg.Cluster, err error) {
	logrus.Debug("Printing the event list gathered from the Agent Rest API")
	cluster.PrintInfraEnvRestAPIEventList()
	err2 := command.LogClusterOperatorConditions(ctx, config)
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

			kubeconfigPath := filepath.Join(assetDir, "auth", "kubeconfig")

			rendezvousIP, sshKey, err := agentpkg.FindRendezvouIPAndSSHKeyFromAssetStore(assetDir)
			if err != nil {
				logrus.Fatal(err)
			}

			ctx := context.Background()
			cluster, err := agentpkg.NewCluster(ctx, assetDir, rendezvousIP, kubeconfigPath, sshKey, workflow.AgentWorkflowTypeInstall)
			if err != nil {
				logrus.Exit(command.ExitCodeBootstrapFailed)
			}

			if err := agentpkg.WaitForBootstrapComplete(cluster); err != nil {
				handleBootstrapError(ctx, cluster.API.Kube.Config, cluster, err)
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

			kubeconfigPath := filepath.Join(assetDir, "auth", "kubeconfig")

			rendezvousIP, sshKey, err := agentpkg.FindRendezvouIPAndSSHKeyFromAssetStore(assetDir)
			if err != nil {
				logrus.Fatal(err)
			}

			ctx := context.Background()
			cluster, err := agentpkg.NewCluster(ctx, assetDir, rendezvousIP, kubeconfigPath, sshKey, workflow.AgentWorkflowTypeInstall)
			if err != nil {
				logrus.Exit(command.ExitCodeBootstrapFailed)
			}

			if err := agentpkg.WaitForBootstrapComplete(cluster); err != nil {
				handleBootstrapError(ctx, cluster.API.Kube.Config, cluster, err)
			}

			assetStore, err := assetstore.NewStore(command.RootOpts.Dir)
			if err != nil {
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallFailed)
			}

			if err = command.WaitForInstallComplete(ctx, cluster.API.Kube.Config, assetStore); err != nil {
				logrus.Error(err)
				err2 := command.LogClusterOperatorConditions(ctx, cluster.API.Kube.Config)
				if err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				command.LogTroubleshootingLink()
				logrus.Error(err)
				logrus.Exit(command.ExitCodeInstallFailed)
			}
		},
	}
}
