package main

import (
	"io"
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/data"
)

func newRhcosMetadataCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "rhcos-metadata",
		Short: "Output Red Hat Enterprise Linux CoreOS metaadata",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE:  runRhcosMetadata,
	}
}

// runRhcosMetadata simply dumps the embedded RHCOS metadata to
// stdout. See also https://github.com/openshift/installer/issues/1399
func runRhcosMetadata(cmd *cobra.Command, args []string) error {
	file, err := data.Assets.Open("rhcos.json")
	if err != nil {
		return err
	}
	defer file.Close()
	_, err = io.Copy(os.Stdout, file)
	if err != nil {
		return err
	}
	return nil
}
