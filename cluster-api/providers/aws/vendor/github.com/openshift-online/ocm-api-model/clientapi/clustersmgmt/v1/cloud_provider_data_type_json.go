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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalCloudProviderData writes a value of the 'cloud_provider_data' type to the given writer.
func MarshalCloudProviderData(object *CloudProviderData, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCloudProviderData(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCloudProviderData writes a value of the 'cloud_provider_data' type to the given stream.
func WriteCloudProviderData(object *CloudProviderData, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = len(object.fieldSet_) > 0 && object.fieldSet_[0] && object.aws != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("aws")
		WriteAWS(object.aws, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 1 && object.fieldSet_[1] && object.gcp != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("gcp")
		WriteGCP(object.gcp, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 2 && object.fieldSet_[2] && object.availabilityZones != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("availability_zones")
		WriteStringList(object.availabilityZones, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_location")
		stream.WriteString(object.keyLocation)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("key_ring_name")
		stream.WriteString(object.keyRingName)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.region != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("region")
		WriteCloudRegion(object.region, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.subnets != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("subnets")
		WriteStringList(object.subnets, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7] && object.version != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("version")
		WriteVersion(object.version, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8] && object.vpcIds != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("vpc_ids")
		WriteStringList(object.vpcIds, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalCloudProviderData reads a value of the 'cloud_provider_data' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCloudProviderData(source interface{}) (object *CloudProviderData, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadCloudProviderData(iterator)
	err = iterator.Error
	return
}

// ReadCloudProviderData reads a value of the 'cloud_provider_data' type from the given iterator.
func ReadCloudProviderData(iterator *jsoniter.Iterator) *CloudProviderData {
	object := &CloudProviderData{
		fieldSet_: make([]bool, 9),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "aws":
			value := ReadAWS(iterator)
			object.aws = value
			object.fieldSet_[0] = true
		case "gcp":
			value := ReadGCP(iterator)
			object.gcp = value
			object.fieldSet_[1] = true
		case "availability_zones":
			value := ReadStringList(iterator)
			object.availabilityZones = value
			object.fieldSet_[2] = true
		case "key_location":
			value := iterator.ReadString()
			object.keyLocation = value
			object.fieldSet_[3] = true
		case "key_ring_name":
			value := iterator.ReadString()
			object.keyRingName = value
			object.fieldSet_[4] = true
		case "region":
			value := ReadCloudRegion(iterator)
			object.region = value
			object.fieldSet_[5] = true
		case "subnets":
			value := ReadStringList(iterator)
			object.subnets = value
			object.fieldSet_[6] = true
		case "version":
			value := ReadVersion(iterator)
			object.version = value
			object.fieldSet_[7] = true
		case "vpc_ids":
			value := ReadStringList(iterator)
			object.vpcIds = value
			object.fieldSet_[8] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
