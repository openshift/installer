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

// AddOnConfigKind is the name of the type used to represent objects
// of type 'add_on_config'.
const AddOnConfigKind = api_v1alpha1.AddOnConfigKind

// AddOnConfigLinkKind is the name of the type used to represent links
// to objects of type 'add_on_config'.
const AddOnConfigLinkKind = api_v1alpha1.AddOnConfigLinkKind

// AddOnConfigNilKind is the name of the type used to nil references
// to objects of type 'add_on_config'.
const AddOnConfigNilKind = api_v1alpha1.AddOnConfigNilKind

// AddOnConfig represents the values of the 'add_on_config' type.
//
// Representation of an add-on config.
// The attributes under it are to be used by the addon once its installed in the cluster.
type AddOnConfig = api_v1alpha1.AddOnConfig

// AddOnConfigListKind is the name of the type used to represent list of objects of
// type 'add_on_config'.
const AddOnConfigListKind = api_v1alpha1.AddOnConfigListKind

// AddOnConfigListLinkKind is the name of the type used to represent links to list
// of objects of type 'add_on_config'.
const AddOnConfigListLinkKind = api_v1alpha1.AddOnConfigListLinkKind

// AddOnConfigNilKind is the name of the type used to nil lists of objects of
// type 'add_on_config'.
const AddOnConfigListNilKind = api_v1alpha1.AddOnConfigListNilKind

type AddOnConfigList = api_v1alpha1.AddOnConfigList
