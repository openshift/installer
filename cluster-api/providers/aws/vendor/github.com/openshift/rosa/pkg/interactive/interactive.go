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

package interactive

import (
	"fmt"
	"net"
	"os"
	"strconv"

	survey "github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	"github.com/AlecAivazis/survey/v2/terminal"
	commonUtils "github.com/openshift-online/ocm-common/pkg/utils"

	"github.com/openshift/rosa/pkg/color"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/interactive/consts"
)

type Input struct {
	Question       string
	Help           string
	Options        []string
	Default        interface{}
	DefaultMessage string
	Required       bool
	Validators     []Validator
}

// Gets string input from the command line
func GetString(input Input) (a string, err error) {
	transformer := survey.TransformString(helper.HandleEscapedEmptyString)
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(string)
	if !ok {
		dflt = ""
	}
	question := input.Question
	if !input.Required && dflt == "" {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Default: dflt,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &a, survey.WithValidator(compose(input.Validators)))
	a = transformer(a).(string)
	return
}

// Gets int number input from the command line
func GetInt(input Input) (a int, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(int)
	if !ok {
		dflt = 0
	}
	dfltStr := fmt.Sprintf("%d", dflt)
	if dfltStr == "0" && input.Required {
		dfltStr = ""
	}
	question := input.Question
	if !input.Required && dfltStr == "" {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Default: dfltStr,
	}
	var str string
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &str, survey.WithValidator(compose(input.Validators)))
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	return parseInt(str)
}

func parseInt(str string) (num int, err error) {
	return strconv.Atoi(str)
}

// Gets float number input from the command line
func GetFloat(input Input) (a float64, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(float64)
	if !ok {
		dflt = 0
	}
	dfltStr := fmt.Sprintf("%f", dflt)
	if dfltStr == "0" {
		dfltStr = ""
	}
	question := input.Question
	if !input.Required && dfltStr == "" {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
	}
	if input.Default != nil {
		prompt.Default = dfltStr
	}
	var str string
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &str, survey.WithValidator(compose(input.Validators)))
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	return parseFloat(str)
}

func parseFloat(str string) (num float64, err error) {
	return strconv.ParseFloat(str, commonUtils.MaxByteSize)
}

// Asks for multiple options selection
func GetMultipleOptions(input Input) ([]string, error) {
	core.DisableColor = !color.UseColor()
	var err error
	res := make([]string, 0)
	dflt, ok := input.Default.([]string)
	if !ok {
		dflt = []string{}
	}
	question := input.Question
	if !input.Required && len(dflt) == 0 {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.MultiSelect{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Options: input.Options,
		Default: dflt,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &res, survey.WithValidator(compose(input.Validators)))
	return res, err
}

// Asks for option selection in the command line
func GetOption(input Input) (a string, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(string)
	if !ok {
		dflt = ""
	}
	defaultMessage := ""
	if dflt != "" {
		if input.DefaultMessage != "" {
			defaultMessage = input.DefaultMessage
		} else {
			defaultMessage = fmt.Sprintf("default = '%s'", dflt)
		}
	}
	question := input.Question
	optionalMessage := ""
	if !input.Required {
		optionalMessage = fmt.Sprintf("optional, choose '%s' to skip selection", consts.SkipSelectionOption)
		input.Options = append([]string{consts.SkipSelectionOption}, input.Options...)
		if dflt == "" {
			dflt = consts.SkipSelectionOption
		} else {
			optionalMessage += ". The default value will be provided"
		}
	}
	if optionalMessage != "" || defaultMessage != "" {
		question = fmt.Sprintf("%s (%s", question, optionalMessage)
		separator := ""
		if optionalMessage != "" {
			separator = "; "
		}
		question = fmt.Sprintf("%s%s%s)", question, separator, defaultMessage)
	}
	// if default is empty or not in the options, default to the first available option
	if (dflt == "" || !containsString(input.Options, dflt)) && len(input.Options) > 0 {
		dflt = input.Options[0]
	}
	prompt := &survey.Select{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Options: input.Options,
		Default: dflt,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &a, survey.WithValidator(compose(input.Validators)))
	if a == consts.SkipSelectionOption {
		return "", nil
	}
	return
}

// containsString checks is a string is present inside a slice
func containsString(s []string, input string) bool {
	for _, a := range s {
		if a == input {
			return true
		}
	}
	return false
}

// Asks for true/false value in the command line
func GetBool(input Input) (a bool, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(bool)
	if !ok {
		dflt = false
	}
	prompt := &survey.Confirm{
		Message: fmt.Sprintf("%s:", input.Question),
		Help:    input.Help,
		Default: dflt,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &a, survey.WithValidator(compose(input.Validators)))
	return
}

// Asks for CIDR value in the command line
func GetIPNet(input Input) (a net.IPNet, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(net.IPNet)
	if !ok {
		dflt = net.IPNet{}
	}
	dfltStr := dflt.String()
	if dfltStr == "<nil>" {
		dfltStr = ""
	}
	question := input.Question
	if !input.Required && dfltStr == "" {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Default: dfltStr,
	}
	var str string
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &str, survey.WithValidator(compose(input.Validators)), survey.WithValidator(IsCIDR))
	if err != nil {
		return
	}
	if str == "" {
		return
	}
	_, cidr, err := net.ParseCIDR(str)
	if err != nil {
		return
	}
	if cidr != nil {
		a = *cidr
	}
	return
}

// Gets password input from the command line
func GetPassword(input Input) (a string, err error) {
	core.DisableColor = !color.UseColor()
	question := input.Question
	if !input.Required {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Password{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &a, survey.WithValidator(compose(input.Validators)))
	return
}

// Gets path to certificate file from the command line
func GetCert(input Input) (a string, err error) {
	core.DisableColor = !color.UseColor()
	dflt, ok := input.Default.(string)
	if !ok {
		dflt = ""
	}
	question := input.Question
	if !input.Required && dflt == "" {
		question = fmt.Sprintf("%s (optional)", question)
	}
	prompt := &survey.Input{
		Message: fmt.Sprintf("%s:", question),
		Help:    input.Help,
		Default: dflt,
	}
	if input.Required {
		input.Validators = append([]Validator{required}, input.Validators...)
	}
	err = survey.AskOne(prompt, &a, survey.WithValidator(compose(input.Validators)), survey.WithValidator(IsCert))
	return
}

var helpTemplate = `{{color "cyan"}}? {{.Message}}
{{range .Steps}}  - {{.}}{{"\n"}}{{end}}{{color "reset"}}`

type Help struct {
	Message string
	Steps   []string
}

func PrintHelp(help Help) error {
	core.DisableColor = !color.UseColor()
	out, _, err := core.RunTemplate(helpTemplate, help)
	if err != nil {
		return err
	}

	fmt.Fprint(terminal.NewAnsiStdout(os.Stdout), out)
	return nil
}
