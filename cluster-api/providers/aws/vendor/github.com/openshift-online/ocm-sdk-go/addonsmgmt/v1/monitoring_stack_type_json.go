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

// MarshalMonitoringStack writes a value of the 'monitoring_stack' type to the given writer.
func MarshalMonitoringStack(object *MonitoringStack, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeMonitoringStack(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeMonitoringStack writes a value of the 'monitoring_stack' type to the given stream.
func writeMonitoringStack(object *MonitoringStack, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.resources != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("resources")
		writeMonitoringStackResources(object.resources, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMonitoringStack reads a value of the 'monitoring_stack' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMonitoringStack(source interface{}) (object *MonitoringStack, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readMonitoringStack(iterator)
	err = iterator.Error
	return
}

// readMonitoringStack reads a value of the 'monitoring_stack' type from the given iterator.
func readMonitoringStack(iterator *jsoniter.Iterator) *MonitoringStack {
	object := &MonitoringStack{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 1
		case "resources":
			value := readMonitoringStackResources(iterator)
			object.resources = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
