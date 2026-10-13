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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalClusterAutoNode writes a value of the 'cluster_auto_node' type to the given writer.
func MarshalClusterAutoNode(object *ClusterAutoNode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterAutoNode(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterAutoNode writes a value of the 'cluster_auto_node' type to the given stream.
func WriteClusterAutoNode(object *ClusterAutoNode, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("mode")
		stream.WriteString(object.mode)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteClusterAutoNodeStatus(object.status, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterAutoNode reads a value of the 'cluster_auto_node' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterAutoNode(source interface{}) (object *ClusterAutoNode, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterAutoNode(iterator)
	err = iterator.Error
	return
}

// ReadClusterAutoNode reads a value of the 'cluster_auto_node' type from the given iterator.
func ReadClusterAutoNode(iterator *jsoniter.Iterator) *ClusterAutoNode {
	object := &ClusterAutoNode{
		fieldSet_: make([]bool, 2),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "mode":
			value := iterator.ReadString()
			object.mode = value
			object.fieldSet_[0] = true
		case "status":
			value := ReadClusterAutoNodeStatus(iterator)
			object.status = value
			object.fieldSet_[1] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
