package coreoscli

import (
	"context"
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
)

// printStreamJSON is the implementation of print-stream-json
func printStreamJSON(cmd *cobra.Command, _ []string) error {
	osImageStream := rhcos.DefaultOSImageStream
	streamFlag, err := cmd.Flags().GetString("stream")
	if err != nil {
		return err
	}
	if streamFlag != "" {
		s := types.OSImageStream(streamFlag)
		valid := false
		for _, v := range types.OSImageStreamValues {
			if s == v {
				valid = true
				break
			}
		}
		if !valid {
			return fmt.Errorf("invalid value %q for --stream; must be one of %v", streamFlag, types.OSImageStreamValues)
		}
		osImageStream = s
	}
	streamData, err := rhcos.FetchRawCoreOSStream(context.Background(), osImageStream)
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

	var stream string
	printStreamCmd := &cobra.Command{
		Use:   "print-stream-json",
		Short: "Outputs the CoreOS stream metadata for the bootimages",
		Args:  cobra.ExactArgs(0),
		RunE:  printStreamJSON,
	}
	printStreamCmd.Flags().StringVar(&stream, "stream", "", fmt.Sprintf("OS image stream to use (one of %v)", types.OSImageStreamValues))
	cmd.AddCommand(printStreamCmd)

	return cmd
}
