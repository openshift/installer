/*
Copyright 2026 The Kubernetes Authors.

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

package awshelpers

import (
	"fmt"

	rosaaws "github.com/openshift/rosa/pkg/aws"
)

// GetNoConsolePolicyName returns the no-console policy name for a role.
// This function is not available in the current vendored ROSA version.
// TODO: Remove once ROSA dependency is bumped to include this function (added in commit e14522438).
func GetNoConsolePolicyName(name string) string {
	return fmt.Sprintf("%s-NoConsole-Policy", name)
}

// GetNoConsolePolicyARN returns the ARN for the no-console policy.
// This function is not available in the current vendored ROSA version.
// TODO: Remove once ROSA dependency is bumped to include this function (added in commit e14522438).
func GetNoConsolePolicyARN(partition, accountID, name, path string) string {
	return rosaaws.GetPolicyArn(partition, accountID, GetNoConsolePolicyName(name), path)
}
