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

// MarshalWifSupport writes a value of the 'wif_support' type to the given writer.
func MarshalWifSupport(object *WifSupport, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteWifSupport(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteWifSupport writes a value of the 'wif_support' type to the given stream.
func WriteWifSupport(object *WifSupport, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("principal")
		stream.WriteString(object.principal)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.roles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("roles")
		WriteWifRoleList(object.roles, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalWifSupport reads a value of the 'wif_support' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalWifSupport(source interface{}) (object *WifSupport, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadWifSupport(iterator)
	err = iterator.Error
	return
}

// ReadWifSupport reads a value of the 'wif_support' type from the given iterator.
func ReadWifSupport(iterator *jsoniter.Iterator) *WifSupport {
	object := &WifSupport{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "principal":
			value := iterator.ReadString()
			object.principal = value
			object.fieldSet_[0] = true
		case "roles":
			value := ReadWifRoleList(iterator)
			object.roles = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
