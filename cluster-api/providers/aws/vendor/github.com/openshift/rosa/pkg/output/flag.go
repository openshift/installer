/*
Copyright (c) 2021 Red Hat, Inc.

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

// This file contains functions used to implement the '--output' command line option.

package output

import (
	"fmt"

	"github.com/spf13/cobra"
)

var o string

var formats = []string{"json", "yaml"}

// AddFlag adds the interactive flag to the given set of command line flags.
func AddFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&o,
		"output",
		"o",
		"",
		fmt.Sprintf("Output format. Allowed formats are %s", formats),
	)

	cmd.RegisterFlagCompletionFunc("output", completion)
}

func completion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return formats, cobra.ShellCompDirectiveDefault
}

func HasFlag() bool {
	return o != ""
}

// Enabled retursn a boolean flag that indicates if the interactive mode is enabled.
func Output() string {
	return o
}
