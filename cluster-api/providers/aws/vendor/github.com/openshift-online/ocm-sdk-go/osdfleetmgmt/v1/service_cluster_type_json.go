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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalServiceCluster writes a value of the 'service_cluster' type to the given writer.
func MarshalServiceCluster(object *ServiceCluster, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeServiceCluster(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeServiceCluster writes a value of the 'service_cluster' type to the given stream.
func writeServiceCluster(object *ServiceCluster, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ServiceClusterLinkKind)
	} else {
		stream.WriteString(ServiceClusterKind)
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
	present_ = object.bitmap_&64 != 0 && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		writeLabelList(object.labels, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.provisionShardReference != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_shard_reference")
		writeProvisionShardReference(object.provisionShardReference, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sector")
		stream.WriteString(object.sector)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
	}
	stream.WriteObjectEnd()
}

// UnmarshalServiceCluster reads a value of the 'service_cluster' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalServiceCluster(source interface{}) (object *ServiceCluster, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readServiceCluster(iterator)
	err = iterator.Error
	return
}

// readServiceCluster reads a value of the 'service_cluster' type from the given iterator.
func readServiceCluster(iterator *jsoniter.Iterator) *ServiceCluster {
	object := &ServiceCluster{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ServiceClusterLinkKind {
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
		case "labels":
			value := readLabelList(iterator)
			object.labels = value
			object.bitmap_ |= 64
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 128
		case "provision_shard_reference":
			value := readProvisionShardReference(iterator)
			object.provisionShardReference = value
			object.bitmap_ |= 256
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.bitmap_ |= 512
		case "sector":
			value := iterator.ReadString()
			object.sector = value
			object.bitmap_ |= 1024
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 2048
		default:
			iterator.ReadAny()
		}
	}
	return object
}
