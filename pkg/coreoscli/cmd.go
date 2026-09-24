package coreoscli

import (
	"context"
	"fmt"
	"os"
	"slices"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/rhcos"
	"github.com/openshift/installer/pkg/types"
	"github.com/openshift/installer/pkg/version"
)

func selectOSImageStream(streamFlag string, releaseVersionInjected bool) (types.OSImageStream, error) {
	validStreams := types.OSImageStreamValues()
	if streamFlag != "" {
		s := types.OSImageStream(streamFlag)
		if !slices.Contains(validStreams, s) {
			return "", fmt.Errorf("invalid value %q for --stream; must be one of %v", streamFlag, validStreams)
		}
		return s, nil
	}

	if !releaseVersionInjected {
		return "", fmt.Errorf("release version metadata was not injected into the installer; specify --stream with one of %v", validStreams)
	}

	return rhcos.BuildDefaultOSImageStream(), nil
}

// printStreamJSON is the implementation of print-stream-json
func printStreamJSON(cmd *cobra.Command, _ []string) error {
	streamFlag, err := cmd.Flags().GetString("stream")
	if err != nil {
		return err
	}

	osImageStream, err := selectOSImageStream(streamFlag, version.IsReleaseVersionInjected())
	if err != nil {
		return err
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
	printStreamCmd.Flags().StringVar(&stream, "stream", "", fmt.Sprintf("OS image stream to use (one of %v)", types.OSImageStreamValues()))
	cmd.AddCommand(printStreamCmd)

	return cmd
}
