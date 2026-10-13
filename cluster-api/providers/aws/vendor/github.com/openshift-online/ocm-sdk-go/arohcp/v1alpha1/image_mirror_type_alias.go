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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	api_v1alpha1 "github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1"
)

// ImageMirrorKind is the name of the type used to represent objects
// of type 'image_mirror'.
const ImageMirrorKind = api_v1alpha1.ImageMirrorKind

// ImageMirrorLinkKind is the name of the type used to represent links
// to objects of type 'image_mirror'.
const ImageMirrorLinkKind = api_v1alpha1.ImageMirrorLinkKind

// ImageMirrorNilKind is the name of the type used to nil references
// to objects of type 'image_mirror'.
const ImageMirrorNilKind = api_v1alpha1.ImageMirrorNilKind

// ImageMirror represents the values of the 'image_mirror' type.
//
// ImageMirror represents a container image mirror configuration for a cluster.
// This enables Day 2 image mirroring configuration for ROSA HCP clusters using
// HyperShift's native imageContentSources mechanism.
type ImageMirror = api_v1alpha1.ImageMirror

// ImageMirrorListKind is the name of the type used to represent list of objects of
// type 'image_mirror'.
const ImageMirrorListKind = api_v1alpha1.ImageMirrorListKind

// ImageMirrorListLinkKind is the name of the type used to represent links to list
// of objects of type 'image_mirror'.
const ImageMirrorListLinkKind = api_v1alpha1.ImageMirrorListLinkKind

// ImageMirrorNilKind is the name of the type used to nil lists of objects of
// type 'image_mirror'.
const ImageMirrorListNilKind = api_v1alpha1.ImageMirrorListNilKind

type ImageMirrorList = api_v1alpha1.ImageMirrorList
