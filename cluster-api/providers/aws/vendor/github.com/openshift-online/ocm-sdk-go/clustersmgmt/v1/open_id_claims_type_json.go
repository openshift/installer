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
	WriteOpenIDClaims(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteOpenIDClaims writes a value of the 'open_ID_claims' type to the given stream.
func WriteOpenIDClaims(object *OpenIDClaims, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.email != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("email")
		WriteStringList(object.email, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.groups != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("groups")
		WriteStringList(object.groups, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.name != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		WriteStringList(object.name, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.preferredUsername != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("preferred_username")
		WriteStringList(object.preferredUsername, stream)
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
	object = ReadOpenIDClaims(iterator)
	err = iterator.Error
	return
}

// ReadOpenIDClaims reads a value of the 'open_ID_claims' type from the given iterator.
func ReadOpenIDClaims(iterator *jsoniter.Iterator) *OpenIDClaims {
	object := &OpenIDClaims{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "email":
			value := ReadStringList(iterator)
			object.email = value
			object.bitmap_ |= 1
		case "groups":
			value := ReadStringList(iterator)
			object.groups = value
			object.bitmap_ |= 2
		case "name":
			value := ReadStringList(iterator)
			object.name = value
			object.bitmap_ |= 4
		case "preferred_username":
			value := ReadStringList(iterator)
			object.preferredUsername = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
