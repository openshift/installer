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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalAddOnParameter writes a value of the 'add_on_parameter' type to the given writer.
func MarshalAddOnParameter(object *AddOnParameter, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnParameter(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnParameter writes a value of the 'add_on_parameter' type to the given stream.
func WriteAddOnParameter(object *AddOnParameter, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(AddOnParameterLinkKind)
	} else {
		stream.WriteString(AddOnParameterKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.addon != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("addon")
		WriteAddOn(object.addon, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.conditions != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("conditions")
		WriteAddOnRequirementList(object.conditions, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("default_value")
		stream.WriteString(object.defaultValue)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("description")
		stream.WriteString(object.description)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("editable")
		stream.WriteBool(object.editable)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("editable_direction")
		stream.WriteString(object.editableDirection)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.options != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("options")
		WriteAddOnParameterOptionList(object.options, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("required")
		stream.WriteBool(object.required)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("validation")
		stream.WriteString(object.validation)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("validation_err_msg")
		stream.WriteString(object.validationErrMsg)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value_type")
		stream.WriteString(object.valueType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnParameter reads a value of the 'add_on_parameter' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnParameter(source interface{}) (object *AddOnParameter, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOnParameter(iterator)
	err = iterator.Error
	return
}

// ReadAddOnParameter reads a value of the 'add_on_parameter' type from the given iterator.
func ReadAddOnParameter(iterator *jsoniter.Iterator) *AddOnParameter {
	object := &AddOnParameter{
		fieldSet_: make([]bool, 16),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == AddOnParameterLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "addon":
			value := ReadAddOn(iterator)
			object.addon = value
			object.fieldSet_[3] = true
		case "conditions":
			value := ReadAddOnRequirementList(iterator)
			object.conditions = value
			object.fieldSet_[4] = true
		case "default_value":
			value := iterator.ReadString()
			object.defaultValue = value
			object.fieldSet_[5] = true
		case "description":
			value := iterator.ReadString()
			object.description = value
			object.fieldSet_[6] = true
		case "editable":
			value := iterator.ReadBool()
			object.editable = value
			object.fieldSet_[7] = true
		case "editable_direction":
			value := iterator.ReadString()
			object.editableDirection = value
			object.fieldSet_[8] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[9] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[10] = true
		case "options":
			value := ReadAddOnParameterOptionList(iterator)
			object.options = value
			object.fieldSet_[11] = true
		case "required":
			value := iterator.ReadBool()
			object.required = value
			object.fieldSet_[12] = true
		case "validation":
			value := iterator.ReadString()
			object.validation = value
			object.fieldSet_[13] = true
		case "validation_err_msg":
			value := iterator.ReadString()
			object.validationErrMsg = value
			object.fieldSet_[14] = true
		case "value_type":
			value := iterator.ReadString()
			object.valueType = value
			object.fieldSet_[15] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
