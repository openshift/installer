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

package hub

import (
	vmoprv1alpha5common "github.com/vmware-tanzu/vm-operator/api/v1alpha5/common"
)

// PolicySpec and PolicyStatus must stay aligned with vm-operator:
// https://github.com/vmware-tanzu/vm-operator/blob/main/api/v1alpha5/virtualmachine_policy_types.go
// Copy definitions from that file when vm-operator bumps; do not invent parallel shapes here.

type PolicySpec vmoprv1alpha5common.LocalObjectRef

type PolicyStatus struct {
	PolicySpec `json:",inline"`

	// Generation describes the observed generation of the policy applied to
	// this VM.
	Generation int64 `json:"generation"`
}
