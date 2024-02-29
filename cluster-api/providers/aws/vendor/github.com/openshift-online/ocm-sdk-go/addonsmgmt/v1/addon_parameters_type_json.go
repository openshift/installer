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

package v1 // github.com/openshift-online/ocm-sdk-go/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonParameters writes a value of the 'addon_parameters' type to the given writer.
func MarshalAddonParameters(object *AddonParameters, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeAddonParameters(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeAddonParameters writes a value of the 'addon_parameters' type to the given stream.
func writeAddonParameters(object *AddonParameters, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.items != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("items")
		writeAddonParameterList(object.items, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonParameters reads a value of the 'addon_parameters' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonParameters(source interface{}) (object *AddonParameters, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readAddonParameters(iterator)
	err = iterator.Error
	return
}

// readAddonParameters reads a value of the 'addon_parameters' type from the given iterator.
func readAddonParameters(iterator *jsoniter.Iterator) *AddonParameters {
	object := &AddonParameters{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "items":
			value := readAddonParameterList(iterator)
			object.items = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
