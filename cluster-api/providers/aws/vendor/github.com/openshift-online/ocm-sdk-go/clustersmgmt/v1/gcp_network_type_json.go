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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalGCPNetwork writes a value of the 'GCP_network' type to the given writer.
func MarshalGCPNetwork(object *GCPNetwork, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeGCPNetwork(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeGCPNetwork writes a value of the 'GCP_network' type to the given stream.
func writeGCPNetwork(object *GCPNetwork, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_name")
		stream.WriteString(object.vpcName)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_project_id")
		stream.WriteString(object.vpcProjectID)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_subnet")
		stream.WriteString(object.computeSubnet)
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("control_plane_subnet")
		stream.WriteString(object.controlPlaneSubnet)
	}
	stream.WriteObjectEnd()
}

// UnmarshalGCPNetwork reads a value of the 'GCP_network' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalGCPNetwork(source interface{}) (object *GCPNetwork, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readGCPNetwork(iterator)
	err = iterator.Error
	return
}

// readGCPNetwork reads a value of the 'GCP_network' type from the given iterator.
func readGCPNetwork(iterator *jsoniter.Iterator) *GCPNetwork {
	object := &GCPNetwork{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "vpc_name":
			value := iterator.ReadString()
			object.vpcName = value
			object.bitmap_ |= 1
		case "vpc_project_id":
			value := iterator.ReadString()
			object.vpcProjectID = value
			object.bitmap_ |= 2
		case "compute_subnet":
			value := iterator.ReadString()
			object.computeSubnet = value
			object.bitmap_ |= 4
		case "control_plane_subnet":
			value := iterator.ReadString()
			object.controlPlaneSubnet = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
