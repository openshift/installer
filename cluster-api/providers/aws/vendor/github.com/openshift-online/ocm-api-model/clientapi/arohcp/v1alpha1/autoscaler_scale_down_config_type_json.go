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

// MarshalAutoscalerScaleDownConfig writes a value of the 'autoscaler_scale_down_config' type to the given writer.
func MarshalAutoscalerScaleDownConfig(object *AutoscalerScaleDownConfig, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteAutoscalerScaleDownConfig(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteAutoscalerScaleDownConfig writes a value of the 'autoscaler_scale_down_config' type to the given stream.
func WriteAutoscalerScaleDownConfig(object *AutoscalerScaleDownConfig, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("delay_after_add")
		stream.WriteString(object.delayAfterAdd)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("delay_after_delete")
		stream.WriteString(object.delayAfterDelete)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("delay_after_failure")
		stream.WriteString(object.delayAfterFailure)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("unneeded_time")
		stream.WriteString(object.unneededTime)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("utilization_threshold")
		stream.WriteString(object.utilizationThreshold)
	}
	stream.WriteObjectEnd()
}

// UnmarshalAutoscalerScaleDownConfig reads a value of the 'autoscaler_scale_down_config' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalAutoscalerScaleDownConfig(source interface{}) (object *AutoscalerScaleDownConfig, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadAutoscalerScaleDownConfig(iterator)
	err = iterator.Error
	return
}

// ReadAutoscalerScaleDownConfig reads a value of the 'autoscaler_scale_down_config' type from the given iterator.
func ReadAutoscalerScaleDownConfig(iterator *jsoniter.Iterator) *AutoscalerScaleDownConfig {
	object := &AutoscalerScaleDownConfig{
		fieldSet_: make([]bool, 6),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "delay_after_add":
			value := iterator.ReadString()
			object.delayAfterAdd = value
			object.fieldSet_[0] = true
		case "delay_after_delete":
			value := iterator.ReadString()
			object.delayAfterDelete = value
			object.fieldSet_[1] = true
		case "delay_after_failure":
			value := iterator.ReadString()
			object.delayAfterFailure = value
			object.fieldSet_[2] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[3] = true
		case "unneeded_time":
			value := iterator.ReadString()
			object.unneededTime = value
			object.fieldSet_[4] = true
		case "utilization_threshold":
			value := iterator.ReadString()
			object.utilizationThreshold = value
			object.fieldSet_[5] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
