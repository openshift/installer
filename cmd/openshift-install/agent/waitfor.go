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

	cmd.AddCommand(newWaitForInstallCompleteCmd())
	return cmd
}

func newWaitForInstallCompleteCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "install-complete",
		Short: "Wait until the cluster is ready",
		Args:  cobra.ExactArgs(0),
		Run: func(_ *cobra.Command, _ []string) {
			err := runWaitForInstallCompleteCmd()
			if err != nil {
				logrus.Fatal(err)
			}
		},
	}

}

func runWaitForInstallCompleteCmd() error {
	return agentcmd.WaitFor()
}
