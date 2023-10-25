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

// MarshalPendingDeleteCluster writes a value of the 'pending_delete_cluster' type to the given writer.
func MarshalPendingDeleteCluster(object *PendingDeleteCluster, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writePendingDeleteCluster(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writePendingDeleteCluster writes a value of the 'pending_delete_cluster' type to the given stream.
func writePendingDeleteCluster(object *PendingDeleteCluster, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(PendingDeleteClusterLinkKind)
	} else {
		stream.WriteString(PendingDeleteClusterKind)
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
		stream.WriteObjectField("best_effort")
		stream.WriteBool(object.bestEffort)
		count++
	}
	present_ = object.bitmap_&16 != 0 && object.cluster != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster")
		writeCluster(object.cluster, stream)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalPendingDeleteCluster reads a value of the 'pending_delete_cluster' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalPendingDeleteCluster(source interface{}) (object *PendingDeleteCluster, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readPendingDeleteCluster(iterator)
	err = iterator.Error
	return
}

// readPendingDeleteCluster reads a value of the 'pending_delete_cluster' type from the given iterator.
func readPendingDeleteCluster(iterator *jsoniter.Iterator) *PendingDeleteCluster {
	object := &PendingDeleteCluster{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == PendingDeleteClusterLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "best_effort":
			value := iterator.ReadBool()
			object.bestEffort = value
			object.bitmap_ |= 8
		case "cluster":
			value := readCluster(iterator)
			object.cluster = value
			object.bitmap_ |= 16
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
