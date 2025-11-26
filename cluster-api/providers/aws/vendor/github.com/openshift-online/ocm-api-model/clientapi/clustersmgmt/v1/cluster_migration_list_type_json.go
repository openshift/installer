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

// MarshalClusterMigrationList writes a list of values of the 'cluster_migration' type to
// the given writer.
func MarshalClusterMigrationList(list []*ClusterMigration, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterMigrationList(list, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterMigrationList writes a list of value of the 'cluster_migration' type to
// the given stream.
func WriteClusterMigrationList(list []*ClusterMigration, stream *jsoniter.Stream) {
	stream.WriteArrayStart()
	for i, value := range list {
		if i > 0 {
			stream.WriteMore()
		}
		WriteClusterMigration(value, stream)
	}
	stream.WriteArrayEnd()
}

// UnmarshalClusterMigrationList reads a list of values of the 'cluster_migration' type
// from the given source, which can be a slice of bytes, a string or a reader.
func UnmarshalClusterMigrationList(source interface{}) (items []*ClusterMigration, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	items = ReadClusterMigrationList(iterator)
	err = iterator.Error
	return
}

// ReadClusterMigrationList reads list of values of the ”cluster_migration' type from
// the given iterator.
func ReadClusterMigrationList(iterator *jsoniter.Iterator) []*ClusterMigration {
	list := []*ClusterMigration{}
	for iterator.ReadArray() {
		item := ReadClusterMigration(iterator)
		list = append(list, item)
	}
	return list
}
