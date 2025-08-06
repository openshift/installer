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

package confirm

import (
	"fmt"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/spf13/pflag"
)

var yes bool

// AddFlag adds the --yes flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.BoolVarP(
		&yes,
		"yes",
		"y",
		false,
		"Automatically answer yes to confirm operation.",
	)
}

func Yes() bool {
	return yes
}

// Asks the user to confirm the operation, using the specified string as the message
func ConfirmRaw(q string) bool {
	return Prompt(false, q)
}

func Confirm(q string, v ...interface{}) bool {
	msg := fmt.Sprintf("Are you sure you want to %s?", fmt.Sprintf(q, v...))
	return ConfirmRaw(msg)
}

func Prompt(dflt bool, q string, v ...interface{}) bool {
	if yes {
		return yes
	}
	prompt := &survey.Confirm{
		Message: fmt.Sprintf(q, v...),
		Default: dflt,
	}
	response := false
	survey.AskOne(prompt, &response, survey.WithValidator(survey.Required))
	return response
}
