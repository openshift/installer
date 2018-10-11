package main

import (
	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
)

var (
	rootOpts struct {
		dir      string
		logLevel string
	}
)

func main() {
	rootCmd := newRootCmd()

	var subCmds []*cobra.Command
	for _, cmd := range newTargetsCmd() {
		subCmds = append(subCmds, cmd)
	}
	subCmds = append(subCmds,
		newDestroyCmd(),
		newVersionCmd(),
		newGraphCmd(),
	)
	for _, subCmd := range subCmds {
		rootCmd.AddCommand(subCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		cause := errors.Cause(err)
		logrus.Fatalf("Error executing openshift-intall: %v", cause)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:               "openshift-install",
		Short:             "Creates OpenShift clusters",
		Long:              "",
		PersistentPreRunE: runRootCmd,
		SilenceErrors:     true,
		SilenceUsage:      true,
	}
	cmd.PersistentFlags().StringVar(&rootOpts.dir, "dir", ".", "assets directory")
	cmd.PersistentFlags().StringVar(&rootOpts.logLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")
	return cmd
}

func runRootCmd(cmd *cobra.Command, args []string) error {
	logrus.SetFormatter(&logrus.TextFormatter{
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	})
	level, err := logrus.ParseLevel(rootOpts.logLevel)
	if err != nil {
		return errors.Wrap(err, "invalid log-level")

	}
	logrus.SetLevel(level)
	return nil
}
