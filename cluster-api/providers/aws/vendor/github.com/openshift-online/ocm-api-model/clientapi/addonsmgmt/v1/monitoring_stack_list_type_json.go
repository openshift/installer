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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/addonsmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalMonitoringStackList writes a list of values of the 'monitoring_stack' type to
// the given writer.
func MarshalMonitoringStackList(list []*MonitoringStack, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMonitoringStackList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMonitoringStackList writes a list of value of the 'monitoring_stack' type to
// the given stream.
func WriteMonitoringStackList(list []*MonitoringStack, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteMonitoringStack(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalMonitoringStackList reads a list of values of the 'monitoring_stack' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalMonitoringStackList(source interface{}) (items []*MonitoringStack, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadMonitoringStackList(iterator)
	err = iterator.Error
	return
}

// ReadMonitoringStackList reads list of values of the ”monitoring_stack' type from
// the given iterator.
func ReadMonitoringStackList(iterator *jsoniter.Iterator) []*MonitoringStack {
	list := []*MonitoringStack{}
	for iterator.ReadArray() {
		item := ReadMonitoringStack(iterator)
		list = append(list, item)
	}
	return list
}
