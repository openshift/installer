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

// MarshalMachinePoolAutoscaling writes a value of the 'machine_pool_autoscaling' type to the given writer.
func MarshalMachinePoolAutoscaling(object *MachinePoolAutoscaling, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMachinePoolAutoscaling(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMachinePoolAutoscaling writes a value of the 'machine_pool_autoscaling' type to the given stream.
func WriteMachinePoolAutoscaling(object *MachinePoolAutoscaling, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(MachinePoolAutoscalingLinkKind)
	} else {
		stream.WriteString(MachinePoolAutoscalingKind)
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
		stream.WriteObjectField("max_replicas")
		stream.WriteInt(object.maxReplicas)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("min_replicas")
		stream.WriteInt(object.minReplicas)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMachinePoolAutoscaling reads a value of the 'machine_pool_autoscaling' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMachinePoolAutoscaling(source interface{}) (object *MachinePoolAutoscaling, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadMachinePoolAutoscaling(iterator)
	err = iterator.Error
	return
}

// ReadMachinePoolAutoscaling reads a value of the 'machine_pool_autoscaling' type from the given iterator.
func ReadMachinePoolAutoscaling(iterator *jsoniter.Iterator) *MachinePoolAutoscaling {
	object := &MachinePoolAutoscaling{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == MachinePoolAutoscalingLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "max_replicas":
			value := iterator.ReadInt()
			object.maxReplicas = value
			object.fieldSet_[3] = true
		case "min_replicas":
			value := iterator.ReadInt()
			object.minReplicas = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
