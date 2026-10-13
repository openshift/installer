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

// MarshalMachinePool writes a value of the 'machine_pool' type to the given writer.
func MarshalMachinePool(object *MachinePool, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteMachinePool(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteMachinePool writes a value of the 'machine_pool' type to the given stream.
func WriteMachinePool(object *MachinePool, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(MachinePoolLinkKind)
	} else {
		stream.WriteString(MachinePoolKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3] && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAWSMachinePool(object.aws, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4] && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCPMachinePool(object.gcp, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.autoscaling != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("autoscaling")
		WriteMachinePoolAutoscaling(object.autoscaling, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.availabilityZones != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zones")
		WriteStringList(object.availabilityZones, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("instance_type")
		stream.WriteString(object.instanceType)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.labels != nil
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
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("replicas")
		stream.WriteInt(object.replicas)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10] && object.rootVolume != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("root_volume")
		WriteRootVolume(object.rootVolume, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11] && object.securityGroupFilters != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("security_group_filters")
		WriteMachinePoolSecurityGroupFilterList(object.securityGroupFilters, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.subnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnets")
		WriteStringList(object.subnets, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13] && object.taints != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("taints")
		WriteTaintList(object.taints, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalMachinePool reads a value of the 'machine_pool' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalMachinePool(source interface{}) (object *MachinePool, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadMachinePool(iterator)
	err = iterator.Error
	return
}

// ReadMachinePool reads a value of the 'machine_pool' type from the given iterator.
func ReadMachinePool(iterator *jsoniter.Iterator) *MachinePool {
	object := &MachinePool{
		fieldSet_: make([]bool, 14),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == MachinePoolLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "aws":
			value := ReadAWSMachinePool(iterator)
			object.aws = value
			object.fieldSet_[3] = true
		case "gcp":
			value := ReadGCPMachinePool(iterator)
			object.gcp = value
			object.fieldSet_[4] = true
		case "autoscaling":
			value := ReadMachinePoolAutoscaling(iterator)
			object.autoscaling = value
			object.fieldSet_[5] = true
		case "availability_zones":
			value := ReadStringList(iterator)
			object.availabilityZones = value
			object.fieldSet_[6] = true
		case "instance_type":
			value := iterator.ReadString()
			object.instanceType = value
			object.fieldSet_[7] = true
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
			object.fieldSet_[8] = true
		case "replicas":
			value := iterator.ReadInt()
			object.replicas = value
			object.fieldSet_[9] = true
		case "root_volume":
			value := ReadRootVolume(iterator)
			object.rootVolume = value
			object.fieldSet_[10] = true
		case "security_group_filters":
			value := ReadMachinePoolSecurityGroupFilterList(iterator)
			object.securityGroupFilters = value
			object.fieldSet_[11] = true
		case "subnets":
			value := ReadStringList(iterator)
			object.subnets = value
			object.fieldSet_[12] = true
		case "taints":
			value := ReadTaintList(iterator)
			object.taints = value
			object.fieldSet_[13] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
