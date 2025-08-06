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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalSocketTotalNodeRoleOSMetricNode writes a value of the 'socket_total_node_role_OS_metric_node' type to the given writer.
func MarshalSocketTotalNodeRoleOSMetricNode(object *SocketTotalNodeRoleOSMetricNode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteSocketTotalNodeRoleOSMetricNode(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteSocketTotalNodeRoleOSMetricNode writes a value of the 'socket_total_node_role_OS_metric_node' type to the given stream.
func WriteSocketTotalNodeRoleOSMetricNode(object *SocketTotalNodeRoleOSMetricNode, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.nodeRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_roles")
		WriteStringList(object.nodeRoles, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operating_system")
		stream.WriteString(object.operatingSystem)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("socket_total")
		stream.WriteFloat64(object.socketTotal)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("time")
		stream.WriteString((object.time).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalSocketTotalNodeRoleOSMetricNode reads a value of the 'socket_total_node_role_OS_metric_node' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSocketTotalNodeRoleOSMetricNode(source interface{}) (object *SocketTotalNodeRoleOSMetricNode, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadSocketTotalNodeRoleOSMetricNode(iterator)
	err = iterator.Error
	return
}

// ReadSocketTotalNodeRoleOSMetricNode reads a value of the 'socket_total_node_role_OS_metric_node' type from the given iterator.
func ReadSocketTotalNodeRoleOSMetricNode(iterator *jsoniter.Iterator) *SocketTotalNodeRoleOSMetricNode {
	object := &SocketTotalNodeRoleOSMetricNode{
		fieldSet_: make([]bool, 4),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "node_roles":
			value := ReadStringList(iterator)
			object.nodeRoles = value
			object.fieldSet_[0] = true
		case "operating_system":
			value := iterator.ReadString()
			object.operatingSystem = value
			object.fieldSet_[1] = true
		case "socket_total":
			value := iterator.ReadFloat64()
			object.socketTotal = value
			object.fieldSet_[2] = true
		case "time":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.time = value
			object.fieldSet_[3] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
