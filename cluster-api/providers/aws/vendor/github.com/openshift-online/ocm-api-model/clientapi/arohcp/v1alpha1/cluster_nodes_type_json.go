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

// MarshalClusterNodes writes a value of the 'cluster_nodes' type to the given writer.
func MarshalClusterNodes(object *ClusterNodes, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteClusterNodes(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteClusterNodes writes a value of the 'cluster_nodes' type to the given stream.
func WriteClusterNodes(object *ClusterNodes, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.autoscaleCompute != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscale_compute")
		WriteMachinePoolAutoscaling(object.autoscaleCompute, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.availabilityZones != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zones")
		WriteStringList(object.availabilityZones, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute")
		stream.WriteInt(object.compute)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.computeLabels != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_labels")
		if object.computeLabels != nil {
			stream.WriteObjectStart()
			keys := make([]string, len(object.computeLabels))
			i := 0
			for key := range object.computeLabels {
				keys[i] = key
				i++
			}
			sort.Strings(keys)
			for i, key := range keys {
				if i > 0 {
					stream.WriteMore()
				}
				item := object.computeLabels[key]
				stream.WriteObjectField(key)
				stream.WriteString(item)
			}
			stream.WriteObjectEnd()
		} else {
			stream.WriteNil()
		}
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.computeMachineType != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_machine_type")
		WriteMachineType(object.computeMachineType, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.computeRootVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("compute_root_volume")
		WriteRootVolume(object.computeRootVolume, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra")
		stream.WriteInt(object.infra)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.infraMachineType != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("infra_machine_type")
		WriteMachineType(object.infraMachineType, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master")
		stream.WriteInt(object.master)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9] && object.masterMachineType != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("master_machine_type")
		WriteMachineType(object.masterMachineType, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.securityGroupFilters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("security_group_filters")
		WriteMachinePoolSecurityGroupFilterList(object.securityGroupFilters, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("total")
		stream.WriteInt(object.total)
	}
	stream.WriteObjectEnd()
}

// UnmarshalClusterNodes reads a value of the 'cluster_nodes' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalClusterNodes(source interface{}) (object *ClusterNodes, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadClusterNodes(iterator)
	err = iterator.Error
	return
}

// ReadClusterNodes reads a value of the 'cluster_nodes' type from the given iterator.
func ReadClusterNodes(iterator *jsoniter.Iterator) *ClusterNodes {
	object := &ClusterNodes{
		fieldSet_: make([]bool, 12),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "autoscale_compute":
			value := ReadMachinePoolAutoscaling(iterator)
			object.autoscaleCompute = value
			object.fieldSet_[0] = true
		case "availability_zones":
			value := ReadStringList(iterator)
			object.availabilityZones = value
			object.fieldSet_[1] = true
		case "compute":
			value := iterator.ReadInt()
			object.compute = value
			object.fieldSet_[2] = true
		case "compute_labels":
			value := map[string]string{}
			for {
				key := iterator.ReadObject()
				if key == "" {
					break
				}
				item := iterator.ReadString()
				value[key] = item
			}
			object.computeLabels = value
			object.fieldSet_[3] = true
		case "compute_machine_type":
			value := ReadMachineType(iterator)
			object.computeMachineType = value
			object.fieldSet_[4] = true
		case "compute_root_volume":
			value := ReadRootVolume(iterator)
			object.computeRootVolume = value
			object.fieldSet_[5] = true
		case "infra":
			value := iterator.ReadInt()
			object.infra = value
			object.fieldSet_[6] = true
		case "infra_machine_type":
			value := ReadMachineType(iterator)
			object.infraMachineType = value
			object.fieldSet_[7] = true
		case "master":
			value := iterator.ReadInt()
			object.master = value
			object.fieldSet_[8] = true
		case "master_machine_type":
			value := ReadMachineType(iterator)
			object.masterMachineType = value
			object.fieldSet_[9] = true
		case "security_group_filters":
			value := ReadMachinePoolSecurityGroupFilterList(iterator)
			object.securityGroupFilters = value
			object.fieldSet_[10] = true
		case "total":
			value := iterator.ReadInt()
			object.total = value
			object.fieldSet_[11] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
