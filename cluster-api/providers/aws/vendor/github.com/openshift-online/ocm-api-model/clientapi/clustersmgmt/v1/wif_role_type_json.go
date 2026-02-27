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

// MarshalWifRole writes a value of the 'wif_role' type to the given writer.
func MarshalWifRole(object *WifRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifRole writes a value of the 'wif_role' type to the given stream.
func WriteWifRole(object *WifRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.permissions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("permissions")
		WriteStringList(object.permissions, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("predefined")
		stream.WriteBool(object.predefined)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.resourceBindings != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource_bindings")
		WriteWifResourceBindingList(object.resourceBindings, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("role_id")
		stream.WriteString(object.roleId)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifRole reads a value of the 'wif_role' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifRole(source interface{}) (object *WifRole, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifRole(iterator)
	err = iterator.Error
	return
}

// ReadWifRole reads a value of the 'wif_role' type from the given iterator.
func ReadWifRole(iterator *jsoniter.Iterator) *WifRole {
	object := &WifRole{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "permissions":
			value := ReadStringList(iterator)
			object.permissions = value
			object.fieldSet_[0] = true
		case "predefined":
			value := iterator.ReadBool()
			object.predefined = value
			object.fieldSet_[1] = true
		case "resource_bindings":
			value := ReadWifResourceBindingList(iterator)
			object.resourceBindings = value
			object.fieldSet_[2] = true
		case "role_id":
			value := iterator.ReadString()
			object.roleId = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
