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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalClusterOperatorInfo writes a value of the 'cluster_operator_info' type to the given writer.
func MarshalClusterOperatorInfo(object *ClusterOperatorInfo, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterOperatorInfo(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterOperatorInfo writes a value of the 'cluster_operator_info' type to the given stream.
func writeClusterOperatorInfo(object *ClusterOperatorInfo, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("condition")
		stream.WriteString(string(object.condition))
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("reason")
		stream.WriteString(object.reason)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("time")
		stream.WriteString((object.time).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterOperatorInfo reads a value of the 'cluster_operator_info' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterOperatorInfo(source interface{}) (object *ClusterOperatorInfo, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readClusterOperatorInfo(iterator)
	err = iterator.Error
	return
}

// readClusterOperatorInfo reads a value of the 'cluster_operator_info' type from the given iterator.
func readClusterOperatorInfo(iterator *jsoniter.Iterator) *ClusterOperatorInfo {
	object := &ClusterOperatorInfo{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "condition":
			text := iterator.ReadString()
			value := ClusterOperatorState(text)
			object.condition = value
			object.bitmap_ |= 1
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 2
		case "reason":
			value := iterator.ReadString()
			object.reason = value
			object.bitmap_ |= 4
		case "time":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.time = value
			object.bitmap_ |= 8
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
