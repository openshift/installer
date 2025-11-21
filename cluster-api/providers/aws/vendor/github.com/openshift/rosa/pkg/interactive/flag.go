/*
Copyright (c) 2020 Red Hat, Inc.

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

// This file contains functions used to implement the '--interactive' command line option.

package interactive

import (
	"github.com/spf13/pflag"
)

// AddFlag adds the interactive flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVarP(
		&enabled,
		"interactive",
		"i",
		false,
		"Enable interactive mode.",
	)
}

// Enabled returns a boolean flag that indicates if the interactive mode is enabled.
func Enabled() bool {
	return enabled
}

// Enable enables the interactive mode
func Enable() {
	enabled = true
}

func SetEnabled(setEnabled bool) {
	enabled = setEnabled
}

// enabled is a boolean flag that indicates that the interactive mode is enabled.
var enabled bool
