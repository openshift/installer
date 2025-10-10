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

// MarshalAddonParameterOption writes a value of the 'addon_parameter_option' type to the given writer.
func MarshalAddonParameterOption(object *AddonParameterOption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonParameterOption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonParameterOption writes a value of the 'addon_parameter_option' type to the given stream.
func WriteAddonParameterOption(object *AddonParameterOption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rank")
		stream.WriteInt(object.rank)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.requirements != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requirements")
		WriteAddonRequirementList(object.requirements, stream)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.value)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonParameterOption reads a value of the 'addon_parameter_option' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonParameterOption(source interface{}) (object *AddonParameterOption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonParameterOption(iterator)
	err = iterator.Error
	return
}

// ReadAddonParameterOption reads a value of the 'addon_parameter_option' type from the given iterator.
func ReadAddonParameterOption(iterator *jsoniter.Iterator) *AddonParameterOption {
	object := &AddonParameterOption{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 1
		case "rank":
			value := iterator.ReadInt()
			object.rank = value
			object.bitmap_ |= 2
		case "requirements":
			value := ReadAddonRequirementList(iterator)
			object.requirements = value
			object.bitmap_ |= 4
		case "value":
			value := iterator.ReadString()
			object.value = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
