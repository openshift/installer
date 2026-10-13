package interactive

import (
	"fmt"
	"net"
	"regexp"
	"strconv"
	"strings"

	asv1 "github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1"
)

func GetAddonArgument(param asv1.AddonParameter, dflt string) (string, error) {
	var input = Input{
		Question: param.Name(),
		Help:     fmt.Sprintf("%s: %s", param.ID(), param.Description()),
		Required: param.Required(),
		Options:  getOptionNames(param),
		Default:  dflt,
	}

	if len(input.Options) > 0 {
		optionName, err := GetOption(input)
		if err != nil {
			return "", fmt.Errorf("expected a valid option for '%s': %v", param.ID(), err)
		}

		var optionValue string
		for _, paramOption := range param.Options() {
			if strings.Compare(paramOption.Name(), optionName) == 0 {
				optionValue = paramOption.Value()
				break
			}
		}
		return optionValue, nil
	}

	if param.Validation() != "" {
		input.Validators = []Validator{
			func(ans interface{}) error {
				strAns := ans.(string)
				if strAns == input.Default {
					return nil
				}
				if isValid, err := regexp.MatchString(param.Validation(), strAns); err != nil || !isValid {
					if param.ValidationErrMsg() != "" {
						return fmt.Errorf("%s", param.ValidationErrMsg())
					}
					return fmt.Errorf("expected %q to match /%s/", strAns, param.Validation())
				}
				return nil
			},
		}
	}

	switch param.ValueType() {
	case "boolean":
		var boolVal bool
		input.Default, _ = strconv.ParseBool(dflt)
		// Set default value based on existing parameter, otherwise use parameter default
		// add a prompt to question name to indicate if the boolean param is required and check validation
		if param.Validation() == "^true$" && param.Required() {
			input.Question = fmt.Sprintf("%s (required)", param.Name())
			input.Validators = []Validator{
				RegExpBoolean(param.Validation()),
			}
		}
		boolVal, err := GetBool(input)
		if err != nil {
			return "", fmt.Errorf("expected a valid boolean value for '%s': %v", param.ID(), err)
		}
		if boolVal {
			return "true", nil
		}
		return "false", nil
	case "cidr":
		var cidrVal net.IPNet
		if dflt != "" {
			_, defaultIDR, _ := net.ParseCIDR(dflt)
			input.Default = *defaultIDR
		}
		cidrVal, err := GetIPNet(input)
		if err != nil {
			return "", fmt.Errorf("expected a valid CIDR value for '%s': %v", param.ID(), err)
		}
		if cidrVal.String() == "<nil>" {
			return "", nil
		}
		return cidrVal.String(), nil
	case "number", "resource":
		input.Default, _ = strconv.Atoi(dflt)
		numVal, err := GetInt(input)
		if err != nil {
			return "", fmt.Errorf("expected a valid numerical value for '%s': %v", param.ID(), err)
		}
		return fmt.Sprintf("%d", numVal), nil

	case "string":
		input.Default = dflt
		value, err := GetString(input)
		if err != nil {
			return "", fmt.Errorf("expected a valid string value for '%s': %v", param.ID(), err)
		}
		return value, nil
	}

	// If the parameter type wasn't handled in the above switch statement
	// then this parameter type is not supported by interactive mode yet.
	return "", fmt.Errorf("the parameter '%s' does not support interactive mode", param.ID())
}

func getOptionNames(param asv1.AddonParameter) []string {
	var optionNames []string
	options, _ := param.GetOptions()
	for _, option := range options {
		optionNames = append(optionNames, option.Name())
	}
	return optionNames
}
