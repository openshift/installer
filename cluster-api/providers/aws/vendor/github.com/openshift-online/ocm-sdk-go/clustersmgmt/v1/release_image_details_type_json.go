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

// MarshalReleaseImageDetails writes a value of the 'release_image_details' type to the given writer.
func MarshalReleaseImageDetails(object *ReleaseImageDetails, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeReleaseImageDetails(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeReleaseImageDetails writes a value of the 'release_image_details' type to the given stream.
func writeReleaseImageDetails(object *ReleaseImageDetails, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.availableUpgrades != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("available_upgrades")
		writeStringList(object.availableUpgrades, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("release_image")
		stream.WriteString(object.releaseImage)
	}
	stream.WriteObjectEnd()
}

// UnmarshalReleaseImageDetails reads a value of the 'release_image_details' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalReleaseImageDetails(source interface{}) (object *ReleaseImageDetails, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readReleaseImageDetails(iterator)
	err = iterator.Error
	return
}

// readReleaseImageDetails reads a value of the 'release_image_details' type from the given iterator.
func readReleaseImageDetails(iterator *jsoniter.Iterator) *ReleaseImageDetails {
	object := &ReleaseImageDetails{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "available_upgrades":
			value := readStringList(iterator)
			object.availableUpgrades = value
			object.bitmap_ |= 1
		case "release_image":
			value := iterator.ReadString()
			object.releaseImage = value
			object.bitmap_ |= 2
		default:
			iterator.ReadAny()
		}
	}
	return object
}
