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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
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
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(ProvisionShardLinkKind)
	} else {
		stream.WriteString(ProvisionShardKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.awsAccountOperatorConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_account_operator_config")
		WriteServerConfig(object.awsAccountOperatorConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_base_domain")
		stream.WriteString(object.awsBaseDomain)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_base_domain")
		stream.WriteString(object.gcpBaseDomain)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.gcpProjectOperator != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp_project_operator")
		WriteServerConfig(object.gcpProjectOperator, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.hiveConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hive_config")
		WriteServerConfig(object.hiveConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.hypershiftConfig != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hypershift_config")
		WriteServerConfig(object.hypershiftConfig, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_update_timestamp")
		stream.WriteString((object.lastUpdateTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_cluster")
		stream.WriteString(object.managementCluster)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		WriteCloudRegion(object.region, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
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
	object := &ProvisionShard{
		fieldSet_: make([]bool, 15),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == ProvisionShardLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "aws_account_operator_config":
			value := ReadServerConfig(iterator)
			object.awsAccountOperatorConfig = value
			object.fieldSet_[3] = true
		case "aws_base_domain":
			value := iterator.ReadString()
			object.awsBaseDomain = value
			object.fieldSet_[4] = true
		case "gcp_base_domain":
			value := iterator.ReadString()
			object.gcpBaseDomain = value
			object.fieldSet_[5] = true
		case "gcp_project_operator":
			value := ReadServerConfig(iterator)
			object.gcpProjectOperator = value
			object.fieldSet_[6] = true
		case "cloud_provider":
			value := ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.fieldSet_[7] = true
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.fieldSet_[8] = true
		case "hive_config":
			value := ReadServerConfig(iterator)
			object.hiveConfig = value
			object.fieldSet_[9] = true
		case "hypershift_config":
			value := ReadServerConfig(iterator)
			object.hypershiftConfig = value
			object.fieldSet_[10] = true
		case "last_update_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpdateTimestamp = value
			object.fieldSet_[11] = true
		case "management_cluster":
			value := iterator.ReadString()
			object.managementCluster = value
			object.fieldSet_[12] = true
		case "region":
			value := ReadCloudRegion(iterator)
			object.region = value
			object.fieldSet_[13] = true
		case "status":
			value := iterator.ReadString()
			object.status = value
			object.fieldSet_[14] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
