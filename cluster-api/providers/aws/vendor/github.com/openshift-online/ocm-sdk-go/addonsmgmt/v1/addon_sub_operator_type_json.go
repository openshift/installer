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

// MarshalAddonSubOperator writes a value of the 'addon_sub_operator' type to the given writer.
func MarshalAddonSubOperator(object *AddonSubOperator, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddonSubOperator(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddonSubOperator writes a value of the 'addon_sub_operator' type to the given stream.
func WriteAddonSubOperator(object *AddonSubOperator, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.addon != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon")
		WriteAddon(object.addon, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_name")
		stream.WriteString(object.operatorName)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operator_namespace")
		stream.WriteString(object.operatorNamespace)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddonSubOperator reads a value of the 'addon_sub_operator' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddonSubOperator(source interface{}) (object *AddonSubOperator, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddonSubOperator(iterator)
	err = iterator.Error
	return
}

// ReadAddonSubOperator reads a value of the 'addon_sub_operator' type from the given iterator.
func ReadAddonSubOperator(iterator *jsoniter.Iterator) *AddonSubOperator {
	object := &AddonSubOperator{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "addon":
			value := ReadAddon(iterator)
			object.addon = value
			object.bitmap_ |= 1
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 2
		case "operator_name":
			value := iterator.ReadString()
			object.operatorName = value
			object.bitmap_ |= 4
		case "operator_namespace":
			value := iterator.ReadString()
			object.operatorNamespace = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
