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

package v1 // github.com/openshift-online/ocm-sdk-go/servicemgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalOperatorIAMRole writes a value of the 'operator_IAM_role' type to the given writer.
func MarshalOperatorIAMRole(object *OperatorIAMRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeOperatorIAMRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeOperatorIAMRole writes a value of the 'operator_IAM_role' type to the given stream.
func writeOperatorIAMRole(object *OperatorIAMRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("namespace")
		stream.WriteString(object.namespace)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_arn")
		stream.WriteString(object.roleARN)
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
	object = readOperatorIAMRole(iterator)
	err = iterator.Error
	return
}

// readOperatorIAMRole reads a value of the 'operator_IAM_role' type from the given iterator.
func readOperatorIAMRole(iterator *jsoniter.Iterator) *OperatorIAMRole {
	object := &OperatorIAMRole{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 1
		case "namespace":
			value := iterator.ReadString()
			object.namespace = value
			object.bitmap_ |= 2
		case "role_arn":
			value := iterator.ReadString()
			object.roleARN = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
