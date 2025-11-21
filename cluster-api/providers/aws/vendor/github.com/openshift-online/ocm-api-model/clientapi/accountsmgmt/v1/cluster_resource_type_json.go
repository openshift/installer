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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/accountsmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalClusterResource writes a value of the 'cluster_resource' type to the given writer.
func MarshalClusterResource(object *ClusterResource, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterResource(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterResource writes a value of the 'cluster_resource' type to the given stream.
func WriteClusterResource(object *ClusterResource, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.total != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("total")
		WriteValueUnit(object.total, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_timestamp")
		stream.WriteString((object.updatedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.used != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("used")
		WriteValueUnit(object.used, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterResource reads a value of the 'cluster_resource' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterResource(source interface{}) (object *ClusterResource, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterResource(iterator)
	err = iterator.Error
	return
}

// ReadClusterResource reads a value of the 'cluster_resource' type from the given iterator.
func ReadClusterResource(iterator *jsoniter.Iterator) *ClusterResource {
	object := &ClusterResource{
		fieldSet_: make([]bool, 3),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "total":
			value := ReadValueUnit(iterator)
			object.total = value
			object.fieldSet_[0] = true
		case "updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedTimestamp = value
			object.fieldSet_[1] = true
		case "used":
			value := ReadValueUnit(iterator)
			object.used = value
			object.fieldSet_[2] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
