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

// MarshalSocketTotalNodeRoleOSMetricNode writes a value of the 'socket_total_node_role_OS_metric_node' type to the given writer.
func MarshalSocketTotalNodeRoleOSMetricNode(object *SocketTotalNodeRoleOSMetricNode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSocketTotalNodeRoleOSMetricNode(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSocketTotalNodeRoleOSMetricNode writes a value of the 'socket_total_node_role_OS_metric_node' type to the given stream.
func writeSocketTotalNodeRoleOSMetricNode(object *SocketTotalNodeRoleOSMetricNode, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.nodeRoles != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_roles")
		writeStringList(object.nodeRoles, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("operating_system")
		stream.WriteString(object.operatingSystem)
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("socket_total")
		stream.WriteFloat64(object.socketTotal)
		count++
	}
	present_ = object.bitmap_&8 != 0
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
	object = readSocketTotalNodeRoleOSMetricNode(iterator)
	err = iterator.Error
	return
}

// readSocketTotalNodeRoleOSMetricNode reads a value of the 'socket_total_node_role_OS_metric_node' type from the given iterator.
func readSocketTotalNodeRoleOSMetricNode(iterator *jsoniter.Iterator) *SocketTotalNodeRoleOSMetricNode {
	object := &SocketTotalNodeRoleOSMetricNode{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "node_roles":
			value := readStringList(iterator)
			object.nodeRoles = value
			object.bitmap_ |= 1
		case "operating_system":
			value := iterator.ReadString()
			object.operatingSystem = value
			object.bitmap_ |= 2
		case "socket_total":
			value := iterator.ReadFloat64()
			object.socketTotal = value
			object.bitmap_ |= 4
		case "time":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.time = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
