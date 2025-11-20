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

// MarshalGCPNetwork writes a value of the 'GCP_network' type to the given writer.
func MarshalGCPNetwork(object *GCPNetwork, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteGCPNetwork(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteGCPNetwork writes a value of the 'GCP_network' type to the given stream.
func WriteGCPNetwork(object *GCPNetwork, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_name")
		stream.WriteString(object.vpcName)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_project_id")
		stream.WriteString(object.vpcProjectID)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_subnet")
		stream.WriteString(object.computeSubnet)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
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
	object = ReadGCPNetwork(iterator)
	err = iterator.Error
	return
}

// ReadGCPNetwork reads a value of the 'GCP_network' type from the given iterator.
func ReadGCPNetwork(iterator *jsoniter.Iterator) *GCPNetwork {
	object := &GCPNetwork{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "vpc_name":
			value := iterator.ReadString()
			object.vpcName = value
			object.fieldSet_[0] = true
		case "vpc_project_id":
			value := iterator.ReadString()
			object.vpcProjectID = value
			object.fieldSet_[1] = true
		case "compute_subnet":
			value := iterator.ReadString()
			object.computeSubnet = value
			object.fieldSet_[2] = true
		case "control_plane_subnet":
			value := iterator.ReadString()
			object.controlPlaneSubnet = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
