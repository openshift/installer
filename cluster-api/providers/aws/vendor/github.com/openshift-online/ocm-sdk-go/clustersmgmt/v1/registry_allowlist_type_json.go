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

// MarshalRegistryAllowlist writes a value of the 'registry_allowlist' type to the given writer.
func MarshalRegistryAllowlist(object *RegistryAllowlist, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	writeRegistryAllowlist(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// writeRegistryAllowlist writes a value of the 'registry_allowlist' type to the given stream.
func writeRegistryAllowlist(object *RegistryAllowlist, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	stream.WriteObjectField("kind")
	if object.bitmap_&1 != 0 {
		stream.WriteString(RegistryAllowlistLinkKind)
	} else {
		stream.WriteString(RegistryAllowlistKind)
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
	present_ = object.bitmap_&8 != 0 && object.cloudProvider != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("cloud_provider")
		writeCloudProvider(object.cloudProvider, stream)
		count++
	}
	present_ = object.bitmap_&16 != 0
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("creation_timestamp")
		stream.WriteString((object.creationTimestamp).Format(time.RFC3339))
		count++
	}
	present_ = object.bitmap_&32 != 0 && object.registries != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("registries")
		writeStringList(object.registries, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRegistryAllowlist reads a value of the 'registry_allowlist' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRegistryAllowlist(source interface{}) (object *RegistryAllowlist, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = readRegistryAllowlist(iterator)
	err = iterator.Error
	return
}

// readRegistryAllowlist reads a value of the 'registry_allowlist' type from the given iterator.
func readRegistryAllowlist(iterator *jsoniter.Iterator) *RegistryAllowlist {
	object := &RegistryAllowlist{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "kind":
			value := iterator.ReadString()
			if value == RegistryAllowlistLinkKind {
				object.bitmap_ |= 1
			}
		case "id":
			object.id = iterator.ReadString()
			object.bitmap_ |= 2
		case "href":
			object.href = iterator.ReadString()
			object.bitmap_ |= 4
		case "cloud_provider":
			value := readCloudProvider(iterator)
			object.cloudProvider = value
			object.bitmap_ |= 8
		case "creation_timestamp":
			text := iterator.ReadString()
			value, err := time.Parse(time.RFC3339, text)
			if err != nil {
				iterator.ReportError("", err.Error())
			}
			object.creationTimestamp = value
			object.bitmap_ |= 16
		case "registries":
			value := readStringList(iterator)
			object.registries = value
			object.bitmap_ |= 32
		default:
			iterator.ReadAny()
		}
	}
	return object
}
