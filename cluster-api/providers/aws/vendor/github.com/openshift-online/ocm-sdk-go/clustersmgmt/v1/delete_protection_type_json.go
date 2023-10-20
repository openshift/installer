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

// MarshalDeleteProtection writes a value of the 'delete_protection' type to the given writer.
func MarshalDeleteProtection(object *DeleteProtection, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeDeleteProtection(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeDeleteProtection writes a value of the 'delete_protection' type to the given stream.
func writeDeleteProtection(object *DeleteProtection, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
	}
	stream.WriteObjectEnd()
}

// UnmarshalDeleteProtection reads a value of the 'delete_protection' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalDeleteProtection(source interface{}) (object *DeleteProtection, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readDeleteProtection(iterator)
	err = iterator.Error
	return
}

// readDeleteProtection reads a value of the 'delete_protection' type from the given iterator.
func readDeleteProtection(iterator *jsoniter.Iterator) *DeleteProtection {
	object := &DeleteProtection{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
