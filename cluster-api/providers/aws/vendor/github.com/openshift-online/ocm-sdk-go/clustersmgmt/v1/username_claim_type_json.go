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

// MarshalUsernameClaim writes a value of the 'username_claim' type to the given writer.
func MarshalUsernameClaim(object *UsernameClaim, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteUsernameClaim(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteUsernameClaim writes a value of the 'username_claim' type to the given stream.
func WriteUsernameClaim(object *UsernameClaim, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("claim")
		stream.WriteString(object.claim)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("prefix")
		stream.WriteString(object.prefix)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("prefix_policy")
		stream.WriteString(object.prefixPolicy)
	}
	stream.WriteObjectEnd()
}

// UnmarshalUsernameClaim reads a value of the 'username_claim' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalUsernameClaim(source interface{}) (object *UsernameClaim, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadUsernameClaim(iterator)
	err = iterator.Error
	return
}

// ReadUsernameClaim reads a value of the 'username_claim' type from the given iterator.
func ReadUsernameClaim(iterator *jsoniter.Iterator) *UsernameClaim {
	object := &UsernameClaim{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "claim":
			value := iterator.ReadString()
			object.claim = value
			object.bitmap_ |= 1
		case "prefix":
			value := iterator.ReadString()
			object.prefix = value
			object.bitmap_ |= 2
		case "prefix_policy":
			value := iterator.ReadString()
			object.prefixPolicy = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
