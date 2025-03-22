/*
Copyright 2021 The Kubernetes Authors.

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

package addons

import (
	"github.com/google/go-cmp/cmp"

	infrav1 "sigs.k8s.io/cluster-api-provider-aws/v2/api/v1beta2"
)

// EKSAddon represents an EKS addon.
type EKSAddon struct {
	Name                  *string
	Version               *string
	ServiceAccountRoleARN *string
	Configuration         *string
	Tags                  infrav1.Tags
	ResolveConflict       *string
	ARN                   *string
	Status                *string
}

// IsEqual determines if 2 EKSAddon are equal.
func (e *EKSAddon) IsEqual(other *EKSAddon, includeTags bool) bool {
	//NOTE: we do not compare the ARN as that is only for existing addons
	if e == other {
		return true
	}
	if !cmp.Equal(e.Version, other.Version) {
		return false
	}
	if !cmp.Equal(e.ServiceAccountRoleARN, other.ServiceAccountRoleARN) {
		return false
	}
	if !cmp.Equal(e.Configuration, other.Configuration) {
		return false
	}
	if !cmp.Equal(e.ResolveConflict, other.ResolveConflict) {
		return false
	}

	if includeTags {
		diffTags := e.Tags.Difference(other.Tags)
		if len(diffTags) > 0 {
			return false
		}
	}

	return true
}
