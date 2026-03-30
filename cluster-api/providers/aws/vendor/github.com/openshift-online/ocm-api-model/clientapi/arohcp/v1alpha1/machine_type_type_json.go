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

package v1alpha1 // github.com/openshift-online/ocm-api-model/clientapi/arohcp/v1alpha1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalMachineType writes a value of the 'machine_type' type to the given writer.
func MarshalMachineType(object *MachineType, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMachineType(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMachineType writes a value of the 'machine_type' type to the given stream.
func WriteMachineType(object *MachineType, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(MachineTypeLinkKind)
	} else {
		stream.WriteString(MachineTypeKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs_only")
		stream.WriteBool(object.ccsOnly)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.cpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu")
		WriteValue(object.cpu, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("architecture")
		stream.WriteString(string(object.architecture))
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("category")
		stream.WriteString(string(object.category))
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("generic_name")
		stream.WriteString(object.genericName)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.memory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		WriteValue(object.memory, stream)
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
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("size")
		stream.WriteString(string(object.size))
	}
	stream.WriteObjectEnd()
}

// UnmarshalMachineType reads a value of the 'machine_type' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMachineType(source interface{}) (object *MachineType, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadMachineType(iterator)
	err = iterator.Error
	return
}

// ReadMachineType reads a value of the 'machine_type' type from the given iterator.
func ReadMachineType(iterator *jsoniter.Iterator) *MachineType {
	object := &MachineType{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == MachineTypeLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "ccs_only":
			value := iterator.ReadBool()
			object.ccsOnly = value
			object.fieldSet_[3] = true
		case "cpu":
			value := ReadValue(iterator)
			object.cpu = value
			object.fieldSet_[4] = true
		case "architecture":
			text := iterator.ReadString()
			value := ProcessorType(text)
			object.architecture = value
			object.fieldSet_[5] = true
		case "category":
			text := iterator.ReadString()
			value := MachineTypeCategory(text)
			object.category = value
			object.fieldSet_[6] = true
		case "cloud_provider":
			value := ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.fieldSet_[7] = true
		case "generic_name":
			value := iterator.ReadString()
			object.genericName = value
			object.fieldSet_[8] = true
		case "memory":
			value := ReadValue(iterator)
			object.memory = value
			object.fieldSet_[9] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[10] = true
		case "size":
			text := iterator.ReadString()
			value := MachineTypeSize(text)
			object.size = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
