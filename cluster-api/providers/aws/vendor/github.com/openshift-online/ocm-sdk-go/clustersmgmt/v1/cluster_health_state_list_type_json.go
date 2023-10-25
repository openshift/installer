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

// MarshalClusterHealthStateList writes a list of values of the 'cluster_health_state' type to
// the given writer.
func MarshalClusterHealthStateList(list []ClusterHealthState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeClusterHealthStateList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeClusterHealthStateList writes a list of value of the 'cluster_health_state' type to
// the given stream.
func writeClusterHealthStateList(list []ClusterHealthState, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		stream.WriteString(string(value))
	}
	stream.WriteArrayEnd()
}

// UnmarshalClusterHealthStateList reads a list of values of the 'cluster_health_state' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalClusterHealthStateList(source interface{}) (items []ClusterHealthState, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = readClusterHealthStateList(iterator)
	err = iterator.Error
	return
}

// readClusterHealthStateList reads list of values of the ”cluster_health_state' type from
// the given iterator.
func readClusterHealthStateList(iterator *jsoniter.Iterator) []ClusterHealthState {
	list := []ClusterHealthState{}
	for iterator.ReadArray() {
		text := iterator.ReadString()
		item := ClusterHealthState(text)
		list = append(list, item)
	}
	return list
}
