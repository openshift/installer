package rosa

import (
	"context"
	"os"

	"github.com/spf13/cobra"
)

// RuntimeVisitor are functions that configure the Runtime for a command.
type RuntimeVisitor func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string)

// CommandRunner is a function supplied by Commands of the ROSA CLI that perform the actual logic of running
// the command
type CommandRunner func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string) error

// DefaultRunner is a centralised implementation of the default Cobra Command.run function that takes care
// of instantiating several key resources on behalf of a command
func DefaultRunner(visitor RuntimeVisitor, runner CommandRunner) func(command *cobra.Command, args []string) {
	return func(command *cobra.Command, args []string) {
		ctx := context.Background()
		r := NewRuntime()
		defer r.Cleanup()

		if visitor != nil {
			visitor(ctx, r, command, args)
		}

		err := runner(ctx, r, command, args)
		if err != nil {
			r.Reporter.Errorf(err.Error())
			os.Exit(1)
		}
	}
}

// DefaultRuntime returns a Runtime with the most basic of setups. None of the clients are initialised.
func DefaultRuntime() RuntimeVisitor {
	return func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string) {}
}

// RuntimeWithOCM configures the Runtime with an OCM Client
func RuntimeWithOCM() RuntimeVisitor {
	return func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string) {
		runtime.WithOCM()
	}
}

// RuntimeWithOCMAndAWS configures the Runtime with an OCM Client and AWS client
func RuntimeWithOCMAndAWS() RuntimeVisitor {
	return func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string) {
		runtime.WithOCM().WithAWS()
	}
}

// RuntimeWithAWS configures the Runtime with an AWS client
func RuntimeWithAWS() RuntimeVisitor {
	return func(ctx context.Context, runtime *Runtime, command *cobra.Command, args []string) {
		runtime.WithAWS()
	}
}
