package main

import (
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/nodejoiner"
)

func main() {
	wd, err := os.Getwd()
	if err != nil {
		logrus.Fatal(err)
	}

	nodesAddCmd := &cobra.Command{
		Use:   "add-nodes",
		Short: "Generates an ISO that could be used to boot the configured nodes to let them join an existing cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			kubeConfig, err := cmd.Flags().GetString("kubeconfig")
			if err != nil {
				return err
			}
			return nodejoiner.NewAddNodesCommand(wd, kubeConfig)
		},
	}

	nodesMonitorCmd := &cobra.Command{
		Use:   "monitor-add-nodes",
		Short: "Monitors the configured nodes while they are joining an existing cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			return nodejoiner.NewMonitorAddNodesCommand(wd)
		},
	}

	rootCmd := &cobra.Command{
		Use: "node-joiner",
	}
	rootCmd.PersistentFlags().String("kubeconfig", "", "Path to the kubeconfig file.")

	rootCmd.AddCommand(nodesAddCmd)
	rootCmd.AddCommand(nodesMonitorCmd)
	if err := rootCmd.Execute(); err != nil {
		logrus.Fatal(err)
	}
}
