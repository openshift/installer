package main

import (
	"fmt"
	"io"
	"os"

	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	terminal "golang.org/x/term"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/nodejoiner"
)

func main() {
	if err := nodeJoiner(); err != nil {
		logrus.Fatal(err)
	}
}

func nodeJoiner() error {
	nodesAddCmd := &cobra.Command{
		Use:   "add-nodes",
		Short: "Generates an ISO that can be used to boot the configured nodes to let them join an existing cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, kubeConfig, err := getCommonFlags(cmd)
			if err != nil {
				return err
			}
			generatePXE, err := cmd.Flags().GetBool("pxe")
			if err != nil {
				return err
			}
			generateConfigISO, err := cmd.Flags().GetBool("config-iso")
			if err != nil {
				return err
			}
			return nodejoiner.NewAddNodesCommand(dir, kubeConfig, generatePXE, generateConfigISO)
		},
	}
	nodesAddCmd.Flags().BoolP("pxe", "p", false, "Instead of an ISO, generates PXE artifacts that can be used to boot the configured nodes to let them join an existing cluster")
	nodesAddCmd.Flags().BoolP("config-iso", "", false, "Generates the config ISO instead of the standard ISO")

	nodesMonitorCmd := &cobra.Command{
		Use:   "monitor-add-nodes <ip-addresses>",
		Short: "Monitors the configured nodes while they are joining an existing cluster",
		RunE: func(cmd *cobra.Command, args []string) error {
			dir, kubeConfig, err := getCommonFlags(cmd)
			if err != nil {
				return err
			}

			ips := args
			logrus.Infof("Monitoring IPs: %v", ips)
			if len(ips) == 0 {
				logrus.Fatal("At least one IP address must be specified")
			}
			return nodejoiner.NewMonitorAddNodesCommand(dir, kubeConfig, ips)
		},
	}

	rootCmd := &cobra.Command{
		Use:              "node-joiner",
		PersistentPreRun: runRootCmd,
	}
	rootCmd.PersistentFlags().String("kubeconfig", "", "Path to the kubeconfig file.")
	rootCmd.PersistentFlags().String("dir", ".", "assets directory")
	rootCmd.PersistentFlags().String("log-level", "info", "log level (e.g. \"debug | info | warn | error\")")

	rootCmd.AddCommand(nodesAddCmd)
	rootCmd.AddCommand(nodesMonitorCmd)

	return rootCmd.Execute()
}

func getCommonFlags(cmd *cobra.Command) (string, string, error) {
	kubeConfig, err := cmd.Flags().GetString("kubeconfig")
	if err != nil {
		return "", "", err
	}
	dir, err := cmd.Flags().GetString("dir")
	if err != nil {
		return "", "", err
	}
	return dir, kubeConfig, nil
}

func runRootCmd(cmd *cobra.Command, args []string) {
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	logLevel, err := cmd.Flags().GetString("log-level")
	if err != nil {
		logrus.Fatal(err)
	}

	level, err := logrus.ParseLevel(logLevel)
	if err != nil {
		level = logrus.InfoLevel
	}

	logrus.AddHook(command.NewFileHookWithNewlineTruncate(os.Stderr, level, &logrus.TextFormatter{
		// Setting ForceColors is necessary because logrus.TextFormatter determines
		// whether or not to enable colors by looking at the output of the logger.
		// In this case, the output is io.Discard, which is not a terminal.
		// Overriding it here allows the same check to be done, but against the
		// hook's output instead of the logger's output.
		ForceColors:            terminal.IsTerminal(int(os.Stderr.Fd())),
		DisableLevelTruncation: true,
		DisableTimestamp:       false,
		FullTimestamp:          true,
		DisableQuote:           true,
	}))

	if err != nil {
		logrus.Fatal(fmt.Errorf("invalid log-level: %w", err))
	}
}
