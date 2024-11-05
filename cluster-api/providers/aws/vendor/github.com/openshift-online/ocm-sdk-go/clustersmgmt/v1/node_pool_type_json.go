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
	"sort"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalNodePool writes a value of the 'node_pool' type to the given writer.
func MarshalNodePool(object *NodePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeNodePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeNodePool writes a value of the 'node_pool' type to the given stream.
func writeNodePool(object *NodePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(NodePoolLinkKind)
	} else {
		stream.WriteString(NodePoolKind)
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
	present_ = object.bitmap_&8 != 0 && object.awsNodePool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws_node_pool")
		writeAWSNodePool(object.awsNodePool, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("auto_repair")
		stream.WriteBool(object.autoRepair)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.autoscaling != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscaling")
		writeNodePoolAutoscaling(object.autoscaling, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zone")
		stream.WriteString(object.availabilityZone)
		count++
	}
	present_ = object.bitmap_&128 != 0 && object.azureNodePool != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("azure_node_pool")
		writeAzureNodePool(object.azureNodePool, stream)
		count++
	}
	present_ = object.bitmap_&256 != 0 && object.kubeletConfigs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kubelet_configs")
		writeStringList(object.kubeletConfigs, stream)
		count++
	}
	present_ = object.bitmap_&512 != 0 && object.labels != nil
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
	present_ = object.bitmap_&1024 != 0 && object.managementUpgrade != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("management_upgrade")
		writeNodePoolManagementUpgrade(object.managementUpgrade, stream)
		count++
	}
	present_ = object.bitmap_&2048 != 0 && object.nodeDrainGracePeriod != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("node_drain_grace_period")
		writeValue(object.nodeDrainGracePeriod, stream)
		count++
	}
	present_ = object.bitmap_&4096 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("replicas")
		stream.WriteInt(object.replicas)
		count++
	}
	present_ = object.bitmap_&8192 != 0 && object.status != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("status")
		writeNodePoolStatus(object.status, stream)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnet")
		stream.WriteString(object.subnet)
		count++
	}
	present_ = object.bitmap_&32768 != 0 && object.taints != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("taints")
		writeTaintList(object.taints, stream)
		count++
	}
	present_ = object.bitmap_&65536 != 0 && object.tuningConfigs != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("tuning_configs")
		writeStringList(object.tuningConfigs, stream)
		count++
	}
	present_ = object.bitmap_&131072 != 0 && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		writeVersion(object.version, stream)
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
	object = readNodePool(iterator)
	err = iterator.Error
	return
}

// readNodePool reads a value of the 'node_pool' type from the given iterator.
func readNodePool(iterator *jsoniter.Iterator) *NodePool {
	object := &NodePool{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == NodePoolLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "aws_node_pool":
			value := readAWSNodePool(iterator)
			object.awsNodePool = value
			object.bitmap_ |= 8
		case "auto_repair":
			value := iterator.ReadBool()
			object.autoRepair = value
			object.bitmap_ |= 16
		case "autoscaling":
			value := readNodePoolAutoscaling(iterator)
			object.autoscaling = value
			object.bitmap_ |= 32
		case "availability_zone":
			value := iterator.ReadString()
			object.availabilityZone = value
			object.bitmap_ |= 64
		case "azure_node_pool":
			value := readAzureNodePool(iterator)
			object.azureNodePool = value
			object.bitmap_ |= 128
		case "kubelet_configs":
			value := readStringList(iterator)
			object.kubeletConfigs = value
			object.bitmap_ |= 256
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
			object.bitmap_ |= 512
		case "management_upgrade":
			value := readNodePoolManagementUpgrade(iterator)
			object.managementUpgrade = value
			object.bitmap_ |= 1024
		case "node_drain_grace_period":
			value := readValue(iterator)
			object.nodeDrainGracePeriod = value
			object.bitmap_ |= 2048
		case "replicas":
			value := iterator.ReadInt()
			object.replicas = value
			object.bitmap_ |= 4096
		case "status":
			value := readNodePoolStatus(iterator)
			object.status = value
			object.bitmap_ |= 8192
		case "subnet":
			value := iterator.ReadString()
			object.subnet = value
			object.bitmap_ |= 16384
		case "taints":
			value := readTaintList(iterator)
			object.taints = value
			object.bitmap_ |= 32768
		case "tuning_configs":
			value := readStringList(iterator)
			object.tuningConfigs = value
			object.bitmap_ |= 65536
		case "version":
			value := readVersion(iterator)
			object.version = value
			object.bitmap_ |= 131072
		default:
			iterator.ReadAny()
		}
	}
	return object
}
