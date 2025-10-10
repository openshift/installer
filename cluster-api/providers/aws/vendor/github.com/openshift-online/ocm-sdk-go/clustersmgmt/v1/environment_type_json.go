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

// MarshalEnvironment writes a value of the 'environment' type to the given writer.
func MarshalEnvironment(object *Environment, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteEnvironment(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteEnvironment writes a value of the 'environment' type to the given stream.
func WriteEnvironment(object *Environment, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("backplane_url")
		stream.WriteString(object.backplaneURL)
		count++
	}
	present_ = object.bitmap_&2 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_limited_support_check")
		stream.WriteString((object.lastLimitedSupportCheck).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&4 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("last_upgrade_available_check")
		stream.WriteString((object.lastUpgradeAvailableCheck).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&8 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("name")
		stream.WriteString(object.name)
	}
	stream.WriteObjectEnd()
}

// UnmarshalEnvironment reads a value of the 'environment' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalEnvironment(source interface{}) (object *Environment, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadEnvironment(iterator)
	err = iterator.Error
	return
}

// ReadEnvironment reads a value of the 'environment' type from the given iterator.
func ReadEnvironment(iterator *jsoniter.Iterator) *Environment {
	object := &Environment{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "backplane_url":
			value := iterator.ReadString()
			object.backplaneURL = value
			object.bitmap_ |= 1
		case "last_limited_support_check":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastLimitedSupportCheck = value
			object.bitmap_ |= 2
		case "last_upgrade_available_check":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.lastUpgradeAvailableCheck = value
			object.bitmap_ |= 4
		case "name":
			value := iterator.ReadString()
			object.name = value
			object.bitmap_ |= 8
		default:
			iterator.ReadAny()
		}
	}
	return object
}
