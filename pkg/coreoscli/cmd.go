package coreoscli

import (
	"context"
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/rhcos"
)

// printStreamJSON is the implementation of print-stream-json
func printStreamJSON(cmd *cobra.Command, _ []string) error {
	streamData, err := rhcos.FetchRawCoreOSStream(context.Background())
	if err != nil {
		return err
	}
	os.Stdout.Write(streamData)
	return nil
}

// NewCmd returns a subcommand for explain
func NewCmd() *cobra.Command {
	cmd := &cobra.Command{
		Use:   "coreos",
		Short: "Commands for operating on CoreOS boot images",
		Args:  cobra.ExactArgs(0),
		RunE: func(cmd *cobra.Command, args []string) error {
			return cmd.Help()
		},
	}

	printStreamCmd := &cobra.Command{
		Use:   "print-stream-json",
		Short: "Outputs the CoreOS stream metadata for the bootimages",
		Args:  cobra.ExactArgs(0),
		RunE:  printStreamJSON,
	}
	cmd.AddCommand(printStreamCmd)

	return cmd
}
