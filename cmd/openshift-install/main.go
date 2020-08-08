package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"math"
	"os"
	"path/filepath"
	"runtime"
	"strings"
	"time"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	"golang.org/x/crypto/ssh/terminal"
	"k8s.io/klog"

	"github.com/openshift/installer/pkg/terraform/exec/plugins"
	"github.com/openshift/installer/pkg/version"
)

var (
	rootOpts struct {
		dir      string
		logLevel string
	}
)

func main() {
	// This attempts to configure klog (used by vendored Kubernetes code) not
	// to log anything.
	var fs flag.FlagSet
	klog.InitFlags(&fs)
	fs.Set("stderrthreshold", "4")
	klog.SetOutput(ioutil.Discard)

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
		newMigrateCmd(),
		newExplainCmd(),
	} {
		rootCmd.AddCommand(subCmd)
	}

	if err := rootCmd.Execute(); err != nil {
		logrus.Fatalf("Error executing openshift-install: %v", err)
	}
}

func newRootCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:              filepath.Base(os.Args[0]),
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

	if version.CommitDate != "" {
		logOldInstallerNotification(version.CommitDate)
	}
}

func logOldInstallerNotification(commitDate string) {
	commit, err := time.Parse("2006-01-02", commitDate)
	if err != nil {
		logrus.Warnf("Failed to parse commit date %q: %v", version.CommitDate, err)
		return
	}

	releaseName, err := version.Version()
	if err != nil {
		logrus.Warnf("Failed to parse version: %v", err)
		return
	}

	now := time.Now()
	age := now.Sub(commit)
	days := int(math.Floor(age.Hours() / 24))
	arch := runtime.GOARCH
	switch arch {
	case "amd64":
		arch = "x86_64"
	}
	uri := fmt.Sprintf("https://mirror.openshift.com/pub/openshift-v4/%s/clients/ocp/latest/", arch)
	if days > 30 {
		logrus.Infof("The OpenShift Installer installs a version of OpenShift which is matched to the Installer's version, %s. This Installer is more than %d days old. Red Hat frequently updates OpenShift with bug fixes, please consider using a newer installer. You may download the latest from %s", releaseName, days, uri)
	}
}
