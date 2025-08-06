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

// MarshalTokenIssuer writes a value of the 'token_issuer' type to the given writer.
func MarshalTokenIssuer(object *TokenIssuer, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteTokenIssuer(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteTokenIssuer writes a value of the 'token_issuer' type to the given stream.
func WriteTokenIssuer(object *TokenIssuer, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ca")
		stream.WriteString(object.ca)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("url")
		stream.WriteString(object.url)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.audiences != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("audiences")
		WriteStringList(object.audiences, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalTokenIssuer reads a value of the 'token_issuer' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalTokenIssuer(source interface{}) (object *TokenIssuer, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadTokenIssuer(iterator)
	err = iterator.Error
	return
}

// ReadTokenIssuer reads a value of the 'token_issuer' type from the given iterator.
func ReadTokenIssuer(iterator *jsoniter.Iterator) *TokenIssuer {
	object := &TokenIssuer{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "ca":
			value := iterator.ReadString()
			object.ca = value
			object.bitmap_ |= 1
		case "url":
			value := iterator.ReadString()
			object.url = value
			object.bitmap_ |= 2
		case "audiences":
			value := ReadStringList(iterator)
			object.audiences = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
