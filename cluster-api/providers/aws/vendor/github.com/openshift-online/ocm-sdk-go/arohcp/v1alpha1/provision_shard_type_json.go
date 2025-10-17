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

package v1alpha1 // github.com/openshift-online/ocm-sdk-go/arohcp/v1alpha1

import (
	"io"
	"time"

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/openshift-online/ocm-sdk-go/clustersmgmt/v1"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalProvisionShard writes a value of the 'provision_shard' type to the given writer.
func MarshalProvisionShard(object *ProvisionShard, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteProvisionShard(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteProvisionShard writes a value of the 'provision_shard' type to the given stream.
func WriteProvisionShard(object *ProvisionShard, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(ProvisionShardLinkKind)
	} else {
		stream.WriteString(ProvisionShardKind)
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
	present_ = object.bitmap_&8 != 0 && object.awsAccountOperatorConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_account_operator_config")
		WriteServerConfig(object.awsAccountOperatorConfig, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_base_domain")
		stream.WriteString(object.awsBaseDomain)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_base_domain")
		stream.WriteString(object.gcpBaseDomain)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.gcpProjectOperator != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_project_operator")
		WriteServerConfig(object.gcpProjectOperator, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		v1.WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.hiveConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hive_config")
		WriteServerConfig(object.hiveConfig, stream)
		count++
	}
	present_ = object.bitmap_&1024 != 0 && object.hypershiftConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hypershift_config")
		WriteServerConfig(object.hypershiftConfig, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_cluster")
		stream.WriteString(object.managementCluster)
		count++
	}
	present_ = object.bitmap_&8192 != 0 && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		v1.WriteCloudRegion(object.region, stream)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		stream.WriteString(object.status)
	}
	stream.WriteObjectEnd()
}

// UnmarshalProvisionShard reads a value of the 'provision_shard' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalProvisionShard(source interface{}) (object *ProvisionShard, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadProvisionShard(iterator)
	err = iterator.Error
	return
}

// ReadProvisionShard reads a value of the 'provision_shard' type from the given iterator.
func ReadProvisionShard(iterator *jsoniter.Iterator) *ProvisionShard {
	object := &ProvisionShard{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ProvisionShardLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "aws_account_operator_config":
			value := ReadServerConfig(iterator)
			object.awsAccountOperatorConfig = value
			object.bitmap_ |= 8
		case "aws_base_domain":
			value := iterator.ReadString()
			object.awsBaseDomain = value
			object.bitmap_ |= 16
		case "gcp_base_domain":
			value := iterator.ReadString()
			object.gcpBaseDomain = value
			object.bitmap_ |= 32
		case "gcp_project_operator":
			value := ReadServerConfig(iterator)
			object.gcpProjectOperator = value
			object.bitmap_ |= 64
		case "cloud_provider":
			value := v1.ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.bitmap_ |= 128
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 256
		case "hive_config":
			value := ReadServerConfig(iterator)
			object.hiveConfig = value
			object.bitmap_ |= 512
		case "hypershift_config":
			value := ReadServerConfig(iterator)
			object.hypershiftConfig = value
			object.bitmap_ |= 1024
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.bitmap_ |= 2048
		case "management_cluster":
			value := iterator.ReadString()
			object.managementCluster = value
			object.bitmap_ |= 4096
		case "region":
			value := v1.ReadCloudRegion(iterator)
			object.region = value
			object.bitmap_ |= 8192
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.bitmap_ |= 16384
		default:
			iterator.ReadAny()
		}
	}
	return object
}
