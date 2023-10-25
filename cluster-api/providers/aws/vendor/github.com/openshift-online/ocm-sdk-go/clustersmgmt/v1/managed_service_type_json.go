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

// MarshalManagedService writes a value of the 'managed_service' type to the given writer.
func MarshalManagedService(object *ManagedService, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeManagedService(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeManagedService writes a value of the 'managed_service' type to the given stream.
func writeManagedService(object *ManagedService, stream *jsoniter.Stream) {
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

// UnmarshalManagedService reads a value of the 'managed_service' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalManagedService(source interface{}) (object *ManagedService, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readManagedService(iterator)
	err = iterator.Error
	return
}

// readManagedService reads a value of the 'managed_service' type from the given iterator.
func readManagedService(iterator *jsoniter.Iterator) *ManagedService {
	object := &ManagedService{}
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
