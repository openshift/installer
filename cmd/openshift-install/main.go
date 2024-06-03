package main

import (
	"context"
	"flag"
	"io"
	"os"
	"path/filepath"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"
	terminal "golang.org/x/term"
	"k8s.io/klog"
	klogv2 "k8s.io/klog/v2"
	"sigs.k8s.io/controller-runtime/pkg/manager/signals"

	"github.com/openshift/installer/cmd/openshift-install/command"
	"github.com/openshift/installer/pkg/clusterapi"
)

func main() {
	// This attempts to configure klog (used by vendored Kubernetes code) not
	// to log anything.
	// Handle k8s.io/klog
	var fs flag.FlagSet
	klog.InitFlags(&fs)
	fs.Set("stderrthreshold", "4")
	klog.SetOutput(io.Discard)
	// Handle k8s.io/klog/v2
	var fsv2 flag.FlagSet
	klogv2.InitFlags(&fsv2)
	fsv2.Set("stderrthreshold", "4")
	klogv2.SetOutput(io.Discard)

	installerMain()
}

func installerMain() {
	rootCmd := newRootCmd()

	// Perform a graceful shutdown upon interrupt or at exit.
	ctx := handleInterrupt(signals.SetupSignalHandler())
	logrus.RegisterExitHandler(shutdown)

	for _, subCmd := range []*cobra.Command{
		newCreateCmd(ctx),
		newDestroyCmd(),
		newWaitForCmd(),
		newGatherCmd(ctx),
		newAnalyzeCmd(),
		newVersionCmd(),
		newGraphCmd(),
		newCoreOSCmd(),
		newCompletionCmd(),
		newExplainCmd(),
		newAgentCmd(ctx),
		newListFeaturesCmd(),
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
	cmd.PersistentFlags().StringVar(&command.RootOpts.Dir, "dir", ".", "assets directory")
	cmd.PersistentFlags().StringVar(&command.RootOpts.LogLevel, "log-level", "info", "log level (e.g. \"debug | info | warn | error\")")
	return cmd
}

func runRootCmd(cmd *cobra.Command, args []string) {
	logrus.SetOutput(io.Discard)
	logrus.SetLevel(logrus.TraceLevel)

	level, err := logrus.ParseLevel(command.RootOpts.LogLevel)
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
		DisableTimestamp:       true,
		DisableLevelTruncation: true,
		DisableQuote:           true,
	}))

	if err != nil {
		logrus.Fatal(errors.Wrap(err, "invalid log-level"))
	}
}

// handleInterrupt executes a graceful shutdown then exits in
// the case of a user interrupt. It returns a new context that
// will be cancelled upon interrupt.
func handleInterrupt(signalCtx context.Context) context.Context {
	ctx, cancel := context.WithCancel(context.Background())

	// If the context from the signal handler is done,
	// an interrupt has been received, so shutdown & exit.
	go func() {
		<-signalCtx.Done()
		logrus.Warn("Received interrupt signal")
		shutdown()
		cancel()
		logrus.Exit(exitCodeInterrupt)
	}()

	return ctx
}

func shutdown() {
	clusterapi.System().Teardown()
}
