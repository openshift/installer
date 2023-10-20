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

package v1 // github.com/openshift-online/ocm-sdk-go/statusboard/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalServiceInfo writes a value of the 'service_info' type to the given writer.
func MarshalServiceInfo(object *ServiceInfo, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeServiceInfo(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeServiceInfo writes a value of the 'service_info' type to the given stream.
func writeServiceInfo(object *ServiceInfo, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("fullname")
		stream.WriteString(object.fullname)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status_type")
		stream.WriteString(object.statusType)
	}
	stream.WriteObjectEnd()
}

// UnmarshalServiceInfo reads a value of the 'service_info' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalServiceInfo(source interface{}) (object *ServiceInfo, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readServiceInfo(iterator)
	err = iterator.Error
	return
}

// readServiceInfo reads a value of the 'service_info' type from the given iterator.
func readServiceInfo(iterator *jsoniter.Iterator) *ServiceInfo {
	object := &ServiceInfo{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "fullname":
			value := iterator.ReadString()
			object.fullname = value
			object.bitmap_ |= 1
		case "status_type":
			value := iterator.ReadString()
			object.statusType = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
