/*
Copyright (c) 2022 Red Hat, Inc.

Licensed under the Apache License, Version 2.0 (the "License");
you may not use this file except in compliance with the License.
You may obtain a copy of the License at

  http://www.apache.org/licenses/LICENSE-2.0

Unless required by applicable law or agreed to in writing, software
distributed under the License is distributed on an "AS IS" BASIS,
WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
See the License for the specific language governing permissions and
limitations under the License.
*/

// This file contains functions used to implement the '--color' command line option.

package color

import (
	"fmt"
	"os"
	"runtime"

	"github.com/spf13/cobra"
)

var color string

var options = []string{"auto", "never", "always"}

// AddFlag adds the interactive flag to the given set of command line flags.
func AddFlag(cmd *cobra.Command) {
	cmd.PersistentFlags().StringVar(
		&color,
		"color",
		"auto",
		fmt.Sprintf("Surround certain characters with escape sequences to display them in color "+
			"on the terminal. Allowed options are %s", options),
	)

	cmd.RegisterFlagCompletionFunc("color", completion)
}

func completion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return options, cobra.ShellCompDirectiveDefault
}

// UseColor returns a bool that indicates whether the color is enabled
func UseColor() bool {
	switch color {
	case "never":
		return false
	case "always":
		return true
	case "auto":
		fallthrough
	default:
		if runtime.GOOS == "windows" {
			return false
		}
		stdout, err := os.Stdout.Stat()
		if err != nil {
			return true
		}
		return (stdout.Mode()&os.ModeDevice != 0) && (stdout.Mode()&os.ModeNamedPipe == 0)
	}
}

func SetColor(colorOption string) {
	color = colorOption
}
