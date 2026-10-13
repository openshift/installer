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

package interactive

import (
	"fmt"

	"github.com/spf13/cobra"

	"github.com/openshift/rosa/pkg/arguments"
)

var mode string

const (
	Mode       = "mode"
	ModeAuto   = "auto"
	ModeManual = "manual"
)

var Modes = []string{ModeAuto, ModeManual}

func AddModeFlag(cmd *cobra.Command) {
	cmd.Flags().StringVarP(
		&mode,
		"mode",
		"m",
		"",
		"How to perform the operation. Valid options are:\n"+
			"auto: Resource changes will be automatic applied using the current AWS account\n"+
			"manual: Commands necessary to modify AWS resources will be output to be run manually",
	)
	cmd.RegisterFlagCompletionFunc("mode", modeCompletion)
}

func SetModeKey(key string) {
	mode = key
}

func GetMode() (string, error) {
	if mode == "" {
		return "", nil
	}
	if !arguments.IsValidMode(Modes, mode) {
		return "", fmt.Errorf("Invalid mode. Allowed values are %s", Modes)
	}
	return mode, nil
}

func modeCompletion(cmd *cobra.Command, args []string, toComplete string) ([]string, cobra.ShellCompDirective) {
	return Modes, cobra.ShellCompDirectiveDefault
}

func GetOptionMode(cmd *cobra.Command, mode string, question string) (string, error) {
	mode, err := GetOption(Input{
		Question: question,
		Help:     cmd.Flags().Lookup(Mode).Usage,
		Default:  mode,
		Options:  Modes,
		Required: true,
	})
	if err != nil {
		return mode, fmt.Errorf("invalid mode: %v", err)
	}
	SetModeKey(mode)
	return mode, nil
}
