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

// MarshalDataPlaneOperatorIdentityRequirement writes a value of the 'data_plane_operator_identity_requirement' type to the given writer.
func MarshalDataPlaneOperatorIdentityRequirement(object *DataPlaneOperatorIdentityRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteDataPlaneOperatorIdentityRequirement(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteDataPlaneOperatorIdentityRequirement writes a value of the 'data_plane_operator_identity_requirement' type to the given stream.
func WriteDataPlaneOperatorIdentityRequirement(object *DataPlaneOperatorIdentityRequirement, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_openshift_version")
		stream.WriteString(object.maxOpenShiftVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("min_openshift_version")
		stream.WriteString(object.minOpenShiftVersion)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_name")
		stream.WriteString(object.operatorName)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("required")
		stream.WriteString(object.required)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.roleDefinitions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_definitions")
		WriteRoleDefinitionOperatorIdentityRequirementList(object.roleDefinitions, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.serviceAccounts != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_accounts")
		WriteK8sServiceAccountOperatorIdentityRequirementList(object.serviceAccounts, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalDataPlaneOperatorIdentityRequirement reads a value of the 'data_plane_operator_identity_requirement' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalDataPlaneOperatorIdentityRequirement(source interface{}) (object *DataPlaneOperatorIdentityRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadDataPlaneOperatorIdentityRequirement(iterator)
	err = iterator.Error
	return
}

// ReadDataPlaneOperatorIdentityRequirement reads a value of the 'data_plane_operator_identity_requirement' type from the given iterator.
func ReadDataPlaneOperatorIdentityRequirement(iterator *jsoniter.Iterator) *DataPlaneOperatorIdentityRequirement {
	object := &DataPlaneOperatorIdentityRequirement{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "max_openshift_version":
			value := iterator.ReadString()
			object.maxOpenShiftVersion = value
			object.fieldSet_[0] = true
		case "min_openshift_version":
			value := iterator.ReadString()
			object.minOpenShiftVersion = value
			object.fieldSet_[1] = true
		case "operator_name":
			value := iterator.ReadString()
			object.operatorName = value
			object.fieldSet_[2] = true
		case "required":
			value := iterator.ReadString()
			object.required = value
			object.fieldSet_[3] = true
		case "role_definitions":
			value := ReadRoleDefinitionOperatorIdentityRequirementList(iterator)
			object.roleDefinitions = value
			object.fieldSet_[4] = true
		case "service_accounts":
			value := ReadK8sServiceAccountOperatorIdentityRequirementList(iterator)
			object.serviceAccounts = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
