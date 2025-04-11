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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalWifRole writes a value of the 'wif_role' type to the given writer.
func MarshalWifRole(object *WifRole, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeWifRole(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeWifRole writes a value of the 'wif_role' type to the given stream.
func writeWifRole(object *WifRole, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.permissions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("permissions")
		writeStringList(object.permissions, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("predefined")
		stream.WriteBool(object.predefined)
		count++
	}
	present_ = object.bitmap_&4 != 0
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
	object = readWifRole(iterator)
	err = iterator.Error
	return
}

// readWifRole reads a value of the 'wif_role' type from the given iterator.
func readWifRole(iterator *jsoniter.Iterator) *WifRole {
	object := &WifRole{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "permissions":
			value := readStringList(iterator)
			object.permissions = value
			object.bitmap_ |= 1
		case "predefined":
			value := iterator.ReadBool()
			object.predefined = value
			object.bitmap_ |= 2
		case "role_id":
			value := iterator.ReadString()
			object.roleId = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
