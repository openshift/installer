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

// MarshalOpenIDClaims writes a value of the 'open_ID_claims' type to the given writer.
func MarshalOpenIDClaims(object *OpenIDClaims, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeOpenIDClaims(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeOpenIDClaims writes a value of the 'open_ID_claims' type to the given stream.
func writeOpenIDClaims(object *OpenIDClaims, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.email != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("email")
		writeStringList(object.email, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		writeStringList(object.groups, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.name != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		writeStringList(object.name, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.preferredUsername != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("preferred_username")
		writeStringList(object.preferredUsername, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalOpenIDClaims reads a value of the 'open_ID_claims' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalOpenIDClaims(source interface{}) (object *OpenIDClaims, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readOpenIDClaims(iterator)
	err = iterator.Error
	return
}

// readOpenIDClaims reads a value of the 'open_ID_claims' type from the given iterator.
func readOpenIDClaims(iterator *jsoniter.Iterator) *OpenIDClaims {
	object := &OpenIDClaims{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "email":
			value := readStringList(iterator)
			object.email = value
			object.bitmap_ |= 1
		case "groups":
			value := readStringList(iterator)
			object.groups = value
			object.bitmap_ |= 2
		case "name":
			value := readStringList(iterator)
			object.name = value
			object.bitmap_ |= 4
		case "preferred_username":
			value := readStringList(iterator)
			object.preferredUsername = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
