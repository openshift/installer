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

// MarshalRegistrySources writes a value of the 'registry_sources' type to the given writer.
func MarshalRegistrySources(object *RegistrySources, writer io.Writer) error {
	stream := helpers.NewStream(writer)
	WriteRegistrySources(object, stream)
	err := stream.Flush()
	if err != nil {
		return err
	}
	return stream.Error
}

// WriteRegistrySources writes a value of the 'registry_sources' type to the given stream.
func WriteRegistrySources(object *RegistrySources, stream *jsoniter.Stream) {
	count := 0
	stream.WriteObjectStart()
	var present_ bool
	present_ = object.bitmap_&1 != 0 && object.allowedRegistries != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("allowed_registries")
		WriteStringList(object.allowedRegistries, stream)
		count++
	}
	present_ = object.bitmap_&2 != 0 && object.blockedRegistries != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("blocked_registries")
		WriteStringList(object.blockedRegistries, stream)
		count++
	}
	present_ = object.bitmap_&4 != 0 && object.insecureRegistries != nil
	if present_ {
		if count > 0 {
			stream.WriteMore()
		}
		stream.WriteObjectField("insecure_registries")
		WriteStringList(object.insecureRegistries, stream)
	}
	stream.WriteObjectEnd()
}

// UnmarshalRegistrySources reads a value of the 'registry_sources' type from the given
// source, which can be an slice of bytes, a string or a reader.
func UnmarshalRegistrySources(source interface{}) (object *RegistrySources, err error) {
	iterator, err := helpers.NewIterator(source)
	if err != nil {
		return
	}
	object = ReadRegistrySources(iterator)
	err = iterator.Error
	return
}

// ReadRegistrySources reads a value of the 'registry_sources' type from the given iterator.
func ReadRegistrySources(iterator *jsoniter.Iterator) *RegistrySources {
	object := &RegistrySources{}
	for {
		field := iterator.ReadObject()
		if field == "" {
			break
		}
		switch field {
		case "allowed_registries":
			value := ReadStringList(iterator)
			object.allowedRegistries = value
			object.bitmap_ |= 1
		case "blocked_registries":
			value := ReadStringList(iterator)
			object.blockedRegistries = value
			object.bitmap_ |= 2
		case "insecure_registries":
			value := ReadStringList(iterator)
			object.insecureRegistries = value
			object.bitmap_ |= 4
		default:
			iterator.ReadAny()
		}
	}
	return object
}
