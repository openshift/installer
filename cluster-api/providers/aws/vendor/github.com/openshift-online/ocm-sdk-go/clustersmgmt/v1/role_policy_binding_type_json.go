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
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalRolePolicyBinding writes a value of the 'role_policy_binding' type to the given writer.
func MarshalRolePolicyBinding(object *RolePolicyBinding, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteRolePolicyBinding(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteRolePolicyBinding writes a value of the 'role_policy_binding' type to the given stream.
func WriteRolePolicyBinding(object *RolePolicyBinding, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("arn")
		stream.WriteString(object.arn)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.policies != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("policies")
		WriteRolePolicyList(object.policies, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteRolePolicyBindingStatus(object.status, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRolePolicyBinding reads a value of the 'role_policy_binding' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRolePolicyBinding(source interface{}) (object *RolePolicyBinding, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadRolePolicyBinding(iterator)
	err = iterator.Error
	return
}

// ReadRolePolicyBinding reads a value of the 'role_policy_binding' type from the given iterator.
func ReadRolePolicyBinding(iterator *jsoniter.Iterator) *RolePolicyBinding {
	object := &RolePolicyBinding{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "arn":
			value := iterator.ReadString()
			object.arn = value
			object.bitmap_ |= 1
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 2
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.bitmap_ |= 4
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 8
		case "policies":
			value := ReadRolePolicyList(iterator)
			object.policies = value
			object.bitmap_ |= 16
		case "status":
			value := ReadRolePolicyBindingStatus(iterator)
			object.status = value
			object.bitmap_ |= 32
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.bitmap_ |= 64
		default:
			iterator.ReadAny()
		}
	}
	return object
}
