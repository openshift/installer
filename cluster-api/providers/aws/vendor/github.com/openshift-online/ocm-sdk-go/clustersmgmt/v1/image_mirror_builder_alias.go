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

// IMPORTANT: This file has been generated automatically, refrain from modifying it manually as all
// your changes will be lost when the file is generated again.

package v1 // github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
)

// ImageMirrorBuilder contains the data and logic needed to build 'image_mirror' objects.
//
// ImageMirror represents a container image mirror configuration for a cluster.
// This enables Day 2 image mirroring configuration for ROSA HCP clusters using
// HyperShift's native imageContentSources mechanism.
type ImageMirrorBuilder = api_v1.ImageMirrorBuilder

// NewImageMirror creates a new builder of 'image_mirror' objects.
var NewImageMirror = api_v1.NewImageMirror
