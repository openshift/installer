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

// MarshalClusterHealthStateList writes a list of values of the 'cluster_health_state' type to
// the given writer.
func MarshalClusterHealthStateList(list []ClusterHealthState, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterHealthStateList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterHealthStateList writes a list of value of the 'cluster_health_state' type to
// the given stream.
func WriteClusterHealthStateList(list []ClusterHealthState, stream *jsoniter.Stream) {
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
	items = ReadClusterHealthStateList(iterator)
	err = iterator.Error
	return
}

// ReadClusterHealthStateList reads list of values of the ”cluster_health_state' type from
// the given iterator.
func ReadClusterHealthStateList(iterator *jsoniter.Iterator) []ClusterHealthState {
	list := []ClusterHealthState{}
	for iterator.ReadArray() {
		text := iterator.ReadString()
		item := ClusterHealthState(text)
		list = append(list, item)
	}
	return list
}
