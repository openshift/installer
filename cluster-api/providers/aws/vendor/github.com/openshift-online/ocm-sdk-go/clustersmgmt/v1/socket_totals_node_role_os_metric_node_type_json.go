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

// MarshalSocketTotalsNodeRoleOSMetricNode writes a value of the 'socket_totals_node_role_OS_metric_node' type to the given writer.
func MarshalSocketTotalsNodeRoleOSMetricNode(object *SocketTotalsNodeRoleOSMetricNode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeSocketTotalsNodeRoleOSMetricNode(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeSocketTotalsNodeRoleOSMetricNode writes a value of the 'socket_totals_node_role_OS_metric_node' type to the given stream.
func writeSocketTotalsNodeRoleOSMetricNode(object *SocketTotalsNodeRoleOSMetricNode, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.socketTotals != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("socket_totals")
		writeSocketTotalNodeRoleOSMetricNodeList(object.socketTotals, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalSocketTotalsNodeRoleOSMetricNode reads a value of the 'socket_totals_node_role_OS_metric_node' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalSocketTotalsNodeRoleOSMetricNode(source interface{}) (object *SocketTotalsNodeRoleOSMetricNode, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readSocketTotalsNodeRoleOSMetricNode(iterator)
	err = iterator.Error
	return
}

// readSocketTotalsNodeRoleOSMetricNode reads a value of the 'socket_totals_node_role_OS_metric_node' type from the given iterator.
func readSocketTotalsNodeRoleOSMetricNode(iterator *jsoniter.Iterator) *SocketTotalsNodeRoleOSMetricNode {
	object := &SocketTotalsNodeRoleOSMetricNode{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "socket_totals":
			value := readSocketTotalNodeRoleOSMetricNodeList(iterator)
			object.socketTotals = value
			object.bitmap_ |= 1
		default:
			iterator.ReadAny()
		}
	}
	return object
}
