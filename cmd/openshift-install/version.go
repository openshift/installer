package main

import (
	"fmt"
	"os"

	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/version"
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "",
		Args:  cobra.ExactArgs(0),
		RunE:  runVersionCmd,
	}
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s %s\n", os.Args[0], version.Raw)
	return nil
}
