package main

import (
	"flag"
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
	// This attempts to configure glog (used by vendored Kubernetes code) not
	// to log anything. Nobody likes you, glog. Go away.
	flag.CommandLine.Parse([]string{})
	flag.CommandLine.Set("stderrthreshold", "4")

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
		newWaitForCmd(),
		newGatherCmd(),
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
		Use:              "openshift-install",
		Short:            "Creates OpenShift clusters",
		Long:             "",
		PersistentPreRun: runRootCmd,
		SilenceErrors:    true,
		SilenceUsage:     true,
	}
	cmd.PersistentFlags().StringVar(&rootOpts.dir, "dir", ".", "assets directory")
	cmd.PersistentFlags().StringVar(&rootOpts.logLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")
	return cmd
}

func runRootCmd(cmd *cobra.Command, args []string) {
	logrus.SetOutput(ioutil.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	level, err := logrus.ParseLevel(rootOpts.logLevel)
	if err != nil {
		level = logrus.InfoLevel
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

	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log-level"))
	}
}
