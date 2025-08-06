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

	jsoniter "github.com/json-iterator/go"
	v1 "github.com/openshift-online/ocm-api-model/clientapi/clustersmgmt/v1"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalCloudRegion writes a value of the 'cloud_region' type to the given writer.
func MarshalCloudRegion(object *CloudRegion, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteCloudRegion(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteCloudRegion writes a value of the 'cloud_region' type to the given stream.
func WriteCloudRegion(object *CloudRegion, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(CloudRegionLinkKind)
	} else {
		stream.WriteString(CloudRegionKind)
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
	present_ = len(object.fieldSet_) > 3 && object.fieldSet_[3]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs_only")
		stream.WriteBool(object.ccsOnly)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_location_id")
		stream.WriteString(object.kmsLocationID)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_location_name")
		stream.WriteString(object.kmsLocationName)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6] && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		v1.WriteCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = len(object.fieldSet_) > 8 && object.fieldSet_[8]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = len(object.fieldSet_) > 9 && object.fieldSet_[9]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("govcloud")
		stream.WriteBool(object.govCloud)
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("supports_hypershift")
		stream.WriteBool(object.supportsHypershift)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("supports_multi_az")
		stream.WriteBool(object.supportsMultiAZ)
	}
	stream.WriteObjectEnd()
}

// UnmarshalCloudRegion reads a value of the 'cloud_region' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalCloudRegion(source interface{}) (object *CloudRegion, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadCloudRegion(iterator)
	err = iterator.Error
	return
}

// ReadCloudRegion reads a value of the 'cloud_region' type from the given iterator.
func ReadCloudRegion(iterator *jsoniter.Iterator) *CloudRegion {
	object := &CloudRegion{
		fieldSet_: make([]bool, 13),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == CloudRegionLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "ccs_only":
			value := iterator.ReadBool()
			object.ccsOnly = value
			object.fieldSet_[3] = true
		case "kms_location_id":
			value := iterator.ReadString()
			object.kmsLocationID = value
			object.fieldSet_[4] = true
		case "kms_location_name":
			value := iterator.ReadString()
			object.kmsLocationName = value
			object.fieldSet_[5] = true
		case "cloud_provider":
			value := v1.ReadCloudProvider(iterator)
			object.cloudProvider = value
			object.fieldSet_[6] = true
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.fieldSet_[7] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[8] = true
		case "govcloud":
			value := iterator.ReadBool()
			object.govCloud = value
			object.fieldSet_[9] = true
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.fieldSet_[10] = true
		case "supports_hypershift":
			value := iterator.ReadBool()
			object.supportsHypershift = value
			object.fieldSet_[11] = true
		case "supports_multi_az":
			value := iterator.ReadBool()
			object.supportsMultiAZ = value
			object.fieldSet_[12] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
