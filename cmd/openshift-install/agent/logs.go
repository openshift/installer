package agent

import (
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	agentpkg "github.com/openshift/installer/pkg/agent"
)

// NewLogsCmd Add commands to gather logs in agent based installations
func NewLogsCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "logs",
		Short: "Retrieve installation logs",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	cmd.AddCommand(newLogsEventsCmd())
	cmd.AddCommand(newLogsClusterCmd())
	return cmd
}

func newLogsEventsCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "event",
		Short: "Print the installation events log",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			agentpkg.RetrieveEventsLog(assetDir)
		},
	}
}

func newLogsClusterCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "cluster",
		Short: "Download the cluster installation logs to a file",
		Args:  cobra.ExactArgs(0),
		Run: func(cmd *cobra.Command, args []string) {
			assetDir := cmd.Flags().Lookup("dir").Value.String()
			logrus.Debugf("asset directory: %s", assetDir)
			if len(assetDir) == 0 {
				logrus.Fatal("No cluster installation directory found")
			}
			agentpkg.DownloadClusterLogs(assetDir)
		},
	}
}
