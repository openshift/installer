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

// This file contains functions used to implement the '--region' command line option.

package region

import (
	"os"

	"github.com/spf13/pflag"

	"github.com/openshift/rosa/pkg/constants"
	"github.com/openshift/rosa/pkg/helper"
)

// AddFlag adds the region flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&region,
		"region",
		"",
		"Use a specific AWS region, overriding the AWS_REGION environment variable.",
	)
}

// Region returns a string with the name of the AWS region being used.
func Region() string {
	if helper.HandleEscapedEmptyString(region) != "" {
		return region
	}
	awsRegion := os.Getenv(constants.AwsRegion)
	if helper.HandleEscapedEmptyString(awsRegion) != "" {
		return awsRegion
	}
	return ""
}

// region is a string flag that indicates which AWS region is being used.
var region string
