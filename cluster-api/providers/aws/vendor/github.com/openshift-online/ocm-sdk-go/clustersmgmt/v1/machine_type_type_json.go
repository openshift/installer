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

// MarshalMachineType writes a value of the 'machine_type' type to the given writer.
func MarshalMachineType(object *MachineType, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeMachineType(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeMachineType writes a value of the 'machine_type' type to the given stream.
func writeMachineType(object *MachineType, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(MachineTypeLinkKind)
	} else {
		stream.WriteString(MachineTypeKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs_only")
		stream.WriteBool(object.ccsOnly)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.cpu != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cpu")
		writeValue(object.cpu, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("architecture")
		stream.WriteString(string(object.architecture))
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("category")
		stream.WriteString(string(object.category))
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		writeCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("generic_name")
		stream.WriteString(object.genericName)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.memory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		writeValue(object.memory, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&2048 != 0
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
	object = readMachineType(iterator)
	err = iterator.Error
	return
}

// readMachineType reads a value of the 'machine_type' type from the given iterator.
func readMachineType(iterator *jsoniter.Iterator) *MachineType {
	object := &MachineType{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == MachineTypeLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "ccs_only":
			value := iterator.ReadBool()
			object.ccsOnly = value
			object.bitmap_ |= 8
		case "cpu":
			value := readValue(iterator)
			object.cpu = value
			object.bitmap_ |= 16
		case "architecture":
			text := iterator.ReadString()
			value := ProcessorType(text)
			object.architecture = value
			object.bitmap_ |= 32
		case "category":
			text := iterator.ReadString()
			value := MachineTypeCategory(text)
			object.category = value
			object.bitmap_ |= 64
		case "cloud_provider":
			value := readCloudProvider(iterator)
			object.cloudProvider = value
			object.bitmap_ |= 128
		case "generic_name":
			value := iterator.ReadString()
			object.genericName = value
			object.bitmap_ |= 256
		case "memory":
			value := readValue(iterator)
			object.memory = value
			object.bitmap_ |= 512
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 1024
		case "size":
			text := iterator.ReadString()
			value := MachineTypeSize(text)
			object.size = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
