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

// This file contains functions used to implement the '--debug' command line option.

package debug

import (
	"github.com/spf13/pflag"
)

// AddFlag adds the debug flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVar(
		&enabled,
		"debug",
		false,
		"Enable debug mode.",
	)
}

// Enabled returns a boolean flag that indicates if the debug mode is enabled.
func Enabled() bool {
	return enabled
}

func SetEnabled(debugEnabled bool) {
	enabled = debugEnabled
}

// enabled is a boolean flag that indicates that the debug mode is enabled.
var enabled bool
