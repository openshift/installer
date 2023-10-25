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

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalCloudRegion writes a value of the 'cloud_region' type to the given writer.
func MarshalCloudRegion(object *CloudRegion, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeCloudRegion(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeCloudRegion writes a value of the 'cloud_region' type to the given stream.
func writeCloudRegion(object *CloudRegion, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(CloudRegionLinkKind)
	} else {
		stream.WriteString(CloudRegionKind)
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
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("ccs_only")
		stream.WriteBool(object.ccsOnly)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_location_id")
		stream.WriteString(object.kmsLocationID)
		count++
	}
	present_ = object.bitmap_&32 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("kms_location_name")
		stream.WriteString(object.kmsLocationName)
		count++
	}
	present_ = object.bitmap_&64 != 0 && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		writeCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("display_name")
		stream.WriteString(object.displayName)
		count++
	}
	present_ = object.bitmap_&256 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("enabled")
		stream.WriteBool(object.enabled)
		count++
	}
	present_ = object.bitmap_&512 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("govcloud")
		stream.WriteBool(object.govCloud)
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("supports_hypershift")
		stream.WriteBool(object.supportsHypershift)
		count++
	}
	present_ = object.bitmap_&4096 != 0
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
	object = readCloudRegion(iterator)
	err = iterator.Error
	return
}

// readCloudRegion reads a value of the 'cloud_region' type from the given iterator.
func readCloudRegion(iterator *jsoniter.Iterator) *CloudRegion {
	object := &CloudRegion{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == CloudRegionLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "ccs_only":
			value := iterator.ReadBool()
			object.ccsOnly = value
			object.bitmap_ |= 8
		case "kms_location_id":
			value := iterator.ReadString()
			object.kmsLocationID = value
			object.bitmap_ |= 16
		case "kms_location_name":
			value := iterator.ReadString()
			object.kmsLocationName = value
			object.bitmap_ |= 32
		case "cloud_provider":
			value := readCloudProvider(iterator)
			object.cloudProvider = value
			object.bitmap_ |= 64
		case "display_name":
			value := iterator.ReadString()
			object.displayName = value
			object.bitmap_ |= 128
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 256
		case "govcloud":
			value := iterator.ReadBool()
			object.govCloud = value
			object.bitmap_ |= 512
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 1024
		case "supports_hypershift":
			value := iterator.ReadBool()
			object.supportsHypershift = value
			object.bitmap_ |= 2048
		case "supports_multi_az":
			value := iterator.ReadBool()
			object.supportsMultiAZ = value
			object.bitmap_ |= 4096
		default:
			iterator.ReadAny()
		}
	}
	return object
}
