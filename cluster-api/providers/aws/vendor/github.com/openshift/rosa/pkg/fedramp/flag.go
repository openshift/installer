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

// This file contains functions used to implement the '--govcloud' command line option.

package fedramp

import (
	"github.com/spf13/cobra"
	"github.com/spf13/pflag"

	"github.com/openshift/rosa/pkg/config"
)

// AddFlag adds the govcloud flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVar(
		&enabled,
		"govcloud",
		false,
		"Uses the FedRAMP High OpenShift Cluster Manager API for creating clusters in AWS GovCloud regions",
	)

	flags.BoolVar(
		&enabled,
		"admin",
		false,
		"Uses the FedRAMP High OpenShift Cluster Manager API Endpoint for Administrator Access",
	)
	flags.MarkHidden("admin")
}

func HasFlag(cmd *cobra.Command) bool {
	flag := cmd.Flags().Lookup("govcloud")
	if flag == nil {
		return false
	}
	return flag.Changed
}

func HasAdminFlag(cmd *cobra.Command) bool {
	flag := cmd.Flags().Lookup("admin")
	if flag == nil {
		return false
	}
	return flag.Changed
}

// Enabled returns a boolean flag that indicates if the fedramp mode is enabled.
func Enabled() bool {
	if enabled {
		return true
	}
	cfg, err := config.Load()
	if err != nil {
		return false
	}
	if cfg != nil && cfg.FedRAMP {
		Enable()
	}
	return enabled
}

// Enable sets the flag for the rest of the command
func Enable() {
	enabled = true
}

func Disable() {
	enabled = false
	cfg, err := config.Load()
	if err != nil || cfg == nil {
		return
	}
	cfg.FedRAMP = false
	config.Save(cfg)
}

// enabled is a boolean flag that indicates that the govcloud mode is enabled.
var enabled bool
