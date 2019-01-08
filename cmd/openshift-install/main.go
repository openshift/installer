package main

import (
	"io/ioutil"
	"os"
	"path/filepath"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"

	"github.com/openshift/installer/pkg/terraform/exec/plugins"
)

var (
	rootOpts struct {
		dir      string
		logLevel string
	}
)

func main() {
	if len(os.Args) > 0 {
		base := filepath.Base(os.Args[0])
		cname := strings.TrimSuffix(base, filepath.Ext(base))
		if pluginRunner, ok := plugins.KnownPlugins[cname]; ok {
			pluginRunner()
			return
		}
	}

	installerMain()
}

func installerMain() {
	rootCmd := newRootCmd()

	for _, subCmd := range []*cobra.Command{
		newCreateCmd(),
		newDestroyCmd(),
		newVersionCmd(),
		newGraphCmd(),
		newCompletionCmd(),
	} {
		rootCmd.AddCommand(subCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error executing openshift-install: %v", err)
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
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	level, err := logrus.ParseLevel(rootOpts.logLevel)
	if err != nil {
		return errors.Wrap(err, "invalid log-level")
	}

	logrus.AddHook(newFileHook(os.Stderr, level, &logrus.TextFormatter{
		// Setting ForceColors is necessary because logrus.TextFormatter determines
		// whether or not to enable colors by looking at the output of the logger.
		// In this case, the output is ioutil.Discard, which is not a terminal.
		// Overriding it here allows the same check to be done, but against the
		// hook's output instead of the logger's output.
		ForceColors:            terminal.IsTerminal(int(os.Stderr.Fd())),
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
	}))

	return nil
}
