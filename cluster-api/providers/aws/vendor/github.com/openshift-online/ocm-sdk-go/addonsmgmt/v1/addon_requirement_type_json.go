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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalAddonRequirement writes a value of the 'addon_requirement' type to the given writer.
func MarshalAddonRequirement(object *AddonRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonRequirement(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonRequirement writes a value of the 'addon_requirement' type to the given stream.
func WriteAddonRequirement(object *AddonRequirement, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.data != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("data")
		if object.data != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.data))
			i := 0
			for key := range object.data {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.data[key]
				stream.WriteObjectField(key)
				stream.WriteVal(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource")
		stream.WriteString(string(object.resource))
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteAddonRequirementStatus(object.status, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonRequirement reads a value of the 'addon_requirement' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonRequirement(source interface{}) (object *AddonRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonRequirement(iterator)
	err = iterator.Error
	return
}

// ReadAddonRequirement reads a value of the 'addon_requirement' type from the given iterator.
func ReadAddonRequirement(iterator *jsoniter.Iterator) *AddonRequirement {
	object := &AddonRequirement{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.bitmap_ |= 1
		case "data":
			value := map[string]interface{}{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				var item interface{}
				iterator.ReadVal(&item)
				value[key] = item
			}
			object.data = value
			object.bitmap_ |= 2
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 4
		case "resource":
			text := iterator.ReadString()
			value := AddonRequirementResource(text)
			object.resource = value
			object.bitmap_ |= 8
		case "status":
			value := ReadAddonRequirementStatus(iterator)
			object.status = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
