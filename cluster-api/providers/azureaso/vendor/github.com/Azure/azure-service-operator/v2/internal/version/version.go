/*
Copyright (c) Microsoft Corporation.
Licensed under the MIT license.
*/

package version

import (
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/spf13/cobra"
)

var (
	// This is populated from the build (see Taskfile.yml)
	BuildVersion string = ""
)

// NewCommand creates a new reusable cobra command to display the current version of the tool
func NewCommand() (*cobra.Command, error) {
	cmd := &cobra.Command{
		Use:   "version",
		Short: "Display version information",
		RunE:  versionCommand,
	}

	return cmd, nil
}

func versionCommand(cmd *cobra.Command, args []string) error {
	ver := BuildVersion
	if ver == "" {
		ver = "dev"
	}

	path, err := os.Executable()
	if err != nil {
		return err
	}

	fmt.Printf(
		"%s %s %s\n",
		filepath.Base(path),
		ver,
		runtime.GOOS)

	return nil
}
