// Copyright 2024 Google LLC. All Rights Reserved.
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
// Package dataproc contains methods and types for managing dataproc GCP resources.
package dataproc

import (
	"github.com/GoogleCloudPlatform/declarative-resource-client-library/dcl"
)

func encodeJobCreateRequest(m map[string]any) map[string]any {
	req := make(map[string]any, 1)
	dcl.PutMapEntry(req, []string{"job"}, m)
	return req
}

func expandClusterProject(_ *Client, project *string, _ *Cluster) (*string, error) {
	return dcl.SelfLinkToName(project), nil
}

// CompareClusterConfigMasterConfigNewStyle exposes the compareClusterConfigMasterConfigNewStyle function for testing.
func CompareClusterConfigMasterConfigNewStyle(d, a any, fn dcl.FieldName) ([]*dcl.FieldDiff, error) {
	return compareClusterConfigMasterConfigNewStyle(d, a, fn)
}

func canonicalizeSoftwareConfigProperties(o, n any) bool {
	// This field is a map that contains both client provided and server provided values. It
	// is also immutable, so always return "no diff".
	return true
}
