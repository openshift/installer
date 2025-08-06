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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalRoleDefinitionOperatorIdentityRequirementList writes a list of values of the 'role_definition_operator_identity_requirement' type to
// the given writer.
func MarshalRoleDefinitionOperatorIdentityRequirementList(list []*RoleDefinitionOperatorIdentityRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteRoleDefinitionOperatorIdentityRequirementList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteRoleDefinitionOperatorIdentityRequirementList writes a list of value of the 'role_definition_operator_identity_requirement' type to
// the given stream.
func WriteRoleDefinitionOperatorIdentityRequirementList(list []*RoleDefinitionOperatorIdentityRequirement, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteRoleDefinitionOperatorIdentityRequirement(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalRoleDefinitionOperatorIdentityRequirementList reads a list of values of the 'role_definition_operator_identity_requirement' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalRoleDefinitionOperatorIdentityRequirementList(source interface{}) (items []*RoleDefinitionOperatorIdentityRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadRoleDefinitionOperatorIdentityRequirementList(iterator)
	err = iterator.Error
	return
}

// ReadRoleDefinitionOperatorIdentityRequirementList reads list of values of the ”role_definition_operator_identity_requirement' type from
// the given iterator.
func ReadRoleDefinitionOperatorIdentityRequirementList(iterator *jsoniter.Iterator) []*RoleDefinitionOperatorIdentityRequirement {
	list := []*RoleDefinitionOperatorIdentityRequirement{}
	for iterator.ReadArray() {
		item := ReadRoleDefinitionOperatorIdentityRequirement(iterator)
		list = append(list, item)
	}
	return list
}
