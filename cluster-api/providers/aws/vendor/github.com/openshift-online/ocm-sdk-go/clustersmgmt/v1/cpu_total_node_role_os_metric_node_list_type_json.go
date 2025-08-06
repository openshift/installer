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

// MarshalCPUTotalNodeRoleOSMetricNodeList writes a list of values of the 'CPU_total_node_role_OS_metric_node' type to
// the given writer.
func MarshalCPUTotalNodeRoleOSMetricNodeList(list []*CPUTotalNodeRoleOSMetricNode, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCPUTotalNodeRoleOSMetricNodeList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCPUTotalNodeRoleOSMetricNodeList writes a list of value of the 'CPU_total_node_role_OS_metric_node' type to
// the given stream.
func WriteCPUTotalNodeRoleOSMetricNodeList(list []*CPUTotalNodeRoleOSMetricNode, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteCPUTotalNodeRoleOSMetricNode(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalCPUTotalNodeRoleOSMetricNodeList reads a list of values of the 'CPU_total_node_role_OS_metric_node' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalCPUTotalNodeRoleOSMetricNodeList(source interface{}) (items []*CPUTotalNodeRoleOSMetricNode, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadCPUTotalNodeRoleOSMetricNodeList(iterator)
	err = iterator.Error
	return
}

// ReadCPUTotalNodeRoleOSMetricNodeList reads list of values of the ”CPU_total_node_role_OS_metric_node' type from
// the given iterator.
func ReadCPUTotalNodeRoleOSMetricNodeList(iterator *jsoniter.Iterator) []*CPUTotalNodeRoleOSMetricNode {
	list := []*CPUTotalNodeRoleOSMetricNode{}
	for iterator.ReadArray() {
		item := ReadCPUTotalNodeRoleOSMetricNode(iterator)
		list = append(list, item)
	}
	return list
}
