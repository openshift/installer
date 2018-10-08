package main

import (
	"fmt"
	"os"
	"os/exec"
	"strings"

	"github.com/pkg/errors"
	"github.com/sirupsen/logrus"
	"github.com/spf13/cobra"

	"github.com/openshift/installer/pkg/terraform"
)

var (
	version = "was not built correctly" // set in hack/build.sh
)

func newVersionCmd() *cobra.Command {
	return &cobra.Command{
		Use:   "version",
		Short: "Print version information",
		Long:  "",
		RunE:  runVersionCmd,
	}
}

func runVersionCmd(cmd *cobra.Command, args []string) error {
	fmt.Printf("%s %s\n", os.Args[0], version)
	terraformVersion, err := terraform.Version()
	if err != nil {
		exitError, ok := err.(*exec.ExitError)
		if ok && len(exitError.Stderr) > 0 {
			logrus.Error(strings.Trim(string(exitError.Stderr), "\n"))
		}
		return errors.Wrap(err, "Failed to calculate Terraform version")
	}
	fmt.Println(terraformVersion)
	return nil
}
