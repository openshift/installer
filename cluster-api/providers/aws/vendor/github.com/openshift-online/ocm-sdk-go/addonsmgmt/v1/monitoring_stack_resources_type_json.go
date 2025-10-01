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

// MarshalMonitoringStackResources writes a value of the 'monitoring_stack_resources' type to the given writer.
func MarshalMonitoringStackResources(object *MonitoringStackResources, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMonitoringStackResources(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMonitoringStackResources writes a value of the 'monitoring_stack_resources' type to the given stream.
func WriteMonitoringStackResources(object *MonitoringStackResources, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.limits != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("limits")
		WriteMonitoringStackResource(object.limits, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.requests != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("requests")
		WriteMonitoringStackResource(object.requests, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMonitoringStackResources reads a value of the 'monitoring_stack_resources' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMonitoringStackResources(source interface{}) (object *MonitoringStackResources, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadMonitoringStackResources(iterator)
	err = iterator.Error
	return
}

// ReadMonitoringStackResources reads a value of the 'monitoring_stack_resources' type from the given iterator.
func ReadMonitoringStackResources(iterator *jsoniter.Iterator) *MonitoringStackResources {
	object := &MonitoringStackResources{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "limits":
			value := ReadMonitoringStackResource(iterator)
			object.limits = value
			object.bitmap_ |= 1
		case "requests":
			value := ReadMonitoringStackResource(iterator)
			object.requests = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
