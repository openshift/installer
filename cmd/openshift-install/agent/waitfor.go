package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	agentcmd "github.com/openshift/installer/pkg/agent"
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

	cmd.AddCommand(newWaitForClusterValidationSuccessCmd())
	cmd.AddCommand(newWaitForBootstrapCompleteCmd())
	cmd.AddCommand(newWaitForInstallCompleteCmd())
	return cmd
}

func newWaitForClusterValidationSuccessCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster-validated",
		Short: "Wait until the cluster manifests are validated for install",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runWaitForClusterValidationSuccessCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newWaitForBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until the cluster bootstrap is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runWaitForInstallCompleteCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster installation is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runWaitForInstallCompleteCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

}

func runWaitForClusterValidationSuccessCmd() error {
	return agentcmd.WaitForClusterValidationSuccess()
}

func runWaitForBootstrapCompleteCmd() error {
	return agentcmd.WaitForBootstrapComplete()
}

func runWaitForInstallCompleteCmd() error {
	return agentcmd.WaitForInstallComplete()
}
