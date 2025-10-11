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

// MarshalClusterMigration writes a value of the 'cluster_migration' type to the given writer.
func MarshalClusterMigration(object *ClusterMigration, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterMigration(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterMigration writes a value of the 'cluster_migration' type to the given stream.
func WriteClusterMigration(object *ClusterMigration, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ClusterMigrationLinkKind)
	} else {
		stream.WriteString(ClusterMigrationKind)
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
		stream.WriteObjectField("cluster_id")
		stream.WriteString(object.clusterID)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.sdnToOvn != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sdn_to_ovn")
		WriteSdnToOvnClusterMigration(object.sdnToOvn, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.state != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("state")
		WriteClusterMigrationState(object.state, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("type")
		stream.WriteString(string(object.type_))
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("updated_timestamp")
		stream.WriteString((object.updatedTimestamp).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterMigration reads a value of the 'cluster_migration' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterMigration(source interface{}) (object *ClusterMigration, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterMigration(iterator)
	err = iterator.Error
	return
}

// ReadClusterMigration reads a value of the 'cluster_migration' type from the given iterator.
func ReadClusterMigration(iterator *jsoniter.Iterator) *ClusterMigration {
	object := &ClusterMigration{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ClusterMigrationLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "cluster_id":
			value := iterator.ReadString()
			object.clusterID = value
			object.bitmap_ |= 8
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 16
		case "sdn_to_ovn":
			value := ReadSdnToOvnClusterMigration(iterator)
			object.sdnToOvn = value
			object.bitmap_ |= 32
		case "state":
			value := ReadClusterMigrationState(iterator)
			object.state = value
			object.bitmap_ |= 64
		case "type":
			text := iterator.ReadString()
			value := ClusterMigrationType(text)
			object.type_ = value
			object.bitmap_ |= 128
		case "updated_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updatedTimestamp = value
			object.bitmap_ |= 256
		default:
			iterator.ReadAny()
		}
	}
	return object
}
