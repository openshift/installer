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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalNodePool writes a value of the 'node_pool' type to the given writer.
func MarshalNodePool(object *NodePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteNodePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteNodePool writes a value of the 'node_pool' type to the given stream.
func WriteNodePool(object *NodePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(NodePoolLinkKind)
	} else {
		stream.WriteString(NodePoolKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.awsNodePool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_node_pool")
		WriteAWSNodePool(object.awsNodePool, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auto_repair")
		stream.WriteBool(object.autoRepair)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.autoscaling != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscaling")
		WriteNodePoolAutoscaling(object.autoscaling, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.azureNodePool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("azure_node_pool")
		WriteAzureNodePool(object.azureNodePool, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.kubeletConfigs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubelet_configs")
		WriteStringList(object.kubeletConfigs, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.labels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("labels")
		if object.labels != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.labels))
			i := 0
			for key := range object.labels {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.labels[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.managementUpgrade != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_upgrade")
		WriteNodePoolManagementUpgrade(object.managementUpgrade, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.nodeDrainGracePeriod != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_drain_grace_period")
		WriteValue(object.nodeDrainGracePeriod, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("replicas")
		stream.WriteInt(object.replicas)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		WriteNodePoolStatus(object.status, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet")
		stream.WriteString(object.subnet)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15] && object.taints != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("taints")
		WriteTaintList(object.taints, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16] && object.tuningConfigs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("tuning_configs")
		WriteStringList(object.tuningConfigs, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 17 && object.fieldSet_[17] && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		WriteVersion(object.version, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalNodePool reads a value of the 'node_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalNodePool(source interface{}) (object *NodePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadNodePool(iterator)
	err = iterator.Error
	return
}

// ReadNodePool reads a value of the 'node_pool' type from the given iterator.
func ReadNodePool(iterator *jsoniter.Iterator) *NodePool {
	object := &NodePool{
		fieldSet_: make([]bool, 18),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == NodePoolLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "aws_node_pool":
			value := ReadAWSNodePool(iterator)
			object.awsNodePool = value
			object.fieldSet_[3] = true
		case "auto_repair":
			value := iterator.ReadBool()
			object.autoRepair = value
			object.fieldSet_[4] = true
		case "autoscaling":
			value := ReadNodePoolAutoscaling(iterator)
			object.autoscaling = value
			object.fieldSet_[5] = true
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.fieldSet_[6] = true
		case "azure_node_pool":
			value := ReadAzureNodePool(iterator)
			object.azureNodePool = value
			object.fieldSet_[7] = true
		case "kubelet_configs":
			value := ReadStringList(iterator)
			object.kubeletConfigs = value
			object.fieldSet_[8] = true
		case "labels":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.labels = value
			object.fieldSet_[9] = true
		case "management_upgrade":
			value := ReadNodePoolManagementUpgrade(iterator)
			object.managementUpgrade = value
			object.fieldSet_[10] = true
		case "node_drain_grace_period":
			value := ReadValue(iterator)
			object.nodeDrainGracePeriod = value
			object.fieldSet_[11] = true
		case "replicas":
			value := iterator.ReadInt()
			object.replicas = value
			object.fieldSet_[12] = true
		case "status":
			value := ReadNodePoolStatus(iterator)
			object.status = value
			object.fieldSet_[13] = true
		case "subnet":
			value := iterator.ReadString()
			object.subnet = value
			object.fieldSet_[14] = true
		case "taints":
			value := ReadTaintList(iterator)
			object.taints = value
			object.fieldSet_[15] = true
		case "tuning_configs":
			value := ReadStringList(iterator)
			object.tuningConfigs = value
			object.fieldSet_[16] = true
		case "version":
			value := ReadVersion(iterator)
			object.version = value
			object.fieldSet_[17] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
