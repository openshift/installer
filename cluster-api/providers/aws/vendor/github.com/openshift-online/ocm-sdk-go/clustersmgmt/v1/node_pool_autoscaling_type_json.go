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

// MarshalNodePoolAutoscaling writes a value of the 'node_pool_autoscaling' type to the given writer.
func MarshalNodePoolAutoscaling(object *NodePoolAutoscaling, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNodePoolAutoscaling(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNodePoolAutoscaling writes a value of the 'node_pool_autoscaling' type to the given stream.
func WriteNodePoolAutoscaling(object *NodePoolAutoscaling, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(NodePoolAutoscalingLinkKind)
	} else {
		stream.WriteString(NodePoolAutoscalingKind)
	}
	count++
	if object.bitmap_&2 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if object.bitmap_&4 != 0 {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("max_replica")
		stream.WriteInt(object.maxReplica)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("min_replica")
		stream.WriteInt(object.minReplica)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNodePoolAutoscaling reads a value of the 'node_pool_autoscaling' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNodePoolAutoscaling(source interface{}) (object *NodePoolAutoscaling, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadNodePoolAutoscaling(iterator)
	err = iterator.Error
	return
}

// ReadNodePoolAutoscaling reads a value of the 'node_pool_autoscaling' type from the given iterator.
func ReadNodePoolAutoscaling(iterator *jsoniter.Iterator) *NodePoolAutoscaling {
	object := &NodePoolAutoscaling{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == NodePoolAutoscalingLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "max_replica":
			value := iterator.ReadInt()
			object.maxReplica = value
			object.bitmap_ |= 8
		case "min_replica":
			value := iterator.ReadInt()
			object.minReplica = value
			object.bitmap_ |= 16
		default:
			iterator.ReadAny()
		}
	}
	return object
}
