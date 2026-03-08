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

// This file contains functions used to implement the '--interactive' command line option.

package interactive

import (
	"fmt"
	"net"
	"net/url"
	"os"
	"regexp"
	"strconv"
	"strings"

	"k8s.io/apimachinery/pkg/util/validation"
	netutils "k8s.io/utils/net"

	"github.com/AlecAivazis/survey/v2"
	"github.com/AlecAivazis/survey/v2/core"
	clustervalidations "github.com/openshift-online/ocm-common/pkg/cluster/validations"
	diskValidator "github.com/openshift-online/ocm-common/pkg/machinepool/validations"

	"github.com/openshift/rosa/pkg/aws"
	"github.com/openshift/rosa/pkg/helper"
	"github.com/openshift/rosa/pkg/ocm"
)

const doubleQuotesToRemove = "\"\""

type Validator survey.Validator

var required = survey.Required

var MaxLength = func(length int) Validator {
	return Validator(survey.MaxLength(length))
}

func compose(validators []Validator) survey.Validator {
	surveyValidators := []survey.Validator{}
	for _, validator := range validators {
		surveyValidators = append(surveyValidators, survey.Validator(validator))
	}
	return survey.ComposeValidators(surveyValidators...)
}

// IsURL validates whether the given value is a valid URL
func IsURL(val interface{}) error {
	_, err := _isUrl(val)
	return err
}

// IsValidHostname is same validation as in the Open Shift GitHub IDP CRD
// Validates the hostname DNS1123 format
// https://github.com/openshift/kubernetes/blob/91607f5d750ba4002f87d34a12ae1cfd45b45b81/openshift-kube-apiserver/admission/customresourcevalidation/oauth/helpers.go#L13
// and denies any [*.]github.com hostnames
// https://github.com/openshift/kubernetes/blob/258f1d5fb6491ba65fd8201c827e179432430627/openshift-kube-apiserver/admission/customresourcevalidation/oauth/validate_github.go#L49
// nolint
func IsValidHostname(val interface{}) error {
	hostname := val.(string)
	if hostname == "" {
		return nil
	}
	if hostname == "github.com" || strings.HasSuffix(hostname, ".github.com") {
		return fmt.Errorf("%s", fmt.Sprintf("'%s' hostname cannot be equal to [*.]github.com", hostname))
	}
	if !(len(validation.IsDNS1123Subdomain(hostname)) == 0 || netutils.ParseIPSloppy(hostname) != nil) {
		return fmt.Errorf("%s", fmt.Sprintf("'%s' hostname must be a valid DNS subdomain or IP address", hostname))
	}
	return nil
}

func IsURLHttps(val interface{}) error {
	parsedUri, err := _isUrl(val)
	if err != nil {
		return err
	}
	if parsedUri.Scheme != helper.ProtocolHttps {
		return fmt.Errorf("expect URL '%s' to use an 'https://' scheme", val.(string))
	}
	return nil
}

func _isUrl(val interface{}) (*url.URL, error) {
	if val == nil {
		return nil, nil
	}
	s, ok := val.(string)
	if !ok {
		return nil, fmt.Errorf("can only validate strings, got %v", val)
	}
	if s == "" {
		return nil, nil
	}
	parsedUri, err := url.ParseRequestURI(fmt.Sprintf("%v", val))
	return parsedUri, err
}

// IsCert validates whether the given filepath is a valid cert file
func IsCert(filepath interface{}) error {
	if filepath == nil {
		return nil
	}
	if s, ok := filepath.(string); ok {
		if s == "" {
			return nil
		}
		if s == doubleQuotesToRemove {
			return nil
		}
		validExtension, err := regexp.MatchString("\\.(pem|ca-bundle|ce?rt?|key)$", s)
		if err != nil {
			return err
		}
		if !validExtension {
			return fmt.Errorf("file '%s' does not have a valid file extension", s)
		}
		if _, err := os.Stat(s); !os.IsNotExist(err) {
			// path to file exist
			return nil
		}
		return fmt.Errorf("file '%s' does not exist on the file system", s)
	}
	return fmt.Errorf("can only validate strings, got %v", filepath)
}

func IsCIDR(val interface{}) error {
	if s, ok := val.(string); ok {
		_, _, err := net.ParseCIDR(s)
		if err != nil {
			return err
		}
		return nil
	}
	return fmt.Errorf("can only validate strings, got %v", val)
}

