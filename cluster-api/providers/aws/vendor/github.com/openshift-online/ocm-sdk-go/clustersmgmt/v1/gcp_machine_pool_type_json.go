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

// MarshalGCPMachinePool writes a value of the 'GCP_machine_pool' type to the given writer.
func MarshalGCPMachinePool(object *GCPMachinePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCPMachinePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCPMachinePool writes a value of the 'GCP_machine_pool' type to the given stream.
func WriteGCPMachinePool(object *GCPMachinePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("secure_boot")
		stream.WriteBool(object.secureBoot)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCPMachinePool reads a value of the 'GCP_machine_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCPMachinePool(source interface{}) (object *GCPMachinePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadGCPMachinePool(iterator)
	err = iterator.Error
	return
}

// ReadGCPMachinePool reads a value of the 'GCP_machine_pool' type from the given iterator.
func ReadGCPMachinePool(iterator *jsoniter.Iterator) *GCPMachinePool {
	object := &GCPMachinePool{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "secure_boot":
			value := iterator.ReadBool()
			object.secureBoot = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
