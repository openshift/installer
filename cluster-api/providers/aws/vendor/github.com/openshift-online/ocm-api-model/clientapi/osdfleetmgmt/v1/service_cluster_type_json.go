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

package v1 // github.com/openshift-online/ocm-api-model/clientapi/osdfleetmgmt/v1

import (
	"io"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalServiceCluster writes a value of the 'service_cluster' type to the given writer.
func MarshalServiceCluster(object *ServiceCluster, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteServiceCluster(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteServiceCluster writes a value of the 'service_cluster' type to the given stream.
func WriteServiceCluster(object *ServiceCluster, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ServiceClusterLinkKind)
	} else {
		stream.WriteString(ServiceClusterKind)
	}
	count++
	if len(object.fieldSet_) > 1 && object.fieldSet_[1] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("id")
		stream.WriteString(object.id)
		count++
	}
	if len(object.fieldSet_) > 2 && object.fieldSet_[2] {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("href")
		stream.WriteString(object.href)
		count++
	}
	var present_ bool
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.dns != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("dns")
		WriteDNS(object.dns, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		stream.WriteString(object.cloudProvider)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.clusterManagementReference != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cluster_management_reference")
		WriteClusterManagementReference(object.clusterManagementReference, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		WriteLabelList(object.labels, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.provisionShardReference != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("provision_shard_reference")
		WriteProvisionShardReference(object.provisionShardReference, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		stream.WriteString(object.region)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("sector")
		stream.WriteString(object.sector)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
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
	object = ReadServiceCluster(iterator)
	err = iterator.Error
	return
}

// ReadServiceCluster reads a value of the 'service_cluster' type from the given iterator.
func ReadServiceCluster(iterator *jsoniter.Iterator) *ServiceCluster {
	object := &ServiceCluster{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ServiceClusterLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "dns":
			value := ReadDNS(iterator)
			object.dns = value
			object.fieldSet_[3] = true
		case "cloud_provider":
			value := iterator.ReadString()
			object.cloudProvider = value
			object.fieldSet_[4] = true
		case "cluster_management_reference":
			value := ReadClusterManagementReference(iterator)
			object.clusterManagementReference = value
			object.fieldSet_[5] = true
		case "labels":
			value := ReadLabelList(iterator)
			object.labels = value
			object.fieldSet_[6] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[7] = true
		case "provision_shard_reference":
			value := ReadProvisionShardReference(iterator)
			object.provisionShardReference = value
			object.fieldSet_[8] = true
		case "region":
			value := iterator.ReadString()
			object.region = value
			object.fieldSet_[9] = true
		case "sector":
			value := iterator.ReadString()
			object.sector = value
			object.fieldSet_[10] = true
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