// MaxValue returns a Validator that validates the entered number is less than or equal to max
func MaxValue(max int) Validator {
	return func(ans interface{}) error {
		if str, ok := ans.(string); ok {
			val, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("please enter an integer value, you entered '%s'", str)
			}
			if val > max {
				return fmt.Errorf(
					"'%d' is greater than the permitted maximum of '%d'", val, max)
			}

			return nil
		}

		return fmt.Errorf("can only validate strings, got %v", ans)
	}
}

// MinValue returns a validator that validates the entered number is greater than or equal to min
func MinValue(min int) Validator {
	return func(ans interface{}) error {
		if str, ok := ans.(string); ok {
			val, err := strconv.Atoi(str)
			if err != nil {
				return fmt.Errorf("please enter an integer value, you entered '%s'", str)
			}
			if val < min {
				return fmt.Errorf(
					"'%d' is less than the permitted minimum of '%d'", val, min)
			}

			return nil
		}

		return fmt.Errorf("can only validate strings, got %v", ans)
	}
}

func RegExp(restr string) Validator {
	re := regexp.MustCompile(restr)
	return func(val interface{}) error {
		if str, ok := val.(string); ok {
			if str == "" {
				return nil
			}
			if !re.MatchString(str) {
				return fmt.Errorf("%s does not match regular expression %s", str, re.String())
			}
			return nil
		}
		return fmt.Errorf("can only validate strings, got %v", val)
	}
}

func RegExpBoolean(restr string) Validator {
	re := regexp.MustCompile(restr)
	return func(val interface{}) error {
		if boolVal, ok := val.(bool); ok {
			var val string
			if boolVal {
				val = "true"
			} else {
				val = "false"
			}
			if !re.MatchString(val) {
				return fmt.Errorf("%s does not match regular expression %s", val, re.String())
			}
			return nil
		}
		return fmt.Errorf("can only validate boolean values, got %v", val)
	}
}

// SubnetsCountValidator get a slice of `[]core.OptionAnswer` as an interface.
// e.g. core.OptionAnswer { Value: subnet-04f67939f44a97dbe (us-west-2b), Index: 0 }
func SubnetsValidator(awsClient aws.Client, multiAZ bool, privateValue bool, hostedCP bool,
	privateIngress bool) Validator {
	return func(input interface{}) (err error) {
		if answers, ok := input.([]core.OptionAnswer); ok {
			if hostedCP {
				subnetIDs := make([]string, len(answers))
				for i, subnet := range answers {
					subnetIDs[i] = aws.ParseOption(subnet.Value)
				}

				_, err = ocm.ValidateHostedClusterSubnets(awsClient, privateValue, subnetIDs, privateIngress)
				return err
			}
			return ocm.ValidateSubnetsCount(multiAZ, privateValue, len(answers))
		}
		return fmt.Errorf("can only validate a slice of string, got %v", input)
	}
}

func AvailabilityZonesCountValidator(multiAZ bool) Validator {
	return func(input interface{}) error {
		if answers, ok := input.([]core.OptionAnswer); ok {
			return clustervalidations.ValidateAvailabilityZonesCount(multiAZ, len(answers))
		}

		return fmt.Errorf("can only validate a slice of string, got %v", input)
	}
}

func MachinePoolRootDiskSizeValidator(version string) Validator {
	return func(val interface{}) error {
		// We expect GigiByte as the unit for the root volume size

		// Validate the worker root volume size is an integer
		machinePoolRootDiskSize, ok := val.(string)
		if !ok {
			return fmt.Errorf("machine pool root disk size must be an string, got %T", machinePoolRootDiskSize)
		}

		// parse it to validate it is a valid unit
		size, err := ocm.ParseDiskSizeToGigibyte(machinePoolRootDiskSize)
		if err != nil {
			return fmt.Errorf("failed to parse machine pool root disk size '%s': %v", machinePoolRootDiskSize, err)
		}
		return diskValidator.ValidateMachinePoolRootDiskSize(version, size)
	}
}

func NodePoolRootDiskSizeValidator() Validator {
	return func(val interface{}) error {
		// We expect GigiByte as the unit for the root volume size

		// Validate the worker root volume size is an integer
		nodePoolRootDiskSize, ok := val.(string)
		if !ok {
			return fmt.Errorf("node pool root disk size must be an string, got %T", nodePoolRootDiskSize)
		}

		// parse it to validate it is a valid unit
		size, err := ocm.ParseDiskSizeToGigibyte(nodePoolRootDiskSize)
		if err != nil {
			return fmt.Errorf("failed to parse machine pool root disk size '%s': %v", nodePoolRootDiskSize, err)
		}
		return diskValidator.ValidateNodePoolRootDiskSize(size)
	}
}
