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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalOperatorIAMRole writes a value of the 'operator_IAM_role' type to the given writer.
func MarshalOperatorIAMRole(object *OperatorIAMRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteOperatorIAMRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOperatorIAMRole writes a value of the 'operator_IAM_role' type to the given stream.
func WriteOperatorIAMRole(object *OperatorIAMRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleARN)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("service_account")
		stream.WriteString(object.serviceAccount)
	}
	stream.WriteObjectEnd()
}

// UnmarshalOperatorIAMRole reads a value of the 'operator_IAM_role' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalOperatorIAMRole(source interface{}) (object *OperatorIAMRole, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadOperatorIAMRole(iterator)
	err = iterator.Error
	return
}

// ReadOperatorIAMRole reads a value of the 'operator_IAM_role' type from the given iterator.
func ReadOperatorIAMRole(iterator *jsoniter.Iterator) *OperatorIAMRole {
	object := &OperatorIAMRole{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[0] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[1] = true
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.fieldSet_[2] = true
		case "role_arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.fieldSet_[3] = true
		case "service_account":
			value := iterator.ReadString()
			object.serviceAccount = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
