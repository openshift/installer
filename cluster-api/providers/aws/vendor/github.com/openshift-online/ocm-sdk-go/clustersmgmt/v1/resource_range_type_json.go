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

// MarshalResourceRange writes a value of the 'resource_range' type to the given writer.
func MarshalResourceRange(object *ResourceRange, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeResourceRange(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeResourceRange writes a value of the 'resource_range' type to the given stream.
func writeResourceRange(object *ResourceRange, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max")
		stream.WriteInt(object.max)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("min")
		stream.WriteInt(object.min)
	}
	stream.WriteObjectEnd()
}

// UnmarshalResourceRange reads a value of the 'resource_range' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalResourceRange(source interface{}) (object *ResourceRange, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readResourceRange(iterator)
	err = iterator.Error
	return
}

// readResourceRange reads a value of the 'resource_range' type from the given iterator.
func readResourceRange(iterator *jsoniter.Iterator) *ResourceRange {
	object := &ResourceRange{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "max":
			value := iterator.ReadInt()
			object.max = value
			object.bitmap_ |= 1
		case "min":
			value := iterator.ReadInt()
			object.min = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
