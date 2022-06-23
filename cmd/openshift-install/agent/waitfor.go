package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

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

func newWaitForBootstrapCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "bootstrap-complete",
		Short: "Wait until the cluster bootstrap is complete",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			err := runWaitForBootstrapCompleteCmd(assetDir)
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
		Run: func(cmd *cobra.Command, _ []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			err := runWaitForInstallCompleteCmd(assetDir)
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

}

func runWaitForBootstrapCompleteCmd(directory string) error {
	_, err := agentpkg.WaitForBootstrapComplete(directory)
	return err
}

func runWaitForInstallCompleteCmd(directory string) error {
	return agentpkg.WaitForInstallComplete(directory)
}
