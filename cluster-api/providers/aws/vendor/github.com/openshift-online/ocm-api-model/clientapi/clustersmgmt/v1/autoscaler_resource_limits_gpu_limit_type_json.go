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

// MarshalAutoscalerResourceLimitsGPULimit writes a value of the 'autoscaler_resource_limits_GPU_limit' type to the given writer.
func MarshalAutoscalerResourceLimitsGPULimit(object *AutoscalerResourceLimitsGPULimit, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAutoscalerResourceLimitsGPULimit(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAutoscalerResourceLimitsGPULimit writes a value of the 'autoscaler_resource_limits_GPU_limit' type to the given stream.
func WriteAutoscalerResourceLimitsGPULimit(object *AutoscalerResourceLimitsGPULimit, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.range_ != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("range")
		WriteResourceRange(object.range_, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(object.type_)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAutoscalerResourceLimitsGPULimit reads a value of the 'autoscaler_resource_limits_GPU_limit' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAutoscalerResourceLimitsGPULimit(source interface{}) (object *AutoscalerResourceLimitsGPULimit, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAutoscalerResourceLimitsGPULimit(iterator)
	err = iterator.Error
	return
}

// ReadAutoscalerResourceLimitsGPULimit reads a value of the 'autoscaler_resource_limits_GPU_limit' type from the given iterator.
func ReadAutoscalerResourceLimitsGPULimit(iterator *jsoniter.Iterator) *AutoscalerResourceLimitsGPULimit {
	object := &AutoscalerResourceLimitsGPULimit{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "range":
			value := ReadResourceRange(iterator)
			object.range_ = value
			object.fieldSet_[0] = true
		case "type":
			value := iterator.ReadString()
			object.type_ = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
