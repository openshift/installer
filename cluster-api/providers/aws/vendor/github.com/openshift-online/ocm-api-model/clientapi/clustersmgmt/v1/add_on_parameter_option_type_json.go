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

// MarshalAddOnParameterOption writes a value of the 'add_on_parameter_option' type to the given writer.
func MarshalAddOnParameterOption(object *AddOnParameterOption, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAddOnParameterOption(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAddOnParameterOption writes a value of the 'add_on_parameter_option' type to the given stream.
func WriteAddOnParameterOption(object *AddOnParameterOption, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rank")
		stream.WriteInt(object.rank)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.requirements != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requirements")
		WriteAddOnRequirementList(object.requirements, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.value)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAddOnParameterOption reads a value of the 'add_on_parameter_option' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAddOnParameterOption(source interface{}) (object *AddOnParameterOption, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAddOnParameterOption(iterator)
	err = iterator.Error
	return
}

// ReadAddOnParameterOption reads a value of the 'add_on_parameter_option' type from the given iterator.
func ReadAddOnParameterOption(iterator *jsoniter.Iterator) *AddOnParameterOption {
	object := &AddOnParameterOption{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[0] = true
		case "rank":
			value := iterator.ReadInt()
			object.rank = value
			object.fieldSet_[1] = true
		case "requirements":
			value := ReadAddOnRequirementList(iterator)
			object.requirements = value
			object.fieldSet_[2] = true
		case "value":
			value := iterator.ReadString()
			object.value = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
