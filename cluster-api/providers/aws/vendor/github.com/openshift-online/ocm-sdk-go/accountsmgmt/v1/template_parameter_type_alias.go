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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	api_v1 "github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1"
)

// TemplateParameter represents the values of the 'template_parameter' type.
//
// A template parameter is used in an email to replace placeholder content with
// values specific to the email recipient.
type TemplateParameter = api_v1.TemplateParameter

// TemplateParameterListKind is the name of the type used to represent list of objects of
// type 'template_parameter'.
const TemplateParameterListKind = api_v1.TemplateParameterListKind

// TemplateParameterListLinkKind is the name of the type used to represent links to list
// of objects of type 'template_parameter'.
const TemplateParameterListLinkKind = api_v1.TemplateParameterListLinkKind

// TemplateParameterNilKind is the name of the type used to nil lists of objects of
// type 'template_parameter'.
const TemplateParameterListNilKind = api_v1.TemplateParameterListNilKind

type TemplateParameterList = api_v1.TemplateParameterList
