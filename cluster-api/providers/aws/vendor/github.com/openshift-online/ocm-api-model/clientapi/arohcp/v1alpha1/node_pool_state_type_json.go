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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalNodePoolState writes a value of the 'node_pool_state' type to the given writer.
func MarshalNodePoolState(object *NodePoolState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNodePoolState(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNodePoolState writes a value of the 'node_pool_state' type to the given stream.
func WriteNodePoolState(object *NodePoolState, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(NodePoolStateLinkKind)
	} else {
		stream.WriteString(NodePoolStateKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_updated_timestamp")
		stream.WriteString((object.lastUpdatedTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("value")
		stream.WriteString(object.nodePoolStateValue)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNodePoolState reads a value of the 'node_pool_state' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNodePoolState(source interface{}) (object *NodePoolState, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadNodePoolState(iterator)
	err = iterator.Error
	return
}

// ReadNodePoolState reads a value of the 'node_pool_state' type from the given iterator.
func ReadNodePoolState(iterator *jsoniter.Iterator) *NodePoolState {
	object := &NodePoolState{
		fieldSet_: make([]bool, 5),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == NodePoolStateLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "last_updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdatedTimestamp = value
			object.fieldSet_[3] = true
		case "value":
			value := iterator.ReadString()
			object.nodePoolStateValue = value
			object.fieldSet_[4] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
