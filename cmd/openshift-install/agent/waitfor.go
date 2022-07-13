package agent

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	agentpkg "github.com/openshift/installer/pkg/agent"
)

const (
	exitCodeInstallConfigError = iota + 3
	exitCodeInfrastructureFailed
	exitCodeBootstrapFailed
	exitCodeInstallFailed
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

func newWaitForBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until the cluster bootstrap is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			cluster, err := runWaitForBootstrapCompleteCmd(assetDir)
			if err != nil {
				logrus.Debug("Printing the event list gathered from the Agent Rest API")
				cluster.PrintInfraEnvRestAPIEventList()
				err2 := cluster.API.OpenShift.LogClusterOperatorConditions()
				if err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				logrus.Info("Use the following commands to gather logs from the cluster")
				logrus.Info("openshift-install gather bootstrap --help")
				logrus.Error(errors.Wrap(err, "Bootstrap failed to complete: "))
				logrus.Exit(exitCodeBootstrapFailed)
			}
		},
	}
}

func runWaitForBootstrapCompleteCmd(directory string) (*agentpkg.Cluster, error) {
	cluster, err := agentpkg.WaitForBootstrapComplete(directory)
	return cluster, err
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster installation is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			cluster, err := runWaitForInstallCompleteCmd(assetDir)
			if err != nil {
				logrus.Debug("Printing the event list gathered from the Agent Rest API")
				cluster.PrintInfraEnvRestAPIEventList()
				err2 := cluster.API.OpenShift.LogClusterOperatorConditions()
				if err2 != nil {
					logrus.Error("Attempted to gather ClusterOperator status after wait failure: ", err2)
				}
				logrus.Error(`Cluster initialization failed because one or more operators are not functioning properly.
				The cluster should be accessible for troubleshooting as detailed in the documentation linked below,
				https://docs.openshift.com/container-platform/latest/support/troubleshooting/troubleshooting-installations.html`)
				logrus.Exit(exitCodeInstallFailed)
			}
			cluster.PrintInstallationComplete()
		},
	}
}

func runWaitForInstallCompleteCmd(directory string) (*agentpkg.Cluster, error) {
	cluster, err := agentpkg.WaitForInstallComplete(directory)
	return cluster, err
}
