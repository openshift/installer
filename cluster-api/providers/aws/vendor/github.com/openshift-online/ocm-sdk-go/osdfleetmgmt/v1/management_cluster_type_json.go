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

package v1 // github.com/openshift-online/ocm-sdk-go/osdfleetmgmt/v1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalManagementCluster writes a value of the 'management_cluster' type to the given writer.
func MarshalManagementCluster(object *ManagementCluster, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeManagementCluster(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeManagementCluster writes a value of the 'management_cluster' type to the given stream.
func writeManagementCluster(object *ManagementCluster, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ManagementClusterLinkKind)
	} else {
		stream.WriteString(ManagementClusterKind)
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
	present_ = object.bitmap_&8 != 0 && object.dns != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns")
		writeDNS(object.dns, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.clusterManagementReference != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_management_reference")
		writeClusterManagementReference(object.clusterManagementReference, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		writeLabelList(object.labels, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.parent != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("parent")
		writeManagementClusterParent(object.parent, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sector")
		stream.WriteString(object.sector)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("update_timestamp")
		stream.WriteString((object.updateTimestamp).Format(time.RFC3339))
	}
	stream.WriteObjectEnd()
}

// UnmarshalManagementCluster reads a value of the 'management_cluster' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalManagementCluster(source interface{}) (object *ManagementCluster, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readManagementCluster(iterator)
	err = iterator.Error
	return
}

// readManagementCluster reads a value of the 'management_cluster' type from the given iterator.
func readManagementCluster(iterator *jsoniter.Iterator) *ManagementCluster {
	object := &ManagementCluster{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ManagementClusterLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "dns":
			value := readDNS(iterator)
			object.dns = value
			object.bitmap_ |= 8
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.bitmap_ |= 16
		case "cluster_management_reference":
			value := readClusterManagementReference(iterator)
			object.clusterManagementReference = value
			object.bitmap_ |= 32
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 64
		case "labels":
			value := readLabelList(iterator)
			object.labels = value
			object.bitmap_ |= 128
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 256
		case "parent":
			value := readManagementClusterParent(iterator)
			object.parent = value
			object.bitmap_ |= 512
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.bitmap_ |= 1024
		case "sector":
			value := iterator.ReadString()
			object.sector = value
			object.bitmap_ |= 2048
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 4096
		case "update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.updateTimestamp = value
			object.bitmap_ |= 8192
		default:
			iterator.ReadAny()
		}
	}
	return object
}
