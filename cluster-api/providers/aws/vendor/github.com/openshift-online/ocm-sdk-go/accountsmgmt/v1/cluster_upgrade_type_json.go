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

package v1 // github.com/openshift-online/ocm-sdk-go/accountsmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalClusterUpgrade writes a value of the 'cluster_upgrade' type to the given writer.
func MarshalClusterUpgrade(object *ClusterUpgrade, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterUpgrade(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterUpgrade writes a value of the 'cluster_upgrade' type to the given stream.
func WriteClusterUpgrade(object *ClusterUpgrade, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("available")
		stream.WriteBool(object.available)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		stream.WriteString(object.state)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_timestamp")
		stream.WriteString((object.updatedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		stream.WriteString(object.version)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterUpgrade reads a value of the 'cluster_upgrade' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterUpgrade(source interface{}) (object *ClusterUpgrade, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterUpgrade(iterator)
	err = iterator.Error
	return
}

// ReadClusterUpgrade reads a value of the 'cluster_upgrade' type from the given iterator.
func ReadClusterUpgrade(iterator *jsoniter.Iterator) *ClusterUpgrade {
	object := &ClusterUpgrade{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "available":
			value := iterator.ReadBool()
			object.available = value
			object.bitmap_ |= 1
		case "state":
			value := iterator.ReadString()
			object.state = value
			object.bitmap_ |= 2
		case "updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedTimestamp = value
			object.bitmap_ |= 4
		case "version":
			value := iterator.ReadString()
			object.version = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
