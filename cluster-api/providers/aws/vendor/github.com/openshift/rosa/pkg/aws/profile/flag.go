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

// This file contains functions used to implement the '--profile' command line option.

package profile

import (
	"os"

	"github.com/spf13/pflag"
)

// AddFlag adds the debug flag to the given set of command line flags.
func AddFlag(flags *pflag.FlagSet) {
	flags.StringVar(
		&profile,
		"profile",
		"",
		"Use a specific AWS profile from your credential file.",
	)
}

// Profile returns a string with the name of the AWS profile being used.
func Profile() string {
	if profile != "" {
		return profile
	}
	awsProfile := os.Getenv("AWS_PROFILE")
	if awsProfile != "" {
		return awsProfile
	}
	return ""
}

// profile is a string flag that indicates which AWS profile is being used.
var profile string
