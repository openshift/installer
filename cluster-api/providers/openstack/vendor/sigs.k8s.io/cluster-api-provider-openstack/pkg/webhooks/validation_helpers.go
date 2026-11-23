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

package webhooks

import (
	"k8s.io/apimachinery/pkg/util/validation/field"
)

const rootVolumeName = "root"

// validateSecurityGroupRulesRemoteMutualExclusion validates that remote* fields
// are mutually exclusive in security group rules. The getRemoteFields function
// extracts whether each remote field is set from a rule of any version.
func validateSecurityGroupRulesRemoteMutualExclusion[T any](
	rules []T,
	fldPath *field.Path,
	getRemoteFields func(*T) (hasRemoteManagedGroups, hasRemoteGroupID, hasRemoteIPPrefix bool),
) field.ErrorList {
	var allErrs field.ErrorList
	for i := range rules {
		hasRMG, hasRGID, hasRIP := getRemoteFields(&rules[i])
		count := 0
		if hasRMG {
			count++
		}
		if hasRGID {
			count++
		}
		if hasRIP {
			count++
		}
		if count > 1 {
			allErrs = append(allErrs, field.Forbidden(fldPath.Index(i), "only one of remoteManagedGroups, remoteGroupID, or remoteIPPrefix can be set"))
		}
	}
	return allErrs
}
