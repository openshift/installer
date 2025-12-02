// Copyright 2023 Google LLC. All Rights Reserved.
// 
// Licensed under the Apache License, Version 2.0 (the "License");
// you may not use this file except in compliance with the License.
// You may obtain a copy of the License at
// 
//     http://www.apache.org/licenses/LICENSE-2.0
// 
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS,
// WITHOUT WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied.
// See the License for the specific language governing permissions and
// limitations under the License.
// Package firebaserules provides Utilities for Firebase Rules custom overrides.
package firebaserules

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

// EncodeReleaseUpdateRequest encapsulates fields in a release {} block, as expected
// by https://firebase.google.com/docs/reference/rules/rest/v1/projects.releases/patch
func EncodeReleaseUpdateRequest(m map[string]interface{}) map[string]interface{} {
	req := make(map[string]interface{})
	dcl.PutMapEntry(req, []string{"release"}, m)
	return req
}
