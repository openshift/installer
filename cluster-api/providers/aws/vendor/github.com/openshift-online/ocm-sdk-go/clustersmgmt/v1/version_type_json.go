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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-sdk-go/helpers"
)

// MarshalVersion writes a value of the 'version' type to the given writer.
func MarshalVersion(object *Version, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeVersion(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeVersion writes a value of the 'version' type to the given stream.
func writeVersion(object *Version, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(VersionLinkKind)
	} else {
		stream.WriteString(VersionKind)
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
		stream.WriteObjectField("gcp_marketplace_enabled")
		stream.WriteBool(object.gcpMarketplaceEnabled)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rosa_enabled")
		stream.WriteBool(object.rosaEnabled)
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.availableUpgrades != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("available_upgrades")
		writeStringList(object.availableUpgrades, stream)
		count++
	}
	present_ = object.bitmap_&64 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("channel_group")
		stream.WriteString(object.channelGroup)
		count++
	}
	present_ = object.bitmap_&128 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("default")
		stream.WriteBool(object.default_)
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
		stream.WriteObjectField("end_of_life_timestamp")
		stream.WriteString((object.endOfLifeTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&1024 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hosted_control_plane_default")
		stream.WriteBool(object.hostedControlPlaneDefault)
		count++
	}
	present_ = object.bitmap_&2048 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hosted_control_plane_enabled")
		stream.WriteBool(object.hostedControlPlaneEnabled)
		count++
	}
	present_ = object.bitmap_&4096 != 0 && object.imageOverrides != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("image_overrides")
		writeImageOverrides(object.imageOverrides, stream)
		count++
	}
	present_ = object.bitmap_&8192 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("raw_id")
		stream.WriteString(object.rawID)
		count++
	}
	present_ = object.bitmap_&16384 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("release_image")
		stream.WriteString(object.releaseImage)
		count++
	}
	present_ = object.bitmap_&32768 != 0 && object.releaseImages != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("release_images")
		writeReleaseImages(object.releaseImages, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalVersion reads a value of the 'version' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalVersion(source interface{}) (object *Version, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readVersion(iterator)
	err = iterator.Error
	return
}

// readVersion reads a value of the 'version' type from the given iterator.
func readVersion(iterator *jsoniter.Iterator) *Version {
	object := &Version{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == VersionLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "gcp_marketplace_enabled":
			value := iterator.ReadBool()
			object.gcpMarketplaceEnabled = value
			object.bitmap_ |= 8
		case "rosa_enabled":
			value := iterator.ReadBool()
			object.rosaEnabled = value
			object.bitmap_ |= 16
		case "available_upgrades":
			value := readStringList(iterator)
			object.availableUpgrades = value
			object.bitmap_ |= 32
		case "channel_group":
			value := iterator.ReadString()
			object.channelGroup = value
			object.bitmap_ |= 64
		case "default":
			value := iterator.ReadBool()
			object.default_ = value
			object.bitmap_ |= 128
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.bitmap_ |= 256
		case "end_of_life_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endOfLifeTimestamp = value
			object.bitmap_ |= 512
		case "hosted_control_plane_default":
			value := iterator.ReadBool()
			object.hostedControlPlaneDefault = value
			object.bitmap_ |= 1024
		case "hosted_control_plane_enabled":
			value := iterator.ReadBool()
			object.hostedControlPlaneEnabled = value
			object.bitmap_ |= 2048
		case "image_overrides":
			value := readImageOverrides(iterator)
			object.imageOverrides = value
			object.bitmap_ |= 4096
		case "raw_id":
			value := iterator.ReadString()
			object.rawID = value
			object.bitmap_ |= 8192
		case "release_image":
			value := iterator.ReadString()
			object.releaseImage = value
			object.bitmap_ |= 16384
		case "release_images":
			value := readReleaseImages(iterator)
			object.releaseImages = value
			object.bitmap_ |= 32768
		default:
			iterator.ReadAny()
		}
	}
	return object
}
