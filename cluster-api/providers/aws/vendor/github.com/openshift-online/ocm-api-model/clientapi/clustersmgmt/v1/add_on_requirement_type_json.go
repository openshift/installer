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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddOnRequirement writes a value of the 'add_on_requirement' type to the given writer.
func MarshalAddOnRequirement(object *AddOnRequirement, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnRequirement(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnRequirement writes a value of the 'add_on_requirement' type to the given stream.
func WriteAddOnRequirement(object *AddOnRequirement, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.data != nil
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
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resource")
		stream.WriteString(object.resource)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteAddOnRequirementStatus(object.status, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnRequirement reads a value of the 'add_on_requirement' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnRequirement(source interface{}) (object *AddOnRequirement, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOnRequirement(iterator)
	err = iterator.Error
	return
}

// ReadAddOnRequirement reads a value of the 'add_on_requirement' type from the given iterator.
func ReadAddOnRequirement(iterator *jsoniter.Iterator) *AddOnRequirement {
	object := &AddOnRequirement{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "id":
			value := iterator.ReadString()
			object.id = value
			object.fieldSet_[0] = true
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
			object.fieldSet_[1] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[2] = true
		case "resource":
			value := iterator.ReadString()
			object.resource = value
			object.fieldSet_[3] = true
		case "status":
			value := ReadAddOnRequirementStatus(iterator)
			object.status = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
