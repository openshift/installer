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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalServiceClusterRequestPayload writes a value of the 'service_cluster_request_payload' type to the given writer.
func MarshalServiceClusterRequestPayload(object *ServiceClusterRequestPayload, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeServiceClusterRequestPayload(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeServiceClusterRequestPayload writes a value of the 'service_cluster_request_payload' type to the given stream.
func writeServiceClusterRequestPayload(object *ServiceClusterRequestPayload, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		writeLabelRequestPayloadList(object.labels, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
	}
	stream.WriteObjectEnd()
}

// UnmarshalServiceClusterRequestPayload reads a value of the 'service_cluster_request_payload' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalServiceClusterRequestPayload(source interface{}) (object *ServiceClusterRequestPayload, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readServiceClusterRequestPayload(iterator)
	err = iterator.Error
	return
}

// readServiceClusterRequestPayload reads a value of the 'service_cluster_request_payload' type from the given iterator.
func readServiceClusterRequestPayload(iterator *jsoniter.Iterator) *ServiceClusterRequestPayload {
	object := &ServiceClusterRequestPayload{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.bitmap_ |= 1
		case "labels":
			value := readLabelRequestPayloadList(iterator)
			object.labels = value
			object.bitmap_ |= 2
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
