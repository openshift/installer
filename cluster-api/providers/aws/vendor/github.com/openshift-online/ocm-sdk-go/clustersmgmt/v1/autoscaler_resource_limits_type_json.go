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

// MarshalAutoscalerResourceLimits writes a value of the 'autoscaler_resource_limits' type to the given writer.
func MarshalAutoscalerResourceLimits(object *AutoscalerResourceLimits, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAutoscalerResourceLimits(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAutoscalerResourceLimits writes a value of the 'autoscaler_resource_limits' type to the given stream.
func WriteAutoscalerResourceLimits(object *AutoscalerResourceLimits, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.gpus != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gpus")
		WriteAutoscalerResourceLimitsGPULimitList(object.gpus, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.cores != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cores")
		WriteResourceRange(object.cores, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_nodes_total")
		stream.WriteInt(object.maxNodesTotal)
		count++
	}
	present_ = object.bitmap_&8 != 0 && object.memory != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("memory")
		WriteResourceRange(object.memory, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAutoscalerResourceLimits reads a value of the 'autoscaler_resource_limits' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAutoscalerResourceLimits(source interface{}) (object *AutoscalerResourceLimits, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAutoscalerResourceLimits(iterator)
	err = iterator.Error
	return
}

// ReadAutoscalerResourceLimits reads a value of the 'autoscaler_resource_limits' type from the given iterator.
func ReadAutoscalerResourceLimits(iterator *jsoniter.Iterator) *AutoscalerResourceLimits {
	object := &AutoscalerResourceLimits{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "gpus":
			value := ReadAutoscalerResourceLimitsGPULimitList(iterator)
			object.gpus = value
			object.bitmap_ |= 1
		case "cores":
			value := ReadResourceRange(iterator)
			object.cores = value
			object.bitmap_ |= 2
		case "max_nodes_total":
			value := iterator.ReadInt()
			object.maxNodesTotal = value
			object.bitmap_ |= 4
		case "memory":
			value := ReadResourceRange(iterator)
			object.memory = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
