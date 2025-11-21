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
	"time"

	jsoniter "github.com/json-iterator/go"
	"github.com/openshift-online/ocm-api-model/clientapi/helpers"
)

// MarshalVersion writes a value of the 'version' type to the given writer.
func MarshalVersion(object *Version, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteVersion(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteVersion writes a value of the 'version' type to the given stream.
func WriteVersion(object *Version, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if len(object.fieldSet_) > 0 && object.fieldSet_[0] {
		stream.WriteString(VersionLinkKind)
	} else {
		stream.WriteString(VersionKind)
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
		stream.WriteObjectField("gcp_marketplace_enabled")
		stream.WriteBool(object.gcpMarketplaceEnabled)
		count++
	}
	present_ = len(object.fieldSet_) > 4 && object.fieldSet_[4]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("rosa_enabled")
		stream.WriteBool(object.rosaEnabled)
		count++
	}
	present_ = len(object.fieldSet_) > 5 && object.fieldSet_[5] && object.availableUpgrades != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("available_upgrades")
		WriteStringList(object.availableUpgrades, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 6 && object.fieldSet_[6]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("channel_group")
		stream.WriteString(object.channelGroup)
		count++
	}
	present_ = len(object.fieldSet_) > 7 && object.fieldSet_[7]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("default")
		stream.WriteBool(object.default_)
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
		stream.WriteObjectField("end_of_life_timestamp")
		stream.WriteString((object.endOfLifeTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = len(object.fieldSet_) > 10 && object.fieldSet_[10]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hosted_control_plane_default")
		stream.WriteBool(object.hostedControlPlaneDefault)
		count++
	}
	present_ = len(object.fieldSet_) > 11 && object.fieldSet_[11]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("hosted_control_plane_enabled")
		stream.WriteBool(object.hostedControlPlaneEnabled)
		count++
	}
	present_ = len(object.fieldSet_) > 12 && object.fieldSet_[12] && object.imageOverrides != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("image_overrides")
		WriteImageOverrides(object.imageOverrides, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 13 && object.fieldSet_[13]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("raw_id")
		stream.WriteString(object.rawID)
		count++
	}
	present_ = len(object.fieldSet_) > 14 && object.fieldSet_[14]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("release_image")
		stream.WriteString(object.releaseImage)
		count++
	}
	present_ = len(object.fieldSet_) > 15 && object.fieldSet_[15] && object.releaseImages != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("release_images")
		WriteReleaseImages(object.releaseImages, stream)
		count++
	}
	present_ = len(object.fieldSet_) > 16 && object.fieldSet_[16]
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("wif_enabled")
		stream.WriteBool(object.wifEnabled)
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
	object = ReadVersion(iterator)
	err = iterator.Error
	return
}

// ReadVersion reads a value of the 'version' type from the given iterator.
func ReadVersion(iterator *jsoniter.Iterator) *Version {
	object := &Version{
		fieldSet_: make([]bool, 17),
	}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == VersionLinkKind {
				object.fieldSet_[0] = true
			}
		case "id":
			object.id = iterator.ReadString()
			object.fieldSet_[1] = true
		case "href":
			object.href = iterator.ReadString()
			object.fieldSet_[2] = true
		case "gcp_marketplace_enabled":
			value := iterator.ReadBool()
			object.gcpMarketplaceEnabled = value
			object.fieldSet_[3] = true
		case "rosa_enabled":
			value := iterator.ReadBool()
			object.rosaEnabled = value
			object.fieldSet_[4] = true
		case "available_upgrades":
			value := ReadStringList(iterator)
			object.availableUpgrades = value
			object.fieldSet_[5] = true
		case "channel_group":
			value := iterator.ReadString()
			object.channelGroup = value
			object.fieldSet_[6] = true
		case "default":
			value := iterator.ReadBool()
			object.default_ = value
			object.fieldSet_[7] = true
		case "enabled":
			value := iterator.ReadBool()
			object.enabled = value
			object.fieldSet_[8] = true
		case "end_of_life_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.endOfLifeTimestamp = value
			object.fieldSet_[9] = true
		case "hosted_control_plane_default":
			value := iterator.ReadBool()
			object.hostedControlPlaneDefault = value
			object.fieldSet_[10] = true
		case "hosted_control_plane_enabled":
			value := iterator.ReadBool()
			object.hostedControlPlaneEnabled = value
			object.fieldSet_[11] = true
		case "image_overrides":
			value := ReadImageOverrides(iterator)
			object.imageOverrides = value
			object.fieldSet_[12] = true
		case "raw_id":
			value := iterator.ReadString()
			object.rawID = value
			object.fieldSet_[13] = true
		case "release_image":
			value := iterator.ReadString()
			object.releaseImage = value
			object.fieldSet_[14] = true
		case "release_images":
			value := ReadReleaseImages(iterator)
			object.releaseImages = value
			object.fieldSet_[15] = true
		case "wif_enabled":
			value := iterator.ReadBool()
			object.wifEnabled = value
			object.fieldSet_[16] = true
		default:
			iterator.ReadAny()
		}
	}
	return object
}
